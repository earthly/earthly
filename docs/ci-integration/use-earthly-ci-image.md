# Building An Earthly CI Image

## Introduction

This guide is intended to help you use the Earthly image for your containerized CI workflows.

## Prerequisites

The `earthly/earthly` image requires that it is run as `--privileged`, if it is meant to be used with the embedded BuiltKit daemon. This is mainly due to the use of overlayfs and the internal networking setup.

## Getting Started

Please see the reference documentation of the [`earthly/earthly` image on DockerHub](https://hub.docker.com/r/earthly/earthly).

It is recommended that the `earthly/earthly` image is used with a pinned version when used in the context of a CI, in order to avoid accidental future breakage as `earthly` evolves.

#### Using `/usr/bin/earthly-entrypoint.sh` as the entrypoint

The `earthly/earthly` image comes with an entrypoint that first starts up BuildKit and then issues an `earthly` command that makes use of it. You may use the image just as you would use `earthly` itself otherwise. Any arguments are passed into the `earthly` command directly.

{% hint style='danger' %}
##### Important
Note that using the `earthly` binary as the entrypoint will not start up BuildKit within the same container and will instead attempt to use the Docker Daemon (assuming one is available via `/var/run/docker.sock`) to start up BuildKit.
{% endhint %}

#### Remote Daemon

An alternative option is to use the `earthly/earthly` image in conjunction with a remote BuildKit Daemon. You may use the environment variable `BUILDKIT_HOST` to specify the hostname of the remote BuildKit Daemon. When this environment variable is set, the `earthly/earthly` image will not attempt to start BuildKit and will instead use the remote BuildKit Daemon.

For more details on using a remote BuildKit daemon, [see our guide](./remote-buildkit.md).

#### Mounting the source code

The image expects the source code of the application you are building in the current working directory (by default `/workspace`). You will need to copy or mount the necessary files to that directory prior to invoking the entrypoint.

```bash
docker run --privileged --rm -v "$PWD":/workspace earthly/earthly:v0.6.22 +my-target
```

Or, if you would like to use an alternative directory:

```bash
docker run --privileged --rm -v "$PWD":/my-dir -w /my-dir earthly/earthly:v0.6.22 +my-target
```

#### `NO_DOCKER` Environment Variable

In many CI use-cases outputting images locally is not necessary. In fact, Earthly's `--ci` flag disables output by default. In such circumstances, you can use the `NO_DOCKER` environment variable to disable checking for the presence of Docker. This will disable some warnings that would otherwise be printed to the console as Earthly starts up.

## An important note about running the image

When running the built image in your CI of choice, if you're not using a remote daemon, Earthly will start Buildkit within the same container. In this case, it is important to ensure that the directory used by Buildkit to cache the builds is mounted as a Docker volume. Failing to do so may result in excessive disk usage, slow builds, or Earthly not functioning properly.

{% hint style='danger' %}
##### Important
We *strongly* recommend using a Docker volume for mounting `/tmp/earthly`. If you do not, Buildkit can consume excessive disk space, operate very slowly, or it might not function correctly.
{% endhint %}

In some environments, not mounting `/tmp/earthly` as a Docker volume results in the following error:

```
--> WITH DOCKER RUN --privileged ...
...
rm: can't remove '/var/earthly/dind/...': Resource busy
```

In EKS, users reported that mounting an EBS volume, instead of a Kubernetes `emptyDir` worked.

This part of our documentation needs improvement. If you have a Kubernetes-based setup, please [let us know](https://earthly.dev/slack) how you have mounted `/tmp/earthly` and whether `WITH DOCKER` worked well for you.

For more information, see the [documentation for `earthly/earthly` on DockerHub](https://hub.docker.com/r/earthly/earthly).
