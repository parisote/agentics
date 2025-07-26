package agentics

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/valyala/fasttemplate"
)

type AgentInterface interface {
	Run(ctx context.Context, bag *Bag[any], mem Memory) AgentResponse
}

type Agent struct {
	Name             string
	Client           *ModelClient
	Model            string
	Instructions     string
	Branchs          []string
	Conditional      func(bag *Bag[any]) string
	Tools            []ToolInterface
	OutputGuardrails []string
	OutputType       string
	hooks            []struct {
		kind Kind
		fn   Func
	}
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
		Model:        "", // Se usará el modelo por defecto del proveedor
	}

	for _, option := range options {
		option(agent)
	}

	return agent
}

func WithClient(client ModelClient) AgentOption {
	return func(a *Agent) {
		a.Client = &client
		// Si ya tenemos un modelo configurado, aplicarlo al nuevo cliente
		if a.Model != "" {
			a.Client.provider.SetModel(a.Model)
		}
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
		if a.Client != nil {
			a.Client.provider.SetModel(model)
		}
	}
}

func WithConditional(conditional func(bag *Bag[any]) string) AgentOption {
	return func(a *Agent) {
		a.Conditional = conditional
	}
}

func WithHooks(kind Kind, name string) AgentOption {
	return func(a *Agent) {
		if fn, ok := getHook(name); ok {
			a.hooks = append(a.hooks, struct {
				kind Kind
				fn   Func
			}{kind, fn})
		} else {
			panic("hook no registrado: " + name)
		}
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

	var c *Context
	if len(a.hooks) > 0 {
		c = &Context{
			Bag:    bag,
			Memory: mem,
			HTTP:   http.DefaultClient,
			DB:     nil,
		}
	}

	for _, h := range a.hooks {
		if h.kind == PreHook {
			h.fn(ctx, c)
		}
	}

	tpl := fasttemplate.New(a.Instructions, "{{", "}}")

	// Convertir valores no string a string
	bagValues := bag.All()
	stringValues := make(map[string]interface{})
	for k, v := range bagValues {
		switch val := v.(type) {
		case string:
			stringValues[k] = val
		case fmt.Stringer:
			stringValues[k] = val.String()
		default:
			stringValues[k] = fmt.Sprintf("%v", val)
		}
	}

	prompt := tpl.ExecuteString(stringValues)

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

					// Crear mensaje con el resultado de la herramienta
					toolResultMessage := fmt.Sprintf("I used the %s tool with the arguments %s and got this result: %s. Please provide a final response based on this information.",
						toolCall.Name, toolCall.Arguments, output.Output)

					// Usar el método ExecuteWithFollowUp del proveedor
					followUpResponse, err := a.Client.provider.ExecuteWithFollowUp(
						ctx,
						prompt,
						mem.All(),
						a.Tools,
						toolResultMessage,
					)
					if err != nil {
						fmt.Println("Error executing tool follow-up:", err)
						return AgentResponse{
							Content:   "",
							Error:     err,
							NextAgent: "",
						}
					}

					return AgentResponse{
						Content:   followUpResponse.GetContent(),
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

	for _, h := range a.hooks {
		if h.kind == PostHook {
			h.fn(ctx, c)
		}
	}

	return AgentResponse{
		Content:   response.GetContent(),
		NextAgent: nextAgent,
	}
}
