# Ligo App

> Go web application built with [Ligo](https://github.com/linkeunid/ligo) framework.

## Quick Start

```bash
# Run
go run ./cmd/api/

# Build
go build -o bin/app ./cmd/api/

# Run binary
./bin/app
```

Server runs on `http://localhost:8080`.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Hello World |

## Project Structure

```
cmd/api/              # Entry point
internal/
├── infrastructure/   # External concerns
│   └── http/         # Controllers, middleware
├── module/           # Module wiring
└── config/           # Application configuration
```

## License

MIT
