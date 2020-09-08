package handlers

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/log"
	"net/http"
)

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

/* My head hurts...
ControllerHandler is called from an app which then creates a function of type ControllerHandlerFunc which returns a http.Handler (and error)
The former is then passed to the ControllerHandlerMiddleware function which checks all fields are present and accounted for
then calls the ControllerHandlerFunc implementation which sets the headers as required to finally return a http.Handler back to the original app that called ControllerHandler
*/
//TODO all the error handling

// TODO move type to top of file
type ControllerHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) (http.Handler, error)

/* Should ControllerHandlerFunc type be implemented here too or are the individual apps providing the
'handlerFunc ControllerHandlerFunc' as seen on line 79? As the headers will have to be set somewhere. If here then line 68 shows how this might roughly be done
*/

//WIP - as above not 100% if this should be done here or in app but this is roughly what it could look like
func ControllerHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) (http.Handler, error) {
	var controllerHandlerFunc ControllerHandlerFunc = func(ctx, w, r, lang, collectionID, accessToken) (http.Handler, error) {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// TODO set headers here
			//	h.ServeHTTP(w, req) // The app should handle the serving right? Ok to remove this.
		})
	}
	return ControllerHandlerMiddleware(ctx, w, r, controllerHandlerFunc)
}

// TODO possibly rename but can't think of a better name atm.
func ControllerHandlerMiddleware(ctx context.Context, w http.ResponseWriter, r *http.Request, handlerFunc ControllerHandlerFunc) (http.Handler, error) {
	// This function checks subdomain and cookie then defaults to English if not set
	// - Should this be grabbing it instead? It could then pass it as an argument in the 'lang' field instead shown below on line 85 4th argument as ""?
	request.SetLocaleCode(r)
	accessToken := getUserAccessTokenFromContext(ctx)
	collectionID := getCollectionIDFromContext(ctx)
	x, y := handlerFunc(ctx, w, r, "", collectionID, accessToken)
	return x, y
}

func getUserAccessTokenFromContext(ctx context.Context) string {
	if ctx.Value(request.FlorenceIdentityKey) != nil {
		accessToken, ok := ctx.Value(request.FlorenceIdentityKey).(string)
		if !ok {
			log.Event(ctx, "error retrieving user access token", log.WARN, log.Error(errors.New("error casting access token context value to string")))
		}
		return accessToken
	}
	return ""
}

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
