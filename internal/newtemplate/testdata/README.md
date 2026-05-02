# Ligo Boilerplate

> Go web application boilerplate using [Ligo](https://github.com/linkeunid/ligo) framework with Clean Architecture + Repository Pattern.

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

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| GET | `/` | API info | - |
| GET | `/health` | Health check | - |
| GET | `/users` | List users | - |
| GET | `/users/:id` | Get user | Bearer token |
| POST | `/users` | Create user | Bearer token |
| PUT | `/users/:id` | Update user | Bearer token |
| DELETE | `/users/:id` | Delete user | Admin only |
| POST | `/files/upload` | Upload file | - |
| GET | `/files/:id` | Download file | - |
| GET | `/files` | List files | - |
| DELETE | `/files/:id` | Delete file | - |

## Example Auth Tokens

```
user:secret   (regular user)
admin:secret  (admin user)
```

```bash
curl -H "Authorization: Bearer user:secret" http://localhost:8080/users/1
```

## Project Structure

```
cmd/api/              # Entry point — wires all layers together
internal/
├── domain/           # Core business (no external deps)
│   ├── entity/       # Business entities
│   ├── repository/   # Repository interfaces
│   └── service/      # Domain service interfaces
├── usecase/          # Application business logic
│   └── dto/          # Data Transfer Objects
├── infrastructure/   # External concerns
│   ├── auth/         # JWT implementation + guards
│   ├── http/         # Controllers, middleware, presenters
│   └── persistence/  # Repository implementations (memory, postgres, ...)
├── module/           # Module wiring (connects layers)
└── config/           # Application configuration
```

## License

MIT
