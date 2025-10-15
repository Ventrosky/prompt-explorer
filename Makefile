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

# Database commands
db-up: ## Start PostgreSQL database
	@echo "Starting PostgreSQL database..."
	cd docker && docker-compose up -d

db-down: ## Stop PostgreSQL database
	@echo "Stopping PostgreSQL database..."
	cd docker && docker-compose down

db-logs: ## View database logs
	@echo "Showing database logs..."
	cd docker && docker-compose logs -f postgres

db-connect: ## Connect to database
	@echo "Connecting to database..."
	docker exec -it prompt-explorer-db psql -U postgres -d prompt_explorer

db-reset: ## Reset database (⚠️ deletes all data)
	@echo "Resetting database..."
	cd docker && docker-compose down -v && docker-compose up -d

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@make deps
	@make db-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Development environment ready!"

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
