package handlers

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-api-clients-go/headers"
	"github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

// ControllerHandlerFunc is a function type that accepts arguments required for logical flow in handlers Implementation of function should set headers as needed
type ControllerHandlerFunc func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string)

// CheckHeader is a wrapper which adds a value from the request header (if found) to the request context.
// This function complies with alice middleware Constructor type: func(http.Handler) -> (http.Handler)
func CheckHeader(key Key) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return DoCheckHeader(h, key)
	}
}

// DoCheckHeader returns a handler that performs the CheckHeader logic
func DoCheckHeader(h http.Handler, key Key) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if headerValue := req.Header.Get(key.Header()); headerValue != "" {
			req = req.WithContext(context.WithValue(req.Context(), key.Context(), headerValue))
		}
		h.ServeHTTP(w, req)
	})
}

// CheckCookie is a wrapper which adds a cookie value (if found) to the request context.
// This function complies with alice middleware Constructor type: func(http.Handler) -> (http.Handler))
func CheckCookie(key Key) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return DoCheckCookie(h, key)
	}
}

// DoCheckCookie returns a handler that performs the CheckCookie logic
func DoCheckCookie(h http.Handler, key Key) http.Handler {
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

// ControllerHandler is a middleware handler that ensures all required logical arguments are being passed along and then returns a generic http.Handler
func ControllerHandler(controllerHandlerFunc ControllerHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		locale := request.GetLocaleCode(r)
		collectionID := getCollectionIDFromContext(ctx)
		accessToken := getUserAccessTokenFromContext(ctx, r)

		controllerHandlerFunc(w, r, locale, collectionID, accessToken)
	}
}

// getUserAccessTokenFromContext will get the Florence Identity Key aka access token from a context or header
func getUserAccessTokenFromContext(ctx context.Context, r *http.Request) string {
	var err error = nil
	accessToken := ""
	ok := false
	if ctx.Value(request.FlorenceIdentityKey) != nil {
		accessToken, ok = ctx.Value(request.FlorenceIdentityKey).(string)
		if !ok {
			log.Event(ctx, "error retrieving user access token", log.WARN, log.Error(errors.New("error casting access token context value to string")))
		}
	} else {
		accessToken, err = headers.GetUserAuthToken(r)
		if err != nil {
			log.Event(ctx, "access token not found", log.WARN, log.Error(err))
		}
	}
	return accessToken
}

// getCollectionIDFromContext will get the Collection ID token from a context
func getCollectionIDFromContext(ctx context.Context) string {
	if ctx.Value(request.CollectionIDHeaderKey) != nil {
		collectionID, ok := ctx.Value(request.CollectionIDHeaderKey).(string)
		if !ok {
			log.Event(ctx, "error retrieving collection ID", log.WARN, log.Error(errors.New("error casting collection ID context value to string")))
		}
		return collectionID
	}
	return ""
}
