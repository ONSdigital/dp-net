package awsauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	"github.com/ONSdigital/dp-net/v2/http/httptest"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSSignerRoundTripper(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := NewAWSSignerRoundTripper("some_filename", "some_profile", "some_region", "some_service")

	assert.Nil(t, err, "error should be nil")
	assert.NotNilf(t, awsSignerRT, "aws signer roundtripper should  not return nil")
}

func TestNewAWSSignerRoundTripper_WhenAWSRegionIsEmpty_Returns(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := NewAWSSignerRoundTripper("some_filename", "some_profile", "", "some_service")

	assert.NotNil(t, err, "error should not be nil")
	assert.Nil(t, awsSignerRT, "aws signer roundtripper should return nil")
}

func TestNewAWSSignerRoundTripper_WhenAWSServiceIsEmpty_Returns(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := NewAWSSignerRoundTripper("some_filename", "", "some_region", "")

	assert.NotNil(t, err, "error should not be nil")
	assert.Nil(t, awsSignerRT, "aws signer roundtripper should return nil")
}

func TestNewClientWithTransport(t *testing.T) {
	Convey("Given a default client and awsauth round tripper", t, func() {

		awsSignerRT, err := NewAWSSignerRoundTripper("some_filename", "some_profile", "some_region", "some_service")
		if err != nil {
			t.Errorf(fmt.Sprintf("unable to implement roundtripper for test, error: %v", err))
		}

		httpClient := dphttp.NewClientWithTransport(awsSignerRT.roundTripper)

		ts := httptest.NewTestServer(200)
		expectedCallCount := 0
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
	})
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
