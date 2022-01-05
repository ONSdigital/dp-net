package http_test

import (
	"testing"

	"github.com/ONSdigital/dp-net/v2/http"
	"github.com/stretchr/testify/assert"
)

func TestNewAWSSignerRoundTripper(t *testing.T) {
	t.Parallel()

	awsSignerRT := http.NewAWSSignerRoundTripper("some_filename", "some_profile", "some_region", "some_service", nil)

	assert.NotNilf(t, awsSignerRT, "aws signer roundtripper should  not return nil")
}
