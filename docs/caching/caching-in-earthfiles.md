# Caching in Earthfiles

Caching is at the heart of how Earthly works. This page will walk you through the key concepts of caching in Earthfiles.

There are three main ways in which Earthly performs caching of builds:

1. **Layer-based caching**. If an Earthfile command is run again, and the inputs to that command are the same, then the cache layer is reused. This allows Earthly to skip re-executing parts of the build that have not changed.
2. **Cache mounts**. Earthly allows you to mount directories into the build environment - either via [`RUN --mount type=cache`](../earthfile/earthfile.md#run), or via the [`CACHE`](../earthfile/earthfile.md#cache) command. These directories are persisted between runs, and can be used to store intermediate build files for incremental compilers, or dependencies that are downloaded from the internet.
3. **Auto-skip**. Earthly allows you to skip large parts of a build in certain situations via `earthly --auto-skip` or `BUILD --auto-skip`. This is especially useful in monorepo setups, where you are building multiple projects at once, and only one of them has changed.

## 1. Layer-based caching

Most commands in an Earthfile create a cache layer as part of the way they execute. You can think of a target in an Earthfile as a cake with multiple layers. If a layer's ingredients change, you need to redo the affected layer, plus any layer on top. Similarly, in an Earthfile target, if the input to a command is different (different ARG values, different source files being COPY'd, or the command itself is different), then Earthly can reuse the layers from a previous run up to that command, but it would have to re-execute that command and what follows after it.

If you happen to be familiar with Dockerfile layer caching, then layer caching in Earthly targets will be very familiar to you as it works the same way.

Earthly supports inheriting from other targets, copying artifacts that result from them, or simply issuing the build of another target. These various target cross-references result in a build graph underneath. Thus, one target could influence whether another target is executed - for example, if a source file changes and that results in rebuilding an artifact in `target1`, but then `target2` performs a `COPY` of that artifact, then at least part of `target2` will need to be re-executed too as a result. Earthly deals with all of this automatically.

Because of how layer caching works, it is best to organize builds in a manner that best utilizes the cache. A common strategy is to download and install dependencies early on in the build. Since the list of dependencies doesn't change very often, this expensive operation will usually be cached. To achieve this, it is important to copy the minimal amount of source files (usually just the file that defines what the dependencies are) before issuing the command that installs the dependencies.

Here is a practical example:

```Dockerfile
# Avoid
COPY . .
RUN go mod download
RUN go build ...
```

In the above example, changing the project's `README.md` or running `git fetch` might cause slow commands like `go mod download` to be re-executed.

Earthly uses `COPY` commands (among other things) to mark certain files as inputs to the build. If any file included in a `COPY` changes, then the build will continue from that `COPY` command onwards. For this reason, you want to be as specific as possible when including files in a `COPY` command. In some cases, you might even have to list files individually.

Here are some possible ways to improve the above example:

```Dockerfile
# Better
COPY go.mod go.sum ./*.go ./
RUN go mod download
RUN go build ...
```

The above is better, as it avoids reacting to changes in `.git` or to unrelated files, like `README.md`. However, this can be arranged even better, to avoid downloading all the dependencies on every `*.go` file change.

```Dockerfile
# Best
COPY go.mod go.sum ./
RUN go mod download
COPY ./*.go ./
RUN go build ...
```

In general, including the smallest set of input files as possible at every step will result in the best cache performance.

## 2. Cache mounts

Sometimes layer caching is not enough to properly express the best way to cache something. Cache mounts help complement layer caching, by allowing the contents of a directory to be reused across multiple builds. Cache mounts can be helpful in cases where the tool you're using to build within Earthly is able to leverage incremental caching on its own. Some package managers are able to do that for downloaded dependencies.

<!-- TODO: It would be nice to include a practical example from a programming language -->

