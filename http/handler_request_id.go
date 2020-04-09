package http

import (
	"net/http"
)

// HandlerRequestID is a wrapper which adds an X-Request-Id header if one does not yet exist
func HandlerRequestID(size int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			requestID := req.Header.Get(RequestHeaderKey)

			if len(requestID) == 0 {
				requestID = NewRequestID(size)
				AddRequestIdHeader(req, requestID)
			}

			h.ServeHTTP(w, req.WithContext(WithRequestId(req.Context(), requestID)))
		})
	}
}
