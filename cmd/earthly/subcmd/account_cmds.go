package subcmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/helper"
)

var (
	errLoginFlagsHaveNoEffect            = errors.New("account login flags have no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
	errLogoutHasNoEffectWhenAuthTokenSet = errors.New("account logout has no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
)

type Account struct {
	cli CLI

	email                  string
	token                  string
	password               string
	termsConditionsPrivacy bool
	registrationPublicKey  string
	writePermission        bool
	expiry                 string
	overWrite              bool
}

func NewAccount(cli CLI) *Account {
	return &Account{
		cli: cli,
	}
}

func (a *Account) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "account",
			Usage:       "Create or manage an Earthly account",
			Description: "Create or manage an Earthly account.",
			Subcommands: []*cli.Command{
				{
					Name:        "register",
					Usage:       "Register for an Earthly account",
					Description: "Register for an Earthly account.",
					UsageText: "You may register using GitHub OAuth, by visiting https://ci.earthly.dev\n" +
						"   Once authenticated, a login token will be displayed which can be used to login:\n" +
						"\n" +
						"       earthly [options] account login --token <token>\n" +
						"\n" +
						"   Alternatively, you can register using an email:\n" +
						"       first, request a token with:\n" +
						"\n" +
						"           earthly [options] account register --email <email>\n" +
						"\n" +
						"       then check your email to retrieve the token, then continue by running:\n" +
						"\n" +
						"           earthly [options] account register --token <token>\n",
					Action: a.actionRegister,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "email",
							Usage:       "Email address to use for register for your Earthly account",
							Destination: &a.email,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "Email verification token, retreived from the email you used to register your account",
							Destination: &a.token,
						},
						&cli.StringFlag{
							Name:        "password",
							EnvVars:     []string{"EARTHLY_PASSWORD"},
							Usage:       "Specify password on the command line instead of interactively being asked",
							Destination: &a.password,
						},
						&cli.StringFlag{
							Name:        "public-key",
							EnvVars:     []string{"EARTHLY_PUBLIC_KEY"},
							Usage:       "Path to public key to register",
							Destination: &a.registrationPublicKey,
						},
						&cli.BoolFlag{
							Name:        "accept-terms-of-service-privacy",
							EnvVars:     []string{"EARTHLY_ACCEPT_TERMS_OF_SERVICE_PRIVACY"},
							Usage:       "Accept the Terms & Conditions, and Privacy Policy",
							Destination: &a.termsConditionsPrivacy,
						},
					},
				},
				{
					Name:        "login",
					Usage:       "Login to an Earthly account",
					Description: "Login to an Earthly account.",
					UsageText: "earthly [options] account login\n" +
						"   earthly [options] account login --email <email>\n" +
						"   earthly [options] account login --email <email> --password <password>\n" +
						"   earthly [options] account login --token <token>\n",
					Action: a.actionLogin,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "email",
							Usage:       "Pass in email address connected with your Earthly account",
							Destination: &a.email,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "Authentication token",
							Destination: &a.token,
						},
						&cli.StringFlag{
							Name:        "password",
							EnvVars:     []string{"EARTHLY_PASSWORD"},
							Usage:       "Specify password on the command line instead of interactively being asked",
							Destination: &a.password,
						},
					},
				},
				{
					Name:        "logout",
					Usage:       "Logout of an Earthly account",
					Description: "Logout of an Earthly account; this has no effect for ssh-based authentication.",
					Action:      a.actionLogout,
				},
				{
					Name:        "list-keys",
					Usage:       "List associated public keys used for authentication",
					UsageText:   "earthly [options] account list-keys",
					Description: "Lists all public keys that are authorized to login to the current Earthly account.",
					Action:      a.actionListKeys,
				},
				{
					Name:        "add-key",
					Usage:       "Authorize a new public key to login with with the current Earthly account",
					UsageText:   "earthly [options] add-key [<key>]",
					Description: "Authorize a new public key to login to the current Earthly account. If 'key' is omitted, an interactive prompt is displayed to select a public key to add.",
					Action:      a.actionAddKey,
				},
				{
					Name:        "remove-key",
					Usage:       "Removes an authorized public key from the current Earthly account",
					UsageText:   "earthly [options] remove-key <key>",
					Description: "Removes an authorized public key from accessing the current Earthly account.",
					Action:      a.actionRemoveKey,
				},
				{
					Name:        "list-tokens",
					Usage:       "List associated tokens used for authentication with the current Earthly account",
					UsageText:   "earthly [options] account list-tokens",
					Description: "List account tokens associated with the current Earthly account. A token is useful for environments where the ssh-agent is not accessible (e.g. a CI system).",
					Action:      a.actionListTokens,
				},
				{
					Name:      "create-token",
					Usage:     "Create a new authentication token for your account",
					UsageText: "earthly [options] account create-token [options] <token name>",
					Description: `Creates a new authentication token. A read-only token is created by default, If the '--write' flag is specified the token will have read+write access.
				The token will never expire, unless an expiry date is supplied via the '--expiry' option.`,
					Action: a.actionCreateToken,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "write",
							Usage:       "Grant write permissions in addition to read",
							Destination: &a.writePermission,
						},
						&cli.StringFlag{
							Name:        "expiry",
							Usage:       "Set token expiry date in the form YYYY-MM-DD or never (default never)",
							Destination: &a.expiry,
						},
						&cli.BoolFlag{
							Name:        "overwrite",
							Usage:       "Overwrite the token if it already exists",
							Destination: &a.overWrite,
						},
					},
				},
				{
					Name:        "remove-token",
					Usage:       "Remove an authentication token from your account",
					UsageText:   "earthly [options] account remove-token <token>",
					Description: "Removes a token from the current Earthly account.",
					Action:      a.actionRemoveToken,
				},
				{
					Name:  "reset",
					Usage: "Reset Earthly account password",
					UsageText: `earthly [options] account reset --email <email>
			earthly [options] account reset --email <email> --token <token>`,
					Description: `Reset the password associated with the provided email.
			The command should first be run without a token, which will cause a token to be emailed to you.
			Once the command is re-run with the provided token, it will prompt you for a new password.`,
					Action: a.actionReset,
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "email",
							Usage:       "Email address for which to reset the password",
							Destination: &a.email,
						},
						&cli.StringFlag{
							Name:        "token",
							Usage:       "Authentication token with with to rerun the command with your email to reset your password",
							Destination: &a.token,
						},
					},
				},
			},
		},
	}
}

