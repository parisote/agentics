package agentic

import (
	"fmt"
	"math/rand"
)

type Agent struct {
	Name             string
	Model            string
	Instructions     string
	Branchs          []string
	Tools            []string
	OutputGuardrails []string
	OutputType       string
}

type AgentResponse struct {
	Content   string
	NextAgent string
}

func (a *Agent) Run(state *State) AgentResponse {
	state.Name = a.Name
	fmt.Printf("State: %+v\n", state)
	fmt.Println("Running agent", a.Name)

	nextAgent := ""
	if a.Branchs != nil {
		nextAgent = a.Branchs[rand.Intn(len(a.Branchs))]
	}

	return AgentResponse{
		Content:   "Hello, world!",
		NextAgent: nextAgent,
	}
}
