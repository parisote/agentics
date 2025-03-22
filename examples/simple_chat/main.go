package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/parisote/agentics/agentics"
	"github.com/subosito/gotenv"
)

func main() {
	fmt.Println("¡Bienvenido al chat simple de consola!")
	fmt.Println("Escribe ':q' para terminar la conversación.")
	fmt.Print("Por favor, ingresa tu nombre: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	graph := createGraph()
	state := &agentics.InputState{}

	fmt.Printf("Hola %s! Puedes comenzar a chatear.\n", username)

	for {
		fmt.Print("> ")
		scanner.Scan()
		userInput := scanner.Text()

		userInput = strings.TrimSpace(userInput)

		if strings.ToLower(userInput) == ":q" {
			fmt.Println("¡Hasta luego!")
			break
		}

		fmt.Printf("[%s]: %s\n", username, userInput)

		state.AddMessages([]string{userInput})
		response := graph.Run(context.Background(), state)

		fmt.Printf("[Sistema]: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])
	}
}

func createGraph() agentics.Graph {
	err := gotenv.Load()
	if err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	agent := agentics.NewAgent("agent", "You are Tomas, a helpful assistant.")

	graph := agentics.Graph{}

	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)
	return graph
}
