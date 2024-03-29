package reverseproxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func Create(proxyURL *url.URL, directorFunc func(*http.Request), modifyResponseFunc func(*http.Response) error) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	director := proxy.Director
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	proxy.Director = func(req *http.Request) {
		director(req)
		req.Host = proxyURL.Host
		if directorFunc != nil {
			directorFunc(req)
		}
	}
	if modifyResponseFunc != nil {
		proxy.ModifyResponse = modifyResponseFunc
	}

	return proxy
}
