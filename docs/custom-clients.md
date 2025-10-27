# Custom kingdom-auth implementations

## Clientside
kingdom-auth is designed to be flexible and can be integrated into custom clients beyond the provided SDKs. Below are some guidelines for building your own client implementations.

### Authentication Flow

1. Get a list of available OAuth providers from the kingdom-auth service: `GET /providers` ->
    `{ "providers": [ "provider_name1", "provider_name2", ... ] }`
2. Open a new window to start the auth-flow. Let it visit `GET /auth/begin/{provider_name}`
   This will redirect the user to the provider's login page, which will then redirect back to kingdom-auth. kingdom-auth will conclude with a simple page that should automatically close the window via browser APIs. Once the window is closed, you can start checking for the auth result.
   If your client doesn't do windows (or doesn't allow closing them with javascript) you can also look for the redirect url /auth/end. If the page is fully loaded with this URL, the auth flow is complete.
   If you set up a custom redirect URL in your kingdom-auth configuration... You know your app best.

### Tokens
After a successful authentication, kingdom-auth will issue a Refresh Token. This token is saved as a cookie and shall not be directly used with your services.

Instead, the secondary "Auth Token" should be used for service authentication. This token is short-lived and can be refreshed using the Refresh Token.

### Refreshing Auth Tokens
To refresh an Auth Token, make a request to the kingdom-auth service: `GET /token/refresh` with the cookie containing the Refresh Token. The response will include a new Auth Token:
```json
{
  "exp": 1700000000,
  "token": "ey..."
}
```

The exp field indicates the expiration time of the new Auth Token (in Unix Seconds).

To complete your implementation, you should refresh the used Auth Token before it expires. It will usually expire within a few minutes. You can use the exp field to determine when to refresh it.

When refreshing the auth token, the refresh token will also be rotated if it's close to expiration. Make sure to update the stored refresh token cookie accordingly if your client doesn't automatically handle cookies.
