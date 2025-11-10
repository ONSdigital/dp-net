package fallback

import (
	"net/http"
)

// AlternativeBuilder is a struct that helps to build a fallback handler by the use of a narrative structure
// for example using the following:
//
//	fallback.Try(handler1).WhenStatus(http.StatusNotFound).Then(handler2)
type AlternativeBuilder struct {
	tryHandler http.Handler
	whenStatus *int
}

// Try takes a `http.Handler` and returns an incomplete AlternativeBuilder
func Try(h http.Handler) *AlternativeBuilder {
	return &AlternativeBuilder{tryHandler: h}
}

// WhenStatus extends an AlternativeBuilder and includes a condition on the http status
func (ab *AlternativeBuilder) WhenStatus(status int) *AlternativeBuilder {
	ab.whenStatus = &status
	return ab
}

// Then takes an AlternativeBuilder and returns a complete Alternative with the fallback handler provided
func (ab *AlternativeBuilder) Then(handler http.Handler) *Alternative {
	return &Alternative{
		TryHandler:   ab.tryHandler,
		WhenStatusIs: ab.whenStatus,
		ThenHandler:  handler,
	}
}

// Alternative implements a [http.Handler] which wraps another http.Handler and tries serving it first. Depending on the
// status returned by the first handler it will either return the response to the caller immediately or it will pass a
// copy of the request to the second handler instead, retuning that handlers response in that case.
type Alternative struct {
	TryHandler   http.Handler
	WhenStatusIs *int
	ThenHandler  http.Handler
}
type responseWriter struct {
	header     http.Header
	statusCode int
	body       []byte
}

var _ http.ResponseWriter = &responseWriter{}

func (t *responseWriter) Header() http.Header {
	m := t.header
	if m == nil {
		m = make(http.Header)
		t.header = m
	}
	return m
}

func (t *responseWriter) Write(bytes []byte) (int, error) {
	t.body = append(t.body, bytes...)
	return len(bytes), nil
}

func (t *responseWriter) WriteHeader(statusCode int) {
	t.statusCode = statusCode
}

func (t *responseWriter) Body() []byte {
	return t.body
}

func (t *responseWriter) StatusCode() int {
	return t.statusCode
}

func (alternative *Alternative) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Split the request's body reader so that it can be read by both handlers
	readClosers := ReadCloserSplit(r.Body, 2)
	bodyReadCloser1 := readClosers[0]
	bodyReadCloser2 := readClosers[1]
	defer bodyReadCloser1.Close() // TODO not sure we even need to do this, will handlers defer close anyway?
	defer bodyReadCloser2.Close()

	w1 := responseWriter{} // TODO don't slurp response body of first request
	r1 := r.Clone(r.Context())
	r1.Body = bodyReadCloser1
	alternative.TryHandler.ServeHTTP(&w1, r1)

	if alternative.WhenStatusIs == nil || w1.statusCode != *alternative.WhenStatusIs {
		for k, vs := range w1.header {
			w.Header()[k] = vs
		}
		w.WriteHeader(w1.statusCode)
		w.Write(w1.body)
		return
	}

	r2 := r.Clone(r.Context())
	r2.Body = bodyReadCloser2
	alternative.ThenHandler.ServeHTTP(w, r2)
}

func (alternative *Alternative) WhenStatus(status int) *AlternativeBuilder {
	return &AlternativeBuilder{
		tryHandler: alternative,
		whenStatus: &status,
	}
}
