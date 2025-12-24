# Repository Guidelines

## Project Structure & Module Organization
- `cmd/server`: main API entrypoint.
- `cmd/seed`: database seeding tool.
- `internal/`: application logic, handlers, and domain code.
- `common/`: shared utilities (e.g., config helpers).
- `migrations/`: Goose SQL migrations.
- `docs/`: generated Swagger docs (from `swag`).
- `sqlc.yaml`: SQLC configuration for repository codegen.
- `bin/`: build output (`bin/server`).

## Build, Test, and Development Commands
- `make build`: compile the server to `bin/server`.
- `make run`: build and run the API locally.
- `make test`: run Go tests (`go test -v ./...`).
- `make seed`: load test data via `cmd/seed`.
- `make swagger`: generate Swagger docs into `docs/`.
- `make sqlc`: regenerate SQLC repository code.
- `make migrate-up|migrate-down|migrate-status`: manage database migrations.
- `make docker-up|docker-down`: start/stop `docker-compose` services.
- `make tools`: install dev tools (`swag`, `goose`, `sqlc`).

## Coding Style & Naming Conventions
- Language: Go. Use `gofmt` formatting and idiomatic Go naming.
- Packages: keep package names short and lowercase.
- Files: use `*_test.go` for tests; keep command binaries in `cmd/<name>`.
- Config: `.env` is loaded by the Makefile; `source-env.sh` can export it for shells.

## Testing Guidelines
- Tests are Go tests (`go test`). Current tests live in `common/`.
- Name tests with `TestXxx` and files as `*_test.go`.
- Prefer fast, deterministic unit tests; add integration tests when touching DB logic.

## Commit & Pull Request Guidelines
- Commit messages in this repo are short, imperative, and capitalized (e.g., “Add swagger to the tokens”).
- PRs should include: a clear summary, the rationale for changes, and any relevant migration/SQLC/Swagger updates.
- Link related issues or TODOs when applicable; include example requests/responses for API changes.

## Security & Configuration Tips
- Do not commit secrets or real credentials in `.env`.
- When changing schema: add a migration in `migrations/` and run `make sqlc`.
