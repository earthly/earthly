# Building An Earthly CI Image

## Introduction

This guide is intended to help you create your own Docker image with Earthly inside it for your containerized CI workflows.

## Getting Started

There are two ways to build a containerized CI image with Earthly:

- Extending the `earthly/earthly` image with an external runner/agent
- Adding Earthly to an existing image

This guide will cover both approaches to constructing your image. 

### Extending The `earthly/earthly` Image

This is the recommended approach when adopting Earthly into your containerized CI. Start by basing your custom image on ours:

```docker
FROM earthly/earthly:v0.8.6
RUN ... # Add your agent, certificates, tools...
```

When extending our image, be sure to pin to a specific version to avoid accidental future breakage as `earthly` evolves.

The `earthly/earthly` image is Alpine Linux based. To add tools to the image, you can use `apk`:

```docker
apk add --no-cache my-cool-tool
```

If you are adding a tool from outside the Alpine Linux repositories, test it to ensure it is compatible. Alpine uses `musl`, which *can* create incompatibilities with some software. 

Also, you should embed any configuration that your Earthly image might need (to avoid having it in your build scripts, or mounted from a host somewhere). You can do this in-line with the [`earthly config` command](../earthly-command/earthly-command.md#earthly-config).

### Adding Earthly To An Existing Image

This section will cover adding Earthly to an existing image when:

- Docker-In-Docker is configured for the base image
- Earthly will be connecting to a remote `earthly/buildkitd` instance

While it is possible to configure a locally-ran `earthly/buildkitd` instance within an image (it's how `earthly/earthly` works), the steps and tweaks are beyond the scope of this guide.

#### Docker-In-Docker

In this setup, Earthly will be allowed to manage an instance of its `earthly/buildkitd` daemon over a live Docker socket.

To enable this, simply follow the installation instructions within your Dockerfile/Earthfile as you would on any other host. An example of installing this can be found below.

```docker
RUN wget https://github.com/earthly/earthly/releases/download/v0.8.6/earthly-linux-amd64 -O /usr/local/bin/earthly && \
    chmod +x /usr/local/bin/earthly && \
    /usr/local/bin/earthly bootstrap
```

As with the Docker containers, be sure to pin the version in the download URL to avoid any accidental future breakage. Assuming Docker is also installed and available, you should be able to invoke Earthly without any additional configuration.

#### Remote Daemon

When connecting to a remote daemon, follow the Docker-In-Docker installation instructions above to get the binary. Then you'll need to issue a few `earthly config` commands to ensure the container is set up to automatically use the remote daemon. It might look something like this:

```docker
RUN earthly config global.buildkit_host buildkit_host: 'tcp://myhost:8372'
```

For more details on using a remote BuildKit daemon, [see our guide](./remote-buildkit.md).

## cgroups v2 Considerations

When cgroups v2 is detected by the `earthly/earthly` image's default entrypoint, it moves it's process under an isolated cgroup. If a different entrypoint is used (i.e. a custom user supplied script),
the root process must be moved into a separate cgroup, for example:

```bash
if [ -f "/sys/fs/cgroup/cgroup.controllers" ]; then
    echo "detected cgroups v2; moving pid $$ to subgroup"

    # move the process under a new cgroup to prevent buildkitd/entrypoint.sh
    # from getting a "sh: write error: Resource busy" error while enabling controllers
    # via echo +pids > /sys/fs/cgroup/cgroup.subtree_control
    mkdir -p /sys/fs/cgroup/my-entrypoint
    echo "$$" > /sys/fs/cgroup/my-entrypoint/cgroup.procs
fi
```

If this step is not performed before the buildkitd process starts up, buildkitd will be unable to initialize it's own cgroup (due to the container's root cgroup already having processes directly under it), and will
fail with the error: `sh: write error: Resource busy`.

## An important note about running the image

When running the built image in your CI of choice, if you're not using a remote daemon, Earthly will start BuildKit within the same container. In this case, it is important to ensure that the directory used by BuildKit to cache the builds is mounted as a Docker volume. Failing to do so may result in excessive disk usage, slow builds, or Earthly not functioning properly.

{% hint style='danger' %}
##### Important
We *strongly* recommend using a Docker volume for mounting `/tmp/earthly`. If you do not, BuildKit can consume excessive disk space, operate very slowly, or it might not function correctly.
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
