package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ONSdigital/dp-net/v2/http/httptest"
	request "github.com/ONSdigital/dp-net/v2/request"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHappyPaths(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a default client and happy paths", t, func() {
		httpClient := NewClient()

		Convey("When Get() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.Get(context.Background(), ts.URL)
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees a GET with no body", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "GET")
				So(call.Body, ShouldEqual, "")
				So(call.Error, ShouldEqual, "")
				So(resp.Header.Get("Content-Type"), ShouldContainSubstring, "text/plain")
			})
		})

		Convey("When Post() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.Post(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(`{"dummy":"ook"}`))
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees a POST with that body as JSON", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "POST")
				So(call.Body, ShouldEqual, `{"dummy":"ook"}`)
				So(call.Error, ShouldEqual, "")
				So(call.Headers[httptest.ContentTypeHeader], ShouldResemble, []string{httptest.JsonContentType})
			})
		})

		Convey("When Put() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.Put(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(`{"dummy":"ook2"}`))
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees a PUT with that body as JSON", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "PUT")
				So(call.Body, ShouldEqual, `{"dummy":"ook2"}`)
				So(call.Error, ShouldEqual, "")
				So(call.Headers[httptest.ContentTypeHeader], ShouldResemble, []string{httptest.JsonContentType})
			})
		})

		Convey("When PostForm() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.PostForm(context.Background(), ts.URL, url.Values{"ook": {"koo"}, "zoo": {"ooz"}})
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees a POST with those values encoded", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "POST")
				So(call.Body, ShouldEqual, "ook=koo&zoo=ooz")
				So(call.Error, ShouldEqual, "")
				So(call.Headers[httptest.ContentTypeHeader], ShouldResemble, []string{httptest.FormEncodedType})
			})
		})
	})
}

func TestClientDoesRetry(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with small client timeout", t, func() {
		// force client to abandon requests before the requested one second delay on the (next) server response
		httpClient := ClientWithTimeout(nil, 100*time.Millisecond)

		Convey("When Post() is called on a URL with a delay on the first response", func() {
			delayByOneSecondOnNext := delayByOneSecondOn(expectedCallCount + 1)
			/// XXX this is two for the retry due to the delayed response to first POST
			expectedCallCount += 2
			resp, err := httpClient.Post(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(delayByOneSecondOnNext))
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees two POST calls", func() {
				So(ts.GetCalls(0), ShouldEqual, expectedCallCount)
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "POST")
				So(call.Body, ShouldEqual, delayByOneSecondOnNext)
				So(call.Error, ShouldEqual, "")
				So(resp.Header.Get(httptest.ContentTypeHeader), ShouldContainSubstring, "text/plain")
			})
		})
	})
}

func TestClientDoesRetryAndContextCancellation(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with small client timeout", t, func() {
		// force client to abandon requests before the requested one second delay on the (next) server response
		httpClient := ClientWithTimeout(nil, 500*time.Millisecond)
		Convey("When Post() is called on a URL with a delay on the first response", func() {
			delayByOneSecondOnNext := delayByOneSecondOn(expectedCallCount + 1)
			expectedCallCount++

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()

			resp, err := httpClient.Post(ctx, ts.URL, httptest.JsonContentType, strings.NewReader(delayByOneSecondOnNext))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "context canceled")
			So(resp, ShouldBeNil)

			Convey("Then the server sees two POST calls", func() {
				So(ts.GetCalls(0), ShouldEqual, expectedCallCount)
			})
		})
	})
}

func TestClientDoesRetryAndContextTimeout(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with small client timeout", t, func() {
		// force client to abandon requests before the requested one second delay on the (next) server response
		httpClient := ClientWithTimeout(nil, 500*time.Millisecond)
		Convey("When Post() is called on a URL with a delay on the first response", func() {
			delayByOneSecondOnNext := delayByOneSecondOn(expectedCallCount + 1)
			expectedCallCount++

			ctx, _ := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))

			resp, err := httpClient.Post(ctx, ts.URL, httptest.JsonContentType, strings.NewReader(delayByOneSecondOnNext))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "context deadline exceeded")
			So(resp, ShouldBeNil)

			Convey("Then the server sees two POST calls", func() {
				So(ts.GetCalls(0), ShouldEqual, expectedCallCount)
			})
		})
	})

	Convey("Client with total timeout", t, func() {
		// total timeout = 500msec, sleep handler = 1sec
		httpClient := ClientWithTotalTimeout(nil, 500*time.Millisecond)
		Convey("When Post() is called on a URL with a delay on the first response", func() {
			delayByOneSecondOnNext := delayByOneSecondOn(expectedCallCount + 1)
			expectedCallCount++

			resp, err := httpClient.Post(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(delayByOneSecondOnNext))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "context deadline exceeded")
			So(resp, ShouldBeNil)

			Convey("Then the server sees two POST calls", func() {
				So(ts.GetCalls(0), ShouldEqual, expectedCallCount)
			})
		})
	})
}

