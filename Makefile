.PHONY: help build run dev test clean fe-deps fe-dev fe-build fe-preview fe-check dev-full

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

test: ## Run all tests
	@echo "Running all tests..."
	go test ./...

lint: ## Run linters (using built-in Go tools)
	@echo "Running Go linters..."
	go fmt ./...
	go vet ./...
	go mod tidy
	@echo "Linting complete!"

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
	@make setup-env
	@make deps
	@make db-up
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Development environment ready!"
	@echo "To start the full application:"
	@echo "  Backend only: make dev"
	@echo "  Frontend only: make fe-dev"
	@echo "  Both together: make dev-full"

setup-env: ## Setup environment file
	@echo "Setting up environment configuration..."
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "Created .env file from env.example"; \
		echo "Please edit .env file with your configuration"; \
	else \
		echo ".env file already exists"; \
	fi

deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Installing React frontend dependencies..."
	cd fe && npm install

# Frontend commands
fe-deps: ## Install React frontend dependencies
	@echo "Installing React frontend dependencies..."
	@if ! command -v npm >/dev/null 2>&1; then \
		echo "Error: npm is not installed. Please install Node.js and npm first."; \
		echo "Visit: https://nodejs.org/"; \
		exit 1; \
	fi
	cd fe && npm install

fe-dev: ## Start React development server
	@echo "Starting React development server..."
	cd fe && npm run dev

fe-build: ## Build React frontend for production
	@echo "Building React frontend..."
	cd fe && npm run build

fe-preview: ## Preview React production build
	@echo "Previewing React production build..."
	cd fe && npm run preview

fe-check: ## Check if React frontend is properly set up
	@echo "Checking React frontend setup..."
	@if [ ! -d "fe" ]; then \
		echo "Error: fe directory not found"; \
		exit 1; \
	fi
	@if [ ! -f "fe/package.json" ]; then \
		echo "Error: fe/package.json not found"; \
		exit 1; \
	fi
	@if [ ! -d "fe/node_modules" ]; then \
		echo "Warning: node_modules not found. Run 'make fe-deps' first"; \
	fi
	@echo "React frontend setup looks good!"

# Full development setup
dev-full: ## Start both backend and frontend in development mode
	@echo "Starting full development environment..."
	@echo "Backend will run on http://localhost:8081"
	@echo "Frontend will run on http://localhost:3000"
	@echo "Starting backend..."
	@make dev &
	@sleep 3
	@echo "Starting frontend..."
	@make fe-dev