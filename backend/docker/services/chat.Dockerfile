FROM social-network/go-base AS build

WORKDIR /app/backend

COPY backend/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o chat_service ./services/chat/cmd/server

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o migrate ./services/chat/cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=build /app/backend/chat_service .
COPY --from=build /app/backend/migrate .
COPY --from=build /app/backend/services/chat/internal/db/migrations ./migrations


CMD ["./chat_service"]
