package agentics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/subosito/gotenv"
)

type JsonGraph struct {
	Entry    string                 `json:"entry"`
	State    []State                `json:"state"`
	Nodes    []Node                 `json:"nodes"`
	Edges    []Edge                 `json:"edges"`
	Metadata map[string]interface{} `json:"metadata"`
}

type State struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Node struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Prompt    string     `json:"prompt"`
	Branches  []string   `json:"branches,omitempty"` // solo existe en el orquestador
	Functions []Function `json:"functions,omitempty"`
}

type Function struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func FromJson(file *os.File) *Graph {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	var jsonGraph JsonGraph

	if err := json.NewDecoder(file).Decode(&jsonGraph); err != nil {
		log.Fatalf("error al decodificar: %v", err)
	}

	mem := NewSliceMemory(10)
	mem.Add("user", "Hi, i interesing in buy a new car")
	bag := NewBag[any]()

	for _, s := range jsonGraph.State {
		switch s.Type {
		case "string":
			bag.Set(s.Name, "")
		case "int":
			bag.Set(s.Name, 0)
		case "bool":
			bag.Set(s.Name, false)
		}
	}

	graph := NewGraph(bag, mem)
	for _, node := range jsonGraph.Nodes {
		var opts []AgentOption

		for _, fn := range node.Functions {
			switch fn.Type {
			case "pre":
				opts = append(opts, WithHooks(PreHook, fn.Name))
			case "post":
				opts = append(opts, WithHooks(PostHook, fn.Name))
			}
		}

		if node.Type == "orchestrator" {
			opts = append(opts, WithBranchs(node.Branches))
		}

		a := NewAgent(node.Name, node.Prompt, opts...)
		graph.AddAgent(a)
	}

	graph.SetEntrypoint("detect_intent")

	for _, edge := range jsonGraph.Edges {
		graph.AddRelation(edge.Source, edge.Target)
	}

	return graph
}