Cache mounts can be used either via [`RUN --mount type=cache`](../earthfile/earthfile.md#run), or via the [`CACHE`](../earthfile/earthfile.md#cache) command. Although both allow you to define a path in your build environment where the cache directory would be mounted, there are a few important differences:

* Scope
  * `RUN --mount type=cache` only mounts the cache for that single `RUN` command.
  * `CACHE` mounts it for any `RUN` command that follows in the same target
* Final image
  * With `RUN --mount type=cache`, the contents of the cache are NOT persisted in the final image.
  * With `CACHE`, the contents of the cache are copied into in the final image, and also, as a result will be available to be read in targets inheriting from the original target
* Performance
  * `RUN --mount type=cache` is very performant as it does not require transferring contents at the end
  * `CACHE` can be slow in certain cases, if the contents are large, due to the need to copy the contents into the final image
* Consistency
  * `RUN --mount type=cache` is isolated to a single command, making it more difficult (but not impossible) to pass along files between steps via the cache
  * `CACHE` is available to all commands in the target, making it easier to pass along files between steps via the cache, and thus also easier to run into race conditions, if a parallel build changes the contents of the cache in unexpected ways

Cache mounts, by default, are only available within the same target. So if both `target1` and `target2` define `RUN --mount type=cache,target=/my-cache`, the contents would not be shared. If you would like to share the contents, you can use the `id` option. Setting the `id` makes the cache mount global, allowing any target to access the same contents, as long as they both use the same `id`: `RUN --mount type=cache,id=my-cache-id,target=/my-cache`.

Parallel builds using the same cache mount (or the same build where the mount is used in multiple targets) pose another aspect to be aware of: accessing the cache mount concurrently. By default, sharing is set to `locked` - meaning that parallel executions will wait for each other to complete, thus allowing access by one process at a time. While this is the safest option, it is also the slowest. Keep in mind that this will limit your build parallelism significantly if you overuse global cache mounts. Other possible options are `shared` (allows concurrent access), or `private` (if a parallel execution occurs, a new empty mount is created).

### Drawbacks of cache mounts

Cache mounts can be a versatile tool for controlling caching in ways that layer caching cannot. There are, however, important limitations to understand.

The most important limitation to be aware of is that reusing state from a previous run can be a source of build inconsistency. A test passing just because it starts off with the right contents in cache could later result in deploying a broken application to production.

Another limitation is that cache mounts are not great for passing files from one build step to another. This is because a parallel build could interfere with the cache between steps in ways that are difficult to debug. Be especially mindful that builds from different development branches might interact with each other unexpectedly in this situation. It is therefore best to avoid using cache mounts as a mechanism to pass along information. It is best to extract the result of an operation out of the cache mount within the same operation, to ensure that the cache is locked during this time.

Finally, another important limitation is the fact that cache mounts can grow in size indefinitely. While Earthly does garbage-collect layers and cache mounts on a least-recently-used basis, a cache mount that is used frequently could grow more than expected. In such situations, you should consider managing the lifecycle of the cache contents yourself, by removing unused files from it every few runs. A good place for such cleanup operations is within the same layer (same `RUN` command) that uses the contents, at the end.

## 3. Auto-skip

Auto-skip is a feature that allows Earthly to skip large parts of a build in certain situations. This is especially useful in monorepo setups, where you are building multiple projects at once, and only one of them has changed.

Unlike layer caching and cache mounts (which store cache local to the runner), auto-skip is a global cache stored in a cloud database. In order to use auto-skip, you will need an [Earthly Cloud](../cloud/overview.md) account.

Auto-skip can be enabled for either an entire run, via `earthly --auto-skip` (*experimental*), or for a specific target, via `BUILD --auto-skip` (*coming soon*).

Unlike layer caching, auto-skip is an all-or-nothing type of cache. Either the entire target is skipped, or none of it is. This is because Earthly does not know which parts of the target are affected by the change. If auto-skip does not deem the run to be skipped, then Earthly will fallback to the other forms of caching to run the build as efficiently as possible.

### When auto-skip is not supported

As auto-skip relies on statically analyzing the structure of the build upfront, including the inter-dependencies between targets across multiple Earthfiles, it is not always possible to use it. If a target being involved has a dynamic name that would only be known at run-time, then auto-skip would have no way of knowing it upfront. In such cases, the build fails with an error message when `--auto-skip` is enabled.

#### Static inference of ARG values

For basic `ARG` operations, auto-skip is able to infer the value of the `ARG` statically, and therefore, it is able to support it. Here is a practical example.

```
# Supported
ARG MY_ARG=foo
BUILD $MY_ARG

# Not supported
ARG MY_ARG=$(cat ./file)
BUILD $MY_ARG
```

In the first case, the value of `MY_ARG` is known statically as its value can be propagated by the auto-skip algorithm. In the second case, the value of `MY_ARG` is not known statically as it depends on a file in the build environment. In such a case, auto-skip is not supported. Note that defining such dynamic `ARG`s is generally allowed, however, as long as the value of the `ARG` is not used in a way that would prevent auto-skip from working.

Similarly, the auto-skip algorithm is able to propagate `ARG`s across targets, as long as the value of the `ARG` is known statically. Here is a practical example:

```
# Supported
ARG MY_ARG=foo
BUILD +target --arg=$MY_ARG

# Might not be supported (depending on how the target uses the arg)
ARG MY_ARG=$(cat ./file)
BUILD +target --arg=$MY_ARG
```

#### Static inference of conditions

`IF` statements are generally supported by auto-skip, however there is a difference in behavior, depending on whether the outcome of the condition can be inferred statically. If the outcome of the condition can be inferred statically, then auto-skip uses the correct `IF` or `ELSE` block for analysis. If the outcome of the condition cannot be inferred statically, then the auto-skip algorithm will analyze both `IF` and `ELSE` blocks, and will only skip the target if both blocks are skipped.

Here is a practical example:

```
# Supported and efficient (only +target2 is analyzed)
ARG MY_ARG=bar
IF [ $MY_ARG = "foo" ]
  BUILD +target1
ELSE
  BUILD +target2
END

# Supported but inefficient (both +target1 and +target2 are analyzed)
ARG MY_ARG=$(cat ./file)
IF [ $MY_ARG = "foo" ]
  BUILD +target1
ELSE
  BUILD +target2
END

# Supported but inefficient (both +target1 and +target2 are analyzed)
IF grep ./file -e "foo" 
  BUILD +target1
ELSE
  BUILD +target2
END
```

#### Unsupported Earthfile features in auto-skip

Here is a list of unsupported features when `--auto-skip` is enabled:

* Dynamic target names, such as `BUILD $MY_TARGET`, `FROM $MY_TARGET`, or `BUILD $(...)`, unless the target name can be inferred statically.
* Dynamic `COPY` commands, such as `COPY $MY_FILE .`, or `COPY $(...) .`, unless the source can be inferred statically.
* Remote references (such as `BUILD github.com/foo/bar+my-target`), unless the remote reference is pinned to a specific SHA, or to an explicit tag expressed as a `tags/...` git reference. For example, `BUILD github.com/foo/bar:tags/v1.0.0+my-target` is supported, but `BUILD github.com/foo/bar:v1.0.0+my-target` is not.
* Remote imports, unless the remote reference is pinned to a specific SHA, or to an explicit tag expressed as a `tags/...` git reference. For example, `IMPORT github.com/foo/bar:tags/v1.0.0` is supported, but `IMPORT github.com/foo/bar:v1.0.0` is not.
* `GIT CLONE`, unless the remote reference is pinned to a specific SHA, or to an explicit tag expressed as a `tags/...` git reference.
* `FOR` loops, unless the list being iterated can be inferred statically.

### Auto-skipping `RUN --no-cache` and `RUN --push`

Unlike layer caching, auto-skip may also skip `RUN --no-cache` and `RUN --push` commands. This can be useful in situations when you would like to skip a deployment, if nothing has changed. Please note that rollbacks may be unintentionally skipped, if attempting to deploy on older version of the codebase that had been previously deployed. In such cases, you can either (a) remove the `--auto-skip` flag to force the deployment to occur, (b) [clear the auto-skip cache](./managing-cache.md#auto-skip-cache), or (c) use roll-forward practices (reverting problematic code changes, and re-deploying as a net-new version).

## Disabling caching

In certain situations, you might want to disable caching either for a specific command, or for the entire build.

To disable layer caching, you can use the `--no-cache` flag. For example, `RUN --no-cache echo "Hello"` will always execute the `echo` command, even if the `RUN` command was executed before with the same arguments. Note that this does not disable cache mounts, or auto-skip. A `RUN --no-cache` command can still be skipped by auto-skip.

To disable layer caching and mount caching for an entire run, you can use `earthly --no-cache +my-target`.

Another way to disable layer caching is to use the `RUN --push` flag. This flag is useful when you want to perform an operation with external effects (e.g. deploying to production). By default Earthly does not run `--push` commands unless the `--push` flag is also specified when invoking Earthly itself (`earthly --push +my-target`). `RUN --push` commands are never cached.

To disable auto-skip, simply remove the `--auto-skip` flag.

## Troubleshooting and gotchas

Debugging caching issues can be tricky. Here are some common issues that you might face and how to resolve them.

### Cache size

If the configured cache size is too small, then Earthly might garbage-collect cached layers more often than you might expect. This can manifest in builds randomly not using cache for certain layers. Usually it is the biggest layers that suffer from this (and oftentimes the biggest layers are the most expensive to recreate). This problem is most common on Mac and Windows, where Docker uses a VM with limited disk size. To resolve this, either configure a larger cache size if you are running local builds, or launch a larger Satellite if you are using remote builds via Earthly Satellites. For more information see the [managing cache page](./managing-cache.md).

### ARGs

In Earthly, like in Dockerfiles, ARGs declared in Earthfiles also behave as environment variables within the target they are declared in. This means that if you declare an ARG, and then use it in a `RUN` command, then the `RUN` command will be invalidated if the ARG changes. This is sometimes not very obvious, especially if you are not actually using the value of that ARG.

For this reason, it is best to declare ARGs as late as possible within the target they are used in, and try to avoid declaring `--global` ARGs as much as possible. If an ARG is not yet declared, it will not influence the cache state of a layer, allowing for more cache hits. Limiting the scope of ARGs as much as possible will yield better cache performance.

Watch out especially for ARGs that change often, such as the built-in ARG `EARTHLY_GIT_HASH`. Declaring this ARG as late as possible in the build will cause less cache misses.

### Secrets

Note that secrets, unlike ARGs, do NOT contribute to the cache state of a layer. This means that if you use a secret in a `RUN` command, and the secret changes, the `RUN` command will not be invalidated.

### Force a build step to always cache

If you have already optimized your cache by maximizing its size, declaring arguments as late as possible, and implementing the other recommendations provided here, but you still encounter performance bottlenecks due to computationally intensive tasks being evicted from the cache, consider employing `SAVE IMAGE` commands at strategic points. These images can serve as manual caches and can improve efficiency at the cost of simplicity. For additional details, refer to the [Best Practices](../guides/best-practices.md#use-save-image-to-always-cache) section.

### Debugging tips

If you are experiencing caching issues and have ruled out the above common situations, we would love to hear from you. Please open an issue in the [Earthly GitHub repository](https://github.com/earthly/earthly).
