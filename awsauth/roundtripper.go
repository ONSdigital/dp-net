package awsauth

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
)

type AwsSignerRoundTripper struct {
	signer       *Signer
	roundTripper http.RoundTripper
}

type Options struct {
	// InsecureSkipVerify controls whether a client verifies the server's certificate
	// chain and host name. If InsecureSkipVerify is true, crypto/tls accepts any
	// certificate presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to machine-in-the-middle attacks unless custom
	// verification is used. This should be used only for testing or in combination
	// with VerifyConnection or VerifyPeerCertificate.
	TlsInsecureSkipVerify bool
}

var defaultAWSTransport = dphttp.DefaultTransport

func NewAWSSignerRoundTripper(ctx context.Context, awsFilename, awsProfile, awsRegion, awsService string, options ...Options) (*AwsSignerRoundTripper, error) {
	if awsRegion == "" || awsService == "" {
		return nil, fmt.Errorf("aws region and service should be valid options")
	}

	// Create the AWS signer
	awsSigner, err := NewAwsSigner(ctx, awsFilename, awsProfile, awsRegion, awsService)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws v4 signer: %w", err)
	}

	// Create the transport and modify TLS settings if needed
	transport := http.DefaultTransport.(*http.Transport).Clone() // Clone the default transport
	if len(options) > 0 && options[0].TlsInsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Return the custom round tripper
	return &AwsSignerRoundTripper{
		signer:       awsSigner,
		roundTripper: transport,
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
