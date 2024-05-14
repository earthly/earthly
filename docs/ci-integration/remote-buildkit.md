# Remote BuildKit 

{% hint style='info' %}
##### Note
This guide is related to self-hosting a remote BuildKit, however, Self-Hosted Satellites **beta** are now available. Self-Hosted Satellites provide more features, have better security, and are easier to deploy than remote BuildKit. Check out the [Self-Hosted Satellites Guide](../cloud/satellites/self-hosted.md) for more details. If your use case cannot tolerate a cloud-based control plane, however, then self-hosting a remote BuildKit is still supported.
{% endhint %}

## Introduction

In some cases, you may want to run a remote instance of [`earthly/buildkitd`](https://hub.docker.com/r/earthly/buildkitd). This guide is intended to help you identify if you might benefit from this configuration, and to help you set it up correctly.

If you are looking for a way to use remote runners without the complexities of managing them yourself, you may want to consider [Earthly Satellites](../cloud/satellites.md): remote runners managed by Earthly.

### Why Remote?

Running a remote daemon is a unique feature of Earthly. It allows the build to happen elsewhere; even when executing it from your local development machine. The key benefit of remote execution is having instant access to the cache, thus making builds dramatically faster compared to many traditional CI setups that require uploading and downloading the cache.

### Running Remote BuildKit

To run a remote BuildKit instance, deploy and configure the image [`earthly/buildkitd`](https://hub.docker.com/r/earthly/buildkitd).

#### Networking

A remote daemon should be reachable by all clients intending to use it. Earthly uses ports `8371-8373` to communicate, so these should be open and available.

#### Mounts

**`/tmp/earthly`**

This path within the container is the location that BuildKit uses for storing the cache. Because this folder sees _a lot_ of traffic, its important that it remains fast.

{% hint style='danger' %}
##### Important
We *strongly* recommend using a Docker volume for mounting `/tmp/earthly`. If you do not, `buildkitd` can consume excessive disk space, operate very slowly, or it might not function correctly.
{% endhint %}

In some environments, not mounting `/tmp/earthly` as a Docker volume results in the following error:

```
--> WITH DOCKER RUN --privileged ...
...
rm: can't remove '/var/earthly/dind/...': Resource busy
```

In EKS, users reported that mounting an EBS volume, instead of a Kubernetes `emptyDir` worked.

This part of our documentation needs improvement. If you have a Kubernetes-based setup, please [let us know](https://earthly.dev/slack) how you have mounted `EARTHLY_TMP_DIR` and whether `WITH DOCKER` worked well for you.

#### Daemon

To configure an `earthly/buildkitd` daemon as a remotely available daemon, you will need to start the container yourself. See our [configuration docs](../earthly-config/earthly-config.md) for more details on all the options available; but here are the ones you need to know:

**`BUILDKIT_TCP_TRANSPORT_ENABLED`**

This will configure `buildkitd` to listen on port `8372`. If you would like it to be externally available on a different port, you will need to handle that at the port mapping level. TCP is required for remotely sharing a daemon.

**`BUILDKIT_TLS_ENABLED`**

Set this to `true` for all daemons that will handle production workloads. This daemon *by design* is an arbitrary code execution machine, and running it without any kind of mTLS configuration is not recommended.

Make sure you mount your certificates and keys in the correct location (`/etc/*.pem`).

For complete details, see the [documentation for `earthly/buildkitd` on DockerHub](https://hub.docker.com/r/earthly/buildkitd).

#### Client

Normally, Earthly will try to start and manage its own `earthly/buildkitd` daemon. However, when relying on a remote `earthly/buildkitd` instance, Earthly will not attempt to manage this daemon. Here are the configuration options needed to use a remote instance:

**`buildkit_host`**

This is the address of the remote daemon. It should look something like this: `tcp://my-cool-remote-daemon:8372`. If the hostname is considered to be a "local" one, Earthly will fall back to the Local-Remote behaviors described below. For reference; all IPv6 Loopback addresses, `127.0.0.1`, and `[localhost](http://localhost)` are considered to be "local". The machine's hostname is not considered "local".

**`tlsca` / `tlscert` / `tlskey`**

These are the paths to the certificates and keys used by the client when communicating with an mTLS-enabled daemon. These paths are relative to the Earthly config (usually `~/.earthly/config.yaml`, unless absolute paths are specified.

**`tls_enabled`**

TLS will be enabled by default (unless using a local buildkit container).  
Set this to `false` when using TLS is not desired.

### Local-Remote

It is also possible to use the remote protocols (TCP and mTLS) locally, while still letting Earthly manage the daemon container.  
Earthly will (optionally) generate its own certificates, and connect to the daemon using `tcp://127.0.0.1:8372`. This is a great way to test some of the remote capabilities without having to generate certificates or manage a separate machine.
