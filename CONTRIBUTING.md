# Contributing to FlowHTTP

Thank you for your interest in contributing!  
FlowHTTP is an open-source Go framework, and contributions of all kinds — code, docs, or ideas — are welcome.

> **Note:**  
> FlowHTTP’s top priorities are **developer experience** and **performance**.  
> Keep every contribution **lightweight, minimal, and fast** — avoid unnecessary abstractions or dependencies.  
> The goal is to keep FlowHTTP simple, predictable, and efficient.

---

## 1. Getting Started

### Prerequisites
- Go 1.23 or later
- Git
- VS Code (recommended)
- Docker (for Dev Container setup)

---

## 2. Local Development Setup

### Option 1: Using Dev Container (Recommended)
The repository includes a full **Dev Container** setup for VS Code.

1. Clone the repository:
   ```bash
   git clone https://github.com/datanadhi/flowhttp.git
   cd flowhttp
   ```
2. Open in VS Code and reopen in Dev Container when prompted.
   VS Code will automatically:
   - Install Go and tools (`gopls`, `golangci-lint`, etc.)
   - Set up linting and formatting (`goimports`, `auto-fix` on save)
   - Forward port `8080`
3. Run the example you want:
   ```bash
   go run server/example/main.go
   ```
   or
   ```bash
   go run client/example/main.go
   ```
4. If you are trying to run server you call access the app at:
   ```bash
   http://localhost:8080
   ```

---

### Option 2: Manual Setup (Without Dev Container)
1. Install Go manually:
   ```bash
   brew install go   # macOS
   sudo apt install golang-go   # Ubuntu/Debian
   ```
2. Clone and initialize:
   ```bash
   git clone https://github.com/datanadhi/flowhttp.git
   cd flowhttp
   go mod tidy
   ```
3. Run the example:
   ```bash
   go run server/example/main.go
   ```
   or
   ```bash
   go run client/example/main.go
   ```

---

## 3. Project Structure
```
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
```

---

## 4. Coding Standards
- Use `goimports` for formatting (`Ctrl + S` auto-formats in Dev Container).
- Lint before committing:
  ```bash
  golangci-lint run
  ```
- Keep function and file names descriptive.
- Prefer composition over inheritance — keep the “flow” readable.
- Add doc comments (`// Comment`) for exported functions.

---

## 5. Testing Changes
You can write and run tests directly inside the container:
```bash
go test ./... -v
```
If your change affects example behavior, ensure the examples run correctly and routes respond as expected.

---

## 6. Submitting a Pull Request
1. Fork the repository and create a new branch:
   ```bash
   git checkout -b feature/my-feature
   ```
2. Make your changes and test locally.
3. Run formatting and lint checks:
   ```bash
   go fmt ./...
   golangci-lint run
   ```
4. Commit using clear messages:
   ```bash
   git commit -m "Add middleware chaining example to docs"
   ```
5. Push and open a Pull Request (PR) on GitHub:
   ```bash
   git push origin feature/my-feature
   ```

---

## 7. Reporting Issues or Ideas
If you encounter bugs or have ideas for improvements:
- Open an issue on GitHub with:
  - Clear title
  - Steps to reproduce (if bug)
  - Expected behavior
  - Actual behavior or logs

---

## 8. License
By contributing, you agree that your contributions will be licensed under the same [MIT License](/LICENSE) as the project.

---

## 9. Quick Reference

| Task | Command |
|------|----------|
| Run server example | `go run server/example/main.go` |
| Run client example | `go run client/example/main.go` |
| Run tests | `go test ./... -v` |
| Format code | `go fmt ./...` |
| Lint code | `golangci-lint run` |
| Run inside container | Open in VS Code → “Reopen in Container” |
| Sync dependencies | `go mod tidy` |
| Check module graph | `go mod graph` |


---

Thanks for helping make FlowHTTP better!