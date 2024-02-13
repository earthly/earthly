# Caching

Caching is at the heart of how Earthly works. It is what makes Earthly builds fast. This page provides a high-level understanding of the main concepts.

1. **When is Earthly fast** - in what situations Earthly will be fast
2. **How caching works in Earthfiles**
3. **How to share cache** between machines, or between runs in ephemeral CIs
4. **Managing cache** - how to reset it, how to configure its size, etc.

## When is Earthly fast

The word "build" can mean many things across many different contexts. When we say that it makes builds faster, we generally mean CI/CD builds.

Here are contexts in which Earthly does a particularly good job in, thanks to its caching:

* Making CI builds faster, especially in these circumstances
  * The CI performs many redundant tasks upfront, like installing dependencies and pulling container images.
  * The CI is a sandboxed CI, where no state is transferred over from one run to the next without explicitly uploading / downloading on each run (e.g. GitHub Actions, Circle CI)
  * Monorepos and Polyrepos: The CI builds multiple interconnected projects or sub-projects at a time
* Making local builds faster, especially in these circumstances
  * The build being executed is the CI build itself (and not just of the component you’re working on)
  * The build is complex, involving multiple projects or sub-projects at a time, possibly using multiple programming languages, where some of the projects could be rebuilt with a lot of cache shared with the CI or with teammates
  * Your internet connection is slow, and you need to perform a lot of image pushes and/or pulls

Here are examples where Earthly doesn’t improve performance:

* Local builds, when you’re iterating in a tight loop in a single programming language. Usually the tools of that programming language are already highly optimized for this use-case and often work better natively.
* CI builds, when the environment is shared between runs (unsafe), and you’re building programming languages with good built-in caching.
* CI builds, when the redundant parts of the build, like installing dependencies, are cached, AND the CI setup preserves the cache well, WITHOUT the need for downloading or uploading.
* CI builds that involve working with large files (i.e. >1 GB files), due to some internal transferring of files that Earthly relies on.

Now all this might be too complicated to remember, so here’s a simplified version. Earthly is:

* Almost always faster in CI, and especially faster in sandboxed CI environments.
* Usually not faster for local builds where you’re iterating in a single programming language in a tight loop.
* Often faster locally, when intending to run the same build as the CI.

The sections below go into more detail about how you are able to get faster builds with Earthly.

## Caching in Earthfiles

Main article: [Caching in Earthfiles](./caching-in-earthfiles.md)

There are three main ways in which Earthly performs caching of builds:

1. **Layer-based caching**. If an Earthfile command is run again, and the inputs to that command are the same, then the cache layer is reused. This allows Earthly to skip re-executing parts of the build that have not changed.
2. **Cache mounts**. Earthly allows you to mount directories into the build environment - either via [`RUN --mount type=cache`](../earthfile/earthfile.md#run), or via the [`CACHE`](../earthfile/earthfile.md#cache) command. These directories are persisted between runs, and can be used to store intermediate build files for incremental compilers, or dependencies that are downloaded from the internet.
3. **Auto-skip**. Earthly allows you to skip large parts of a build in certain situations via `earthly --auto-skip` (*experimental*) or `BUILD --auto-skip` (*coming-soon*). This is especially useful in monorepo setups, where you are building multiple projects at once, and only one of them has changed.

## Sharing Cache

The above capabilities can make your builds very fast. However, if you are using ephemeral CI runners, all of that valuable context can be lost between runs, resulting in poor build performance. Earthly's remote runners and caching via a registry capabilities solve this problem.

Since most CI platforms do not allow reusing state between runs efficiently, passing Earthly's cache via traditional CI cache constructs that rely on an upload and a download is too inefficient to be practical. Earthly's remote caching via a registry helps by optimizing what is uploaded and downloaded for maximum efficiency, although it does require experimentation to get right, and there are a number of limitations.

The most effective means of sharing cache between runs is to execute the Earthly builds remotely. This allows Earthly maintain the cache close to where it executes, thus being able to access it instantly without the need for an upload/download step. Because all Earthly builds are containerized, you still get the ephemeral nature of the CI runner, allowing for build repeatability, but you also get the benefits of a fast cache that is local to the execution environment.

Below is a comparison between remote runners, such as [Earthly Satellites](../cloud/satellites.md), and remote caching via a registry.

| Cache characteristic | Remote runners (e.g. Satellite) | Remote Cache via registry |
| --- | --- | --- |
| Storage location | Runner (e.g. Satellite) | A container registry of your choice |
| Proximity to compute | ✅ Same machine | ❌ Performing upload/download is required |
| Just works, no configuration necessary | ✅ Yes | ❌ Requires experimentation with the various settings |
| Concurrent access | ✅ Yes | 🟡 Concurrent read access only |
| Retains entire cache of the build | ✅ Yes | ❌ Usually no, due to prohibitive upload time |
| Retains cache for multiple historical builds | ✅ Yes | ❌ No, only one build retained |
| Cache mounts (`RUN --mount type=cache` and `CACHE`) included | ✅ Yes | ❌ No |

To read more, check out the [remote runners page](../remote-runners.md), and the [caching via a registry](./caching-via-registry.md).

## Managing Cache

For information on how to manage cache either locally, or on a remote runner, like a satellite, see the [Managing Cache guide](./managing-cache.md).
