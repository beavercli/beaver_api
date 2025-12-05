.PHONY: help run build seed swagger sqlc migrate-up migrate-down migrate-create docker-up docker-down test

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

## seed: Seed the database with test data
seed:
	go run ./cmd/seed

## swagger: Generate Swagger documentation
swagger:
	swag init -g cmd/server/main.go -o docs

## sqlc: Generate repository code with sqlc
sqlc:
	sqlc generate

## migrate-up: Run all pending migrations
migrate-up:
	goose up

## migrate-down: Rollback the last migration
migrate-down:
	goose down

## migrate-create: Create a new migration file (usage: make migrate-create name=<migration_name>)
migrate-create:
	goose create $(name) sql

## migrate-status: Show migration status
migrate-status:
	goose status

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
