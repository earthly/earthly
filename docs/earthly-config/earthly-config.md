# Earthly configuration file

Global configuration values for earthly can be stored on disk in the configuration file.

By default, earthly reads the configuration file `~/.earthly/config.yml`; however, it can also be
overridden with the `--config` command flag option.

## Format

The earthly config file is a [YAML](https://yaml.org/) formatted file that looks like:

```yaml
global:
  cache_size_mb: <cache_size_mb>
git:
    global:
        url_instead_of: <url_instead_of>
    <site>:
        auth: https|ssh
        user: <username>
        password: <password>
    <site2>:
        ...
```

Example:

```yaml
global:
    cache_size_mb: 20000
git:
    global:
        url_instead_of: "git@example.com:=https://localmirror.example.com/"
    github.com:
        auth: https
        user: alice
        password: itsasecret
```

{% hint style='info' %}
##### Tip
To quickly change a configuration item via the `earthly` command, you can use [`earthly config`](../earthly-command/earthly-command.md#earthly-config).

```bash
earthly config <key> <value>
```

For example

```bash
earthly config global.cache_size_mb 20000
```
{% endhint %}

## Global configuration reference

### cache_size_mb

Specifies the total size of the BuildKit cache, in MB. The BuildKit daemon uses this setting to configure automatic garbage collection of old cache.
Setting this to 0, either explicitly or by omission, will cause buildkit to use its internal default of 10% of the root filesystem.

### cache_size_pct

Specifies the total size of the BuildKit cache, as a percentage (0-100) of the total filesystem size.
When used in combination with `cache_size_mb`, the lesser of the two values will be used. This limit is ignored when set to 0.

### secret_provider (experimental)

A custom user-supplied program to call which returns a secret for use by earthly. The secret identifier is passed as the first argument to the program.

If no secret is found, the program can instruct earthly to continue searching for secrets under `.secret`, by exiting with a status code of `2`, all other non-zero
status codes will cause earthly to exit.

For example, if you have:

```yaml
config:
  secret_provider: my-secret-provider
```

and `my-secret-provider` (which is accessible on your `PATH`):

```bash
#!/bin/sh
set -e

if [ "$1" = "mysecret" ]; then
    echo -n "open sesame"
    exit 0
fi

exit 2
```

Then when earthly encounters a command that requires a secret, such as

```Dockerfile
RUN --secret mysecret echo "the passphrase is $mysecret."
```

earthly will request the secret for `mysecret` by calling `my-secret_provider mysecret`.

{% hint style='info' %}
##### Note

All stdout data will be used as the secret value, including whitespace (and newlines).
You may want to use `echo -n` to prevent returning a newline.

Any data sent to stderr will be displayed on the earthly console, this makes it possible
to insert commands such as `echo >&2 "here is some debug text"` without affecting the contents
of the secret.

{% endhint %}

### disable_analytics

When set to true, disables collecting command line analytics; otherwise, earthly will report anonymized analytics for invocation of the earthly command. For more information see the [data collection page](../data-collection/data-collection.md).

### disable_log_sharing

When set to true, disables sharing build logs after each build. This setting applies to logged-in users only.

### conversion_parallelism

The number of concurrent converters for speeding up build targets that use blocking commands like `IF`, `WITH DOCKER --load`, `FROM DOCKERFILE` and others.

### buildkit_max_parallelism

The maximum parallelism configured for the buildkit daemon workers. The default is 20.

{% hint style='info' %}
##### Note

Set this configuration to a lower value if your machine is resource constrained and performs poorly when running too many builds in parallel.

{% endhint %}

### buildkit_additional_args

This option allows you to pass additional options to Docker when starting up the Earthly BuildKit daemon. 
Note that changes to these values will trigger earthly to restart buildkit on the next run.

#### Bypass User Namespacing

The `--userns` flag can be set as follows:

```yaml
global:
  buildkit_additional_args: ["--userns", "host"]
```

#### Session Timeout

By default, Buildkit will automatically cancel sessions (i.e. individual builds) after 24 hours.
This value can be overriden using the following option:

```yaml
global:
  buildkit_additional_args: ["-e", "BUILDKIT_SESSION_TIMEOUT=72h"]
```

Note that setting a value of zero `0` here will disable the feature entirely.
This can be useful in cases where long-lived interactive sessions are used.

### buildkit_additional_config

This option allows you to pass additional options to BuildKit.
Note that changes to these values will trigger earthly to restart buildkit on the next run.


#### Additional CA Certificates

Additional CA certificates can be passed in to buildkit. This also requires a corresponding change in `buildkit_additional_args`.

```yaml
global:
  buildkit_additional_args: ["-v", "<absolute-path-to-ca-file>:/etc/config/add.ca"]
  buildkit_additional_config: |
    [registry."<registry-hostname>"]
      ca=["/etc/config/add.ca"]
```

### cni_mtu

Allows overriding Earthly's automatic MTU detection. This is used when configuring the BuildKit internal CNI network. MTU must be between 64 and 65,536.

### ip_tables

Allows overriding Earthly's automatic `ip_tables` module detection. Valid choices are `iptables-legacy` or `iptables-nft`.

### no_loop_device (obsolete)

This option is obsolete and it is ignored. Earthly no longer uses a loop device for its cache.

### git_image

Allows to override the image used to run internal `git` commands (e.g. during `GIT CLONE` or `IMPORT`). This defaults to `alpine/git:v2.30.1`.

### org

The default organization to use when performing Earthly operations that require an organization. Ignored when  the `--org` CLI option is present, or when the `EARTHLY_ORG` environment variable are set.

### Frontend configuration

This option allows you to specify what supported frontend you are using (Docker / Podman).
By default, Earthly will attempt to discover the frontend in this order: Docker -> Podman -> None

For Docker:
```yaml
global:
  container_frontend: docker-shell
```

For Podman:
```yaml
global:
  container_frontend: podman-shell
```

You can use the following command to set the configuration option using the earthly CLI:

```bash
# Docker
earthly config 'global.container_frontend' 'docker-shell'

# Podman
earthly config 'global.container_frontend' 'podman-shell'
```

## Git configuration reference

All git configuration is contained under site-specific options.

### site-specific options

#### site

The git repository hostname. For example `github.com`, or `gitlab.com`

#### auth

Either `ssh`, `https`, or `auto` (default). If `https` is specified, user and password fields are used
to authenticate over HTTPS when pulling from git for the corresponding site. If `auto` is specified
earthly will use `ssh` when the ssh-agent is running and has at least one key loaded, and will fallback
to using `https` when no ssh-keys are present.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication.

#### user

The HTTPS username to use when auth is set to `https`. This setting is ignored when auth is `ssh`.

#### password

The HTTPS password to use when auth is set to `https`. This setting is ignored when auth is `ssh`.

#### strict_host_key_checking

The `strict_host_key_checking` option can be used to control access to ssh-based repos whose key is not known or has changed.
Strict host key checking is enabled by default, setting it to `false` disables host key checking.
This setting is only used when auth is `ssh`.

{% hint style='info' %}
##### Tip
Disabling strict host key checking is a bad security practice (as it makes a man-in-the-middle attack possible).
Instead, it's recommended to record the host's ssh key to `~/.ssh/known_hosts`; this can be done by running

```bash
ssh-keyscan <hostname> >> ~/.ssh/known_hosts
```
{% endhint %}

#### ssh_command

The `ssh_command` option can be used to override the ssh command that is used by `git` when connecting to an ssh-based repository.
For example, if you need to connect to an outdated sshd-server which only supports the insecure RSA signature algorithm, you could set the `ssh_command` to `ssh -o 'PubKeyAcceptedKeyTypes +ssh-rsa'`.

#### port

Connect using a non-standard git port, e.g. `2222`.

#### prefix

The `prefix` option is used to indicate where git repositories are stored on the server, e.g. `/var/git/`.

#### pattern

A regular expression defined to match git URLs, defaults to the `<site>/([^/]+)/([^/]+)`. For example if the site is `github.com`, then the default pattern will
match `github.com/<user>/<repo>`.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication with self-hosted git repositories.

See the [RE2 docs](https://github.com/google/re2/wiki/Syntax) for a complete definition of the supported regular expression syntax.


#### substitute

If specified, a regular expression substitution will be performed to determine which URL is cloned by git. Values like `$1`, `$2`, ... will be replaced
with matched subgroup data. If no substitute is given, a URL will be created based on the requested SSH authentication mode.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication with self-hosted git repositories.
