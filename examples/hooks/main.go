package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/parisote/agentics/agentics"
	_ "github.com/parisote/agentics/examples/hooks/hooks"
)

func main() {

	file, err := os.Open("weather.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	graph := agentics.FromJson(file)
	graph.Mem.Add("user", "Cual es la temperatura en la ciudad de Buenos Aires?")

	response := graph.Run(context.Background())
	fmt.Printf("Response: %s\n", response.Mem.LastN(1)[0].Content)
	fmt.Println("weather_c = ", response.Bag.Get("weather_c"))
}
