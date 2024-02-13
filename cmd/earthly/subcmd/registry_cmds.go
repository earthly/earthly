package subcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/helper"
)

const (
	gcpServiceAccountKeyFlag      = "gcp-service-account-key"
	gcpServiceAccountKeyPathFlag  = "gcp-service-account-key-path"
	gcpServiceAccountKeyStdinFlag = "gcp-service-account-key-stdin"
)

var (
	errMultipleGCPServiceAccountFlags = fmt.Errorf("the --%s --%s --%s flags are mutually exclusive", gcpServiceAccountKeyFlag, gcpServiceAccountKeyPathFlag, gcpServiceAccountKeyStdinFlag)
)

type Registry struct {
	cli CLI

	CredHelper                string
	Username                  string
	Password                  string
	PasswordStdin             bool
	awsAccessKeyID            string
	awsSecretAccessKey        string
	gcpServiceAccountKey      string
	gcpServiceAccountKeyPath  string
	gcpServiceAccountKeyStdin bool
}

func NewRegistry(cli CLI) *Registry {
	return &Registry{
		cli: cli,
	}
}

func (a *Registry) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "registry",
			Aliases:     []string{"registries"},
			Description: "*beta* Manage registry access.",
			Usage:       "*beta* Manage registry access",
			UsageText:   "earthly [options] registry [--org <organization-name>, --project <project>] (setup|list|remove) [<flags>]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "org",
					EnvVars:     []string{"EARTHLY_ORG"},
					Usage:       "The organization to which the project belongs",
					Required:    false,
					Destination: &a.cli.Flags().OrgName,
				},
				&cli.StringFlag{
					Name:        "project",
					EnvVars:     []string{"EARTHLY_PROJECT"},
					Usage:       "The organization project in which to store registry credentials",
					Required:    false,
					Destination: &a.cli.Flags().ProjectName,
				},
			},
			Subcommands: []*cli.Command{
				{
					Name:        "setup",
					Usage:       "*beta* Setup, and store, registry credentials in the earthly-cloud",
					Description: "*beta* Setup, and store, registry credentials in the earthly-cloud.",
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
					Action: a.actionSetup,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "cred-helper",
							EnvVars:     []string{"EARTHLY_REGISTRY_CRED_HELPER"},
							Usage:       "Use a credential helper when logging into the registry (ecr-login, gcloud)",
							Required:    false,
							Destination: &a.CredHelper,
						},
						&cli.StringFlag{
							Name:        "username",
							EnvVars:     []string{"EARTHLY_REGISTRY_USERNAME"},
							Usage:       "The username to use when logging into the registry; if omitted, earthly will prompt for a username via stdin.",
							Required:    false,
							Destination: &a.Username,
						},
						&cli.StringFlag{
							Name:    "password",
							EnvVars: []string{"EARTHLY_REGISTRY_PASSWORD"},
							Usage: `The password to use when logging into the registry
						(use --password-stdin to prevent leaking your password via your shell history)`,
							Required:    false,
							Destination: &a.Password,
						},
						&cli.BoolFlag{
							Name:        "password-stdin",
							EnvVars:     []string{"EARTHLY_REGISTRY_PASSWORD_STDIN"},
							Usage:       "(Deprecated) Read the password from stdin (and wait for an EOF)",
							Required:    false,
							Destination: &a.PasswordStdin,
						},
						&cli.StringFlag{
							Name:        "aws-access-key-id",
							EnvVars:     []string{"AWS_ACCESS_KEY_ID"},
							Usage:       "AWS access key ID to use for elastic-container-registry",
							Required:    false,
							Destination: &a.awsAccessKeyID,
						},
						&cli.StringFlag{
							Name:        "aws-secret-access-key",
							EnvVars:     []string{"AWS_SECRET_ACCESS_KEY"},
							Usage:       "AWS secret access key to use for elastic-container-registry",
							Required:    false,
							Destination: &a.awsSecretAccessKey,
						},
						&cli.StringFlag{
							Name:        gcpServiceAccountKeyFlag,
							EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY"},
							Usage:       "GCP service account key to use for artifact or container registry",
							Required:    false,
							Destination: &a.gcpServiceAccountKey,
						},
						&cli.StringFlag{
							Name:        gcpServiceAccountKeyPathFlag,
							EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY_PATH", "GOOGLE_APPLICATION_CREDENTIALS"},
							Usage:       "path to GCP service account key to use for artifact or container registry",
							Required:    false,
							Destination: &a.gcpServiceAccountKeyPath,
						},
						&cli.BoolFlag{
							Name:        gcpServiceAccountKeyStdinFlag,
							EnvVars:     []string{"GCP_SERVICE_ACCOUNT_KEY_STDIN"},
							Usage:       "GCP service account key to use for artifact or container registry, read from stdin",
							Required:    false,
							Destination: &a.gcpServiceAccountKeyStdin,
						},
					},
				},
				{
					Name:  "list",
					Usage: "*beta* List configured registries",
					UsageText: "earthly registry list\n" +
						"	earthly registry --org <org> --project <project> list\n",
					Description: "*beta* List configured registries.",
					Action:      a.actionList,
				},
				{
					Name:  "remove",
					Usage: "*beta* Remove stored registry credentials",
					UsageText: "earthly registry remove [<host>]\n" +
						"	earthly registry [--org <org> --project <project>] remove <host>\n",
					Description: "*beta* Remove stored registry credentials.",
					Action:      a.actionRemove,
				},
			},
		},
	}
}

