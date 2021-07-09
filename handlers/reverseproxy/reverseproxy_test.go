package reverseproxy_test

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/ONSdigital/dp-net/v2/handlers/reverseproxy"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDirectorFunc(t *testing.T) {
	proxyURL, _ := url.Parse("https://www.ons.gov.uk")
	Convey("Create proxy", t, func() {
		reverseProxy := reverseproxy.Create(proxyURL, nil, nil)

		So(reverseProxy, ShouldNotBeNil)
		So(reverseProxy, ShouldImplement, (*http.Handler)(nil))

		req, _ := http.NewRequest(`GET`, `https://cy.ons.gov.uk`, nil)
		So(func() { reverseProxy.(*httputil.ReverseProxy).Director(req) }, ShouldNotPanic)
		So(req.URL.Host, ShouldEqual, `www.ons.gov.uk`)
	})

	Convey("Create proxy with director func", t, func() {
		var directorCalled bool
		reverseProxy := reverseproxy.Create(proxyURL, func(req *http.Request) {
			directorCalled = true
			req.URL.Host = `host`
		}, func(req *http.Response) error {
			return nil
		})

		So(reverseProxy, ShouldNotBeNil)
		So(reverseProxy, ShouldImplement, (*http.Handler)(nil))

		req, _ := http.NewRequest(`GET`, `https://cy.ons.gov.uk`, nil)
		So(func() { reverseProxy.(*httputil.ReverseProxy).Director(req) }, ShouldNotPanic)
		So(req.URL.Host, ShouldEqual, `host`)
		So(directorCalled, ShouldBeTrue)
	})
}
