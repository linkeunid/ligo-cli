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
| `ligo new-module <name>` | `ligo nm` | Scaffold a new Ligo extension module |
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

Every scaffold ships:

- **`CLAUDE.md`** — project conventions: Go + Ligo best practices, do/don't, lint/test toolchain install steps, pre-merge checklist.
- **`.golangci.yml`** — shared v2 linter config (`gci`, `gofumpt`, `tagalign`, `errorlint`, `govet` with `shadow`/`nilness`, `staticcheck`, `revive`, …). The `gci` local-prefix is templated to your module path.
- **`.github/workflows/ci.yml`** — `golangci-lint`, `go test -race`, `govulncheck` on push to `main` and on PRs. Pinned to Node-24 action versions.
- **`.gitignore`** — Go binaries, IDE files, coverage / profile artifacts, `.env`, common OS noise. Skipped under `--no-git`.

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
| `runner` | `run` | `cmd/runner/<name>/main.go` + `internal/<name>/module.go` + `internal/<name>/usecase.go` + `internal/<name>/worker.go` | — |
| `wired` | `wire` | `<pkg>/wired_gen.go` (compile-time DI wiring from a `//go:build wireinject` injector — see [Wired codegen](#wired-codegen)) | — |

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
ligo g wired                    # codegen DI from internal/wired/inject.go
ligo g wired internal/wired     # explicit pkg path (default: internal/wired)
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
| `internal/email/module.go` | `func EmailModule() ligo.Module` |
| `internal/email/usecase.go` | `EmailUseCase` with `Execute()` method |
| `internal/email/worker.go` | `Controller` with hooks (Initialize, Start, Drain, Stop) |

Each runner runs independently. Start multiple workers by running: `ligo work email` & `ligo work process-orders`

### Interactive mode

Running without a name drops into prompts:

```bash
ligo g res
# ? Resource name: product
# ? Register Product module in internal/module/main.go? (Y/n)
```

## `ligo new-module`

```bash
ligo new-module my-extension
ligo nm my-extension
ligo nm --module github.com/example/my-ext
ligo nm my-ext --no-git
```

Scaffolds a new Ligo extension module into `./my-extension/`, following the structure of [ligo-memory](https://github.com/linkeunid/ligo-memory). Generates a complete extension package with:

- `go.mod` with correct module path and ligo dependency
- `README.md` with installation and usage instructions
- Core service file (`<name>.go`) with "Hello World" placeholder
- DI integration (`module.go`) with `Provider()` and `Module()` functions
- Unit tests (`<name>_test.go`, `module_test.go`)

**Interactive mode:**
```bash
ligo nm
# ? Extension name: myext
# ? Go module path: github.com/myext
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--module <path>` | Custom Go module path (default: github.com/<name>) |
| `--no-git` | Skip git initialization |

---

## `ligo serve`

```bash
ligo serve             # go run ./cmd/api/
ligo serve -w          # restart on .go file changes (debounced 500ms)
ligo serve -nw         # watch on, wired auto-regen off
```

`--watch` (`-w`) monitors `internal/` and `cmd/` recursively via `fsnotify`. `--no-wired` (`-n`) skips wired regen. See [docs/features/wired.md](docs/features/wired.md) for the codegen details.

---

## `ligo work`

```bash
ligo work email             # run cmd/runner/email/main.go
ligo work email -w          # auto-reload on .go changes
ligo work email -nw         # watch on, wired auto-regen off
```

Equivalent to `go run cmd/runner/<name>/main.go`. `--watch` (`-w`) watches `cmd/runner/<name>/` and `internal/<name>/`. `--no-wired` (`-n`) skips wired regen.

## `ligo build`

```bash
ligo build             # go build -o bin/app ./cmd/api/
ligo build --out dist/server
```

## Wired codegen

`ligo g wired` (alias `wire`) emits a compile-time wired DI graph from a `//go:build wireinject` injector — zero runtime reflection, errors caught at generation time, no `github.com/linkeunid/ligo` import required.

```bash
ligo g wired              # scans internal/wired/
ligo g wired --dry-run    # preview
```

`ligo serve` / `ligo work` auto-regen `wired_gen.go` before launch (and on each watch restart). Opt out with `--no-wired` / `-n`.

Full guide: **[docs/features/wired.md](docs/features/wired.md)**.
