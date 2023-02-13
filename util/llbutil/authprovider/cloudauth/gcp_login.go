package cloudauth

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/cli/cli/config/types"
	"golang.org/x/oauth2/google"
)

const (
	OauthScope     = "https://www.googleapis.com/auth/cloud-platform"
	OauthTokenUser = "oauth2accesstoken"
)

func (ap *authProvider) getAuthConfigGCR(ctx context.Context, fullPathPrefix, pathPrefix, org, project, host string) (*authConfig, error) {
	jsonKeyPath := pathPrefix + "GCP_KEY"
	ap.console.VerbosePrintf("looking up GCP_KEY", jsonKeyPath)
	keyJSON, err := ap.getProjectOrUserSecret(ctx, org, project, jsonKeyPath)
	if err != nil {
		return nil, err
	}
	keyJSON = strings.TrimSpace(keyJSON)
	if keyJSON == "" {
		return nil, fmt.Errorf("%sGCP_KEY is missing (or empty), but %scred_helper was set to gcp-login", fullPathPrefix, fullPathPrefix)
	}

	ap.console.VerbosePrintf("creating new a new JWT config for %s", host)

	jwtCFG, err := google.JWTConfigFromJSON([]byte(keyJSON), OauthScope)

	token, err := jwtCFG.TokenSource(ctx).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get a new token using credentials from %s failed: %w", fullPathPrefix, err)
	}
	ap.console.VerbosePrintf("gcp-login succeeded using gcloud credentials stored under %s", fullPathPrefix)
	cfg := &authConfig{
		ac: &types.AuthConfig{
			Username:      OauthTokenUser,
			Password:      token.AccessToken,
			ServerAddress: host,
		},
		loc: fullPathPrefix,
	}
	return cfg, nil
}
