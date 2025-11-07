![# kingdom-auth](./assets/banner.png)

[![Go](https://github.com/5000K/kingdom-auth/actions/workflows/go.yml/badge.svg)](https://github.com/5000K/kingdom-auth/actions/workflows/go.yml) [![Docker](https://github.com/5000K/kingdom-auth/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/5000K/kingdom-auth/actions/workflows/docker-publish.yml)

> ! the current version of kingdom-auth is a open preview. API and features are not complete or final at this point in time. Don't use it for important projects yet.

## About - a manifesto
kingdom-auth is a minimal authentication backend written in go (golang). It aims to provide the simplest dev-experience possible without sacrificing on user experience or security. Selfhosted. No lock-in. If you can scale postgres, you can scale kingdom-auth.

The single core mission of kingdom-auth is to provide a seamless, quick to set up and integrate way to authenticate your users. It's JWTs are simple to validate in your own services, it's slim typescript library is quickly integrated into your frontend. Even for things beyond that: implementing a usable kingdom-auth client consists of implementing three endpoints and doing some basic, scheduled token-refreshing in the background.

### Understanding kingdom-auths minimalism
"Minimal authentication backend" could mean a lot of things. In this case, it means the following:

1. Everything that has to be done can be done in exactly one good way.
2. **All** things necessary to fulfill the main task without sacrifices are present.
3. **Nothing beyond** these things is present.

There are tradeoffs - not every project needs basic auth and nothing beyond. To be upfront with them:

1. kingdom-auth relies on OAuth to authenticate users. No username/email/phone-number based authentication. Because that would open a whole new book with verifying mails or phone numbers, dealing with reset-links, TFA, ... OAuth is easy to set up, everyone has accounts for the big providers nowadays.
2. kingdom-auth uses JWTs. Two to be exact. Read more about them below.


### The two tokens
kingdom-auth uses two tokens, commonly named Refresh-Token and Auth-Token. The Refresh-Token is only used to generate the Auth-Token. This is a pretty common pattern.
The reason behind this is pretty simple too: The Refresh-Token lives in the users browser as a cookie. But as an http-only cookie, it might (and usually will) not be able to be sent to your other services (if kingdom-auth lives on its own domain).
To deal with this, the kingdom-auth typescript client will send a request to your kingdom-auth instance, to generate a short-lived JWT (merely valid for a minute). This JWT now lives within your browsers memory.
Usually, we don't want this to happen to avoid attack surfaces opened by XSS. This is why this token only lives 60 seconds by default. The kingdom-auth typescript client automatically refreshes the token in the background before it is invalidated.

## Configuration

kingdom-auth is configured via a YAML configuration file combined with environment variables. The path to the config file is set via the `CONFIG_PATH` environment variable (defaults to `config.yml`).

### Configuration File

Create a `config.yml` file with your settings:

```yaml
# Cookie settings
cookie_name: katok
cookie_domain: localhost  # Change to your domain in production

# CORS - Allowed origins for cross-origin requests
allowed_origins:
  - http://localhost:5173
  - https://yourdomain.com

# Database configuration
db:
  type: sqlite  # Options: sqlite, mysql, postgres
  dsn: kingdom-auth.db  # For SQLite: filepath; for others: connection string

# OAuth providers - Add your OAuth providers here
providers:
  - name: github
    url: https://github.com/
    client_id: your_github_client_id
    client_secret: your_github_client_secret
    scopes:
      - user:email
  
  - name: forgejo
    url: https://forgejo.example.com/
    client_id: your_forgejo_client_id
    client_secret: your_forgejo_client_secret

# Token configuration
token:
  # RSA key paths for JWT signing (RS512)
  private_key_path: private_key.pem
  public_key_path: public_key.pem
  
  # Token lifetimes (in seconds)
  refresh_token_ttl: 864000  # 10 days
  auth_token_ttl: 90  # 1.5 minutes
  
  # JWT settings
  issuer: kingdom-auth
  default_audience: default-audience

# Main service configuration
main_service:
  port: 14414
  public_url: http://localhost:14414

# System service (admin API)
system_service:
  port: 14415
  tokens:
    - name: admin
      token: your_secure_system_token_here
```

### Environment Variables

Some configuration options can also be set with environment variables.
Make sure to only set an option either in the config file or via environment variable to avoid conflicts.

| Environment Variable     | Description                                  | Default                  |
|--------------------------|----------------------------------------------|--------------------------|
| `CONFIG_PATH`            | Path to the config file                      | `config.yml`             |
| `COOKIE_NAME`            | Name of the authentication cookie            | `katok`                  |
| `COOKIE_DOMAIN`          | Domain for the authentication cookie         | `localhost`              |
| `DB_TYPE`                | Database type (sqlite, mysql, postgres)      | `sqlite`                 |
| `DB_DSN`                 | Database connection string                   | `kingdom-auth.db`        |
| `DB_RUN_MIGRATIONS`      | Automatically run migrations                 | `true`                   |
| `PRIVATE_KEY_PATH`       | Path to RSA private key for JWT signing      | `private_key.pem`        |
| `PUBLIC_KEY_PATH`        | Path to RSA public key for JWT verification  | `public_key.pem`         |
| `REFRESH_TOKEN_TTL`      | Refresh token lifetime in seconds            | `864000` (10 days)       |
| `AUTH_TOKEN_TTL`         | Auth token lifetime in seconds               | `90` (1.5 min)           |
| `JWT_ISSUER`             | JWT issuer claim                             | `kingdom-auth`           |
| `JWT_DEFAULT_AUDIENCE`   | Default JWT audience claim                   | `default-audience`       |
| `MAIN_PORT`              | Main service port                            | `14414`                  |
| `MAIN_PUBLIC_URL`        | Public URL of the main service               | `http://localhost:14414` |
| `SYSTEM_PORT`            | System service port                          | `14415`                  |


### Examples: Database Connection Strings

#### SQLite
```yaml
db:
  type: sqlite
  dsn: kingdom-auth.db  # or /path/to/database.db
```

#### PostgreSQL
```yaml
db:
  type: postgres
  dsn: "host=localhost user=username password=password dbname=kingdom_auth port=5432 sslmode=disable"
```

#### MySQL
```yaml
db:
  type: mysql
  dsn: "username:password@tcp(localhost:3306)/kingdom_auth?charset=utf8mb4&parseTime=True&loc=Local"
```

### Example Files

The repository includes example configuration files to help you get started:

- **`config.example.yml`** - Complete configuration example with comments
- **`docker-compose.yml`** - Production setup with PostgreSQL
- **`docker-compose.dev.yml`** - Simple development setup with SQLite

Copy and customize these files for your deployment.
