package awsauth

import (
	"context"
	"errors"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	testAccessKey       = "TEST_ACCESS_KEY"
	testSecretAccessKey = "TEST_SECRET_KEY"
	testCredFile        = "testdata/.aws/credentials"
	testProfile         = "default"
)

var ctx = context.Background()

func TestCreateNewSigner(t *testing.T) {
	Convey("Given that we want to create the aws sdk signer", t, func() {
		Convey("When the region is set to an empty string", func() {
			Convey("Then an error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner(ctx, "", "", "", "es")
				So(err, ShouldResemble, errors.New("no AWS region was provided. Cannot sign request"))
				So(signer, ShouldBeNil)
			})
		})

		Convey("When the service is set to an empty string", func() {
			Convey("Then an error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner(ctx, "", "", "eu-west-1", "")
				So(err, ShouldResemble, errors.New("no AWS service was provided. Cannot sign request"))
				So(signer, ShouldBeNil)
			})
		})

		Convey("When the service and region are set and credentials are set in environment variables", func() {
			accessKeyID, secretAccessKey := setEnvironmentVars()

			Convey("Then no error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner(ctx, "", "", "eu-west-1", "es")
				So(err, ShouldBeNil)
				So(signer, ShouldNotBeNil)

				Convey("And no error is returned when attempting to Sign the request", func() {
					req := httptest.NewRequest("GET", "http://test-url", nil)

					err := signer.Sign(req, nil, time.Now())
					So(err, ShouldBeNil)
				})
			})

			removeTestEnvironmentVariables(accessKeyID, secretAccessKey)
		})
	})
}

func TestSignFunc(t *testing.T) {
	Convey("Given that we want to use the aws sdk signer to sign request", t, func() {
		Convey("When the signer is nil", func() {
			Convey("Then an error is returned when attempting to Sign the request", func() {
				var signer *Signer
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldResemble, errors.New("v4 signer missing. Cannot sign request"))
			})
		})

		Convey("When the signer.v4 is nil", func() {
			Convey("Then an error is returned when attempting to Sign the request", func() {
				signer := &Signer{
					v4: nil,
				}
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldResemble, errors.New("v4 signer missing. Cannot sign request"))
			})
		})

		Convey("When the signer.v4 is a valid aws v4 signer", func() {
			// Create valid v4 signer
			accessKeyID, secretAccessKey := setEnvironmentVars()

			signer, err := NewAwsSigner(ctx, "", "", "eu-west-1", "es")
			So(err, ShouldBeNil)
			So(signer, ShouldNotBeNil)
			So(signer.v4, ShouldNotBeNil)

			Convey("Then the request successfully signs and does not return an error", func() {

				req := httptest.NewRequest("GET", "http://test-url", nil)

				err = signer.Sign(req, nil, time.Now())
				So(err, ShouldBeNil)
			})

			removeTestEnvironmentVariables(accessKeyID, secretAccessKey)
		})

		Convey("When the signer.v4 is a valid aws file and profile", func() {
			signer, err := NewAwsSigner(ctx, testCredFile, testProfile, "eu-west-1", "es")
			So(err, ShouldBeNil)
			So(signer, ShouldNotBeNil)
			So(signer.v4, ShouldNotBeNil)

			Convey("Then the request successfully signs and does not return an error", func() {

				req := httptest.NewRequest("GET", "http://test-url", nil)

				err = signer.Sign(req, nil, time.Now())
				So(err, ShouldBeNil)
			})
		})
	})
}

func setEnvironmentVars() (accessKeyID, secretAccessKey string) {
	accessKeyID = os.Getenv(envAccessKeyID)
	secretAccessKey = os.Getenv(envSecretAccessKey)

	os.Setenv(envAccessKeyID, testAccessKey)
	os.Setenv(envSecretAccessKey, testSecretAccessKey)

	return
}

func removeTestEnvironmentVariables(accessKeyID, secretAccessKey string) {
	os.Setenv(envAccessKeyID, accessKeyID)
	os.Setenv(envSecretAccessKey, secretAccessKey)
}
