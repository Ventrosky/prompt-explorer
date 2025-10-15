.PHONY: help build run dev test clean

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building Prompt Explorer..."
	go build -o bin/prompt-explorer ./cmd/server

run: build ## Build and run the application
	@echo "Running Prompt Explorer..."
	./bin/prompt-explorer

dev: ## Run the application in development mode
	@echo "Starting development server..."
	go run ./cmd/server

test: ## Run tests
	@echo "Running tests..."
	go test ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/ tmp/
