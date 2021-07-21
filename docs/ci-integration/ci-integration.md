# Earthly in CI

Continuous Integration systems are as varied as the companies that use them. Fortunately, Earthly is flexible enough to fit into most (and where we don't, let us know!). This document serves as a startping point to configuring Earthly in your CI environment.

## Getting Started

### Dependencies

Earthly has two software dependencies: `docker` and `git`. 

 `docker` is used to glean information about the containerization environment, and manage our `earthly-buildkitd` daemon. It is also used to do things like save images locally on your machine after they have been built by Earthly.

Currently, the `earthly-buildkitd` daemon requires running in `--privileged` mode, which means that the `docker` daemon needs to be configured to allow this as well. Rootless configurations are currently unsupported.

`git` is used to help fetch remote targets, and also provides metadata for Earthly during your build.

Because `earthly` will not install these for you, please ensure they are present before proceeding. These tools are very common, so many environments will already have them installed. If they are not, here are some installation instructions that may help:

#### Linux

To install `git`, you can typically use your distributions package manager. [This page](https://git-scm.com/download/linux) has installation instructions for most distributions.

To install `docker`, use the most recent versions [directly from Docker](https://docs.docker.com/engine/install/#server). The versions packaged for many distributions tend to fall behind.

#### macOS

To install `git`, the easiest way is to install the "XCode Command Line Tools". If you open up `Terminal`, and type:

```go
git --version
```

Then macOS will prompt you to install these tools. You can also use the `git` provided installer or Homebrew, if you prefer. [Details can be found here](https://git-scm.com/download/mac).

To install `docker`, [download and install Docker CE](https://hub.docker.com/editions/community/docker-ce-desktop-mac). Be sure to grab the correct installer depending on your CPU architecture.

#### Windows

To install `git`, use the  [MSI installer](https://gitforwindows.org/). This will provide `git`, and a Bash shell; which may prove more natural for using Earthly. You may also use your package manager of choice.

To install `docker`, [download and install Docker CE](https://hub.docker.com/editions/community/docker-ce-desktop-windows). Both the HyperV and WSL2 backends are supported, but the WSL2 one is very likely to be faster.

### Installation

Once you have ensured that the dependencies are available, you'll need to install `earthly` itself.

#### Bare Metal

This is the simplest method for adding `earthly` to your CI. It will work best on dedicated computers, or in scripted/auto-provisioned build environments. You can follow our [regular installation guide](https://earthly.dev/get-earthly) to add `earthly` if desired.

If you are going to mostly be working from a WSL2 prompt in Windows, you might want to consider following the Linux instructions for installation. This will help prevent any cross-subsystem file transfers and keep your builds fast. Note that the "original" WSL is unsupported.

If your build host is persistent, it is recommended to install `earthly` as part of the new host's configuration, and not as part of your build. This will speed up your builds, since you do not need to download `earthly` each time; and it will also provide stability in case a future version of `earthly` changes the behavior of a command.

#### Containers

Earthly currently offers two official images:

- [`earthly/earthly`](https://hub.docker.com/r/earthly/earthly), which is a 1-stop shop. It includes a built-in `earthly-buildkitd` daemon, and accepts a target to be built as a parameter. It requires a mount for your source code, and an accessible `DOCKER_HOST`. When building a runner image for your CI; it is usually easier to start from this image, and add the pieces you need. See [this Jenkins agent configuration](https://github.com/earthly/ci-examples/blob/ce20840cffd2a8b04a8bd5dce477751adac3f490/jenkins/Earthfile#L48-L54) for an example.
- [`earthly/buildkitd`](https://hub.docker.com/r/earthly/buildkitd), which is the same `earthly-buildkitd` container that `earthly` will run on your host. This is useful in more advanced configurations, such as sharing a single `buildkitd` machine across many workers, or isolating the privileged parts of builds. See [this Kubernetes configuration](https://github.com/earthly/ci-examples/blob/main/kubernetes/buildkit.yaml) for an example.

If you need to provide additional configuration, [consider building your own image for CI](building-an-image.md).

### Configuration

While `earthly` is fairly configurable by itself, it also depends on the configuration of its dependencies. In a CI environment, you will need to ensure all of them are configured correctly.

#### Git

If you plan to build any private, or otherwise secure repositories, `git` will need to be configured to have access to these repositories. Please see the [`git` documentation for how to configure access](https://git-scm.com/docs/gitcredentials).

#### Docker

Like `git`, `docker` also needs to be configured to have access to any private repositories referenced in the `Earthfiles` you want to build. Please see [`docker`'s documentation for how to log in](https://docs.docker.com/engine/reference/commandline/login/), and our examples for pushing to many popular repositories.

If your private registry can use a [credential helper](https://docs.docker.com/engine/reference/commandline/login/#credential-helpers), configure it according to your vendor's instructions. `earthly` can also make use of these to provide access when needed. If you need help configuring `docker` for use with Earthly, see our [guides on configuring many popular registries](https://docs.earthly.dev/docs/guides/configuring-registries) for details.

#### Earthly

`earthly` has quite a few configuration options that can either be set through a configuration file or environment variables. See our [configuration reference](../earthly-config/earthly-config.md) for a complete list of options.

You can also configure `earthly` by using the [`earthly config` command](../earthly-command/earthly-command.md#earthly-config) from within a script. This can be useful for some dynamic configuration.

Some options that may make sense in a CI environment are:

|          Variable          |                                                                                                                 Description                                                                                                                      |
|----------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `CNI_MTU`                  | In some environments, the MTU externally may be different than the MTU on the internal CNI network, causing the internet to be unavailable. This lets you configure the internal network for when `earthly` auto-configures the MTU incorrectly. |
| `NO_COLOR` / `FORCE_COLOR` | Lets you force on/off the ANSI color codes. Use this when `earthly` misinterprets the presence of a terminal. Set either one to `1` to enable or disable colors.                                                                                 |
| `EARTHLY_BUILDKIT_HOST`    | Use this when you have an external buildkit instance you would like to use instead of the one `earthly` manages.                                                                                                                                 |

Earthly also has some special command-line switches to ensure best practices are followed within your CI. These come *highly* recommended. Enable these with the [`--ci`](../earthly-command/earthly-command.md#--ci-experimental) option,  which is shorthand for [`--use-inline-cache`](../earthly-command/earthly-command.md#--use-inline-cache-experimental) [`--save-inline-cache`](../earthly-command/earthly-command.md#--save-inline-cache-experimental) [`--strict`](../earthly-command/earthly-command.md#--strict) [`--no-output`](../earthly-command/earthly-command.md#--no-output).

Earthly also has a special [`--push`](../earthfile/earthfile.md#--push) option that can be used when invoking a target. In a CI, you may want to ensure this flag is present to push images or run commands that are not typically done as part of a normal development workflow.

If you would like to do cross-platform builds, you will need to install some [`binfmt_misc`](https://github.com/multiarch/qemu-user-static) entries. This can be done by running: `docker run --rm --privileged multiarch/qemu-user-static --reset -p yes`. This installs the needed entries and `qemue-user-static` binaries on your system. This will need to be repeated on each physical box (only once, since its a kernel level change, and the kernel is shared across containers).

To share secrets with `earthly`, use the [`--secret`](../earthfile/earthfile.md#--secret-env-varsecret-ref) option to inject secrets into your builds. You could also use our [cloud secrets](../guides/cloud-secrets.md), for a more seamless experience.

#### Networking & Security

Upon invocation, `earthly` depends on the availability of an `earthly-buildkit` daemon to perform its build. This daemon has some networking and security considerations.

Large builds can generate many `docker` pull requests for certain images. You can setup and use a [pull through cache](pull-through-cache.md) to circumvent this.

If `earthly` is running on a dedicated host, the only consideration to take is the ability to run the container in a `--privileged` mode. Typical installations *should* support this out of the box. We also support running under user namespaces, [when `earthly` is configured to start the `earthly-buildkit` container with the `--userns host` option](../earthly-config/earthly-config.md#buildkit_additional_args). Rootless configurations are currently unsupported.

If `earthly` is connecting to a remote `earthly-buildkitd`, then you will need to take additional steps. See this article for [running a remote buildkit instance](guides/remote-buildkit.md).

### Examples

Below are links to CI systems that we have produced more specific guides for. If you run into anything in your CI that wasn't covered here, we would love to add it to our documentation. Pull requests are welcome!

 * [Jenkins](guides/jenkins.md)
 * [Kubernetes](guides/kubernetes.md)
 * [Circle CI](guides/circle-integration.md)
 * [AWS CodeBuild](guides/codebuild-integration.md)
 * [GitHub Actions](guides/gh-actions-integration.md)
