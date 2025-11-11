package links

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	defaultInternalURL = &url.URL{Scheme: "http", Host: "localhost:8080"}
	defaultDownloadURL = &url.URL{Scheme: "https", Host: "download.api.host"}
)

const (
	invalidURL = ":invalid/url"
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

			builder := FromHeadersOrDefault(&h, defaultInternalURL)
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
		newurl, err := builder.BuildLink(invalidURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})
}

func Test_BuildDownloadLink(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			oldLink string
			want    string
		}{
			// Empty old link
			{
				"",
				"https://download.api.host/downloads",
			},
			// Old link with only /downloads
			{
				"https://download.api.host/downloads",
				"https://download.api.host/downloads",
			},
			// Old link with /downloads and path
			{
				"https://download.api.host/downloads/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old link with no path
			{
				"http://localhost:23600",
				"https://download.api.host/downloads",
			},
			// Old link with different base url
			{
				"https://localhost:23600",
				"https://download.api.host/downloads",
			},
			// Old link with path
			{
				"http://localhost:23600/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old link without base url
			{
				"/some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old link without base url and / prefix
			{
				"some/path",
				"https://download.api.host/downloads/some/path",
			},
			// Old link with query params
			{
				"http://localhost:23600/some/path?param1=value1&param2=value2",
				"https://download.api.host/downloads/some/path?param1=value1&param2=value2",
			},
			// Old link with multiple /downloads
			{
				"https://download.api.host/downloads/downloads/downloads/some/path",
				"https://download.api.host/downloads/some/path",
			},
		}

		for _, tt := range tests {
			newurl, err := BuildDownloadLink(tt.oldLink, defaultDownloadURL)
			So(err, ShouldBeNil)
			So(newurl, ShouldEqual, tt.want)
		}
	})

	Convey("When an invalid old URL is provided", t, func() {
		newurl, err := BuildDownloadLink(invalidURL, defaultDownloadURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})
}

func Test_BuildDownloadNewLink(t *testing.T) {
	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			oldLink string
			want    string
		}{
			// Empty old link
			{
				"",
				"https://download.api.host/downloads-new",
			},
			// Old link with only /downloads-new
			{
				"https://download.api.host/downloads-new",
				"https://download.api.host/downloads-new",
			},
			// Old link with /downloads-new and path
			{
				"https://download.api.host/downloads-new/some/path",
				"https://download.api.host/downloads-new/some/path",
			},
			// Old link with no path
			{
				"http://localhost:23600",
				"https://download.api.host/downloads-new",
			},
			// Old link with different base url
			{
				"https://localhost:23600",
				"https://download.api.host/downloads-new",
			},
			// Old link with path
			{
				"http://localhost:23600/some/path",
				"https://download.api.host/downloads-new/some/path",
			},
			// Old link without base url
			{
				"/some/path",
				"https://download.api.host/downloads-new/some/path",
			},
			// Old link without base url and / prefix
			{
				"some/path",
				"https://download.api.host/downloads-new/some/path",
			},
			// Old link with query params
			{
				"http://localhost:23600/some/path?param1=value1&param2=value2",
				"https://download.api.host/downloads-new/some/path?param1=value1&param2=value2",
			},
			// Old link with multiple /downloads-new
			{
				"https://download.api.host/downloads-new/downloads-new/downloads-new/some/path",
				"https://download.api.host/downloads-new/some/path",
			},
		}

		for _, tt := range tests {
			newurl, err := BuildDownloadNewLink(tt.oldLink, defaultDownloadURL)
			So(err, ShouldBeNil)
			So(newurl, ShouldEqual, tt.want)
		}
	})

	Convey("When an invalid old URL is provided", t, func() {
		newurl, err := BuildDownloadNewLink(invalidURL, defaultDownloadURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})
}

func Test_RemovePrefixFromPath(t *testing.T) {
	Convey("RemovePrefixFromPath removes all leading prefix instances from the path", t, func() {
		tests := []struct {
			path   string
			prefix string
			want   string
		}{
			{"", "/prefix", ""},
			{"/prefix", "/prefix", ""},
			{"/prefix/some/path", "/prefix", "/some/path"},
			{"/prefix/prefix/prefix/some/path", "/prefix", "/some/path"},
			{"/some/path", "/prefix", "/some/path"},
			{"some/path", "/prefix", "some/path"},
		}

		for _, tt := range tests {
			result := RemovePrefixFromPath(tt.path, tt.prefix)
			So(result, ShouldEqual, tt.want)
		}
	})
}
