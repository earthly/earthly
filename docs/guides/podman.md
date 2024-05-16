# Podman
[Podman](https://podman.io/) is an alternative to docker; 
it's a daemonless container engine for developing, managing and running OCI containers on a Linux system.
Podman also works on Mac using a [podman machine](https://docs.podman.io/en/latest/markdown/podman-machine.1.html).

## Prerequisites
 - [Install podman](https://podman.io/getting-started/installation)
 - Mac: ensure a [podman machine](https://docs.podman.io/en/latest/markdown/podman-machine.1.html) is running.
 - Linux: for [multi-platform builds](https://docs.earthly.dev/docs/guides/multi-platform), install [qemu-user-static](https://github.com/multiarch/qemu-user-static).
 - [WITH DOCKER](https://docs.earthly.dev/docs/earthfile#with-docker) requires rootful mode.
   - Linux: run with `sudo` (i.e., `sudo earthly -P +with-docker-target`)
   - Mac: run a [rootful machine](https://docs.podman.io/en/latest/markdown/podman-machine-set.1.html#rootful).

## Getting started
When earthly starts a check is done to determine what frontend is available.
By default, earthly will attempt to use docker and then fall back to podman.
If you wish to change the behavior of the startup check, run the following command:

```bash
# Configure earthly to use podman
earthly config global.container_frontend podman-shell

# Configure earthly to use docker
earthly config global.container_frontend docker-shell
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

           buildkitd | Starting buildkit daemon as a podman container (earthly-buildkitd)...
           buildkitd | ...Done
```

If instead you see `No frontend initialized`, and you're using Mac, it may mean your podman machine is not running.

## Known limitations / troubleshooting
### Builds running slowly
There are a few steps you should take to rule out common performance bottlenecks.

#### Mac: check podman resources
At the time of writing this, podman machines use a single core and 2GB of RAM by default. 
Depending on what you're doing you may need more resources.

Resources can be adjusted by using one of these commands:
```bash
# Initialize a new default machine with 5 CPUs, 128GB disk space, 8196 MB of memory, and start it
podman machine init --now --cpus 5 --disk-size 128 --memory 8196 

# Adjust the current default podman machine to use 5 CPUs, 128GB disk space, and 8196 MB of memory
podman machine stop ; podman machine set --cpus 5 --disk-size 128 --memory 8196 && podman machine start
```

### Mac: check machine architecture
Running `podman version` will display the specifications of your podman client and server (machine).
You should ensure the architecture in OS/Arch is the same between client and server.
This will rule out emulation as a performance bottleneck.

The output may look like this:
```bash
> podman version
Client:       Podman Engine
Version:      4.2.1
API Version:  4.2.1
Go Version:   go1.18.6
Built:        Tue Sep  6 13:16:02 2022
OS/Arch:      darwin/arm64

Server:       Podman Engine
Version:      4.2.0
API Version:  4.2.0
Go Version:   go1.18.4
Built:        Thu Aug 11 08:43:11 2022
OS/Arch:      linux/arm64
```
In this example, the client us running on an M1 Mac and both the client and server are using arm64.

### Check graph driver
Running `podman info --debug` will show your current podman configuration.
VFS and other drivers can perform poorly when compared to overlay and 
[are not recommended by the podman community](https://github.com/containers/podman/issues/13226).
Ensure overlay is used by looking for the following in the podman info output:
```bash
> podman info --debug

...
graphDriverName: overlay  # or something similar
...
```

### Mac: docker-credential-desktop: executable file not found in $PATH
This error typically occurs when switching from docker desktop to podman without docker installed.
There may be a lingering configuration file that will be read by the attachable used to authenticate calls to BuildKit.

To fix this issue, try removing or renaming the `~/.docker/config.json` file.

### Earthly CLI - no frontend initialized
Seeing the error on startup means the check for podman has failed.
```bash
> earthly github.com/earthly/hello-world:main+hello
 1. Init 🚀
————————————————————————————————————————————————————————————————————————————————

            frontend | No frontend initialized.
```

Ensure you have correctly installed podman and, if you are using a Mac, the podman machine is running.
```bash
> podman machine start
```

### Rootless podman
Running podman in rootless mode is not supported due to the [earthly/dind](https://hub.docker.com/r/earthly/dind) and 
[earthly/buildkit](https://hub.docker.com/r/earthly/buildkitd) because they [require privileged access](https://docs.earthly.dev/docs/guides/using-the-earthly-docker-images/buildkit-standalone#requirements).
Specifically, [WITH DOCKER](https://docs.earthly.dev/docs/earthfile#with-docker) will fail.
You must use `sudo` on Linux or [set your podman machine to rootful mode on Mac](https://docs.podman.io/en/latest/markdown/podman-machine-set.1.html#rootful) to use [WITH DOCKER](https://docs.earthly.dev/docs/earthfile#with-docker).

### Podman within WITH DOCKER
[WITH DOCKER](https://docs.earthly.dev/docs/earthfile#with-docker) starts a container with a docker installation. 
You can only use the podman CLI in the RUN statement if you specify [LOCALLY](https://docs.earthly.dev/best-practices#pattern-optionally-locally)
to run it on the host machine; otherwise, you will need to use the docker CLI.

```bash
docker-locally:
   LOCALLY
   WITH DOCKER
     RUN podman ps
   END
```

```bash
docker:
   WITH DOCKER
     RUN docker ps
   END
```

### Cross-image targets
You need to configure QEMU if you are running a cross-platform target.
If you haven't properly configured QEMU you will receive an error message containing the following message:
```bash
> earthly +cross-platform
...
exec /bin/sh: exec format error
...
```

We've found installing [qemu-user-static](https://github.com/multiarch/qemu-user-static) will allow cross-platform targets tun run on Linux.
```bash
> apt-get install qemu-user-static
# or
> yum install qemu-user-static
```

### crun: open executable: Permission denied: OCI permission denied
This can happen if you attempt to run (or the `ENTRYPOINT` references) a binary without the execution permission.
https://github.com/containers/podman/issues/9377
https://github.com/signalwire/freeswitch/pull/1748
