package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/earthly/earthly/cloud"
	"github.com/moby/buildkit/session/secrets"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	gcpServiceAccountKeyFlag      = "gcp-service-account-key"
	gcpServiceAccountKeyPathFlag  = "gcp-service-account-key-path"
	gcpServiceAccountKeyStdinFlag = "gcp-service-account-key-stdin"
)

var (
	errMultipleGCPServiceAccountFlags = fmt.Errorf("the --%s --%s --%s flags are mutually exclusive", gcpServiceAccountKeyFlag, gcpServiceAccountKeyPathFlag, gcpServiceAccountKeyStdinFlag)
)

func (app *earthlyApp) registryCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "setup",
			Usage:       "setup, and store, registry credentials in the earthly-cloud *beta*",
			Description: "setup, and store, registry credentials in the earthly-cloud *beta*",
			UsageText: "earthly registry setup [--org <org> --project <project>] [--cred-helper <none|ecr-login|gcloud>] ...\n\n" +
				"username/password based registry (--cred-helper=none):\n" +
				"	earthly registry setup --username <username> --password <password> [<host>]\n" +
				"	earthly registry --org <org> --project <project> setup --username <username> --password <password> [<host>]\n\n" +
				"AWS elastic container registry (--cred-helper=ecr-login):\n" +
				"	earthly registry setup --cred-helper ecr-login --aws-access-key-id <key> --aws-secret-access-key <secret> <host>\n" +
				"	earthly registry --org <org> --project <project> setup --cred-helper ecr-login --aws-access-key-id <key> --aws-secret-access-key <secret> <host>\n\n" +
				"GCP artifact or container registry (--cred-helper=gcloud):\n" +
				"	earthly registry setup --cred-helper gcloud --gcp-key <key> <host>\n" +
				"	earthly registry --org <org> --project <project> setup --cred-helper gcloud --gcp-key <key> <host>\n" +
				"",
			Action: app.actionRegistrySetup,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "cred-helper",
					EnvVars:     []string{"EARTHLY_REGISTRY_CRED_HELPER"},
					Usage:       "Use a credential helper when logging into the registry (ecr-login, gcloud).",
					Required:    false,
					Destination: &app.registryCredHelper,
				},
				&cli.StringFlag{
					Name:        "username",
					EnvVars:     []string{"EARTHLY_REGISTRY_USERNAME"},
					Usage:       "The username to use when logging into the registry.",
					Required:    false,
					Destination: &app.registryUsername,
				},
				&cli.StringFlag{
					Name:        "password",
					EnvVars:     []string{"EARTHLY_REGISTRY_PASSWORD"},
					Usage:       "The password to use when logging into the registry (use --password-stdin to prevent leaking your password via your shell history).",
					Required:    false,
					Destination: &app.registryPassword,
				},
				&cli.BoolFlag{
					Name:        "password-stdin",
					EnvVars:     []string{"EARTHLY_REGISTRY_PASSWORD_STDIN"},
					Usage:       "Read the password from stdin (recommended).",
					Required:    false,
					Destination: &app.registryPasswordStdin,
				},
				&cli.StringFlag{
					Name:        "aws-access-key-id",
					EnvVars:     []string{"AWS_ACCESS_KEY_ID"},
					Usage:       "AWS access key ID to use for elastic-container-registry.",
					Required:    false,
					Destination: &app.awsAccessKeyID,
				},
				&cli.StringFlag{
					Name:        "aws-secret-access-key",
					EnvVars:     []string{"AWS_SECRET_ACCESS_KEY"},
					Usage:       "AWS secret access key to use for elastic-container-registry.",
					Required:    false,
					Destination: &app.awsSecretAccessKey,
				},
				&cli.StringFlag{
					Name:        gcpServiceAccountKeyFlag,
					EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY"},
					Usage:       "GCP service account key to use for artifact or container registry.",
					Required:    false,
					Destination: &app.gcpServiceAccountKey,
				},
				&cli.StringFlag{
					Name:        gcpServiceAccountKeyPathFlag,
					EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY_PATH"},
					Usage:       "path to GCP service account key to use for artifact or container registry.",
					Required:    false,
					Destination: &app.gcpServiceAccountKeyPath,
				},
				&cli.BoolFlag{
					Name:        gcpServiceAccountKeyStdinFlag,
					EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY_STDIN"},
					Usage:       "GCP service account key to use for artifact or container registry, read from stdin.",
					Required:    false,
					Destination: &app.gcpServiceAccountKeyStdin,
				},
			},
		},
		{
			Name:  "list",
			Usage: "List configured registries *beta*",
			UsageText: "earthly registry list\n" +
				"	earthly registry --org <org> --project <project> list\n",
			Action: app.actionRegistryList,
		},
		{
			Name:  "remove",
			Usage: "Remove stored registry credentials *beta*",
			UsageText: "earthly registry remove [<host>]\n" +
				"	earthly registry [--org <org> --project <project>] remove <host>\n",
			Action: app.actionRegistryRemove,
		},
	}
}

