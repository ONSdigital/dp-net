package http

import (
	"bytes"
	"fmt"
	awsAuth "github.com/ONSdigital/dp-net/awsauth"
	"io"
	"net/http"
	"time"
)

type AwsSignerRoundTripper struct {
	signer       *awsAuth.Signer
	roundTripper http.RoundTripper
}

func NewAWSSignerRoundTripper(awsFilename, awsProfile, awsRegion, awsService string, customTransport http.RoundTripper) (*AwsSignerRoundTripper, error) {
	fmt.Println("inside aws signer..........................")
	var roundTripper http.RoundTripper
	if awsRegion == "" || awsService == "" {
		return nil, fmt.Errorf("aws region and service should be valid options")
	}
	awsSigner, err := awsAuth.NewAwsSigner(awsFilename, awsProfile, awsRegion, awsService)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws v4 signer: %w", err)
	}

	if customTransport == nil {
		roundTripper = http.DefaultTransport
	} else {
		roundTripper = customTransport
	}

	return &AwsSignerRoundTripper{
		roundTripper: roundTripper,
		signer:       awsSigner,
	}, nil
}

func (srt *AwsSignerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println("inside round tripper..........................")
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
