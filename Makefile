# Makefile for Go project

# Variables
BINARY_NAME=hello-world-app
BUILD_DIR=bin

# Default target
all: tidy build test

# Tidy up the module dependencies
tidy:
	go mod tidy

# Build the Go application
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Run tests
test:
	go test ./...

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Run the application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

.PHONY: all tidy build test clean run fmt lint