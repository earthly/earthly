This image contains `earthly`, `buildkit`, and some extra configuration to enable the two to work together. All that's missing is your source code! This image is mainly intended for use in containerized CI scenarios, or where maintaining a persistent installation of `earthly` isn't possible.

## Tags

Currently, the `latest` tag is `v0.8.4`.  
For other available tags, please check out https://hub.docker.com/r/earthly/earthly/tags

## Quickstart

Want to get started? Here are a couple sample `docker run` commands that cover the most common use-cases:

### Usage with Docker Socket

This example shows how to use the Earthly container in conjunction with a Docker socket that Earthly can use to start up the Buildkit daemon.

```bash
docker run -t -v $(pwd):/workspace -v /var/run/docker.sock:/var/run/docker.sock -e NO_BUILDKIT=1 earthly/earthly:v0.8.4 +for-linux
```

Here's a quick breakdown:

- `-t` tells Docker to emulate a TTY. This makes the `earthly` log output colorized.
- `-v $(pwd):/workspace` mounts the source code into the conventional location within the docker container. Earthly is executed from this directory when starting the container. Any artifacts saved within this folder remain on your local machine.
- `-v /var/run/docker.sock:/var/run/docker.sock` mounts the Docker socket such that Earthly can start Buildkit as a Docker container in the host's Docker.
- `-e NO_BUILDKIT=1` tells the Earthly container not to start en embedded buildkit. A Buildkit daemon will instead be started via the Docker socket provided.
- `+for-linux` is the target to be invoked. All arguments specified after the image tag will be passed to `earthly`.

### Usage with Embedded Buildkit

This example shows how the Earthly image can start a Buildkit daemon within the same container. A Docker socket is not needed in this case, however the container will need to be run with the `--privileged` flag.

```bash
docker run --privileged -t -v $(pwd):/workspace -v earthly-tmp:/tmp/earthly:rw earthly/earthly:v0.8.4 +for-linux
```

Here's a quick breakdown:

- `--privileged` is required when you are using the internal, embedded `buildkit`. This is because `buildkit` currently requires it for OverlayFS support and for network configuration.
- `-t` tells Docker to emulate a TTY. This makes the `earthly` log output colorized.
- `-v $(pwd):/workspace` mounts the source code into the conventional location within the docker container. Earthly is executed from this directory when starting the container. Any artifacts saved within this folder remain on your local machine.
- `-v earthly-tmp:/tmp/earthly:rw` mounts (and creates, if necessary) the `earthly-tmp` Docker volume into the containers `/tmp/earthly`. This is used as a temporary/working directory for `buildkitd` during builds.
- `+for-linux` is the target to be invoked. All arguments specified after the image tag will be passed to `earthly`.

### Usage with Satellites and No Local Code

