# Earthly Cloud

Earthly Cloud is a collection of features that enrich the Earthly experience via cloud-based services. These include:

* [Earthly Cloud Secrets](./cloud-secrets.md): A secret management system that allows you to store secrets in a cloud-based service and use them across builds.
* [Earthly Satellites](./satellites.md): Cloud-based Buildkit instances managed by the Earthly team.
* **Earthly CI** (coming soon): A cloud-based continuous integration / continuous delivery system that allows you to continuously build your code in the cloud.

## Getting started

### Creating an account

To get started with Earthly Cloud, you'll need to register an Earthly account. You can do so by visiting [Earthly CI](https://ci.earthly.dev), or by using the CLI as described below.

```bash
earthly account register --email <email>
```

An email will be sent to you containing a verification token. Next run:

```bash
earthly account register --email <email> --token <token>
```

This command will prompt you to set a password, and to optionally register a public-key for password-less authentication.

### Creating or joining an Earthly org

An Earthly org allows you to share projects, secrets, satellites and pipelines with colleagues. To create an Earthly org you can run:

```bash
earthly org create <org-name>
```

To invite another user to join your org, run:

```bash
earthly org invite /<org-name>/ <email>
```

Note the slashes around the org name. Also, please note that the user must have an account on Earthly before they can be invited. (This is a temporary limitation which will be addressed in the future.)

You can join an Earthly org by following the steps outlined in the invitation email sent to you by an Earthly admin.
