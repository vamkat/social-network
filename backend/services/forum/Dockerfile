FROM social-network/go-base AS build

WORKDIR /app

COPY . .

RUN go build -o forum_service ./services/forum/cmd/server

RUN go build -o migrate ./services/forum/cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=build /app/forum_service .
COPY --from=build /app/migrate .
COPY --from=build /app/services/users/internal/db/migrations ./migrations
COPY --from=build /app/services/users/internal/db/seeds ./seeds

RUN chmod +x /app/seeds/seed.sh

CMD ["./forum_service"]
