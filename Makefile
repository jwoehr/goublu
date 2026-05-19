# Makefile for Goublu
# A Go frontend for Ublu providing enhanced console interface

# Variables
BINARY_NAME := goublu
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.5.0-dev")
BUILD_TIME := $(shell date -u +%Y-%m-%d.%H:%M:%S\(UTC\))
LDFLAGS := -ldflags "-X main.GoubluVersion=$(VERSION) -X main.CompileDate=$(BUILD_TIME)"
GOPATH := $(shell go env GOPATH)
INSTALL_PATH := $(GOPATH)/bin

# Default target
.DEFAULT_GOAL := build

# Phony targets (not actual files)
.PHONY: all build install clean test fmt vet lint help run

# Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "Build complete: ./$(BINARY_NAME)"

# Install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	go install $(LDFLAGS) .
	@echo "Installed to $(INSTALL_PATH)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	go clean
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Run linter (if golangci-lint is installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Run the binary (for testing)
run: build
	./$(BINARY_NAME)

# Build and run all checks
all: fmt vet test build
	@echo "All checks passed!"

# Show help
help:
	@echo "Goublu Makefile targets:"
	@echo "  make build    - Build the binary (default)"
	@echo "  make install  - Install to GOPATH/bin"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make lint     - Run golangci-lint (if installed)"
	@echo "  make run      - Build and run"
	@echo "  make all      - Run fmt, vet, test, and build"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Current version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"

# Made with Bob
