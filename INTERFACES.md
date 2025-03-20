I# fetchr.sh Interfaces

This document outlines the interfaces and commands for the fetchr.sh CLI tool.

## Core Commands

### `fetchr.sh serve`

Starts the fetchr.sh proxy server.

```bash
fetchr.sh serve [flags]
```

Flags:
- `--port, -p`: Port to listen on (default: 8080)
- `--config, -c`: Path to configuration file
- `--log-level`: Logging level (debug, info, warn, error)
- `--metrics`: Enable metrics endpoint
- `--health`: Enable health check endpoint

### `fetchr.sh request`

Makes a request through the proxy.

```bash
fetchr.sh request [flags]
```

Flags:
- `--target, -t`: Target service URL (required)
- `--method, -m`: HTTP method (default: GET)
- `--headers, -H`: HTTP headers (key=value)
- `--body, -b`: Request body
- `--timeout`: Request timeout
- `--no-cache`: Disable caching for this request
- `--async`: Make request asynchronously
- `--output, -o`: Output file for response

### `fetchr.sh replay`

Replays captured requests.

```bash
fetchr.sh replay [flags]
```

Flags:
- `--from`: Directory containing replay data
- `--target`: Target service URL
- `--filter`: Filter requests (method=GET, status=200, etc.)
- `--time-range`: Time range for replay
- `--concurrent`: Number of concurrent replays
- `--modify`: Modify requests before replay
- `--output`: Output directory for replay results

### `fetchr.sh capture`

Captures requests for later replay.

```bash
fetchr.sh capture [flags]
```

Flags:
- `--output`: Output directory for captured requests
- `--filter`: Filter requests to capture
- `--time-range`: Time range for capture
- `--max-size`: Maximum size of capture data
- `--compress`: Compress captured data

### `fetchr.sh stats`

Shows statistics and metrics.

```bash
fetchr.sh stats [flags]
```

Flags:
- `--target`: Filter stats by target service
- `--time-range`: Time range for stats
- `--format`: Output format (text, json)
- `--metrics`: Show detailed metrics
- `--export`: Export stats to file

### `fetchr.sh config`

Manages configuration.

```bash
fetchr.sh config [flags] <command>
```

Commands:
- `show`: Show current configuration
- `edit`: Edit configuration file
- `validate`: Validate configuration
- `reset`: Reset to default configuration

### `fetchr.sh cache`

Manages the cache.

```bash
fetchr.sh cache [flags] <command>
```

Commands:
- `clear`: Clear the cache
- `invalidate`: Invalidate specific cache entries
- `stats`: Show cache statistics
- `export`: Export cache contents
- `import`: Import cache contents

## Global Flags

These flags are available for all commands:

- `--verbose, -v`: Enable verbose output
- `--debug`: Enable debug mode
- `--quiet, -q`: Suppress output
- `--help, -h`: Show help
- `--version`: Show version

## Examples

### Making a Request

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

### Capturing and Replaying Requests

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

### Managing Cache

```bash
# Clear cache for specific target
fetchr.sh cache clear --target https://api.example.com

# Show cache statistics
fetchr.sh cache stats --format json
```

## Future Enhancements

1. **Interactive Mode**
   - Add an interactive shell mode for easier command execution
   - Support command history and auto-completion

2. **Advanced Filtering**
   - Add support for complex query filters
   - Support regular expressions in filters
   - Add support for custom filter plugins

3. **Request Templates**
   - Support for saving and reusing request templates
   - Template variables and substitution
   - Template versioning

4. **Batch Operations**
   - Support for batch request processing
   - Parallel request execution
   - Batch result aggregation

5. **Integration Features**
   - Export/import in various formats (HAR, Postman, etc.)
   - Integration with popular API testing tools
   - Support for CI/CD pipelines

6. **Advanced Analytics**
   - Request/response analysis
   - Performance metrics
   - Error pattern detection
   - Custom metric collection 