func TestClientNoRetries(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with no retries", t, func() {
		httpClient := ClientWithTimeout(nil, 100*time.Millisecond)
		httpClient.SetMaxRetries(0)

		Convey("When Post() is called on a URL with a delay on the first call", func() {
			delayByOneSecondOnNext := delayByOneSecondOn(expectedCallCount + 1)
			resp, err := httpClient.Post(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(delayByOneSecondOnNext))
			Convey("Then the server has no response", func() {
				So(resp, ShouldBeNil)
				So(err.Error(), ShouldContainSubstring, "Timeout exceeded")
			})
		})
	})
}

func TestClientHandlesUnsuccessfulRequests(t *testing.T) {

	Convey("Given a client with no retries", t, func() {
		httpClient := ClientWithTimeout(nil, 5*time.Second)
		httpClient.SetMaxRetries(0)

		Convey("When the server tries to make a request to a service it is unable to connect to", func() {
			ts := httptest.NewTestServer(500)
			defer ts.Close()

			Convey("Then the server responds with a internal server error", func() {
				resp, err := httpClient.Get(context.Background(), ts.URL)

				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 500)
				So(err, ShouldBeNil)

				call, err := unmarshallResp(resp)
				So(err, ShouldBeNil)

				Convey("And the server sees one GET call", func() {
					So(call.CallCount, ShouldEqual, 1)
					So(call.Method, ShouldEqual, "GET")
					So(call.Error, ShouldEqual, "")
					So(resp.Header.Get(httptest.ContentTypeHeader), ShouldContainSubstring, "text/plain")
				})
			})
		})

		Convey("When the server tries to make a request to a service that currently denying its services", func() {
			ts := httptest.NewTestServer(429)
			defer ts.Close()

			Convey("Then the server responds with too many requests", func() {
				resp, err := httpClient.Get(context.Background(), ts.URL)

				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 429)
				So(err, ShouldBeNil)

				call, err := unmarshallResp(resp)
				So(err, ShouldBeNil)

				Convey("And the server sees one GET call", func() {
					So(call.CallCount, ShouldEqual, 1)
					So(call.Method, ShouldEqual, "GET")
					So(call.Error, ShouldEqual, "")
					So(resp.Header.Get(httptest.ContentTypeHeader), ShouldContainSubstring, "text/plain")
				})
			})
		})
	})

	Convey("Given a client with retries", t, func() {
		httpClient := ClientWithTimeout(nil, 5*time.Second)
		httpClient.SetMaxRetries(1)

		Convey("When the server tries to make a request to a service it is unable to"+
			"connect to and is a path that should not handle retries", func() {
			ts := httptest.NewTestServer(500)
			defer ts.Close()

			path := "/testing"
			httpClient.SetPathsWithNoRetries([]string{path})

			Convey("Then the server responds with a internal server error", func() {
				resp, err := httpClient.Get(context.Background(), ts.URL+path)

				So(resp, ShouldNotBeNil)
				So(resp.StatusCode, ShouldEqual, 500)
				So(err, ShouldBeNil)

				call, err := unmarshallResp(resp)
				So(err, ShouldBeNil)

				Convey("And the server sees one GET call", func() {
					So(call.CallCount, ShouldEqual, 1)
					So(call.Method, ShouldEqual, "GET")
					So(call.Path, ShouldEqual, path)
					So(call.Error, ShouldEqual, "")
					So(resp.Header.Get(httptest.ContentTypeHeader), ShouldContainSubstring, "text/plain")
				})
			})
		})
	})
}

