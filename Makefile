.PHONY: dev build run test clean tidy lint docker-build docker-up docker-down help

APP_NAME=halolight-api
PORT ?= 8000

# Development
dev:
	go run ./cmd/server

# Build
build:
	go build -ldflags="-w -s" -o bin/$(APP_NAME) ./cmd/server

# Run binary
run: build
	./bin/$(APP_NAME)

# Test
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run ./...

# Tidy dependencies
tidy:
	go mod tidy
	go mod verify

# Docker
docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Help
help:
	@echo "Available commands:"
	@echo "  make dev            - Run development server"
	@echo "  make build          - Build binary"
	@echo "  make run            - Build and run binary"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make lint           - Run linter"
	@echo "  make tidy           - Tidy go modules"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make clean          - Clean build artifacts"
