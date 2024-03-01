package secretprovider

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

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
		"AWS_ACCESS_KEY_ID":     {},
		"AWS_SECRET_ACCESS_KEY": {},
		"AWS_SESSION_TOKEN":     {},
		"AWS_DEFAULT_REGION":    {},
	}

	q, err := url.ParseQuery(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse secret info")
	}

	secretName := q.Get("name")

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

	return nil, secrets.ErrNotFound
}

func (c *AWSCredentialProvider) loadFromEnv(ctx context.Context, name string) (string, bool) {
	return os.LookupEnv(name)
}

func (c *AWSCredentialProvider) loadFromConfig(ctx context.Context, name string) (string, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config == nil || c.credsConfig == nil {
		homeDir, _ := fileutil.HomeDir()

		var err error

		configFile := filepath.Join(homeDir, ".aws", "config")
		c.config, err = ini.Load(configFile)
		if err != nil {
			return "", false, errors.Wrapf(err, "failed to load %s", configFile)
		}

		credsFile := filepath.Join(homeDir, ".aws", "credentials")
		c.credsConfig, err = ini.Load(credsFile)
		if err != nil {
			return "", false, errors.Wrap(err, "failed to decode credentials file")
		}
	}

	switch name {
	case "AWS_ACCESS_KEY_ID":
		v, ok := iniKey(c.credsConfig, "default", "aws_access_key_id")
		return v, ok, nil
	case "AWS_SECRET_ACCESS_KEY":
		v, ok := iniKey(c.credsConfig, "default", "aws_secret_access_key")
		return v, ok, nil
	case "AWS_DEFAULT_REGION":
		v, ok := iniKey(c.config, "default", "region")
		return v, ok, nil
	default:
		return "", false, nil
	}
}

func iniKey(f *ini.File, section, name string) (string, bool) {
	if !f.HasSection(section) {
		return "", false
	}

	s := f.Section(section)

	if !s.HasKey(name) {
		return "", false
	}

	return s.Key(name).Value(), true
}
