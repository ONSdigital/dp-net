package headers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetCacheControl(t *testing.T) {
	Convey("When Cache-Control header exists", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		value := "no-store"
		request.Header.Set(HeaderCacheControl, value)

		result, err := GetCacheControlHeader(&request)

		Convey("Then the value should be returned", func() {
			So(err, ShouldBeNil)
			So(result, ShouldEqual, &value)
		})
	})

	Convey("When Cache-Control header doesn't exist", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		result, err := GetCacheControlHeader(&request)

		Convey("Then the an error should be returned", func() {
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestCacheControl_Set(t *testing.T) {
	Convey("CacheControl set should set request header", t, func() {
		cacheControl := CacheControlNoStore

		rr := httptest.NewRecorder()
		cacheControl.Set(rr)

		So(rr.Header(), ShouldContainKey, HeaderCacheControl)
		So(rr.Header().Get(HeaderCacheControl), ShouldEqual, cacheControl.String())
	})
}
