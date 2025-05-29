.PHONY: build test lint clean coverage release version install check ci help all dev commit-feat commit-fix commit-refactor commit-perf commit-security

# Build variables
BINARY_NAME=fetchr
GO=go
VERSION := $(shell git describe --tags --always --dirty)
GOLANGCI_LINT_VERSION=v2.1.6

# Show help
help:
	@echo "Available targets:"
	@echo "  help        - Show this help message"
	@echo "  install     - Install development tools (golangci-lint, goimports)"
	@echo "  deps        - Download and tidy Go modules"
	@echo "  build       - Build the application"
	@echo "  test        - Run unit tests with coverage"
	@echo "  e2e         - Run end-to-end tests"
	@echo "  lint        - Run golangci-lint"
	@echo "  check       - Run all checks (deps, test, e2e, lint, build)"
	@echo "  ci          - Run all CI checks (same as check)"
	@echo "  coverage    - Show coverage report in browser"
	@echo "  clean       - Clean build artifacts"
	@echo "  version     - Show current version"
	@echo "  release     - Create a new release tag (usage: make release version=1.0.0)"
	@echo ""
	@echo "Commit helpers:"
	@echo "  commit-feat      - Create a feature commit (usage: make commit-feat msg='your message')"
	@echo "  commit-fix       - Create a fix commit (usage: make commit-fix msg='your message')"
	@echo "  commit-refactor  - Create a refactor commit (usage: make commit-refactor msg='your message')"
	@echo "  commit-perf      - Create a performance commit (usage: make commit-perf msg='your message')"
	@echo "  commit-security  - Create a security commit (usage: make commit-security msg='your message')"

# Install development tools
install:
	@echo "Installing development tools..."
	@echo "Checking Go version..."
	@go_version=$$(go version | sed 's/go version go\([0-9.]*\).*/\1/'); \
	required_version="1.24"; \
	if ! printf '%s\n%s\n' "$$required_version" "$$go_version" | sort -V -C; then \
		echo "❌ Error: Go $$go_version detected, but Go $$required_version or higher is required"; \
		echo ""; \
		echo "Please install Go 1.24 or higher:"; \
		echo "  - Download from: https://golang.org/dl/"; \
		echo "  - Or use Homebrew: brew install go"; \
		echo "  - Or use official installer script"; \
		echo ""; \
		exit 1; \
	else \
		echo "✅ Go $$go_version detected (meets requirement: >=$$required_version)"; \
	fi
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	@echo "✅ golangci-lint installed successfully"
	@echo "Installing other Go tools..."
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ goimports installed successfully"
	@echo "✅ All development tools installed!"
	@echo ""
	@echo "Make sure $$(go env GOPATH)/bin is in your PATH:"
	@echo "  export PATH=\$$PATH:\$$(go env GOPATH)/bin"

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

# Run all checks (deps, test, e2e, lint, build)
check: deps test e2e lint build

# Run all checks for CI (same as check but with explicit steps)
ci: deps test e2e lint build
	@echo "✅ All CI checks passed!"

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
all: help

# Development workflow
dev: deps check 