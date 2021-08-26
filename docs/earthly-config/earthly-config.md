# Earthly configuration file

Global configuration values for earthly can be stored on disk in the configuration file.

By default, earthly reads the configuration file `~/.earthly/config.yml`; however, it can also be
overridden with the `--config` command flag option.

## Format

The earthly config file is a [yaml](https://yaml.org/) formatted file that looks like:

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

Specifies the total size of the BuildKit cache, in MB. The BuildKit daemon uses this setting to configure automatic garbage collection of old cache. A value of 0 causes the size to be adaptive depending on how much space is available on your system. The default is 0.

### disable_analytics

When set to true, disables collecting command line analytics; otherwise, earthly will report anonymized analytics for invokation of the earthly command. For more information see the [data collection page](../data-collection/data-collection.md).

### buildkit_additional_args

This option allows you to pass additional options to Docker when starting up the Earthly buildkit daemon. For example, this can be used to bypass user namespacing like so:

```yaml
global:
  buildkit_additional_args: ["--userns", "host"]
```

### buildkit_additional_config

This option allows you to pass additional options to Buildkit. For example, this can be used to specify additional CA certificates:

```yaml
global:
  buildkit_additional_args: ["-v", "<absolute-path-to-ca-file>:/etc/config/add.ca"]
  buildkit_additional_config: |
    [registry."<registry-hostname>"]
      ca=["/etc/config/add.ca"]
```

### cni_mtu

Allows overriding Earthly's automatic MTU detection. This is used when configuring the Buildkit internal CNI network. MTU must be between 64 and 65,536.

### ip_tables

Allows overriding Earthly's automatic `ip_tables` module detection. Valid choices are `iptables-legacy` or `iptables-nft`.

### no_loop_device (obsolete)

This option is obsolete and it is ignored. Earthly no longer uses a loop device for its cache.

### cache_path (obsolete)

This option is obsolete and it is ignored. Earthly cache has moved to a Docker volume. For more information see the [page on managing cache](../guides/cache.md).

## Git configuration reference

All git configuration is contained under site-specific options.

### site-specific options

#### site

The git repository hostname. For example `github.com`, or `gitlab.com`

#### auth

Either `ssh`, `https`, or `auto` (default). If `https` is specified, user and password fields are used
to authenticate over https when pulling from git for the corresponding site. If `auto` is specified
earthly will use `ssh` when the ssh-agent is running and has at least one key loaded, and will fallback
to using `https` when no ssh-keys are present.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication.

#### user

The https username to use when auth is set to `https`. This setting is ignored when auth is `ssh`.

#### password

The https password to use when auth is set to `https`. This setting is ignored when auth is `ssh`.

#### pattern

A regular expression defined to match git URLs, defaults to the `<site>/([^/]+)/([^/]+)`. For example if the site is `github.com`, then the default pattern will
match `github.com/<user>/<repo>`.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication with self-hosted git repositories.

See the [RE2 docs](https://github.com/google/re2/wiki/Syntax) for a complete definition of the supported regular expression syntax.


#### substitute

If specified, a regular expression substitution will be preformed to determine which URL is cloned by git. Values like `$1`, `$2`, ... will be replaced
with matched subgroup data. If no substitute is given, a URL will be created based on the requested SSH authentication mode.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication with self-hosted git repositories.
