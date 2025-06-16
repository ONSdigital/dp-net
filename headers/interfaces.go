package headers

import "net/http"

type Header interface {
	String() string
	Set(w http.ResponseWriter)
}