func (a *Account) actionRegister(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountRegister")
	if a.email == "" {
		return errors.New("no email given")
	}

	if !strings.Contains(a.email, "@") {
		return errors.New("email is invalid")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	if a.token == "" {
		err := cloudClient.RegisterEmail(cliCtx.Context, a.email)
		if err != nil {
			return errors.Wrap(err, "failed to register email")
		}
		fmt.Printf("An email has been sent to %q containing a registration token\n", a.email)
		return nil
	}

	var publicKeys []*agent.Key
	if a.cli.Flags().SSHAuthSock != "" {
		var err error
		publicKeys, err = cloudClient.GetPublicKeys(cliCtx.Context)
		if err != nil {
			return err
		}
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass a.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	pword := a.password
	if pword == "" {
		pword, err = promptNewPassword()
		if err != nil {
			return err
		}
	}

	var interactiveAccept bool
	if !a.termsConditionsPrivacy {
		rawAccept, err := common.PromptInput(cliCtx.Context, "I acknowledge Earthly Technologies’ Privacy Policy (https://earthly.dev/privacy-policy) and agree to Earthly Technologies Terms of Service (https://earthly.dev/tos) [y/N]: ")
		if err != nil {
			return err
		}
		if rawAccept == "" {
			rawAccept = "n"
		}
		accept := strings.ToLower(rawAccept)[0]
		interactiveAccept = (accept == 'y')
	}
	termsConditionsPrivacy := a.termsConditionsPrivacy || interactiveAccept

	var publicKey string
	if a.registrationPublicKey == "" {
		if len(publicKeys) > 0 {
			rawIsRegisterSSHKey, err := common.PromptInput(cliCtx.Context, "Would you like to enable password-less login using public key authentication (https://earthly.dev/public-key-auth)? [Y/n]: ")
			if err != nil {
				return err
			}
			if rawIsRegisterSSHKey == "" {
				rawIsRegisterSSHKey = "y"
			}
			isRegisterSSHKey := (strings.ToLower(rawIsRegisterSSHKey)[0] == 'y')
			if isRegisterSSHKey {
				fmt.Printf("Which of the following keys do you want to register?\n")
				fmt.Printf("0) none\n")
				for i, key := range publicKeys {
					fmt.Printf("%d) %s\n", i+1, key.String())
				}
				keyNum, err := common.PromptInput(cliCtx.Context, "enter key number (1=default): ")
				if err != nil {
					return err
				}
				if keyNum == "" {
					keyNum = "1"
				}
				i, err := strconv.Atoi(keyNum)
				if err != nil {
					return errors.Wrap(err, "invalid key number")
				}
				if i < 0 || i > len(publicKeys) {
					return errors.Errorf("invalid key number")
				}
				if i > 0 {
					publicKey = publicKeys[i-1].String()
				}
			}
		}
	} else {
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(a.registrationPublicKey))
		if err == nil {
			// supplied public key is valid
			publicKey = a.registrationPublicKey
		} else {
			// otherwise see if it matches the name (Comment) of a key known by the ssh agent
			for _, key := range publicKeys {
				if key.Comment == a.registrationPublicKey {
					publicKey = key.String()
					break
				}
			}
			if publicKey == "" {
				return errors.Errorf("failed to find key in ssh agent's known keys")
			}
		}
	}

	err = cloudClient.CreateAccount(cliCtx.Context, a.email, a.token, pword, publicKey, termsConditionsPrivacy)
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	fmt.Println("Account registration complete")
	return nil
}

// promptHiddenText prompts the user to enter a value without echoing it to stdout
// the hiddenTextPrompt is displayed to the user followed by a colon, as an interactive prompt.
func promptHiddenText(hiddenTextPrompt string) (string, error) {
	fmt.Printf("%s: ", hiddenTextPrompt)
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println("")
	return string(password), nil
}
func promptPassword() (string, error) {
	return promptHiddenText("password")
}

func promptNewPassword() (string, error) {
	enteredPassword, err := promptHiddenText("New password")
	if err != nil {
		return "", err
	}
	enteredPassword2, err := promptHiddenText("Confirm password")
	if err != nil {
		return "", err
	}
	fmt.Println("")
	if string(enteredPassword) != string(enteredPassword2) {
		return "", errors.Errorf("passwords do not match")
	}
	return string(enteredPassword), nil
}

func (a *Account) actionListKeys(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountListKeys")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	keys, err := cloudClient.ListPublicKeys(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to list account keys")
	}
	for _, key := range keys {
		fmt.Printf("%s\n", key)
	}
	return nil
}

