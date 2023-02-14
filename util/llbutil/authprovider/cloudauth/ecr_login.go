package cloudauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/docker/cli/cli/config/types"
)

func (ap *authProvider) getAuthConfigECR(ctx context.Context, host, org, project string) (*authConfig, error) {
	registry, err := api.ExtractRegistry(host)
	if err != nil {
		return nil, fmt.Errorf("failed to extract aws ecr registry details from %s: %w", host, err)
	}

	if registry.FIPS {
		return nil, fmt.Errorf("FIPS registry %s not supported", host)
	}

	accessKeyIDPath := getRegistrySecret(host, org, project, "AWS_ACCESS_KEY_ID")
	secretAccessKeyPath := getRegistrySecret(host, org, project, "AWS_SECRET_ACCESS_KEY")
	credHelperPath := getRegistrySecret(host, org, project, "cred_helper")
	registryPath := getRegistrySecretPrefix(host, org, project)

	ap.console.VerbosePrintf("looking up %s", accessKeyIDPath)
	accessKeyIDSecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, accessKeyIDPath)
	if err != nil {
		return nil, err
	}
	accessKeyID := strings.TrimSpace(accessKeyIDSecret.Value)
	if accessKeyID == "" {
		return nil, fmt.Errorf("%s is missing (or empty), but %s was set to ecr-login", accessKeyIDPath, credHelperPath)
	}

	ap.console.VerbosePrintf("looking up %s", secretAccessKeyPath)
	secretAccessKeySecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, secretAccessKeyPath)
	if err != nil {
		return nil, err
	}
	secretAccessKey := strings.TrimSpace(secretAccessKeySecret.Value)
	if secretAccessKey == "" {
		return nil, fmt.Errorf("%s is missing (or empty), but %s was set to ecr-login", secretAccessKeyPath, credHelperPath)
	}

	clientFactory := api.DefaultClientFactory{}
	client := clientFactory.NewClient(aws.Config{
		Region: registry.Region,
		Credentials: awsAccessKeyCredentials{
			accessKeyID:     accessKeyID,
			secretAccessKey: secretAccessKey,
		},
	})

	ap.console.VerbosePrintf("calling ecr-login GetCredentials for %s", host)
	auth, err := client.GetCredentials(host)
	if err != nil {
		return nil, fmt.Errorf("ecr-login using credentials from %s failed: %w", registryPath, err)
	}

	ap.console.VerbosePrintf("ecr-login succeeded using aws credentials stored under %s", registryPath)
	return &authConfig{
		ac: &types.AuthConfig{
			ServerAddress: host,
			Username:      auth.Username,
			Password:      auth.Password,
		},
		loc: registryPath,
	}, nil
}

type awsAccessKeyCredentials struct {
	accessKeyID     string
	secretAccessKey string
}

// Retrieve implements the CredentialsProvider interface, and only allows for auth based on accessKeyID and secretAccessKey
func (hc awsAccessKeyCredentials) Retrieve(context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     hc.accessKeyID,
		SecretAccessKey: hc.secretAccessKey,
	}, nil
}
