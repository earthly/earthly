# Building An Image

## Introduction

This guide is intended to help you create your own Earthly-enabled image for your containerized CI workflows.

## Getting Started

There are two ways to build a containerized CI image with Earthly:

- Extending `earthly/earthly` with your runner/agent
- Adding Earthly to an existing image

This guide will cover both approaches to constructing your image. 

### Extending `earthly/earthly`

This is the recommended approach when adopting Earthly into your containerized CI. Start by basing your custom image on ours:

```docker
FROM earthly/earthly:v0.5.18
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
RUN wget https://github.com/earthly/earthly/releases/download/v0.5.18/earthly-linux-amd64 -O /usr/local/bin/earthly && \
    chmod +x /usr/local/bin/earthly && \
    /usr/local/bin/earthly bootstrap
```

As with the Docker containers, be sure to pin the version in the download URL to avoid any accidental future breakage. Assuming Docker is also installed and available, you should be able to invoke Earthly without any additional configuration.

#### Remote Daemon

When connecting to a remote daemon, follow the Docker-In-Docker installation instructions above to get the binary. Then you'll need to issue a few `earthly config` commands to ensure the container is set up to automatically use the remote daemon. It might look something like this:

```docker
RUN earthly config global "{buildkit_host: 'tcp://myhost:8372', buildkit_transport: 'tcp'}"
```

For more details on using a remote buildkit daemon, see our guide here.