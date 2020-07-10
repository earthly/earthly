# Earthly configuration file

Global configuration values for earth can be stored on disk in the configuration file.

By default, earth reads the configuration file `~/.earthly/config.yaml`; however, it can also be
overridden with the `--config` command flag option.

Currently only https authentication configuration is configured in the config file; ssh authentication
is additionally supported (and the default when the no overrides have been configured)

## Format

The earthly config file is a [yaml](https://yaml.org/) formatted file that looks like:

```yaml
git:
    <site>:
        auth: https|ssh
        user: <username>
        password: <password>
    <site2>:
        ...
```

Example:

```yaml
git:
    github.com:
        auth: https
        user: alice
        password: itsasecret
```

## Git configuration reference

### site

The git repo hostname; e.g. `github.com`, or `gitlab.com`

### auth

Either `https` or `ssh` (default). If https is specified, user and password fields are used
to authenticate over https when pulling from git for the corresponding site.

See the [Authentication guide](guides/auth.md) for a guide on setting up authentication.

### user

The https username to use when auth is set to https; this is ignored when auth is ssh.

### password

The https password to use when auth is set to https; this is ignored when auth is ssh.

