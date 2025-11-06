# Quick Start Guide

This is a quick reference for running kingdom-auth. For full documentation, see [README.md](README.md).

## Prerequisites

1. Generate RSA keys for JWT signing:
```bash
openssl genrsa -out private_key.pem 4096
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

2. Copy the example config and customize it:
```bash
cp config.example.yml config.yml
```

3. Edit `config.yml` to set your OAuth providers, cookie domain, CORS origins, database settings, token lifetimes, and service ports/URLs.

4. Startup kingdom-auth.

## Option 1: Run Locally

```bash
# Build
go build -o kingdom-auth ./main

# Run
export CONFIG_PATH=config.yml
./kingdom-auth
```

## Option 2: Docker

The repository provides two example docker-compose files you can base your setup on.

```bash
# Ready to use, with SQLite
docker-compose -f docker-compose.sqlite.yml up -d

# Ready to use, with Postgres
docker-compose -f docker-compose.postgresql.yml up -d
```

## Access

- **Main Service (User API)**: http://localhost:14414
- **System Service (Admin API)**: http://localhost:14415
