package agentics

import (
	"context"
	"fmt"
)

type ToolFunc func(ctx context.Context, bag *Bag[any], input *ToolParams) interface{}

type ToolInterface interface {
	GetName() string
	GetDescription() string
	GetParameters() []DescriptionParams
	Run(ctx context.Context, bag *Bag[any], input *ToolParams) *ToolResponse
}

type ToolResponse struct {
	Output string
}

type ToolParams struct {
	Params map[string]interface{}
}

type Tool struct {
	Name        string
	Description string
	Parameters  []DescriptionParams
	Function    ToolFunc
}

type DescriptionParams struct {
	Name string
	Type string
}

var toolRegistry = make(map[string]ToolFunc)

func RegisterTool(name string, fn ToolFunc) {
	toolRegistry[name] = fn
}

func getTool(name string) (ToolFunc, bool) {
	fn, ok := toolRegistry[name]
	return fn, ok
}

func NewTool(name string, description string, parameters []DescriptionParams, function func(ctx context.Context, bag *Bag[any], input *ToolParams) interface{}) ToolInterface {
	return &Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
		Function:    function,
	}
}

func (t *Tool) GetName() string {
	return t.Name
}

func (t *Tool) GetDescription() string {
	return t.Description
}

func (t *Tool) GetParameters() []DescriptionParams {
	return t.Parameters
}

func (t *Tool) Run(ctx context.Context, bag *Bag[any], input *ToolParams) *ToolResponse {
	output := t.Function(ctx, bag, input)
	var outputString string

	switch output := output.(type) {
	case string:
		outputString = output
	case int:
		outputString = fmt.Sprintf("%d", output)
	}

	return &ToolResponse{
		Output: outputString,
	}
}
