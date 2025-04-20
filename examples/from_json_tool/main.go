package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/parisote/agentics/agentics"
	_ "github.com/parisote/agentics/examples/from_json_tool/tools"
)

func main() {

	file, err := os.Open("with_tool.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	graph := agentics.FromJson(file)
	graph.Mem.Add("user", "Cuanto es 30 / 3?")

	response := graph.Run(context.Background())
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("result = ", response.Bag.Get("result"))

	graph.Mem.Add("user", "Cuanto es 30 * 3?")
	response = graph.Run(context.Background())
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("result = ", response.Bag.Get("result"))
}
