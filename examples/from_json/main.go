package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
)

type Graph struct {
	Entry    string                 `json:"entry"`
	State    []State                `json:"state"`
	Nodes    []Node                 `json:"nodes"`
	Edges    []Edge                 `json:"edges"`
	Metadata map[string]interface{} `json:"metadata"` // flexible
}

type State struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Node struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Prompt   string   `json:"prompt"`
	Branches []string `json:"branches,omitempty"` // solo existe en el orquestador
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func main() {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	file, err := os.Open("simple_chat.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	var jsonGraph Graph

	if err := json.NewDecoder(file).Decode(&jsonGraph); err != nil {
		log.Fatalf("error al decodificar: %v", err)
	}

	graph := agentics.Graph{}
	for _, node := range jsonGraph.Nodes {
		var a *agentics.Agent
		if node.Type == "orchestrator" {
			a = agentics.NewAgent(node.Name, node.Prompt, agentics.WithBranchs(node.Branches))
		} else {
			a = agentics.NewAgent(node.Name, node.Prompt)
		}

		graph.AddAgent(a)
	}

	state := &agentics.InputState{
		Messages: []string{"hello world"},
	}

	graph.SetEntrypoint("orchestrator")

	for _, edge := range jsonGraph.Edges {
		graph.AddRelation(edge.Source, edge.Target)
	}

	response := graph.Run(context.Background(), state)
	fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])

	state.AddMessages([]string{"Hola mundo"})
	response = graph.Run(context.Background(), state)
	fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])
}
