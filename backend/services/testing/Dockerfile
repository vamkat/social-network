# ---- Build stage ----
FROM social-network/go-base AS build

WORKDIR /app

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build -o tester ./services/testing/cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/tester .

CMD ["./tester"]