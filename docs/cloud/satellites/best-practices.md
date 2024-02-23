# Satellites Best Practices

Earthly Satellites are [remote build runners](../../remote-runners.md) hosted by Earthly, which enables building in the cloud with instantly available cache. This best practices guide aims to help users optimize their workflow and leverage the full potential of Earthly Satellites.

## Size Satellites for Builds

Consider the following when picking a satellite size:

### 1. Resource Requirements of a Single Build

Understand the resource demands of your build, including memory usage, CPU requirements, and required cache size on disk.

For example, builds utilizing computationally intensive compilation steps, or large integration tests may require additional CPU and memory. Builds that download many gigabytes of dependencies may require a large amount of cache on disk.

{% hint style='info' %}
##### Tip
Use the `earthly --exec-stats` flag to view CPU and memory usage during a build.
{% endhint %}

### 2. Concurrent Build Capacity

Take into account the number of simultaneous builds that are run by your development team. The Satellite size should be chosen to handle the expected concurrency level effectively. 

One metric to consider is the number of commits made per hour that will trigger builds on the satellite. For example, an xsmall or small satellites might handle a couple of concurrent builds; however, larger teams contributing more frequently may require a larger instance to accommodate the extra load.

## Splitting Builds and Using More Satellites

Large projects will often benefit from using multiple satellites in a single CI run, balancing the load across different machines in parallel. With this strategy, it’s best to optimize for cache reuse, by running the same target on the same satellite rather than alternating runs of the target across different satellites.

As an example, consider an Earthfile that builds multiple libraries:

```
build-all:
  BUILD +library-1
  BUILD +library-2
```

If the `+build-all` target is too intensive for a single satellite, dedicated satellites can be created for each library:

```
earthly --satellite sat-1 +library-1
earthly --satellite sat-2 +library-2
```

This strategy is easiest to implement when there are distinct components that can be split apart. This is natural in monorepos, but the same considerations on cache reuse should be made in any case.

## Cancelling Builds from Stale Commits

When running Earthly on pull requests, it is common for new commits to be pushed before a previous build has finished. Since the previous results may no longer be useful, it’s good practice to cancel the previous build – freeing up resources on the satellite for the current build.

Many CI systems allow this behavior to be configured. Github Actions, for example, can be configured with:

```yaml
name: Github Actions CI

concurrency: 
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true

jobs:
  ... 
```

## Determining when Satellites are Overloaded

A satellite may become sluggish, or even fail a build (in extreme cases) if it has taken on more jobs than it has the resources to handle. If you are experiencing a drop in performance, consider using a larger instance, or splitting up your build across multiple satellites. Note that in some cases, performance issues can be the result of a problem in the build script or external factors such as network operations.

There are some metrics printed by the Earthly CLI that can be used to gauge the performance and overall health of a satellite. These metrics are printed at the start of the build during the "Init" phase, when running the satellite inspect command, or sometimes printed midway through the build if Earthly detects a performance problem. Here is an example of the metrics printed during Init:

```
 Init 🚀
————————————————————————————————————————————————————————————————————————————————

           satellite | Connecting to my-satellite...
           satellite | ...Done
           satellite | Version github.com/earthly/buildkit v0.7.23 086f60eecf5a2261bbd53299614f88b3def12746
           satellite | Platforms: linux/amd64 (native) linux/amd64/v2 linux/amd64/v3 linux/amd64/v4 linux/arm64 linux/riscv64 linux/ppc64 linux/ppc64le linux/s390x linux/386 linux/mips64le linux/mips64 linux/arm/v7 linux/arm/v6
           satellite | Utilization: 0 other builds, 0/12 op load
           satellite | GC stats: 2.1 MB cache, avg GC duration 0s, all-time GC duration 23ms, last GC duration 0s, last cleared 0 B
```

In this output, there are a couple of metrics which are worth keeping an eye on:
* `Utilization` – The number of other concurrent builds and op load usage. The op load indicates the number of concurrent operations being executed, compared to the max allowed by the satellite, where an operation is a step within an Earthfile. When the op load is at the max, operations may begin to queue, resulting in slower performance.
* `GC Stats` – The amount of cache on disk, and time spent doing garbage collection. Compare this to the amount of [disk available per the satellite’s size](https://earthly.dev/pricing). The cache typically hovers around 50% of the total disk size, however, if the cache size is close to the max, then aggressive garbage collection can cause degraded performance.

On top of the metrics mentioned above, these other symptoms may indicate an overloaded satellite:
* "No space left on device" errors
* "Out of memory" errors
* Op load metrics printing during the build
* Noticeably long delays between steps in the build
* New builds failing to connect to the satellite
* Unexpected build failures or crashes during a build

If you are unsure that your satellite is overloaded, please reach out via [email](mailto:support@earthly.dev) or our [Community Slack channel](https://earthly.dev/slack). An Earthly team member will be happy to investigate your instance and provide advice.

## Limiting Unnecessary Outputs

By default, Earthly outputs artifacts or images locally at the end of the build. When running in CI, however, these artifacts may not actually be needed locally. For example, when your earthly build pushes an image to an external registry, but you do not actually need a copy of that image to be downloaded back to your CI runner.

Using the [`--ci` flag](../../../guides/best-practices.md#use-ci-when-running-in-ci) can result in better performance, especially when using satellites, since the output will run over the network.
