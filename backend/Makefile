SHELL := /bin/bash
.DEFAULT_GOAL := dev-up

ENV_FILE := .env

# Load env vars if .env exists
ifneq ("$(wildcard $(ENV_FILE))","")
	include $(ENV_FILE)
	export
endif

.PHONY: down rebuild migrate dev-up

## Stop and clean up all containers and volumes
down:
	@echo "==> Shutting down and removing all containers, networks, and volumes..."
	docker-compose down -v --remove-orphans

## Rebuild all images from scratch
rebuild:
	@echo "==> Building images from scratch..."
	docker-compose build --no-cache

## Run database migrations via the docker-compose-defined migrate service
migrate:
	@echo "==> Running migrations..."
	docker-compose up migrate

## Start up DB, run migrations, then launch the app
dev-up: down rebuild
	@echo "==> Starting database..."
	docker-compose up -d db
	@echo "==> Waiting for DB to be ready..."
	sleep 5
	@$(MAKE) migrate
	@echo "==> Starting app..."
	docker-compose up app

.PHONY: lint
lint:
	@echo "==> Running linters..."
	golangci-lint run ./...

.PHONY: test
test:
	@echo "==> Running unit tests..."
	go test ./... -v -cover -race -coverprofile=coverage.out
	go tool cover -func=coverage.out
