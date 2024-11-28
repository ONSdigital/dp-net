package links

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestURLBuild(t *testing.T) {
	Convey("Given a valid old URL and context with protocol, host, port, and path prefix", t, func() {
		oldURL := "https://api.beta.ons.gov.uk/v1/"
		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
		ctx = context.WithValue(ctx, ctxPort, "8080")
		ctx = context.WithValue(ctx, ctxPathPrefix, "")

		Convey("When URLBuild is called", func() {
			newURL, err := URLBuild(ctx, oldURL)

			Convey("Then it should return the correctly updated URL", func() {
				So(err, ShouldBeNil)
				So(newURL, ShouldEqual, "https://dataset-api:8080/v1/")
			})
		})
	})

	Convey("Given an invalid old URL", t, func() {
		oldURL := ":/invalid-url"
		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
		ctx = context.WithValue(ctx, ctxPort, "8080")
		ctx = context.WithValue(ctx, ctxPathPrefix, "/v1")

		Convey("When URLBuild is called", func() {
			newURL, err := URLBuild(ctx, oldURL)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(newURL, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a context without port and path prefix", t, func() {
		oldURL := "https://api.beta.ons.gov.uk/v1/"
		ctx := context.WithValue(context.Background(), ctxProtocol, "https")
		ctx = context.WithValue(ctx, ctxHost, "dataset-api")
		ctx = context.WithValue(ctx, ctxPort, "")
		ctx = context.WithValue(ctx, ctxPathPrefix, "")

		Convey("When URLBuild is called", func() {
			newURL, err := URLBuild(ctx, oldURL)

			Convey("Then it should return the updated URL without port and path prefix", func() {
				So(err, ShouldBeNil)
				So(newURL, ShouldEqual, "https://dataset-api/v1/")
			})
		})
	})

	Convey("Given a context with some missing or default values", t, func() {
		oldURL := "https://api.beta.ons.gov.uk/v1/"
		ctx := context.WithValue(context.Background(), ctxProtocol, "")
		ctx = context.WithValue(ctx, ctxHost, "")
		ctx = context.WithValue(ctx, ctxPort, "")
		ctx = context.WithValue(ctx, ctxPathPrefix, "")

		Convey("When URLBuild is called", func() {
			newURL, err := URLBuild(ctx, oldURL)

			Convey("Then it should return a URL with the old protocol and host, and no changes applied", func() {
				So(err, ShouldBeNil)
				So(newURL, ShouldEqual, "/v1/")
			})
		})
	})
}
