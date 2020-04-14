# Handlers

This package contains handlers to forward values coming from a request header or a cookie, to a request context.

The mapping is done using enumeration of possible header and cookie keys, and their corresponding context key.

## Usage

- Read a header from a request and add its value to the output request context:

```go
    // Access token
    handler := handlers.CheckHeader(h, handlers.UserAccessHeaderKey)

    // Locale
    handler := handlers.CheckHeader(h, handlers.LocaleHeaderKey)

    // CollectionID
    handler := handlers.CheckHeader(h, handlers.CollectionIDHeaderKey)
```

- Read a cookie and add its value to the output request context:

```go
    // Access token
    handler := handlers.CheckCookie(h, handlers.UserAccessCookieKey)

    // Locale
    handler := handlers.CheckCookie(h, handlers.LocaleCookieKey)

    // CollectionID
    handler := handlers.CheckCookie(h, handlers.CollectionIDCookieKey)
```
