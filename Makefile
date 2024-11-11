# Makefile for Go project

# Variables
BINARY_NAME=hello-world-app
BUILD_DIR=bin

# Default target
all: tidy build test

# Tidy up the module dependencies
tidy:
	go mod tidy

copy-hugo:
	cp -r ui/public public

build-hugo:
	rm -rf public
	cd ui && hugo

docker-deploy:
	docker buildx build --platform linux/amd64 -t  $(DOCKER_TAG) .
	docker push $(DOCKER_TAG)
	
# Build only the Go application
build-go:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build the Go application
build: build-hugo copy-hugo
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

genkit_mode:
	export TEST_MODE=true
	genkit start

# Run tests
test:
	go test -cover ./...

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

integration-test: build
	export PORT=8080
	$(BUILD_DIR)/$(BINARY_NAME) &
	sleep 5
	BASE_URL=http://localhost:8080 go test -tags=integration ./...
	pkill -f $(BINARY_NAME)

.PHONY: all tidy build test clean run fmt lint 
