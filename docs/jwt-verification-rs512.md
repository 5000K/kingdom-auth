# RS512 JWT Verification - Quick Reference

## For Service Developers

If you're building a service that needs to verify kingdom-auth tokens, you need the **public key**.

### Getting the Public Key

The public key (typically `public_key.pem`) is generated when following our [Setup Guide](setup.md). It's safe to copy and distribute.

### Verifying Tokens
All service-focused libraries for kingdom-auth have validation helpers built in too, but even without them:

Verifying a JWT is a well documented problem in various languages. This is probably a good example for when you should just ask an LLM.

## Token Claims

The Auth tokens contain the following claims:

### Refresh Token Claims:
- `sub` - User ID (as string)
- `iss` - Issuer (configured in kingdom-auth)
- `exp` - Expiration time (Unix timestamp)
- `iat` - Issued at time (Unix timestamp)
- `aud` - Audience (can be customized per user via an authorized service)
- `public-data` - User's public data (JSON string). The user can't edit this, but the service can.
- `kaver` - Version of kingdom-auth (**K**ingdom **A**uth **Ver**sion; used to handle breaking changes
