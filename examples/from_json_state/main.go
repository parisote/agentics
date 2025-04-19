package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/parisote/agentics/agentics"
)

func main() {

	file, err := os.Open("customer_service.json")
	if err != nil {
		log.Fatalf("no pude abrir el archivo: %v", err)
	}
	defer file.Close()

	agentics.RegisterHook("changeIntent", changeIntent)
	agentics.RegisterHook("fetchAlgo", fetchAlgo)

	graph := agentics.FromJson(file)

	response := graph.Run(context.Background())
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
