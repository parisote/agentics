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

	graph := agentics.Graph{}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(orchestrator)
	graph.SetEntrypoint(orchestrator.Name)
	graph.AddRelation("orchestrator", "english_agent")
	graph.AddRelation("orchestrator", "spanish_agent")
	mem := agentics.NewSliceMemory(10)
	bag := agentics.NewBag[any]()
	mem.Add("user", "Hello, how are you?")
	response := graph.Run(context.Background(), bag, mem)
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)

	mem.Add("user", "Hola mundo")
	response = graph.Run(context.Background(), bag, mem)
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
}
