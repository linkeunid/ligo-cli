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
| `ligo work <runner-name>` | — | Run a background worker/runner |

## `ligo new`

```bash
ligo new my-app
ligo new my-app --module github.com/acme/my-app
ligo new my-app --no-git
```

Scaffolds a new Ligo project into `./my-app/`, runs `go mod tidy`, and initialises a git repo.

By default generates a minimal project with one Hello endpoint. Use `--full` for the complete boilerplate with users, file upload, and JWT auth. Use `--runner` for a background worker/worker service.

```bash
ligo new my-app                                      # simple boilerplate (default)
ligo new my-app --full                               # full boilerplate
ligo new my-app --runner                             # background worker/runner
ligo new my-app --module github.com/acme/my-app      # custom module path
ligo new my-app --no-git                             # skip git init
ligo new my-app --pre-release                        # use local ../ligo sibling dirs
```

## `ligo generate`

### Schematics

| Schematic | Alias | Default (simple) | `--full` |
|-----------|-------|-----------------|---------|
| `resource` | `res` | entity + dto + usecase + errors.go + repo + uuid.go + controller + presenter + module | — |
| `module` | `mo` | `internal/module/<name>.go` (wires simple usecase+controller) | `internal/module/<name>.go` (wires repo+usecase+controller) |
| `controller` | `co` | `internal/infrastructure/http/controller/<name>.go` (Hello endpoint) | Full controller (CRUD + presenter) |
| `usecase` | `uc` | `internal/usecase/<name>.go` (Hello method) | Full usecase (CRUD) + DTOs + `errors.go` |
| `entity` | `en` | `internal/domain/entity/<name>.go` + `internal/domain/repository/<name>.go` | — |
| `repository` | `rep` | `internal/infrastructure/persistence/memory/<name>_repo.go` + `uuid.go` | — |
| `dto` | `dto` | `internal/usecase/dto/create_<name>.go` + `update_<name>.go` | — |
| `presenter` | `pre` | `internal/infrastructure/http/presenter/<name>.go` | — |
| `runner` | `run` | `cmd/runner/<name>/main.go` + `internal/<name>/module/` + `internal/<name>/usecase/` + `internal/<name>/worker/` | — |

### Examples

```bash
ligo g co product               # simple controller (Hello endpoint)
ligo g co product --full        # full controller with repo/presenter
ligo g co order --full --with-auth  # full controller with AuthGuard
ligo g uc product               # simple use case (Hello method)
ligo g uc product --full        # full use case with DTOs
ligo g mo product               # simple module
ligo g mo product --full        # full module with repo wiring
ligo g res product              # all layers (entity, dto, usecase, repo, controller, presenter, module)
ligo g res product --dry-run    # preview files without writing
ligo g run email                # background worker/runner
ligo g run process-orders       # another runner
```

### Flags

| Flag | Applies to | Description |
|------|------------|-------------|
| `--full` | `resource`, `controller`, `usecase`, `module` | Generate full template with repo/DTOs instead of simple |
| `--with-auth` | `controller` (with `--full`) | Wire `AuthGuard` into the generated controller |
| `--dry-run` | all schematics | Print `CREATE`/`UPDATE` actions without touching the filesystem |

### Generated type names

**`ligo g res product`**:

| File | Type |
|------|------|
| `domain/entity/product.go` | `Product` |
| `domain/repository/product.go` | `ProductRepository` (interface) |
| `usecase/errors.go` | `ErrNotFound`, `ErrValidation`, ... |
| `usecase/product.go` | `ProductUseCase` with CRUD methods |
| `usecase/dto/create_product.go` | `CreateProductInput` |
| `usecase/dto/update_product.go` | `UpdateProductInput` |
| `infrastructure/persistence/memory/uuid.go` | `newUUID()` helper |
| `infrastructure/persistence/memory/product_repo.go` | `ProductRepository` (impl) |
| `infrastructure/http/controller/product.go` | `ProductController` |
| `infrastructure/http/presenter/product.go` | `ProductPresenter` |
| `module/product.go` | `func Product() ligo.Module` |

**`ligo g run email`**:

| File | Type |
|------|------|
| `cmd/runner/email/main.go` | Entry point with `OnStart`/`OnStop` hooks |
| `internal/email/module/email.go` | `func Module() ligo.Module` |
| `internal/email/usecase/email.go` | `EmailUseCase` with `Execute()` method |
| `internal/email/worker/controller.go` | `Controller` with `Start()`/`Stop()` lifecycle |

Each runner runs independently. Start multiple workers by running: `go run cmd/runner/email/main.go` & `go run cmd/runner/process-orders/main.go`

### Interactive mode

Running without a name drops into prompts:

```bash
ligo g res
# ? Resource name: product
# ? Register Product module in internal/module/main.go? (Y/n)
```

## `ligo serve`

```bash
ligo serve             # go run ./cmd/api/
ligo serve --watch     # restart on .go file changes (debounced 500ms)
```

`--watch` monitors `internal/` and `cmd/` recursively via `fsnotify`. `Ctrl+C` cleanly kills the child process.

---

## `ligo work`

```bash
ligo work email             # run cmd/runner/email/main.go
ligo work process-orders    # run cmd/runner/process-orders/main.go
```

Convenience command to run a background worker/runner. Equivalent to `go run cmd/runner/<name>/main.go`.

## `ligo build`

```bash
ligo build             # go build -o bin/app ./cmd/api/
ligo build --out dist/server
```
