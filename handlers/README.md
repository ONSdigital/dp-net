# Handlers

This package contains handlers to forward values coming from a request header or a cookie, to a request context.

The mapping is done using enumeration of possible keys, and their mappings to header, context or cookie keys.

## Usage

- Read a header from a request and add its value to the underlying request context:

```go
    // Access token
    handler := handlers.CheckHeader(handlers.UserAccess)

    // Locale
    handler := handlers.CheckHeader(handlers.Locale)

    // CollectionID
    handler := handlers.CheckHeader(handlers.CollectionID)

    // RequestID
    handler := handlers.CheckHeader(handlers.RequestID)

    // User Identity
    handler := handlers.CheckHeader(handlers.UserIdentity)
```

- Read a cookie and add its value to the output request context:

```go
    // Access token
    handler := handlers.CheckCookie(handlers.UserAccess)

    // Locale
    handler := handlers.CheckCookie(handlers.Locale)

    // CollectionID
    handler := handlers.CheckCookie(handlers.CollectionID)
```

- Access the context value:

```go
    // Access token
    accessToken := ctx.Value(handlers.UserAccess.Context())

    // Locale
    locale := ctx.Value(handlers.Locale.Context())

    // CollectionID
    collectionID := ctx.Value(handlers.CollectionID.Context())

    // RequestID
    requestID := ctx.Value(handlers.RequestID.Context())

    // User Identity
    handler := ctx.Value(handlers.UserIdentity.Context())
```
