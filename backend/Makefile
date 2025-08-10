# Makefile for running tests

.PHONY: test test-unit test-integration test-benchmark test-all test-coverage clean docker-test help

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Run all unit tests
test-unit: ## Run unit tests
	go test -v ./book/... ./process-url/... ./helper/...

# Run integration tests
test-integration: ## Run HTTP integration tests
	go test -v ./server/... -run TestHTTP

# Run benchmark tests
test-benchmark: ## Run benchmark tests
	go test -v ./server/... -bench=. -benchmem -run=^Benchmark

# Run all tests
test-all: ## Run all tests (unit + integration + benchmark)
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	go test -v ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests in parallel
test-parallel: ## Run tests in parallel
	go test -v ./... -parallel 4

# Run tests with race detection
test-race: ## Run tests with race detection
	go test -v ./... -race

# Run integration tests only
test-integration-only: ## Run integration tests only
	go test -v ./server/http_integration_test.go ./server/server.go ./server/router.go ./server/db.go -run TestHTTP

# Run specific test
test-specific: ## Run specific test (usage: make test-specific TEST=TestName)
	go test -v ./server/... -run $(TEST)

# Clean test artifacts
clean: ## Clean test artifacts
	rm -f coverage.out coverage.html
	go clean -testcache

# Docker test environment
docker-test: ## Run tests in Docker environment
	docker build -t byfood-api-test .
	docker run --rm byfood-api-test go test -v ./...

# Lint code
lint: ## Run linter
	golangci-lint run

# Format code
fmt: ## Format Go code
	go fmt ./...

# Download dependencies
deps: ## Download dependencies
	go mod download
	go mod tidy

# Run development server
run: ## Run development server
	go run cmd/main.go

# Build binary
build: ## Build binary
	go build -o bin/api cmd/main.go

# Run tests with verbose output and short mode
test-short: ## Run tests in short mode
	go test -v -short ./...

# Generate test coverage badge
coverage-badge: test-coverage ## Generate coverage badge
	@if command -v coverage-badge >/dev/null 2>&1; then \
		coverage-badge -f coverage.out -o coverage-badge.svg; \
		echo "Coverage badge generated: coverage-badge.svg"; \
	else \
		echo "coverage-badge tool not found. Install with: go install github.com/AlexBeauchemin/coverage-badge@latest"; \
	fi

# Continuous integration target
ci: clean deps fmt lint test-race test-coverage ## Run CI pipeline

# Development setup
setup: deps ## Setup development environment
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development environment setup complete!"
