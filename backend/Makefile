SHELL := /bin/bash
.DEFAULT_GOAL := dev-up
ENV_FILE := .env

# Load env vars if .env exists
ifneq ("$(wildcard $(ENV_FILE))","")
	include $(ENV_FILE)
	export
endif

# === 🛠 Core Commands ===

.PHONY: dev-up
dev-up: ## Start up DB, run migrations, then launch the app
	@$(MAKE) db-up
	@$(MAKE) migrate
	@$(MAKE) app-up

.PHONY: dev-reset
dev-reset: down rebuild dev-up ## Full clean up, rebuild, then launch dev services

.PHONY: db-up
db-up: ## Start up database only
	@echo "==> Starting database..."
	docker compose up -d db
	@echo "==> Waiting for DB to be ready..."
	sleep 3

.PHONY: app-up
app-up: ## Start up app service only (assuming DB is running)
	@echo "==> Starting app..."
	docker compose up app

.PHONY: migrate
migrate: ## Run database migrations via the migrate container
	@echo "==> Running migrations..."
	docker compose run --rm migrate

.PHONY: rebuild
rebuild: ## Rebuild all Docker images (no-cache)
	@echo "==> Building images from scratch..."
	docker compose build --no-cache

.PHONY: down
down: ## Stop and clean up all containers, networks, and volumes
	@echo "==> Shutting down and removing all containers, networks, and volumes..."
	docker compose down -v --remove-orphans
	@echo "==> Removing dangling images..."
	docker image prune -f
	@echo "==> Removing unused volumes..."
	docker volume prune -f


# === ✅ QA Commands ===

.PHONY: test
test: ## Run unit tests with race and coverage
	@echo "==> Running unit tests..."
	go test ./... -v -cover -race -coverprofile=coverage.out
	@go tool cover -func=coverage.out

.PHONY: lint
lint: ## Run Go linters
	@echo "==> Running linters..."
	golangci-lint run ./...

# === 📘 Help ===

.PHONY: help
help: ## Show list of backend make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
