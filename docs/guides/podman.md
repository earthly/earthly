# Podman
[Podman](https://podman.io/) is an alternative to docker; 
it's a daemonless container engine for developing, managing and running OCI containers on a Linux system.
Podman also works on Mac using a [podman machine](https://docs.podman.io/en/latest/markdown/podman-machine.1.html).

## Prerequisites
 - [Install podman](https://podman.io/getting-started/installation)
 - Mac users should ensure a [podman machine](podman machine) is running.
 - For [multi-platform builds](https://docs.earthly.dev/docs/guides/multi-platform) on Linux, install [qemu-user-static](https://github.com/multiarch/qemu-user-static).
 - Usage of the [WITH DOCKER](https://docs.earthly.dev/docs/earthfile#with-docker) command requires rootful mode.
   - Linux: run with `sudo` (i.e., `sudo earthly -P +with-docker-target`)
   - Mac: run a [rootful machine](https://docs.podman.io/en/latest/markdown/podman-machine-set.1.html#rootful).

## Getting Started
When earthly starts it performs a check to determine what frontend is available.
By default, earthly will attempt to use docker and then fall back to podman.
If you wish to change the behavior of the startup check, run the following command:

```bash
# Configure earthly to use podman
> earthly config global.container_frontend podman-shell

# Configure earthly to use docker
> earthly config global.container_frontend docker-shell
```

You can verify the command worked by checking the `~/.earthly/config.yml` file and verifying it contains a `container_frontend` entry.
```bash
> cat ~/.earthly/config.yml
global:
container_frontend: podman-shell
```

Then, you can run a basic hello world example to see earthly using the appropriate container frontend.
```bash
> earthly github.com/earthly/hello-world:main+hello
 1. Init 🚀
————————————————————————————————————————————————————————————————————————————————

           buildkitd | Starting buildkit daemon as a **podman** container (earthly-buildkitd)...
           buildkitd | ...Done
```

If instead you see 