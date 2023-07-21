# Earthly CI

{% hint style='danger' %}
##### Important

This feature is currently in **Beta** stage

* The feature may break or change significantly in future versions of Earthly.
* Give us feedback on
  * [Slack](https://earthly.dev/slack)
  * [GitHub issues](https://github.com/earthly/cloud-issues/issues)
  * [Emailing support](mailto:support+ci@earthly.dev)
{% endhint %}

Earthly CI is a hosted CI service that allows you to run your Earthly builds in the cloud. Earthly CI gives teams repeatable pipelines that run exactly the same in CI as on your laptop; has an automatic and instantly available build cache that makes builds faster; and is super simple to use.

Earthly CI uses Earthfiles as the build configuration language.

## Get started with Earthly CI for free!
Earthly CI is included with [Earthly Cloud](https://docs.earthly.dev/earthly-cloud/overview). Earthly Cloud is a SaaS build automation platform with consistent builds, ridiculous speed, and a next-gen developer experience that works seamlessly with any CI. It includes the functionality of both Earthly CI and [Earthly Satellites](https://docs.earthly.dev/earthly-cloud/satellites). *Get 6,000 build minutes/month as part of Earthly Cloud's no time limit free tier.* ***[Sign up today](https://cloud.earthly.dev/login).***

## Benefits

* **Ridiculously fast** - Earthly CI uses the same build cache and build parallelization technology as Earthly Satellites, so builds are 2-20X faster compared to a traditional CI.
* **Super simple** - Earthfiles have a super simple, instantly recognizable syntax – like Dockerfile and Makefile had a baby.
* **Great for Monorepos and Polyrepos** - Earthly CI is great for both monorepos and polyrepos. You can organize your build logic however makes the most sense for your project. The caching ensures that only what has changed is rebuilt.
* **Remote build runners** - Earthly CI comes with access to Earthly Satellites. This means that you can run ad-hoc remote builds in CI from your laptop without having to commit code to Git with every attempt.

## Getting started

### 1. Sign up for Earthly Cloud (free)

Earthly CI is part of Earthly Cloud. You can use it for free as part of our free tier. Get started with Earthly Cloud by visiting the [sign up](https://cloud.earthly.dev/login) page, and get 6,000 build minutes/month for free.

### 2. Create your first Earthly pipeline

In the Earthly Web UI, under the Organizations drop-down in the top-left of the screen, select the organization that was granted access to Earthly CI. Then follow the instructions on the screen to create a new Earthly CI project, add your first repository, create your pipeline via a new or an existing Earthfile, and then run your first build.

For your very first CI pipeline you can use the following example:

```Earthfile
VERSION 0.7
PROJECT my-org/my-project

FROM alpine:3.15

my-pipeline:
  PIPELINE
  TRIGGER push main
  TRIGGER pr main
  BUILD +my-build

my-build:
  RUN echo Hello world
```

This example shows a simple pipeline called `my-pipeline`, which is triggered on either a push to the `main` branch, or a pull request against the `main` branch. The pipeline executes the target `my-build`, which simply prints `Hello world`.

Pipelines and their definitions, including their triggers must be merged into the primary branch (which, unless overridden, is the default branch on GitHub -- usually `main`) in order for the triggers to take effect.

Please note that the `PROJECT` declaration needs to match the name of your organization and the name of your Earthly project that you create in Earthly Cloud.

### 3. Set up registry access

During your experimentation with Earthly CI, you may encounter DockerHub rate limiting errors. To avoid this, you can setup your DockerHub account by using the command `earthly registry setup`.

```bash
earthly registry --org <org-name> --project <project-name> \
setup --username <registry-user-name> --password-stdin \
<host>
```

If the registry is DockerHub, then you can leave out the registry host argument.

You may additionally log into other registries, such as AWS ECR, or GCP Artifact Registry, by using the following:

```bash
# For AWS ECR
earthly registry --org <org-name> --project <project-name> setup --cred-helper ecr-login --aws-access-key-id <key> --aws-secret-access-key <secret> <host>
# For GCP Artifact Registry
earthly registry --org <org-name> --project <project-name> setup --cred-helper gcloud --gcp-service-account-key <key> <host>
```

### 4. Invite your team

A final step is to invite your team to use Earthly Cloud. This can be done by running:

```bash
earthly --org <org-name> org invite <email>
```

Or by using the Earthly Web UI.

## Known limitations

Please note that we are aware of the following ongoing issues:

* GitHub only for now. If you are using GitLab, using [Earthly Satellite](./satellites.md) on top of GitLab CI is a great alternative.

If you run into any issues please let us know either via [Slack](https://earthly.dev/slack), [GitHub issues](https://github.com/earthly/cloud-issues/issues) or by [emailing support](mailto:support+ci@earthly.dev).

## Next

Check out the following additinal resources:

* [Earthly Cloud overview](./overview.md)
* [PIPELINE command reference](../earthfile/earthfile.md#pipeline-beta)
* [TRIGGER command reference](../earthfile/earthfile.md#trigger-beta)
* [PROJECT command reference](../earthfile/earthfile.md#project)
* [Earthly CI secrets](./cloud-secrets.md)
* [Earthly registry command reference](../earthly-command/earthly-command.md#earthly-registry)
