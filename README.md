# dp-net
Network library, containing an HTTP client and Server, handlers and other utilities for network communications.

## http

Http package contains a base http Server and Client to be used by all ONS digital publishing services that require HTTP communication.

The package also includes http utilities like constants, error definitions, a requestID handler required by the HTTP server, and validation utilities for identity, locale and models.

### Client

The dp-net/v2/http Client provides for robust contextual HTTP, and a default client
that inherits the methods associated with the standard HTTP client,
but with the addition of production-ready timeouts and context-sensitivity,
and the ability to perform exponential backoff when calling another HTTP server.

#### How to use

This client should have a familiar feel to it when it is used.
There are a few options available when you want to create a client
- Simple client with no timeout value (shown below)
- With timeout (both per request AND total timeout (spans all retries)) can be set
  - ClientWithTimeout (sets per try timeout duration)
  - ClientWithTotalTimeout (total client timeout)
  - ClientWithTimeouts (allows you to set both the above timeouts)

An example given below:

```go
import http "github.com/ONSdigital/dp-net/v2/http"

func httpHandlerFunc(w http.ResponseWriter, req *http.Request) {
    client := http.NewClient()

    resp, err := client.Get(req.Context(), "https://www.google.com")
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

In this case, in the unlikely event of https://www.google.com returning a status
of 500 or above, the client will retry at exponentially-increasing intervals, until
the max retries (10 by default is reached).

Also, if the inbound request is cancelled, for example, its context will be closed
and this will be noticed by the client.

You also do not have to use the default client if you don't like the configured
timeouts or do not wish to use exponential backoff. The following example shows
how to configure your own dp-net/v2/http client:

```go
import (
    "net/http"
    dphttp "github.com/ONSdigital/dp-net/v2/http"
)

func main() {
    rcClient := &dphttp.Client{
        // MaxRetries is the maximum number of retries you wish to
        // wait for, the retry method implements exponential backoff
        MaxRetries:         10,
        // RetryTime is the gap before (any) first retry (increases for second retry, and so on)
        RetryTime:          1 * time.Second,
        // PathsWithNoRetries is a list of all paths that you do not wish to retry call on failure,
        // the path should be set to true (default client has empty map)
        PathsWithNoRetries: map[string]bool{
            "/health": true,
        },
        // Create your own http client with configured timeouts
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
            Transport: &http.Transport{
                DialContext: (&net.Dialer{
                    Timeout: 5 * time.Second,
                }).DialContext,
                TLSHandshakeTimeout: 5 * time.Second,
                MaxIdleConns:        10,
                IdleConnTimeout:     30 * time.Second,
            },
        },
    }
}
```

### Server

The Server extends the default golang HTTP Server by adding a requestID and logger middleware. By default it handles the OSSignals, and it has a default shutdown timeout of 10 seconds.

This Server is intended to be used by all ONS digital publishing services that require to serve HTTP. The following example shows how to use the Server:

#### Creation

You have 2 options available (depending on if you want to specify a request timeout)
 - NewServer(bindAddr string, router http.Handler)
 - NewServerWithTimeout(bindAddr string, router http.Handler, timeout time.Duration, timeoutMessage string)
 
Assuming you have created a router with your API handlers, you can create the http server like so:

```go
import http "github.com/ONSdigital/dp-net/v2/http"
    ...
    httpServer := http.NewServer(bindAddr, router)
    httpServer.HandleOSSignals = false
    ...
```
Note that HandleOSSignal is set to false, so that the main thread will be responsible to shutdown the server during graceful shutdown.

#### Start

Start the server in a new go-routine, because this operation is blocking:
```go
    ...
    go func() {
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Error(ctx, "error starting http server", err)
            return
        }
    }()
    ...
```
Note that we ignore ErrServerClosed, because this is a valid scenario during graceful shutdown.

#### Shutdown

Shutdown the server when you no longer require it. Usually you will need to do this as part of the service graceful shutdown, after receiving a SIGINT or SIGTERM system call in your signal channel:
```go
    ...
    err := httpServer.Shutdown(shutdownCtx)
    if err != nil {
        log.Error(shutdownCtx, "http server shutdown error", err)
    } else {
        log.Info(shutdownCtx, "http server successful shutdown")
    }
    ...
```

## Handlers

This module includes handlers for accessToken, collectionID, localeCode, and finally a JSON response writer and a Proxy creation utility.