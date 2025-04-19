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

	agent1 := agentics.NewAgent("agent1", "You are an agent that will perform a task in English.")

	agent2 := agentics.NewAgent("agent2",
		"You are an agent that will translate the message to Spanish.",
		agentics.WithPostStateFunction(func(bag *agentics.Bag[any]) error {
			v := bag.Get("step").(int)
			bag.Set("step", v+1)
			return nil
		}),
	)

	agent3 := agentics.NewAgent("agent3",
		"You are an agent that will translate the message to Deutsch.",
		agentics.WithPostStateFunction(func(bag *agentics.Bag[any]) error {
			v := bag.Get("step").(int)
			bag.Set("step", v+1)
			return nil
		}),
	)

	graph := agentics.Graph{}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(agent3)
	graph.SetEntrypoint(agent1.Name)
	graph.AddRelation("agent1", "agent2")
	graph.AddRelation("agent2", "agent3")
	bag := agentics.NewBag[any]()
	bag.Set("step", 0)
	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "hello world")
	response := graph.Run(context.Background(), bag, mem)

	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Printf("Steps: %d\n", response.Bag.Get("step"))
}
