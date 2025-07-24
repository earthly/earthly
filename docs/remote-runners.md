# Remote runners

Earthly supports running builds remotely via remote runners. Remote runners allow you to benefit from sharing the cache with other users of that remote runner. This is especially useful in CI environments where you want to share the cache between runs.

## Benefits

Typical use cases for remote runners include:

* Speeding up CI builds in sandboxed CI environments such as GitHub Actions, GitLab, CircleCI, and others. Most CI build times are improved by a factor of 2-20X.
* Executing builds on AMD64/Intel architecture natively when working from an Apple Silicon machine (Apple M1/M2) and vice versa.
* Sharing compute and cache with coworkers or with the CI.
* Benefiting from high-bandwidth internet access from the cloud, thus allowing for fast downloads of dependencies and fast pushes for deployments. This is particularly useful if operating from a location with slow internet.
* Using EarthBuild from environments where privileged access or docker-in-docker are not supported, or from environments where Docker is not installed. Note that the remote runner itself still requires privileged access.

## How remote runners work

When using remote runners, even though the build executes remotely, the following pieces of functionality are still available:

* Build logs are streamed to your local machine in real-time, just as if you were running the build locally
* Outputs (images and artifacts) resulting from the build, if any, are transferred back to your local machine
* Commands under `LOCALLY` execute on your local machine
* Secrets available locally, including Docker/Podman credentials are passed to the satellite whenever needed by the build
* Any images to be pushed are pushed directly from the satellite, using any Docker/Podman credentials available on the local system.

## Get started

Self-hosting a remote BuildKit is the best approach. To get started self-hosting BuildKit, see the [remote BuildKit page](ci-integration/remote-buildkit.md).

### Types of Remote Runners

Below is a comparison of the different features available with each kind of remote runner.

| Feature                                                                     | Remote BuildKit        |
|-----------------------------------------------------------------------------|------------------------|
| Managed By                                                                  | You                    |
| Cache Persistence                                                           | 🟡 Needs configuration |
| Cloud Control-Plane                                                         | ❌ No                   |
| Managed TLS Certificates                                                    | ❌ No                   |
| Auto-Sleep                                                                  | ❌ No                   |
| Auto-Updates                                                                | ❌ No                   |
| [GitHub Actions integration](cloud/satellites/gha-runners.md)               | ❌ No                   |
