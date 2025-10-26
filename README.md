![# kingdom-auth](./assets/banner.png)

## About - a manifesto
> Authentication has gotten too complicated. I can't be the only one thinking that. I don't need all those features bloating my dev-experience. Yes, for big projects they are well worth it. But for my side-projects? Sometimes, I just want an easy to set up service that supports OAuth and allows my other services in the project to authenticate my users reliably.  
> ~ me, just now

kingdom-auth is that.

kingdom-auth is a minimal authentication backend written in go (golang). It aims to provide the simplest dev-experience possible without sacrificing on user experience or security.

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

