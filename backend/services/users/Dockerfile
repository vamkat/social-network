# ---- Build stage ----
FROM social-network/go-base AS build

WORKDIR /app

# Copy entire repository so all shared/internal code is included
COPY . .

# Build the users service binary
RUN go build -o users_service ./services/users/cmd/server

RUN go build -o migrate ./services/users/cmd/migrate

# ---- Runtime stage ----
FROM alpine:3.20

WORKDIR /app

# Only copy the compiled binary (tiny image!)
COPY --from=build /app/users_service .
COPY --from=build /app/migrate .
COPY --from=build /app/services/users/internal/db/migrations ./migrations
COPY --from=build /app/services/users/internal/db/seeds ./seeds

# Ensure seed script is executable
RUN chmod +x /app/seeds/seed.sh

CMD ["./users_service"]
