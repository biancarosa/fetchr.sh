# fetchr.sh

A powerful HTTP proxy and request capture tool with a modern web dashboard.

## Features

- **HTTP Proxy**: Basic HTTP proxy with request forwarding
- **Web Dashboard**: Modern web interface for monitoring and managing proxy traffic
- **Request Capture**: Capture and analyze HTTP requests and responses
- **Configurable**: Configurable port and logging levels
- **Full HTTP Support**: Support for all HTTP methods with header preservation
- **Graceful Shutdown**: Clean shutdown with connection draining
- **Monitoring**: Optional metrics and health check endpoints
- **Real-time Updates**: Live dashboard updates for request monitoring

## Installation

```bash
go install github.com/biancarosa/fetchr.sh/cmd/fetchr@latest
```

## Usage

### Starting the Proxy Server

```bash
fetchr serve --port 8080 --log-level info --admin-port 8081
```

Available flags:
- `--port`: Port to listen on (default: 8080)
- `--log-level`: Logging level (debug, info, warn, error)
- `--admin-port`: Admin port for health checks and metrics (0 to disable, default: 0)
- `--dashboard`: Enable web dashboard (default: false)
- `--dashboard-port`: Dashboard port (default: 3000)

### Accessing the Dashboard

When the dashboard is enabled, you can access it at `http://localhost:3000` (or whatever port you specified with `--dashboard-port`).

The dashboard provides:
- Real-time request monitoring
- Request/response inspection
- Proxy configuration management
- Performance metrics and analytics
- Request filtering and search
- Export functionality for captured data

### Making Requests Through the Proxy

Configure your HTTP client to use the proxy at `http://localhost:8080` (or whatever port you specified).

Example using curl:
```bash
curl -x http://localhost:8080 https://api.example.com/endpoint
```

## Development

### Prerequisites

- Go 1.24 or later
- Node.js 18 or later (for dashboard development)
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

3. Install dashboard dependencies:
```bash
cd dashboard
npm install --legacy-peer-deps
cd ..
```

4. Run all checks to ensure everything works:
```bash
make check
```

### Development Commands

The project uses a comprehensive Makefile for development tasks. Run `make` or `make help` to see all available commands:

#### Setup Commands
- `make install` - Install Go development tools (golangci-lint, goimports, etc.)
- `make deps` - Download and tidy Go modules
- `make install-dashboard` - Install dashboard dependencies

#### Build and Test Commands
- `make build` - Build the Go application
- `make build-dashboard` - Build the dashboard for production
- `make build-all` - Build both backend and dashboard
- `make test` - Run unit tests with coverage
- `make test-dashboard` - Run dashboard tests
- `make e2e` - Run end-to-end tests
- `make lint` - Run golangci-lint for code quality
- `make lint-dashboard` - Run dashboard linting
- `make check` - Run all checks (deps, test, e2e, lint, build)
- `make ci` - Run all CI checks (same as check)

#### Development Server Commands
- `make dev` - Start both backend and dashboard in development mode
- `make serve-dev` - Start only the backend in development mode
- `make dashboard-dev` - Start only the dashboard in development mode

#### Utility Commands
- `make coverage` - Show Go coverage report in browser
- `make clean` - Clean build artifacts
- `make clean-dashboard` - Clean dashboard build artifacts
- `make clean-all` - Clean all build artifacts
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
   make install          # Install Go dev tools
   make deps             # Download Go dependencies
   make install-dashboard # Install dashboard dependencies
   ```

2. **During development:**
   ```bash
   make dev              # Start both backend and dashboard
   # or run separately:
   make serve-dev        # Terminal 1: Start backend
   make dashboard-dev    # Terminal 2: Start dashboard
   ```

3. **Before committing:**
   ```bash
   make check            # Run all Go checks
   make test-dashboard   # Run dashboard tests
   make lint-dashboard   # Run dashboard linting
   ```

4. **Before pushing:**
   ```bash
   make ci               # Ensure CI will pass
   ```

### Building from Source

#### Full Build
```bash
git clone https://github.com/biancarosa/fetchr.sh.git
cd fetchr.sh
make build-all
```

#### Backend Only
```bash
make build
```

The Go binary will be created as `fetchr` in the project root.

#### Dashboard Only
```bash
make build-dashboard
```

The dashboard build will be created in `dashboard/dist/`.

### Running Tests

#### Unit Tests (Go)
```bash
make test
```

#### Dashboard Tests
```bash
make test-dashboard
```

#### End-to-End Tests
```bash
make e2e
```

#### All Tests
```bash
make check  # Runs all Go tests plus linting
make test-dashboard  # Run dashboard tests separately
```

### Code Quality

#### Go Code Quality
The project uses golangci-lint for Go code quality checks:

```bash
make lint
```

#### Dashboard Code Quality
The dashboard uses ESLint and Prettier:

```bash
make lint-dashboard
```

To see the Go coverage report in your browser:

```bash
make coverage
```

### Dashboard Development

The dashboard is built with modern web technologies:
- **Framework**: React with TypeScript
- **Styling**: Tailwind CSS
- **State Management**: Zustand
- **API Client**: Axios with proper TypeScript types
- **Testing**: Jest and React Testing Library
- **Build Tool**: Vite

#### Dashboard Structure
```
dashboard/
├── src/
│   ├── components/     # Reusable UI components
│   ├── pages/         # Page components
│   ├── services/      # API services
│   ├── types/         # TypeScript type definitions
│   ├── hooks/         # Custom React hooks
│   ├── utils/         # Utility functions
│   └── store/         # State management
├── public/            # Static assets
└── tests/            # Test files
```

#### Dashboard Commands
All dashboard commands should be run from the `dashboard/` directory:

```bash
npm install --legacy-peer-deps    # Install dependencies
npm run dev                       # Start development server
npm run build                     # Build for production
npm run test                      # Run tests
npm run lint                      # Run linting
npm run preview                   # Preview production build
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Guidelines
- Follow the cursor rules in `.cursorrules`
- Use conventional commits
- Write tests for new features
- Update documentation
- Run `make check` before submitting PRs
- For dashboard changes, also run `make test-dashboard` and `make lint-dashboard`

## License

This project is licensed under the MIT License - see the LICENSE file for details.
