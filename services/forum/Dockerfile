# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build forum service binary
RUN go build -o forum_service ./services/forum

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binary and migrations
COPY --from=build /app/forum_service .
COPY --from=build /app/services/forum/migrations ./migrations

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@forum-db:5432/social_forum?sslmode=disable

CMD ["./forum_service"]
