# Rate Limiter

A simple and configurable rate limiter for Go APIs. This project demonstrates how to apply rate limiting to different domains and actions (like messaging or authentication) using a YAML-based configuration.

## Features

- Per-domain and per-descriptor rate limiting
- Configurable via YAML file
- Middleware for easy integration with Go HTTP servers

## How It Works

The rate limiter reads rules from a `config.yaml` file, where you can specify different rate limits for different domains and descriptors (e.g., message type, auth type). Each rule defines:

- **Domain**: The logical area (e.g., `messaging`, `auth`)
- **Descriptors**: Key-value pairs to further specify the action (e.g., `message_type: marketing`)
- **Rate Limit**: The allowed number of requests per time unit (e.g., 5 requests per minute)

When a request comes in, the middleware checks the domain and descriptors, matches them to a rule, and enforces the specified rate limit.

## Algorithm Used

This rate limiter uses the **Token Bucket** algorithm:

- Each rule has a bucket with a fixed capacity (number of tokens = allowed requests per unit)
- Tokens are refilled at a constant rate (e.g., every minute)
- Each request consumes a token; if no tokens are left, the request is rejected (rate limited)

This approach allows for short bursts of traffic while enforcing an average rate over time.

## Configuration Example

See `config.example.yaml` for a sample configuration:

```yaml
rate_limits:
  - domain: messaging
    descriptors:
      message_type: marketing
    rate_limit:
      unit: minute
      requests_per_unit: 5
  - domain: auth
    descriptors:
      auth_type: login
    rate_limit:
      unit: minute
      requests_per_unit: 5
```

## Usage

1. **Copy the example config:**
   ```sh
   cp config.example.yaml config.yaml
   # Edit config.yaml as needed
   ```
2. **Build and run the server:**
   ```sh
   go build -o rate-limiter
   ./rate-limiter
   ```
3. **Integrate the middleware:**
   Import and use the middleware in your Go HTTP server to protect endpoints.

## License

MIT
