package links

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FromHeadersOrDefault(t *testing.T) {

	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			defaultURL    string
			fwdProto      string
			fwdHost       string
			fwdPathPrefix string
			want          string
		}{
			// Without any forwarded headers
			{
				"http://localhost:8080/",
				"",
				"",
				"",
				"http://localhost:8080/",
			},
			// With all forwarded headers
			{
				"http://localhost:8080/",
				"https",
				"api.external.host",
				"/prefix",
				"https://api.external.host/prefix",
			},
			// With only forwarded proto
			{
				"http://localhost:8080/",
				"https",
				"",
				"",
				"http://localhost:8080/",
			},
			// With only forwarded host
			{
				"http://localhost:8080/",
				"",
				"api.external.host",
				"",
				"https://api.external.host/",
			},
			// With only forwarded path prefix
			{
				"http://localhost:8080/",
				"",
				"",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// Without all headers except forwarded proto
			{
				"http://localhost:8080/",
				"",
				"api.external.host",
				"/prefix",
				"https://api.external.host/prefix",
			},
			// Without all headers except forwarded host
			{
				"http://localhost:8080/",
				"https",
				"",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// Without all headers except forwarded path prefix
			{
				"http://localhost:8080/",
				"https",
				"api.external.host",
				"",
				"https://api.external.host/",
			},
			// With only forwarded proto and host
			{
				"http://localhost:8080/",
				"https",
				"api.external.host",
				"",
				"https://api.external.host/",
			},
			// With only forwarded prefix and host
			{
				"http://localhost:8080/",
				"",
				"api.external.host",
				"/prefix",
				"https://api.external.host/prefix",
			},
			// With only forwarded proto and prefix
			{
				"http://localhost:8080/",
				"https",
				"",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// With non-external forwarded host
			{
				"http://localhost:8080/",
				"",
				"internalhost",
				"",
				"http://localhost:8080/",
			},
		}

		for _, tt := range tests {
			du, err := url.Parse(tt.defaultURL)
			So(err, ShouldBeNil)

			h := http.Header{}
			if tt.fwdProto != "" {
				h.Add("X-Forwarded-Proto", tt.fwdProto)
			}
			if tt.fwdHost != "" {
				h.Add("X-Forwarded-Host", tt.fwdHost)
			}
			if tt.fwdPathPrefix != "" {
				h.Add("X-Forwarded-Path-Prefix", tt.fwdPathPrefix)
			}

			du.JoinPath()
			r := &http.Request{
				URL:  &url.URL{},
				Host: "localhost:8080",
			}
			builder := FromHeadersOrDefault(&h, r, du)
			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)

		}

	})

	Convey("Given an empty incoming request host", t, func() {
		r := &http.Request{
			URL:  &url.URL{},
			Host: "",
		}

		Convey("When the builder is created with forwarded headers", func() {
			h := &http.Header{}
			h.Add("X-Forwarded-Proto", "https")
			h.Add("X-Forwarded-Host", "api.newhost")
			h.Add("X-Forwarded-Path-Prefix", "v1")

			defaultURL, err := url.Parse("http://localhost:8080/")
			So(err, ShouldBeNil)

			builder := FromHeadersOrDefault(h, r, defaultURL)

			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)

			Convey("Then the builder URL should be the default URL with the path prefix", func() {
				So(builder.URL.String(), ShouldEqual, "http://localhost:8080/v1")
			})
		})

		Convey("When the builder is created without forwarded headers", func() {
			h := &http.Header{}

			defaultURL, err := url.Parse("http://localhost:8080/")
			So(err, ShouldBeNil)

			builder := FromHeadersOrDefault(h, r, defaultURL)

			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)

			Convey("Then the builder URL should be the default URL", func() {
				So(builder.URL.String(), ShouldEqual, "http://localhost:8080/")
			})
		})
	})

}

func TestBuilder_BuildLink(t *testing.T) {

	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			builderURL string
			oldLink    string
			want       string
		}{
			// Empty old link
			{
				"http://localhost:8080/",
				"",
				"http://localhost:8080/",
			},
			// Old link with no path
			{
				"http://localhost:8080/",
				"http://localhost:8080/",
				"http://localhost:8080/",
			},
			// Old link with different base url
			{
				"http://localhost:8080/",
				"https://oldHost:1000/",
				"http://localhost:8080/",
			},
			// Old link with path
			{
				"http://localhost:8080/",
				"http://localhost:8080/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link with path and different base url
			{
				"http://localhost:8080/",
				"http://oldHost:1000/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link without base url
			{
				"http://localhost:8080/",
				"/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link with query params
			{
				"http://localhost:8080/",
				"http://localhost:8080/some/path?param1=value1&param2=value2",
				"http://localhost:8080/some/path?param1=value1&param2=value2",
			},
			// Old external link to new internal url
			{
				"http://localhost:8080/",
				"https://some.api.host/v1/some/path",
				"http://localhost:8080/some/path",
			},
			// Old external link to new external url
			{
				"https://some.api.host/v1",
				"https://some.api.host/v1/some/path",
				"https://some.api.host/v1/some/path",
			},
			// Old internal link to new external url
			{
				"https://some.api.host/v1",
				"http://localhost:8080/some/path",
				"https://some.api.host/v1/some/path",
			},
			// Old internal link to new external url with query params
			{
				"https://some.api.host/v1",
				"http://localhost:8080/some/path?param1=value1&param2=value2",
				"https://some.api.host/v1/some/path?param1=value1&param2=value2",
			},
		}

		for _, tt := range tests {

			bu, err := url.Parse(tt.builderURL)
			So(err, ShouldBeNil)
			builder := &Builder{URL: bu}

			newurl, err := builder.BuildLink(tt.oldLink)
			So(err, ShouldBeNil)
			So(newurl, ShouldEqual, tt.want)

			// Check that the function hasn't modified the builder's internal URL
			So(builder.URL.String(), ShouldEqual, tt.builderURL)
		}

	})

	Convey("When an invalid old URL is provided", t, func() {
		builder := &Builder{URL: &url.URL{}}
		invalidURL := ":invalid/url"
		newurl, err := builder.BuildLink(invalidURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})

}
func Test_FromHeadersOrDefault_HeaderRemoval(t *testing.T) {
	Convey("Given a request with headers that should be removed", t, func() {
		h := &http.Header{}
		h.Add("Authorization", "Bearer token")
		h.Add("X-Florence-Token", "florence-token")
		h.Add("X-Forwarded-Proto", "https")
		h.Add("X-Forwarded-Host", "api.external.host")
		h.Add("X-Forwarded-Path-Prefix", "/prefix")

		r := &http.Request{
			URL:    &url.URL{},
			Host:   "localhost:8080",
			Header: *h,
		}

		defaultURL, err := url.Parse("http://localhost:8080/")
		So(err, ShouldBeNil)

		Convey("When the builder is created", func() {
			builder := FromHeadersOrDefault(h, r, defaultURL)

			Convey("Then the builder URL should be correct", func() {
				So(builder, ShouldNotBeNil)
				So(builder.URL, ShouldNotBeNil)
				So(builder.URL.String(), ShouldEqual, "https://api.external.host/prefix")
			})

			Convey("And the headers should not contain Authorization and X-Florence-Token", func() {
				So(h.Get("Authorization"), ShouldBeEmpty)
				So(h.Get("X-Florence-Token"), ShouldBeEmpty)
				So(r.Header.Get("Authorization"), ShouldBeEmpty)
				So(r.Header.Get("X-Florence-Token"), ShouldBeEmpty)
			})
		})
	})
}
