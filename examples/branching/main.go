package main

import (
	"context"
	"fmt"

	"github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
)

func main() {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	agent1 := agentics.NewAgent("english_agent", "You are an agent that will perform a task in English.")

	agent2 := agentics.NewAgent("spanish_agent", "You are an agent that will perform a task in Spanish.")

	orchestrator := agentics.NewAgent("orchestrator",
		"Your job is to decide which agent to use based on the task.",
		agentics.WithBranchs([]string{"english_agent", "spanish_agent"}))

	state := &agentics.InputState{
		Messages: []string{"hello world"},
	}

	graph := agentics.Graph{}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(orchestrator)
	graph.SetEntrypoint(orchestrator.Name)
	graph.AddRelation("orchestrator", "english_agent")
	graph.AddRelation("orchestrator", "spanish_agent")
	response := graph.Run(context.Background(), state)
	fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])

	state.AddMessages([]string{"Hola mundo"})
	response = graph.Run(context.Background(), state)
	fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])
}
