package http

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	awsAuth "github.com/ONSdigital/dp-elasticsearch/v2/awsauth"
)

type awsSignerRoundTripper struct {
	signer       *awsAuth.Signer
	roundTripper http.RoundTripper
}

func NewAWSSignerRoundTripper(awsFilename, awsProfile, awsRegion, awsService string, customTransport http.RoundTripper) http.RoundTripper {
	var roundTripper http.RoundTripper
	if awsRegion == "" || awsService == "" {
		log.Fatal(context.Background(), "aws region and service should be valid options")
	}
	awsSigner, err := awsAuth.NewAwsSigner(awsFilename, awsProfile, awsRegion, awsService)
	if err != nil {
		log.Fatal(context.Background(), "failed to create aws v4 signer", err)
	}

	if customTransport == nil {
		roundTripper = http.DefaultTransport
	} else {
		roundTripper = customTransport
	}

	return &awsSignerRoundTripper{
		signer:       awsSigner,
		roundTripper: roundTripper,
	}
}

func (srt *awsSignerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}
	if err := srt.signer.Sign(req, bytes.NewReader(body), time.Now()); err != nil {
		return nil, err
	}

	return srt.roundTripper.RoundTrip(req)
}
