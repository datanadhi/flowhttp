# FlowHTTP

FlowHTTP is a lightweight, flow-based HTTP framework for Go that lets you structure web servers as composable flows.  
It focuses on clarity, middleware chaining, and flexible route grouping — without adding heavy abstractions.

---

## Overview

FlowHTTP introduces two simple but powerful ideas:

- **Flow** — The core router. It holds route definitions and can be forked to create nested groups.
- **Step** — A middleware unit. Steps can be chained, cleared, or overridden per branch to control request flow.

Together, these allow you to build readable, modular HTTP servers.

---

## Key Features

- Middleware chaining (pre- and post-handlers)
- Branch-level routing with inheritance
- Short-circuiting steps for validations and auth
- Dynamic path parameters (`ctx.Param`)
- Context storage (`ctx.Set` / `ctx.Get`)
- Built-in JSON helpers (`ctx.JSON`, `ctx.BindJSON`)
- Route groups and wildcard paths
- Graceful shutdown support

---

## Installation

```bash
go get github.com/datanadhi/flowhttp
```

## Example

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/datanadhi/flowhttp"
)

func main() {
	f := flowhttp.NewFlow()

	// Middleware: request logger
	logger := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		start := time.Now()
		fmt.Printf("[REQ] %s %s\n", ctx.Request.Method, ctx.Request.URL.Path)
		next(ctx)
		fmt.Printf("[END] %s took %v\n", ctx.Request.URL.Path, time.Since(start))
	})

	root := f.Fork("/", []flowhttp.Step{logger})

	root.Stream("GET", "/", nil, func(ctx *flowhttp.FlowContext) {
		ctx.JSON(http.StatusOK, map[string]string{"message": "Welcome to FlowHTTP"})
	})

	fmt.Println("Server running on :8080")
	f.Run(8080)
}
```

---

## Concepts

### Flow
A **Flow** is the root router of the application.
It defines routes and acts as the entry point for all branches.
Use `NewFlow()` to create one, then fork it to build subroutes.

```go
f := flowhttp.NewFlow()
root := f.Fork("/", nil)
```

A Flow can be started with:
```go
f.Run(8080)
```

---

### Branch

A **Branch** is a segment of a flow, typically representing a grouped set of routes.
Each branch can:
- Inherit its parent’s middleware (Steps)
- Define its own routes and sub-branches
- Optionally clear inherited Steps with .ClearSteps()

```go
api := root.Fork("/api", nil)          // inherits middleware
auth := api.Fork("/auth", nil).ClearSteps() // isolated branch
```

---

### Step

A **Step** is a middleware-like unit.
It runs before and/or after the main handler and can modify the request or short-circuit execution.

```go
authCheck := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
	if ctx.Request.Header.Get("X-Auth") != "secret" {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	next(ctx)
})
```
You can attach Steps to a branch or individual route.

---

### Sink

A **Sink** represents the next function in the chain — similar to a handler in other frameworks.
It’s automatically managed by FlowHTTP and passed into each Step.
When you call `next(ctx)`, control moves forward to the next Step or final route handler.

---

### FlowContext

`FlowContext` carries everything about a request — similar to Go’s `http.Request` but with added helpers.

It includes:
- `Request` — the raw `*http.Request`
- `Response` — the `http.ResponseWriter`
- `Set(key, value)` / `Get(key)` — store and retrieve data within a request’s lifetime
- `Param(name)` — extract path parameters
- `JSON(status, value)` — send JSON responses
- `BindJSON(target)` — decode JSON request bodies

Example:
```go
ctx.Set("user", "Alice")
fmt.Println(ctx.Get("user")) // "Alice"

id := ctx.Param("id")
ctx.JSON(http.StatusOK, map[string]string{"id": id})
```

---

## JSON Utilities

```go
// Send JSON
ctx.JSON(http.StatusOK, map[string]string{"message": "OK"})

// Parse request JSON
var payload struct {
    Name string `json:"name"`
}
ctx.BindJSON(&payload)
```

---

## Full Example

A complete demonstration is available in [example/main.go](/example/main.go)
It includes middleware chaining, branch-level routing, JSON parsing, short-circuiting, and wildcard routes.

Run it:
```go
go run example/main.go
```

---

## License
FlowHTTP is licensed under the [MIT License](/LICENSE)

---

## Contributing
See [CONTRIBUTING.md](/CONTRIBUTING.md) for setup and contribution guidelines.

---
