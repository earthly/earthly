# Remote Buildkit

## Introduction

In some cases, you may want to run a remote instance of `earthly/buildkitd`. This guide is intended to help you identify if you might benefit from this configuration, and to help you set it up correctly.

### Why Remote?

Running a remote daemon is a unique feature of Earthly. It allows the build to happen elsewhere; even when executing it from your local development machine. However, it is not always the best option. Before setting up a remote daemon, first look into Earthly's shared caching capabilities and see if those can get you the boost you need. In our experience, shared caching is usually enough.

However, there are instances where a remote daemon can make the most sense. Here are some examples:

- You have a single, powerful build machine you would like to share with your development team
- There is data closer to a remote machine than your development/CI environment, so you bring the build to the data
- You are using Earthly in Kubernetes, and want to isolate the containers doing the actual building because they require privileged mode
- You want to share a build machine (or cluster) with your CI environment and your developers
- Your local computer does not have the capabilities to build the software (`docker` is missing, or you lack sufficient privileges, or it is simply not powerful enough)

### Configuring A Remote Cache

#### Networking

A remote daemon should be reachable by all clients intending to use it. Earthly uses ports `8371-8373` to communicate, so these should be open and available to communicate.

#### Daemon

To configure an `earhly/buildkitd` daemon as a remotely available daemon, you will need to start the container yourself. See our configuration docs for more details on all the options available; but here are the ones you need to know:

**`BUILDKIT_TCP_TRANSPORT_ENABLED`**

This will configure `buildkitd` to listen on port `8372`. If you would like it to be externally available on a different port, you will need to handle that at the port mapping level. TCP is required for remotely sharing a daemon.

**`BUILDKIT_TLS_ENABLED`**

Set this to `true` for all daemons that will handle production workloads. This daemon *by design* is an arbitrary code execution machine, and running it without any kind of mTLS configuration is not recommended.

Make sure you mount your certificates and keys in the correct location (`/etc/*.pem`).

#### Client

Normally, Earthly will try to start and manage its own `earthly/buildkitd` daemon. However, when relying on a remote `earthly/buildkitd` instance, Earthly will not attempt to manage this daemon. Here are the configuration options needed to use a remote instance:

**`buildkit_host`**

This is the address of the remote daemon. It should look something like this: `tcp://my-cool-remote-daemon:8372`. If the hostname is considered to be a "local" one, Earthly will fall back to the Local-Remote behaviors described below. For reference; all IPv6 Loopback addresses, `127.0.0.1`, and `[localhost](http://localhost)` are considered to be "local". The machine's hostname is not considered "local".

**`buildkit_transport`**

Set this to `tcp` for connection to the remote daemon.

**`tlsca` / `tlscert` / `tlskey`**

These are the paths to the certificates and keys used by the client when communicating with an mTLS-enabled daemon. These paths are relative to the Earthly config (usually `~/.earthly/config.yaml`, unless absolute paths are specified.

**`tls_enabled`**

Set this to `true` when using TLS is desired.

### Local-Remote

It is also possible to use the remote protocols (TCP and mTLS) locally, while still letting Earthly manage the daemon container. You can do this by enabling TCP transport(`buildkit_transport`), and enabling mTLS(`tls_enabled`).

By doing this, Earthly will (optionally) generate its own certificates, and connect to the daemon using `tcp://127.0.0.1:8372`. This is a great way to test some of the remote capabilities without having to generate certificates or manage a separate machine.