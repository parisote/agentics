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

	toolMultiply := agentics.NewTool("multiply",
		"Use this tool to multiply two numbers. "+
			"The input is a map with two keys, 'a' and 'b', "+
			"which are the integers to multiply.",
		[]agentics.DescriptionParams{
			{Name: "a", Type: "integer"},
			{Name: "b", Type: "integer"},
		},
		func(ctx context.Context, bag *agentics.Bag[any], input *agentics.ToolParams) interface{} {
			result := input.Params["a"].(int) * input.Params["b"].(int)
			return result
		})

	toolDivide := agentics.NewTool("divide",
		"Use this tool to divide two numbers. "+
			"The input is a map with two keys, 'a' and 'b', "+
			"which are the integers to divide.",
		[]agentics.DescriptionParams{
			{Name: "a", Type: "integer"},
			{Name: "b", Type: "integer"},
		},
		func(ctx context.Context, bag *agentics.Bag[any], input *agentics.ToolParams) interface{} {
			if input.Params["b"].(int) == 0 {
				return "Error: Division by zero"
			}
			result := input.Params["a"].(int) / input.Params["b"].(int)
			return result
		})

	agent := agentics.NewAgent("agent",
		"You are Tomas, a helpful assistant.",
		agentics.WithTools([]agentics.ToolInterface{toolMultiply, toolDivide}),
	)

	graph := agentics.Graph{}
	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	bag := agentics.NewBag[any]()
	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "how many is 3 * 2?")
	response := graph.Run(ctx, bag, mem)
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)

	mem2 := agentics.NewSliceMemory(10)
	mem2.Add("user", "how many is 30 / 10?")
	response = graph.Run(ctx, bag, mem2)
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
}
