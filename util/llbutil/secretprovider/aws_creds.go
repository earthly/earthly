package secretprovider

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
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
type AWSCredentialProvider struct {
	mu          sync.Mutex
	config      *ini.File
	credsConfig *ini.File
}

// NewAWSCredentialProvider creates and returns a credential provider for AWS.
func NewAWSCredentialProvider() *AWSCredentialProvider {
	return &AWSCredentialProvider{}
}

// GetSecret attempts to find an AWS secret from either the environment or a local config file.
func (c *AWSCredentialProvider) GetSecret(ctx context.Context, name string) ([]byte, error) {

	names := map[string]struct{}{
		awsAccessKey:    {},
		awsSecretKey:    {},
		awsSessionToken: {},
		awsRegion:       {},
	}

	q, err := url.ParseQuery(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse secret info")
	}

	secretName := q.Get("name")

	// This provider only deals with secrets prefixed with "aws:".
	if !strings.HasPrefix(secretName, "aws:") {
		return nil, secrets.ErrNotFound
	}

	if _, ok := names[secretName]; ok {
		if val, ok := c.loadFromEnv(ctx, secretName); ok {
			return []byte(val), nil
		}

		val, ok, err := c.loadFromConfig(ctx, secretName)
		if err != nil {
			return nil, err
		}

		if ok {
			return []byte(val), nil
		}
	}

	// Use a custom error here as not to fall back on other secret providers.
	return nil, errors.Errorf("missing AWS credential: %s", secretName)
}

func (c *AWSCredentialProvider) loadFromEnv(ctx context.Context, name string) (string, bool) {
	envName, ok := awsEnvNames[name]
	if !ok {
		return "", false
	}
	return os.LookupEnv(envName)
}

func (c *AWSCredentialProvider) loadFromConfig(ctx context.Context, name string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config == nil || c.credsConfig == nil {
		homeDir, _ := fileutil.HomeDir()

		var err error

		configFile := filepath.Join(homeDir, ".aws", "config")
		c.config, err = ini.Load(configFile)
		if err != nil && !os.IsNotExist(err) {
			return "", false, errors.Wrapf(err, "failed to load %s", configFile)
		}

		credsFile := filepath.Join(homeDir, ".aws", "credentials")
		c.credsConfig, err = ini.Load(credsFile)
		if err != nil && !os.IsNotExist(err) {
			return "", false, errors.Wrapf(err, "failed to load %s", credsFile)
		}
	}

	switch name {
	case awsAccessKey:
		v, ok := iniKey(c.credsConfig, "default", "aws_access_key_id")
		return v, ok, nil
	case awsSecretKey:
		v, ok := iniKey(c.credsConfig, "default", "aws_secret_access_key")
		return v, ok, nil
	case awsSessionToken:
		v, ok := iniKey(c.credsConfig, "default", "aws_session_token")
		return v, ok, nil
	case awsRegion:
		v, ok := iniKey(c.config, "default", "region")
		return v, ok, nil
	default:
		return "", false, nil
	}
}

func iniKey(f *ini.File, section, name string) (string, bool) {
	if f == nil {
		return "", false
	}

	if !f.HasSection(section) {
		return "", false
	}

	s := f.Section(section)

	if !s.HasKey(name) {
		return "", false
	}

	return s.Key(name).Value(), true
}
