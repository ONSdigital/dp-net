package links

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Middleware struct {
	handler           http.Handler
	DefaultProtocol   string
	DefaultHost       string
	DefaultPort       string
	DefaultUrlVersion string
}

func NewMiddleWare(defaultURL string) (func(http.Handler) http.Handler, error) {
	parsedURL, err := url.Parse(defaultURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	defaultProtocol := parsedURL.Scheme
	defaultHost := parsedURL.Hostname()
	defaultPort := parsedURL.Port()
	defaultUrlVersion := ""

	return func(h http.Handler) http.Handler {
		return &Middleware{
			DefaultProtocol:   defaultProtocol,
			DefaultHost:       defaultHost,
			DefaultPort:       defaultPort,
			DefaultUrlVersion: defaultUrlVersion,
			handler:           h,
		}
	}, nil
}

func (l *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	forwardedProto, _ := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Proto", l.DefaultProtocol)
	forwardedHost, forwardedHostFound := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Host", l.DefaultHost)
	forwardedPort, _ := getForwardedHeaderElseDefault(r.Header, "X-Forwarded-Port", l.DefaultPort)
	forwardedUrlVersion, _ := getForwardedHeaderElseDefault(r.Header, "urlVersion", l.DefaultUrlVersion)

	// No forwarded host so default port is required
	if !forwardedHostFound {
		forwardedPort = l.DefaultPort
	}

	ctx := context.WithValue(r.Context(), ctxProtocol, forwardedProto)
	ctx = context.WithValue(ctx, ctxHost, forwardedHost)
	ctx = context.WithValue(ctx, ctxPort, forwardedPort)
	ctx = context.WithValue(ctx, ctxUrlVersion, forwardedUrlVersion)

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
