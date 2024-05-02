# Managing Satellites

This page describes how to use [Earthly Satellites](../satellites.md).

## Prerequisites

In order to use Earthly Satellites, you must have an Earthly account and you must be a member of an Earthly organization. For more information, see the [Earthly Cloud overview](../overview.md) and the [Satellites page](../satellites.md).

If you are new to Earthly or to Earthly Cloud, get started by signing up for free at [Earthly Cloud](https://cloud.earthly.dev/login).

## Background

Earthly Satellites allow Earthly to execute builds in the cloud seamlessly. You execute build commands in the terminal, like you always have (for example, `earthly +build`), and Earthly takes care of running the build in the cloud in real time, instead of your local machine.

It uploads parts of your working directory, passes along any secrets, executes the build in the cloud while streaming the build log in real-time back to you, and then downloads the resulting build images and artifacts back to your computer.

For more information about how Earthly Satellites work, see the [Satellites page](../satellites.md).

## Using satellites

When you are added to an Earthly organization, you get access to its satellites. You can use the following command to select the organization:

```bash
earthly org select <org-name>
```

To view the satellites currently available in the organization, you can run:

```bash
earthly sat ls
```

### Selecting a satellite

In order to start using satellites, you can select one for use. Selecting a satellite causes Earthly to use that satellite for any builds from that point onwards.

```bash
earthly sat select <satellite-name>
```

Any build performed after selecting a satellite will be performed in the cloud on that satellite. You will notice that the output of the build contains information about the satellite that is being used:

```
$ earthly +build

 1. Init 🚀
————————————————————————————————————————————————————————————————————————————————

Note: the interactive debugger, interactive RUN commands, and inline caching do not yet work on Earthly Satellites.

The following feature flags are recommended for use with Satellites and will be auto-enabled:
  --new-platform, --use-registry-for-with-docker

           satellite | Connecting to core-test...
           satellite | ...Done
           satellite | Version github.com/earthly/buildkit v0.6.21 7a6f9e1ab2a3a3ddec5f9e612ef390af218a32bd
           satellite | Info: BuildKit version (v0.6.21) is different from Earthly version (prerelease)
           satellite | Platforms: linux/amd64 (native) linux/amd64/v2 linux/amd64/v3 linux/amd64/v4 linux/arm64 linux/riscv64 linux/ppc64le linux/s390x linux/386 linux/mips64le linux/mips64 linux/arm/v7 linux/arm/v6
           satellite | Utilization: 0 other builds, 0/12 op load
           satellite | GC stats: 9.0 GB cache, avg GC duration 275ms, all-time GC duration 2.754s, last GC duration 0s, last cleared 0 B

...
```

To go back to using your local machine for builds, instead of the satellite, you can unselect the satellite by running:

```bash
earthly sat unselect
```

### Specifying a satellite for one build only

If a satellite is not currently selected, you can still use it for a specific build by using the `--sat` flag or the corresponding `EARTHLY_SATELLITE` environment variable.

```bash
earthly --sat <satellite-name> +build
```

Conversely, if a satellite is currently selected, you can choose to use the local machine for a specific build using the `--no-sat` flag.

```bash
earthly --no-sat +build
```

### Managing performance

As satellites run the execution in the cloud, behind the scenes, they require upload of the current directory contents that may be needed as part of the build, and download of the results of the build. This is performed automatically by Earthly, however, if the file transfers are large and/or if the network bandwidth is low, the performance impact can be noticeable.

Oftentimes, you will find that running a build with the flag `--no-output` executes significantly faster. This flag disables downloading the build results from the satellite at the end of a build.

The `--no-output` flag can still be combined with `--push`, thus allowing Earthly Satellites to be used as a highly performant deployment tool.
