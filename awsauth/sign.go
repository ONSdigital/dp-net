package awsauth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	signerV4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
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

	// Create credentials cache using a custom credentials provider
	creds := aws.NewCredentialsCache(
		aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			// First try retrieving credentials from the loaded config
			if cfg.Credentials != nil {
				creds, err := cfg.Credentials.Retrieve(ctx)
				if err == nil {
					return creds, nil
				}
			}

			// Fallback to EC2 Role credentials if no valid credentials were found
			ec2Provider := ec2rolecreds.New()
			ec2Creds, err := ec2Provider.Retrieve(ctx)
			if err == nil {
				return ec2Creds, nil
			}

			// Return error if no valid credential provider found
			return aws.Credentials{}, errors.New("no valid credential provider found")
		}),
	)

	// Create the signer
	signer = &Signer{
		awsRegion:  awsRegion,
		awsService: awsService,
		v4:         signerV4.NewSigner(),
		creds:      creds,
	}

	return signer, nil
}

func (s *Signer) Sign(req *http.Request, bodyReader io.ReadSeeker, currentTime time.Time) error {
	if s == nil || s.v4 == nil {
		return errors.New("v4 signer missing. Cannot sign request")
	}

	credentials, err := s.creds.Retrieve(req.Context())
	if err != nil {
		return err
	}

	return s.v4.SignHTTP(req.Context(), credentials, req, "", s.awsService, s.awsRegion, currentTime)
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
