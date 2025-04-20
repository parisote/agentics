package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/parisote/agentics/agentics"
)

func main() {

	file, err := os.Open("with_tool.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	agentics.RegisterTool("divide", divideTool)
	agentics.RegisterTool("multiply", multiplyTool)

	graph := agentics.FromJson(file)
	graph.Mem.Add("user", "Cuanto es 30 / 3?")

	response := graph.Run(context.Background())
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("result = ", response.Bag.Get("result"))

	graph.Mem.Add("user", "Cuanto es 30 * 3?")
	response = graph.Run(context.Background())
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("result = ", response.Bag.Get("result"))
}

func divideTool(ctx context.Context, bag *agentics.Bag[any], input *agentics.ToolParams) interface{} {
	if input.Params["b"].(int) == 0 {
		return "Error: Division by zero"
	}
	result := input.Params["a"].(int) / input.Params["b"].(int)
	bag.Set("result", result)
	return result
}

func multiplyTool(ctx context.Context, bag *agentics.Bag[any], input *agentics.ToolParams) interface{} {
	result := input.Params["a"].(int) * input.Params["b"].(int)
	bag.Set("result", result)
	return result
}
