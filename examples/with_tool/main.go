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

	ctx := context.Background()

	tool := agentics.NewTool("multiply",
		"Use this tool to multiply two numbers. "+
			"The input is a map with two keys, 'a' and 'b', "+
			"which are the integers to multiply.",
		[]agentics.DescriptionParams{
			{Name: "a", Type: "integer"},
			{Name: "b", Type: "integer"},
		},
		func(ctx context.Context, state agentics.State, input *agentics.ToolParams) interface{} {
			result := input.Params["a"].(int) * input.Params["b"].(int)
			return result
		})

	agent := agentics.NewAgent("agent",
		"You are Tomas, a helpful assistant.",
		agentics.WithTools([]agentics.ToolInterface{tool}),
	)

	state := &agentics.InputState{
		Messages: []string{"how many is 3 * 2?"},
	}

	graph := agentics.Graph{}
	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	graph.Run(ctx, state)

	state.AddMessages([]string{"how many is 3 * 10?"})

	graph.Run(ctx, state)

}
