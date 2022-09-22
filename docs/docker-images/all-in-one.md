This image contains `earthly`, `buildkit`, and some extra configuration to enable the two to work together. All that's missing is your source code! This image is mainly intended for use in containerized CI scenarios, or where maintaining a persistent installation of `earthly` isn't possible.

## Tags

* `v0.6.24`, `latest`
* `v0.6.23`
* `v0.6.22`

## Quickstart

Want to just get started? Here are a couple sample `docker run` commands that cover the most common use-cases:

### Simple Usage

```bash
docker run --privileged -t -e NO_DOCKER=1 -v $(pwd):/workspace -v earthly-tmp:/tmp/earthly:rw earthly/earthly:v0.6.24 +for-linux
```

Heres a quick breakdown:

- `--privileged` is required when you are using the internal `buildkit`. This is because `buildkit` currently requires it for OverlayFS support.
- `-t` tells Docker to emulate a TTY. This makes the `earthly` log output colorized.
- `-e NO_DOCKER=1` skips the check for a functional Docker daemon. If you are not exporting any images, you do not need Docker.
- `-v $(pwd):/workspace` mounts the source code into the conventional location within the docker container. Earthly is executed from this directory when starting the container. Any artifacts saved within this folder remain on your local machine.
- `-v earthly-tmp:/tmp/earthly:rw` mounts (and creates, if necessary) the `earthly-tmp` Docker volume into the containers `/tmp/earthly`. This is used as a temporary/working directory for `buildkitd` during builds.
- `+for-linux` is the target to be invoked. All arguments specified after the image tag will be passed to `earthly`.

### More Complicated Usage

```bash
docker run -t --privileged -v $(pwd):/workspace:rw -v earthly-tmp:/tmp/earthly:rw --network=host -v ~/.earthly/config.yml:/etc/.earthly/config.yml -e DOCKER_HOST="tcp://192.168.1.234:2375" earthly/earthly:v0.6.24 --ci -P +for-linux
```

Omitting the options already discussed from the simple example:

- `--network=host` runs the container on the host machines network. This is necessary so (in this example) `earthly` can reach the Docker host.
- `-v ~/.earthly/config.yml:/etc/.earthly/config.yml` mounts a custom `earthly` configuration file in the conventional location.
- `-e DOCKER_HOST="tcp://192.168.1.234:2375"` specifies the external Docker host that will be used. This endpoint will be checked for accessibility before the target is built.
- `--ci` run `earthly` in CI mode.
- `-P` runs `earthly` in privileged mode.

## Using This Image

### Requirements

There are a couple requirements this image expects you to follow when using it. These requirements streamline usage of the image and save configuration effort.

#### Privileged Mode

If you are using the baked-in `buildkitd`, then this image needs to be run as a privileged container. This is because `buildkitd` needs appropriate access to use `overlayfs`.

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
| NO_DOCKER                           |                                | Disables the check for a working Docker Daemon. Setting this _at all_ disables the check.                                                                                                                     |
| DOCKER_HOST                         | `/var/run/docker.sock`         | From Docker's CLI.                                                                                                                                                                                            |
| BUILDKIT_HOST                       | `tcp://<hostname>:8372`        | The address of your BuildKit host. Use this when you have a remote `buildkitd` you would like to connect to.                                                                                                  |
| EARTHLY_ADDITIONAL_BUILDKIT_CONFIG  |                                | Additional `buildkitd` config to append to the generated configuration file.                                                                                                                                  |
| BUILDKIT_TCP_TRANSPORT_ENABLED      |                                | Required to be set to `true` when using an external `buildkitd` via `BUILDKIT_HOST`. `true` when using the baked-in `buildkitd`.                                                                              |
| BUILDKIT_TLS_ENABLED                |                                | Required when using an external `buildkitd` via `BUILDKITD_HOST`, and the external `buildkitd` requires mTLS. You will also need to mount certificates into the right place (`/etc/.earthly/certs`).          |
| CNI_MTU                             | MTU of first default interface | Set this when we autodetect the MTU incorrectly. The device used for autodetection can be shown by the command  `ip route show \| grep default \| cut -d' ' -f5 \| head -n 1`                                 |
| EARTHLY_RESET_TMP_DIR               | `false`                        | Cleans out `/tmp/earthly` before running, if set to `true`. Useful when you host-mount an temporary directory across runs.                                                                                 |
| NETWORK_MODE                        | `cni`                          | Specifies the networking mode of `buildkitd`. Default uses a CNI bridge network, configured with the `CNI_MTU`.                                                                                               |
| CACHE_SIZE_MB                       | `0`                            | How big should the `buildkitd` cache be allowed to get, in MiB? A value of 0 sets the cache size to "adaptive", causing BuildKit to detect the available size of the system and choose a limit automatically. |
| GIT_URL_INSTEAD_OF                  |                                | Configure `git config --global url.<url>.insteadOf` rules to be used by `buildkitd`.                                                                                                                          |
| IP_TABLES                           |                                | Override which binary (`iptables_nft` or `iptables_legacy`) is used for configuring `ip_tables`. Only set this if autodetection fails for your platform.                                                      |
