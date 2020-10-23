# CI integration

Integrating Earthly into your CI is simply a matter of automating the same steps you would use for your local installation. In this guide, we will walk through this process.

## Step 1: Ensure pre-requisites are available

### Docker and Git

The first step is to ensure that Earthly's pre-requisites, Docker and Git, are available. On many CI systems both of these already exist in the default base image or environment. Refer to your provider's documentation.

Vendors known to include these dependencies:

* CircleCI image `ubuntu-1604:201903-01`
* GitHub actions `ubuntu-latest`
* Travis dist `xenial`
* GitLab image `docker:git` with service `docker:dind` added.
* Azure DevOps vmImage `Ubuntu-16.04`
* AWS CodeBuild image `aws/codebuild/amazonlinux2-x86_64-standard:3.0`

### Privileged mode

In addition to Docker and Git, Earthly also requires privileged mode as it executes container builds under the hood. In most linux-based CI environments, this is readily available and no special setting is necessary. GitLab CI requires using a compatible runner (eg Docker) and explicitly enabling [privileged mode](https://docs.gitlab.com/runner/executors/docker.html#the-privileged-mode).

## Step 2: Install earth command

The next step is to install the `earth` command. For this, you need to run the command:

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.3.10/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth'
```

{% hint style='info' %}
##### Note

The above command installs a specific version of `earth`. It is not advisable to always download the very latest, as any possible backwards incompatible changes of Earthly (rare) could cause your builds to fail unexpectedly. Pinning to a specific version, as shown here, is recommended.
{% endhint %}

In certain CI environments, such as Jenkins, where you have access to the host, it may be more convenient to install Earthly on the host in advance, so that your builds do not need to download `earth` every time they run.

## Step 3: Configure earth

Depending on your needs, you may need to ensure that Git has authenticated access and / or that Docker is logged in so that it has access to private repositories.

To authenticate Git, you may either use SSH-based authentication, or username-password-based authentication. See the [Authentication page for more information](./auth.md). In case you don't need any Git authentication, you might want to force all GitHub URLs to be transformed to `https://github.com/...` instead of `git@github.com:...`. For this, you can add an environment variable to configure this behavior:

```bash
export GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:"
```

The way you configure environment variables in your CI will vary.

To log in Docker, simply run

```bash
docker login --username '<username>' --password '<password>'
```

{% hint style='info' %}
##### Note

Make sure that secrets (like `<password>` above) are not exposed in plain text. You may need to configure an environment variable with your CI vendor.
{% endhint %}

## Step 4: (Optional) Choose log formatting options

There are a few options you can configure to make the log formatting look best in your CI environment.

### Force or disable color output

The CLI `earth` automatically detects the presence of a TTY for the purpose of deciding whether to use colorized output or not. In some CI environments, this kind of detection is not enough in order to infer support for colorized output. However, two environment variables can be used to either disable or force it:

* `NO_COLOR=1` disables the use of color.
* `FORCE_COLOR=1` forces the use of color.

The following environments are known to require additional settings:

* GitHub Actions: requires `FORCE_COLOR=1`
* Jenkins: requires `NO_COLOR=1`

### Tweak target name column

The CLI defaults to a right-aligned, 20-character width to display the name of the target that the log originated from. When a target name exceeds 20 characters, the name is truncated and the last 20 characters are used.

 * `TARGET_PADDING=n` will set the column to the width of `n` characters.
 * `FULL_TARGET=1` will always print the full target name, and nothing more. This may result in unaligned logs.

## Step 5: Run the build

This is often as simple as

```bash
earth +target-name
```

If you would like to enable pushing Docker images to registries and also running `RUN --push` commands, you might use

```bash
earth --push +target-name
```

If you need to pass secrets to the Earthly build, you might also use the `--secret` flag, mentioning the env var where the secret is kept.

```bash
earth --secret SOME_SECRET_ENV_VAR +target-name
```

For more information see the [earth command reference](../earth-command/earth-command.md).

## Complete examples

A couple of build examples are available for

* [Circle CI](../examples/circle-integration.md)
* [GitHub Actions](../examples/gh-actions-integration.md)
* [AWS CodeBuild](../examples/codebuild-integration.md)
