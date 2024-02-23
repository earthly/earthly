# Remote runners

Earthly supports running builds remotely via remote runners. Remote runners allow you to benefit from sharing the cache with other users of that remote runner. This is especially useful in CI environments where you want to share the cache between runs.

## Benefits

Typical use cases for remote runners include:

* Speeding up CI builds in sandboxed CI environments such as GitHub Actions, GitLab, CircleCI, and others. Most CI build times are improved by a factor of 2-20X via Satellites.
* Executing builds on AMD64/Intel architecture natively when working from an Apple Silicon machine (Apple M1/M2) and vice versa.
* Sharing compute and cache with coworkers or with the CI.
* Benefiting from high-bandwidth internet access from the cloud, thus allowing for fast downloads of dependencies and fast pushes for deployments. This is particularly useful if operating from a location with slow internet.
* Using Earthly from environments where privileged access or docker-in-docker are not supported, or from environments where Docker is not installed. Note that the remote runner itself still requires privileged access.

## How remote runners work

When using remote runners, even though the build executes remotely, the following pieces of functionality are still available:

* Build logs are streamed to your local machine in real-time, just as if you were running the build locally
* Outputs (images and artifacts) resulting from the build, if any, are transferred back to your local machine
* Commands under `LOCALLY` execute on your local machine
* Secrets available locally, including Docker/Podman credentials are passed to the satellite whenever needed by the build
* Any images to be pushed are pushed directly from the satellite, using any Docker/Podman credentials available on the local system.

## Get started

To get started with free remote runners managed by Earthly, check out [Earthly Satellites](cloud/satellites.md).

To get started with self-hosted runners, see the [Self-Hosted Satellites Guide](cloud/satellites/self-hosted.md).

If your use case cannot tolerate a cloud-based control plane, then self-hosting a remote Buildkit is the best approach. Remote Buildkit has less features, is less secure, and is more difficult to deploy than Self-Hosted Satellites (see diagram below for comparison). To get started self-hosting Buildkit, see the [remote buildkit page](ci-integration/remote-buildkit.md).

### Types of Remote Runners

Below is a comparison of the different features available with each kind of remote runner.

| Feature | Cloud Satellites | Self-Hosted Satellites | Remote Buildkit |
| --- | --- | --- | --- |
| Managed By | Earthly | You | You |
| Cache Persistence | ✅ Yes | 🟡 Needs configuration | 🟡 Needs configuration |
| Cloud Control-Plane | ✅ Yes | ✅ Yes | ❌ No |
| Managed TLS Certificates |  ✅ Yes | ✅ Yes | ❌ No |
| Auto-Sleep | ✅ Yes | ❌ No | ❌ No | 
| Auto-Updates | ✅ Yes | ❌ No | ❌ No |
