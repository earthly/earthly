# Caching via a registry (advanced)

{% hint style='danger' %}
##### Important
**Remote caching is an advanced feature that requires significant experimentation to get right**. Remote caching is generally not recommended unless remote runners are not possible in your setup. Remote caching does not cache everything in your build, and is very difficult to debug when it doesn't work. Remote caching also does not allow concurrent writes, and can be tricky to use in a mix of main vs PR CI builds.

[Remote runners](../remote-runners.md) (either self-managed, or [Earthly Satellites](../cloud/satellites.md)) are faster, and easier to set up because the cache does not need any transferring between runs (allowing for the entire build to be cached efficiently), cache mounts are correctly persisted, and concurrent access is built-in.

For a comparison between remote runners and remote caching, see the [Caching page](./caching.md).
{% endhint %}

Earthly has the ability to share cache between different isolated CI runs and even with developers via remote caching. This page goes through the available features, common use-cases and situations where remote caching is most useful.

Remote caching is made possible by storing intermediate steps of a build in a cloud-based Docker registry. This cache can then be downloaded on another machine in order to skip common parts.

## When to use remote caching

Remote caching should be used in the following circumstances:

* When remote runners are not possible in your setup
* When preserving cache mounts is not needed
* When caching only part of the build is sufficient
* When you want to cache compute-intensive steps, but not download-intensive steps
* When concurrent writing to the cache is not important
* When you are willing to experiment with the various settings to get the best performance

## When NOT to use remote caching

Remote caching should NOT be used in the following circumstances:

* When remote runners are possible in your setup
* When preserving cache mounts is needed
* When caching the entire build is useful
* When you want to cache download-intensive steps too, not just compute-intensive steps
* If you want a simple setup that just works and does not require significant debugging

## Types of remote cache

Earthly makes available two types of remote caching:

