package headers

import "net/http"

type ContentType string

func (c *ContentType) String() string {
	return string(*c)
}

const (
	ContentTypeJSON ContentType = "application/json; charset=utf-8"
)

var _ Header = (*ContentType)(nil)

func GetContentTypeHeader(r *http.Request) (*string, error) {
	return GetHeader(r, HeaderContentType)
}

func (c *ContentType) Set(w http.ResponseWriter) {
	w.Header().Set(HeaderContentType, c.String())
}
