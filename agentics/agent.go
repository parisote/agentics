package agentic

import (
	"context"
	"fmt"
	"strings"
)

type Agent struct {
	Name             string
	Client           *ModelClient
	Model            string
	Instructions     string
	Branchs          []string
	Tools            []string
	OutputGuardrails []string
	OutputType       string
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
		Model:        "gpt-4o",
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

func (a *Agent) Run(ctx context.Context, state *State) AgentResponse {
	nextAgent := ""
	fmt.Println("Running agent", a.Name)

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
	if strings.Contains(ressult, "agent1") {
		nextAgent = "agent1"
	} else if strings.Contains(ressult, "agent2") {
		nextAgent = "agent2"
	}

	state.Messages = append(state.Messages, response.GetContent())

	return AgentResponse{
		Content:   response.GetContent(),
		NextAgent: nextAgent,
	}
}
