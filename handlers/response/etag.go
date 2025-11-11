package response

import (
	//nolint:gosec // SHA-1 is used for ETag generation only, not for security purposes
	"crypto/sha1"
	"fmt"
	"net/http"
)

const (
	ETagHeader = "ETag"
)

// GenerateETag generates a SHA-1 hash of the body with type []byte. SHA-1 is not cryptographically safe,
// but it has been selected for performance as we are only interested in uniqueness.
// A strong or weak eTag can be generated. Please note that ETags are surrounded by double quotes.
//
// Example: ETag = `"24decf55038de874bc6fa9cf0930adc219f15db1"`
//
// The definition of ETag is explained in https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
func GenerateETag(body []byte, weak bool) (etag string) {
	//nolint:gosec // SHA-1 is used for ETag generation only, not for security purposes
	hash := sha1.Sum(body)

	if weak {
		etag = fmt.Sprintf(`W/"%x"`, hash)
	} else {
		etag = fmt.Sprintf(`"%x"`, hash)
	}

	return etag
}

// SetETag sets the new eTag of the resource in the header. Please note that ETags are surrounded by double quotes.
//
// Example: ETag = `"24decf55038de874bc6fa9cf0930adc219f15db1"`
func SetETag(w http.ResponseWriter, newETag string) {
	w.Header().Add(ETagHeader, newETag)
}
