# Cloud Secrets

{% hint style='danger' %}
##### Important

This feature is currently in **Beta** stage

* If you encounter any issues, please give us feedback on [Slack](https://earthly.dev/slack) in the `#cloud-secrets` channel.
{% endhint %}

Earthly has the ability to use secure cloud-based storage for build secrets. This page goes through the basic setup and usage examples.

Cloud secrets can be used to share secrets between team members or across multiple computers and a CI systems.

## Introduction

This document covers the use of cloud-hosted secrets. It builds upon the understanding of [build arguments and locally-supplied secrets](../guides/build-args.md).

## Managing secrets

In order to be able to use cloud secrets, you need to first register an Earthly cloud account and create an Earthly org. Follow the steps in the [Earthly Cloud overview](../overview.md#getting-started) to get started.

Then, you need create an Earthly project. To do that, you may use the command

```bash
earthly project --org <org-name> create <project-name>
```

Or alternatively, launch the Earthly web interface by running `earthly web`, and click on **New Project**.

### Listing secrets

Each Earthly project has its own isolated secret store. To view the secrets within a given project, you can run

```bash
earthly secret --org <org-name> --project <project-name> ls
```

### Setting a value

To set a secret value, use the `secret set` command:

```bash
earthly secret --org <org-name> --project <project-name> set my_key 'hello world'
```

### Getting a value

To view a secret value, use the `secret get` command:

```bash
earthly secret --org <org-name> --project <project-name> ls
earthly secret --org <org-name> --project <project-name> get my_key
```

### User secrets

If a secret key starts with `/user/`, then the secret is stored in a special location accessible only by the current user. These secrets can never be shared. This may be useful, in cases where builds require that each developer uses their own set of credentials to access certain resources during builds.

When using the `/user/` prefix, the org name and project name are no longer required. These secrets may otherwise be accessed the same way.

```bash
earthly secret ls /user
earthly secret set /user/my_private_key 'hello private world'
earthly secret ls /user
earthly secret get /user/my_private_key
```

## Using cloud secrets in builds

When secrets need to be referenced in an Earthfile, you need to declare the project the secrets belong to at the top of the Earthfile, after the `VERSION` declaration.

```Dockerfile
VERSION 0.7
PROJECT <org-name>/<project-name>
```

Then, cloud secrets can be referenced in a similar way to [locally-defined secrets](build-args.md).

For example:

```Dockerfile
RUN --secret MY_KEY=my_key echo $MY_KEY
```

The env variable `MY_KEY` will be set with the value stored under the secret key `my_key`.

Or, to reference a user secret:

```Dockerfile
RUN --secret MY_KEY=/user/my_private_key echo $MY_KEY
```

## Security Details

The Earthly command uses HTTPS to communicate with the cloud secrets server. The server encrypts all secrets using OpenPGP's implementation of AES256 before storing it in a database. We use industry-standard security practices for managing our encryption keys in the cloud. For more information see our [Security page](https://earthly.dev/security).

## Migrating from the old 0.6 experimental version of Earthly secrets

The 0.6 version of Earthly cloud secrets is no longer supported. To help migrate to the new version, a migration command has been made available.

In the 0.6 version, the secrets were stored globally as part of an Earthly organization. In the new version, secrets are stored per project. To migrate, you need to first create a new project using `earthly` 0.7+ (if you haven't already), and then run the migration command.

You will first need to have read access to the source organization, and write access to the project you will be migrating to. The source organization can be the same as the destination one.

```bash
earthly project --org <org-name> create <project-name>
earthly secret --org <org-name> --project <project-name> migrate <source-org-name>
```

Once migration is complete, you can view the secrets in the new project using `earthly secret ls`.

```bash
earthly secret --org <org-name> --project <project-name> ls
```

To update your Earthfile to use the new secrets, you need to add the `PROJECT` declaration at the top of any Earthfile that needs secret access.

```Dockerfile
VERSION 0.7
PROJECT <org-name>/<project-name>
```

Secret references then need to be changed from `RUN --secret <env-var>=+secrets/<org-name>/<secret-key>` to `RUN --secret <env-var>=<secret-key>`. So the prefix `+secrets/<org-name>/` needs to be removed.
