# Authenticating Git and image registries

This page guides you through passing Git and Docker authentication to Earthly builds, to empower related Earthly features, like `GIT CLONE` or `FROM`.

{% hint style='danger' %}
##### Important

This page is NOT about passing Git or Docker credentials for your own custom commands within builds. For those cases, use the [`RUN --secret`](../earthfile/earthfile.md#run) feature.
{% endhint %}

## Git authentication

A number of Earthly features use Git credentials to perform remote Git operations:

* Resolving a build context when referencing remote targets
* The `GIT CLONE` command

{% hint style='info' %}
##### Note

Currently, only `github.com` is supported as an SCM provider. If you need support for others, please [open a new GitHub issue](https://github.com/earthly/earthly/issues/new).
{% endhint %}

There are two possible ways to pass Git authentication to Earthly builds:

* Via SSH agent socket (for SSH-based authentication)
* Via username-password (usually for https Git URLs)

#### SSH agent socket

SSH agent socket passing is configured by default and it should just work. It uses the environment variable `SSH_AUTH_SOCK` to detect where the SSH agent socket is located and mounts that socket to the BuildKit daemon container. (As an exception, on Mac, Docker's compatibility SSH auth socket is used instead).

If you need to override the SSH agent socket, you can set the environment variable `EARTHLY_SSH_AUTH_SOCK` to point to an alternative SSH agent.

In order for the SSH agent to have the right credentials available, make sure you run `ssh-add` before executing Earthly builds.

Another key setting available, is `GIT_URL_INSTEAD_OF`. It allows for `https://github.com` URLs to be translated into `git@github.com` URLs, effectively forcing all github references to be interpreted through SSH-based authentication. This setting defaults to `git@github.com:=https://github.com/`, making SSH-based authentication work out of the box.

#### Username-password authentication

Username-password based authentication can be configured in the [earthly config file](../earth-config/earth-config.md) under the git section: 

```yaml
git:
    github.com:
        auth: https
        user: <username>
        password: <password>
    gitlab.com:
        auth: https
        user: <username>
        password: <password>

```

Alternatively, environment variables can be set which will be override all host entries from the config file:

* `GIT_USERNAME`
* `GIT_PASSWORD`

However, environment variable authentication are now deprecated in favor of using the configuration file instead.

## Docker authentication

Docker credentials are used in Earthly for inheriting from private images (via `FROM`) and for pushing images (via `SAVE IMAGE --push`).

Docker authentication works automatically out of the box. It uses the same Docker libraries to infer the location of the credentials on the system and optionally invoke any necessary credentials store helper to decrypt them.

All you have to do as a user is issue the command

```bash
docker login --username <username>
```

before issuing earth commands, if you have not already done so in the past.

## See also

* The [earth command reference](../earth-command/earth-command.md)
