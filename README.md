# fetchr.sh

A modern HTTP proxy and request capture tool with a sleek web dashboard.

## Features

- **HTTP Proxy Server**: Forward HTTP requests through the proxy
- **Request History**: Comprehensive tracking of all proxied requests
- **Web Dashboard**: Modern React-based interface for managing requests
- **Real-time Statistics**: Monitor proxy performance and request metrics
- **Admin API**: RESTful endpoints for health checks, metrics, and history
- **Request Replay**: Load and replay requests from history
- **Timing Metrics**: Detailed breakdown of proxy overhead and upstream latency

## Quick Start

### 1. Build the Project

```bash
make build
```

### 2. Start the Proxy Server

```bash
# Start with request history and admin endpoints
./fetchr serve --admin-port 8081 --history-size 1000

# Or use the development command
make dev
```

### 3. Access the Dashboard

```bash
# Build and serve the dashboard
cd dashboard
npm install --legacy-peer-deps
npm run build
npm run dev
```

Open http://localhost:3000 to access the dashboard.

## Dashboard Features

### Request Builder
- **HTTP Method Selection**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **URL Input**: Full URL with validation
- **Headers Management**: Dynamic header pairs with enable/disable toggles
- **Request Body**: JSON editor with syntax highlighting
- **Send Requests**: Execute requests through the proxy or directly

### Request History
The dashboard provides two types of request history:

#### Local History
- Requests made from the dashboard interface
- Stored in browser localStorage
- Includes request/response data and timing
- Replay functionality to load previous requests

#### Proxy History
- All HTTP requests that pass through the proxy server
- Stored in server memory (configurable size)
- Detailed timing metrics (proxy overhead, upstream latency)
- Request/response size tracking
- Success/error status tracking

### Statistics Panel
Real-time statistics showing:
- **Request Counts**: Success vs error rates
- **Timing Metrics**: Average durations and latency breakdown
- **Data Transfer**: Total request/response sizes
- **Status Codes**: Distribution of HTTP status codes
- **Methods**: Usage breakdown by HTTP method

## API Endpoints

### Proxy Server (Port 8080)
- **HTTP Proxy**: All HTTP traffic

### Admin Server (Port 8081)
- `GET /healthz` - Health check
- `GET /metrics` - Prometheus-style metrics
- `GET /requests` - Request history (JSON)
- `GET /requests/stats` - Request statistics
- `POST /requests/clear` - Clear request history

## Configuration

### Command Line Flags

```bash
./fetchr serve [flags]

Flags:
  --port int              Proxy server port (default 8080)
  --admin-port int        Admin server port (enables admin endpoints)
  --history-size int      Maximum requests to keep in history (default 1000)
  --log-level string      Log level: debug, info, warn, error (default "info")
```

### Environment Variables

Dashboard configuration (in `dashboard/.env.local`):
```bash
NEXT_PUBLIC_PROXY_HOST=localhost
NEXT_PUBLIC_PROXY_PORT=8080
NEXT_PUBLIC_ADMIN_PORT=8081
```

## Usage Examples

### Using as HTTP Proxy

```bash
# Configure your application to use the proxy
curl -x http://localhost:8080 http://httpbin.org/get

# Or set environment variables
export HTTP_PROXY=http://localhost:8080
export HTTPS_PROXY=http://localhost:8080
```

### Request History API

```bash
# Get request history
curl http://localhost:8081/requests

# Get statistics
curl http://localhost:8081/requests/stats

# Clear history
curl -X POST http://localhost:8081/requests/clear
```

## Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- npm

### Development Commands

```bash
# Install Go development tools
make install

# Run tests
make test
make e2e

# Run linting
make lint
make check

# Start development servers
make dev

# Build for production
make build
make build-dashboard
```

### Project Structure

```
fetchr.sh/
├── cmd/fetchr/          # Main application entry point
├── internal/
│   ├── proxy/           # Proxy server implementation
│   └── api/             # Admin API handlers
├── dashboard/           # React dashboard
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── services/    # API service layer
│   │   └── hooks/       # Custom hooks and state management
│   └── public/          # Static assets
└── Makefile            # Development commands
```

## Request History Details

### HTTP vs HTTPS Requests

- **HTTP Requests**: Fully captured with complete request/response data
- **HTTPS Requests**: Only CONNECT tunnel establishment is visible (encrypted content cannot be captured)

### Data Captured

For each HTTP request through the proxy:
- Request method, URL, headers, body
- Response status, headers, body
- Timing breakdown:
  - Proxy overhead (time spent in proxy code)
  - Upstream latency (time waiting for target server)
  - Total duration
- Data sizes (request and response bytes)
- Success/error status

### Performance Considerations

- Request history is stored in memory
- Configurable maximum size (default: 1000 requests)
- Automatic cleanup of oldest requests when limit reached
- Minimal performance impact on proxy operations

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check` to ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
