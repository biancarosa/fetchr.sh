.PHONY: build test lint clean coverage release version install check ci help all dev commit-feat commit-fix commit-refactor commit-perf commit-security
.PHONY: build-dashboard test-dashboard lint-dashboard clean-dashboard install-dashboard dashboard-dev serve-dev build-all clean-all

# Build variables
BINARY_NAME=fetchr
GO=go
VERSION := $(shell git describe --tags --always --dirty)
GOLANGCI_LINT_VERSION=v2.1.6
DASHBOARD_DIR=dashboard

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Setup commands:"
	@echo "  install          - Install Go development tools (golangci-lint, goimports)"
	@echo "  deps             - Download and tidy Go modules"
	@echo "  install-dashboard - Install dashboard dependencies"
	@echo ""
	@echo "Build commands:"
	@echo "  build            - Build the Go application"
	@echo "  build-dashboard  - Build the dashboard for production"
	@echo "  build-all        - Build both backend and dashboard"
	@echo ""
	@echo "Test commands:"
	@echo "  test             - Run Go unit tests with coverage"
	@echo "  test-dashboard   - Run dashboard tests"
	@echo "  e2e              - Run end-to-end tests"
	@echo ""
	@echo "Lint commands:"
	@echo "  lint             - Run golangci-lint for Go code"
	@echo "  lint-dashboard   - Run ESLint for dashboard code"
	@echo ""
	@echo "Development commands:"
	@echo "  dev              - Start both backend and dashboard in development mode (fast start)"
	@echo "  dev-safe         - Start development mode after running all checks (tests, lint, build)"
	@echo "  serve-dev        - Start only the backend in development mode"
	@echo "  dashboard-dev    - Start only the dashboard in development mode"
	@echo ""
	@echo "Quality commands:"
	@echo "  check            - Run all Go checks (deps, test, e2e, lint, build)"
	@echo "  ci               - Run all CI checks (same as check)"
	@echo "  coverage         - Show Go coverage report in browser"
	@echo ""
	@echo "Utility commands:"
	@echo "  clean            - Clean Go build artifacts"
	@echo "  clean-dashboard  - Clean dashboard build artifacts"
	@echo "  clean-all        - Clean all build artifacts"
	@echo "  version          - Show current version"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Release commands:"
	@echo "  release          - Create a new release tag (usage: make release version=1.0.0)"
	@echo ""
	@echo "Commit helpers:"
	@echo "  commit-feat      - Create a feature commit (usage: make commit-feat msg='your message')"
	@echo "  commit-fix       - Create a fix commit (usage: make commit-fix msg='your message')"
	@echo "  commit-refactor  - Create a refactor commit (usage: make commit-refactor msg='your message')"
	@echo "  commit-perf      - Create a performance commit (usage: make commit-perf msg='your message')"
	@echo "  commit-security  - Create a security commit (usage: make commit-security msg='your message')"

# Install development tools
install:
	@echo "Installing Go development tools..."
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
	@echo "✅ All Go development tools installed!"
	@echo ""
	@echo "Make sure $$(go env GOPATH)/bin is in your PATH:"
	@echo "  export PATH=\$$PATH:\$$(go env GOPATH)/bin"

# Install dashboard dependencies
install-dashboard:
	@echo "Installing dashboard dependencies..."
	@if [ ! -d "$(DASHBOARD_DIR)" ]; then \
		echo "❌ Error: Dashboard directory '$(DASHBOARD_DIR)' not found"; \
		echo "Please ensure the dashboard code is in the '$(DASHBOARD_DIR)' directory"; \
		exit 1; \
	fi
	@echo "Checking Node.js version..."
	@node_version=$$(node --version | sed 's/v\([0-9]*\).*/\1/'); \
	required_version="18"; \
	if [ "$$node_version" -lt "$$required_version" ]; then \
		echo "❌ Error: Node.js $$node_version detected, but Node.js $$required_version or higher is required"; \
		echo ""; \
		echo "Please install Node.js 18 or higher:"; \
		echo "  - Download from: https://nodejs.org/"; \
		echo "  - Or use nvm: nvm install 18"; \
		echo "  - Or use Homebrew: brew install node"; \
		echo ""; \
		exit 1; \
	else \
		echo "✅ Node.js $$node_version detected (meets requirement: >=$$required_version)"; \
	fi
	@cd $(DASHBOARD_DIR) && npm install --legacy-peer-deps
	@echo "✅ Dashboard dependencies installed successfully!"

