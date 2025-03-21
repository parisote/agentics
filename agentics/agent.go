package agentic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type Agent struct {
	Name             string
	Client           *ModelClient
	Model            string
	Instructions     string
	Branchs          []string
	Conditional      func(state *State) string
	Tools            []string
	OutputGuardrails []string
	OutputType       string
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

func WithTools(tools []string) AgentOption {
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

func WithConditional(conditional func(state *State) string) AgentOption {
	return func(a *Agent) {
		a.Conditional = conditional
	}
}

func (a *Agent) Run(ctx context.Context, state *State) AgentResponse {
	nextAgent := ""
	fmt.Println("Running agent", a.Name)

	if a.Conditional != nil {
		return AgentResponse{
			NextAgent: a.Conditional(state),
		}
	}

	response, err := a.Client.provider.Execute(
		ctx,
		a.Instructions,
		state.Messages,
	)
	if err != nil {
		return AgentResponse{
			Content:   "",
			Error:     err,
			NextAgent: "",
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

	state.Messages = append(state.Messages, response.GetContent())

	return AgentResponse{
		Content:   response.GetContent(),
		NextAgent: nextAgent,
	}
}
