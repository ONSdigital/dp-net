package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetContentType(t *testing.T) {
	Convey("When Content-Type header exists", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		value := "application/json"
		request.Header.Set(HeaderContentType, value)

		result, err := GetContentTypeHeader(&request)

		Convey("Then the value should be returned", func() {
			So(err, ShouldBeNil)
			So(result, ShouldEqual, &value)
		})
	})

	Convey("When Content-Type header doesn't exist", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		result, err := GetContentTypeHeader(&request)

		Convey("Then the an error should be returned", func() {
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestContentType_Set(t *testing.T) {
	Convey("ContentType set should set request header", t, func() {
		contentType := ContentTypeJSON

		rr := httptest.NewRecorder()
		contentType.Set(rr)

		So(rr.Header(), ShouldContainKey, HeaderContentType)
		So(rr.Header().Get(HeaderContentType), ShouldEqual, contentType.String())
	})
}
