# ---- Build stage ----
FROM base AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build notifications service binaries
RUN go build -o notifications_server ./services/notifications/cmd/server
RUN go build -o notifications_migrate ./services/notifications/cmd/migrate

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binaries and migrations
COPY --from=build /app/notifications_server .
COPY --from=build /app/notifications_migrate .
COPY --from=build /app/services/notifications/migrations ./migrations

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@notifications-db:5432/social_notifications?sslmode=disable

CMD ["./notifications_server"]
