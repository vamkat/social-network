# Run from light-infra
# make docker-up
# make run-migrations
# make run-all

.PHONY: \
	run-services run-all \
	run-api run-media run-chat run-posts run-users run-notifications\
	run-migrations \
	migrate-media migrate-users migrate-posts migrate-chat migrate-notifications\
	create-network \
	docker-up docker-down

######################################
# Shared environment (all services)
######################################

export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=
export REDIS_DB=0

export USERS_GRPC_ADDR=localhost:50051
export POSTS_GRPC_ADDR=localhost:50052
export CHAT_GRPC_ADDR=localhost:50053
export NOTIFICATIONS_ADDRESS=localhost:50054
export MEDIA_GRPC_ADDR=localhost:50055

export HTTP_ADDR=0.0.0.0:8081
export SHUTDOWN_TIMEOUT_SECONDS=5
export ENABLE_DEBUG_LOGS=true
export ENABLE_SIMPLE_PRINT=true



######################################
# Top-level commands
######################################

run-services: run-migrations run-all

######################################
# Run ALL services (parallel, blocking)
######################################

run-all:
	@echo "Starting all services..."
	@trap 'kill 0' INT TERM; \
	\
	ENC_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	JWT_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	PASSWORD_SECRET=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	TELEMETRY_COLLECTOR_ADDR=localhost:4317 \
	go run -C ../backend ./services/gateway/cmd/main.go 2>&1 | sed 's/^/[API] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5437/social_media?sslmode=disable \
	MIGRATE_PATH=services/media/internal/db/migrations \
	GRPC_SERVER_PORT=:50055 \
	MINIO_ENDPOINT=localhost:9000 \
	MINIO_ACCESS_KEY=minioadmin \
	MINIO_SECRET_KEY=minioadmin \
	go run -C ../backend ./services/media/cmd/server/main.go 2>&1 | sed 's/^/[MEDIA] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5435/social_chat?sslmode=disable \
	MIGRATE_PATH=services/chat/internal/db/migrations \
	GRPC_SERVER_PORT=:50053 \
	go run -C ../backend ./services/chat/cmd/server/main.go 2>&1 | sed 's/^/[CHAT] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5434/social_posts?sslmode=disable \
	GRPC_SERVER_PORT=:50052 \
	go run -C ../backend ./services/posts/cmd/server/main.go 2>&1 | sed 's/^/[POSTS] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5433/social_users?sslmode=disable \
	MIGRATE_PATH=services/users/internal/db/migrations \
	GRPC_SERVER_PORT=:50051 \
	go run -C ../backend ./services/users/cmd/server/main.go 2>&1 | sed 's/^/[USERS] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5436/social_notifications?sslmode=disable \
	MIGRATE_PATH=services/notifications/internal/db/migrations \
	GRPC_SERVER_PORT=:50054 \
	go run -C ../backend ./services/notifications/cmd/server/main.go 2>&1 | sed 's/^/[NOTIFICATIONS] /' & \
	\
	wait


build-all:
	@echo "Starting building..."
	@trap 'kill 0' INT TERM; \
	OUTDIR="$(CURDIR)"; \
	cd ../backend && \
	go build -o "$$OUTDIR/gateway.social" ./services/gateway/cmd/main.go ; \
	go build -o "$$OUTDIR/media.social" ./services/media/cmd/server/main.go ; \
	go build -o "$$OUTDIR/chat.social" ./services/chat/cmd/server/main.go ; \
	go build -o "$$OUTDIR/posts.social" ./services/posts/cmd/server/main.go ; \
	go build -o "$$OUTDIR/users.social" ./services/users/cmd/server/main.go ; \
# 	go build -o "$$OUTDIR/notifs.social" ./services/notifications/cmd/server/main.go ; \
	@echo "Finished building!"
	wait



run-built:
	@echo "Starting all services..."
	@trap 'kill 0' INT TERM; \
	\

	ENC_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	JWT_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	PASSWORD_SECRET=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	TELEMETRY_COLLECTOR_ADDR=localhost:4317 \
	./gateway.social 2>&1 | sed 's/^/[API] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5437/social_media?sslmode=disable \
	MIGRATE_PATH=services/media/internal/db/migrations \
	GRPC_SERVER_PORT=:50055 \
	MINIO_ENDPOINT=localhost:9000 \
	MINIO_ACCESS_KEY=minioadmin \
	MINIO_SECRET_KEY=minioadmin \
	./media.social 2>&1 | sed 's/^/[MEDIA] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5435/social_chat?sslmode=disable \
	MIGRATE_PATH=services/chat/internal/db/migrations \
	GRPC_SERVER_PORT=:50053 \
	./chat.social 2>&1 | sed 's/^/[CHAT] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5434/social_posts?sslmode=disable \
	GRPC_SERVER_PORT=:50052 \
	./posts.social 2>&1 | sed 's/^/[POSTS] /' & \
	\
	DATABASE_URL=postgres://postgres:secret@localhost:5433/social_users?sslmode=disable \
	MIGRATE_PATH=services/users/internal/db/migrations \
	GRPC_SERVER_PORT=:50051 \
	./users.social 2>&1 | sed 's/^/[USERS] /' & \
	\
# 	DATABASE_URL=postgres://postgres:secret@localhost:5436/social_notifications?sslmode=disable \
# 	MIGRATE_PATH=services/notifications/internal/db/migrations \
# 	GRPC_SERVER_PORT=:50054 \
# 	./notifs.social 2>&1 | sed 's/^/[NOTIFICATIONS] /' & \
	\
	wait

