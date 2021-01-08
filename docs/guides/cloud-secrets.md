# Cloud Secrets

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/575) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#cloud-secrets` channel.
{% endhint %}

Earthly has the ability to use secure cloud-based storage for build secrets. This page goes through the basic setup and usage examples.

Cloud secrets can be used to share secrets between team members or across multiple computers and a CI systems.

## Managing secrets

In order to be able to use cloud secrets, you need to first register an Earthly cloud account.

```bash
earthly account register --email <email>
```

An email will be sent to you containing a verification token, next run:

```bash
earthly account register --email <email> --token <token>
```

This command will prompt you to set a password, and to optionally register a public-key for password-less authentication.

### Login / Logout

It is recommended that you register a public RSA key during registration; if this is done, you will be logged in automatically whenever
earthly needs to authenticate you. If you did not supply a public key, then your plain-text password will be cached on your local disk under
`~/.earthly/auth-token`, which will be used to log you in. If this file is deleted, you will need to run `earthly account login` to re-create it.

To logout, you can run `earthly account logout`, which deletes the `~/.earthly/auth-token` file from your disk.

### Interacting with the private user secret store from the command line

Each user has a non-sharable private userspace which can be referenced by `/user/...`; this can be thought of as your home directory.
To view this workspace, try running:

```bash
earthly secrets ls
earthly secrets ls /user
```

Secrets are referenced by a path, and can contain up to 512 bytes.

### Setting a value

To set a secret value, use the `secrets set` command:

```bash
earthly secrets set /user/my_key 'hello world'
```

### Getting a value

To view a secret value, use the `secrets get` command:

```bash
earthly secrets ls /user
earthly secrets get /user/my_key
```

## Using cloud secrets in builds

Secrets can also be referenced in an Earthfile:

```Dockerfile
FROM alpine:latest

build:
    RUN --secret MY_KEY=+secrets/user/my_key echo $MY_KEY
    SAVE IMAGE myimage:latest
```

The env variable `MY_KEY` will be set with the value stored under your private `/user/my_key` secret.

You can build it via:

```bash
earthly +build
```

## Sharing secrets

To share secrets between teams, an organization must first be created:

```bash
earthly org create <org-name>
```

Then additional users can be invited into the organization:

```bash
earthly org invite /<org-name>/ <email>
```

By default this will grant the invited user read privileges to all keys under the organization. It's also possible to
use the `--write` flag to grant write permission too. Additionally, the permissions can be set to lower paths.

### Sharing example

Alice and Bob sign up for earthly accounts using alice@example.com and bob@example.com respectively:

```bash
earthly account register --email alice@example.com --token ...
earthly account register --email bob@example.com --token ...
```

Alice then creates an organization called hush-co:

```bash
earthly org create hush-co
```

Alice then creates a secret under the `project-zulu` sub directory:

```bash
earthly secrets set /hush-co/project-zulu/transponder-code peanut
```

Alice then grants Bob read permission on all of `project-zulu`:

```bash
earthly org invite /hush-co/project-zulu/ bob@example.com
```

Bob now has permission to everything under the `/hush-co/project-zulu/` directory. If he runs

```bash
earthly secrets ls /hush-co/
```

he will see:

```
/hush-co/project-zulu/transponder-code
```

However if Alice were to create any secrets outside of `project-zulu`, Bob would not be able to list or retrieve them.

## Using cloud secrets in CI

To reference secrets from a CI environment, you can make use of the password or ssh-key authentication referenced under the login/logout section, or you can generate an authentication token by running:

```bash
earthly account create-token [--write] <token-name>
```

This token can then be exported as

```
EARTHLY_TOKEN=...
```

Which will then force Earthly to use that token when accessing secrets. This is useful for cases
where running an ssh-agent is impractical.

# Feedback

The secrets store is still an experimental feature, we would love to hear feedback in our 
[Slack](https://earthly.dev/slack) community.
