BINARY_NAME=algopay
BUILD_DIR=build
MAIN_PATH=cmd/algopay

.PHONY: build run test clean install-deps tidy

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Run with custom port
run-port:
	@echo "Running $(BINARY_NAME) on port 9000..."
	@go run $(MAIN_PATH) -port=9000

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Initialize environment file
init-env:
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		cp .env.example .env; \
	else \
		echo ".env file already exists"; \
	fi

# Run development server with hot reload (requires air)
dev:
	@echo "Starting development server..."
	@air

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Builds completed in $(BUILD_DIR)/"