package cloudauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/cli/cli/config/types"
	"golang.org/x/oauth2/google"
)

const (
	GCPCredHelper  = "gcp-login"
	oauthScope     = "https://www.googleapis.com/auth/cloud-platform"
	oauthTokenUser = "oauth2accesstoken"
)

func (ap *authProvider) getAuthConfigGCP(ctx context.Context, host, org, project string) (*authConfig, error) {
	gcpJSONPath := getRegistrySecret(host, org, project, "GCP_KEY")
	credHelperPath := getRegistrySecret(host, org, project, "cred_helper")
	registryPath := getRegistrySecretPrefix(host, org, project)

	ap.console.VerbosePrintf("looking up %s", gcpJSONPath)
	gcpJSONSecret, err := ap.cloudClient.GetUserOrProjectSecret(ctx, gcpJSONPath)
	if err != nil {
		return nil, err
	}
	gcpJSON := strings.TrimSpace(gcpJSONSecret.Value)
	if gcpJSON == "" {
		return nil, fmt.Errorf("%s is missing (or empty), but %s was set to %s", gcpJSONPath, credHelperPath, GCPCredHelper)
	}

	ap.console.VerbosePrintf("creating a new JWT config for %s", host)

	jwtCFG, err := google.JWTConfigFromJSON([]byte(gcpJSON), oauthScope)
	if err != nil {
		return nil, fmt.Errorf("failed to get a jwt cfg from the service account json in %s", gcpJSONPath)
	}

	token, err := jwtCFG.TokenSource(ctx).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get a new token using credentials from %s failed: %w", registryPath, err)
	}
	ap.console.VerbosePrintf("%s succeeded using gcloud credentials stored under %s", GCPCredHelper, registryPath)
	cfg := &authConfig{
		ac: &types.AuthConfig{
			Username:      oauthTokenUser,
			Password:      token.AccessToken,
			ServerAddress: host,
		},
		loc: registryPath,
	}
	return cfg, nil
}
