# Load Tester

A simple concurrent HTTP load tester that hammers a URL multiple times and reports performance metrics.

## Usage

```bash
# Default: 10 requests
go run main.go <url>

# Custom number of requests
go run main.go <url> <iterations>
```

## Example

```bash
go run main.go https://example.com 50
```

## Output

Displays:
- Total execution time for all requests
- Per-request status and duration
- Average response time
- P95 latency

All requests run in parallel using goroutines.
