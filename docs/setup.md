# Auth setup

Prerequisites: kingdom-auth heavily depends on oauth to function. It does not offer user/password authentication - by design.

Your oauth provider needs to support OIDC-connect (which all well implemented providers do nowadays).

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

# Key Setup

kingdom-auth uses RSA keys (RS512) to sign its JWTs. You need to generate a private and public key pair:

```bash
# Generate private key (4096-bit RSA)
openssl genrsa -out private_key.pem 4096

# Extract public key from private key
openssl rsa -in private_key.pem -pubout -out public_key.pem
```

**Important:** Keep your `private_key.pem` secure and never commit it to version control!  
The public key (`public_key.pem`) can be shared with services that need to verify tokens.

By default, kingdom-auth looks for these files in the project root directory. You can customize the paths in your config.yml:

```yml
token:
  private_key_path: /path/to/private_key.pem  # defaults to: private_key.pem
  public_key_path: /path/to/public_key.pem    # defaults to: public_key.pem
  # ... other token config
```
