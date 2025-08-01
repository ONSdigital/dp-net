# How to Use the Links Middleware

## Purpose of the Links Middleware

Many of our services include links in their code e.g. links that appear in the json or html responses. Quite often these services can be run in different environments with different protocols, hosts, and path prefixes (e.g. version numbers). For example, a 'Self' link for a redirect, in the Redirect API service, when run locally, needs to look something like this:

```json
"http://localhost:29900/v1/redirects/L2J1c2luZXNzaW5kdXN0cnlhbmR0cmFkZS9pdGFuZGludGVybmV0aW5kdXN0cnkvbXlidWxsZXRpbg=="
```

But when run in the api beta domain, it needs to look something like this instead:

```json
"https://api.beta.ons.gov.uk/v1/redirects/L2J1c2luZXNzaW5kdXN0cnlhbmR0cmFkZS9pdGFuZGludGVybmV0aW5kdXN0cnkvbXlidWxsZXRpbg=="
```

The links middleware allows you to pass it any link, when handling a request, and it will use the request header information to determine the correct protocol, host, and path prefix values to put in the link. It will then rewrite the link using those correct values. This means that the service can potentially be run in any environment and will automatically give the correct link value.

## Example of Rewriting a Link

In the case of the self link above, the original (local) link is created by the service using default values of protocol and host, which are given by its default API URL:

 `"http://localhost:29900"`

NB. It may not always be desirable for links to be rewritten. Therefore a config setting, 'EnableURLRewriting', is often used to switch the functionality on or off.

If 'EnableURLRewriting' is switched on then the following things happen next:

A pointer to the request header, and the value of the default API URL (such as the one specified above) are passed into the `links.FromHeadersOrDefault` function, which returns a `links.Builder`. E.g.

```redirectLinkBuilder := links.FromHeadersOrDefault(&req.Header, api.apiUrl)```

The self link value can then be rewritten by passing its original value into the BuildLink function of the `links.Builder` E.g.

```newSelf, err = redirectLinkBuilder.BuildLink(originalSelf)```
