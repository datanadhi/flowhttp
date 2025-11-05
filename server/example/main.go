package main

import (
	"fmt"
	"net/http"
	"time"

	flowhttp "github.com/datanadhi/flowhttp/server"
)

/*
FlowHTTP Full Example
------------------------

Demonstrates almost everything FlowHTTP can do:

- Middleware chaining (pre + post)
- Branch-level routing with inheritance and ClearSteps()
- Short-circuit middleware (auth check)
- Context storage (ctx.Set / ctx.Get)
- Dynamic path params and wildcards
- JSON response (ctx.JSON)
- JSON parsing from request body (ctx.BindJSON)
- Header manipulation middleware
- Nested forks and isolated branches
- Graceful shutdown support

--------------------------------------
Run the demo:
  go run example/main.go

--------------------------------------
Routes Overview:
--------------------------------------
1. GET  /                 → Welcome message with context value
2. GET  /api/user/:id     → Dynamic route param
3. GET  /api/users        → JSON array response
4. POST /api/auth/login   → Short-circuit auth + JSON binding
5. GET  /files/*          → Wildcard route demo
6. GET  /admin/plain/ping → ClearSteps() example (no globals)

--------------------------------------
CURL Examples:
--------------------------------------
# 1. Root route (basic JSON)
curl -v http://localhost:8080/

# 2. Dynamic user
curl -v http://localhost:8080/api/user/42

# 3. List users (JSON array)
curl -v http://localhost:8080/api/users

# 4. Unauthorized login (no X-Auth)
curl -v -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"123"}'

# 5. Authorized login (short-circuit passes)
curl -v -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -H "X-Auth: secret" \
     -d '{"username":"alice","password":"wonder"}'

# 6. Wildcard route
curl -v http://localhost:8080/files/images/cat.jpg

# 7. Admin plain route (no global middleware)
curl -v http://localhost:8080/admin/plain/ping
*/

func main() {
	f := flowhttp.NewFlow()

	// Global middleware: logger
	logger := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		start := time.Now()
		fmt.Printf("[REQ] %s %s\n", ctx.Request.Method, ctx.Request.URL.Path)
		next(ctx)
		fmt.Printf("[END] %s took %v\n", ctx.Request.URL.Path, time.Since(start))
	})

	// Adds common header
	addHeader := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		ctx.Response.Header().Set("X-App", "FlowHTTP")
		next(ctx)
	})

	// Sets app name in context
	setApp := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		ctx.Set("appName", "FlowHTTP Full Demo")
		next(ctx)
	})

	// Post-processing middleware
	after := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		next(ctx)
		ctx.Response.Header().Set("X-Processed-At", time.Now().Format(time.RFC3339))
	})

	// Root branch with global middleware
	root := f.Fork("/", []flowhttp.Step{logger, addHeader, setApp, after})

	// -------------------------------
	// 1. Route: GET /
	// -------------------------------
	root.Stream("GET", "/", nil, func(ctx *flowhttp.FlowContext) {
		ctx.JSON(http.StatusOK, map[string]any{
			"message": "Welcome to FlowHTTP!",
			"app":     ctx.Get("appName"),
		})
	})

	// -------------------------------
	// 2. /api branch with version middleware
	// -------------------------------
	api := root.Fork("/api", []flowhttp.Step{
		flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
			ctx.Set("version", "v1")
			next(ctx)
		}),
	})

	api.Stream("GET", "/user/:id", nil, func(ctx *flowhttp.FlowContext) {
		id := ctx.Param("id")
		ctx.JSON(http.StatusOK, map[string]any{
			"user_id": id,
			"version": ctx.Get("version"),
		})
	})

	api.Stream("GET", "/users", nil, func(ctx *flowhttp.FlowContext) {
		users := []map[string]string{
			{"id": "1", "name": "Alice"},
			{"id": "2", "name": "Bob"},
		}
		ctx.JSON(http.StatusOK, map[string]any{
			"count": len(users),
			"users": users,
		})
	})

	// -------------------------------
	// 3. Auth branch with short-circuit middleware
	// -------------------------------
	authCheck := flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
		if ctx.Request.Header.Get("X-Auth") != "secret" {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		next(ctx)
	})

	auth := api.Fork("/auth", nil).ClearSteps()
	auth.Stream("POST", "/login", []flowhttp.Step{authCheck}, func(ctx *flowhttp.FlowContext) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := ctx.BindJSON(&payload); err != nil {
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{
			"status": "ok",
			"user":   payload.Username,
		})
	})

	// -------------------------------
	// 4. Wildcard route
	// -------------------------------
	root.Stream("GET", "/files/*", nil, func(ctx *flowhttp.FlowContext) {
		ctx.JSON(http.StatusOK, map[string]string{"path": ctx.Request.URL.Path})
	})

	// -------------------------------
	// 5. Nested branch with ClearSteps()
	// -------------------------------
	admin := root.Fork("/admin", []flowhttp.Step{
		flowhttp.CreateStep(func(next flowhttp.Sink, ctx *flowhttp.FlowContext) {
			ctx.Set("admin-check", true)
			next(ctx)
		}),
	})
	plain := admin.Fork("/plain", nil).ClearSteps()
	plain.Stream("GET", "/ping", nil, func(ctx *flowhttp.FlowContext) {
		ctx.JSON(http.StatusOK, map[string]any{
			"ping":         "pong",
			"admin-check":  ctx.Get("admin-check"),
			"globalHeader": ctx.Response.Header().Get("X-App"),
		})
	})

	fmt.Println("FlowHTTP server running at http://localhost:8080")
	if err := f.Run(8080); err != nil {
		fmt.Println("Server error:", err)
	}
}