######################################
# Individual services
######################################

run-api:
	ENC_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	JWT_KEY=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	PASSWORD_SECRET=a2F0LWFsZXgtdmFnLXlwYXQtc3RhbS16b25lMDEtZ28= \
	go run -C ../backend ./services/gateway/cmd/main.go 2>&1 | sed 's/^/[API] /'

run-users:
	DATABASE_URL=postgres://postgres:secret@localhost:5433/social_users?sslmode=disable \
	MIGRATE_PATH=services/users/internal/db/migrations \
	GRPC_SERVER_PORT=:50051 \
	go run -C ../backend ./services/users/cmd/server/main.go 2>&1 | sed 's/^/[USERS] /'

run-posts:
	DATABASE_URL=postgres://postgres:secret@localhost:5434/social_posts?sslmode=disable \
	GRPC_SERVER_PORT=:50052 \
	go run -C ../backend ./services/posts/cmd/server/main.go 2>&1 | sed 's/^/[POSTS] /'

run-chat:
	DATABASE_URL=postgres://postgres:secret@localhost:5435/social_chat?sslmode=disable \
	MIGRATE_PATH=services/chat/internal/db/migrations \
	GRPC_SERVER_PORT=:50053 \
	go run -C ../backend ./services/chat/cmd/server/main.go 2>&1 | sed 's/^/[CHAT] /'

run-media:
	DATABASE_URL=postgres://postgres:secret@localhost:5437/social_media?sslmode=disable \
	MIGRATE_PATH=services/media/internal/db/migrations \
	GRPC_SERVER_PORT=:50055 \
	MINIO_ENDPOINT=localhost:9000 \
	MINIO_ACCESS_KEY=minioadmin \
	MINIO_SECRET_KEY=minioadmin \
	go run -C ../backend ./services/media/cmd/server/main.go 2>&1 | sed 's/^/[MEDIA] /'

run-notifications:
	DATABASE_URL=postgres://postgres:secret@localhost:5436/social_notifications?sslmode=disable \
	MIGRATE_PATH=services/notifications/internal/db/migrations \
	GRPC_SERVER_PORT=:50054 \
	go run -C ../backend ./services/notifications/cmd/server/main.go 2>&1 | sed 's/^/[NOTIFICATIONS] /'

######################################
# Migrations
######################################

run-migrations: migrate-media migrate-users migrate-posts migrate-chat

migrate-media:
	DATABASE_URL=postgres://postgres:secret@localhost:5437/social_media?sslmode=disable \
	MIGRATE_PATH=services/media/internal/db/migrations \
	go run -C ../backend ./services/media/cmd/migrate

migrate-chat:
	DATABASE_URL=postgres://postgres:secret@localhost:5435/social_chat?sslmode=disable \
	MIGRATE_PATH=services/chat/internal/db/migrations \
	go run -C ../backend ./services/chat/cmd/migrate

migrate-users:
	DATABASE_URL=postgres://postgres:secret@localhost:5433/social_users?sslmode=disable \
	MIGRATE_PATH=services/users/internal/db/migrations \
	go run -C ../backend ./services/users/cmd/migrate

migrate-posts:
	DATABASE_URL=postgres://postgres:secret@localhost:5434/social_posts?sslmode=disable \
	MIGRATE_PATH=services/posts/internal/db/migrations \
	go run -C ../backend ./services/posts/cmd/migrate

migrate-notifications:
	DATABASE_URL=postgres://postgres:secret@localhost:5436/social_notifications?sslmode=disable \
	MIGRATE_PATH=services/notifications/internal/db/migrations \
	go run -C ../backend ./services/notifications/cmd/migrate

######################################
# Docker helpers
######################################

docker-up:
	$(MAKE) create-network
	docker compose -f docker-compose.infra.yml up -d

docker-down:
	docker compose -f docker-compose.infra.yml down

create-network:
	@docker network inspect social-network >/dev/null 2>&1 || \
	docker network create social-network

boom:
	$(MAKE) docker-up
	$(MAKE) run-migrations
	$(MAKE) run-all
