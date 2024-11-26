package links

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Middleware struct {
	handler         http.Handler
	DefaultProtocol string
	DefaultHost     string
	DefaultPort     string
}

func NewMiddleWare(defaultURL string) (func(http.Handler) http.Handler, error) {
	parsedURL, err := url.Parse(defaultURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}

	defaultProtocol := parsedURL.Scheme
	defaultHost := parsedURL.Hostname()
	defaultPort := parsedURL.Port()

	return func(h http.Handler) http.Handler {
		return &Middleware{
			DefaultProtocol: defaultProtocol,
			DefaultHost:     defaultHost,
			DefaultPort:     defaultPort,
			handler:         h,
		}
	}, nil
}

func (l *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forwardedProto, _ := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Proto", l.DefaultProtocol)
	forwardedHost, forwardedHostFound := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Host", l.DefaultHost)
	forwardedPort := l.DefaultPort
	if forwardedHostFound {
		forwardedPort, _ = getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Port", "")
	}
	forwardedPathPrefix, _ := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Path-Prefix", "")

	ctx := context.WithValue(r.Context(), ctxProtocol, forwardedProto)
	ctx = context.WithValue(ctx, ctxHost, forwardedHost)
	ctx = context.WithValue(ctx, ctxPort, forwardedPort)
	ctx = context.WithValue(ctx, ctxPathPrefix, forwardedPathPrefix)

	r2 := r.WithContext(ctx)

	l.handler.ServeHTTP(w, r2)
}

func getForwardedHeaderElseDefault(header http.Header, key, defaultValue string) (value string, found bool) {
	value = header.Get(key)
	if value == "" {
		fmt.Printf("\n%s value not found, using default: %s\n", key, defaultValue)
		return defaultValue, false
	}
	fmt.Printf("\n%s value found, using forwarded: %s\n", key, value)
	return value, true
}