func (app *earthlyApp) isUserRegistryLocation() (bool, error) {
	if app.orgName == "" && app.projectName == "" {
		return true, nil
	}
	if app.orgName == "" {
		return false, fmt.Errorf("--project was specified without an --org value")
	}
	if app.projectName == "" {
		return false, fmt.Errorf("--org was specified without a --project value")
	}
	return false, nil
}

func (app *earthlyApp) getRegistriesPath() (string, error) {
	user, err := app.isUserRegistryLocation()
	if err != nil {
		return "", err
	}
	if user {
		return "/user/std/registry/", nil
	}
	return fmt.Sprintf("/%s/%s/std/registry/", app.orgName, app.projectName), nil
}

func (app *earthlyApp) actionRegistrySetup(cliCtx *cli.Context) error {
	app.commandName = "registrySetup"

	regPath, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	if cliCtx.NArg() > 1 {
		return fmt.Errorf("only a single host can be given")
	}

	host := cliCtx.Args().Get(0)
	if host == "" {
		host = "registry-1.docker.io"
	}

	if strings.Contains(host, "/") {
		return fmt.Errorf("hosts is malformed")
	}

	switch app.registryCredHelper {
	case "", "none":
		return app.actionRegistrySetupUsernamePassword(cliCtx, regPath, cloudClient, host)
	case "ecr-login":
		return app.actionRegistrySetupECRLogin(cliCtx, regPath, cloudClient, host)
	case "gcloud":
		return app.actionRegistrySetupGCloud(cliCtx, regPath, cloudClient, host)
	default:
		return fmt.Errorf("unsupported credential helper %s", app.registryCredHelper)
	}
}

func (app *earthlyApp) actionRegistrySetupECRLogin(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	if app.awsAccessKeyID == "" {
		return fmt.Errorf("--aws-access-key-id is missing (or empty)")
	}
	if app.awsSecretAccessKey == "" {
		return fmt.Errorf("--aws-secret-access-key is missing (or empty)")
	}
	err := cloudClient.SetSecret(cliCtx.Context, regPath+host+"/cred_helper", []byte("ecr-login"))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, regPath+host+"/AWS_ACCESS_KEY_ID", []byte(app.awsAccessKeyID))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, regPath+host+"/AWS_SECRET_ACCESS_KEY", []byte(app.awsSecretAccessKey))
	if err != nil {
		return err
	}
	return nil
}

func (app *earthlyApp) actionRegistrySetupGCloud(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	serviceAccountKey := app.gcpServiceAccountKey
	if app.gcpServiceAccountKeyPath != "" {
		if serviceAccountKey != "" {
			return errMultipleGCPServiceAccountFlags
		}
		data, err := os.ReadFile(app.gcpServiceAccountKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", app.gcpServiceAccountKeyPath, err)
		}
		if len(data) == 0 {
			return fmt.Errorf("service account file %s is empty", app.gcpServiceAccountKeyPath)
		}
		serviceAccountKey = string(data)
	}
	if app.gcpServiceAccountKeyStdin {
		if serviceAccountKey != "" {
			return errMultipleGCPServiceAccountFlags
		}
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
		if len(data) == 0 {
			return fmt.Errorf("no data was read from stdin")
		}
		serviceAccountKey = string(data)
	}
	if serviceAccountKey == "" {
		return fmt.Errorf("no gcp service key was provided")
	}
	err := cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "cred_helper"), []byte("gcloud"))
	if err != nil {
		return err
	}
	return cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "GCP_KEY"), []byte(serviceAccountKey))
}

