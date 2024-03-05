package secretprovider

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
)

// Internal reserved credentials names used to acquire the equivalent values
// from the environment.
const (
	awsAccessKey    = "aws:access_key"
	awsSecretKey    = "aws:secret_key"
	awsSessionToken = "aws:session_token"
	awsRegion       = "aws:region"
)

// AWSCredentials contains the basic set of credentials that users will need to
// use AWS tools.
var AWSCredentials = []string{
	awsAccessKey,
	awsSecretKey,
	awsSessionToken,
	awsRegion,
}

var awsEnvNames = map[string]string{
	awsAccessKey:    "AWS_ACCESS_KEY_ID",
	awsSecretKey:    "AWS_SECRET_ACCESS_KEY",
	awsSessionToken: "AWS_SESSION_TOKEN",
	awsRegion:       "AWS_REGION",
}

// AWSEnvName converts and internal AWS secret name to the equivalent official
// environmental variable.
func AWSEnvName(name string) (string, bool) {
	envName, ok := awsEnvNames[name]
	return envName, ok
}

// AWSCredentialProvider can load AWS settings from the environment.
type AWSCredentialProvider struct{}

// NewAWSCredentialProvider creates and returns a credential provider for AWS.
func NewAWSCredentialProvider() *AWSCredentialProvider {
	return &AWSCredentialProvider{}
}

// GetSecret attempts to find an AWS secret from either the environment or a local config file.
func (c *AWSCredentialProvider) GetSecret(ctx context.Context, name string) ([]byte, error) {

	q, err := url.ParseQuery(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse secret info")
	}

	secretName := q.Get("name")

	// This provider only deals with secrets prefixed with "aws:".
	if !strings.HasPrefix(secretName, "aws:") {
		return nil, secrets.ErrNotFound
	}

	// By default, the AWS config loader will attempt to query for EC2 instance
	// metadata which can be used to store some config details. Unfortunately
	// this can only be set as an environmental variable.
	os.Setenv("AWS_EC2_METADATA_DISABLED", "true")

	// Note: results of this call are cached.
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultsMode(aws.DefaultsModeStandard))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load AWS config")
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		// Suppress verbose warning about the EC2 metadata feature being disabled.
		if strings.Contains(err.Error(), "AWS_EC2_METADATA_DISABLED") {
			return nil, errors.New("failed to load AWS credentials from environment")
		}
		return nil, errors.Wrap(err, "failed to load AWS credentials")
	}

	var val string

	switch secretName {
	case awsAccessKey:
		val = creds.AccessKeyID
	case awsSecretKey:
		val = creds.SecretAccessKey
	case awsSessionToken:
		val = creds.SessionToken
	case awsRegion:
		val = cfg.Region
	default:
		return nil, errors.Errorf("unexpected secret: %s", secretName)
	}

	if val == "" {
		// Use a custom error here as not to fall back on other secret providers.
		return nil, errors.Errorf("AWS setting %s not found in environmental variables or ~/.aws", secretName)
	}

	return []byte(val), nil
}
