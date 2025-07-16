package secretprovider

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/moby/buildkit/session/secrets"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"

	"github.com/earthly/earthly/util/hint"
	"github.com/earthly/earthly/util/oidcutil"
)

// Internal reserved credentials names used to acquire the equivalent values
// from the environment.
const (
	awsAccessKey    = "aws:access_key"
	awsSecretKey    = "aws:secret_key"
	awsSessionToken = "aws:session_token"
	awsRegion       = "aws:region"

	roleARNURLParam         = "role-arn"
	regionURLParam          = "region"
	sessionDurationURLParam = "session-duration"
	sessionNameURLParam     = "session-name"
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

// AWSCredentialProvider can load AWS settings from the environment or oidc provider
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

	cfg, err := getCFG(ctx)
	if err != nil {
		return nil, err
	}
	creds, err := cfg.Credentials.Retrieve(ctx)

	if err = handleError(err, cfg.Region); err != nil {
		return nil, err
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

	if val == "" && secretName != awsRegion {
		// the region may be provided by a separate arg/env if it's not provided in the local env or oidc configuration
		// Use a custom error here as not to fall back on other secret providers.
		return nil, errors.Errorf("AWS setting %s not found in environment", secretName)
	}

	return []byte(val), nil
}

// getCFG returns a configuration that can provide credentials and region
// The cfg is host environment based (e.g. ~/.aws).
func getCFG(ctx context.Context) (aws.Config, error) {
	// Get the secrets from the host environment
	// Note: results of this call are cached.
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultsMode(aws.DefaultsModeStandard))
	if err != nil {
		return aws.Config{}, errors.Wrap(err, "failed to load AWS config")
	}
	return cfg, nil
}

// SetURLValuesFunc returs a function that takes url.Values and sets oidc values.
// This is used by SecretID() to be able to identify secrets from this provider
func SetURLValuesFunc(awsInfo *oidcutil.AWSOIDCInfo) func(values url.Values) {
	return func(values url.Values) {
		values.Set(sessionNameURLParam, awsInfo.SessionName)
		values.Set(roleARNURLParam, awsInfo.RoleARN.String())
		values.Set(regionURLParam, awsInfo.Region)
		if awsInfo.SessionDuration != nil {
			values.Set(sessionDurationURLParam, awsInfo.SessionDuration.String())
		}
	}
}

func handleError(err error, region string) error {
	if err == nil {
		return nil
	}
	if grpcErr, ok := grpcerrors.AsGRPCStatus(err); ok {
		switch grpcErr.Code() {
		case codes.InvalidArgument:
			return hint.Wrapf(err, `is %q a valid AWS region?`, region)
		}
	}
	return errors.Wrap(err, "failed to load AWS credentials")
}
