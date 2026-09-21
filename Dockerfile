# syntax=docker/dockerfile:1


FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION=dev

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
        -trimpath \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o /out/app ./cmd/app

# ---------- runtime ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -g 10001 -S app \
 && adduser -u 10001 -S -G app app

WORKDIR /app

COPY --from=builder /out/app /app/app
COPY internal/platform/config/config.yaml /app/config.yaml

USER app:app

ENV SERVER_HOST=0.0.0.0 \
    SERVER_PORT=3001 \
    APP_ENV=production

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -q --spider "http://127.0.0.1:${SERVER_PORT}/health" || exit 1

ENTRYPOINT ["/app/app"]
