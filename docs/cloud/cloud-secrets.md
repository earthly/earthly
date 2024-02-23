# Cloud Secrets

Earthly has the ability to use secure cloud-based storage for build secrets. This page goes through the basic setup and usage examples.

Cloud secrets can be used to share secrets between team members or across multiple computers and a CI systems.

## Introduction

This page covers the use of cloud-hosted secrets. It builds upon the understanding of [build arguments and locally-supplied secrets](../guides/build-args.md).

## Managing secrets

In order to be able to use cloud secrets, you need to first register an Earthly Cloud account. Visit [Earthly Cloud](https://cloud.earthly.dev) to sign up for free.

Then, you need create an Earthly project. To do that, you may use the web interface or the command...

```bash
earthly project create <project-name>
```

Access to secrets is controlled by the project they belong to. Anyone with at least `read+secrets` access level for the org or the specific project will be able to see and use the secrets in their builds. Anyone with `write` access level will be able to create, modify and delete secrets. For more information on managing permissions see the [Managing Permissions page](./managing-permissions.md).

### Listing secrets

Each Earthly project has its own isolated secret store. Multiple code repositories may be associated with a single Earthly project. To view the secrets within a given project, you can run

```bash
earthly secret --project <project-name> ls
```

### Setting a value

To set a secret value, use the `secret set` command:

```bash
earthly secret --project <project-name> set my_key 'hello world'
```

### Getting a value

To view a secret value, use the `secret get` command:

```bash
earthly secret --project <project-name> ls
earthly secret --project <project-name> get my_key
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
VERSION 0.8
PROJECT <org-name>/<project-name>
```

Then, cloud secrets can be referenced in a similar way to [locally-defined secrets](../guides/build-args.md).

For example:

```Dockerfile
RUN --secret MY_KEY=my_key echo $MY_KEY
```

The env variable `MY_KEY` will be set with the value stored under the secret key `my_key`.

Or, to reference a user secret:

```Dockerfile
RUN --secret MY_KEY=/user/my_private_key echo $MY_KEY
```

## Using Cloud Secrets to log into registries

Cloud secrets can also be used to log into container registries such as DockerHub, AWS ECR, or GCP Artifact Registry. Here is the command to use:

```bash
earthly registry --org <org-name> --project <project-name> \
setup --username <registry-user-name> --password-stdin \
<host>
```

This command stores the username and password within the cloud secrets store. Earthly picks these up automatically when pulling or pushing images. If the registry is DockerHub, then you can leave out the registry host argument.

You may additionally log into other registries, such as AWS ECR, or GCP Artifact Registry, by using the following:

```bash
# For AWS ECR
earthly registry --org <org-name> --project <project-name> setup --cred-helper ecr-login --aws-access-key-id <key> --aws-secret-access-key <secret> <host>
# For GCP Artifact Registry
earthly registry --org <org-name> --project <project-name> setup --cred-helper gcloud --gcp-service-account-key <key> <host>
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
earthly secret --project <project-name> ls
```

To update your Earthfile to use the new secrets, you need to add the `PROJECT` declaration at the top of any Earthfile that needs secret access.

```Dockerfile
VERSION 0.8
PROJECT <org-name>/<project-name>
```

Secret references then need to be changed from `RUN --secret <env-var>=+secrets/<org-name>/<secret-key>` to `RUN --secret <env-var>=<secret-key>`. So the prefix `+secrets/<org-name>/` needs to be removed. For user secrets, only the prefix `+secrets` needs to be removed, such that the key of the secret contains the `/user/` prefix.

Please note that user secrets are not migrated by the `migrate` command. You will need to manually re-create them.
