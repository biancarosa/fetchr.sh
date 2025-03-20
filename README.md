# fetchr.sh

A powerful HTTP request proxy written in Go that provides advanced request management capabilities including caching, throttling, rate limiting, and asynchronous support.

## Overview

fetchr.sh is designed to be a robust sidecar proxy for your applications, handling HTTP requests with built-in performance and reliability features. It supports multiple API endpoints in a single container, allowing you to manage and optimize requests to various services through a unified interface. It's particularly useful for applications that need to manage multiple API calls while maintaining control over request patterns and resource usage.

## Features

- **Multi-API Support**: Handle multiple API endpoints in a single container
- **Caching**: Intelligent response caching to reduce redundant API calls
- **Throttling**: Control request rates to prevent overwhelming target services
- **Rate Limiting**: Configurable rate limits per endpoint or globally
- **Asynchronous Support**: Handle requests concurrently with controlled concurrency
- **Request Logging**: Comprehensive request/response logging with detailed metadata
- **Request Replay**: Ability to replay captured requests for testing and debugging
- **Sidecar Architecture**: Easy integration with existing applications
- **Cloud-Ready**: Designed for cloud environments (coming soon)

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker (optional, for containerized deployment)

### Installation

```bash
go install github.com/yourusername/fetchr.sh@latest
```

### Basic Usage

1. Start fetchr.sh as a sidecar to your application:

```bash
fetchr.sh serve --port 8080
```

2. Make requests through fetchr.sh using the CLI:

```bash
# Simple GET request
fetchr.sh request --target https://api.example.com/users

# POST request with body
fetchr.sh request \
  --target https://api.example.com/users \
  --method POST \
  --headers "Content-Type=application/json" \
  --body '{"name": "John Doe"}'

# Asynchronous request
fetchr.sh request \
  --target https://api.example.com/process \
  --method POST \
  --async \
  --output result.json
```

3. Configure your application to make requests through fetchr.sh instead of directly to the target services.

### Configuration

fetchr.sh can be configured through environment variables or a configuration file:

```yaml
port: 8080
cache:
  enabled: true
  ttl: 5m
rate_limit:
  requests_per_second: 100
throttling:
  enabled: true
  max_concurrent_requests: 50
logging:
  enabled: true
  level: debug
  format: json
  output: logs/fetchr.log
  retention: 7d
replay:
  enabled: true
  storage: ./replay-data
  max_size: 1GB
```

## Request Logging and Replay

### Logging Capabilities

fetchr.sh provides comprehensive request logging that captures:
- Request headers and body
- Response headers and body
- Timing information
- Cache hits/misses
- Rate limit status
- Error details
- Request metadata (client IP, user agent, etc.)
- Target service information

Logs can be output in various formats (JSON, structured text) and can be configured for different retention periods.

### Request Replay

The replay feature allows you to:
- Capture and store requests for later analysis
- Replay requests with original timing and headers
- Modify requests before replay
- Compare responses between original and replayed requests
- Create test scenarios from production traffic

Example capture and replay commands:
```bash
# Capture requests
fetchr.sh capture \
  --output ./captures \
  --filter "method=POST" \
  --time-range "2024-03-20T00:00:00Z/2024-03-21T00:00:00Z"

# Replay captured requests
fetchr.sh replay \
  --from ./captures \
  --target https://api.example.com \
  --concurrent 5 \
  --output ./replay-results
```

## Cache Management

Manage the cache using the following commands:

```bash
# Clear cache for specific target
fetchr.sh cache clear --target https://api.example.com

# Show cache statistics
fetchr.sh cache stats --format json

# Invalidate specific cache entries
fetchr.sh cache invalidate --pattern "users/*"
```

## Monitoring and Statistics

View statistics and metrics:

```bash
# Show general statistics
fetchr.sh stats

# Show detailed metrics for a specific target
fetchr.sh stats --target https://api.example.com --metrics

# Export statistics to file
fetchr.sh stats --export stats.json
```

## Architecture

fetchr.sh operates as a reverse proxy, intercepting requests and applying configured policies before forwarding them to the specified target services. It maintains its own connection pool and implements circuit breaking to prevent cascading failures. The proxy supports multiple target services simultaneously, with independent rate limiting, caching, and monitoring for each service.

## Roadmap

- [ ] Cloud provider integrations (AWS, GCP, Azure)
- [ ] Metrics and monitoring
- [ ] Advanced caching strategies
- [ ] Request/Response transformation
- [ ] Authentication and authorization

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers.
