This image contains `buildkit` with some Earthly-specific setup. This is what Earthly will start when using a local daemon. You can also start it up yourself and use it as a remote/shared BuildKit daemon.

*Note that versions of this container have only ever been tested with their corresponding version of `earthly`.* Mismatched versions are unsupported.

## Tags

Currently, the `latest` tag is `v0.8.10`.  
For other available tags, please check out https://hub.docker.com/r/earthly/buildkitd/tags

## Quickstart

Want to just get started? Here are a couple sample `docker run` commands that cover the most common use-cases:

### Simple Usage (Use Locally)

```bash
docker run --privileged -t -v earthly-tmp:/tmp/earthly:rw earthly/buildkitd:v0.7.19
```

Heres a quick breakdown:

- `--privileged` is required. This is because `earthly` needs some privileged `buildkit` functionality.
- `-t` tells Docker to emulate a TTY. This makes the `buildkit` log output colorized.
- `-v earthly-tmp:/tmp/earthly:rw` mounts (and creates, if necessary) the `earthly-tmp` Docker volume into the containers `/tmp/earthly`. This is used as a temporary/working directory for `buildkitd` during builds.

Assuming you are running this on your machine, you could use this `buildkitd` by setting `EARTHLY_BUILDKIT_HOST=docker-container://<container-name>`, or by specifying the appropriate values in `config.yml`.

### Usage (Use As Remote)

```bash
docker run --privileged -t -v earthly-tmp:/tmp/earthly:rw -e BUILDKIT_TCP_TRANSPORT_ENABLED=true -p 8372:8372 earthly/buildkitd:v0.7.19
```

Omitting the options already discussed from the simple example:

- `-e BUILDKIT_TCP_TRANSPORT_ENABLED=true` makes `buildkitd` listen on a TCP port instead of a Unix socket.
- `-p 8372:8372` forwards the hosts port 8372 to the containers port 8372. When using TCP, `buildkit` will always listen on 8372, but you can configure the apparent port by forwarding a different port on your host.

Assuming you ran this on another machine named `fast-builder`, you could use this remote `buildkitd` by setting `EARTHLY_BUILDKIT_HOST=tcp://fast-builder:8372`, or by specifying the address in your `config.yml`.

## Using This Image

### Requirements

#### Privileged Mode

This image needs to be run as a privileged container. This is because `buildkitd` needs appropriate access to start and run additional containers itself via `runc`.

#### `/tmp/earthly`

Because this folder sees _a lot_ of traffic, its important that it remains fast. We *strongly* recommend using a Docker volume for mounting `/tmp/earthly`. If you do not, `buildkitd` can consume excessive disk space, operate very slowly, or it might not function correctly.

In some environments, not mounting `/tmp/earthly` as a Docker volume results in the following error:

```
--> WITH DOCKER RUN --privileged ...
...
rm: can't remove '/var/earthly/dind/...': Resource busy
```

#### External Usage

To use this image externally, it requires you to forward a port on your machine to the containers port 8372. You will need to ensure that external access to the machine on the port you chose is possible as well.

When using this container locally with `earthly`, please note that setting `EARTHLY_BUILDKIT_HOST` values with hosts `127.0.0.1`, ` ::1/128`, or `localhost` are considered local and will result in Earthly attempting to manage the BuildKit container itself. Consider using your hostname, or another alternative name in these cases.

### Supported Environment Variables

| Variable Name                       | Default Values                 | Description                                                                                                                                                                   |
|-------------------------------------|--------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| EARTHLY_ADDITIONAL_BUILDKIT_CONFIG  |                                | Additional `buildkitd` config to append to the generated configuration file.                                                                                                  |
| BUILDKIT_TCP_TRANSPORT_ENABLED      |                                | Set to `true` when the `buildkitd` instance is going to be used remotely                                                                                                      |
| BUILDKIT_TLS_ENABLED                |                                | Set to `true` when the `buildkitd` instance will require mTLS from the clients. You will also need to mount certificates into the right place (`/etc/*.pem`).                |
| CNI_MTU                             | MTU of first default interface | Set this when we autodetect the MTU incorrectly. The device used for autodetection can be shown by the command  `ip route show \| grep default \| cut -d' ' -f5 \| head -n 1` |
| EARTHLY_RESET_TMP_DIR               | `false`                        | Cleans out `/tmp/earthly` before running, if set to `true`. Useful when you host-mount an temporary directory across runs.                                                            |
| NETWORK_MODE                        | `cni`                          | Specifies the networking mode of `buildkitd`. Default uses a CNI bridge network, configured with the `CNI_MTU`.                                                               |
| CACHE_SIZE_MB                       | `0`                            | How big should the `buildkitd` cache be allowed to get, in MiB? 0 is unbounded.                                                                                               |
| GIT_URL_INSTEAD_OF                  |                                | Configure `git config --global url.<url>.insteadOf` rules used by `buildkitd`                                                                                                 |
| IP_TABLES                           |                                | Override which binary (`iptables_nft` or `iptables_legacy`) is used for configuring `ip_tables`. Only set this if autodetection fails for your platform.                          |
