# FlowHTTP

FlowHTTP is a lightweight and composable HTTP framework written in Go.
It now provides both server-side flow handling and client-side HTTP utilities, allowing you to build and consume APIs using a consistent and minimal interface.

---

## Overview

FlowHTTP helps you handle HTTP in Go with two complementary modules:
| Module   | Description |
|-----------|-------------|
| `server/` | A composable HTTP flow framework built around middleware chains, contextual storage, and expressive routing. |
| `client/` | A minimal, reliable HTTP client wrapper for making requests, parsing JSON, and handling responses easily. |

Each module is independent and can be used standalone or together in the same project.

---

## Features

### Server
- Middleware chaining and flow-based routing
- Context-aware request handling (`ctx.Set`, `ctx.Get`)
- Built-in JSON binding and response helpers
- Dynamic routing with parameters and wildcards
- Route grouping with `Fork()` and `ClearSteps()`
- Graceful shutdown support
- Minimal dependencies and clean structure

### Client
- Simple GET and POST helpers
- Built-in JSON, string, and byte parsing from responses
- Configurable timeout support
- Header and query parameter helpers
- Status helpers (`IsSuccess`, `StatusText`)
- Automatic body caching for multiple reads

---

## Repository Structure

```
flowhttp/
├── client/
│   ├── client.go
│   ├── response.go
│   ├── utils.go
│   └── example/
│       └── main.go
│
├── server/
│   ├── flow.go
│   ├── context.go
│   ├── middleware.go
│   ├── routing.go
│   ├── server.go
│   └── example/
│       └── main.go
│
└── README.md
```

---

## Getting Started

### Installation

```bash
go get github.com/datanadhi/flowhttp
```
Import the package you need:
```go
import "github.com/datanadhi/flowhttp/server"
import "github.com/datanadhi/flowhttp/client"
```

---

## Usage Examples

### 1. Server Example
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/datanadhi/flowhttp/server"
)

func main() {
	f := server.NewFlow()

	logger := server.CreateStep(func(next server.Sink, ctx *server.FlowContext) {
		start := time.Now()
		fmt.Printf("[REQ] %s %s\n", ctx.Request.Method, ctx.Request.URL.Path)
		next(ctx)
		fmt.Printf("[END] %s took %v\n", ctx.Request.URL.Path, time.Since(start))
	})

	root := f.Fork("/", []server.Step{logger})

	root.Stream("GET", "/", nil, func(ctx *server.FlowContext) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "Welcome to FlowHTTP"})
	})

	fmt.Println("Server running at http://localhost:8080")
	f.Run(8080)
}
```

Run it with:

```bash
go run <path-to-file>/filename.go
```

---

### 2. Client Example
```go
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/datanadhi/flowhttp/client"
)

func main() {
	c := client.NewClient(5 * time.Second)

	// GET
	resp, err := c.Get("https://httpbin.org/get", map[string]string{"q": "flowhttp"}, nil)
	if err != nil {
		panic(err)
	}
	data, _ := resp.Json()
	fmt.Println("GET:", data["url"])

	// POST
	payload := strings.NewReader(`{"hello": "world"}`)
	resp2, err := c.Post("https://httpbin.org/post", nil, nil, payload, "application/json")
	if err != nil {
		panic(err)
	}
	fmt.Println("POST success:", resp2.IsSuccess())
}
```

Run it with:

```bash
go run <path-to-file>/filename.go
```

---

## Philosophy
FlowHTTP is designed to simplify how developers handle HTTP in Go without adding heavy abstractions.
Both `server` and `client` follow Go’s standard library design, keeping everything explicit, composable, and minimal.

You can choose to:
- Use only the `server` package for routing and flow-based APIs.
- Use only the `client` package for outgoing requests.
- Or combine both for full API integration testing and communication.

---

## License
FlowHTTP is licensed under the [MIT License](/LICENSE)

---

## Contributing
See [CONTRIBUTING.md](/CONTRIBUTING.md) for setup and contribution guidelines.

---
