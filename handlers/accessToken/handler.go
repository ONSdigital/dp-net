package accessToken

import (
	"net/http"

	"github.com/ONSdigital/dp-net/handlers"
)

// CheckHeaderValueAndForwardWithRequestContext is a wrapper which adds a accessToken from the request header to context if one does not yet exist
func CheckHeaderValueAndForwardWithRequestContext(h http.Handler) http.Handler {
	return handlers.CheckHeader(h, handlers.UserAccessHeaderKey)
}

// CheckCookieValueAndForwardWithRequestContext is a wrapper which adds a accessToken from the cookie to context if one does not yet exist
func CheckCookieValueAndForwardWithRequestContext(h http.Handler) http.Handler {
	return handlers.CheckCookie(h, handlers.UserAccessCookieKey)
}
