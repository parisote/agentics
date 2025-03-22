package agentics

import (
	"context"
	"fmt"
)

const (
	Entrypoint = "START"
	Exitpoint  = "END"
)

type Graph struct {
	Entrypoint string
	Agents     map[string]AgentInterface
	Relations  [][]string
	State      State
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

func (g *Graph) Run(ctx context.Context) {
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
		response := agent.Run(ctx, g.State)

		if response.NextAgent != "" {
			queue = append([]string{response.NextAgent}, queue...)
		} else {
			for _, relation := range g.Relations {
				if relation[0] == currentAgent && !visited[relation[1]] {
					queue = append(queue, relation[1])
				}
			}
		}
	}

	fmt.Printf("State: %+v\n", g.State)
	fmt.Println("Graph finished")
}
