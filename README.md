# Prompt Explorer

A web application built with **Go**, **HTMX**, and **Tailwind CSS** for testing and managing system prompts with the Gemini API. It allows developers to modify prompts, interact with an LLM, and store full conversation histories — ideal for **prompt engineering**, **SQL query generation**, and **manual validation**.

## Features

- 🧠 **System Prompt Editor** - Large, scrollable text area for prompt definition
- 💬 **Chat Interface** - Interactive chat with AI models
- 🗂️ **Conversation Navigation** - Navigate through previous prompts and conversations
- ⭐ **Favorites** - Mark conversations as favorites
- 🆕 **New Conversation** - Start fresh conversations while keeping system prompts
- 📊 **Conversation History** - Full conversation tracking and storage

## Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make

### Setup

```bash
# Clone the repository
git clone https://github.com/ventrosky/prompt-explorer.git
cd prompt-explorer

# Setup everything (dependencies + database)
make setup

# Run the application
make dev
```

The application will be available at `http://localhost:8081`

### Manual Setup

If you prefer step-by-step:

```bash
# 1. Install dependencies
make deps

# 2. Configure environment variables
cp env.example .env
# Edit .env file with your configuration

# 3. Start PostgreSQL database
make db-up

# 4. Run the application
make dev
```

### Environment Configuration

The application supports configuration via environment variables or a `.env` file:

**Option 1: Environment Variables**
```bash
export GEMINI_API_KEY=your_api_key_here
export SERVER_PORT=8081
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/prompt_explorer?sslmode=disable
```

**Option 2: .env File**
```bash
# Copy the example file
cp env.example .env

# Edit .env with your configuration
nano .env
```

**Required Environment Variables:**
- `GEMINI_API_KEY` - Your Gemini API key (get from [Google AI Studio](https://aistudio.google.com/app/apikey))
- `SERVER_PORT` - Server port (default: 8081)
- `DATABASE_URL` - PostgreSQL connection string (default: localhost)
- `ENVIRONMENT` - Environment mode (development/production)

**AI Model:**
- Uses **Gemini 2.5 Pro** - Google's most capable model for complex reasoning and analysis

## Available Commands

Run `make help` to see all available commands.

### Development Commands

```bash
make dev    # Run development server
make build  # Build the application
make run    # Build and run the application
make test   # Run all tests
make clean  # Clean build artifacts
```

### Database Commands

```bash
make db-up        # Start PostgreSQL database
make db-down      # Stop PostgreSQL database
make db-logs      # View database logs
make db-connect   # Connect to database
make db-reset     # Reset database (⚠️ deletes all data)
```

### Setup Commands

```bash
make setup        # Complete development setup
make deps         # Install dependencies
```

## Project Structure

```
prompt-explorer/
├── cmd/server/           # Application entry point
├── api/                  # HTTP handlers and routing
├── app/                  # Application orchestration
├── business/             # Domain logic
│   ├── models/          # Domain models
│   ├── services/        # Business services
│   └── repositories/    # Data access layer
├── foundation/           # Infrastructure utilities
│   ├── config/          # Configuration management
│   ├── database/        # Database layer
│   └── logger/          # Logging utilities
├── web/                  # Frontend assets
│   ├── templates/       # HTML templates
│   ├── static/          # Static files
│   └── assets/          # Web assets
├── docker/               # Docker configuration
│   ├── docker-compose.yml
│   └── init.sql
└── tests/               # Test files
```

## Database

The application uses PostgreSQL running in Docker with two databases:

- **prompt_explorer** - Main application database
- **test_prompt_explorer** - Test database for integration tests

The database includes:

- **conversations** - Conversation tracking with titles and favorites
- **messages** - Chat messages with conversation references

### Database Schema

**Conversations Table:**
- `id` (UUID) - Unique identifier
- `title` (TEXT) - Conversation title
- `is_favorite` (BOOLEAN) - Favorite status
- `created_at` (TIMESTAMP) - Creation time

**Messages Table:**
- `id` (UUID) - Unique identifier
- `conversation_id` (UUID) - Reference to conversation
- `sender` (ENUM) - Message type: 'system', 'user', 'assistant'
- `content` (TEXT) - Message content
- `created_at` (TIMESTAMP) - Creation time

### Testing

The project includes both unit tests and integration tests:

- **Unit Tests** - Fast, isolated tests that don't require a database
- **Integration Tests** - Tests that run against a real PostgreSQL database

To run all tests, ensure the database is running:

```bash
# Start the database
make db-up

# Run all tests (unit + integration)
make test
```

## Configuration

Configuration is handled through environment variables:

```bash
# Server Configuration
SERVER_PORT=8081
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/prompt_explorer?sslmode=disable
```

## Development

### Running Tests

```bash
make test
```

### Database Management

```bash
# Start database
make db-up

# View database logs
make db-logs

# Connect to database
make db-connect

# Reset database (⚠️ deletes all data)
make db-reset
```

### Hot Reload

The application supports hot reload during development:

```bash
make dev
```
