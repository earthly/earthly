# Earthly command reference

## earthly

#### Synopsis

* Target form
  ```
  earthly [options...] <target-ref> [build-args...]
  ```
* Artifact form
  ```
  earthly [options...] --artifact|-a <target-ref>/<artifact-path> [<dest-path>]
  earthly [options...] --artifact|-a (<target-ref>/<artifact-path> [build-args...]) [<dest-path>]
  ```
* Image form
  ```
  earthly [options...] --image <target-ref> [build-args...]
  ```

#### Description

The command executes a build referenced by `<target-ref>` (*target form* and *image form*) or `<artifact-ref>` (*artifact form*). In the *target form*, the referenced target and its dependencies are built. In the *artifact form*, the referenced artifact and its dependencies are built, but only the specified artifact is output. The output path of the artifact can be optionally overridden by `<dest-path>`. In the *image form*, the image produced by the referenced target and its dependencies are built, but only the specified image is output.

If a BuildKit daemon has not already been started, and the option `--buildkit-host` is not specified, this command also starts up a container named `earthly-buildkitd` to act as a build daemon.

The execution has four phases:

* Init
* Build
* Push (optional - disabled by default)
* Local output (optional - enabled by default)

During the init phase the configuration is interpreted and the BuildKit daemon is started (if applicable). During the build phase, the referenced target and all its direct or indirect dependencies are executed. During the push phase, when enabled, Earthly performs image pushes and it also runs `RUN --push` commands.  During the local output phase, all applicable artifacts with an `AS LOCAL` specification are written to the specified output location, and all applicable docker images are loaded onto the host's docker daemon.

If the build phase does not succeed, no output is produced and no push instruction is executed. In this case, the command exits with a non-zero exit code.

#### Target and Artifact Reference

The `<target-ref>` can reference both local and remote targets.

##### Local Reference

`+<target-name>` will reference a target in the local Earthfile in the current directory.

`<local-path>+<target-name>` will reference a local Earthfile in a different directory as
specified by `<local-path>`, which must start with `./`, `../`, or `/`.

##### Remote Reference

`<gitvendor>/<namespace>/<project>/path/in/project[:some-tag]+<target-name>` will access a remote git repository.

##### Artifact Reference

The `<artifact-ref>` can reference artifacts built by targets. `<target-ref>/<artifact-path>` will reference a build target's artifact.

##### Examples

See the [Target, artifact, and image referencing guide](../guides/target-ref.md) for more details and examples.

#### Build args

Synopsis:

  * Target form `earthly <target-ref> [--<build-arg-key>=<build-arg-value>...]`
  * Artifact form `earthly --artifact (<target-ref>/<artifact-path> [--<build-arg-key>=<build-arg-value>...]) <dest-path>`
  * Image form `earthly --image <target-ref> [--<build-arg-key>=<build-arg-value>...]`

Also available as an env var setting: `EARTHLY_BUILD_ARGS="<build-arg-key>=<build-arg-value>,<build-arg-key>=<build-arg-value>,..."`.

Build arg overrides may be specified as part of the Earthly command. The value of the build arg `<build-arg-key>` is set to `<build-arg-value>`.

