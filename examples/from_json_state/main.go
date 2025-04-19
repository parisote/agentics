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
	"sync"
	"maps"
)

type Bag[T any] struct {
	m  map[string]T
	mu sync.RWMutex
}

func NewBag[T any]() *Bag[T]        { return &Bag[T]{m: make(map[string]T)} }
func (b *Bag[T]) Get(k string) T    { b.mu.RLock(); v := b.m[k]; b.mu.RUnlock(); return v }
func (b *Bag[T]) Set(k string, v T) { b.mu.Lock(); b.m[k] = v; b.mu.Unlock() }
func (b *Bag[T]) All() map[string]T { b.mu.RLock(); defer b.mu.RUnlock(); return maps.Clone(b.m) }

func Run(bag interface{}) error {
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

	mem := agentics.NewSliceMemory(10)
	mem.Add("user", "Hi, i interesing in buy a new car")
	bag := agentics.NewBag[any]()

	for _, s := range jsonGraph.State {

		var t reflect.Type
		switch s.Type {
		case "string":
			t = reflect.TypeOf("")
		case "int":
			t = reflect.TypeOf(0)
		case "bool":
			t = reflect.TypeOf(false)
		}

		sf := reflect.StructField{
			Name: strings.Title(s.Name),
			Type: t,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, s.Name)),
		}

		dynType := reflect.StructOf([]reflect.StructField{sf})
		v := reflect.New(dynType).Elem()
		bag.Set(s.Name, v)
	}

	graph := agentics.Graph{}
	for _, node := range jsonGraph.Nodes {
		var a *agentics.Agent
		if node.Type == "orchestrator" {
			a = agentics.NewAgent(node.Name, node.Prompt, agentics.WithBranchs(node.Branches))
		} else {
			var prePostFn map[string]func(bag *agentics.Bag[any]) error
			for _, function := range node.Functions {
				if prePostFn == nil {
					prePostFn = make(map[string]func(bag *agentics.Bag[any]) error)
				}
				fnRaw := indent(function.Code, 1)
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

				adaptedFn := func(bag *agentics.Bag[any]) error {
					stateMap := make(map[string]interface{})
					for attName, attValue := range bag.All() {
						stateMap[attName] = attValue
					}

					stateMap["messages"] = mem.ToArrayString()
					return fn(stateMap)
				}

				switch function.Type {
				case "post":
					prePostFn["post"] = adaptedFn
				case "pre":
					prePostFn["pre"] = adaptedFn
				}
			}
			if len(prePostFn) > 0 {
				a = agentics.NewAgent(node.Name, node.Prompt, agentics.WithPreStateFunction(prePostFn["pre"]), agentics.WithPostStateFunction(prePostFn["post"]))
			} else {
				a = agentics.NewAgent(node.Name, node.Prompt)
			}
		}

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
