# Earthly in CI

Continuous Integration systems are as varied as the companies that use them. Fortunately, Earthly is flexible enough to fit into most (and where we don't, let us know!). This document serves as a starting point to configuring Earthly in your CI environment.

Setting up Earthly is as easy as three steps:

 1. [Installing Dependencies](#dependencies)
 2. [Installing Earthly](#installation)
 3. [Configuration](#configuration)

We also have instructions for [specific CI systems](#examples); and special-case instructions for other scenarios (explore the "CI Integrations" category.)

## Dependencies

Earthly has two software dependencies: `docker` and `git`. Because `earthly` will not install these for you, please ensure they are present before proceeding. These tools are very common, so many environments will already have them installed. If you choose to use our prebuilt containers, these dependencies are already included.

`docker` is used to glean information about the containerization environment, and manage our `earthly-buildkitd` daemon. It is also used to do things like save images locally on your machine after they have been built by Earthly. To install `docker`, use the most recent versions [directly from Docker](https://docs.docker.com/engine/install/#server). The versions packaged for many distributions tend to fall behind.

`git` is used to help fetch remote targets, and also provides metadata for Earthly during your build. To install `git`, [you can typically use your distributions package manager](https://git-scm.com/download/linux).

## Installation

Once you have ensured that the dependencies are available, you'll need to install `earthly` itself.

### Option 1: Direct install

This is the simplest method for adding `earthly` to your CI. It will work best on dedicated computers, or in scripted/auto-provisioned build environments. You can pin it to a specific version like so:

```shell
wget https://github.com/earthly/earthly/releases/download/v0.8.8/earthly-linux-amd64 -O /usr/local/bin/earthly && \
chmod +x /usr/local/bin/earthly && \
/usr/local/bin/earthly bootstrap
```

It is recommended to install `earthly` as part of the new host's configuration, and not as part of your build. This will speed up your builds, since you do not need to download `earthly` each time; and it will also provide stability in case a future version of `earthly` changes the behavior of a command.

Don't forget to run `earthly bootstrap` when you are done to finish configuration!

### Option 2: Image

If a local installation isn't possible, Earthly currently offers two official images:

- [`earthly/earthly`](https://hub.docker.com/r/earthly/earthly), which is a 1-stop shop. It includes a built-in `earthly-buildkitd` daemon, and accepts a target to be built as a parameter. It requires a mount for your source code, and an accessible `DOCKER_HOST`.
- [`earthly/buildkitd`](https://hub.docker.com/r/earthly/buildkitd), which is the same `earthly-buildkitd` container that `earthly` will run on your host. This is useful in more advanced configurations, such as [remotely sharing](./remote-buildkit.md) a single `buildkitd` machine across many workers, or isolating the privileged parts of builds. This feature is experimental.

If you need to provide additional configuration or tools, [consider building your own image for CI](build-an-earthly-ci-image.md).

## Configuration

While `earthly` is fairly configurable by itself, it also depends on the configuration of its dependencies. In a CI environment, you will need to ensure all of them are configured correctly.

### Git

If you plan to build any private, or otherwise secure repositories, `git` will need to be configured to have access to these repositories. Please see our [documentation for how to configure access](../guides/auth.md#git-authentication).

### Docker

Like `git`, `docker` also needs to be configured to have access to any private repositories referenced in the `Earthfiles` you want to build. Please our [documentation for how to log in](../guides/auth.md#docker-authentication), and our examples for pushing to many popular repositories.

If your private registry can use a [credential helper](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers), configure it according to your vendor's instructions. `earthly` can also make use of these to provide access when needed. If you need help configuring `docker` for use with Earthly, see our [guides on configuring many popular registries](https://docs.earthly.dev/docs/guides/configuring-registries) for details.

Finally, the `earthly-buildkitd` daemon requires running in `--privileged` mode, which means that the `docker` daemon needs to be configured to allow this as well. Rootless configurations are currently unsupported.

### Earthly

`earthly` has quite a few configuration options that can either be set through a configuration file or environment variables. See our [configuration reference](../earthly-config/earthly-config.md) for a complete list of options.

You can also configure `earthly` by using the [`earthly config` command](../earthly-command/earthly-command.md#earthly-config) from within a script. This can be useful for some dynamic configuration.

Some options that may make sense in a CI environment are:

|          Variable          |                                                                                                                 Description                                                                                                                      |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `CNI_MTU`                  | In some environments, the MTU externally may be different than the MTU on the internal CNI network, causing the internet to be unavailable. This lets you configure the internal network for when `earthly` auto-configures the MTU incorrectly. |
| `NO_COLOR` / `FORCE_COLOR` | Lets you force on/off the ANSI color codes. Use this when `earthly` misinterprets the presence of a terminal. Set either one to `1` to enable or disable colors.                                                                                 |
| `EARTHLY_BUILDKIT_HOST`    | Use this when you have an external BuildKit instance you would like to use instead of the one `earthly` manages.                                                                                                                                 |

Earthly also has some special command-line switches to ensure best practices are followed within your CI. These come *highly* recommended. Enable these with the [`--ci`](../earthly-command/earthly-command.md#--ci) option, which is shorthand for [`--save-inline-cache`](../earthly-command/earthly-command.md#save-inline-cache) [`--strict`](../earthly-command/earthly-command.md#strict) [`--no-output`](../earthly-command/earthly-command.md#no-output).

Earthly also has a special [`--push`](../earthfile/earthfile.md#push) option that can be used when invoking a target. In a CI, you may want to ensure this flag is present to push images or run commands that are not typically done as part of a normal development workflow.

If you would like to do cross-platform builds, you will need to install some [`binfmt_misc`](https://github.com/multiarch/qemu-user-static) entries. This can be done by running: `docker run --rm --privileged multiarch/qemu-user-static --reset -p yes`. This installs the needed entries and `qemu-user-static` binaries on your system. This will need to be repeated on each physical box (only once, since its a kernel level change, and the kernel is shared across containers).

To share secrets with `earthly`, use the [`--secret`](../earthfile/earthfile.md#secret-less-than-env-var-greater-than-less-than-secret-ref-greater-than) option to inject secrets into your builds. You could also use our [cloud secrets](../cloud/cloud-secrets.md), for a more seamless experience.

### Networking & Security

Upon invocation, `earthly` depends on the availability of an `earthly-buildkit` daemon to perform its build. This daemon has some networking and security considerations.

Large builds can generate many `docker` pull requests for certain images. You can set up and use a [pull through cache](pull-through-cache.md) to circumvent this.

If `earthly` is running on a dedicated host, the only consideration to take is the ability to run the container in a `--privileged` mode. Typical installations *should* support this out of the box. We also support running under user namespaces, [when `earthly` is configured to start the `earthly-buildkit` container with the `--userns host` option](../earthly-config/earthly-config.md#buildkit_additional_args). Rootless configurations are currently unsupported.

If `earthly` is connecting to a remote `earthly-buildkitd`, then you will need to take additional steps. See this article for [running a remote BuildKit instance](remote-buildkit.md).

## Examples

Below are links to CI systems that we have more specific information for. If you run into anything in your CI that wasn't covered here, we would love to add it to our documentation. Pull requests are welcome!

 * [GitHub Actions](guides/gh-actions-integration.md)
 * [Circle CI](guides/circle-integration.md)
 * [GitLab CI/CD](guides/gitlab-integration.md)
 * [Jenkins](guides/jenkins.md)
 * [AWS CodeBuild](guides/codebuild-integration.md)
 * [Google Cloud Build](guides/google-cloud-build.md)
 * [Woodpecker CI](guides/woodpecker-integration.md)
 * [Kubernetes](guides/kubernetes.md)
