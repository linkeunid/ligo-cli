# How to Release ligo-cli

## Prerequisites

- `git` and `gh` CLI installed and authenticated
- All changes committed and pushed to `main`
- All other ligo packages released first (ligo-cli embeds their boilerplate templates)

## Release

Run the release script from the project root:

```bash
./scripts/release.sh          # patch bump: v0.1.4 → v0.1.5
./scripts/release.sh minor    # minor bump: v0.1.4 → v0.2.0
./scripts/release.sh major    # major bump: v0.1.4 → v1.0.0
```

The script will:
1. Detect the latest semver tag automatically
2. Increment the version
3. Push the current branch
4. Create and push the annotated tag

## What gets published

ligo-cli is a Go binary published via the module proxy. Once the tag is pushed, users can install with:

```bash
go install github.com/linkeunid/ligo-cli/cmd/ligo@latest
# or pin to a specific version:
go install github.com/linkeunid/ligo-cli/cmd/ligo@v0.1.5
```

The binary version is embedded automatically via `debug.ReadBuildInfo()` — no manual version file to update.

## Verifying the release

```bash
GOPROXY=direct go install github.com/linkeunid/ligo-cli/cmd/ligo@v0.1.5
ligo --version
```

> The module proxy caches `@latest` for a few minutes. Use an explicit version tag to install immediately after release.

## Dependency order

When releasing multiple ligo packages together, release in this order to avoid dependency resolution issues:

1. `ligo`
2. `ligo-memory`, `ligo-validator` (depend on ligo)
3. `ligo-boilerplate` (depends on all above)
4. `ligo-cli` (references boilerplate templates)
