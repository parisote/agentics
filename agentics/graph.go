package agentic

import "fmt"

const (
	Entrypoint = "START"
	Exitpoint  = "END"
)

type Graph struct {
	Entrypoint string
	Agents     map[string]Agent
	Relations  [][]string
	State      State
}

func (g *Graph) AddAgent(agent Agent) {
	if g.Agents == nil {
		g.Agents = make(map[string]Agent)
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

func (g *Graph) Run() {
	currentAgent := g.Entrypoint
	visited := make(map[string]bool)
	queue := []string{currentAgent}

	// Procesar la cola hasta que esté vacía
	for len(queue) > 0 {
		// Tomar el siguiente agente de la cola
		currentAgent = queue[0]
		queue = queue[1:]

		// Si ya visitamos este agente, continuar
		if visited[currentAgent] {
			continue
		}

		// Marcar como visitado y ejecutar
		visited[currentAgent] = true
		agent := g.Agents[currentAgent]
		response := agent.Run(&g.State)

		// Si el agente especifica el siguiente agente, priorizarlo
		if response.NextAgent != "" {
			// Añadir el agente especificado al principio de la cola
			queue = append([]string{response.NextAgent}, queue...)
		} else {
			// Si no hay un siguiente agente específico, seguir el flujo normal
			// Agregar todos los agentes conectados a la cola
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
