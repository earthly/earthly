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

There are two possible ways to pass Git authentication to Earthly builds:

* Via SSH agent socket (for SSH-based authentication)
* Via username-password (usually for HTTPS Git URLs)

#### Auto authentication

Earthly defaults to an `auto` authentication mode, where ssh-based authentication is automatically attempted, and falls back to https-based cloning.

{% hint style='info' %}
If you are having trouble accessing a private repository and want to use ssh-based authentication, first make sure `ssh-agent` is running and the `SSH_AUTH_SOCK`
environment variable is set. If not, you can start it with `eval $(ssh-agent)`.

Next make sure your private key has been added by running `ssh-add <path to key>`.
{% endhint %}

For users who want explicit control over git authentication, the following sections explain how.

#### SSH agent socket

Earthly uses the environment variable `SSH_AUTH_SOCK` to detect where the SSH agent socket is located and mounts that socket to the BuildKit daemon container.
(As an exception, on Mac, Docker's compatibility SSH auth socket is used instead).

If you need to override the SSH agent socket, you can set the environment variable `EARTHLY_SSH_AUTH_SOCK`, or use the `--ssh-auth-sock` flag to point to an alternative SSH agent.

In order for the SSH agent to have the right credentials available, make sure you run `ssh-add` before executing Earthly builds.

Another key setting is the `auth` mode for the git site that hosts the repository. By default earthly automatically default to `ssh` authentication if the ssh auth agent is running and has at least 1 key loaded, otherwise `earthly` will fallback to using non-authenticated HTTPS.

Sites can be explicitly added to the [earthly config file](../earthly-config/earthly-config.md) under the git section in order to override the auto-authentication mode:

```yaml
git:
    git.example.com:
        auth: ssh
        user: git
```

#### Username-password authentication

Username-password based authentication can be configured in the [earthly config file](../earthly-config/earthly-config.md) under the git section:

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

If no `user` or `password` are found, earthly will check for entries under [`~/.netrc`](https://everything.curl.dev/usingcurl/netrc).

##### Global override via environment variables (deprecated)

Alternatively, environment variables can be set which will be override all host entries from the config file:

* `GIT_USERNAME`
* `GIT_PASSWORD`

However, environment variable authentication are now deprecated in favor of using the configuration file instead.

#### Self-hosted and private Git Repositories

Currently, `github.com`, `gitlab.com`, and `bitbucket.org` have been tested as SCM providers. In order to use self-hosted git repositories
such as GitHub enterprise you will need to configure Earthly to use them using a regular expression:

```yaml
git:
    ghe.internal.mycompany.com:
        pattern: 'ghe.internal.mycompany.com/([^/]+)/([^/]+)'
        substitute: 'ssh://git@ghe.internal.mycompany.com:22/$1/$2.git'
        auth: ssh
```

## Docker authentication

Docker credentials are used in Earthly for inheriting from private images (via `FROM`) and for pushing images (via `SAVE IMAGE --push`).

Docker authentication works automatically out of the box. It uses the same Docker libraries to infer the location of the credentials on the system and optionally invoke any necessary credentials store helper to decrypt them.

### Manually

All you have to do as a user is issue the command

```bash
docker login --username <username>
```

before issuing earthly commands, if you have not already done so in the past. If you run into troubles, [you can find out more about `docker login` here](https://docs.docker.com/engine/reference/commandline/login/).

### Credential Helpers

Docker can use various credential helpers to automatically generate and use credentials on your behalf. These are usually created by cloud providers to allow Docker to authenticate using the cloud providers own credentials.

You can see examples of configuring Docker to use these, and working with Earthly here:
* [Pushing and Pulling Images with AWS ECR](./registries/aws-ecr.md)
* [Pushing and Pulling Images with GCP Artifact Registry](./registries/gcp-artifact-registry.md)
* [Pushing and Pulling Images with Azure ACR](./registries/azure-acr.md)

## See also

* The [earthly command reference](../earthly-command/earthly-command.md)
