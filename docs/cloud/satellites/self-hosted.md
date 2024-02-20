# Self-Hosted Satellites

Self-hosted satellites **Beta** are a good alternative to Earthly Cloud’s [managed satellites](../satellites.md) when CI/CD is required to run in your own environment. Self-hosted satellites provide most of the benefits of Earthly Cloud while ensuring that your code and build-related data never leave your network.

Self-Hosted Satellites have the following features:
* Automatic and instantly available build cache that makes builds faster
* Cloud-enabled control plane
* Encrypted connections by default using mTLS
* Ready to run builds with just a single command

On the other hand, they may have the following drawbacks when compared to Cloud Satellites offered in Earthly Cloud:
* No automatic updates
* Does not automatically sleep to save costs while idle
* Requires more knowledge and tuning to achieve good performance

## Getting Started
Launching a self-hosted satellite is as easy as running the public Docker image. A few environment variables are needed to link the satellite with your Earthly account. The satellite will use these values to automatically register with Earthly Cloud servers once it starts. Note that Earthly Cloud will never connect to your satellite directly. The host and port of your satellite can thus be private to your internal corporate network or VPC, if you wish.

Here is a minimal command to start a self-managed satellite using Docker:

```
docker run --privileged \
    -v satellite-cache:/tmp/earthly:rw \
    -p 8372:8372 \
    -e EARTHLY_TOKEN=GuFna*****nve7e \ 
    -e EARTHLY_ORG=my-org \
    -e SATELLITE_NAME=my-satellite \
    -e SATELLITE_HOST=153.65.8.0 \
  earthly/satellite:v0.8.3
```

