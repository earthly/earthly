# Earthly configuration file

Global configuration values for earth can be stored on disk in the configuration file.

By default, earth reads the configuration file `~/.earthly/config.yaml`; however, it can also be
overridden with the `--config` command flag option.

## Format

The earthly config file is a [yaml](https://yaml.org/) formatted file that looks like:

```yaml
global:
  cache_path: <cache_path>
  cache_size_mb: <cache_size_mb>
  no_loop_device: true|false
  buildkit_image: <buildkit_image>
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

## Global configuration reference

### cache_path

Specifies the location where build data is stored. The default location is `/var/cache/earthly` on Linux and `~/Library/Caches/earthly` on macOS.

### cache_size_mb

Specifies the total size of the buildkit cache, in MB. The buildkit daemon will allocate disk space for this size. Size less than `1000` (1GB) is not recommended. The default size if this option is not set is `10000` (10GB).

This setting is only used when the cache is initialized for the first time. In order to apply the setting immediately, issue the following command after changing the configuration

```bash
earth prune --reset
```

### no_loop_device

When set to true, disables the use of a loop device for storing the cache. By default, Earthly uses a file mounted as a loop device, so that it can control the type of filesystem used for the cache, in order to ensure that overlayfs can be mounted on top of it. If you are already using a filesystem compatible with overlayfs, then you can disable the loop device.

## Git configuration reference

The git configuration is split up into global config options, or site-specific options.

### global options

The global git options 

#### url_instead_of

Rewrites git URLs of a certain pattern. Similar to [`git-config url.<base>.insteadOf`](https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf).
Multiple values can be separated by commas. Format: `<base>=<instead-of>[,...]`.

This setting allows rewriting all git URLs of the form `https://example...` into `git@example.com:...`, or vice-versa.

For example:

* `--git-url-instead-of='git@example.com:=https://example.com/'` forces use of SSH-based URLs rather than HTTPS
* `--git-url-instead-of='https://localmirror.example.com/=git@example.com:'` forces use of HTTPS-based local mirror for ssh-based example.com repositories

NOTE: if the `auth` option is configured under a site-specific configuration, then the appropriate rewriting rule will be automatically applied.

### site-specific options

#### site

The git repository hostname. For example `github.com`, or `gitlab.com`

#### auth

Either `https` or `ssh` (default). If https is specified, user and password fields are used
to authenticate over https when pulling from git for the corresponding site.

See the [Authentication guide](../guides/auth.md) for a guide on setting up authentication.

#### user

The https username to use when auth is set to `https`. This setting is ignored when auth is `ssh`.

#### password

The https password to use when auth is set to `https`. This setting is ignored when auth is `ssh`.
