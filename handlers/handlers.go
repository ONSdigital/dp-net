package handlers

import (
	"context"
	"net/http"

	"github.com/ONSdigital/log.go/log"
)

// CheckHeaderMiddleware returns a CheckHeader function for the provided key as a middleware function
// (alice middleware Constructor, ie. func(http.Handler) -> (http.Handler))
func CheckHeaderMiddleware(key Key) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return CheckHeader(h, key)
	}
}

// CheckHeader is a wrapper which adds a value from the request header to context if one does not yet exist.
func CheckHeader(h http.Handler, key Key) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if headerValue := req.Header.Get(key.Header()); headerValue != "" {
			req = req.WithContext(context.WithValue(req.Context(), key.Context(), headerValue))
		}

		h.ServeHTTP(w, req)
	})
}

// CheckCookieMiddleware returns a CheckCookie function for the provided key as a middleware function
// (alice middleware Constructor, ie. func(http.Handler) -> (http.Handler))
func CheckCookieMiddleware(key Key) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return CheckCookie(h, key)
	}
}

// CheckCookie is a wrapper which adds a cookie value to context if one does not yet exist
func CheckCookie(h http.Handler, key Key) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		cookieValue, err := req.Cookie(key.Cookie())
		if err != nil {
			if err != http.ErrNoCookie {
				log.Event(req.Context(), "unexpected error while extracting value from cookie", log.ERROR, log.Error(err),
					log.Data{"cookie_key": key.Cookie()})
			}
		} else {
			req = req.WithContext(context.WithValue(req.Context(), key.Context(), cookieValue.Value))
		}

		h.ServeHTTP(w, req)
	})
}
