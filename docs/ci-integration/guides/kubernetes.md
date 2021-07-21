# Kubernetes

# Example

You can find our [Kubernetes example here](https://github.com/earthly/ci-examples/tree/main/kubernetes).

To run it yourself, first you will need to install some prerequisites on your machine. This example requires `kind` and `kubectl` to be installed on your system. Here are some links to installation instructions:

- `[kind` Quick Start](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

When you are ready, clone the `ci-examples` repository, and then run (from the root of the repository):

```go
earthly ./kubernetes+run
```

Running this target will:

- Create a `kind` cluster named `earthlydemo`
- Create a `deployment` with 1 replica of `[earthly/buildkitd](https://hub.docker.com/r/earthly/buildkitd)`, (the same container `earthly` will start when running locally)
- Create a `service` to front the `deployment`
- Create & watch a `job` that runs an `earthly` build against this independent, external `earthly-buildkitd` deployment

When the example is complete, the cluster is left up and intact for exploration experimentation. If you would like to clean up the cluster when complete, run (again from the root of the repository):

```go
earthly ./kubernetes+cleanup
```

# Details

`earthly` has been tested in both the all-in-one `earthly/earthly` mode, and in a remote mode using `earthly/buildkitd`. 

If you are using a remote `earthly-buildkitd` daemon, please see that page for configuration considerations.

### `earthly` Configuration Notes

The `earthly/earthly` image will expect to find the source code (with `Earthfile`) rooted in `/workspace`. To configure this, ensure that the `SRC_DIR` environment variable is set correctly. In the case of the example, we are building a remote target, so mounting a dummy volume is needed.

### Kubernetes Configuration Notes

In some instances, notably when using Calico within your cluster, the MTU of the clusters network may end up mismatched with the internal CNI network, preventing external communication. You can set this through the `CNI_MTU` environment variable to force a match.

`earthly/buildkitd` currently requires the use of privileged mode. Use this in your container spec to enable it:

```yaml
securityContext:
  privileged: true
```

The `earthly/buildkitd` container will operate best when provided with decent storage for intermediate operations. Mount in a volume like this:

```yaml
volumeMounts:
  - mountPath: /tmp/earthly
    name: buildkitd-temp
...
volumes:
  - name: buildkitd-temp
    emptyDir: {} # Or other volume type
```

The location within the container for this temp folder is configurable with the `EARTHLY_TMP_DIR` environment variable.

It is possible to run multiple `earthly/buildkitd` instances in Kubernetes, for larger deployments. There are some caveats that come with it, though:

 1. Some local cache is not available across instances, so it may take a while for the cache to become "warm".
 2. Builds that occur across multiple instances simultaneously may fail in odd ways. This is not supported.
 3. The TLS configuration needs to be shared across the entire fleet.

To mitigate some of the issues, it is recommended to run in a "sticky" mode to keep builds pinned to a single instance for the duration. You can see how to do this in our example:

```yaml
# Use session affinity to prevent "roaming" across multiple buildkit instances; if needed.
sessionAffinity: ClientIP
sessionAffinityConfig:
  clientIP:
    timeoutSeconds: 600 # This should be longer than your build.
```