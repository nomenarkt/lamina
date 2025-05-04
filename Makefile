# Makefile — root of project

DOCKER_COMPOSE := docker-compose

# 🧼 Cleanup and rebuild from scratch
dev-up:
	$(DOCKER_COMPOSE) down -v --remove-orphans
	chmod +x wait-for-it.sh
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up app

# 🔁 Rebuild without stopping volumes
rebuild:
	$(DOCKER_COMPOSE) build --no-cache

# 🚀 Run only the app container
app:
	$(DOCKER_COMPOSE) up app

# 🧪 Run database migration
migrate:
	$(DOCKER_COMPOSE) run --rm migrate

# 🧹 Stop and clean everything
down:
	$(DOCKER_COMPOSE) down -v --remove-orphans

# 🐳 Show running containers
ps:
	$(DOCKER_COMPOSE) ps

# 🧭 See logs
logs:
	$(DOCKER_COMPOSE) logs -f app
