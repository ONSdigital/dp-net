package http_test

import (
	"testing"

	"github.com/ONSdigital/dp-net/v2/http"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSSignerRoundTripper(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := http.NewAWSSignerRoundTripper("some_filename", "some_profile", "some_region", "some_service", nil)

	assert.Nil(t, err, "error should be nil")
	assert.NotNilf(t, awsSignerRT, "aws signer roundtripper should  not return nil")
}

func TestNewAWSSignerRoundTripper_WhenAWSRegionIsEmpty_Returns(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := http.NewAWSSignerRoundTripper("some_filename", "some_profile", "", "some_service", nil)

	assert.NotNil(t, err, "error should not be nil")
	assert.Nil(t, awsSignerRT, "aws signer roundtripper should return nil")
}

func TestNewAWSSignerRoundTripper_WhenAWSServiceIsEmpty_Returns(t *testing.T) {
	t.Parallel()

	awsSignerRT, err := http.NewAWSSignerRoundTripper("some_filename", "", "some_region", "", nil)

	assert.NotNil(t, err, "error should not be nil")
	assert.Nil(t, awsSignerRT, "aws signer roundtripper should return nil")
}
