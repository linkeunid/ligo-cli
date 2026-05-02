# ligo-cli

[![Go Version](https://img.shields.io/badge/go-1.21+-blue)](https://go.dev/dl)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen)](https://github.com/linkeunid/ligo-cli)

A NestJS-inspired CLI for [Ligo framework](https://github.com/linkeunid/ligo) projects. Scaffold new projects and generate Clean Architecture boilerplate so you stay in the architecture without memorising file locations or wiring patterns.

## Installation

```bash
go install github.com/linkeunid/ligo-cli/cmd/ligo@latest
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `ligo new <name>` | `ligo n` | Scaffold a new Ligo project |
| `ligo generate <schematic> <name>` | `ligo g` | Generate a schematic into an existing project |
| `ligo build` | `ligo b` | Build the application (`go build -o bin/app ./cmd/api/`) |
| `ligo serve` | `ligo s` | Run the application (`go run ./cmd/api/`) |

## `ligo new`

```bash
ligo new my-app
ligo new my-app --module github.com/acme/my-app
ligo new my-app --no-git
```

Writes the full Ligo boilerplate tree into `./my-app/`, runs `go mod tidy`, and initialises a git repo.

## `ligo generate`

### Schematics

| Schematic | Alias | Default (simple) | `--full` |
|-----------|-------|-----------------|---------|
| `resource` | `res` | All simple layers + register prompt | All full layers + register prompt |
| `module` | `mo` | `internal/module/<name>.go` (simple) | `internal/module/<name>.go` (with repo wiring) |
| `controller` | `co` | `internal/infrastructure/http/controller/<name>.go` (simple) | Full controller with repo/presenter |
| `usecase` | `uc` | `internal/usecase/<name>.go` (simple Hello) | Full usecase + DTOs |
| `entity` | `en` | `internal/domain/entity/<name>.go` + repository interface | — |
| `repository` | `rep` | `internal/infrastructure/persistence/memory/<name>_repo.go` | — |
| `dto` | `dto` | `internal/usecase/dto/create_<name>.go` + `update_<name>.go` | — |

### Examples

```bash
ligo g co product               # simple controller (Hello endpoint)
ligo g co product --full        # full controller with repo/presenter
ligo g co order --full --with-auth  # full controller with AuthGuard
ligo g uc product               # simple use case (Hello method)
ligo g uc product --full        # full use case with DTOs
ligo g mo product               # simple module
ligo g mo product --full        # full module with repo wiring
ligo g res product              # all simple layers
ligo g res product --full       # all full layers
ligo g res product --dry-run    # preview files without writing
```

### Flags

| Flag | Applies to | Description |
|------|------------|-------------|
| `--full` | `resource`, `controller`, `usecase`, `module` | Generate full template with repo/DTOs instead of simple |
| `--with-auth` | `controller` (with `--full`) | Wire `AuthGuard` into the generated controller |
| `--dry-run` | all schematics | Print `CREATE`/`UPDATE` actions without touching the filesystem |

### Generated type names

Given `ligo g res product`:

| File | Type |
|------|------|
| `domain/entity/product.go` | `Product` |
| `domain/repository/product.go` | `ProductRepository` (interface) |
| `usecase/product.go` | `ProductUseCase` |
| `usecase/dto/create_product.go` | `CreateProductInput` |
| `infrastructure/http/controller/product.go` | `ProductController` |
| `infrastructure/http/presenter/product.go` | `ProductPresenter` |
| `infrastructure/persistence/memory/product_repo.go` | `ProductRepository` (impl) |
| `module/product.go` | `func Product() ligo.Module` |

### Interactive mode

Running without a name drops into prompts:

```bash
ligo g res
# ? Resource name: product
# ? Generate with auth guard? (y/N)
# ? Register Product module in internal/module/main.go? (Y/n)
```

## `ligo serve`

```bash
ligo serve             # go run ./cmd/api/
ligo serve --watch     # restart on .go file changes (debounced 500ms)
```

`--watch` monitors `internal/` and `cmd/` recursively via `fsnotify`. `Ctrl+C` cleanly kills the child process.

## `ligo build`

```bash
ligo build             # go build -o bin/app ./cmd/api/
ligo build --out dist/server
```
