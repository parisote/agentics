package agentics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/valyala/fasttemplate"
)

type AgentInterface interface {
	Run(ctx context.Context, bag *Bag[any], mem Memory) AgentResponse
}

type Agent struct {
	Name              string
	Client            *ModelClient
	Model             string
	Instructions      string
	Branchs           []string
	Conditional       func(bag *Bag[any]) string
	Tools             []ToolInterface
	OutputGuardrails  []string
	OutputType        string
	PreStateFunction  func(bag *Bag[any]) error
	PostStateFunction func(bag *Bag[any]) error
}

type NextAgent struct {
	Next string `json:"next"`
}
type AgentOption func(*Agent)

type AgentResponse struct {
	Content   string
	Error     error
	NextAgent string
}

func NewAgent(name string, instructions string, options ...AgentOption) *Agent {
	agent := &Agent{
		Name:         name,
		Instructions: instructions,
		Client:       NewModelClient(OpenAI),
		Model:        "gpt-4o-mini",
	}

	for _, option := range options {
		option(agent)
	}

	return agent
}

func WithClient(client ModelClient) AgentOption {
	return func(a *Agent) {
		a.Client = &client
	}
}

func WithBranchs(branchs []string) AgentOption {
	return func(a *Agent) {
		for _, branch := range branchs {
			a.Instructions = a.Instructions +
				"Response with {\"next\": \"" + branch + "\"} and appropiate agent."
		}
		a.Branchs = branchs
	}
}

func WithTools(tools []ToolInterface) AgentOption {
	return func(a *Agent) {
		a.Tools = tools
	}
}

func WithOutputGuardrails(guardrails []string) AgentOption {
	return func(a *Agent) {
		a.OutputGuardrails = guardrails
	}
}

func WithOutputType(outputType string) AgentOption {
	return func(a *Agent) {
		a.OutputType = outputType
	}
}

func WithModel(model string) AgentOption {
	return func(a *Agent) {
		a.Model = model
	}
}

func WithConditional(conditional func(bag *Bag[any]) string) AgentOption {
	return func(a *Agent) {
		a.Conditional = conditional
	}
}

func WithPreStateFunction(stateFunction func(bag *Bag[any]) error) AgentOption {
	return func(a *Agent) {
		a.PreStateFunction = stateFunction
	}
}

func WithPostStateFunction(stateFunction func(bag *Bag[any]) error) AgentOption {
	return func(a *Agent) {
		a.PostStateFunction = stateFunction
	}
}

func (a *Agent) Run(ctx context.Context, bag *Bag[any], mem Memory) AgentResponse {
	fmt.Printf("Running agent: %s\n", a.Name)
	nextAgent := ""

	if a.Conditional != nil {
		return AgentResponse{
			NextAgent: a.Conditional(bag),
		}
	}

	if a.PreStateFunction != nil {
		a.PreStateFunction(bag)
	}

	tpl := fasttemplate.New(a.Instructions, "{{", "}}")
	prompt := tpl.ExecuteString(bag.All())

	response, err := a.Client.provider.Execute(
		ctx,
		prompt,
		mem.All(),
		a.Tools,
	)
	if err != nil {
		fmt.Println("Error executing agent:", err)
		return AgentResponse{
			Content:   "",
			Error:     err,
			NextAgent: "",
		}
	}
	if response.IsToolCall {
		for _, tool := range a.Tools {
			for _, toolCall := range response.ToolCalls {
				if tool.GetName() == toolCall.Name {
					params := make(map[string]interface{})
					if err := json.Unmarshal([]byte(toolCall.Arguments), &params); err != nil {
						fmt.Println("Error unmarshalling tool call arguments:", err)
						return AgentResponse{
							Content:   "",
							Error:     err,
							NextAgent: "",
						}
					}

					// TODO: Refactor?
					for k, v := range params {
						if floatVal, ok := v.(float64); ok && floatVal == float64(int(floatVal)) {
							params[k] = int(floatVal)
						}
					}

					output := tool.Run(ctx, bag, &ToolParams{Params: params})
					return AgentResponse{
						Content:   output.Output,
						Error:     nil,
						NextAgent: "",
					}
				}
			}
		}
	}

	ressult := response.GetContent()
	if strings.Contains(ressult, "next") {
		var nextAgentStruct NextAgent
		if err := json.Unmarshal([]byte(ressult), &nextAgentStruct); err != nil {
			fmt.Println("Error unmarshalling next agent:", err)
			return AgentResponse{
				Content:   "",
				Error:     err,
				NextAgent: "",
			}
		}
		nextAgent = nextAgentStruct.Next
	}

	mem.Add("assistant", response.GetContent())

	if a.PostStateFunction != nil {
		a.PostStateFunction(bag)
	}

	return AgentResponse{
		Content:   response.GetContent(),
		NextAgent: nextAgent,
	}
}
