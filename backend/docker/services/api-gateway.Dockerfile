FROM social-network/go-base AS build

WORKDIR /app/backend

COPY backend/ .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o api_gateway ./services/gateway/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/backend/api_gateway .

CMD ["./api_gateway"]