func (a *Account) actionAddKey(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountAddKey")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	if cliCtx.NArg() > 0 {
		for _, k := range cliCtx.Args().Slice() {
			err := cloudClient.AddPublicKey(cliCtx.Context, k)
			if err != nil {
				return errors.Wrap(err, "failed to add public key to account")
			}
		}
		return nil
	}

	publicKeys, err := cloudClient.GetPublicKeys(cliCtx.Context)
	if err != nil {
		return err
	}
	if len(publicKeys) == 0 {
		return errors.Errorf("unable to list available public keys, is ssh-agent running?; do you need to run ssh-add?")
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass a.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Which of the following public keys do you want to register (your private key is never sent or even read by earthly)?\n")
	for i, key := range publicKeys {
		fmt.Printf("%d) %s\n", i+1, key.String())
	}
	keyNum, err := common.PromptInput(cliCtx.Context, "enter key number (1=default): ")
	if err != nil {
		return err
	}
	if keyNum == "" {
		keyNum = "1"
	}
	i, err := strconv.Atoi(keyNum)
	if err != nil {
		return errors.Wrap(err, "invalid key number")
	}
	if i <= 0 || i > len(publicKeys) {
		return errors.Errorf("invalid key number")
	}
	publicKey := publicKeys[i-1].String()

	err = cloudClient.AddPublicKey(cliCtx.Context, publicKey)
	if err != nil {
		return errors.Wrap(err, "failed to add public key to account")
	}

	// switch over to new key if the user is currently using password-based auth
	email, authType, _, err := cloudClient.WhoAmI(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to validate auth token")
	}
	if authType == "password" {
		err = cloudClient.SetSSHCredentials(cliCtx.Context, email, publicKey)
		if err != nil {
			a.cli.Console().Warnf("failed to authenticate using newly added public key: %s", err.Error())
			return nil
		}
		fmt.Printf("Switching from password-based login to ssh-based login\n")
	}

	return nil
}

