package agentics

import "context"

type Tool struct {
	Name        string
	Description string
	Function    func(state State) string
}

func (t *Tool) Run(ctx context.Context, state State) string {
	return t.Function(state)
}
