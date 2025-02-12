package links

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	defaultInternalURL         = url.URL{Scheme: "http", Host: "localhost:8080"}
	defaultInternalDownloadURL = url.URL{Scheme: "http", Host: "localhost:23600"}
	defaultExternalDownloadURL = url.URL{Scheme: "https", Host: "download.api.host"}
)

func Test_FromHeadersOrDefault(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			fwdHost       string
			fwdPathPrefix string
			want          string
		}{
			// Without forwarded headers
			{
				"",
				"",
				"http://localhost:8080",
			},
			// With forwarded host and path prefix
			{
				"api.external.host",
				"prefix",
				"https://api.external.host/prefix",
			},
			// With forwarded host
			{
				"api.external.host",
				"",
				"https://api.external.host",
			},
			// With forwarded path prefix
			{
				"",
				"prefix",
				"http://localhost:8080/prefix",
			},
			// With internal forwarded host
			{
				"internalhost",
				"",
				"http://localhost:8080",
			},
			// With internal forwarded host and path prefix
			{
				"internalhost",
				"prefix",
				"http://localhost:8080/prefix",
			},
		}

		for _, tt := range tests {
			h := http.Header{}

			if tt.fwdHost != "" {
				h.Add("X-Forwarded-Host", tt.fwdHost)
			}
			if tt.fwdPathPrefix != "" {
				h.Add("X-Forwarded-Path-Prefix", tt.fwdPathPrefix)
			}

			builder := FromHeadersOrDefault(&h, &defaultInternalURL)
			So(builder, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)
		}
	})
}

func Test_BuildLink(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			builderURL string
			oldLink    string
			want       string
		}{
			// Empty old link
			{
				"http://localhost:8080",
				"",
				"http://localhost:8080",
			},
			// Old link with no path
			{
				"http://localhost:8080",
				"http://localhost:8080",
				"http://localhost:8080",
			},
			// Old link with different base url
			{
				"http://localhost:8080",
				"https://oldHost:1000",
				"http://localhost:8080",
			},
			// Old link with path
			{
				"http://localhost:8080",
				"http://localhost:8080/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link with path and different base url
			{
				"http://localhost:8080",
				"http://oldHost:1000/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link without base url
			{
				"http://localhost:8080",
				"/some/path",
				"http://localhost:8080/some/path",
			},
			// Old link with query params
			{
				"http://localhost:8080",
				"http://localhost:8080/some/path?param1=value1&param2=value2",
				"http://localhost:8080/some/path?param1=value1&param2=value2",
			},
			// Old external link to new internal url
			{
				"http://localhost:8080",
				"https://some.api.host/v1/some/path",
				"http://localhost:8080/some/path",
			},
			// Old external link to new external url
			{
				"https://some.api.host/v1",
				"https://some.api.host/v1/some/path",
				"https://some.api.host/v1/some/path",
			},
			// Old external link to new external url (multiple /v1)
			{
				"https://some.api.host/v1",
				"https://some.api.host/v1/v1/v1/some/path",
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

func Test_FromHeadersOrDefaultDownload(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			fwdHost string
			want    string
		}{
			// Without a forwarded host
			{
				"",
				"http://localhost:23600/downloads",
			},
			// With an external forwarded host
			{
				"api.external.host",
				"https://download.api.host/downloads",
			},
			// With an internal forwarded host
			{
				"internalhost",
				"http://localhost:23600/downloads",
			},
		}

		for _, tt := range tests {
			h := http.Header{}

			if tt.fwdHost != "" {
				h.Add("X-Forwarded-Host", tt.fwdHost)
			}

			builder := FromHeadersOrDefaultDownload(&h, &defaultInternalDownloadURL, &defaultExternalDownloadURL)
			So(builder, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)
		}
	})
}

func Test_BuildDownloadLink(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			builderURL string
			oldLink    string
			want       string
		}{
			// Empty old link
			{
				"http://localhost:23600/downloads",
				"",
				"http://localhost:23600/downloads",
			},
			// Old link with no path
			{
				"http://localhost:23600/downloads",
				"http://localhost:23600",
				"http://localhost:23600/downloads",
			},
			// Old link with different base url
			{
				"http://localhost:23600/downloads",
				"https://oldHost:1000",
				"http://localhost:23600/downloads",
			},
			// Old link with path
			{
				"http://localhost:23600/downloads",
				"http://localhost:23600/some/path",
				"http://localhost:23600/downloads/some/path",
			},
			// Old link with path and different base url
			{
				"http://localhost:23600/downloads",
				"http://oldHost:1000/some/path",
				"http://localhost:23600/downloads/some/path",
			},
			// Old link without base url
			{
				"http://localhost:23600/downloads",
				"/some/path",
				"http://localhost:23600/downloads/some/path",
			},
			// Old link with query params
			{
				"http://localhost:23600/downloads",
				"http://localhost:23600/some/path?param1=value1&param2=value2",
				"http://localhost:23600/downloads/some/path?param1=value1&param2=value2",
			},
			// Old external link to new internal url
			{
				"http://localhost:23600/downloads",
				"https://download.api.host/downloads/some/path",
				"http://localhost:23600/downloads/some/path",
			},
			// Old external link to new external url
			{
				"https://download.api.host/downloads",
				"https://download.api.host/downloads/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old external link to new external url (multiple /downloads)
			{
				"https://download.api.host/downloads",
				"https://download.api.host/downloads/downloads/downloads/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old internal link to new external url
			{
				"https://download.api.host/downloads",
				"http://localhost:23600/downloads/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old internal link to new external url with query params
			{
				"https://download.api.host/downloads",
				"http://localhost:23600/downloads/some/path?param1=value1&param2=value2",
				"https://download.api.host/downloads/some/path?param1=value1&param2=value2",
			},
		}

		for _, tt := range tests {
			bu, err := url.Parse(tt.builderURL)
			So(err, ShouldBeNil)
			builder := &Builder{URL: bu}

			newurl, err := builder.BuildDownloadLink(tt.oldLink)
			So(err, ShouldBeNil)
			So(newurl, ShouldEqual, tt.want)

			// Check that the function hasn't modified the builder's internal URL
			So(builder.URL.String(), ShouldEqual, tt.builderURL)
		}
	})

	Convey("When an invalid old URL is provided", t, func() {
		builder := &Builder{URL: &url.URL{}}
		invalidURL := ":invalid/url"
		newurl, err := builder.BuildDownloadLink(invalidURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})
}
