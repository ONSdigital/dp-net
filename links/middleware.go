package links

import (
	"context"
	"net/http"
	"net/url"
)

type Middleware struct {
	defaultURL *url.URL
	handler    http.Handler
}

func NewMiddleWare(defaultURL *url.URL) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Middleware{
			defaultURL: defaultURL,
		}
	}
}

// HandlerFunc returns a handlerFunc that can be used for individual routes so that the middleware doesn't have to be used globally for all routes.
func HandlerFunc(defaultURL *url.URL, handlerFunc http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := getAPIURLFromHeaderOrDefault(&r.Header, defaultURL)
		handlerFunc(w, r.WithContext(context.WithValue(r.Context(), ctxAPIURL, u)))
	}
}

func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := getAPIURLFromHeaderOrDefault(&r.Header, mw.defaultURL)
	mw.handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxAPIURL, u)))
}

func getAPIURLFromHeaderOrDefault(h *http.Header, defaultURL *url.URL) *url.URL {
	host := h.Get("X-Forwarded-Host")
	if host == "" {
		return defaultURL
	}

	scheme := h.Get("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}

	port := h.Get("X-Forwarded-Port")
	if port != "" {
		host += ":" + port
	}

	path := h.Get("X-Forwarded-Path-Prefix")

	url := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   "/",
	}
	return url.JoinPath(path)
}
