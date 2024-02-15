# Kubernetes

{% hint style='info' %}
##### Note
This guide is related to self-hosting a remote Buildkit, however, Self-Hosted Satellites **beta** are now available. Self-Hosted Satellites provide more features, have better security, and are easier to deploy than remote Buildkit. Check out the [Self-Hosted Satellites Guide](../../cloud/satellites/self-hosted.md) for more details and instructions to deploy in Kubernetes or AWS EC2.
{% endhint %}


## Overview

Kubernetes isn't a CI per-se, but it _can_ serve as the underpinning for many modern CI systems. As such, this example serves as a bare-bones example to base your implementations on.

### Compatibility

`earthly` has been tested with the all-in-one `earthly/earthly` mode, and works as long as the pod runs in a `privileged` mode.

It has also been tested with a _single_ remote `earthly/buildkitd` running in `privileged` mode, and an `earthly/earthly` pod running without any additional security concerns. This configuration is considered experimental. See [these additional instructions](../remote-buildkit.md).

Multi-node `earthly/buildkitd` configurations are currently unsupported.

### Resources

 * [Kubernetes Documentation](https://kubernetes.io/docs/home/supported-doc-versions/)
 * [Kubernetes Taints & Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)

## Setup (`earthly/earthly` Only)

This is the recommended approach when using Earthly within Kubernetes. Assuming you are following the steps outlined in the [overview](../overview.md), here are the additional things you need to configure:

### Dependencies

Your Kubernetes cluster needs to allow `privileged` mode pods. It's possible to use a separate instance group, along with Taints and Tolerations to effectively segregate these pods.

### Installation

The default image from `earthly/earthly` should be sufficient. If you need additional tools or configuration, you can [create your own runner image](../build-an-earthly-ci-image.md).

### Configuration

In some instances, notably when using [Calico](https://www.tigera.io/project-calico/) within your cluster, the MTU of the clusters network may end up mismatched with the internal CNI network, preventing external communication. You can set this through the `CNI_MTU` environment variable to force a match.

`earthly/earthly` currently requires the use of privileged mode. Use this in your container spec to enable it:

```yaml
securityContext:
  privileged: true
```

The `earthly/earthly` container will operate best when provided with decent storage for intermediate operations. Mount a volume like this:

```yaml
volumeMounts:
  - mountPath: /tmp/earthly
    name: buildkitd-temp
...
volumes:
  - name: buildkitd-temp
    emptyDir: {} # Or other volume type
```

The location within the container for this temporary folder is configurable with the `EARTHLY_TMP_DIR` environment variable.

The `earthly/earthly` image will expect to find the source code (with `Earthfile`) rooted in the default working directory, which is set to `/workspace`.

## Setup (Remote `earthly/buildkitd`)

{% hint style='danger' %}
##### Note

This an _experimental_ configuration.

{% endhint %}

It is possible to run multiple `earthly/buildkitd` instances in Kubernetes, for larger deployments. Follow the configuration instructions for using the `earthly/earthly` image above.

There are some caveats that come with this kind of a setup, though:

1. Some local cache is not available across instances, so it may take a while for the cache to become "warm".
2. Builds that occur across multiple instances simultaneously may fail in odd ways. This is not supported.
3. The TLS configuration needs to be shared across the entire fleet.

### Configuration

To mitigate some of the issues, it is recommended to run in a "sticky" mode to keep builds pinned to a single instance for the duration. You can see how to do this in our example:

```yaml
# Use session affinity to prevent "roaming" across multiple buildkit instances; if needed.
sessionAffinity: ClientIP
sessionAffinityConfig:
  clientIP:
    timeoutSeconds: 600 # This should be longer than your build.
```

## Example

{% hint style='danger' %}
##### Note

This example is not production ready, and is intended to showcase configuration needed to get Earthly off the ground. If you run into any issues, or need help, [don't hesitate to reach out](https://github.com/earthly/earthly/issues/new)!

{% endhint %}

See our [Kubernetes examples](https://github.com/earthly/ci-examples/tree/main/kubernetes).

To run it yourself, first you will need to install some prerequisites on your machine. This example requires `kind` and `kubectl` to be installed on your system. Here are some links to installation instructions:

- [`kind` Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

When you are ready, clone the [`ci-examples` repository](https://github.com/earthly/ci-examples), and then run (from the root of the repository):

```go
earthly ./kubernetes+start
```

Running this target will:

- Create a `kind` cluster named `earthlydemo-aio`
- Create & watch a `job` that runs an `earthly` build

When the example is complete, the cluster is left up and intact for exploration and experimentation. If you would like to clean up the cluster when complete, run (again from the root of the repository):

```go
earthly ./kubernetes+clean
```
