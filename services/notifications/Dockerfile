# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build notifications service binary
RUN go build -o notifications_service ./services/notifications

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binary and migrations
COPY --from=build /app/notifications_service .
COPY --from=build /app/services/notifications/migrations ./migrations

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@notifications-db:5432/social_notifications?sslmode=disable

CMD ["./notifications_service"]
