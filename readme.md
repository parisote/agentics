# Agentics

Agentics is a Go library for building and orchestrating AI agent systems. It provides a flexible framework for creating agents that can interact with Large Language Models (LLMs) and execute custom tools.

## Features

- Support for multiple LLM providers (OpenAI, Anthropic)
- Agent orchestration through a graph-based system
- Custom tool integration
- State management
- Branching and conditional flows
- Pre and post-processing hooks for state manipulation

## Installation
```bash
go get github.com/agentics/agentics
```

## Quick Start

Here's a simple example of creating an agent:

```go
package main

import (
    "context"
    "fmt"
    "github.com/parisote/agentics/agentics"
)

func main() {
    // Create a new agent
    agent := agentics.NewAgent("agent", "You are a helpful assistant.")

    // Create a graph
    graph := agentics.Graph{}
    graph.AddAgent(agent)
    graph.SetEntrypoint(agent.Name)

    // Run the agent with a message
    response := graph.Run(context.Background(), &agentics.InputState{
        Messages: []string{"Hello, how are you?"},
    })

    fmt.Printf("Response: %s\n", response.State.GetMessages()[len(response.State.GetMessages())-1])
}
```

## Examples

The repository includes several examples demonstrating different features:

### Simple Chat
A console-based chat application:
```go
agent := agentics.NewAgent("agent", "You are a helpful assistant.")
graph := agentics.Graph{}
graph.AddAgent(agent)
graph.SetEntrypoint(agent.Name)
```

### Tools Integration
Example of using custom tools with agents:
```go
toolMultiply := agentics.NewTool("multiply",
    "Use this tool to multiply two numbers.",
    []agentics.DescriptionParams{
        {Name: "a", Type: "integer"},
        {Name: "b", Type: "integer"},
    },
    func(ctx context.Context, state agentics.State, input *agentics.ToolParams) interface{} {
        return input.Params["a"].(int) * input.Params["b"].(int)
    })

agent := agentics.NewAgent("agent", "You are a helpful assistant.",
    agentics.WithTools([]agentics.ToolInterface{toolMultiply}))
```

### Branching Logic
Example of conditional agent routing:
```go
orchestrator := agentics.NewAgent("orchestrator",
    "Your job is to decide which agent to use based on the task.",
    agentics.WithBranchs([]string{"english_agent", "spanish_agent"}))
```

### Advanced Usage
Multi-agent system with custom state management:
```go
agent1 := agentics.NewAgent("agent1", "Task in English")
agent2 := agentics.NewAgent("agent2", "Translate to Spanish")
agent3 := agentics.NewAgent("agent3", "Translate to Deutsch")

graph := agentics.Graph{}
graph.AddAgent(agent1)
graph.AddAgent(agent2)
graph.AddAgent(agent3)
graph.SetEntrypoint(agent1.Name)
graph.AddRelation("agent1", "agent2")
graph.AddRelation("agent2", "agent3")
```

## Environment Setup

Create a `.env` file in your project root:

```env
OPENAI_API_KEY=your_openai_key
ANTHROPIC_API_KEY=your_anthropic_key
```

## API Reference

### Agent Options
- `WithClient(client ModelClient)`: Set custom model client
- `WithTools(tools []ToolInterface)`: Add tools to agent
- `WithBranchs(branchs []string)`: Configure branching logic
- `WithModel(model string)`: Specify LLM model
- `WithOutputGuardrails(guardrails []string)`: Add output validation
- `WithOutputType(outputType string)`: Set expected output format

### Graph Operations
- `AddAgent(agent *Agent)`: Add an agent to the graph
- `AddRelation(from, to string)`: Create agent relationships
- `SetEntrypoint(agent string)`: Set the starting agent
- `Run(ctx context.Context, state State)`: Execute the agent graph

## License

This project is licensed under the MIT License - see the LICENSE file for details.