package responder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-net/v3/handlers/response"
	"github.com/ONSdigital/dp-net/v3/headers"
)

// Concrete implementation of HTTPResponseBuilder
type ResponseBuilder struct {
	headers    map[string]string // HTTP headers to write in the response
	etag       *string           // ETag value to write in the response headers
	body       *ResponseBody     // Response body to write
	statusCode int               // Response status code
	responder  RequestResponder  // RequestResponder used for writing responses
}

// Interface checks
var _ HTTPResponseBuilder = (*ResponseBuilder)(nil)

// ResponseBody encapsulates the response body content and its content type
type ResponseBody struct {
	contentType headers.ContentType // Content type of the body to write
	body        interface{}         // Actual response body to write to response

}

func createResponseBody(contentType headers.ContentType, body any) *ResponseBody {
	return &ResponseBody{
		contentType: contentType,
		body:        body,
	}
}

func CreateHTTPResponseBuilder() HTTPResponseBuilder {
	return CreateHTTPResponseBuilderWithResponder(&Responder{})
}

func CreateHTTPResponseBuilderWithResponder(responder RequestResponder) HTTPResponseBuilder {
	return &ResponseBuilder{
		headers:   make(map[string]string),
		responder: responder,
	}
}

// WithHeader adds a custom header for writing to the HTTP response
func (r *ResponseBuilder) WithHeader(key, value string) HTTPResponseBuilder {
	r.headers[key] = value
	return r
}

// WithETag sets the ETag header for writing to the HTTP response
func (r *ResponseBuilder) WithETag(value string) HTTPResponseBuilder {
	r.etag = &value
	return r
}

// WithCacheControl sets the Cache-Control header for writing to the HTTP response
func (r *ResponseBuilder) WithCacheControl(value headers.CacheControl) HTTPResponseBuilder {
	r.headers[headers.HeaderCacheControl] = value.String()
	return r
}

// WithJSONBody sets the body for the HTTP response, to be marshalled as a JSON response
func (r *ResponseBuilder) WithJSONBody(body any) HTTPResponseBuilder {
	return r.WithBody(headers.ContentTypeJSON, body)
}

// WithBody sets the body for the HTTP response, with the specified body content type
func (r *ResponseBuilder) WithBody(contentType headers.ContentType, body any) HTTPResponseBuilder {
	r.body = createResponseBody(contentType, body)
	return r
}

// WithStatusCode sets the HTTP status code for the response
func (r *ResponseBuilder) WithStatusCode(statusCode int) HTTPResponseBuilder {
	r.statusCode = statusCode
	return r
}

// Build constructs and writes the HTTP response to the provided ResponseWriter.
func (r *ResponseBuilder) Build(ctx context.Context, w http.ResponseWriter) error {
	if r.etag != nil {
		response.SetETag(w, *r.etag)
	}

	for key, value := range r.headers {
		w.Header().Set(key, value)
	}

	if r.statusCode > 0 {
		r.responder.StatusCode(w, r.statusCode)
	}

	if r.body != nil {
		err := r.Write(ctx, w)
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the ResponseBody's body to the ResponseWriter
func (r *ResponseBuilder) Write(ctx context.Context, w http.ResponseWriter) error {
	w.Header().Set(headers.HeaderContentType, string(r.body.contentType))

	switch r.body.contentType {
	case headers.ContentTypeJSON:
		r.responder.JSON(ctx, w, r.statusCode, r.body.body)
	default:
		return fmt.Errorf("response body type %s is not implemented", r.body.contentType)
	}

	return nil
}
