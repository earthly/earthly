package main

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
)

var (
	errLoginFlagsHaveNoEffect            = errors.New("account login flags have no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
	errLogoutHasNoEffectWhenAuthTokenSet = errors.New("account logout has no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")
)

func (app *earthlyApp) accountCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "register",
			Usage:       "Register for an Earthly account *beta*",
			Description: "Register for an Earthly account *beta*",
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
			Action: app.actionAccountRegister,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "email",
					Usage:       "Email address",
					Destination: &app.email,
				},
				&cli.StringFlag{
					Name:        "token",
					Usage:       "Email verification token",
					Destination: &app.token,
				},
				&cli.StringFlag{
					Name:        "password",
					EnvVars:     []string{"EARTHLY_PASSWORD"},
					Usage:       "Specify password on the command line instead of interactively being asked",
					Destination: &app.password,
				},
				&cli.StringFlag{
					Name:        "public-key",
					EnvVars:     []string{"EARTHLY_PUBLIC_KEY"},
					Usage:       "Path to public key to register",
					Destination: &app.registrationPublicKey,
				},
				&cli.BoolFlag{
					Name:        "accept-terms-of-service-privacy",
					EnvVars:     []string{"EARTHLY_ACCEPT_TERMS_OF_SERVICE_PRIVACY"},
					Usage:       "Accept the Terms & Conditions, and Privacy Policy",
					Destination: &app.termsConditionsPrivacy,
				},
			},
		},
		{
			Name:        "login",
			Usage:       "Login to an Earthly account *beta*",
			Description: "Login to an Earthly account *beta*",
			UsageText: "earthly [options] account login\n" +
				"   earthly [options] account login --email <email>\n" +
				"   earthly [options] account login --email <email> --password <password>\n" +
				"   earthly [options] account login --token <token>\n",
			Action: app.actionAccountLogin,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "email",
					Usage:       "Email address",
					Destination: &app.email,
				},
				&cli.StringFlag{
					Name:        "token",
					Usage:       "Authentication token",
					Destination: &app.token,
				},
				&cli.StringFlag{
					Name:        "password",
					EnvVars:     []string{"EARTHLY_PASSWORD"},
					Usage:       "Specify password on the command line instead of interactively being asked",
					Destination: &app.password,
				},
			},
		},
		{
			Name:        "logout",
			Usage:       "Logout of an Earthly account *beta*",
			Description: "Logout of an Earthly account; this has no effect for ssh-based authentication *beta*",
			Action:      app.actionAccountLogout,
		},
		{
			Name:      "list-keys",
			Usage:     "List associated public keys used for authentication *beta*",
			UsageText: "earthly [options] account list-keys",
			Action:    app.actionAccountListKeys,
		},
		{
			Name:      "add-key",
			Usage:     "Associate a new public key with your account *beta*",
			UsageText: "earthly [options] add-key [<key>]",
			Action:    app.actionAccountAddKey,
		},
		{
			Name:      "remove-key",
			Usage:     "Removes an existing public key from your account *beta*",
			UsageText: "earthly [options] remove-key <key>",
			Action:    app.actionAccountRemoveKey,
		},
		{
			Name:      "list-tokens",
			Usage:     "List associated tokens used for authentication *beta*",
			UsageText: "earthly [options] account list-tokens",
			Action:    app.actionAccountListTokens,
		},
		{
			Name:      "create-token",
			Usage:     "Create a new authentication token for your account *beta*",
			UsageText: "earthly [options] account create-token [options] <token name>",
			Action:    app.actionAccountCreateToken,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "write",
					Usage:       "Grant write permissions in addition to read",
					Destination: &app.writePermission,
				},
				&cli.StringFlag{
					Name:        "expiry",
					Usage:       "Set token expiry date in the form YYYY-MM-DD or never (default 1year)",
					Destination: &app.expiry,
				},
			},
		},
		{
			Name:      "remove-token",
			Usage:     "Remove an authentication token from your account *beta*",
			UsageText: "earthly [options] account remove-token <token>",
			Action:    app.actionAccountRemoveToken,
		},
		{
			Name:  "reset",
			Usage: "Reset Earthly account password",
			UsageText: "earthly [options] account reset --email <email>\n" +
				"   earthly [options] account reset --email <email> --token <token>\n",
			Action: app.actionAccountReset,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "email",
					Usage:       "Email address",
					Destination: &app.email,
				},
				&cli.StringFlag{
					Name:        "token",
					Usage:       "Authentication token",
					Destination: &app.token,
				},
			},
		},
	}
}

