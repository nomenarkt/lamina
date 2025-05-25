MAKEFLAGS += --no-print-directory

BACKEND_DIR := backend
FRONTEND_DIR := frontend-admin

# === 🧱 Backend Commands ===

.PHONY: app-up app-logs dev-up dev-reset dev-restart rebuild down migrate test lint

app-up: ## 🔁 Use after editing Go code (hot-reload via Air enabled)
	$(MAKE) -C $(BACKEND_DIR) app-up

app-logs: ## 📺 Use to tail server logs (debugging)
	$(MAKE) -C $(BACKEND_DIR) app-logs

dev-up: ## 🚀 Use when starting backend from scratch (DB + migrate + app)
	$(MAKE) -C $(BACKEND_DIR) dev-up

dev-reset: ## 💥 Use after schema or Docker config changes (full rebuild)
	$(MAKE) -C $(BACKEND_DIR) dev-reset

dev-restart: ## ♻️ Restart app only (DB unchanged)
	$(MAKE) -C $(BACKEND_DIR) dev-restart

rebuild: ## 🛠 Use after modifying Dockerfile, Go dependencies, or binary flags
	$(MAKE) -C $(BACKEND_DIR) rebuild

down: ## 🧹 Use to stop & clean up all backend containers/volumes
	$(MAKE) -C $(BACKEND_DIR) down

migrate: ## 🗄 Run DB migrations (run this after editing migrations/*.sql)
	$(MAKE) -C $(BACKEND_DIR) migrate

test: ## ✅ Run Go unit tests (always before commit/PR)
	$(MAKE) -C $(BACKEND_DIR) test

lint: ## 🧼 Run Go linters (required before merge)
	$(MAKE) -C $(BACKEND_DIR) lint

# === 🎨 Frontend Commands ===

.PHONY: frontend-dev frontend-lint frontend-test frontend-format

frontend-dev: ## 🧪 Use to launch Vite dev server for the admin UI
	cd $(FRONTEND_DIR) && npm run dev

frontend-lint: ## 🧹 Lint React/TS code (required before PRs)
	cd $(FRONTEND_DIR) && npm run lint

frontend-test: ## ✅ Run frontend unit tests
	cd $(FRONTEND_DIR) && npm test

frontend-format: ## 🧼 Auto-format frontend code (Prettier, etc.)
	cd $(FRONTEND_DIR) && npm run format

# === 📘 Help ===

.PHONY: help

help:
	@printf "\n\033[1m🧱 Backend Commands:\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E '^\s*(dev-|app-|down|test|lint|rebuild|migrate)' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

	@printf "\n\033[1m🎨 Frontend Commands:\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -E 'frontend-' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
