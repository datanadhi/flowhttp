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
`Flow` represents a branch of routes.
Each branch can inherit its parent’s middleware (Steps) or define its own.

```go
api := root.Fork("/api", nil) // inherits global steps
auth := api.Fork("/auth", nil).ClearSteps() // no global steps
```

### Step

A `Step` is a middleware-like unit that receives the request context and decides how to continue the chain.

```go
authCheck := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
	if ctx.Request.Header.Get("X-Auth") != "secret" {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	next(ctx)
})
```

**Steps can:**
- Modify headers
- Inject values into ctx
- Short-circuit the request (e.g., for authentication)
- Run logic before or after next(ctx)

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
