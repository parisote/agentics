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

	graph := agentics.Graph{}

	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "Hello, how are you?")
	response := graph.Run(context.Background(), agentics.NewBag[any](), mem)

	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
}
