package headers

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHeader(t *testing.T) {
	headerKey := "header-key"
	Convey("When the header exists", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		value := "some-test-header"
		request.Header.Set(headerKey, value)

		result, err := GetHeader(&request, headerKey)

		Convey("Then the value should be returned", func() {
			So(err, ShouldBeNil)
			So(result, ShouldEqual, &value)
		})
	})

	Convey("When the header doesn't exist", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		result, err := GetHeader(&request, headerKey)

		Convey("Then the an error should be returned", func() {
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
