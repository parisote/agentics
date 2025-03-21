package main

import agentic "github.com/parisote/agentics/agentics"

func main() {
	agent1 := agentic.Agent{
		Name:         "agent1",
		Instructions: "You are an agent that will perform a task in English.",
	}

	agent2 := agentic.Agent{
		Name:         "agent2",
		Instructions: "You are an agent that will perform a task in Spanish.",
	}

	orchestrator := agentic.Agent{
		Name:         "orchestrator",
		Branchs:      []string{"agent1", "agent2"},
		Instructions: "You are an agent that decides which agent to use based on the task.",
	}

	graph := agentic.Graph{}
	graph.AddAgent(agent1)
	graph.AddAgent(agent2)
	graph.AddAgent(orchestrator)
	graph.SetEntrypoint(orchestrator.Name)
	graph.AddRelation("orchestrator", "agent1")
	graph.AddRelation("orchestrator", "agent2")
	graph.Run()
}
