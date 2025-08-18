package headers

import (
	"fmt"
	"net/http"
)

// Header key constants
const (
	HeaderAuthorization = "Authorization"
	HeaderCacheControl  = "Cache-Control"
	HeaderLocation      = "Location"
	HeaderContentType   = "Content-Type"
	HeaderETag          = "ETag"
	HeaderFlorenceToken = "X-Florence-Token"
	HeaderIfMatch       = "If-Match"
)

// createMissingHeaderError creates an error with a standardised message for a not-found header
func createMissingHeaderError(header string) error {
	return fmt.Errorf("header %s was not found", header)
}

// GetHeader gets the matching header from the request, and returns an error if it is missing or empty
func GetHeader(r *http.Request, header string) (*string, error) {
	value := r.Header.Get(header)
	if value == "" {
		return nil, createMissingHeaderError(header)
	}

	return &value, nil
}
