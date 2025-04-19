package agentics

import (
	"context"
	"fmt"
)

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
	Function    func(ctx context.Context, bag *Bag[any], input *ToolParams) interface{}
}

type DescriptionParams struct {
	Name string
	Type string
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

	switch output := output.(type) {
	case string:
		return &ToolResponse{
			Output: output,
		}
	case int:
		return &ToolResponse{
			Output: fmt.Sprintf("%d", output),
		}
	}

	return nil
}
