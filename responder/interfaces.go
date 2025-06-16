package responder

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-net/v3/headers"
)

// RequestResponder is responsible for responding to HTTP requests, providing methods for responding
// in with JSON and handling errors
type RequestResponder interface {
	// JSON responds to a HTTP request, expecting the response body to be marshall-able into JSON
	JSON(ctx context.Context, w http.ResponseWriter, status int, resp interface{})

	// Error responds with a single error, formatted to fit in ONS's desired error
	// response structure (essentially an array of errors)
	Error(ctx context.Context, w http.ResponseWriter, status int, err error)

	// Error responds with a single error, formatted to fit in ONS's desired error
	// response structure (essentially an array of errors)
	Errors(ctx context.Context, w http.ResponseWriter, status int, errs []error)

	// StatusCode responds with a raw status code
	StatusCode(w http.ResponseWriter, status int)

	// Bytes responds to a http request with the raw bytes of whatever's passed as resp.
	// Can be used to respond with a raw string, bytes, pre-encoded object etc
	Bytes(ctx context.Context, w http.ResponseWriter, status int, resp []byte)
}

// HTTPResponseBuilder provides a fluent interface for building HTTP responses.
// It allows chaining method calls to configure headers, body, ETag, cache control,
// and status code before writing the response.
type HTTPResponseBuilder interface {
	// WithHeader adds a custom header to the HTTP response
	WithHeader(key string, value string) HTTPResponseBuilder

	// WithETag sets the ETag header for the HTTP response.
	WithETag(value string) HTTPResponseBuilder

	// WithCacheControl sets the Cache-Control header for the HTTP response.
	WithCacheControl(value headers.CacheControl) HTTPResponseBuilder

	// WithJSONBody sets the response body as JSON content.
	WithJSONBody(body any) HTTPResponseBuilder

	// WithBody sets the response body with a specific content type.
	WithBody(bodyType headers.ContentType, body any) HTTPResponseBuilder

	// WithStatusCode sets the status code of the HTTP response
	WithStatusCode(statusCode int) HTTPResponseBuilder

	// Build constructs and writes the HTTP response with the values that were set in the various With... methods
	Build(ctx context.Context, w http.ResponseWriter) error
}
