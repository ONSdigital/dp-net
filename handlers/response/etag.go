package response

import (
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
func GenerateETag(body []byte, weak bool) string {
	hash := sha1.Sum(body)
	etag := fmt.Sprintf(`"%x"`, hash)

	if weak {
		etag = "W/" + etag
	}

	return etag
}

// SetETag sets the new eTag of the resource in the header. Please note that ETags are surrounded by double quotes.
//
// Example: ETag = `"24decf55038de874bc6fa9cf0930adc219f15db1"`
func SetETag(w http.ResponseWriter, newETag string) {
	w.Header().Add(ETagHeader, newETag)
}