The following environment variables are required:
* `EARTHLY_TOKEN` - a persistent auth token obtained by running `earthly account create-token <token-name>`
* `EARTHLY_ORG` - the name of your Earthly Organization created on [https://cloud.earthly.dev](https://cloud.earthly.dev)
* `SATELLITE_NAME` - a name chosen by the user to identify the satellite instance
* `SATELLITE_HOST` - Hostname or IP address for the Earthly CLI to connect. Note that users will select the satellite by name when running builds.

Using a volume for cache storage is recommended for preserving cache after the container is destroyed. See [Advanced Configuration](#advanced-configuration) below for more details and a list of optional environment variables that can be used to finetune your deployment. See also [Platform Specific Guides](#platform-specific-guides) for examples of deploying to platforms like Kubernetes or EC2.

The self-hosted satellite will not accept localhost as an address. For instructions on testing locally, please refer to [Running on Localhost](#running-on-localhost) below.


{% hint style='info' %}
##### Security Advisory
The `--privileged` flag is currently required. Privileged mode is used by some features of Earthly, including the `WITH DOCKER` command in Earthfiles. If using privileged mode is a concern, running satellites in a dedicated VM or separate Kubernetes cluster can provide better isolation compared to running the container directly in your production environment. A rootless version will be available in the future.
{% endhint %}

### Auto-discovery of Host _\*experimental\*_

When `SATELLITE_HOST` is left unset, the satellite will attempt to auto-discover its host address. Auto-discovery may not always work, and currently only supports environments using cloud-init (such as EC2).

## Connecting to a Self-Hosted Satellite

Once the satellite has finished registering, it can be selected and managed using typical earthly satellite commands from the CLI. 
For example, to select the satellite for use in subsequent builds:

```
earthly sat select my-satellite
```

Or, to invoke a build once on the satellite, without having it persistently selecting:

```
earthly --sat my-satellite +my-target
```

Note that `earthly` `v0.8.0` or later is required to connect to a self-hosted satellite.

## Managing Self-Hosted Satellites

A list of satellites in your org can be viewed via the satellite ls command. For example:

```
> earthly sat ls
NAME          PLATFORM     SIZE         VERSION  STATE
my-satellite  linux/amd64  self-hosted  v0.8.3   operational
...
```

The size field of a self-hosted satellite will be listed as `self-hosted`, as opposed to one of the sizes of Earthly Cloud’s managed satellites.

The state field will be either `operational` or `offline`. A self-hosted satellite automatically transitions to an `offline` state when it is gracefully terminated, or if no heartbeat has been received by Earthly’s servers for a while.

You can run `satellite rm` on a self-hosted satellite when it is in an offline state to remove it from your account permenantly.

## Platform-Specific Guides

### AWS EC2 (Recommended)

Deploying Satellites in a dedicated VM is the most secure method, since it isolates the satellite process from the rest of your infrastructure.

When launching on EC2, we recommend using the latest version of Amazon Linux 2023. The following cloud-init script can be configured when launching a new EC2 instance so that it automatically starts the satellite on boot.

```yaml
#cloud-config
runcmd:
  - sudo dnf update
  - sudo dnf install -y docker
  - sudo systemctl start docker.service
  - sudo systemctl enable docker.service
  - |-
      sudo docker run -d --privileged \
        --restart always \
        --name satellite \
        -p 8372:8372 \
        -v /earthly-cache:/tmp/earthly:rw \
        -e EARTHLY_TOKEN=GuFna*****nve7e \
        -e EARTHLY_ORG=my-org \
        -e SATELLITE_NAME=my-satellite \
        earthly/satellite:v0.8.3
```

Note that the `SATELLITE_HOST` variable is unset in this example so that the host is auto-discovered by the satellite when it starts. This should result in the instance’s private DNS being used as the host.

If the Earthly CLI is unable to connect to the satellite via the EC2’s private DNS, then `SATELLITE_HOST` should be provided in the `docker run` command with an alternate value.


### Kubernetes

Below is a basic example of how to start a self-hosted satellite as a Kubernetes Pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-satellite
  namespace: earthly-satellites
spec:
  volumes:
  - name: earthly-cache
    emptyDir: {} # or other volume type
  containers:
   - name:
   valueFrom:
     fieldRef:
       fieldPath: metadata.name
      image: earthly/satellite:v0.8.3
      securityContext:
        privileged: true
      ports:
      - containerPort: 8372
      volumeMounts:
      - mountPath: /tmp/earthly
        name: earthly-cache
       env:
       - name: EARTHLY_ORG
         value: my-org
       - name: EARTHLY_TOKEN
         value: u4a*************l92
       - name: SATELLITE_NAME
         valueFrom:
           fieldRef:
             fieldPath: metadata.name
       - name: SATELLITE_HOST
         valueFrom:
           fieldRef:
             fieldPath: status.podIP
```

This example uses the pod’s IP address (via [Downward API](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/)) as the `SATELLITE_HOST` value. The Earthly CLI must be able to reach the IP on its network.

If the Earthly CLI should connect via a different address, such as DNS, then this value should be provided as the `SATELLITE_HOST` instead.

Note: for best security, deploy satellites in a Kubernetes cluster that is separate from production.

## Advanced Configuration

### Cache Disk

Many deployments can benefit from using a dedicated volume as the satellite’s cache directory. Using a separate volume allows for persistent cache, even after the container is destroyed, such as when upgrading to a newer version. The satellite’s cache directory is located at `/tmp/Earthly` within the Docker container.

Here’s an example of how to attach the volume using the Docker command line:

```
docker run -v earthly-cache:/tmp/earthly:rw \
  ...
  earthly/satellite:v0.8.3
```

## Additional Environment Variables

The following environment variables can also be set to tweak the performance of a self-hosted satellite.

* `SATELLITE_PORT` - Sets an alternate port for Earthly CLI clients to connect to. Must be changed if exposing the satellite on a port other than 8372.
* `CACHE_SIZE_PCT` - The amount of disk to use for long-term cache storage. An integer from 1-100. Note that in-progress builds may consume additional disk space, so this value can result in "no space left on device" errors if set too high.
* `BUILDKIT_MAX_PARALLELISM` - The number of concurrent processes the satellite can run at once. Consider the number of cores available to the satellite when configuring this.
* `BUILDKIT_SESSION_TIMEOUT` - The max duration a single build can run for before timing out. Set to 24h by default.
* `CACHE_KEEP_DURATION` - How long idle cache will be retained on disk before being pruned (in seconds).
* `RUNNER_DISABLE_TLS` - Disable TLS on the satellite. Requires Earthly CLI to also disable TLS. Not recommended in most cases.
* `LOG_LEVEL` - The log level of the internal "runner" process within the satellite. Set to INFO by default.
* `BUILDKIT_DEBUG` - Enable buildkit debug-level logs. False by default.
* `EARTHLY_ADDITIONAL_BUILDKIT_CONFIG` - allows additional buildkit configs to be injected in TOML format.

## Running on Localhost

For testing purposes, you may want to try your self-hosted satellite on localhost. Satellites currently do not support using `localhost` or `127.0.0.1` as an address when supplied to `SATELLITE_HOST`, in part because Earthly CLI has reserved this address for use as its own local Buildkit container.

It’s still possible to test self-hosted satellites locally, however, by using an alternate entry in your `/etc/hosts` file that maps to `localhost`. For example, you can try adding this entry:

```
# Local test for Earthly
127.0.0.1	earthly.local
```

Starting your satellite with `SATELLITE_HOST` set to `earthly.local` should allow for a localhost connection using your Earthly CLI.

## Debugging

If you are having problems using or deploying your self-hosted satellite, please refer to the following tips or reach out to us through our community [Slack channel](https://earthly.dev/slack).

### Problem: Satellite is not listed in the output of `earthly satellite ls`

**Resolution:** Check the logs from the satellite’s Docker container. There may be a message containing the phrase `SATELLITE IS NOT REGISTERED`. This usually means there was a problem with the values supplied to the satellite’s run command. Check the error for additional context. Ensure the values provided for account token, earthly org, etc are correct.


### Problem: Satellite is not starting

**Resolution:** There may be some required environment variables missing or they may contain invalid values. Ensure all values are entered correctly. Check the container logs for more information.

### Problem: Earthly client is unable to connect to the satellite

**Resolution:** Check the address of the satellite. This is printed at the start of the build during the "Init" phase, or can be found via the satellite inspect command. Ensure the value here is as expected, and that the earthly client can reach the address on its network. If the port has been remapped from 8372, you may also need to change this via the `SATELLITE_PORT` environment variable.

### Problem: The satellite log says that it is running on port 9372

**Resolution:** The log message `"running server on [::]:9372"` can be misleading, however, the exposed port on the container is still 8372. Multiple processes are running inside the satellite container, including an earthly/buildkit process. This log message comes from the buildkit process, however, a separate process on port 8372 handles the incoming gRPC requests to the container.

### Problem: Satellite shows an `operational` state even though it is no longer running

**Resolution:** It is possible that the satellite did not terminate gracefully, and hence did not automatically deregister as it shutdown. You can run `earthly sat rm --force` to force-remove the satellite from the list, or wait some time for the satellite to automatically be removed. Earthly Cloud will automatically drop the satellite from the list if it detects no heartbeat messages for a while.
