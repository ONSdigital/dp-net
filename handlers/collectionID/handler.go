package collectionID

import (
	"net/http"

	"github.com/ONSdigital/dp-net/handlers"
)

// CheckHeader is a wrapper which adds a CollectionID from the request header to context if one does not yet exist
func CheckHeader(h http.Handler) http.Handler {
	return handlers.CheckHeader(h, handlers.CollectionIDHeaderKey)
}

// CheckCookie is a wrapper which adds a CollectionID from the cookie to context if one does not yet exist
func CheckCookie(h http.Handler) http.Handler {
	return handlers.CheckCookie(h, handlers.CollectionIDCookieKey)
}
