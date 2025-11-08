FROM golang:1.25.3 AS builder

WORKDIR /app

# Copy go modules and vendor directory for reproducible builds
COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

# Build with CGO enabled (required for sqlite) using vendored dependencies
RUN go build -mod=vendor -o kingdomauth ./main

FROM debian:bookworm-slim

# Install necessary runtime dependencies for CGO-enabled binaries
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/kingdomauth /kingdomauth

CMD ["/kingdomauth"]