func (app *earthlyApp) actionAccountRegister(cliCtx *cli.Context) error {
	app.commandName = "accountRegister"
	if app.email == "" {
		return errors.New("no email given")
	}

	if !strings.Contains(app.email, "@") {
		return errors.New("email is invalid")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	if app.token == "" {
		err := cloudClient.RegisterEmail(cliCtx.Context, app.email)
		if err != nil {
			return errors.Wrap(err, "failed to register email")
		}
		fmt.Printf("An email has been sent to %q containing a registration token\n", app.email)
		return nil
	}

	var publicKeys []*agent.Key
	if app.sshAuthSock != "" {
		var err error
		publicKeys, err = cloudClient.GetPublicKeys(cliCtx.Context)
		if err != nil {
			return err
		}
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass app.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	pword := app.password
	if app.password == "" {
		pword, err = promptPassword()
		if err != nil {
			return err
		}
	}

	var interactiveAccept bool
	if !app.termsConditionsPrivacy {
		rawAccept, err := promptInput(cliCtx.Context, "I acknowledge Earthly Technologies’ Privacy Policy (https://earthly.dev/privacy-policy) and agree to Earthly Technologies Terms of Service (https://earthly.dev/tos) [y/N]: ")
		if err != nil {
			return err
		}
		if rawAccept == "" {
			rawAccept = "n"
		}
		accept := strings.ToLower(rawAccept)[0]
		interactiveAccept = (accept == 'y')
	}
	termsConditionsPrivacy := app.termsConditionsPrivacy || interactiveAccept

	var publicKey string
	if app.registrationPublicKey == "" {
		if len(publicKeys) > 0 {
			rawIsRegisterSSHKey, err := promptInput(cliCtx.Context, "Would you like to enable password-less login using public key authentication (https://earthly.dev/public-key-auth)? [Y/n]: ")
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
				keyNum, err := promptInput(cliCtx.Context, "enter key number (1=default): ")
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
		_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(app.registrationPublicKey))
		if err == nil {
			// supplied public key is valid
			publicKey = app.registrationPublicKey
		} else {
			// otherwise see if it matches the name (Comment) of a key known by the ssh agent
			for _, key := range publicKeys {
				if key.Comment == app.registrationPublicKey {
					publicKey = key.String()
					break
				}
			}
			if publicKey == "" {
				return errors.Errorf("failed to find key in ssh agent's known keys")
			}
		}
	}

	err = cloudClient.CreateAccount(cliCtx.Context, app.email, app.token, pword, publicKey, termsConditionsPrivacy)
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	fmt.Println("Account registration complete")
	return nil
}

func promptPassword() (string, error) {
	fmt.Printf("New password: ")
	enteredPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println("")
	fmt.Printf("Confirm password: ")
	enteredPassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println("")
	if string(enteredPassword) != string(enteredPassword2) {
		return "", errors.Errorf("passwords do not match")
	}
	return string(enteredPassword), nil
}

func (app *earthlyApp) actionAccountListKeys(cliCtx *cli.Context) error {
	app.commandName = "accountListKeys"
	cloudClient, err := app.newCloudClient()
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

func (app *earthlyApp) actionAccountAddKey(cliCtx *cli.Context) error {
	app.commandName = "accountAddKey"
	cloudClient, err := app.newCloudClient()
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
	// as there's no way to pass app.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Which of the following public keys do you want to register (your private key is never sent or even read by earthly)?\n")
	for i, key := range publicKeys {
		fmt.Printf("%d) %s\n", i+1, key.String())
	}
	keyNum, err := promptInput(cliCtx.Context, "enter key number (1=default): ")
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
			app.console.Warnf("failed to authenticate using newly added public key: %s", err.Error())
			return nil
		}
		fmt.Printf("Switching from password-based login to ssh-based login\n")
	}

	return nil
}

func (app *earthlyApp) actionAccountRemoveKey(cliCtx *cli.Context) error {
	app.commandName = "accountRemoveKey"
	cloudClient, err := app.newCloudClient()
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
func (app *earthlyApp) actionAccountListTokens(cliCtx *cli.Context) error {
	app.commandName = "accountListTokens"
	cloudClient, err := app.newCloudClient()
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
		lastAccessedAtStr := "\t"
		if !token.LastAccessedAt.UTC().IsZero() {
			lastAccessedAtStr = fmt.Sprintf("\t%s UTC", token.LastAccessedAt.UTC().Format("2006-01-02T15:04"))
		}
		fmt.Fprint(w, lastAccessedAtStr)
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	return nil
}
func (app *earthlyApp) actionAccountCreateToken(cliCtx *cli.Context) error {
	app.commandName = "accountCreateToken"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}

	var expiry *time.Time
	if app.expiry != "" && app.expiry != "never" {
		layouts := []string{
			"2006-01-02",
			time.RFC3339,
		}

		var err error
		for _, layout := range layouts {
			var parsedTime time.Time
			parsedTime, err = time.Parse(layout, app.expiry)
			if err == nil {
				expiry = &parsedTime
				break
			}
		}
		if err != nil {
			return errors.Errorf("failed to parse expiry %q", app.expiry)
		}
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	name := cliCtx.Args().First()
	token, err := cloudClient.CreateToken(cliCtx.Context, name, app.writePermission, expiry)
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
func (app *earthlyApp) actionAccountRemoveToken(cliCtx *cli.Context) error {
	app.commandName = "accountRemoveToken"
	if cliCtx.NArg() != 1 {
		return errors.New("invalid number of arguments provided")
	}
	name := cliCtx.Args().First()
	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}
	err = cloudClient.RemoveToken(cliCtx.Context, name)
	if err != nil {
		return errors.Wrap(err, "failed to remove account tokens")
	}
	return nil
}

func (app *earthlyApp) actionAccountLogin(cliCtx *cli.Context) error {
	app.commandName = "accountLogin"
	email := app.email
	token := app.token
	pass := app.password

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
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
	if app.authToken != "" {
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
		app.console.Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		app.printLogSharingMessage()
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
			err = cloudClient.Authenticate(cliCtx.Context)
			if err != nil {
				return errors.Wrap(err, "authentication with cloud server failed")
			}
			app.console.Printf("Logged in as %q using ssh auth\n", email)
			app.printLogSharingMessage()
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
		app.console.Printf("Logged in as %q using %s auth\n", loggedInEmail, authType)
		app.printLogSharingMessage()
		return nil
	default:
		return err
	}

	if email == "" && token == "" {
		if app.sshAuthSock == "" {
			app.console.Warnf("No ssh auth socket detected; falling back to password-based login\n")
		}

		emailOrToken, err := promptInput(cliCtx.Context, "enter your email or auth token: ")
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
		// as there's no way to pass app.ctx to stdin read calls.
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
		app.console.Printf("Logged in as %q using token auth\n", email) // TODO display if using read-only token
		app.printLogSharingMessage()
	} else {
		err = app.loginAndSavePasswordCredentials(cliCtx.Context, cloudClient, email, string(pass))
		if err != nil {
			return err
		}
		app.printLogSharingMessage()
	}
	err = cloudClient.Authenticate(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "authentication with cloud server failed")
	}
	return nil
}

func (app *earthlyApp) printLogSharingMessage() {
	app.console.Printf("Log sharing is enabled by default. If you would like to disable it, run:\n" +
		"\n" +
		"\tearthly config global.disable_log_sharing true")
}

func (app *earthlyApp) loginAndSavePasswordCredentials(ctx context.Context, cloudClient *cloud.Client, email, password string) error {
	err := cloudClient.SetPasswordCredentials(ctx, email, password)
	if err != nil {
		return err
	}
	app.console.Printf("Logged in as %q using password auth\n", email)
	app.console.Printf("Warning unencrypted password has been stored under ~/.earthly/auth.credentials; consider using ssh-based auth to prevent this.\n")
	return nil
}

func (app *earthlyApp) actionAccountLogout(cliCtx *cli.Context) error {
	app.commandName = "accountLogout"

	if app.authToken != "" {
		return errLogoutHasNoEffectWhenAuthTokenSet
	}

	cloudClient, err := app.newCloudClient()
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

func (app *earthlyApp) actionAccountReset(cliCtx *cli.Context) error {
	app.commandName = "accountReset"

	if app.email == "" {
		return errors.New("no email given")
	}

	cloudClient, err := app.newCloudClient()
	if err != nil {
		return err
	}

	if app.token == "" {
		err = cloudClient.AccountResetRequestToken(cliCtx.Context, app.email)
		if err != nil {
			return errors.Wrap(err, "failed to request account reset token")
		}
		app.console.Printf("An account reset token has been emailed to %q\n", app.email)
		return nil
	}

	// Our signal handling under main() doesn't cause reading from stdin to cancel
	// as there's no way to pass app.ctx to stdin read calls.
	signal.Reset(syscall.SIGINT, syscall.SIGTERM)

	pword, err := promptPassword()
	if err != nil {
		return err
	}

	err = cloudClient.AccountReset(cliCtx.Context, app.email, app.token, pword)
	if err != nil {
		return errors.Wrap(err, "failed to reset account")
	}
	app.console.Printf("Account password has been reset\n")
	err = app.loginAndSavePasswordCredentials(cliCtx.Context, cloudClient, app.email, pword)
	if err != nil {
		return err
	}
	app.printLogSharingMessage()
	return nil
}
