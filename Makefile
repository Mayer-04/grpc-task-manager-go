.PHONY: proto clean build run docker-up docker-down test help

# Variables
PROTO_PATH = proto
PKG_PATH = pkg/taskpb
GO_OUT = .

# Variables
PROTO_DIR = proto
PROTO_FILES = $(PROTO_DIR)/*.proto
PKG_DIR = pkg/taskpb
GO_MODULE = github.com/Mayer-04/grpc-task-manager-go

# Help command
help:
	@echo "Available commands:"
	@echo "  proto       - Generate Go code from protobuf files"
	@echo "  clean       - Clean generated files and build artifacts"
	@echo "  build       - Build the server binary"
	@echo "  run         - Run the server"
	@echo "  docker-up   - Start PostgreSQL with docker-compose"
	@echo "  docker-down - Stop PostgreSQL containers"
	@echo "  test        - Run tests"
	@echo "  help        - Show this help message"

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@mkdir -p $(PKG_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(PKG_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "Protobuf files generated in $(PKG_DIR)!"

# Build server
build: proto
	@echo "Building server..."
	go build -o cmd/server/server ./cmd/server
	@echo "Server built successfully!"

# Build and run client
client: proto
	@echo "Building client..."
	go build -o cmd/client/client ./cmd/client
	@echo "Starting client..."
	./cmd/client/client

# Run server
run: build
	@echo "Starting server..."
	./cmd/server/server

# Start PostgreSQL with docker-compose
docker-up:
	@echo "Starting PostgreSQL..."
	docker compose up -d
	@echo "PostgreSQL started successfully!"

# Stop docker containers
docker-down:
	@echo "Stopping containers..."
	docker compose down
	@echo "Containers stopped!"

# Create database and tables (run after docker-up)
db-setup: docker-up
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "Creating database and tables..."
	docker exec postgres-go psql -U $${POSTGRES_USER:-postgres} -d postgres -c "CREATE DATABASE IF NOT EXISTS taskdb;"
	docker exec postgres-go psql -U $${POSTGRES_USER:-postgres} -d taskdb -c "\
		CREATE TABLE IF NOT EXISTS tasks ( \
			id UUID PRIMARY KEY, \
			user_id VARCHAR(255) NOT NULL, \
			title VARCHAR(500) NOT NULL, \
			description TEXT, \
			completed BOOLEAN DEFAULT FALSE, \
			created_at TIMESTAMP DEFAULT NOW(), \
			updated_at TIMESTAMP DEFAULT NOW() \
		);"
	@echo "Database setup completed!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Development setup (run this first)
dev-setup: deps docker-up db-setup proto
	@echo "Development environment ready!"