* [Inline cache](#inline-cache)
* [Explicit cache](#explicit-cache-advanced) (advanced)

For a summary of the differences see [comparison between inline and explicit cache](#comparison-between-inline-and-explicit-cache).

### Compatibility with major registry providers

Not all registries support the needed manifest formats to allow the usage of each kind of cache. Here is a compatibility matrix for many popular registries:

| Registry                  | Supports Inline Cache | Supports Explicit Cache | Notes                                                                               |
|---------------------------|:---------------------:|:-----------------------:|-------------------------------------------------------------------------------------|
| Google GCR                |           ✅           |            ❌            |                                                                                     |
| AWS ECR                   |           ✅           |            ✅            |                                                                                     |
| Google Artifact Registry  |           ✅           |            ✅            |                                                                                     |
| Azure ACR                 |           ✅           |            ✅            |                                                                                     |
| Docker Hub                |           ✅           |            ✅            |                                                                                     |
| GitHub Container Registry |           ✅           |            ✅            |                                                                                     |
| GitLab Container Registry |           ✅           |            ✅            |                                                                                     |
| Self-Hosted  `registry:2` |           ✅           |            ✅            |                                                                                     |
| JFrog Artifactory         |           ✅           |            ✅            | Only versions > 7.31.10 work, due to https://www.jfrog.com/jira/browse/RTFACT-26179 |

### Inline cache

Inline caching is the easiest to configure. It essentially makes use of any image already being pushed to the registry and adds some very small metadata (a few KiB) as part of its configuration about how Earthly is able to reuse that image for future runs.

The key benefit of this approach is that you get the upload for free if you anyway push images to the registry.

#### How to use inline caching

In order to enable inline caching, add `--use-inline-cache --save-inline-cache` in your invocation of `earthly` when you `--push`, or `--use-inline-cache` otherwise. When `--use-inline-cache --save-inline-cache --push` flags are specified, the use of the cache will be read-write.

Use a read-only inline cache (on a developer's computer or in PR builds):

```bash
earthly --use-inline-cache +some-target
```

Use a read-write inline cache (typically in master/main branch builds):

```bash
earthly --use-inline-cache --save-inline-cache --push +some-target
```

The options mentioned above are also available as environment variables. See [Earthly command reference](../earthly-command/earthly-command.md) for more information.

The way this works underneath is that Earthly uses `SAVE IMAGE --push` declarations as source and destination for any inline cache.

In case different Docker tags are used in branch or PR builds, it is possible to use additional cache sources via [`SAVE IMAGE --cache-from=...`](../earthfile/earthfile.md#save-image). This may be useful so that PR builds are able to use the main branch cache. Here is a simple example:

```Dockerfile
FROM ...
...
ARG BRANCH=master
SAVE IMAGE --cache-from=mycompany/myimage:master --push mycompany/myimage:$BRANCH
```

{% hint style='info' %}
The `--use-inline-cache` flag is required to enable importing existing caches, and the `--save-inline-cache` flag is required to enable exporting images to the remote cache.

Since `VERSION 0.6` the inline cache is only exported to images [that are connected to the initial target through a chain of BUILD commands](https://docs.earthly.dev/docs/earthfile#build).
{% endhint %}

#### Optimizing inline cache performance

Inline caching is very easy to use, however it can also turn out to be ineffective for some builds. One limitation is that only the layers that end up in uploaded images are actually used. Certain intermediate layers (e.g. targets only used for compiling binaries) will not exist.

If you find that certain steps could benefit from being cached but are not, you may consider creating additional images for those steps specifically. All you need to do is add the following at the end. Use a Docker tag that is not used for anything else.

```Dockerfile
SAVE IMAGE --push <docker-tag>
```

Note however that adding more images to the build results in additional time spent uploading them. Disregard the performance of the very first upload, as a fresh push is always less performant because there is no commonality with any previous run.

#### Example of using inline caching

Good example uses of inline caching are the Earthly [C++](https://github.com/earthly/earthly/tree/main/examples/cpp) and [Scala](https://github.com/earthly/earthly/tree/main/examples/scala) samples.

In the C++ case, a lot of computation is saved as a result of the `apt-get install` command. Reusing the cache improves performance by a factor of 4X.

In the Scala case, time is saved from processing the dependencies, resulting in a 3X performance improvement.

In both cases, a major benefit is that we are anyway pushing the images to the cloud via the `SAVE IMAGE --push` commands. So there is no performance penalty on the cache upload side. The command that would be used in the CI to execute the builds together with inline caching is

```bash
earthly --use-inline-cache --save-inline-cache --push +docker
```

### Explicit cache (advanced)

Explicit caching requires that you dedicate a Docker tag specifically for cache storage. Unlike inline caching, this tag is not meant to be used for anything else. For this reason, uploading the cache is an added step that takes additional time.

#### How to use explicit caching

To enable explicit caching, use the flag `--remote-cache=...` to specify the Docker tag to use as cache. Make sure that this Docker tag is not used for anything else (e.g. DO NOT use `myimage:latest`, in case `latest` is used in a critical workflow).

For example, if the Docker tag used for explicit caching is `mycompany/myimage:cache`, then the flag can be used as follows.

In CI, read-only inline cache (typically in PR builds):

```bash
earthly --ci --remote-cache=mycompany/myimage:cache +some-target
```

In CI, read-write inline cache (typically in master/main branch builds):

```bash
earthly --ci --remote-cache=mycompany/myimage:cache --push +some-target
```

On developer's computer (optional):

```bash
earthly --remote-cache=mycompany/myimage:cache +some-target
```

##### Additional Image Attributes

Some Registry Providers require setting certain attributes so that they can be used for explicit caching.  
These attributes can be set in the ``--remote-cache` flag, for example:

```bash
earthly --remote-cache=harbor.company.team/earthly-cache/test:cache,image-manifest=true,oci-mediatypes=true +example
```

The options mentioned above are also available as environment variables. See [Earthly command reference](../earthly-command/earthly-command.md) for more information.

{% hint style='info' %}
##### Note

If a project has multiple CI pipelines or `earthly` invocations, it is recommended to use different `--remote-cache` Docker tags for each pipeline or invocation. This will prevent the cache from being overwritten in ways in which it makes it less effective.
{% endhint %}

{% hint style='info' %}
##### Note

It is currently not possible to push both inline and explicit caches in a single run.
{% endhint %}

#### Optimizing explicit cache performance (advanced)

Explicit caching works by storing a cache containing all the layers of the final target, plus any target containing `SAVE IMAGE --push ...`. If additional targets need to be added as part of the cache, it is possible to add `SAVE IMAGE --cache-hint` (no Docker tag necessary) at the end, in order to mark them for explicit caching.

```Dockerfile
deps:
  COPY go.mod go.sum ./
  RUN go mod download
  SAVE IMAGE --cache-hint
```

Making use of explicit caching effectively may not always be possible. Sometimes the overhead of uploading and redownloading the cache defeats the purpose of gaining build performance. Oftentimes, multiple iterations of trial-and-error need to be attempted in order to optimize its effectiveness. Keep in mind that caching compute-heavy targets is more likely to yield results, rather than download-heavy targets.

As an additional setting available, Earthly can be instructed to save all intermediary steps as part of the explicit cache. The setting `--max-remote-cache` can be used to enable this. Note that this results in large uploads and is usually not very effective. An example where this feature is useful, however, is when you would like to optimize CI run times in PRs, and are willing to sacrifice CI run times in default branch builds. This can be achieved by enabling `--push` and `--max-remote-cache` on the default branch builds only.

#### Example of using explicit caching

A good example of using explicit caching is this [integration test example](https://github.com/earthly/earthly/tree/main/examples/integration-test). The target `+project-files` is perfect for introducing a cache hint via `SAVE IMAGE --cache-hint`. The processing that takes place as part of installing Scala and compiling the dependencies is sufficiently compute-intensive to save ~2 min from the total build time in CI. In addition, these dependencies change rarely enough that the cache can be utilized consistently.

A typical invocation of the build to make use of the explicit cache:

```bash
earthly --ci --remote-cache=mycompany/integration-example:cache --push +all
```

### Comparison between inline and explicit cache

Inline and explicit caching have similar traits, but they also have a number of fundamental differences.

The key similarity is that both types of caches make use of Docker tags being pushed to an image registry in order to store the cache.

The most important difference is that inline caching relies on image uploads that are already being made. And as such, the cache may be split across multiple separate images. Every `SAVE IMAGE --push` command adds more cacheable targets in the form of separate images. However, in the case of explicit caching, the entire cache is stored as part of a single Docker tag and every `SAVE IMAGE --cache-hint` command adds more cacheable targets within the image. This final image containing all the explicit cache cannot be used for anything else. So as a user, you incur the performance cost of both the upload and the subsequent download.

Below is a summary of the different characteristics of each type of cache.

#### Key takeaways for inline caching

* Cache is embedded within images that are already being pushed. No new layers are added to the images, only a few KiB of metadata.
* Very easy to use (just add `--ci` to your `earthly` invocations in CI)
* It is usually effective right away, with little modifications
* Typically you incur the performance cost only for the subsequent download. Upload is for free if you are pushing images anyway
* By default, caches only the images being pushed
* You can add more cache via additional `SAVE IMAGE --push <docker-tag>` commands

#### Key takeaways for explicit caching

* Cache is uploaded as part of a new Docker tag that should not be used for anything else
* The only available choice if no images are already pushed during the build
* More control over what is being cached and what is not. However it often requires some level of experimentation to get right.
* Incur the performance cost for both the upload and the download
* By default, caches only the layers of the target being built, and not of any other referenced targets
* You can cache additional targets by adding `SAVE IMAGE --cache-hint` commands

## When to use remote caching

There are several situations where remote caching can provide a significant performance boost. The following are only a few examples of how to get a feel for its usefulness.

### Compute-heavy vs Download-heavy

In general remote caching is very useful when there is a significant computation overhead during the execution of your build. Assuming that the inputs of that computation do not change regularly, then remote caching could be a good candidate. If a time-consuming operation, however, is not compute-heavy, but rather download-heavy, then remote caching may not be as effective (it's one download versus another).

As an example of this distinction, consider the use of the `apk` tool shipped in `alpine` images. Installing packages via `apk` is download-heavy, but usually not very compute-heavy, and so using remote caching to offset `apk` download times might not be as effective. On the other hand, consider `apt-get` tool shipped in `ubuntu` images. Besides performing downloads, `apt-get` also performs additional post-download steps which tend to be compute-intensive. For this reason, remote caching is usually very effective here.

Similarly to the comparison between `apk` and `apt-get`, similar remarks can be made about the various language-specific dependency management tools. Some will be pure download-based (e.g. `go mod download`), while others will be a mix of download and computation (.e.g `sbt`).

### An intermediate result is small and doesn't change much

An area where remote caching is particularly impactful are cases where a rare-changing prerequisite downloads many dependencies and/or performs intensive computation, but the end result is relatively small (e.g. a single binary). Passing this prerequisite over the wire as part of the remote caching is very fast (especially if the downloads required to generate it are not used anywhere else), whereas regenerating it requires a lot of work.

### Monorepo and Polyrepo setups

An excellent example of the above are typical inter-project dependencies. Regardless of whether your layout is a monorepo or a polyrepo, if projects reference artifacts or images from each other, then whatever tools used to generate those artifacts or images are usually not required across projects. In such cases it is possible to prevent entire target trees of downloads and computation and simply download the final result using the remote cache.

A simple way to visualize this use-case is comparing the performance of a build that takes place behind a `FROM +some-target` instruction versus just using the previously built image directly. If `+some-target` has a `SAVE IMAGE --push myimage:latest` instruction, then the performance becomes almost the same to using `FROM myimage:latest` directly.

### CIs that operate in a sandbox

Modern CIs execute in a sandbox. They start with a blank slate and need to download and regenerate everything from scratch. Examples of such CIs: GitHub Actions, Circle CI, DroneCI, GitLab CI. Such CIs benefit greatly from being able to share precomputed steps between runs.

If, however, you are using a CI which reuses the same environment (e.g. Jenkins, BuildKite - depending on how they are configured), then simply relying on the local cache is enough.

### Remote caching for developers

It is possible to use cache in read-only mode for developers to speed up local development. This can be achieved by enabling read-write remote caching in CI and read-only cache for individual developers. Since all Earthly cache is kept in Docker registries, managing access to the cache can be controlled by managing access to individual Docker images.

Note however that there is small performance penalty for regularly checking the remote registry on every run.
