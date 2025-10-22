# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
BINARY_NAME := http-calendar
BINARY_DIR := bin
MAIN_PATH := ./cmd/server
LOG_DIR := logs

# Linting
GOLANGCI_LINT := golangci-lint

.PHONY: all build clean test test-coverage lint fmt mod-tidy run help

all: clean build test lint

build:
	@echo "Building binary..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

lint:
	@echo "Linting code..."
	@if ! command -v $(GOLANGCI_LINT) > /dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	$(GOLANGCI_LINT) run

fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

mod-tidy:
	@echo "Tidying Go modules..."
	$(GOCMD) mod tidy

run:
	@echo "Running application..."
	@mkdir -p $(LOG_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_DIR)/$(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  all           - Clean, build, test, and lint"
	@echo "  build         - Build the binary"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  mod-tidy      - Tidy Go modules"
	@echo "  run           - Build and run the application"
	@echo "  help          - Show this help message"
