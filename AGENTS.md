# Repository Guidelines

## Project Structure & Module Organization
- `cmd/server`: entrypoint for the HTTP API server (Swagger annotations live here).
- `internal/router`: HTTP routing/handlers and OpenAPI comment blocks.
- `internal/service`: business logic and DTOs used by handlers.
- `internal/storage`: sqlc-generated data access layer; `internal/queries/queries.sql` is the source for generation.
- `common`: shared config/database helpers.
- `migrations`: Goose SQL migrations; run against the DB.
- `docs`: generated Swagger assets (`swagger.yaml/json`, `docs.go`).

## Build, Test, and Development Commands
- `make build`: compile the server binary to `bin/server`.
- `make run`: build then start the API locally.
- `make test`: run Go tests with verbose output.
- `make swagger`: regenerate Swagger docs from handler comments.
- `make sqlc`: regenerate storage code from `internal/queries/queries.sql`.
- `make migrate-up` / `make migrate-down`: apply or rollback DB migrations via Goose.
- `make docker-up` / `make docker-down`: start/stop compose services (DB, etc.).

## Coding Style & Naming Conventions
- Go code follows standard formatting; always run `gofmt` (the Make targets rely on it).
- Keep handler comment blocks (`@Summary`, `@Router`, etc.) accurate; they drive Swagger output.
- sqlc queries: add new statements to `internal/queries/queries.sql` with `-- name: <QueryName> :<result>`; run `make sqlc` afterward.
- Prefer clear, descriptive names; exported types/functions use PascalCase, locals use camelCase.

## Testing Guidelines
- Framework: Go’s built-in `testing` package.
- Command: `go test ./...` (or `make test`) before sending changes.
- Add table-driven tests for service logic; mock storage when possible to keep tests fast.

## Commit & Pull Request Guidelines
- Commit messages: short imperative subject (e.g., `Add contributor count query`).
- PRs should include: concise description, linked issue (if any), manual/automated test evidence, and notes on migrations or generated files (`docs`, `internal/storage`).
- If Swagger or sqlc output changes, commit the regenerated artifacts alongside the source changes.

## Security & Configuration Tips
- Configuration via env vars (`config` package); avoid committing secrets.
- Use `make migrate-status` to verify DB state before running data-changing commands.
