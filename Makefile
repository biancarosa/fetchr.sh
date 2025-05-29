.PHONY: build test lint clean coverage release version commit-feat commit-fix commit-refactor commit-perf commit-security

# Build variables
BINARY_NAME=fetchr
GO=go
VERSION := $(shell git describe --tags --always --dirty)

# Build the application
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/fetchr

# Run unit tests with coverage
test:
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./... -tags=unit

e2e:
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./test/e2e/... -tags=e2e

# Run golangci-lint
lint:
	golangci-lint run

# Show coverage report
coverage:
	$(GO) tool cover -html=coverage.txt

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME) coverage.txt
	$(GO) clean

# Install dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Run all checks (test, lint, build)
check: test lint build

# Show current version
version:
	@echo "Current version: $(VERSION)"

# Conventional commit helpers
commit-feat:
	@if [ "$(msg)" = "" ]; then \
		echo "Error: commit message is required. Usage: make commit-feat msg='your message'"; \
		exit 1; \
	fi
	git commit -m "feat: $(msg)"

commit-fix:
	@if [ "$(msg)" = "" ]; then \
		echo "Error: commit message is required. Usage: make commit-fix msg='your message'"; \
		exit 1; \
	fi
	git commit -m "fix: $(msg)"

commit-refactor:
	@if [ "$(msg)" = "" ]; then \
		echo "Error: commit message is required. Usage: make commit-refactor msg='your message'"; \
		exit 1; \
	fi
	git commit -m "refactor: $(msg)"

commit-perf:
	@if [ "$(msg)" = "" ]; then \
		echo "Error: commit message is required. Usage: make commit-perf msg='your message'"; \
		exit 1; \
	fi
	git commit -m "perf: $(msg)"

commit-security:
	@if [ "$(msg)" = "" ]; then \
		echo "Error: commit message is required. Usage: make commit-security msg='your message'"; \
		exit 1; \
	fi
	git commit -m "security: $(msg)"

# Create a new release tag
release:
	@if [ "$(version)" = "" ]; then \
		echo "Error: version is required. Usage: make release version=1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release tag v$(version)..."
	git tag -a "v$(version)" -m "Release v$(version)"
	git push origin "v$(version)"

# Default target
all: deps check 