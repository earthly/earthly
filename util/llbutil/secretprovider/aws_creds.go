package secretprovider

import (
	"context"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
)

type AWSCredentialProvider struct{}

func NewAWSCredentialProvider() *AWSCredentialProvider {
	return &AWSCredentialProvider{}
}

func (c *AWSCredentialProvider) GetSecret(ctx context.Context, name string) ([]byte, error) {

	names := map[string]int{
		"AWS_ACCESS_KEY_ID":     0,
		"AWS_SECRET_ACCESS_KEY": 0,
		"AWS_SESSION_TOKEN":     0,
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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", false, errors.Wrap(err, "failed to determine user home directory")
	}

	credPath := filepath.Join(homeDir, ".aws", "credentials")

	f, err := os.Open(credPath)
	if err != nil {
		return "", false, errors.Wrapf(err, "could not open %s", credPath)
	}

	defer func() {
		_ = f.Close()
	}()

	config := struct {
		Default struct {
			AWSAccessKeyID     string `toml:"aws_access_key_id"`
			AWSSecretAccessKey string `toml:"aws_secret_access_key"`
		}
	}{}

	_, err = toml.NewDecoder(f).Decode(&config)
	if err != nil {
		return "", false, errors.Wrap(err, "failed to decode credentials file")
	}

	switch name {
	case "AWS_ACCESS_KEY_ID":
		return config.Default.AWSAccessKeyID, true, nil
	case "AWS_SECRET_ACCESS_KEY":
		return config.Default.AWSAccessKeyID, true, nil
	default:
		return "", false, nil
	}
}
