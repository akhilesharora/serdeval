.PHONY: all build test lint clean install run-web help

# Variables
BINARY_NAME=serdeval
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"

# Default target
all: lint test build

# Help target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the binary"
	@echo "  make test       - Run all tests"
	@echo "  make lint       - Run linters"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install the binary"
	@echo "  make run-web    - Run web interface"
	@echo "  make coverage   - Generate test coverage report"
	@echo "  make bench      - Run benchmarks"
	@echo "  make pre-commit - Install pre-commit hooks"

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY_NAME} ./cmd/serdeval
	@echo "Build complete: ./${BINARY_NAME}"

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME}
	@rm -f coverage.out coverage.html
	@rm -f benchmark.txt
	@echo "Clean complete"

# Install the binary
install: build
	@echo "Installing ${BINARY_NAME}..."
	@go install ${LDFLAGS}
	@echo "Installation complete"

# Run web interface
run-web: build
	@echo "Starting web interface on http://localhost:8080"
	@./${BINARY_NAME} web --port 8080

# Install pre-commit hooks
pre-commit:
	@echo "Installing pre-commit hooks..."
	@./scripts/setup-hooks.sh

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -w ${GO_FILES}
	@echo "Code formatted"

# Check if code is formatted
check-fmt:
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l ${GO_FILES})" ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		gofmt -d ${GO_FILES}; \
		exit 1; \
	fi
	@echo "Code is properly formatted"

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"