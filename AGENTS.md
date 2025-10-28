# Repository Guidelines

## Project Structure & Module Organization
- `main.go` – entrypoint calling Cobra `cmd.Execute()`.
- `cmd/` – CLI commands (`root.go`, `manage.go`, `config.go`, `models.go`, `version.go`, `update.go`).
- `internal/` – packages by domain:
  - `azure/` (Azure CLI discovery, endpoints)
  - `interactive/` (Bubbletea selector + config wizard)
  - `launcher/` (env setup, launches `codex`)
  - `profiles/` (profile storage in `~/.codzure/`)
  - `updater/` (self‑update via GitHub Releases)
- `Makefile` – local build and clean; output in `bin/`.
- `.goreleaser.yml` – release config; tag `vX.Y.Z` to publish.

## Build, Test, and Development Commands
- `make build` – builds to `bin/codzure` using local `.gocache`.
- `make clean` – removes `bin/` and `.gocache`.
- `go build -o bin/codzure .` – direct build alternative.
- Run locally: `./bin/codzure manage config` then `./bin/codzure`.

## Coding Style & Naming Conventions
- Go 1.23+. Run `go fmt ./...` and `go vet ./...` before PRs.
- Packages and files are lower‑case; Cobra commands live in `cmd/*.go`.
- Prefer small, composable functions; avoid persisting secrets—fetch via `az` at runtime.
- Use clear scopes in names, e.g., `internal/azure`, `internal/profiles`.

## Testing Guidelines
- Use Go’s `testing` package; place tests alongside code as `*_test.go`.
- Run `go test ./... -cover` locally.
- Mock Azure interactions (e.g., command runners) instead of requiring `az login` in unit tests.

## Commit & Pull Request Guidelines
- Commits: concise, imperative, optionally scoped, e.g., `cmd: add models list`, `internal/launcher: pass args`.
- PRs: include description, linked issues, CLI output/screenshot for UX changes, and docs updates when applicable.

## Security & Configuration Tips
- Never commit `~/.codzure` contents or secrets.
- Ensure `az` is installed and authenticated (`az login`); keys/endpoint are resolved at launch.

## Agent‑Specific Instructions
- Keep changes minimal and consistent with existing patterns.
- Wire new commands via `cmd/` and implement logic in `internal/`.
- Update `README.md` and `CONFIGURATION.md` when behavior or flags change.
