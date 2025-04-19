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

	bag := agentics.NewBag[any]()
	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "Hello, how are you?")

	graph := agentics.NewGraph(bag, mem)
	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	response := graph.Run(context.Background())

	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
}
