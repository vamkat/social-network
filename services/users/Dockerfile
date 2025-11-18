# ---- Build stage ----
FROM golang:1.25-alpine AS build

WORKDIR /app

# Copy root module files
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy entire repo so shared/ and service code are available
COPY ../../ .

# Build users service binary
RUN go build -o users_service ./services/users/cmd/server

# ---- Runtime stage ----
FROM alpine:latest

WORKDIR /app

# Copy built binary and migrations
COPY --from=build /app/users_service .
COPY --from=build /app/services/users/internal/db/migrations ./migrations
COPY --from=build /app/services/users/internal/db/seeds ./seeds

# Install PostgreSQL client so `psql` is available for the seed script
RUN apk add --no-cache postgresql-client

# Ensure the seed script is executable inside the image
RUN chmod +x /app/seeds/seed.sh

# Env vars (can also be set in docker-compose)
ENV DATABASE_URL=postgres://postgres:secret@users-db:5432/social_users?sslmode=disable

CMD ["./users_service"]