func TestClientAddsRequestIDHeader(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with no correlation ID in context", t, func() {
		// throw in a check for wrapped client instantiation
		httpClient := ClientWithTimeout(nil, 5*time.Second)

		Convey("When Post() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.Post(context.Background(), ts.URL, httptest.JsonContentType, strings.NewReader(`{"hello":"there"}`))
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees the auth header", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "POST")
				So(call.Body, ShouldEqual, `{"hello":"there"}`)
				So(call.Error, ShouldEqual, "")
				So(len(call.Headers[request.RequestHeaderKey]), ShouldEqual, 1)
				So(len(call.Headers[request.RequestHeaderKey][0]), ShouldEqual, 20)
			})
		})
	})
}

func TestClientAppendsRequestIDHeader(t *testing.T) {
	ts := httptest.NewTestServer(200)
	defer ts.Close()
	expectedCallCount := 0

	Convey("Given a client with existing correlation ID in context", t, func() {
		upstreamRequestID := "call1234"
		// throw in a check for wrapped client instantiation
		httpClient := ClientWithTimeout(nil, 5*time.Second)

		Convey("When Post() is called on a URL", func() {
			expectedCallCount++
			resp, err := httpClient.Post(request.WithRequestId(context.Background(), upstreamRequestID), ts.URL, httptest.JsonContentType, strings.NewReader(`{}`))
			So(resp, ShouldNotBeNil)
			So(err, ShouldBeNil)

			call, err := unmarshallResp(resp)
			So(err, ShouldBeNil)

			Convey("Then the server sees the auth header", func() {
				So(call.CallCount, ShouldEqual, expectedCallCount)
				So(call.Method, ShouldEqual, "POST")
				So(call.Error, ShouldEqual, "")
				So(len(call.Headers[request.RequestHeaderKey]), ShouldEqual, 1)
				So(call.Headers[request.RequestHeaderKey][0], ShouldStartWith, upstreamRequestID+",")
				So(len(call.Headers[request.RequestHeaderKey][0]), ShouldBeGreaterThan, len(upstreamRequestID)*3/2)
			})
		})
	})
}

func TestSetPathsWithNoRetries(t *testing.T) {
	client := NewClient()
	Convey("Successfully create map of paths when SetPathsWithNoRetries is called", t, func() {
		client.SetPathsWithNoRetries([]string{"/health", "/healthcheck"})
		paths := client.GetPathsWithNoRetries()
		sort.Strings(paths) // cannot guarentee order of paths
		So(len(paths), ShouldEqual, 2)
		So(paths[0], ShouldEqual, "/health")
		So(paths[1], ShouldEqual, "/healthcheck")
	})

	Convey("Successfully update client with map of paths with ClientWithListOfNonRetriablePaths", t, func() {
		ClientWithListOfNonRetriablePaths(client, []string{"/test"})
		paths := client.GetPathsWithNoRetries()
		So(len(paths), ShouldEqual, 1)
		So(paths[0], ShouldEqual, "/test")
	})
}

func TestNewClientWithTransport(t *testing.T) {
	t.Parallel()

	Convey("Given a custom http transport", t, func() {
		customTransport := DefaultTransport
		customTransport.IdleConnTimeout = 30 * time.Second

		Convey("And a new http client is created with custom transport", func() {
			httpClient := NewClientWithTransport(customTransport)

			ts := httptest.NewTestServer(200)
			expectedCallCount := 0
			Convey("When Get() is called on a URL", func() {
				expectedCallCount++
				resp, err := httpClient.Get(context.Background(), ts.URL)
				So(err, ShouldBeNil)
				So(resp, ShouldNotBeNil)

				call, err := unmarshallResp(resp)
				So(err, ShouldBeNil)

				Convey("Then the server sees a GET with no body", func() {
					So(call.CallCount, ShouldEqual, expectedCallCount)
					So(call.Method, ShouldEqual, "GET")
					So(call.Body, ShouldEqual, "")
					So(call.Error, ShouldEqual, "")
					So(resp.Header.Get("Content-Type"), ShouldContainSubstring, "text/plain")
				})
			})
		})
	})
}

// end of tests //

// delayByOneSecondOn returns the json which will instruct the server to delay responding on call-number `delayOnCall`
func delayByOneSecondOn(delayOnCall int) string {
	return `{"delay":"1s","delay_on_call":` + strconv.Itoa(delayOnCall) + `}`
}

func unmarshallResp(resp *http.Response) (*httptest.Responder, error) {
	responder := &httptest.Responder{}
	body := httptest.GetBody(resp.Body)
	err := json.Unmarshal(body, responder)
	if err != nil {
		panic(err.Error() + string(body))
	}
	return responder, err
}