func (app *earthlyApp) actionRegistrySetupUsernamePassword(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	var err error
	var password []byte
	if app.registryPasswordStdin {
		if app.registryPassword != "" {
			return fmt.Errorf("only one of  --password or --password-stdin")
		}
		password, err = io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
	} else {
		password = []byte(app.registryPassword)
	}
	if len(password) == 0 {
		return fmt.Errorf("password can not be empty")
	}

	err = cloudClient.RemoveSecret(cliCtx.Context, path.Join(regPath, host, "cred_helper"))
	if err != nil {
		if !errors.Is(err, secrets.ErrNotFound) {
			return err
		}
	}

	err = cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "username"), []byte(app.registryUsername))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "password"), password)
	if err != nil {
		return err
	}

	return nil
}

type registryCredentials struct {
	host           string
	credHelper     string
	username       string
	accessID       string
	gcpClientEmail string
}

func extractGCPClientEmail(gcpKey string) (string, error) {
	data := struct {
		ClientEmail string `json:"client_email"`
	}{}
	err := json.Unmarshal([]byte(gcpKey), &data)
	if err != nil {
		return "", err
	}
	return data.ClientEmail, nil
}

func (app *earthlyApp) secretsToRegistryCredentials(pathPrefix string, secrets []*cloud.Secret) []*registryCredentials {
	credentials := map[string]*registryCredentials{}
	for _, secret := range secrets {
		parts := strings.Split(strings.TrimPrefix(secret.Path, pathPrefix), "/")
		if len(parts) != 2 {
			continue
		}
		host := parts[0]
		key := parts[1]

		rc, ok := credentials[host]
		if !ok {
			rc = &registryCredentials{
				host: host,
			}
			credentials[host] = rc
		}

		switch key {
		case "cred_helper":
			rc.credHelper = secret.Value
		case "username":
			rc.username = secret.Value
		case "AWS_ACCESS_KEY_ID":
			rc.accessID = secret.Value
		case "GCP_KEY":
			var err error
			rc.gcpClientEmail, err = extractGCPClientEmail(secret.Value)
			if err != nil {
				app.console.Warnf("Failed to extract client email from %s: %s\n", path.Join(pathPrefix, "GCP_KEY"), err)
			}
		}
	}

	sortedCredentials := []*registryCredentials{}
	for _, rc := range credentials {
		sortedCredentials = append(sortedCredentials, rc)
	}
	sort.Slice(sortedCredentials, func(i, j int) bool {
		return sortedCredentials[i].host < sortedCredentials[j].host
	})

	return sortedCredentials
}

func (app *earthlyApp) actionRegistryList(cliCtx *cli.Context) error {
	app.commandName = "registryList"

	path, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Registry", "Cred Helper", "Username/Access ID")
	for _, rc := range app.secretsToRegistryCredentials(path, secrets) {
		switch rc.credHelper {
		case "ecr-login":
			fmt.Fprintf(w, "%s\t%s\t%s\n", rc.host, rc.credHelper, rc.accessID)
		case "gcloud":
			fmt.Fprintf(w, "%s\t%s\t%s\n", rc.host, rc.credHelper, rc.gcpClientEmail)
		case "":
			fmt.Fprintf(w, "%s\t%s\t%s\n", rc.host, "none", rc.username)
		default:
			fmt.Fprintf(w, "%s\t%s\t%s\n", rc.host, "unknown: "+rc.credHelper, "")
		}
	}
	w.Flush()

	return nil
}

func (app *earthlyApp) actionRegistryRemove(cliCtx *cli.Context) error {
	app.commandName = "registryRemove"
	path, err := app.getRegistriesPath()
	if err != nil {
		return err
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	host := cliCtx.Args().Get(0)
	if host == "" {
		host = "registry-1.docker.io"
	}

	if cliCtx.NArg() > 1 {
		return fmt.Errorf("only a single registry host can be given (or perhaps the host was before the --flags?)")
	}

	fmt.Printf("Removing registry credentials for %s\n", host)
	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path+host)
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		err = cloudClient.RemoveSecret(cliCtx.Context, secret.Path)
		if err != nil && !errors.Is(err, cloud.ErrNotFound) {
			return err
		}
	}

	return nil
}