func (a *Account) actionRemoveKey(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountRemoveKey")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	for _, k := range cliCtx.Args().Slice() {
		err := cloudClient.RemovePublicKey(cliCtx.Context, k)
		if err != nil {
			return errors.Wrap(err, "failed to add public key to account")
		}
	}
	return nil
}
func (a *Account) actionListTokens(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountListTokens")
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	tokens, err := cloudClient.ListTokens(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to list account tokens")
	}
	if len(tokens) == 0 {
		return nil // avoid printing header columns when there are no tokens
	}

	now := time.Now()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Token Name\tRead/Write\tExpiry\tLast Used\n")
	for _, token := range tokens {
		expired := now.After(token.Expiry)
		fmt.Fprintf(w, "%s", token.Name)
		if token.Write {
			fmt.Fprintf(w, "\trw")
		} else {
			fmt.Fprintf(w, "\tr")
		}
		if token.Indefinite {
			fmt.Fprintf(w, "\tNever")
		} else {
			fmt.Fprintf(w, "\t%s UTC", token.Expiry.UTC().Format("2006-01-02T15:04"))
			if expired {
				fmt.Fprintf(w, " *expired*")
			}
		}
		if token.LastAccessedAt.UTC().IsZero() {
			fmt.Fprint(w, "\tNever")
		} else {
			fmt.Fprintf(w, "\t%s UTC", token.LastAccessedAt.UTC().Format("2006-01-02T15:04"))
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}
func (a *Account) actionCreateToken(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountCreateToken")
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	var expiry *time.Time
	if a.expiry != "" && a.expiry != "never" {
		layouts := []string{
			"2006-01-02",
			time.RFC3339,
		}

		var err error
		for _, layout := range layouts {
			var parsedTime time.Time
			parsedTime, err = time.Parse(layout, a.expiry)
			if err == nil {
				expiry = &parsedTime
				break
			}
		}
		if err != nil {
			return errors.Errorf("failed to parse expiry %q", a.expiry)
		}
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	name := cliCtx.Args().First()
	token, err := cloudClient.CreateToken(cliCtx.Context, name, a.writePermission, expiry, a.overWrite)
	if err != nil {
		return errors.Wrap(err, "failed to create token")
	}

	expiryStr := "will never expire"
	if expiry != nil {
		expiryStr = fmt.Sprintf("will expire %s", humanize.Time(*expiry))
	}

	fmt.Printf("created token %q which %s; save this token somewhere, it can't be viewed again (only reset)\n", token, expiryStr)

	return nil
}
func (a *Account) actionRemoveToken(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountRemoveToken")
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	name := cliCtx.Args().First()
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.RemoveToken(cliCtx.Context, name)
	if err != nil {
		return errors.Wrap(err, "failed to remove account tokens")
	}
	return nil
}

func (a *Account) actionLogin(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountLogin")
	email := a.email
	token := a.token
	pass := a.password

	opts := []cloud.ClientOpt{}
	if token != "" {
		opts = append(opts, cloud.WithAuthToken(token))
	}

	cloudClient, err := helper.NewCloudClient(a.cli, opts...)
	if err != nil {
		return err
	}

	loggedInEmail, authType, _, err := cloudClient.WhoAmI(cliCtx.Context)

	if err == nil && authType == cloud.AuthMethodCachedJWT {
		// already logged in, don't re-attempt a login
		a.cli.Console().Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		return nil
	}

	if err != nil && !errors.Is(err, cloud.ErrUnauthorized) && !errors.Is(err, cloud.ErrAuthTokenExpired) {
		return errors.Wrap(err, "failed to login")
	}

	// if a user explicitly calls Login, we will allow future auto-logins
	// which is required for refreshing tokens
	err = cloudClient.EnableAutoLogin(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to enable auto login")
	}

	if cliCtx.NArg() == 1 {
		emailOrToken := cliCtx.Args().First()
		if token == "" && email == "" {
			if cloud.IsValidEmail(cliCtx.Context, emailOrToken) {
				email = emailOrToken
			} else {
				token = emailOrToken
			}

		} else {
			return errors.New("invalid number of arguments provided")
		}
	} else if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}

	if token != "" && (email != "" || pass != "") {
		return errors.New("--token cannot be used in conjuction with --email or --password")
	}

	// special case where global auth token overrides login logic
	if a.cli.Flags().AuthToken != "" {
		if email != "" || token != "" || pass != "" {
			return errLoginFlagsHaveNoEffect
		}
		loggedInEmail, authType, writeAccess, err := cloudClient.WhoAmI(cliCtx.Context)
		if err != nil {
			return errors.Wrap(err, "failed to validate auth token")
		}
		if !writeAccess {
			authType = "read-only-" + authType
		}
		a.cli.Console().Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		a.printLogSharingMessage()
		return nil
	}

	err = cloudClient.DeleteCachedToken(cliCtx.Context)
	if err != nil {
		return err
	}

	if token != "" || pass != "" {
		err := cloudClient.DeleteAuthCache(cliCtx.Context)
		if err != nil {
			return errors.Wrap(err, "failed to clear cached credentials")
		}
		cloudClient.DisableSSHKeyGuessing(cliCtx.Context)
	} else if email != "" {
		if err = cloudClient.FindSSHCredentials(cliCtx.Context, email); err == nil {
			// if err is not nil, we will try again below via cloudClient.WhoAmI()
			authType, err = cloudClient.Authenticate(cliCtx.Context)
			if err != nil {
				return errors.Wrap(err, "authentication with cloud server failed")
			}
			a.cli.Console().Printf("Logged in as %q using %s auth\n", email, authType)
			a.printLogSharingMessage()
			return nil
		}
	}

	loggedInEmail, authType, writeAccess, err := cloudClient.WhoAmI(cliCtx.Context)
	switch errors.Cause(err) {
	case cloud.ErrUnauthorized:
		break
	case nil:
		if email != "" && email != loggedInEmail {
			break // case where a user has multiple emails and wants to switch
		}
		if !writeAccess {
			authType = "read-only-" + authType
		}
		a.cli.Console().Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		a.printLogSharingMessage()
		return nil
	default:
		return err
	}

	if email == "" && token == "" {
		if a.cli.Flags().SSHAuthSock == "" {
			a.cli.Console().Warnf("No ssh auth socket detected; falling back to password-based login\n")
		}

		emailOrToken, err := common.PromptInput(cliCtx.Context, "enter your email or auth token: ")
		if err != nil {
			return err
		}
		if strings.Contains(emailOrToken, "@") {
			email = emailOrToken
		} else {
			token = emailOrToken
		}
	}

	if email != "" && pass == "" {
		// Our signal handling under main() doesn't cause reading from stdin to cancel
		// as there's no way to pass a.ctx to stdin read calls.
		signal.Reset(syscall.SIGINT, syscall.SIGTERM)

		fmt.Printf("enter your password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		fmt.Println("")
		pass = string(passwordBytes)
		if pass == "" {
			return errors.Errorf("no password entered")
		}
	}

	if token != "" {
		email, err = cloudClient.SetTokenCredentials(cliCtx.Context, token)
		if err != nil {
			return err
		}
	} else {
		err = a.loginAndSavePasswordCredentials(cliCtx.Context, cloudClient, email, string(pass))
		if err != nil {
			return err
		}
	}
	authType, err = cloudClient.Authenticate(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "authentication with cloud server failed")
	}
	a.cli.Console().Printf("Logged in as %q using %s auth\n", email, authType)
	a.printLogSharingMessage()
	return nil
}

func (a *Account) printLogSharingMessage() {
	a.cli.Console().Printf("Log sharing is enabled by default. If you would like to disable it, run:\n" +
		"\n" +
		"\tearthly config global.disable_log_sharing true")
}

func (a *Account) loginAndSavePasswordCredentials(ctx context.Context, cloudClient *cloud.Client, email, password string) error {
	err := cloudClient.SetPasswordCredentials(ctx, email, password)
	if err != nil {
		return err
	}
	a.cli.Console().Printf("Logged in as %q using password auth\n", email)
	a.cli.Console().Printf("Warning unencrypted password has been stored under ~/.earthly/auth.credentials; consider using ssh-based auth to prevent this.\n")
	return nil
}

func (a *Account) actionLogout(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountLogout")

	if a.cli.Flags().AuthToken != "" {
		return errLogoutHasNoEffectWhenAuthTokenSet
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	err = cloudClient.DeleteAuthCache(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to logout")
	}
	err = cloudClient.DisableAutoLogin(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "failed to disable auto login")
	}
	return nil
}

func (a *Account) actionReset(cliCtx *cli.Context) error {
	a.cli.SetCommandName("accountReset")

	if a.email == "" {
		return errors.New("no email given")
	}

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}

	if a.token == "" {
		err = cloudClient.AccountResetRequestToken(cliCtx.Context, a.email)
		if err != nil {
			return errors.Wrap(err, "failed to request account reset token")
		}
		a.cli.Console().Printf("An account reset token has been emailed to %q\n", a.email)
		return nil
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass a.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	pword, err := promptNewPassword()
	if err != nil {
		return err
	}

	err = cloudClient.AccountReset(cliCtx.Context, a.email, a.token, pword)
	if err != nil {
		return errors.Wrap(err, "failed to reset account")
	}
	a.cli.Console().Printf("Account password has been reset\n")
	err = a.loginAndSavePasswordCredentials(cliCtx.Context, cloudClient, a.email, pword)
	if err != nil {
		return err
	}
	a.printLogSharingMessage()
	return nil
}
