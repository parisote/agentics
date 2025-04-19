package agentics

import (
	"context"
)

const (
	Entrypoint = "START"
	Exitpoint  = "END"
)

type Graph struct {
	Entrypoint string
	Agents     map[string]AgentInterface
	Relations  [][]string
	Bag        *Bag[any]
	Mem        Memory
}

type GraphResponse struct {
	Bag *Bag[any]
	Mem Memory
}

func (g *Graph) AddAgent(agent *Agent) {
	if g.Agents == nil {
		g.Agents = make(map[string]AgentInterface)
	}
	g.Agents[agent.Name] = agent
}

func (g *Graph) AddRelation(from, to string) {
	if g.Relations == nil {
		g.Relations = [][]string{}
	}
	g.Relations = append(g.Relations, []string{from, to})
}

func (g *Graph) SetEntrypoint(agent string) {
	g.Entrypoint = agent
}

func NewGraph(bag *Bag[any], mem Memory) *Graph {
	return &Graph{
		Bag: bag,
		Mem: mem,
	}
}
func (g *Graph) Run(ctx context.Context) *GraphResponse {
	currentAgent := g.Entrypoint
	visited := make(map[string]bool)
	queue := []string{currentAgent}

	for len(queue) > 0 {
		currentAgent = queue[0]
		queue = queue[1:]

		if visited[currentAgent] {
			continue
		}

		visited[currentAgent] = true
		agent := g.Agents[currentAgent]
		response := agent.Run(ctx, g.Bag, g.Mem)

		if response.NextAgent != "" {
			queue = append([]string{response.NextAgent}, queue...)
		} else {
			for _, relation := range g.Relations {
				if relation[0] == currentAgent && !visited[relation[1]] {
					queue = append(queue, relation[1])
				}
			}
		}

		g.Mem.Add("assistant", response.Content)
	}

	return &GraphResponse{
		Bag: g.Bag,
		Mem: g.Mem,
	}
}
