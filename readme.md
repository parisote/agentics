# Agentics

> Lightweight graph orchestrator for LLM‑powered agents in Go

Agentics lets you compose **agents**, **tools**, **hooks** and **branching graphs** entirely in Go—no reflection or Python runtime required.

---

## Features

* 🚀  **Tiny core** (< 800 LOC) with zero external runtime deps
* 🧩  **Graph orchestration** – create conditional, multi‑agent pipelines
* 👜  **`Bag`** (thread‑safe generics map) for shared variables
* 💬  **`Memory`** (windowed message history) for LLM context
* 🔌  **Hooks in plain Go** – register once, wire from JSON
* 🛠   Tool calling & function‑calling support
* 🗂   Pluggable providers (OpenAI, Anthropic …)

---

## Installation
```bash
go get github.com/parisote/agentics
```
Requires Go 1.21+.

---

## Quick start
```go
package main

import (
    "context"
    "fmt"

    "github.com/parisote/agentics/agentics"
    _ "myapp/hooks"            // registers fetchWeather()
)

func main() {
    bag := agentics.NewBag[any]()       // shared state
    mem := agentics.NewSliceMemory(32)  // chat history (32 msgs)

    // Agents
    greet := agentics.NewAgent("greeter", "Hello {{name}}!")
    weather := agentics.NewAgent(
        "weather",
        "Temperature in BA: {{weather_c}}°.",
        agentics.WithHook(agentics.Pre, "fetchWeather"),
    )

    // Graph wiring
    graph := agentics.NewGraph(bag, mem)
    graph.AddAgent(greet)
    graph.AddAgent(weather)
    graph.SetEntrypoint(greet.Name)
    graph.AddRelation("greeter", "weather")

    bag.Set("name", "Tomas")

    response := graph.Run(context.Background())
    fmt.Printf("Temp: %.1f\n", bag.Get("weather_c"))
}
```

---

## Examples
* **`examples/simple_agent`** – Basic usage with a single agent
* **`examples/simple_chat`** – REPL chat with a single agent
* **`examples/with_tool`** – define and call custom Go tools
* **`examples/branching`** – orchestrator that routes to English / Spanish agents
* **`examples/from_json_state`** – load graph + hooks from a JSON descriptor
* **`examples/advance`** – more advanced usage
* **`examples/hooks`** – Demonstrates the use of custom hooks

### Tools integration
```go
multiply := agentics.NewTool(
    "multiply",
    "Multiply two integers.",
    []agentics.DescriptionParams{{Name: "a", Type: "integer"}, {Name: "b", Type: "integer"}},
    func(ctx context.Context, bag *agentics.Bag[any], p *agentics.ToolParams) interface{} {
        return p.Params["a"].(int) * p.Params["b"].(int)
    },
)
agent := agentics.NewAgent("calc", "Use multiply when needed.", agentics.WithTools([]agentics.ToolInterface{multiply}))
```

### Branching logic
```go
orch := agentics.NewAgent("orchestrator",
    "Decide which agent should answer.",
    agentics.WithBranchs([]string{"english_agent", "spanish_agent"}),
)
```

---

## JSON configuration
You can ship your pipeline as JSON and wire Go hooks by **name**:
```jsonc
{
  "entry": "detect_intent",
  "nodes": [
    {
      "name": "detect_intent",
      "prompt": "Detect buyer/seller intent and set bag.intent",
      "hooks": [{"type": "post", "name": "detectIntent"}]
    },
    {
      "name": "context_agent",
      "prompt": "Say hello in the right tone (bag.intent)."
    }
  ],
  "edges": [
    {"source": "detect_intent", "target": "context_agent"}
  ]
}
```
The loader in `examples/from_json_state` turns this into a live `Graph`.

---

## Writing hooks
```go
package hooks

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/parisote/agentics/agentics"
)

type weather struct{ Current struct{ Temp float64 `json:"temp_c"` } }

func init() { agentics.RegisterHook("fetchWeather", fetchWeather) }

func fetchWeather(ctx context.Context, c *agentics.Context) error {
    r, err := c.HTTP.Get("https://api.weatherapi.com/v1/current.json?q=Buenos+Aires")
    if err != nil { return err }
    defer r.Body.Close()

    var w weather
    if err := json.NewDecoder(r.Body).Decode(&w); err != nil { return err }
    c.Bag.Set("weather_c", w.Current.Temp)
    return nil
}
```

---

## API reference (core)
### Bag
| Method | Description |
|--------|-------------|
| `NewBag[T]()` | Create a new bag.
| `Get(key)` / `Set(key,val)` | Thread‑safe access.
| `All()` | Shallow clone of all entries.

### Memory
| Method | Description |
|--------|-------------|
| `NewSliceMemory(max int)` | Create windowed memory.
| `Add(role, content)` | Append message (auto‑prune).
| `All()` | Return slice of messages.

### Graph
| Method | Description |
|--------|-------------|
| `NewGraph(bag, mem)` | Instantiate graph with Bag and Memory.
| `AddAgent(agent)` | Add agent to graph.
| `SetEntrypoint(name)` | Set entrypoint agent.
| `AddRelation(a,b)` | Connect nodes.
| `Run(ctx)` | Execute flow and return response.

---

## License
MIT © 2025 Tomas Climente