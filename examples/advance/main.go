package main

import (
	"context"
	"fmt"

	"github.com/parisote/agentics/agentics"
	"github.com/parisote/agentics/examples/advance/ownstate"
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
		agentics.WithPostStateFunction(func(state agentics.State) error {
			newState := state.(*ownstate.OwnState)
			newState.Step++
			return nil
		}),
	)

	agent3 := agentics.NewAgent("agent3",
		"You are an agent that will translate the message to Deutsch.",
		agentics.WithPostStateFunction(func(state agentics.State) error {
			newState := state.(*ownstate.OwnState)
			newState.Step++
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
	graph.Run(context.Background(), &ownstate.OwnState{
		Step:     0,
		Messages: []string{"hello world"},
	})
}
