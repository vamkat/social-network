FROM social-network/go-base AS build

WORKDIR /app/backend

COPY backend/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o live ./services/live/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/backend/live .

CMD ["./live"]