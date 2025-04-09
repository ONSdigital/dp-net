package awsauth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signerV4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Signer struct {
	awsRegion  string
	awsService string
	v4         *signerV4.Signer
	creds      aws.CredentialsProvider
}

func NewAwsSigner(ctx context.Context, awsFilename, awsProfile, awsRegion, awsService string) (signer *Signer, err error) {
	if err = validateAwsSDKSigner(awsRegion, awsService); err != nil {
		return
	}

	// Initialize options for the AWS configuration
	var configOpts []func(*config.LoadOptions) error
	configOpts = append(configOpts, config.WithRegion(awsRegion)) // Always set the region

	// Only add shared config file and profile if they are provided
	if awsFilename != "" {
		configOpts = append(configOpts, config.WithSharedConfigFiles([]string{awsFilename}))
	}
	if awsProfile != "" {
		configOpts = append(configOpts, config.WithSharedConfigProfile(awsProfile))
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, configOpts...)
	if err != nil {
		// Proceed by attempting to load the config without a custom file/profile (i.e., use defaults)
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
		if err != nil {
			return nil, fmt.Errorf("unable to load AWS config even with default credentials: %w", err)
		}
	}

	// Create the signer
	signer = &Signer{
		awsRegion:  awsRegion,
		awsService: awsService,
		v4:         signerV4.NewSigner(),
		creds:      cfg.Credentials,
	}

	return signer, nil
}

func (s *Signer) Sign(req *http.Request, bodyReader io.ReadSeeker, currentTime time.Time) error {
	if s == nil || s.v4 == nil {
		return errors.New("v4 signer missing. Cannot sign request")
	}

	var payloadHash string
	var err error

	if bodyReader != nil {
		payloadHash, err = hashBody(bodyReader)
		if err != nil {
			return err
		}

		// Ensure body is set after hashing
		req.Body = io.NopCloser(bodyReader)
	}

	// Retrieve fresh credentials on every sign
	creds, err := s.creds.Retrieve(req.Context())
	if err != nil {
		return err
	}

	// Sign
	err = s.v4.SignHTTP(req.Context(), creds, req, payloadHash, s.awsService, s.awsRegion, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func validateAwsSDKSigner(awsRegion, awsService string) error {
	if awsRegion == "" {
		return errors.New("no AWS region was provided. Cannot sign request")
	}

	if awsService == "" {
		return errors.New("no AWS service was provided. Cannot sign request")
	}

	return nil
}

// hashBody computes the SHA-256 hash of a bodyReader and resets it
func hashBody(body io.ReadSeeker) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, body); err != nil {
		return "", err
	}
	sum := hasher.Sum(nil)

	// Reset the reader so the request can read from it again
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	return hex.EncodeToString(sum), nil
}
