.PHONY: build run clean

# Binary name
BINARY_NAME=kex

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/kex

# Run the application (builds first)
run: build
	@echo "\nRunning $(BINARY_NAME)...\n"
	@./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME) 