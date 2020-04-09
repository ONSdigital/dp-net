package localeCode

import (
	"context"
	"net/http"

	netHttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/log"
)

// CheckHeaderValueAndForwardWithRequestContext is a wrapper which adds a localeCode from the request header to context if one does not yet exist
func CheckHeaderValueAndForwardWithRequestContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		localeCode := req.Header.Get(netHttp.LocaleHeaderKey)

		if localeCode != "" {
			req = req.WithContext(context.WithValue(req.Context(), netHttp.LocaleHeaderKey, localeCode))
		}

		h.ServeHTTP(w, req)
	})
}

// CheckCookieValueAndForwardWithRequestContext is a wrapper which adds a localeCode from the cookie to context if one does not yet exist
func CheckCookieValueAndForwardWithRequestContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		localeCodeCookie, err := req.Cookie(netHttp.LocaleCookieKey)
		if err == nil {
			localeCode := localeCodeCookie.Value
			req = req.WithContext(context.WithValue(req.Context(), netHttp.LocaleHeaderKey, localeCode))
		} else {
			if err != http.ErrNoCookie {
				log.Event(req.Context(), "unexpected error while extracting language from cookie", log.Error(err))
			}
		}

		h.ServeHTTP(w, req)
	})
}
