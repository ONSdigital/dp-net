package responder

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-net/v3/headers"

	. "github.com/smartystreets/goconvey/convey"
)

type responseBodyTestStruct struct {
	ID     string
	Number int
}

func CreateMockResponder() RequestResponderMock {
	return RequestResponderMock{
		JSONFunc: func(ctx context.Context, w http.ResponseWriter, status int, resp interface{}) {

		},
		StatusCodeFunc: func(w http.ResponseWriter, status int) {
			w.WriteHeader(status)
		},
	}
}

func TestCreateResponseBuilder(t *testing.T) {
	Convey("CreateResponseBuilder", t, func() {
		Convey("Should create a response builder", func() {
			responseBuilder := CreateHTTPResponseBuilder()

			So(responseBuilder, ShouldNotBeNil)
		})
	})

	Convey("CreateResponseBuilder", t, func() {
		Convey("Should create a response builder with default values", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder()

			So(httpResponseBuilder, ShouldNotBeNil)

			responseBuilder := httpResponseBuilder.(*ResponseBuilder)
			So(responseBuilder.headers, ShouldNotBeNil)
			So(responseBuilder.headers, ShouldBeEmpty)

			So(responseBuilder.statusCode, ShouldEqual, 0)
			So(responseBuilder.body, ShouldBeNil)

			So(responseBuilder.etag, ShouldBeNil)
		})
	})

	Convey("WithHeader", t, func() {
		headerKey := "testing header"
		headerValue := "test header value"
		Convey("Should set header", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithHeader(headerKey, headerValue)

			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.headers, ShouldNotBeEmpty)
			So(responseBuilder.headers, ShouldContainKey, headerKey)
			actualValue := responseBuilder.headers[headerKey]
			So(actualValue, ShouldEqual, headerValue)
		})
	})

	Convey("WithETag", t, func() {
		etagValue := "etag value"
		Convey("Should set etag", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithETag(etagValue)
			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.etag, ShouldNotBeNil)

			So(*responseBuilder.etag, ShouldEqual, etagValue)
		})
	})

	Convey("WithCacheControl", t, func() {
		cacheControl := headers.CacheControlNoStore

		Convey("Should set cache control header", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithCacheControl(cacheControl)
			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.headers, ShouldNotBeEmpty)
			So(responseBuilder.headers, ShouldContainKey, headers.HeaderCacheControl)
			actualValue := responseBuilder.headers[headers.HeaderCacheControl]
			So(actualValue, ShouldEqual, cacheControl.String())
		})
	})

	Convey("WithStatusCode", t, func() {
		statusCode := 200

		Convey("Should set status code", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithStatusCode(statusCode)
			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.statusCode, ShouldEqual, statusCode)
		})
	})

	Convey("WithBody", t, func() {
		body := responseBodyTestStruct{
			ID:     "test",
			Number: 1234,
		}

		contentType := headers.ContentTypeJSON

		Convey("Should set body", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithBody(contentType, body)
			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.body, ShouldNotBeNil)
			So(responseBuilder.body.contentType, ShouldEqual, contentType)
			So(responseBuilder.body.body, ShouldEqual, body)
		})
	})

	Convey("WithJSONBody", t, func() {
		body := responseBodyTestStruct{
			ID:     "test",
			Number: 1234,
		}
		contentType := headers.ContentTypeJSON

		Convey("Should set body", func() {
			httpResponseBuilder := CreateHTTPResponseBuilder().WithJSONBody(body)
			responseBuilder := httpResponseBuilder.(*ResponseBuilder)

			So(responseBuilder.body, ShouldNotBeNil)
			So(responseBuilder.body.contentType, ShouldEqual, contentType)
			So(responseBuilder.body.body, ShouldEqual, body)
		})
	})

	Convey("Build", t, func() {
		etagValue := "test-etag-value"
		cacheControl := headers.CacheControlNoStore
		testHeaders := map[string]string{
			"test-header-one":   "test-header-value",
			"other-test-header": "expected-test-header-value",
		}
		statusCode := 400
		responseBody := responseBodyTestStruct{
			ID:     "test",
			Number: 1234,
		}

		ctx := context.Background()

		Convey("Should set response etag header", func() {
			rr := httptest.NewRecorder()
			err := CreateHTTPResponseBuilder().WithETag(etagValue).Build(ctx, rr)
			So(err, ShouldBeNil)
			So(rr.Header().Get(headers.HeaderETag), ShouldEqual, etagValue)
		})

		Convey("Should set response cache-control header", func() {
			rr := httptest.NewRecorder()
			err := CreateHTTPResponseBuilder().WithCacheControl(cacheControl).Build(ctx, rr)

			So(err, ShouldBeNil)
			So(rr.Header().Get(headers.HeaderCacheControl), ShouldEqual, cacheControl.String())
		})

		Convey("Should set response headers", func() {
			rr := httptest.NewRecorder()
			responseBuilder := CreateHTTPResponseBuilder()

			for key, value := range testHeaders {
				responseBuilder = responseBuilder.WithHeader(key, value)
			}

			err := responseBuilder.Build(ctx, rr)

			So(err, ShouldBeNil)
			for key, value := range testHeaders {
				So(rr.Header().Get(key), ShouldEqual, value)
			}
		})

		Convey("Should set status code", func() {
			rr := httptest.NewRecorder()
			err := CreateHTTPResponseBuilder().WithStatusCode(statusCode).Build(ctx, rr)
			So(err, ShouldBeNil)
			So(rr.Result().StatusCode, ShouldEqual, statusCode)
		})

		Convey("Should set response body + content-type header", func() {
			rr := httptest.NewRecorder()

			mockResponder := CreateMockResponder()
			err := CreateHTTPResponseBuilderWithResponder(&mockResponder).
				WithStatusCode(statusCode).
				WithJSONBody(responseBody).
				Build(ctx, rr)

			So(err, ShouldBeNil)
			So(mockResponder.calls.JSON, ShouldNotBeEmpty)
			So(mockResponder.calls.JSON, ShouldHaveLength, 1)
			So(mockResponder.calls.JSON[0].Status, ShouldEqual, statusCode)
			So(mockResponder.calls.JSON[0].Resp, ShouldEqual, responseBody)
		})

		Convey("Should set all", func() {
			rr := httptest.NewRecorder()
			mockResponder := CreateMockResponder()

			builder := CreateHTTPResponseBuilderWithResponder(&mockResponder).
				WithCacheControl(cacheControl).
				WithStatusCode(statusCode).
				WithETag(etagValue).
				WithJSONBody(responseBody)

			for key, value := range testHeaders {
				builder = builder.WithHeader(key, value)
			}

			err := builder.Build(ctx, rr)

			So(err, ShouldBeNil)

			for key, value := range testHeaders {
				So(rr.Header().Get(key), ShouldEqual, value)
			}

			So(mockResponder.calls.JSON, ShouldNotBeEmpty)
			So(mockResponder.calls.JSON, ShouldHaveLength, 1)
			So(mockResponder.calls.JSON[0].Status, ShouldEqual, statusCode)
			So(mockResponder.calls.JSON[0].Resp, ShouldEqual, responseBody)

			So(rr.Result().StatusCode, ShouldEqual, statusCode)

			for key, value := range testHeaders {
				So(rr.Header().Get(key), ShouldEqual, value)
			}

			So(rr.Header().Get(headers.HeaderETag), ShouldEqual, etagValue)
			So(rr.Header().Get(headers.HeaderCacheControl), ShouldEqual, cacheControl.String())
		})
	})
}
