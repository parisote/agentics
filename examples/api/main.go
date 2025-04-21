package main

import (
	"context"
	"net/http"

	"github.com/parisote/agentics/agentics"
)

func main() {

	bag := agentics.NewBag[any]()
	memory := agentics.NewSliceMemory(10)
	agent := agentics.NewAgent("english_agent", "You are an agent that will perform a task in English.")
	graph := agentics.NewGraph(bag, memory)
	graph.AddAgent(agent)
	graph.SetEntrypoint(agent.Name)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		graph.Mem.Add("user", r.URL.Query().Get("input"))
		response := graph.Run(ctx)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response.Mem.LastN(1)[0].Content))
	})
	http.ListenAndServe(":8080", nil)
}