In the target and image forms the build args are passed after the target reference. For example `earthly +some-target --NAME=john --SPECIES=human`. In the artifact form, the build args are passed immediately after the artifact reference, however they are surrounded by parenthesis, similar to a [`COPY` command](../earthfile/earthfile.md#copy). For example `earthly --artifact (+some-target/some-artifact --NAME=john --SPECIES=human) ./dest/path/`.

The build arg overrides only apply to the target being called directly and any other target referenced as part of the same Earthfile. Build arg overrides, will not apply to targets referenced from other directories or other repositories.

For more information about build args see the [`ARG` Earthfile command](../earthfile/earthfile.md#arg).

#### Environment Variables and .env File

As specified under the [options section](#options), all flag options have an environment variable equivalent, which can be used as an alternative.

Furthermore, additional environment variables are also read from a file named `.env`, if one exists in the current directory. The syntax of the `.env` file is of the form

```.env
<NAME_OF_ENV_VAR>=<value>
...
```

as one variable per line, without any surrounding quotes. If quotes are included, they will become part of the value. Lines beginning with `#` are treated as comments. Blank lines are allowed. Here is a simple example:

```.env
# Settings
EARTHLY_ALLOW_PRIVILEGED=true
MY_SETTING=a setting which contains spaces

# Secrets
MY_SECRET=MmQ1MjFlY2UtYzhlNi00YjJkLWI5YTMtNjIzNzJmYjcwOTJk
ANOTHER_SECRET=MjA5YjU2ZTItYmIxOS00MDQ3LWFlNzYtNmQ5NGEyZDFlYTQx
```

{% hint style='info' %}
##### Note
The directory used for loading the `.env` file is the directory where `earthly` is called from and not necessarily the directory where the Earthfile is located in.
{% endhint %}

The additional environment variables specified in the `.env` file are loaded by `earthly` in three distinct ways:

* **Setting options for `earthly` itself** - the settings are loaded if they match the environment variable equivalent of an `earthly` option.
* **Build args** - the settings are passed on to the build and are used to override any [`ARG`](../earthfile/earthfile.md#arg) declaration.
* **Secrets** - the settings are passed on to the build to be referenced via the [`RUN --secret`](../earthfile/earthfile.md#secret-less-than-env-var-greater-than-less-than-secret-ref-greater-than) option.

{% hint style='danger' %}
##### Important
The `.env` file is meant for settings which are specific to the local environment the build executes in. These settings may cause inconsistencies in the way the build executes on different systems, leading to builds that are difficult to reproduce. Keep the contents of `.env` files to a minimum to avoid such issues.
{% endhint %}

#### Global Options

##### `--config <path>`

Also available as an env var setting: `EARTHLY_CONFIG=<path>`.

Overrides the earthly [configuration file](../earthly-config/earthly-config.md), defaults to `~/.earthly/config.yml`.

##### `--installation-name <name>`

Also available as an env var setting: `EARTHLY_INSTALLATION_NAME=<name>`.

Overrides the Earthly installation name. The installation name is used for the Buildkit Daemon name, the cache volume name, the configuration directory (`~/.<installation-name>`) and for the ports used by Buildkit. Using multiple installation names on the same system allows Earthly to run as multiple isolated instances, each with its own configuration, cache and daemon. Defaults to `earthly`.

##### `--ssh-auth-sock <path-to-sock>`

Also available as an env var setting: `EARTHLY_SSH_AUTH_SOCK=<path-to-sock>`.

Sets the path to the SSH agent sock, which can be used for SSH authentication. SSH authentication is used by Earthly in order to perform git clone's underneath.

On Linux systems, this setting defaults to the value of the env var $SSH_AUTH_SOCK. On most systems, the env var `SSH_AUTH_SOCK` env var is already set if an SSH agent is running.

On Mac systems, this setting defaults to `/run/host-services/ssh-auth.sock` to match recommendation in [the official Docker documentation](https://docs.docker.com/docker-for-mac/osxfs/#ssh-agent-forwarding).

For more information see the [Authentication page](../guides/auth.md).

##### `--auth-token <value>`

Also available as an env var setting: `EARTHLY_TOKEN=<value>`.

Force Earthly account login to authenticate with supplied token.

##### `--verbose`

Also available as an env var setting: `EARTHLY_VERBOSE=1`.

Enables verbose logging.

##### `--git-username <git-user>` (**deprecated**)

Also available as an env var setting: `GIT_USERNAME=<git-user>`.

This option is now deprecated. Please use the [configuration file](../earthly-config/earthly-config.md) instead.

##### `--git-password <git-pass>` (**deprecated**)

Also available as an env var setting: `GIT_PASSWORD=<git-pass>`.

This option is now deprecated. Please use the [configuration file](../earthly-config/earthly-config.md) instead.

##### `--git-url-instead-of <git-instead-of>` (**obsolete**)

Also used to be available as an env var setting: `GIT_URL_INSTEAD_OF=<git-instead-of>`.

This option is now obsolete. By default, `earthly` will automatically switch from ssh to HTTPS when no keys are found or the ssh-agent isn't running.
Please use the [configuration file](../earthly-config/earthly-config.md) to override the default behavior.

#### Build Options

Build options are specific to executing Earthly builds; they are simply listed in this section for readability, and can be supplied as global options.

##### `--secret|-s <secret-id>[=<value>]`

Also available as an env var setting: `EARTHLY_SECRETS="<secret-id>=<value>,<secret-id>=<value>,..."`.

Passes a secret with ID `<secret-id>` to the build environments. If `<value>` is not specified, then the value becomes the value of the environment variable with the same name as `<secret-id>`.

The secret can be referenced within Earthfile recipes as `RUN --secret <arbitrary-env-var-name>=<secret-id>`. For more information see the [`RUN --secret` Earthfile command](../earthfile/earthfile.md#run).

##### `--secret-file <secret-id>=<path>`

Also available as an env var setting: `EARTHLY_SECRET_FILES="<secret-id>=<path>,<secret-id>=<path>,..."`.

Loads the contents of a file located at `<path>` into a secret with ID `<secret-id>` for use within the build environments.

The secret can be referenced within Earthfile recipes as `RUN --secret <arbitrary-env-var-name>=<secret-id>`. For more information see the [`RUN --secret` Earthfile command](../earthfile/earthfile.md#run).

##### `--push`

Also available as an env var setting: `EARTHLY_PUSH=true`.

Instructs Earthly to push any docker images declared with the `--push` flag to remote docker registries and to run any `RUN --push` commands. For more information see the [`SAVE IMAGE` Earthfile command](../earthfile/earthfile.md#save-image) and the [`RUN --push` Earthfile command](../earthfile/earthfile.md#run).

Pushing only happens during the output phase, and only if the build has succeeded.

##### `--no-output`

Also available as an env var setting: `EARTHLY_NO_OUTPUT=true`.

Instructs Earthly not to output any images or artifacts. This option cannot be used with the *artifact form* or the *image form*.

##### `--output`

Also available as an env var setting: `EARTHLY_OUTPUT=true`.

Allow artifacts or images to be output, even when running under --ci mode.

##### `--no-cache`

Also available as an env var setting: `EARTHLY_NO_CACHE=true`.

Instructs Earthly to ignore any cache when building. It does, however, continue to store new cache formed as part of the build (to be possibly used on future invocations).

##### `--allow-privileged|-P`

Also available as an env var setting: `EARTHLY_ALLOW_PRIVILEGED=true`.

Permits the build to use the --privileged flag in RUN commands. For more information see the [`RUN --privileged` command](../earthfile/earthfile.md#run).

##### `--use-inline-cache`

Also available as an env var setting: `EARTHLY_USE_INLINE_CACHE=true`

Enables use of inline cache, if available. Any `SAVE IMAGE --push` command is used to inform the system of possible inline cache sources. For more information see the [remote caching guide](../remote-caching.md).

##### `--save-inline-cache`

Also available as an env var setting: `EARTHLY_SAVE_INLINE_CACHE=true`

Enables embedding inline cache in any pushed images. This cache can be used on other systems, if enabled via `--use-inline-cache`. For more information see the [remote caching guide](../remote-caching.md).

##### `--remote-cache <image-tag>`

Also available as an env var setting: `EARTHLY_REMOTE_CACHE=<image-tag>`

Enables use of explicit cache. The provided `<image-tag>` is used for storing and retrieving the cache to/from a Docker registry. Storing explicit cache is only enabled if the option `--push` is also passed in. For more information see the [remote caching guide](../remote-caching.md).

##### `--max-remote-cache`

Also available as an env var setting: `EARTHLY_MAX_REMOTE_CACHE=true`

Enables storing all intermediate layers as part of the explicit cache. Note that this setting is rarely effective due to the excessive upload overhead. For more information see the [remote caching guide](../remote-caching.md).

##### `--ci`

Also available as an env var setting: `EARTHLY_CI=true`

In *target mode*, this option is an alias for

```
--use-inline-cache --save-inline-cache --no-output --strict
```

In *artifact* and *image modes* , this option is an alias for

```
--use-inline-cache --save-inline-cache
```

##### `--platform <platform>`

Also available as an env var setting: `EARTHLY_PLATFORMS=<platform>`.

Sets the platform to build for.

{% hint style='info' %}
##### Note
It is not yet possible to specify multiple platforms through this flag. You may, however, use a wrapping target and a `BUILD` command in your Earthfile:

```Dockerfile
build-all-platforms:
  BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build

build:
  ...
```
{% endhint %}

##### `--build-arg <key>[=<value>]` (**deprecated**)

This option has been deprecated in favor of the new build arg syntax `earthly <target-ref> --<key>=<value>`.

Also available as an env var setting: `EARTHLY_BUILD_ARGS="<key>=<value>,<key>=<value>,..."`.

Overrides the value of the build arg `<key>`. If `<value>` is not specified, then the value becomes the value of the environment variable with the same name as `<key>`. For more information see the [`ARG` Earthfile command](../earthfile/earthfile.md#arg).

##### `--interactive|-i`

Also available as an env var setting: `EARTHLY_INTERACTIVE=true`.

Enable interactive debugging mode. By default when a `RUN` command fails, earthly will display the error and exit. If the interactive mode is enabled and an error occurs, an interactive shell is presented which can be used for investigating the error interactively. Due to technical limitations, only a single interactive shell can be used on the system at any given time.

##### `--strict`

Disallow usage of features that may create unrepeatable builds.

#### Log formatting options

These options can only be set via environment variables, and have no command line equivalent.

| Variable               | Usage                                                                                                                                                                                                      |
|------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| NO_COLOR               | `NO_COLOR=1` disables the use of color.                                                                                                                                                                    |
| FORCE_COLOR            | `FORCE_COLOR=1` forces the use of color.                                                                                                                                                                   |
| EARTHLY_TARGET_PADDING | `EARTHLY_TARGET_PADDING=n` will set the column to the width of `n` characters. If a name is longer than `n`, its path will be truncated and remaining extra length will cause the column to go ragged. |
| EARTHLY_FULL_TARGET    | `EARTHLY_FULL_TARGET=1` will always print the full target name, and leave the target name column ragged.                                                                                                   |

## earthly prune

#### Synopsis

* Standard form
  ```
  earthly [options] prune [--all|-a]
  ```
* Reset form
  ```
  earthly [options] prune --reset
  ```

#### Description

The command `earthly prune` eliminates Earthly cache. In the *standard form* it issues a prune command to the BuildKit daemon. In the *reset form* it restarts the BuildKit daemon, instructing it to completely delete the cache directory on startup, thus forcing it to start from scratch.

#### Options

##### `--all|-a`

Instructs earthly to issue a "prune all" command to the BuildKit daemon.

##### `--reset`

Restarts the BuildKit daemon and completely resets the cache directory.

## earthly config

#### Synopsis

```
earthly [options] config [key] [value]
```

#### Description

Manipulates values in `~/.earthly/config.yml`. It does its best to preserve existing formatting and comments. `[value]` must be a valid YAML literal for the given `[key]`.

#### Options

##### `--help`

Prints help text, along with some examples.

##### `[key] --help`

Prints help for the specific key, including what it is used for and what kind of value it needs to be.

#### Examples

Set your cache size:

```
config global.cache_size_mb 1234
```

Set additional BuildKit args, using a YAML array:

```
config global.buildkit_additional_args ['userns', '--host']
```

Set a key containing a period:

```
config git."example.com".password hunter2
```

Set up a whole custom git repository for a server called example.com, using a single-line YAML literal:
* which stores git repos under /var/git/repos/name-of-repo.git
* allows access over ssh
* using port 2222
* sets the username to git
* is recognized to earthly as example.com/name-of-repo

```
config git "{example: {pattern: 'example.com/([^/]+)', substitute: 'ssh://git@example.com:2222/var/git/repos/\$1.git', auth: ssh, user: git}}"
```

The above command yields the following config file:

```yaml
git:
    example:
        pattern: example.com/([^/]+)
        substitute: ssh://git@example.com:2222/var/git/repos/$1.git
        auth: ssh
        user: git
```

## earthly account

Contains sub-commands for registering and administration an Earthly account.

#### earthly account register

###### Synopsis

* ```
  earthly account register --email <email>
  earthly account register --email <email> --token <email-verification-token> [--password <password>] [--public-key <public-key>] [--accept-terms-conditions-privacy]
  ```

###### Description

Register for an Earthly account. Registration is done in two steps: first run the register command with only the --email argument, this will then send an email to the
supplied email address with a registration token (which is used to verify your email address), second re-run the register command with both the --email and --token arguments
to complete the registration process.

#### earthly account login

###### Synopsis

* ```
  earthly [options] account login
  earthly [options] account login --email <email>
  earthly [options] account login --email <email> --password <password>
  earthly [options] account login --token <token>
  ```

###### Description

Login to an existing Earthly account. If no email or token is given, earthly will attempt to login using [registered public keys](../public-key-auth/public-key-auth.md).

#### earthly account logout

###### Synopsis

* ```
  earthly [options] account logout
  ```

###### Description

Removes cached login information from `~/.earthly/auth.token`.

#### earthly account list-keys

###### Synopsis

* ```
  earthly account list-keys
  ```

###### Description

Lists all public keys that are authorized to login to the current Earthly account.

#### earthly account add-key

###### Synopsis

* ```
  earthly account add-key [<public-key>]
  ```

###### Description

Authorize a new public key to login to the current Earthly account. If `key` is omitted, an interactive prompt is displayed
to select a public key to add.

#### earthly account remove-key

###### Synopsis

* ```
  earthly account remove-key <public-key>
  ```

###### Description

Removes an authorized public key from accessing the current Earthly account.

#### earthly account list-tokens

###### Synopsis

* ```
  earthly account list-tokens
  ```

###### Description

List account tokens associated with Earthly account. A token is useful for environments where the ssh-agent is not accessible (e.g. a CI system).

#### earthly account create-token

###### Synopsis

* ```
  earthly account create-token [--write] [--expiry <expiry>] <token-name>
  ```

###### Description

Creates a new authentication token. A read-only token is created by default, If the `--write` flag is specified the token will have read+write access.
The token will expire in 1 year from creation date unless a different date is supplied via the `--expiry` option.

{% hint style='info' %}
It is then possible to `export EARTHLY_TOKEN=...`, which will force earthly to use this token for all authentication (overriding any other currently-logged in sessions).
{% endhint %}

#### earthly account remove-token

###### Synopsis

* ```
  earthly account remove-token <token>
  ```

###### Description

Removes a token from the current Earthly account.


## earthly org

Contains sub-commands for creating and managing Earthly organizations.

#### earthly org create

###### Synopsis

* ```
  earthly org create <org-name>
  ```

###### Description

Create a new organization, which can be used to share secrets between different user accounts.

#### earthly org list

###### Synopsis

* ```
  earthly org list
  ```

###### Description

List all organizations the current account is a member, or administrator of.

#### earthly org list-permissions

###### Synopsis

* ```
  earthly org list-permissions <org-name>
  ```

###### Description

List all accounts and the paths they have permission to access under a particular organization.

#### earthly org invite

###### Synopsis

* ```
  earthly org invite [--write] <org-path> <email> [<email>, ...]
  ```

###### Description

Invites a user into an organization; `<org-path>` can either be a top-level org access by granting permission on `/<org-name>/`, or finer-grained access can be granted to a subpath e.g. `/<org-name>/path/to/share/`.
By default users are granted read-only access unless the `--write` flag is given.

#### earthly org revoke

###### Synopsis

* ```
  earthly org revoke <org-path> <email> [<email>, ...]
  ```

###### Description

Revokes a previously invited user from an organization.

## earthly secrets

Contains sub-commands for creating and managing Earthly secrets.

#### earthly secrets set

###### Synopsis

* ```
  earthly secrets set <path> <value>
  earthly secrets set --file <local-path> <path>
  ```

###### Description

Stores a secret in the secrets store

#### earthly secrets get

###### Synopsis

* ```
  earthly secrets get [-n] <path>
  ```

###### Description

Retrieve a secret from the secrets store. If `-n` is given, no newline is printed after the contents of the secret.

#### earthly secrets ls

###### Synopsis

* ```
  earthly secrets ls [<path>]
  ```

###### Description

List secrets the current account has access to.

#### earthly secrets rm

###### Synopsis

* ```
  earthly secrets rm <path>
  ```

###### Description

Removes a secret from the secrets store.


## earthly bootstrap

#### Synopsis

* ```
  earthly bootstrap
  ```

#### Description

Performs initialization tasks needed for `earthly` to function correctly. This command can be re-run to fix broken setups. It is recommended to run this with sudo. 

#### Options

##### `--no-buildkit`

Skips setting up the BuildKit container during bootstrapping. If needed, it will also be performed when a build is ran.

##### `--with-autocomplete`

Installs shell autocompletions during bootstrap. Requires `sudo` to install them correctly.

## earthly --help

#### Synopsis

* ```
  earthly --help
  ```
* ```
  earthly <command> --help
  ```

#### Description

Prints help information about earthly.

## earthly --version

#### Synopsis

* ```
  earthly --version
  ```

#### Description

Prints version information about earthly.

## earthly ls

#### Synopsis

* ```
  earthly ls [<project-ref>]
  ```

#### Description

Prints all targets in an `Earthfile` in a project.

#### Options

##### `--args`

Show arguments (`ARG` statements) in the targets.

##### `--long`

Show full, canonical target references (includes the project part of the reference, if applicable).

## earthly doc

#### Synopsis

* ```
  earthly doc [<project-ref>[+<target-ref>]]
  ```

#### Description

Prints documentation comments for documented targets in an `Earthfile` in a
project. Documentation on a target is any comment block that ends on the line
immediately above the target definition and begins with the name of the target.

#### Examples

Given the following `Earthfile`:

```
VERSION 0.7
FROM golang:1.19-alpine3.15

deps:
    COPY go.mod go.sum .
    RUN go mod download

# build runs 'go build' and saves the artifact locally.
build:
    FROM +deps
    COPY . .
    ARG output=./build/something
    RUN go build -o /bin/something
    SAVE ARTIFACT /bin/something AS LOCAL $output

# tidy runs 'go mod tidy' and saves go.mod/go.sum locally.
tidy:
    FROM +deps
    COPY . .
    RUN go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum
```

##### Print the doc comments for all documented targets:

```
$ earthly doc
TARGETS:
  +build
    build runs 'go build' and saves the artifact locally.
  +tidy
    tidy runs 'go mod tidy' and saves go.mod/go.sum locally.
```

Note that, unlike `earthly ls`, `earthly doc` does not mention the `deps`
target. Since it has no documentation, the `deps` target is not included in the
output.

##### Print the doc comments for a specific target:

```
$ earthly doc +build
+build
  build runs 'go build' and saves the artifact locally.
```

## earthly web

#### Synopsis

* ```
  earthly web [--provider=<provider-ref>]]
  ```

#### Description

Prints a url for entering the CI application and attempts to open your default browser with that url.
If the provider argument is given the CI application will automatically begin an OAuth flow with the given provider.
If you are logged into the CLI the url will contain a token used to link your OAuth credentials to your Earthly user.

#### Examples

##### Login to the CI application with GitHub

* ```
  earthly web --provider=github
  ```
