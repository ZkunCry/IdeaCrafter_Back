.PHONY: help run dev build test lint docker-build docker-up docker-down docker-logs

VERSION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
IMAGE   ?= startup-back

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

run: 
	go run ./cmd/app

dev: 
	air

build:
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w -X main.version=$(VERSION)" -o bin/app ./cmd/app

test: 
	go test ./...

lint:
	gofmt -l .
	go vet ./...

docker-build: 
	docker build --build-arg VERSION=$(VERSION) -t $(IMAGE):$(VERSION) -t $(IMAGE):latest .

docker-up:
	docker compose up -d --build

docker-down: 
	docker compose down

docker-logs: 
	docker compose logs -f api
