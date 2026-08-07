package request

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestRequestBody struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
}

func createValidJSON(id string, number int) []byte {
	jsonString := fmt.Sprintf(`{"id":%q,"number":%d}`, id, number)

	return []byte(jsonString)
}

func createGetRequestBodyTestRequest(body []byte) *http.Request {
	var bodyReader *bytes.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	r, _ := http.NewRequest("POST", "/test", bodyReader)
	return r
}

func TestGetRequestBody(t *testing.T) {
	const (
		ID     = "test-id-1234"
		Number = 123450
	)
	Convey("When the request body contains valid JSON", t, func() {
		validJSON := createValidJSON(ID, Number)
		r := createGetRequestBodyTestRequest(validJSON)

		Convey("Then it should parse successfully and return the struct", func() {
			result, err := GetJSONRequestBody[TestRequestBody](r)

			So(err, ShouldBeNil)
			So(result, ShouldNotBeNil)
			So(result.ID, ShouldEqual, ID)
			So(result.Number, ShouldEqual, Number)
		})
	})

	Convey("When the request body is empty", t, func() {
		r := createGetRequestBodyTestRequest([]byte{})

		Convey("Then it should return an error", func() {
			result, err := GetJSONRequestBody[TestRequestBody](r)

			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, errorDescriptionMalformedRequest)
		})
	})

	Convey("When the request body contains invalid JSON", t, func() {
		invalidJSON := []byte(`{"id":"bundle","number":}`)
		r := createGetRequestBodyTestRequest(invalidJSON)

		Convey("Then it should return an error", func() {
			result, err := GetJSONRequestBody[TestRequestBody](r)

			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, errorDescriptionMalformedRequest)
		})
	})

	Convey("When the request body contains malformed JSON", t, func() {
		malformedJSON := []byte(`{this is not valid json}`)
		r := createGetRequestBodyTestRequest(malformedJSON)

		Convey("Then it should return an error", func() {
			result, err := GetJSONRequestBody[TestRequestBody](r)

			So(result, ShouldBeNil)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, errorDescriptionMalformedRequest)
		})
	})
}
