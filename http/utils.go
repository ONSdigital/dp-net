package http

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
)

// DrainBody drains the body of the given of the given HTTP request.
func DrainBody(r *http.Request) {

	if r.Body == nil {
		return
	}

	_, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		log.Error(r.Context(), "error draining request body", err)
	}

	err = r.Body.Close()
	if err != nil {
		log.Error(r.Context(), "error closing request body", err)
	}
}

// GetFreePort is simple utility to find a free port on the "localhost" interface of the host machine
// for a local server to use. This is especially useful for testing purposes
func GetFreePort() (port int, err error) {
	var l *net.TCPListener

	l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(`127.0.0.1`)})
	if err != nil {
		return
	}
	defer func(l *net.TCPListener) {
		err = l.Close()
	}(l)

	return l.Addr().(*net.TCPAddr).Port, nil
}
