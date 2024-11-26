package links

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestNewMiddleware checks that NewMiddleWare initializes correctly.
func TestNewMiddleware(t *testing.T) {
	Convey("Given a default URL", t, func() {
		defaultURL := "http://localhost:22000"

		Convey("When NewMiddleware is created", func() {
			mwFunc, err := NewMiddleWare(defaultURL)

			Convey("Then it should initialize without errors", func() {
				So(err, ShouldBeNil)
				So(mwFunc, ShouldNotBeNil)
			})

			Convey("And it should forward requests to the wrapped handler", func() {
				handler := mwFunc(http.NotFoundHandler())
				rec := httptest.NewRecorder()
				req := httptest.NewRequest("GET", "/datasets", nil)

				handler.ServeHTTP(rec, req)

				So(rec.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

// TestMiddleware_ServeHTTP checks the behavior of ServeHTTP.
func TestMiddleware_ServeHTTP(t *testing.T) {
	Convey("Given a Middleware with default values", t, func() {
		mw := &Middleware{
			DefaultProtocol:   "https",
			DefaultHost:       "localhost",
			DefaultPort:       "22000",
			DefaultUrlVersion: "/v1",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				So(ctx.Value(ctxProtocol), ShouldEqual, "https")
				So(ctx.Value(ctxHost), ShouldEqual, "dataset-api")
				So(ctx.Value(ctxPort), ShouldEqual, "23200")
			}),
		}

		Convey("When a request with forwarded headers is handled", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/datasets", nil)

			req.Header.Set("X-Forwarded-Proto", "https")
			req.Header.Set("X-Forwarded-Host", "dataset-api")
			req.Header.Set("X-Forwarded-Port", "23200")

			mw.ServeHTTP(rec, req)

			Convey("Then the context should contain the correct forwarded values", func() {
				So(rec.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

// TestGetForwardedHeaderElseDefault validates the header fallback logic.
func TestGetForwardedHeaderElseDefault(t *testing.T) {
	Convey("Given HTTP headers", t, func() {
		header := http.Header{}
		header.Set("X-Forwarded-Proto", "https")

		Convey("When getting a header that exists", func() {
			value, found := getForwardedHeaderElseDefault(header, "X-Forwarded-Proto", "http")

			Convey("Then it should return the forwarded value", func() {
				So(value, ShouldEqual, "https")
				So(found, ShouldBeTrue)
			})
		})

		Convey("When getting a header that does not exist", func() {
			value, found := getForwardedHeaderElseDefault(header, "X-Forwarded-Port", "8080")

			Convey("Then it should return the default value", func() {
				So(value, ShouldEqual, "8080")
				So(found, ShouldBeFalse)
			})
		})
	})
}

// TestURLBuild verifies URL building logic.
func TestURLBuild(t *testing.T) {
	Convey("Given a context with protocol, host, port, and version", t, func() {
		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
		ctx = context.WithValue(ctx, ctxHost, "example.org")
		ctx = context.WithValue(ctx, ctxPort, "443")
		ctx = context.WithValue(ctx, ctxUrlVersion, "/v2")

		Convey("When building URLs", func() {
			tests := []struct {
				oldURL   string
				expected string
			}{
				{
					oldURL:   "http://example.com/api/resource",
					expected: "https://example.org:443/v2/api/resource",
				},
				{
					oldURL:   "http://example.com:8080/api/resource",
					expected: "https://example.org:443/v2/api/resource",
				},
			}

			for _, tt := range tests {
				Convey("Then the new URL should match the expected value", func() {
					newURL, err := URLBuild(ctx, tt.oldURL)
					So(err, ShouldBeNil)
					So(newURL, ShouldEqual, tt.expected)
				})
			}
		})
	})
}
