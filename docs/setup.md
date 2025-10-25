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
