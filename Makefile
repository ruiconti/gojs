# Makefile for github.com/ruiconti/gojs

BINARY_NAME = gojs
WORKSPACE = github.com/ruiconti/gojs
BUILD_FLAGS = -race -trimpath

# Build target
build:
	@echo "Building the binary..."
	@go $BUILD_FLAGS build -o $(BINARY_NAME) $(WORKSPACE)

# Test target
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# PHONY targets
.PHONY: build test clean