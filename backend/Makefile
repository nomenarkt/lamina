SHELL := /bin/bash
.DEFAULT_GOAL := dev-up
ENV_FILE := .env

# Load env vars if .env exists
ifneq ("$(wildcard $(ENV_FILE))","")
	include $(ENV_FILE)
	export
endif

# === 🛠 Core Commands ===

.PHONY: dev-up dev-reset db-up app-up dev-restart migrate app-logs rebuild down test lint help

dev-up: ## Start up DB, run migrations, then launch the app (Air-enabled)
	@$(MAKE) db-up
	@$(MAKE) migrate
	@$(MAKE) app-up

dev-reset: down rebuild dev-up ## Full clean up, rebuild, then launch dev services

db-up: ## Start up database only
	@echo "==> Starting database..."
	docker compose up -d db
	@echo "==> Waiting for DB to be ready..."
	sleep 3

app-up: ## Start up app service only (assuming DB is running)
	@echo "==> Starting app..."
	docker compose up -d app

dev-restart: ## Restart only the app container
	docker compose restart app

migrate: ## Run database migrations via the migrate container
	@echo "==> Running migrations..."
	docker compose run --rm migrate

app-logs: ## Tail the app container logs
	docker compose logs -f app

rebuild: ## Rebuild all Docker images (no-cache)
	@echo "==> Building images from scratch..."
	docker compose build --no-cache

down: ## Stop and clean up all containers, networks, and volumes
	@echo "==> Shutting down and removing all containers, networks, and volumes..."
	docker compose down -v --remove-orphans
	@echo "==> Removing dangling images..."
	docker image prune -f
	@echo "==> Removing unused volumes..."
	docker volume prune -f

# === ✅ QA Commands ===

test: ## Run unit tests with race and coverage
	@echo "==> Running unit tests..."
	go test -mod=vendor ./... -v -cover -race -coverprofile=coverage.out
	@go tool cover -func=coverage.out

lint: ## Run Go linters
	@echo "==> Running linters..."
	golangci-lint run ./...

# === 📘 Help ===

help: ## Show list of backend make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
