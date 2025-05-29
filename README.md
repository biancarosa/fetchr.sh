# fetchr.sh

A powerful HTTP proxy and request capture tool.

## Features

- Basic HTTP proxy with request forwarding
- Configurable port and logging levels
- Support for all HTTP methods
- Header preservation
- Graceful shutdown
- Optional metrics and health check endpoints

## Installation

```bash
go install github.com/biancarosa/fetchr.sh/cmd/fetchr@latest
```

## Usage

### Starting the Proxy Server

```bash
fetchr serve --port 8080 --log-level info
```

Available flags:
- `--port`: Port to listen on (default: 8080)
- `--log-level`: Logging level (debug, info, warn, error)
- `--metrics`: Enable metrics endpoint
- `--health`: Enable health check endpoint

### Making Requests Through the Proxy

Configure your HTTP client to use the proxy at `http://localhost:8080` (or whatever port you specified).

Example using curl:
```bash
curl -x http://localhost:8080 https://api.example.com/endpoint
```

## Development

### Prerequisites

- Go 1.24 or later
- Git

### Quick Start

1. Clone the repository:
```bash
git clone https://github.com/biancarosa/fetchr.sh.git
cd fetchr.sh
```

2. Install development tools:
```bash
make install
```

3. Run all checks to ensure everything works:
```bash
make check
```

### Development Commands

The project uses a comprehensive Makefile for development tasks. Run `make` or `make help` to see all available commands:

#### Setup Commands
- `make install` - Install development tools (golangci-lint, goimports, etc.)
- `make deps` - Download and tidy Go modules

#### Build and Test Commands
- `make build` - Build the application
- `make test` - Run unit tests with coverage
- `make e2e` - Run end-to-end tests
- `make lint` - Run golangci-lint for code quality
- `make check` - Run all checks (deps, test, e2e, lint, build)
- `make ci` - Run all CI checks (same as check)

#### Utility Commands
- `make coverage` - Show coverage report in browser
- `make clean` - Clean build artifacts
- `make version` - Show current version
- `make help` - Show all available commands

#### Release Commands
- `make release version=1.0.0` - Create a new release tag

#### Conventional Commit Helpers
- `make commit-feat msg='your message'` - Create a feature commit
- `make commit-fix msg='your message'` - Create a fix commit
- `make commit-refactor msg='your message'` - Create a refactor commit
- `make commit-perf msg='your message'` - Create a performance commit
- `make commit-security msg='your message'` - Create a security commit

### Development Workflow

1. **Initial setup:**
   ```bash
   make install  # Install dev tools
   make deps     # Download dependencies
   ```

2. **During development:**
   ```bash
   make check    # Run all checks before committing
   ```

3. **Before pushing:**
   ```bash
   make ci       # Ensure CI will pass
   ```

### Building from Source

```bash
git clone https://github.com/biancarosa/fetchr.sh.git
cd fetchr.sh
make build
```

The binary will be created as `fetchr` in the project root.

### Running Tests

#### Unit Tests
```bash
make test
```

#### End-to-End Tests
```bash
make e2e
```

#### All Tests
```bash
make check  # Runs both unit and e2e tests plus linting
```

### Code Quality

The project uses golangci-lint for code quality checks. Run:

```bash
make lint
```

To see the coverage report in your browser:

```bash
make coverage
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
