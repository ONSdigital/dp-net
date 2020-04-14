package handlers

import (
	"context"
	"net/http"

	netHttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/log"
)

// HeaderKey - iota enum of possible header keys
type HeaderKey int

// CookieKey - iota enum of possible cookie keys
type CookieKey int

// Possible values for header keys
const (
	UserAccessHeaderKey HeaderKey = iota
	LocaleHeaderKey
	CollectionIDHeaderKey
)

// Possible values for cookie keys
const (
	UserAccessCookieKey CookieKey = iota
	LocaleCookieKey
	CollectionIDCookieKey
)

// Header Keys string representations
var headerKeys = map[HeaderKey]string{
	UserAccessHeaderKey:   netHttp.FlorenceHeaderKey,
	LocaleHeaderKey:       netHttp.LocaleHeaderKey,
	CollectionIDHeaderKey: netHttp.CollectionIDHeaderKey,
}

// Context Keys corresponding to each header key
var headerContextKeys = map[HeaderKey]interface{}{
	UserAccessHeaderKey:   netHttp.FlorenceIdentityKey,
	LocaleHeaderKey:       netHttp.LocaleHeaderKey,
	CollectionIDHeaderKey: netHttp.CollectionIDHeaderKey,
}

// Cookie Keys string representations
var cookieKeys = map[CookieKey]string{
	UserAccessCookieKey:   netHttp.FlorenceCookieKey,
	LocaleCookieKey:       netHttp.LocaleCookieKey,
	CollectionIDCookieKey: netHttp.CollectionIDCookieKey,
}

// Context Keys corresponding to each cookie key
var cookieContextKeys = map[CookieKey]interface{}{
	UserAccessCookieKey:   netHttp.FlorenceIdentityKey,
	LocaleCookieKey:       netHttp.LocaleHeaderKey,
	CollectionIDCookieKey: CollectionIDHeaderKey,
}

// String returns the string representation of a header key
func (hk HeaderKey) String() string {
	return headerKeys[hk]
}

// ContextKey returns the context key corresponding to a header key
func (hk HeaderKey) ContextKey() interface{} {
	return headerContextKeys[hk]
}

// String returns the string representation of a cookie key
func (ck CookieKey) String() string {
	return cookieKeys[ck]
}

// ContextKey returns the context key corresponding to a cookie key
func (ck CookieKey) ContextKey() interface{} {
	return cookieContextKeys[ck]
}

// CheckHeader is a wrapper which adds a value from the request header to context if one does not yet exist.
func CheckHeader(h http.Handler, headerKey HeaderKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if headerValue := req.Header.Get(headerKey.String()); headerValue != "" {
			req = req.WithContext(context.WithValue(req.Context(), headerKey.ContextKey(), headerValue))
		}

		h.ServeHTTP(w, req)
	})
}

// CheckCookie is a wrapper which adds a cookie value to context if one does not yet exist
func CheckCookie(h http.Handler, cookieKey CookieKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		cookieValue, err := req.Cookie(cookieKey.String())
		if err != nil {
			if err != http.ErrNoCookie {
				log.Event(req.Context(), "unexpected error while extracting value from cookie", log.ERROR, log.Error(err),
					log.Data{"cookieKey": cookieKey.String()})
			}
		} else {
			req = req.WithContext(context.WithValue(req.Context(), cookieKey.ContextKey(), cookieValue.Value))
		}

		h.ServeHTTP(w, req)
	})
}
