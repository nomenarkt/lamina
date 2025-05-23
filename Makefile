MAKEFLAGS += --no-print-directory

BACKEND_DIR := backend
FRONTEND_DIR := frontend-admin

# === 🧱 Backend Commands ===

.PHONY: app-up dev-up dev-reset rebuild down migrate test lint
app-up:
	$(MAKE) -C $(BACKEND_DIR) app-up

dev-up:
	$(MAKE) -C $(BACKEND_DIR) dev-up

dev-reset:
	$(MAKE) -C $(BACKEND_DIR) dev-reset

rebuild:
	$(MAKE) -C $(BACKEND_DIR) rebuild

down:
	$(MAKE) -C $(BACKEND_DIR) down

migrate:
	$(MAKE) -C $(BACKEND_DIR) migrate

test:
	$(MAKE) -C $(BACKEND_DIR) test

lint:
	$(MAKE) -C $(BACKEND_DIR) lint

# === 🎨 Frontend Commands ===

.PHONY: frontend-dev frontend-lint frontend-test frontend-format
frontend-dev: ## Start Vite dev server for React admin panel
	cd $(FRONTEND_DIR) && npm run dev

frontend-lint: ## Lint frontend code
	cd $(FRONTEND_DIR) && npm run lint

frontend-test: ## Run frontend unit tests
	cd $(FRONTEND_DIR) && npm test

frontend-format: ## Format frontend code
	cd $(FRONTEND_DIR) && npm run format

# === 📘 Help ===

.PHONY: help
help:
	@printf "\n\033[1m🧱 Backend Commands:\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(BACKEND_DIR)/Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

	@printf "\n\033[1m🎨 Frontend Commands:\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; /frontend-/ {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
