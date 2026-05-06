# {{.ProjectName}}

A background worker/service built with Ligo framework.

## Features

- Simple background worker that executes tasks periodically
- Clean architecture with separated concerns
- Graceful shutdown handling

## Running

```bash
go run cmd/runner/main.go
```

## Project Structure

```
.
├── cmd/
│   └── runner/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration
│   ├── module/
│   │   └── runner.go        # Worker module definition
│   ├── usecase/
│   │   └── worker.go        # Worker business logic
│   └── worker/
│       └── controller.go    # Worker lifecycle management
└── go.mod
```

## Customization

Edit `internal/usecase/worker.go` to implement your worker logic.
