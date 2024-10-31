# Makefile for Go project

# Variables
BINARY_NAME=hello-world-app
BUILD_DIR=bin

# Default target
all: tidy build test

# Tidy up the module dependencies
tidy:
	go mod tidy

build-hugo:
	rm -rf public
	cd ui && hugo && cp -r public ../public
# Build the Go application
build: build-hugo
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

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