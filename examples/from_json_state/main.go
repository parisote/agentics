package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
)

type Graph struct {
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

func main() {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	file, err := os.Open("customer_service.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	var jsonGraph Graph

	if err := json.NewDecoder(file).Decode(&jsonGraph); err != nil {
		log.Fatalf("error al decodificar: %v", err)
	}

	agentics.RegisterHook("changeIntent", changeIntent)
	agentics.RegisterHook("fetchAlgo", fetchAlgo)
	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "Hi, i interesing in buy a new car")
	bag := agentics.NewBag[any]()

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

	graph := agentics.Graph{}
	for _, node := range jsonGraph.Nodes {
		var opts []agentics.AgentOption

		for _, fn := range node.Functions {
			switch fn.Type {
			case "pre":
				opts = append(opts, agentics.WithHooks(agentics.PreHook, fn.Name))
			case "post":
				opts = append(opts, agentics.WithHooks(agentics.PostHook, fn.Name))
			}
		}

		if node.Type == "orchestrator" {
			opts = append(opts, agentics.WithBranchs(node.Branches))
		}

		a := agentics.NewAgent(node.Name, node.Prompt, opts...)
		graph.AddAgent(a)
	}

	graph.SetEntrypoint("detect_intent")

	for _, edge := range jsonGraph.Edges {
		graph.AddRelation(edge.Source, edge.Target)
	}

	response := graph.Run(context.Background(), bag, mem)
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("intent = ", response.Bag.Get("intent"))
	fmt.Println("noIntent = ", response.Bag.Get("noIntent"))
	fmt.Println("step = ", response.Bag.Get("step"))
}

func changeIntent(ctx context.Context, bag *agentics.Bag[any], mem agentics.Memory) error {
	lastMessage := mem.LastN(1)[0].Content
	if strings.Contains(lastMessage, "buyer") {
		bag.Set("intent", "buyer")
		bag.Set("noIntent", "seller")
	} else {
		bag.Set("intent", "seller")
		bag.Set("noIntent", "buyer")
	}
	return nil
}

func fetchAlgo(ctx context.Context, bag *agentics.Bag[any], mem agentics.Memory) error {
	bag.Set("step", 10)
	return nil
}
