# Auth setup

Prerequisites: kingdom-auth heavily depends on oauth to function. It does not offer user/password authentication - by design.

Your oauth provider needs to support OIDC-connect (which all well implemented providers do nowadays).

## Step 1: Configure OAuth Providers

You'll simply add a list of providers to your config.yml:

```yml
# ... other config options

providers:
  - name: provider_name1           # technical identifier you'll use everywhere when dealing with it. needs to be URL-valid 
    url: https://auth.example.com  # the base url. Needed endpoints will be detected via OIDC-connect
    client_id: ...                 # your client_id
    client_secret: ...             # your client_secret
  - name: provider_name2           
    url: https://auth2.example.com
    client_id: ...
    client_secret: ...
    
# ... other config options
```

You need to register at least one OAuth provider before you can use kingdom-auth.

## Step 2: Generate RSA Keys

kingdom-auth uses 4096-bit RSA keys to sign its JWTs. You need to generate a private and public key pair:

```bash
# Generate private key (4096-bit RSA)
openssl genrsa -out private_key.pem 4096

# Extract public key from private key
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

**Important:** Keep your `private_key.pem` secure and never commit it to version control!  
The public key (`public_key.pem`) can and should be shared with services that need to verify tokens.

By default, kingdom-auth looks for these files in the project root directory. You can customize the paths in your config.yml:

```yml
token:
  private_key_path: /path/to/private_key.pem  # defaults to: private_key.pem
  public_key_path: /path/to/public_key.pem    # defaults to: public_key.pem
  # ... other token config
```

## Step 3: Create Your Configuration File

Copy the example configuration and customize it:

```bash
cp config.example.yml config.yml
```

Edit `config.yml` to set your:
- OAuth providers (from Step 1)
- Cookie domain
- Allowed CORS origins
- Database settings
- Token lifetimes
- Service ports and URLs

The configuration file path can be set via the `CONFIG_PATH` environment variable:

```bash
export CONFIG_PATH=/path/to/config.yml
```

## Step 4: Choose Your Database

kingdom-auth supports three database types:

### SQLite
```yaml
db:
  type: sqlite
  dsn: kingdom-auth.db
```

### Postgres
```yaml
db:
  type: postgres
  dsn: "host=localhost user=kingdom_auth password=securepass dbname=kingdom_auth port=5432 sslmode=disable"
```

### MySQL
```yaml
db:
  type: mysql
  dsn: "username:password@tcp(localhost:3306)/kingdom_auth?charset=utf8mb4&parseTime=True&loc=Local"
```

## Step 5: Run kingdom-auth

See [QUICKSTART.md](../QUICKSTART.md) for detailed run instructions, or see the [main README](../README.md) for complete documentation.
