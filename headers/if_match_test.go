package headers

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetIfMatch(t *testing.T) {
	Convey("When the If-Match headers exists", t, func() {
		request := http.Request{
			Header: http.Header{},
		}
		value := "if-match-value"

		request.Header.Set(HeaderIfMatch, value)

		result, err := GetIfMatchHeader(&request)

		Convey("Then the value should be returned", func() {
			So(err, ShouldBeNil)
			So(result, ShouldEqual, &value)
		})
	})

	Convey("When the If-Match header doesn't exist", t, func() {
		request := http.Request{
			Header: http.Header{},
		}

		result, err := GetIfMatchHeader(&request)

		Convey("Then the an error should be returned", func() {
			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
		})
	})
}
