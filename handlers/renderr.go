package handlers

import (
	"net/http"

	"github.com/ONSdigital/dp-cookies/cookies"
	render "github.com/ONSdigital/dp-renderer/v2"
	"github.com/ONSdigital/log.go/v2/log"
)

type responseInterceptor struct {
	http.ResponseWriter
	req            *http.Request
	intercepted    bool
	headersWritten bool
	headerCache    http.Header
	renderClient   *render.Render
}

func (rI *responseInterceptor) WriteHeader(status int) {
	if status >= 400 {
		log.Info(rI.req.Context(), "Intercepted error response", log.Data{"status": status})
		rI.intercepted = true
		if status == 401 || status == 404 || status == 500 {
			rI.renderErrorPage(status)
			return
		}
	}
	rI.writeHeaders()
	rI.ResponseWriter.WriteHeader(status)
}

func (rI *responseInterceptor) renderErrorPage(code int) {
	m := rI.renderClient.NewBasePageModel()

	// add cookie preferences to error page model
	preferencesCookie := cookies.GetCookiePreferences(rI.req)
	m.CookiesPreferencesSet = preferencesCookie.IsPreferenceSet
	m.CookiesPolicy.Essential = preferencesCookie.Policy.Essential
	m.CookiesPolicy.Usage = preferencesCookie.Policy.Usage

	rI.renderClient.BuildErrorPage(rI.ResponseWriter, m, code)
}

func (rI *responseInterceptor) Write(b []byte) (int, error) {
	if rI.intercepted {
		return len(b), nil
	}
	rI.writeHeaders()
	return rI.ResponseWriter.Write(b)
}

func (rI *responseInterceptor) writeHeaders() {

	if rI.headersWritten {
		return
	}

	for k, v := range rI.headerCache {
		for _, v2 := range v {
			rI.ResponseWriter.Header().Add(k, v2)
		}
	}

	rI.headersWritten = true
}

func (rI *responseInterceptor) Header() http.Header {
	return rI.headerCache
}

// Renderr is middleware that renders error pages based on response status codes
func Renderr(rendC *render.Render) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h.ServeHTTP(&responseInterceptor{w, req, false, false, make(http.Header), rendC}, req)
		})
	}
}
