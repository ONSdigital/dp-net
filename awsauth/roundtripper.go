package awsauth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AwsSignerRoundTripper struct {
	signer       *Signer
	roundTripper http.RoundTripper
}

var defaultAWSTransport = http.DefaultTransport

func NewAWSSignerRoundTripper(awsFilename, awsProfile, awsRegion, awsService string) (*AwsSignerRoundTripper, error) {

	if awsRegion == "" || awsService == "" {
		return nil, fmt.Errorf("aws region and service should be valid options")
	}

	awsSigner, err := NewAwsSigner(awsFilename, awsProfile, awsRegion, awsService)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws v4 signer: %w", err)
	}

	return &AwsSignerRoundTripper{
		signer:       awsSigner,
		roundTripper: defaultAWSTransport,
	}, nil
}

func (srt *AwsSignerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
	}

	if err := srt.signer.Sign(req, bytes.NewReader(body), time.Now()); err != nil {
		return nil, fmt.Errorf("failed to sign the request: %w", err)
	}

	return srt.roundTripper.RoundTrip(req)
}
