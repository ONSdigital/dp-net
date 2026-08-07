# How to Use the Links Middleware

## Purpose of the Links Middleware

Many of our services include links in their code e.g. links that appear in the json or html responses. Quite often these services can be run in different environments with different protocols, hosts, and path prefixes (e.g. version numbers). For example, a 'Self' link for a redirect, in the Redirect API service, when run locally, needs to look something like this:

```json
"http://localhost:29900/v1/redirects/some_redirect_id_value"
```

But when run in the api beta domain, it needs to look something like this instead:

```json
"https://api.beta.ons.gov.uk/v1/redirects/some_redirect_id_value"
```

The links middleware, in a request handler, will use the request header information to determine the correct protocol, host, and path prefix values to put in any link. It will then create a builder that can be used to create any link by passing it a relative path. This means that the service can potentially be run in any environment and will automatically give the correct full link value.

Links created this way, which will then be written differently in different environments, are known as HATEOAS links.

## Example of Creating a HATEOAS Link

A pointer to the request header, and the value of the default API URL (such as `"http://localhost:29900"`) are passed into the `links.FromHeadersOrDefault` function, which returns a `links.Builder`. E.g.

```linkBuilder := links.FromHeadersOrDefault(&req.Header, "http://localhost:29900")```

If the request header contains values for the protocol and host then these will be used along with any path prefix provided. However, if these values are not present then the API URL value will be used instead, by default.

In either case, the self link value can then be created by passing its relative path value into the BuildLink function of the `links.Builder` E.g.

```redirectSelf, err := linkBuilder.BuildLink("/v1/redirects/some_redirect_id_value")```

NB. The functionality is often used for rewriting existing links as HATEAOS links.
