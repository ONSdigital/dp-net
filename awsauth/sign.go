package awsauth

import (
	"context"
	"errors"
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

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(awsRegion),
		config.WithSharedConfigFiles([]string{awsFilename}), // Ensure awsFilename is correct
		config.WithSharedConfigProfile(awsProfile),          // Ensure awsProfile exists in the config
	)
	if err != nil {
		return nil, err
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

	credentials, err := s.creds.Retrieve(context.TODO())
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
