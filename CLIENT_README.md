# Kingdom Auth Go Client

You should probably also read the [JWT verification reference](docs/jwt-verification-rs512.md) for more context.

A Go client for kingdom-auth. Targeted towards services that need to interact with kingdom-auth and validate JWT tokens it creates.

## Installation

```bash
go get github.com/5000K/kingdom-auth
```

## Usage

### Initialize the Client

```go
import (
    kingdomauth "github.com/5000K/kingdom-auth"
)

client, err := kingdomauth.NewClient(
    "https://auth.example.com",  // Kingdom Auth service URL
    "your-service-secret",        // Service authentication secret
    "./public_key.pem",           // Path to RSA public key
)
if err != nil {
    log.Fatal(err)
}
```

### Validate JWT Tokens
See [examples/client_usage.go](examples/client_usage.go) for a complete working example.
