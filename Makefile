MAKEFLAGS += --no-print-directory

BACKEND_DIR := backend
FRONTEND_DIR := frontend-admin

# === Backend Commands ===

.PHONY: down rebuild dev-up test lint
down:
	$(MAKE) -C $(BACKEND_DIR) down

rebuild:
	$(MAKE) -C $(BACKEND_DIR) rebuild

dev-up:
	$(MAKE) -C $(BACKEND_DIR) dev-up

test:
	$(MAKE) -C $(BACKEND_DIR) test

lint:
	$(MAKE) -C $(BACKEND_DIR) lint

# === Frontend Commands ===

.PHONY: frontend-dev frontend-install frontend-lint frontend-build

frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-lint:
	cd $(FRONTEND_DIR) && npm run lint

frontend-build:
	cd $(FRONTEND_DIR) && npm run build