# Build the Go application
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/fetchr

# Build the dashboard for production
build-dashboard:
	@if [ ! -d "$(DASHBOARD_DIR)" ]; then \
		echo "❌ Error: Dashboard directory '$(DASHBOARD_DIR)' not found"; \
		exit 1; \
	fi
	@cd $(DASHBOARD_DIR) && npm run build

# Build both backend and dashboard
build-all: build build-dashboard

# Run Go unit tests with coverage
test:
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./... -tags=unit

# Run dashboard tests
test-dashboard:
	@if [ ! -d "$(DASHBOARD_DIR)" ]; then \
		echo "❌ Error: Dashboard directory '$(DASHBOARD_DIR)' not found"; \
		exit 1; \
	fi
	@cd $(DASHBOARD_DIR) && npm run test --legacy-peer-deps

# Run end-to-end tests
e2e:
	$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./test/e2e/... -tags=e2e

# Run golangci-lint for Go code
lint:
	golangci-lint run

# Run ESLint for dashboard code
lint-dashboard:
	@if [ ! -d "$(DASHBOARD_DIR)" ]; then \
		echo "❌ Error: Dashboard directory '$(DASHBOARD_DIR)' not found"; \
		exit 1; \
	fi
	@cd $(DASHBOARD_DIR) && npm run lint --legacy-peer-deps

# Start both backend and dashboard in development mode
dev: deps
	@echo "Starting fetchr.sh in development mode..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Dashboard will be available at http://localhost:3000"
	@echo ""
	@echo "Press Ctrl+C to stop both services"
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
	$(MAKE) serve-dev & \
	$(MAKE) dashboard-dev & \
	wait

# Start development environment after running all checks
dev-safe: check
	@echo "All checks passed! Starting fetchr.sh in development mode..."
	@echo "Backend will be available at http://localhost:8080"
	@echo "Dashboard will be available at http://localhost:3000"
	@echo ""
	@echo "Press Ctrl+C to stop both services"
	@trap 'kill %1 %2 2>/dev/null; exit' INT; \
	$(MAKE) serve-dev & \
	$(MAKE) dashboard-dev & \
	wait

# Start only the backend in development mode
serve-dev:
	@echo "Starting backend in development mode on port 8080..."
	@$(GO) run ./cmd/fetchr serve --port 8080 --log-level debug --health --metrics

# Start only the dashboard in development mode
dashboard-dev:
	@if [ ! -d "$(DASHBOARD_DIR)" ]; then \
		echo "❌ Error: Dashboard directory '$(DASHBOARD_DIR)' not found"; \
		exit 1; \
	fi
	@echo "Starting dashboard in development mode on port 3000..."
	@cd $(DASHBOARD_DIR) && npm run dev

# Show Go coverage report
coverage:
	$(GO) tool cover -html=coverage.txt

# Clean Go build artifacts
clean:
	rm -f $(BINARY_NAME) coverage.txt
	$(GO) clean

# Clean dashboard build artifacts
clean-dashboard:
	@if [ -d "$(DASHBOARD_DIR)" ]; then \
		cd $(DASHBOARD_DIR) && rm -rf dist/ node_modules/.cache/; \
	fi

# Clean all build artifacts
clean-all: clean clean-dashboard
	@if [ -d "$(DASHBOARD_DIR)" ]; then \
		cd $(DASHBOARD_DIR) && rm -rf node_modules/; \
	fi

# Install Go dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

# Run all Go checks (deps, test, e2e, lint, build)
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