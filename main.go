package main

import (
	"context"
	"fmt"

	agentic "github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
)

func main() {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	agent1 := agentic.NewAgent("agent1", "You are an agent that will perform a task in English.")

	agent2 := agentic.NewAgent("agent2", "You are an agent that will perform a task in Spanish.")

	orchestrator := agentic.NewAgent("orchestrator",
		"Your job is to decide which agent to use based on the task. "+
			"If the task is in English, response {\"next\": \"agent1\"}. "+
			"If the task is in Spanish, response {\"next\": \"agent2\"}.",
		agentic.WithBranchs([]string{"agent1", "agent2"}))

	graph := agentic.Graph{
		State: agentic.State{
			Messages: []string{"hello world"},
		},
	}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(orchestrator)
	graph.SetEntrypoint(orchestrator.Name)
	graph.AddRelation("orchestrator", "agent1")
	graph.AddRelation("orchestrator", "agent2")
	graph.Run(context.Background())
}
