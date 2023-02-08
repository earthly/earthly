package cloudauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/docker/cli/cli/config/types"
)

func (ap *authProvider) getAuthConfigECR(ctx context.Context, fullPathPrefix, pathPrefix, org, project, host string) (*authConfig, error) {
	registry, err := api.ExtractRegistry(host)
	if err != nil {
		return nil, fmt.Errorf("failed to extract aws ecr registry details from %s: %w", host, err)
	}

	if registry.FIPS {
		return nil, fmt.Errorf("FIPS registry %s not supported", host)
	}

	accessKeyIDPath := pathPrefix + "AWS_ACCESS_KEY_ID"
	secretAccessKeyPath := pathPrefix + "AWS_SECRET_ACCESS_KEY"

	ap.console.VerbosePrintf("looking up %sAWS_ACCESS_KEY_ID", fullPathPrefix)
	accessKeyID, err := ap.getProjectOrUserSecret(ctx, org, project, accessKeyIDPath)
	if err != nil {
		return nil, err
	}
	accessKeyID = strings.TrimSpace(accessKeyID)
	if accessKeyID == "" {
		return nil, fmt.Errorf("%sAWS_ACCESS_KEY_ID is missing (or empty), but %scred_helper was set to ecr-login", fullPathPrefix, fullPathPrefix)
	}

	ap.console.VerbosePrintf("looking up %sAWS_SECRET_ACCESS_KEY", fullPathPrefix)
	secretAccessKey, err := ap.getProjectOrUserSecret(ctx, org, project, secretAccessKeyPath)
	if err != nil {
		return nil, err
	}
	secretAccessKey = strings.TrimSpace(secretAccessKey)
	if secretAccessKey == "" {
		return nil, fmt.Errorf("%sAWS_SECRET_ACCESS_KEY is missing (or empty), but %scred_helper was set to ecr-login", fullPathPrefix, fullPathPrefix)
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
		return nil, fmt.Errorf("ecr-login using credentials from %s failed: %w", fullPathPrefix, err)
	}

	ap.console.VerbosePrintf("ecr-login succeeded using aws credentials stored under %s", fullPathPrefix)
	return &authConfig{
		ac: &types.AuthConfig{
			ServerAddress: host,
			Username:      auth.Username,
			Password:      auth.Password,
		},
		loc: fullPathPrefix,
	}, nil
}
