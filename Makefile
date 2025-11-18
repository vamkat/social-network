# Makefile for dev workflow: start/stop containers, seed DBs, clear DBs with migrations
# Works on Windows (with Git Bash/WSL), Linux, Mac

DOCKER_COMPOSE=docker-compose
SEED_SCRIPT=/app/seeds/seed.sh
USERS_CONTAINER=users
USERS_DB_CONTAINER=users-db
SHELL=/bin/sh

.PHONY: up down build rebuild seed clear clear-users clear-forum clear-chat clear-notifications reset-users

# ---------------------------
# Start all containers
# ---------------------------
up:
	@echo "Starting all containers..."
	$(DOCKER_COMPOSE) up -d
	@echo "All containers started."

# ---------------------------
# Stop all containers
# ---------------------------
down:
	@echo "Stopping all containers..."
	$(DOCKER_COMPOSE) down
	@echo "All containers stopped."

# ---------------------------
# Build all images
# ---------------------------
build:
	@echo "Building all service images..."
	$(DOCKER_COMPOSE) build
	@echo "Build complete."

# ---------------------------
# Rebuild (down + build + up)
# ---------------------------
rebuild: down build up
	@echo "Rebuild complete."

# ---------------------------
# Seed dev database (users only)
# ---------------------------
seed:
	@echo "Waiting for users DB to be ready..."
	@until $(DOCKER_COMPOSE) exec -T users-db pg_isready -U postgres -d social_users >/dev/null 2>&1; do \
		echo "Waiting for users-db..."; sleep 1; \
	done
	@echo "Users DB is ready. Seeding users database..."
	$(DOCKER_COMPOSE) exec -T users sh -c "$(SEED_SCRIPT)"
	@echo "Seeding complete."

# ---------------------------
# Clear all databases and restart containers
# ---------------------------
clear: clear-users clear-forum clear-chat clear-notifications
	@echo "All databases cleared."
	@echo "Waiting a few seconds for services to restart and migrations to run..."
	@sleep 5
	@echo "Done. You can now run 'make seed' to seed the dev database."

# ---------------------------
# Clear individual databases
# ---------------------------
clear-users:
	@echo "Clearing users database..."
	$(DOCKER_COMPOSE) exec -T $(USERS_DB_CONTAINER) psql -U postgres -d social_users -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "Users database cleared. Restarting users container to apply migrations..."
	$(DOCKER_COMPOSE) restart $(USERS_CONTAINER)
	@sleep 5
	@echo "Users container restarted. Migrations should have run."

clear-forum:
	@echo "Clearing forum database..."
	$(DOCKER_COMPOSE) exec -T forum-db psql -U postgres -d social_forum -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(DOCKER_COMPOSE) restart forum
	@sleep 5
	@echo "Forum container restarted. Migrations should have run."

clear-chat:
	@echo "Clearing chat database..."
	$(DOCKER_COMPOSE) exec -T chat-db psql -U postgres -d social_chat -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(DOCKER_COMPOSE) restart chat
	@sleep 5
	@echo "Chat container restarted. Migrations should have run."

clear-notifications:
	@echo "Clearing notifications database..."
	$(DOCKER_COMPOSE) exec -T notifications-db psql -U postgres -d social_notifications -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	$(DOCKER_COMPOSE) restart notifications
	@sleep 5
	@echo "Notifications container restarted. Migrations should have run."

# ---------------------------
# Reset users DB: clear + seed in one command
# ---------------------------
reset-users: clear-users seed
	@echo "Users database reset and seeded successfully."