func (a *Registry) isUserRegistryLocation() (bool, error) {
	if a.cli.OrgName() == "" && a.cli.Flags().ProjectName == "" {
		return true, nil
	}
	if a.cli.OrgName() == "" {
		return false, fmt.Errorf("--project was specified without an --org value")
	}
	if a.cli.Flags().ProjectName == "" {
		return false, fmt.Errorf("--org was specified without a --project value")
	}
	return false, nil
}

func (a *Registry) getRegistriesPath(ctx context.Context, cloudClient *cloud.Client) (string, error) {
	user, err := a.isUserRegistryLocation()
	if err != nil {
		return "", err
	}
	if user {
		return "/user/std/registry/", nil
	}
	orgName, projectName, _, err := getOrgAndProject(ctx, a.cli.OrgName(), a.cli.Flags().ProjectName, cloudClient, "")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/%s/%s/std/registry/", orgName, projectName), nil
}

func (a *Registry) actionSetup(cliCtx *cli.Context) error {
	a.cli.SetCommandName("registrySetup")

	if cliCtx.NArg() > 1 {
		return fmt.Errorf("only a single host can be given")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	regPath, err := a.getRegistriesPath(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	host := cliCtx.Args().Get(0)
	if host == "" {
		host = "registry-1.docker.io"
	}

	if strings.Contains(host, "/") {
		return fmt.Errorf("hosts is malformed")
	}

	switch a.CredHelper {
	case "", "none":
		return a.actionSetupUsernamePassword(cliCtx, regPath, cloudClient, host)
	case "ecr-login":
		return a.actionSetupECRLogin(cliCtx, regPath, cloudClient, host)
	case "gcloud":
		return a.actionSetupGCloud(cliCtx, regPath, cloudClient, host)
	default:
		return fmt.Errorf("unsupported credential helper %s", a.CredHelper)
	}
}

func (a *Registry) actionSetupECRLogin(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	if a.awsAccessKeyID == "" {
		return fmt.Errorf("--aws-access-key-id is missing (or empty)")
	}
	if a.awsSecretAccessKey == "" {
		return fmt.Errorf("--aws-secret-access-key is missing (or empty)")
	}
	err := cloudClient.SetSecret(cliCtx.Context, regPath+host+"/cred_helper", []byte("ecr-login"))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, regPath+host+"/AWS_ACCESS_KEY_ID", []byte(a.awsAccessKeyID))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, regPath+host+"/AWS_SECRET_ACCESS_KEY", []byte(a.awsSecretAccessKey))
	if err != nil {
		return err
	}
	return nil
}

func (a *Registry) actionSetupGCloud(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	serviceAccountKey := a.gcpServiceAccountKey
	if a.gcpServiceAccountKeyPath != "" {
		if serviceAccountKey != "" {
			return errMultipleGCPServiceAccountFlags
		}
		data, err := os.ReadFile(a.gcpServiceAccountKeyPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", a.gcpServiceAccountKeyPath, err)
		}
		if len(data) == 0 {
			return fmt.Errorf("service account file %s is empty", a.gcpServiceAccountKeyPath)
		}
		serviceAccountKey = string(data)
	}
	if a.gcpServiceAccountKeyStdin {
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

func (a *Registry) actionSetupUsernamePassword(cliCtx *cli.Context, regPath string, cloudClient *cloud.Client, host string) error {
	var err error
	username := a.Username
	if username == "" {
		fmt.Printf("username: ")
		fmt.Scanln(&username)
		if username == "" {
			return fmt.Errorf("username can not be empty")
		}
	}

	password := a.Password
	if a.PasswordStdin {
		a.cli.Console().Warnf("Deprecated: the --password-stdin flag will be removed in the future, to read from stdin simply omit both of the --password-stdin and --password flags\n")
		passwordBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read from stdin")
		}
		password = string(passwordBytes)
	}
	if password == "" {
		// prompt via stdin instead

		// Our signal handling under main() doesn't cause reading from stdin to cancel
		// as there's no way to pass app.ctx to stdin read calls.
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)

		password, err = promptPassword()
		if err != nil {
			return err
		}
	}
	if len(password) == 0 {
		return fmt.Errorf("password can not be empty")
	}

	err = cloudClient.RemoveSecret(cliCtx.Context, path.Join(regPath, host, "cred_helper"))
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return err
		}
	}

	err = cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "username"), []byte(username))
	if err != nil {
		return err
	}
	err = cloudClient.SetSecret(cliCtx.Context, path.Join(regPath, host, "password"), []byte(password))
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

func (a *Registry) secretsToRegistryCredentials(pathPrefix string, secrets []*cloud.Secret) []*registryCredentials {
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
				a.cli.Console().Warnf("Failed to extract client email from %s: %s\n", path.Join(pathPrefix, "GCP_KEY"), err)
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

func (a *Registry) actionList(cliCtx *cli.Context) error {
	a.cli.SetCommandName("registryList")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err := a.getRegistriesPath(cliCtx.Context, cloudClient)
	if err != nil {
		return err
	}

	secrets, err := cloudClient.ListSecrets(cliCtx.Context, path)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Registry", "Cred Helper", "Username/Access ID")
	for _, rc := range a.secretsToRegistryCredentials(path, secrets) {
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

func (a *Registry) actionRemove(cliCtx *cli.Context) error {
	a.cli.SetCommandName("registryRemove")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	path, err := a.getRegistriesPath(cliCtx.Context, cloudClient)
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
