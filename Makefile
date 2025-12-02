.PHONY: help run build swagger sqlc migrate-up migrate-down migrate-create docker-up docker-down test

.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'

# Load .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

## run: Build and run the server
run: build
	./bin/server

## build: Build the server binary
build:
	go build -o bin/server ./cmd/server

## swagger: Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

## sqlc: Generate repository code with sqlc
sqlc:
	sqlc generate

## migrate-up: Run all pending migrations
migrate-up:
	goose -dir migrations postgres "$(DATABASE_URL)" up

## migrate-down: Rollback the last migration
migrate-down:
	goose -dir migrations postgres "$(DATABASE_URL)" down

## migrate-create: Create a new migration file
migrate-create:
	@read -p "Migration name: " name; \
	goose -dir migrations create $$name sql

## migrate-status: Show migration status
migrate-status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

## docker-up: Start docker compose services
docker-up:
	docker compose up -d

## docker-down: Stop docker compose services
docker-down:
	docker compose down

## test: Run all tests in verbose mode
test:
	go test -v ./...

## tools: Install dev dependencies (swag, goose, sqlc)
tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
