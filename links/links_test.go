package links

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FromHeadersOrDefault_With_Forwarded_Headers(t *testing.T) {

	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			defaultURL    string
			fwdProto      string
			fwdHost       string
			fwdPort       string
			fwdPathPrefix string
			want          string
		}{
			// Without any forwarded headers
			{
				"http://localhost:8080/",
				"",
				"",
				"",
				"",
				"http://localhost:8080/",
			},
			// With all forwarded headers
			{
				"http://localhost:8080/",
				"https",
				"forwardedhost",
				"9090",
				"/prefix",
				"https://forwardedhost:9090/prefix",
			},
			// With only forwarded proto
			{
				"http://localhost:8080/",
				"https",
				"",
				"",
				"",
				"http://localhost:8080/",
			},
			// With only forwarded host
			{
				"http://localhost:8080/",
				"",
				"forwardedhost",
				"",
				"",
				"https://forwardedhost/",
			},
			// With only forwarded port
			{
				"http://localhost:8080/",
				"",
				"",
				"9090",
				"",
				"http://localhost:8080/",
			},
			// With only forwarded path prefix
			{
				"http://localhost:8080/",
				"",
				"",
				"",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// Without all headers except forwarded proto
			{
				"http://localhost:8080/",
				"",
				"forwardedhost",
				"9090",
				"/prefix",
				"https://forwardedhost:9090/prefix",
			},
			// Without all headers except forwarded host
			{
				"http://localhost:8080/",
				"https",
				"",
				"9090",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// Without all headers except forwarded port
			{
				"http://localhost:8080/",
				"https",
				"forwardedhost",
				"",
				"/prefix",
				"https://forwardedhost/prefix",
			},
			// Without all headers except forwarded path prefix
			{
				"http://localhost:8080/",
				"https",
				"forwardedhost",
				"9090",
				"",
				"https://forwardedhost:9090/",
			},
			// With only forwarded proto and host
			{
				"http://localhost:8080/",
				"https",
				"forwardedhost",
				"",
				"",
				"https://forwardedhost/",
			},
			// With only forwarded port and host
			{
				"http://localhost:8080/",
				"",
				"forwardedhost",
				"9090",
				"",
				"https://forwardedhost:9090/",
			},
			// With only forwarded prefix and host
			{
				"http://localhost:8080/",
				"",
				"forwardedhost",
				"",
				"/prefix",
				"https://forwardedhost/prefix",
			},
			// With only forwarded proto and port
			{
				"http://localhost:8080/",
				"https",
				"",
				"9090",
				"",
				"http://localhost:8080/",
			},
			// With only forwarded proto and prefix
			{
				"http://localhost:8080/",
				"https",
				"",
				"",
				"/prefix",
				"http://localhost:8080/prefix",
			},
			// With only forwarded port and prefix
			{
				"http://localhost:8080/",
				"",
				"",
				"9090",
				"/prefix",
				"http://localhost:8080/prefix",
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
			if tt.fwdPort != "" {
				h.Add("X-Forwarded-Port", tt.fwdPort)
			}
			if tt.fwdPathPrefix != "" {
				h.Add("X-Forwarded-Path-Prefix", tt.fwdPathPrefix)
			}

			du.JoinPath()
			r := &http.Request{
				URL: &url.URL{Scheme: "http", Host: "example.com"},
			}
			builder := FromHeadersOrDefault(&h, r, du)
			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)

		}

	})

}

func Test_FromHeadersOrDefault_Without_Forwarded_Headers(t *testing.T) {

	Convey("Given a list of test cases", t, func() {
		tests := []struct {
			incomingRequestHost string
			defaultURL          string
			want                string
		}{
			// Without incoming request host
			{
				"",
				"http://localhost:8080/",
				"http://localhost:8080/",
			},
			// With incoming request host and no port
			{
				"localhost",
				"http://localhost:8080/",
				"http://localhost/",
			},
			// With incoming request host and different port
			{
				"localhost:6789",
				"http://localhost:8080/",
				"http://localhost:6789/",
			},
			{
				"10.30.100.123:4567",
				"http://localhost:8080/",
				"http://10.30.100.123:4567/",
			},
			// With incoming request host and default URL with path
			{
				"localhost",
				"http://localhost:8080/some/path",
				"http://localhost/some/path",
			},
		}

		for _, tt := range tests {
			du, err := url.Parse(tt.defaultURL)
			So(err, ShouldBeNil)

			incomingRequest := &http.Request{
				Host: tt.incomingRequestHost,
				URL:  &url.URL{Scheme: "http", Host: "example.com"},
			}

			builder := FromHeadersOrDefault(&http.Header{}, incomingRequest, du)
			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)

		}

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
		builder := &Builder{URL: &url.URL{Scheme: "http", Host: "localhost:8080"}}
		invalidURL := ":invalid/url"
		newurl, err := builder.BuildLink(invalidURL)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "unable to parse link to URL")
		So(newurl, ShouldBeEmpty)
	})

}
