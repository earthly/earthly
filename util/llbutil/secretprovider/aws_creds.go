package secretprovider

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/moby/buildkit/session/secrets"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"

	"github.com/earthly/earthly/cloud"
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
	uuidURLParam            = "uuid"
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

var oidcCredsProviderCache = make(map[string]*aws.Config)
var oidcCredsProviderCacheMU sync.Mutex

// AWSEnvName converts and internal AWS secret name to the equivalent official
// environmental variable.
func AWSEnvName(name string) (string, bool) {
	envName, ok := awsEnvNames[name]
	return envName, ok
}

// AWSCredentialProvider can load AWS settings from the environment or oidc provider
type AWSCredentialProvider struct {
	client *cloud.Client
}

// NewAWSCredentialProvider creates and returns a credential provider for AWS.
func NewAWSCredentialProvider(client *cloud.Client) *AWSCredentialProvider {
	return &AWSCredentialProvider{
		client: client,
	}
}

// GetSecret attempts to find an AWS secret from either the environment or a local config file.
func (c *AWSCredentialProvider) GetSecret(ctx context.Context, name string) ([]byte, error) {

	q, err := url.ParseQuery(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse secret info")
	}

	secretName := q.Get("name")
	orgName := q.Get("org")
	projectName := q.Get("project")

	// This provider only deals with secrets prefixed with "aws:".
	if !strings.HasPrefix(secretName, "aws:") {
		return nil, secrets.ErrNotFound
	}

	oidcInfo := oidcInfoFromValues(q)

	cfg, err := getCFG(ctx, orgName, projectName, oidcInfo, c.client)
	if err != nil {
		return nil, err
	}
	creds, err := cfg.Credentials.Retrieve(ctx)

	if err = handleError(err, oidcInfo.RoleARN.String(), cfg.Region, orgName, projectName); err != nil {
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

func oidcInfoFromValues(values url.Values) *oidcutil.AWSOIDCInfo {
	roleARN := values.Get(roleARNURLParam)
	if roleARN == "" { // no arn implies oidc is not in play
		return nil
	}
	region := values.Get(regionURLParam)
	sessionDuration := values.Get(sessionDurationURLParam)
	// the values are pre validated in the interperter
	parsedARN, _ := arn.Parse(roleARN)
	var duration *time.Duration
	if sessionDuration != "" {
		durVal, _ := time.ParseDuration(sessionDuration)
		duration = &durVal
	}
	return &oidcutil.AWSOIDCInfo{
		RoleARN:         &parsedARN,
		Region:          region,
		SessionDuration: duration,
		Uuid:            values.Get(uuidURLParam),
	}
}

type oidcCredentialsProvider struct {
	client      *cloud.Client
	cache       *aws.Credentials
	oidcInfo    *oidcutil.AWSOIDCInfo
	orgName     string
	projectName string
	cacheMU     sync.Mutex
}

func (p *oidcCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	if p.cache != nil {
		return *p.cache, nil
	}
	p.cacheMU.Lock()
	defer p.cacheMU.Unlock()
	if p.cache != nil {
		return *p.cache, nil
	}
	res, err := p.client.GetAWSCredentials(ctx, p.oidcInfo.RoleARN.String(), p.orgName, p.projectName, p.oidcInfo.Region, p.oidcInfo.SessionDuration)
	if err != nil {
		return aws.Credentials{}, err
	}
	p.cache = &aws.Credentials{
		AccessKeyID:     res.GetCredentials().GetAccessKeyId(),
		SecretAccessKey: res.GetCredentials().GetSecretAccessKey(),
		SessionToken:    res.GetCredentials().GetSessionToken(),
		CanExpire:       true,
		Expires:         res.GetCredentials().GetExpiry().AsTime().UTC(),
	}
	return *p.cache, nil
}

// getCFG returns a configuration that can provide credentials and region
// The cfg is either host environment based (e.g. ~/.aws) or oidc based.
// When based on oidc, it would get a session token from the cloud and cache the result.
// Caching is done so that next calls of GetSecret would get the rest of the matching credentials keys
func getCFG(ctx context.Context, orgName string, projectName string, oidcInfo *oidcutil.AWSOIDCInfo, client *cloud.Client) (aws.Config, error) {
	if oidcInfo == nil {
		// no oidc info implies getting the secrets from the host environment
		// Note: results of this call are cached.
		cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultsMode(aws.DefaultsModeStandard))
		if err != nil {
			return aws.Config{}, errors.Wrap(err, "failed to load AWS config")
		}
		return cfg, nil
	}
	// check if we already have a config for the specified oidc info
	key := oidcInfo.String()
	if cfg, ok := oidcCredsProviderCache[key]; ok {
		return *cfg, nil
	}
	// check one more time, this time with a lock
	oidcCredsProviderCacheMU.Lock()
	defer oidcCredsProviderCacheMU.Unlock()
	if cfg, ok := oidcCredsProviderCache[key]; ok {
		return *cfg, nil
	}
	cfg := &aws.Config{
		Region: oidcInfo.Region,
		Credentials: &oidcCredentialsProvider{
			client:      client,
			oidcInfo:    oidcInfo,
			orgName:     orgName,
			projectName: projectName,
		},
	}
	oidcCredsProviderCache[key] = cfg
	return *cfg, nil
}

// SetURLValuesFunc returs a function that takes url.Values and sets oidc values.
// This is used by SecretID() to be able to identify secrets from this provider
func SetURLValuesFunc(awsInfo *oidcutil.AWSOIDCInfo) func(values url.Values) {
	return func(values url.Values) {
		values.Set(uuidURLParam, awsInfo.Uuid)
		values.Set(roleARNURLParam, awsInfo.RoleARN.String())
		values.Set(regionURLParam, awsInfo.Region)
		if awsInfo.SessionDuration != nil {
			values.Set(sessionDurationURLParam, awsInfo.SessionDuration.String())
		}
	}
}

func handleError(err error, arn, region, orgName, projectName string) error {
	if err == nil {
		return nil
	}
	if grpcErr, ok := grpcerrors.AsGRPCStatus(err); ok {
		switch grpcErr.Code() {
		case codes.InvalidArgument:
			if strings.Contains(grpcErr.Message(), "could not be found") {
				return hint.Wrapf(err, `do the org "%s" and project "%s exist"`, orgName, projectName)
			}
			return hint.Wrapf(err, `is "%s" a valid AWS region?`, region)
		case codes.PermissionDenied, codes.FailedPrecondition:
			return hint.Wrapf(err, `make sure the role %s has a valid trust policy configured in AWS`, arn)
		}
	}
	return errors.Wrap(err, "failed to load AWS credentials")
}
