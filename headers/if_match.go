package headers

import "net/http"

func GetIfMatchHeader(r *http.Request) (*string, error) {
	return GetHeader(r, HeaderIfMatch)
}
