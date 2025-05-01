# Agentics Configuration Generator

Una interfaz web para crear configuraciones JSON para [Agentics](https://github.com/parisote/agentics) - framework de orquestación de agentes LLM en Go.

## Características

- Interfaz intuitiva para definir configuraciones de Agentics
- Crea variables de estado, nodos (agentes), funciones y conexiones
- Vista previa del JSON en tiempo real
- Descarga el archivo de configuración o copia al portapapeles

## Inicio rápido

1. Instala las dependencias:

```bash
npm install
```

2. Inicia el servidor de desarrollo:

```bash
npm run dev
```

3. Abre [http://localhost:3000](http://localhost:3000) en tu navegador.

## Cómo usar el generador

1. **Variables de estado**: Define las variables que estarán disponibles en el contexto.
2. **Nodos (Agentes)**: Crea agentes con prompts, funciones y ramas.
3. **Conexiones**: Define el flujo entre agentes estableciendo conexiones.
4. **Entrypoint**: Selecciona el nodo de entrada que iniciará el flujo.
5. **JSON**: Copia o descarga el archivo de configuración generado.

## Ejemplo de configuración

El JSON generado tendrá un formato similar a este:

```json
{
  "entry": "detect_intent",
  "state": [
    {
      "name": "intent",
      "type": "string"
    }
  ],
  "nodes": [
    {
      "name": "detect_intent",
      "type": "agent",
      "prompt": "Tu trabajo es detectar la intención del cliente...",
      "functions": [
        {
          "type": "post",
          "name": "changeIntent"
        }
      ],
      "tools": []
    },
    {
      "name": "context_agent",
      "type": "agent",
      "prompt": "Tu trabajo es saludar al cliente..."
    }
  ],
  "edges": [
    {
      "source": "detect_intent",
      "target": "context_agent"
    }
  ],
  "metadata": {}
}
```

## Uso en Agentics

Una vez generado el JSON, puedes utilizarlo con [Agentics](https://github.com/parisote/agentics) con el siguiente código:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "os"

    "github.com/parisote/agentics/agentics"
)

func main() {
    // Leer archivo JSON
    data, err := os.ReadFile("agentics_config.json")
    if err != nil {
        panic(err)
    }

    // Crear grafo desde JSON
    bag := agentics.NewBag[any]()
    mem := agentics.NewSliceMemory(10)
    graph, err := agentics.FromJSON(data, bag, mem)
    if err != nil {
        panic(err)
    }

    // Ejecutar el grafo
    mem.Add("user", "Hola, necesito ayuda")
    response := graph.Run(context.Background())
    
    fmt.Println(response.Mem.LastN(1)[0].Content)
}
```
