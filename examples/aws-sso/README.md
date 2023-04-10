# Earthly with AWS SSO

This directory contains an example of using AWS SSO and OIDC with Earthly and GitHub Actions

## Usage

First, edit `sso_config` with your AWS SSO details. Edit `config` so that it points to the AWS region you'd like to use.

When used on a developer's machine, Earthly will run `aws sso login`, open a web browser, and login with AWS SSO. Credentials will be cached both in an image layer, and at `cache`.

When used with the `--ci` flag, e.g. in GitHub Actions, Earthly will accept the credentials passed in with Earthly Secrets.

`earthly +target` will run `aws sts get-caller-identity` showing how to run arbitrary AWS commands using this Earthfile.

The included GitHub Actions workflow authenticates with OIDC and passes the credentials to Earthly.
