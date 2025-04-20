package tools

import (
	"context"

	"github.com/parisote/agentics/agentics"
)

func init() {
	agentics.RegisterTool("divide", divideTool)
	agentics.RegisterTool("multiply", multiplyTool)
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
