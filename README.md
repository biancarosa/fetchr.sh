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

### Building from Source

```bash
git clone https://github.com/biancarosa/fetchr.sh.git
cd fetchr.sh
go build -o fetchr ./cmd/fetchr
```

### Running Tests

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