This example utilizes an [Earthly Satellite](https://docs.earthly.dev/earthly-cloud/satellites) to perform builds. The code to be built is downloaded directly from GitHub.

```bash
docker run -t -e NO_BUILDKIT=1 -e EARTHLY_TOKEN=<my-token> earthly/earthly:v0.8.4 --ci --org <my-org> --sat <my-sat> github.com/earthly/earthly+for-linux
```

Here's what this does:

- `-e EARTHLY_TOKEN=<my-token>` passes along an Earthly token such that Earthly can access satellites. This token can be created via `earthly account create-token`.
- `--org <my-org>` specifies the organization that the satellite belongs to.
- `--sat <my-sat>` specifies the satellite to use.
- `github.com/earthly/earthly+for-linux` specifies the target to build. This target is located on GitHub, and will be pulled from the Satellite.

### Usage for non-build commands

This example shows how to use the Earthly container to run non-build commands. This is useful for running commands like `earthly account`, or `earthly secret`.

```bash
docker run -t -e NO_BUILDKIT=1 -e EARTHLY_TOKEN=<my-token> earthly/earthly:v0.8.4 account list-tokens
```

```bash
docker run -t -e NO_BUILDKIT=1 -e EARTHLY_TOKEN=<my-token> earthly/earthly:v0.8.4 secret get foo
```

## Using This Image

### Requirements

There are a couple requirements this image expects you to follow when using it. These requirements streamline usage of the image and save configuration effort.

#### Privileged Mode

If you are using the embedded `buildkitd`, then this image needs to be run as a privileged container. This is because `buildkitd` needs appropriate access to use `overlayfs`.

#### `/tmp/earthly`

Because this folder sees _a lot_ of traffic, its important that it remains fast. We *strongly* recommend using a Docker volume for mounting `/tmp/earthly`. If you do not, `buildkitd` can consume excessive disk space, operate very slowly, or it might not function correctly.

In some environments, not mounting `/tmp/earthly` as a Docker volume results in the following error:

```
--> WITH DOCKER RUN --privileged ...
...
rm: can't remove '/var/earthly/dind/...': Resource busy
```

#### Source Mounting

Because `earthly` is running inside a container, it does not have access to your source code unless you grant it. This image expects to find a valid `Earthfile` in the working directory, which is set by default to `/workspace`.

#### DOCKER_HOST

This image *does* include a functional Docker CLI, but does not include a full Docker daemon. If your `Earthfile` requires a Docker daemon of any sort, you will need to provide it through this environment variable.

If your daemon is on the same host as this container, you can also volume mount your hosts docker daemon using `-v /var/run/docker.sock:/var/run/docker.sock`. Note that this will cause `earthly` to use your hosts Docker daemon, and could lead to name conflicts if multiple copies of this image are run on the same host.

#### -t

This is the easiest way to ensure you get the nice, colorized output from `earthly`. Alternatively, you could provide the `FORCE_COLOR` environment variable.

### Supported Environment Variables

| Variable Name                       | Default Values                 | Description                                                                                                                                                                                                   |
|-------------------------------------|--------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GLOBAL_CONFIG                       |                                | Any valid YAML for the top-level `global` key in `config.yml`. Example: `{disable_analytics: true, local_registry_host: 'tcp://127.0.0.1:8371'}`                                                              |
| GIT_CONFIG                          |                                | Any valid YAML for the top-level `git` key in `config.yml`. Example: `{example: {pattern: 'example.com/([^/]+)', substitute: 'ssh://git@example.com:2222/var/git/repos/$1.git', auth: ssh}}`                  |
| NO_BUILDKIT                         |                                | Disables the embedded Buildkit daemon.                                                                                                                                                                        |
| DOCKER_HOST                         | `/var/run/docker.sock`         | From Docker's CLI.                                                                                                                                                                                            |
| BUILDKIT_HOST                       | `tcp://<hostname>:8372`        | The address of your BuildKit host. Use this when you have a remote `buildkitd` you would like to connect to.                                                                                                  |
| EARTHLY_ADDITIONAL_BUILDKIT_CONFIG  |                                | Additional `buildkitd` config to append to the generated configuration file.                                                                                                                                  |
| BUILDKIT_TCP_TRANSPORT_ENABLED      |                                | Required to be set to `true` when using an external `buildkitd` via `BUILDKIT_HOST`. `true` when using the embedded `buildkitd`.                                                                              |
| BUILDKIT_TLS_ENABLED                |                                | Required when using an external `buildkitd` via `BUILDKITD_HOST`, and the external `buildkitd` requires mTLS. You will also need to mount certificates into the right place (`/etc/.earthly/certs`).          |
| CNI_MTU                             | MTU of first default interface | Set this when we autodetect the MTU incorrectly. The device used for autodetection can be shown by the command  `ip route show \| grep default \| cut -d' ' -f5 \| head -n 1`                                 |
| EARTHLY_RESET_TMP_DIR               | `false`                        | Cleans out `/tmp/earthly` before running, if set to `true`. Useful when you host-mount an temporary directory across runs.                                                                                    |
| NETWORK_MODE                        | `cni`                          | Specifies the networking mode of `buildkitd`. Default uses a CNI bridge network, configured with the `CNI_MTU`.                                                                                               |
| CACHE_SIZE_MB                       | `0`                            | How big should the `buildkitd` cache be allowed to get, in MiB? A value of 0 sets the cache size to "adaptive", causing BuildKit to detect the available size of the system and choose a limit automatically. |
| GIT_URL_INSTEAD_OF                  |                                | Configure `git config --global url.<url>.insteadOf` rules to be used by `buildkitd`.                                                                                                                          |
| IP_TABLES                           |                                | Override which binary (`iptables_nft` or `iptables_legacy`) is used for configuring `ip_tables`. Only set this if autodetection fails for your platform.                                                      |
