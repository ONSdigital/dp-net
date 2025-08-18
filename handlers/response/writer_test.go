package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-net/v3/handlers/response"
	. "github.com/smartystreets/goconvey/convey"
)

type child struct {
	Name string `json:"value"`
}

type parent struct {
	Name  string `json:"name"`
	Child child  `json:"child"`
}
type mockOnsJSONEncoder struct {
	encodeCalls     int
	mockedBehaviour func(w http.ResponseWriter, value interface{}, status int) error
}

func (mock *mockOnsJSONEncoder) WriteResponseJSON(w http.ResponseWriter, value interface{}, status int) error {
	mock.encodeCalls++
	return mock.mockedBehaviour(w, value, status)
}

func initMock() *mockOnsJSONEncoder {
	actualImpl := &response.OnsJSONEncoder{}
	mock := &mockOnsJSONEncoder{
		encodeCalls: 0,
		mockedBehaviour: func(w http.ResponseWriter, value interface{}, status int) error {
			return actualImpl.WriteResponseJSON(w, value, status)
		},
	}
	response.JsonResponseEncoder = mock
	return mock
}

func TestWriteJSONResponse(t *testing.T) {
	var input parent
	var statusCode int
	var rec *httptest.ResponseRecorder
	mock := initMock()

	Convey("Given a valid responseWriter, response value and status code", t, func() {
		input = parent{Name: "Hello World!", Child: child{Name: "Bob!"}}
		statusCode = http.StatusOK
		rec = httptest.NewRecorder()

		Convey("When the encoder is invoked", func() {
			err := response.WriteJSON(rec, input, http.StatusOK)

			Convey("There is no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the input value is written to the response body.", func() {
				var actual parent
				err := json.Unmarshal(rec.Body.Bytes(), &actual)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, input)
			})

			Convey("And the response http status code matches the parameter passed in.", func() {
				So(rec.Code, ShouldEqual, statusCode)
			})

			Convey("And the response content type header is 'application/json'", func() {
				So(rec.Header().Get(response.ContentTypeHeader), ShouldEqual, response.ContentTypeJSON)
			})

			Convey("And the encoder is invoked the expected number of times.", func() {
				So(mock.encodeCalls, ShouldEqual, 5)
			})
		})
	})
}

func TestWriteJSONResponseWithInvalidData(t *testing.T) {
	var invalidInput interface{}
	var statusCode int
	var rec *httptest.ResponseRecorder
	mock := initMock()

	Convey("Given a valid responseWriter, an invalid response value and a valid status code", t, func() {
		rec = httptest.NewRecorder()
		invalidInput = func() string {
			return "HelloWorld"
		}
		statusCode = http.StatusInternalServerError

		Convey("When the encoder is invoked", func() {
			err := response.WriteJSON(rec, invalidInput, http.StatusOK)

			Convey("And the response content type header is 'application/json'", func() {
				So(rec.Header().Get(response.ContentTypeHeader), ShouldEqual, response.ContentTypeJSON)
			})

			Convey("Then an http internal server error status is written to the response.", func() {
				So(rec.Code, ShouldEqual, statusCode)
			})

			Convey("And the encoder is invoked the expected number of times.", func() {
				So(mock.encodeCalls, ShouldEqual, 3)
			})

			Convey("And an error is returned", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
