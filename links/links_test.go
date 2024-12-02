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
			fwdPort       string
			fwdPathPrefix string
			want          string
		}{
			{"http://localhost:8080/",
				"",
				"",
				"",
				"",
				"http://localhost:8080/",
			},
			{"http://localhost:8080/",
				"",
				"moo.quack",
				"",
				"",
				"http://moo.quack/",
			},
			{"http://localhost:8080/",
				"https",
				"api.blah",
				"",
				"",
				"https://api.blah/",
			},
			{"http://localhost:8080/",
				"http",
				"localhost",
				"50505",
				"",
				"http://localhost:50505/",
			},
			{"http://localhost:8080/",
				"https",
				"api.blah",
				"",
				"v1",
				"https://api.blah/v1",
			},

			// TODO: Add test cases.
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
			builder := FromHeadersOrDefault(&h, du)
			So(builder, ShouldNotBeNil)
			So(builder.URL, ShouldNotBeNil)
			So(builder.URL.String(), ShouldEqual, tt.want)

		}

	})

}

//func TestURLBuild(t *testing.T) {
//	Convey("Given a valid old URL and context with protocol, host, port, and path prefix", t, func() {
//		oldURL := "https://api.beta.ons.gov.uk/v1/"
//		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
//		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
//		ctx = context.WithValue(ctx, ctxPort, "8080")
//		ctx = context.WithValue(ctx, ctxPathPrefix, "")
//
//		Convey("When URLBuild is called", func() {
//			newURL, err := URLBuild(ctx, oldURL)
//
//			Convey("Then it should return the correctly updated URL", func() {
//				So(err, ShouldBeNil)
//				So(newURL, ShouldEqual, "https://dataset-api:8080/v1/")
//			})
//		})
//	})
//
//	Convey("Given an invalid old URL", t, func() {
//		oldURL := ":/invalid-url"
//		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
//		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
//		ctx = context.WithValue(ctx, ctxPort, "8080")
//		ctx = context.WithValue(ctx, ctxPathPrefix, "/v1")
//
//		Convey("When URLBuild is called", func() {
//			newURL, err := URLBuild(ctx, oldURL)
//
//			Convey("Then it should return an error", func() {
//				So(err, ShouldNotBeNil)
//				So(newURL, ShouldBeEmpty)
//			})
//		})
//	})
//
//	Convey("Given a context without port and path prefix", t, func() {
//		oldURL := "https://api.beta.ons.gov.uk/v1/"
//		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
//		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
//		ctx = context.WithValue(ctx, ctxPort, "")
//		ctx = context.WithValue(ctx, ctxPathPrefix, "")
//
//		Convey("When URLBuild is called", func() {
//			newURL, err := URLBuild(ctx, oldURL)
//
//			Convey("Then it should return the updated URL without port and path prefix", func() {
//				So(err, ShouldBeNil)
//				So(newURL, ShouldEqual, "https://dataset-api/v1/")
//			})
//		})
//	})
//
//	Convey("Given a context with some missing or default values", t, func() {
//		oldURL := "https://api.beta.ons.gov.uk/v1/"
//		ctx := context.WithValue(context.Background(), ctxProtocol, "")
//		ctx = context.WithValue(ctx, ctxHost, "")
//		ctx = context.WithValue(ctx, ctxPort, "")
//		ctx = context.WithValue(ctx, ctxPathPrefix, "")
//
//		Convey("When URLBuild is called", func() {
//			newURL, err := URLBuild(ctx, oldURL)
//
//			Convey("Then it should return a URL with the old protocol and host, and no changes applied", func() {
//				So(err, ShouldBeNil)
//				So(newURL, ShouldEqual, "/v1/")
//			})
//		})
//	})
//}
