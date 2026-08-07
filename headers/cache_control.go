package headers

import "net/http"

// CacheControl represents a value for the Cache-Control header
type CacheControl string

const (
	CacheControlNoStore CacheControl = "no-store"
)

func GetCacheControlHeader(r *http.Request) (*string, error) {
	return GetHeader(r, HeaderCacheControl)
}

var _ Header = (*CacheControl)(nil)

func (c *CacheControl) String() string {
	return string(*c)
}

func (c *CacheControl) Set(w http.ResponseWriter) {
	w.Header().Set(HeaderCacheControl, c.String())

}
