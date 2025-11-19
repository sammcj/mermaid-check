.PHONY: help build test lint clean

# Default target
.DEFAULT_GOAL := build

# Binary name
BINARY_NAME=go-mermaid
BUILD_DIR=.

help: ## Display available targets
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the CLI tool for current architecture
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-mermaid
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run all tests with coverage
	@echo "Running tests..."
	@go test -v -race -cover ./...

lint: ## Run golangci-lint on all packages
	@echo "Running moderise..."
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linting complete"

clean: ## Remove build artefacts
	@echo "Cleaning build artefacts..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -rf dist/
	@echo "Clean complete"

tidy: ## Tidy and verify module dependencies
	@echo "Tidying module dependencies..."
	@go mod tidy
	@go mod verify

coverage: ## Generate test coverage report with detailed statistics
	@echo "Generating coverage report..."
	@go test -coverpkg=./... -coverprofile=coverage.txt -covermode=atomic ./... 2>&1 | grep -E "coverage:|ok"
	@echo ""
	@echo "Overall coverage by file:"
	@go tool cover -func=coverage.txt | grep -E "^github.com.*\.go:" | awk '{pkg=$$1; sub(/:[0-9]+:.*/, "", pkg); coverage[pkg]+=$$NF+0; count[pkg]++} END {for (p in coverage) printf "  %-60s %.1f%%\n", p, coverage[p]/count[p]}' | sort
	@echo ""
	@echo "Total coverage:"
	@go tool cover -func=coverage.txt | tail -1
	@echo ""
	@echo "Files below 80% coverage (summary):"
	@go tool cover -func=coverage.txt | grep -E "^github.com.*\.go:" | awk '{pkg=$$1; sub(/:[0-9]+:.*/, "", pkg); coverage[pkg]+=$$NF+0; count[pkg]++} END {for (p in coverage) if (coverage[p]/count[p] < 80.0) printf "  %-60s %.1f%%\n", p, coverage[p]/count[p]}' | sort | head -20
	@echo ""
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "HTML coverage report: coverage.html"
