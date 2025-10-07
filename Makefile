.PHONY: build run clean install help

# Variables
BINARY_NAME=labracodabrador
BUILD_DIR=bin
CMD_DIR=cmd/labracodabrador

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run with custom config
run-config: build
	@echo "Starting $(BINARY_NAME) with custom config..."
	./$(BUILD_DIR)/$(BINARY_NAME) -config=$(CONFIG) -genesis=$(GENESIS)

# Clean build artifacts and data
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -rf data/
	rm -rf logs/
	@echo "Clean complete"

# Clean only data (keep binary)
clean-data:
	@echo "Cleaning data directories..."
	rm -rf data/
	@echo "Data cleaned"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run ./...

# Create data directories
init-dirs:
	@echo "Creating directories..."
	mkdir -p data logs
	@echo "Directories created"

# Help command
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make install      - Install dependencies"
	@echo "  make run          - Build and run the application"
	@echo "  make run-config   - Run with custom config (CONFIG=path GENESIS=path)"
	@echo "  make clean        - Remove build artifacts and data"
	@echo "  make clean-data   - Remove only data directories"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Run linters"
	@echo "  make init-dirs    - Create necessary directories"
	@echo "  make help         - Show this help message"

