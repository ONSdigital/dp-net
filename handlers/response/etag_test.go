package response_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-net/v3/handlers/response"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenerateETag(t *testing.T) {
	t.Parallel()

	Convey("Given an empty body with strong etag generation", t, func() {
		body := []byte{}

		Convey("When GenerateETag is called", func() {
			eTag := response.GenerateETag(body, false)

			Convey("Then the etag should be returned", func() {
				So(eTag, ShouldNotBeEmpty)

				Convey("And the eTag should be a string of ASCII characters placed between double quotes", func() {
					So(eTag[0:1], ShouldEqual, `"`)
					So(eTag[len(eTag)-1:], ShouldEqual, `"`)

					Convey("And the eTag should not contain the weak tag", func() {
						So(eTag[0:2], ShouldNotEqual, `W/`)
					})
				})
			})
		})
	})

	Convey("Given some body with strong etag generation", t, func() {
		body := []byte("test data")

		Convey("When GenerateETag is called", func() {
			eTag := response.GenerateETag(body, false)

			Convey("Then the etag should be returned", func() {
				So(eTag, ShouldNotBeEmpty)

				Convey("And the eTag should be a string of ASCII characters placed between double quotes", func() {
					So(eTag[0:1], ShouldEqual, `"`)
					So(eTag[len(eTag)-1:], ShouldEqual, `"`)

					Convey("And the eTag should not contain the weak tag", func() {
						So(eTag[0:2], ShouldNotEqual, `W/`)
					})
				})
			})
		})
	})

	Convey("Given an empty body with weak etag generation", t, func() {
		body := []byte{}

		Convey("When GenerateETag is called", func() {
			eTag := response.GenerateETag(body, true)

			Convey("Then the etag should be returned", func() {
				So(eTag, ShouldNotBeEmpty)

				Convey("And the eTag should start with the weak tag", func() {
					So(eTag[0:2], ShouldEqual, `W/`)

					Convey("And the remainder of the eTag should be a string of ASCII characters placed between double quotes", func() {
						So(eTag[2:3], ShouldEqual, `"`)
						So(eTag[len(eTag)-1:], ShouldEqual, `"`)
					})
				})
			})
		})
	})

	Convey("Given some body with strong etag generation", t, func() {
		body := []byte(`test data`)

		Convey("When GenerateETag is called", func() {
			eTag := response.GenerateETag(body, true)

			Convey("Then the etag should be returned", func() {
				So(eTag, ShouldNotBeEmpty)

				Convey("And the eTag should start with the weak tag", func() {
					So(eTag[0:2], ShouldEqual, `W/`)

					Convey("And the rest of the eTag should be a string of ASCII characters placed between double quotes", func() {
						So(eTag[2:3], ShouldEqual, `"`)
						So(eTag[len(eTag)-1:], ShouldEqual, `"`)
					})
				})
			})
		})
	})
}

func TestSetETag(t *testing.T) {
	Convey("Given a response and a new eTag", t, func() {
		resp := httptest.NewRecorder()
		newEtag := "33a64df551425fcc55e4d42a148795d9f25f89d4"

		Convey("When SetETag is called", func() {
			response.SetETag(resp, newEtag)

			Convey("Then the new eTag should be added to the ETag header of the response", func() {
				So(resp.Header().Get(response.ETagHeader), ShouldEqual, newEtag)
			})
		})
	})
}
