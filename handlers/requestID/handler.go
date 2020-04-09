package requestID

import (
	"net/http"

	netHttp "github.com/ONSdigital/dp-net/http"
)

// Handler is a wrapper which adds an X-Request-Id header if one does not yet exist
func Handler(size int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			requestID := req.Header.Get(netHttp.RequestHeaderKey)

			if len(requestID) == 0 {
				requestID = netHttp.NewRequestID(size)
				netHttp.AddRequestIdHeader(req, requestID)
			}

			h.ServeHTTP(w, req.WithContext(netHttp.WithRequestId(req.Context(), requestID)))
		})
	}
}
