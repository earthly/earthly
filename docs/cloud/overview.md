# Earthly Cloud

Earthly Cloud is a collection of features that enrich the Earthly experience via cloud-based services. These include:

* [Earthly Satellites](./satellites.md): Cloud-based BuildKit instances managed by the Earthly team.
* [Earthly Cloud Secrets](./cloud-secrets.md): A secret management system that allows you to store secrets in a cloud-based service and use them across builds.
* [Auto-skip](../caching/caching-in-earthfiles.md#auto-skip): A feature that allows you to skip large parts of a build in certain situations.
* **Log sharing**: The ability to share build logs with coworkers.
* [OIDC Authentication](./oidc.md): The ability to authenticate to 3rd-party cloud services without storing long-term credentials.

## Sign up for Earthly Cloud for free!

*Get 6,000 build minutes/month on Satellites as part of Earthly Cloud's no time limit free tier.* ***[Sign up today](https://cloud.earthly.dev/login).***

## Getting started

### Creating an account

To get started with Earthly Cloud, you'll need to register an Earthly account. You can do so by visiting [Earthly Cloud Sign up page](https://cloud.earthly.dev/login), or by using the CLI as described below.

```bash
earthly account register --email <email>
```

An email will be sent to you containing a verification token. Next run:

```bash
earthly account register --email <email> --token <token>
```

This command will prompt you to set a password, and to optionally register a public-key for password-less authentication.

### Creating or joining an Earthly org

An Earthly org allows you to share projects, secrets, and satellites with colleagues. To view the orgs you belong to, run:

```bash
earthly org ls
```

To create an Earthly org you can run:

```bash
earthly org create <org-name>
```

To select the org you would like to use, run:

```bash
earthly org select <org-name>
```

To invite another user to join your org, run:

```bash
earthly org invite <email>
```

You can join an Earthly org by following the steps outlined in the invitation email sent to you by an Earthly admin.

### Creating a project

To use certain features, Earthly Cloud Secrets, you will additionally need to create an Earthly Project. You can create a project by using the CLI as described below.

```bash
earthly project create <project-name>
```

## Logging in from a CI

To be able to use certain Earthly features, such as Cloud Secrets, or Satellites from your CI, you will need to log into Earthly. The easiest way to do that is to create an Earthly authentication token by running

```bash
earthly account create-token [--write] <token-name>
```

This token can then be exported as an environment variable in the CI of choice.

```bash
EARTHLY_TOKEN=...
```

Which will then force Earthly to use that token when accessing secrets or satellites.
