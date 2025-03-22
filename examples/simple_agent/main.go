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

	agent := agentics.NewAgent("agent", "You are Tomas, a helpful assistant.")

	graph := agentics.Graph{
		State: &agentics.InputState{
			Messages: []string{"Hello, how are you?"},
		},
	}
	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	graph.Run(context.Background())
}
