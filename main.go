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

	agent1 := agentic.NewAgent("english_agent", "You are an agent that will perform a task in English.")

	agent2 := agentic.NewAgent("spanish_agent", "You are an agent that will perform a task in Spanish.")

	orchestrator := agentic.NewAgent("orchestrator",
		"Your job is to decide which agent to use based on the task.",
		agentic.WithBranchs([]string{"english_agent", "spanish_agent"}))
	// agentic.WithConditional(func(state *agentic.State) string {
	// 	if strings.Contains(state.Messages[0], "Hola") {
	// 		return "spanish_agent"
	// 	}
	// 	return "english_agent"
	// }))

	graph := agentic.Graph{
		State: agentic.State{
			Messages: []string{"hello world"},
		},
	}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(orchestrator)
	graph.SetEntrypoint(orchestrator.Name)
	graph.AddRelation("orchestrator", "english_agent")
	graph.AddRelation("orchestrator", "spanish_agent")
	graph.Run(context.Background())
}
