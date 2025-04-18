package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/goccy/go-json"
	"github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

const tmpl = `
package dyn
import (
	"fmt"
	"reflect"
	"strings"
)

type IntentState struct {
	Att      map[string]interface{}
	Messages []string
}

func (s *IntentState) GetMessages() []string {
	return s.Messages
}

func (s *IntentState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}

func (s *IntentState) AddAtt(attName string, att interface{}) {
	if s.Att == nil {
		s.Att = make(map[string]interface{})
	}
	s.Att[attName] = att
}

func (s *IntentState) GetAtt(attName string) interface{} {
	return s.Att[attName]
}


func Run(state interface{}) error {
	%s
	return nil
}`

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
	Code string `json:"code"`
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

	state := &IntentState{
		Messages: []string{"Hi, i interesing in sell my car"},
	}

	for _, s := range jsonGraph.State {
		sf := reflect.StructField{
			Name: strings.Title(s.Name),
			Type: reflect.TypeOf(""),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, s.Name)),
		}

		dynType := reflect.StructOf([]reflect.StructField{sf})
		v := reflect.New(dynType).Elem()
		state.AddAtt(s.Name, v)
	}

	graph := agentics.Graph{}
	for _, node := range jsonGraph.Nodes {
		var a *agentics.Agent
		if node.Type == "orchestrator" {
			a = agentics.NewAgent(node.Name, node.Prompt, agentics.WithBranchs(node.Branches))
		} else {
			fnRaw := indent(node.Functions[0].Code, 1)
			fullSource := fmt.Sprintf(tmpl, fnRaw)

			i := interp.New(interp.Options{})
			i.Use(stdlib.Symbols)

			_, err := i.Eval(fullSource)
			if err != nil {
				fmt.Printf("yaegi: %v\n\n----- código generado -----\n%s", err, fullSource)
				return
			}

			v, _ := i.Eval("dyn.Run")
			fn := v.Interface().(func(interface{}) error)

			adaptedFn := func(state agentics.State) error {
				intentState := state.(*IntentState)

				stateMap := map[string]interface{}{
					"intent":   intentState.GetAtt("intent"),
					"noIntent": intentState.GetAtt("noIntent"),
					"messages": intentState.GetMessages(),
				}

				return fn(stateMap)
			}

			a = agentics.NewAgent(node.Name, node.Prompt, agentics.WithPostStateFunction(adaptedFn))
		}

		graph.AddAgent(a)
	}

	graph.SetEntrypoint("detect_intent")

	for _, edge := range jsonGraph.Edges {
		graph.AddRelation(edge.Source, edge.Target)
	}

	response := graph.Run(context.Background(), state)
	fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])
	fmt.Println("intent = ", state.GetAtt("intent"))
	fmt.Println("noIntent = ", state.GetAtt("noIntent"))
}

type IntentState struct {
	Att      map[string]interface{}
	Messages []string
}

func (s *IntentState) GetMessages() []string {
	return s.Messages
}

func (s *IntentState) AddMessages(messages []string) {
	s.Messages = append(s.Messages, messages...)
}

func (s *IntentState) AddAtt(attName string, att interface{}) {
	if s.Att == nil {
		s.Att = make(map[string]interface{})
	}
	s.Att[attName] = att
}

func (s *IntentState) GetAtt(attName string) interface{} {
	return s.Att[attName]
}

func indent(body string, tabs int) string {
	pad := strings.Repeat("\t", tabs)
	return pad + strings.ReplaceAll(body, "\n", "\n"+pad)
}

func IsValid(src string) error {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "snippet.go", src, parser.AllErrors)
	if err != nil {
		return err
	}

	conf := types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {},
	}
	_, err = conf.Check("snippet", fset, []*ast.File{file}, nil)
	return err
}
