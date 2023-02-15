# Earthly CI

This feature is part of the Earthly CI paid plan.

{% hint style='danger' %}
##### Important

This feature is currently in **Alpha** stage. **We do not recommend using Earthly CI Alpha for mission-critical pipelines.**

* The feature may break or change significantly in future versions of Earthly.
* Give us feedback on
  * [Slack](https://earthly.dev/slack)
  * [GitHub issues](https://github.com/earthly/cloud-issues/issues)
  * [Emailing support](mailto:support+ci@earthly.dev)
{% endhint %}

Earthly CI is a hosted CI service that allows you to run your Earthly builds in the cloud. Earthly CI gives teams repeatable pipelines that run exactly the same in CI as on your laptop; has an automatic and instantly available build cache that makes builds faster; and is super simple to use.

Earthly CI uses Earthfiles as the build configuration language.

## Benefits

* **Ridiculously fast** - Earthly CI uses the same build cache and build parallelization technology as Earthly Satellites, so builds are 2-20X faster compared to a traditional CI.
* **Super simple** - Earthfiles have a super simple, instantly recognizable syntax – like Dockerfile and Makefile had a baby.
* **Great for Monorepos and Polyrepos** - Earthly CI is great for both monorepos and polyrepos. You can organize your build logic however makes the most sense for your project. The caching ensures that only what has changed is rebuilt.
* **Remote build runners** - Earthly CI comes with access to Earthly Satellites. This means that you can run builds from your laptop without having to commit code to Git with every attempt.

If you are upgrading from Earthly Satellites, the main benefit of using Earthly CI is that you no longer need to use a traditional CI in combination. This means less moving parts, simpler setup, slightly faster builds (no need to download Earthly during the build), and less bills to pay!

## Getting started

### 1. Register an account and create an org

Follow the steps in the [Earthly Cloud overview](./overview.md#getting-started) to register an account and create an org.

### 2. Gain access to Earthly CI

Fill out the [Earthly CI signup form](https://forms.gle/iwSLazNv7FjYqyqY6). You will need to wait for an Earthly team member to approve your request.

<!-- Use this after alpha.
Fill out the [Earthly CI signup form](https://earthly.dev/signup/earthly-ci/). You will need to wait for an Earthly team member to approve your request.
-->

### 3. Ensure that you have the latest version of Earthly

You need to make sure that you are using the [latest Earthly 0.7 release candidate](https://github.com/earthly/earthly/releases). Since this is a release candidate, it will not show up in the regular update channels, such as brew. You have to download the binary manually and add it to your `PATH`. Please note that there are some key changes from 0.6 to 0.7 - one of the important ones is that the cloud-based secret storage has a new format. See the [secrets migration guide](./cloud-secrets.md#migrating-from-the-old-0.6-experimental-version-of-earthly-secrets) for more information.

<!-- Use this after alpha.
Because this feature is under heavy development right now, it is very important that you use the latest version of Earthly available.

**On Linux**, simply repeat the [installation steps](https://earthly.dev/get-earthly) to upgrade.

**On Mac**, you can perform:

```bash
brew update
brew upgrade earthly/earthly/earthly
```
-->

### 4. Open the Earthly Web UI

You can open the Earthly Web UI by running:

```bash
earthly web
```

This will associate your Earthly account with your GitHub login and then take you to the Earthly Web UI.

### 5. Create your first Earthly pipeline

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

Please note that the `PROJECT` declaration needs to match the name of your organization and the name of your Earthly project that you create in the Earthly Web UI.

### 6. Set up registry access

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

### 7. Invite your team

A final optional step is to invite your team to use Earthly CI. This can be done by running:

```bash
earthly --org <org-name> org invite <email>
```

Or by using the Earthly Web UI.

## Known limitations

Please note that we are aware of the following ongoing issues:

* The logs of certain builds do not show up in the UI sometimes. If you run into this problem, please let us know and provide us with the URL of the page!
* Creating an account, or logging in directly in the web UI is not yet available. Use the `earthly web` command instead.
* GitHub only for now.

If you run into any issues please let us know either via [Slack](https://earthly.dev/slack), [GitHub issues](https://github.com/earthly/cloud-issues/issues) or by [emailing support](mailto:support+ci@earthly.dev).

## Next

Check out the following additinal resources:

* [Earthly Cloud overview](./overview.md)
* [PIPELINE command reference](../earthfile/earthfile.md#pipeline-beta)
* [TRIGGER command reference](../earthfile/earthfile.md#trigger-beta)
* [PROJECT command reference](../earthfile/earthfile.md#project)
* [Earthly CI secrets](./cloud-secrets.md)
* [Earthly registry command reference](../earthly-command/earthly-command.md#earthly-registry)
