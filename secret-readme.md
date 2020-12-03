# Earthly Secret Store

The secret store can be used to share secrets between team members or across multiple computers and CI systems.

This service is currently experimental.

# Usage

Please upgrade to the latest version, at the time of writing this guide, that version is v0.3.19

## Registering for an account

First, you must register for an account:

    earth account register --email <email>

An email will be sent to you containing a verification token, next run:

    earth account register --email <email> --token <token>

This command will prompt you to set a password, and to optionally register a public-key for password-less authentication.

## Login / Logout

It is recommended that you register a public rsa key during registration; if this is done, you will be logged in automatically whenever
earth needs to authenticate you. If you did not supply a public key, then your plain-text password will be cached on your local disk under
`~/.earthly/auth-token`, which will be used to log you in. If this file is deleted, you will need to run `earth account login` to re-create it.

To logout, you can run `earth account logout`, which deletes the `~/.earthly/auth-token` file from your disk.

## Interacting with the private user secret store from the command line

Each user will have a non-sharable private userspace which can be referenced by `/user/...`; this can be thought of as your home directory.
To view this workspace, try running:

    earth secrets ls

Secrets are referenced by a path, and can contain up to 512 bytes.

### Setting a value

To set a secret value, use the `secrets set` command:

    earth secrets set /user/my_key 'hello wolrd'

### Getting a value

To view a secret value, use the `secrets get` command:

    earth secrets get /user/my_key

### Listing secrets

To view secrets, use the `secrets ls` command:

    earth secrets ls /user

## Referencing secrets from an Earthfile

Secrets can also be included under an Earthfile:

    FROM alpine:latest
    
    build:
        RUN --secret MY_KEY=+secrets/user/my_key echo $MY_KEY > /secret-data
        SAVE IMAGE my_image_with_a_secret:latest

The env variable `MY_KEY` will be set with the value stored under your private `/user/my_key` secret.

You can build it via:

    earth +build

Then run a command like `cat` to show that the data was correctly added in the build step:

    docker run --rm my_image_with_a_secret:latest cat /secret-data


# Sharing secrets between teams

To share secrets between teams, an organization must first be created:

    earth org create <org-name>

Then additional users can be invited into the organization:

    earth org invite /<org-name>/ <email>

By default this will grant the invited user read privileges to all keys under the organization. It's also possible to
use the `--write` flag to grant write permission too. Additionally, the permissions can be set to lower paths.

## Example

Alice and Bob sign up for earthly accounts using alice@earthly.dev and bob@earthly.dev respectively:

    earth account register --email alice@earthly.dev --token ...
    earth account register --email bob@earthly.dev --token ...

Alice then creates an organization called hush-co:

    earth org create hush-co

Alice then creates a secret under the project-zulu sub directory:

    earth secrets set /hush-co/project-zulu/transponder-code peanut

Alice then grants Bob read permission on all of project-zulu:

    earth org invite /hush-co/project-zulu/ bob@earthly.dev

Bob now has permission to everything under the `/hush-co/project-zulu/` directory; if he runs

    earth secrets ls /hush-co/

he will see:

    /hush-co/project-zulu/transponder-code

However if Alice were to create any secrets outside of project-zulu, Bob would not be able to list or retrieve them.

# Integrating with CI

To reference secrets from a CI environment, one can make use of the password or ssh-key authentication
referenced under the login/logout section, or you can generate an authentication token by running:

    earth account create-token [--write] <token-name>

This token can then be exported as

    EARTHLY_TOKEN=...

Which will then force earthly to use that token when accessing secrets. This is useful for cases
where running an ssh-agent is impractical.

# Feedback

The secrets store is still an experimental feature, we would love to hear feedback in our 
[Slack](https://join.slack.com/t/earthlycommunity/shared_invite/zt-ix9rtuv8-DUFl8uxe5bFULxyCGGbqJQ) community.

Happy Coding.
