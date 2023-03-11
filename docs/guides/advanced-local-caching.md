# Advanced caching techniques

In this article we're going to discuss various caching techniques that may be used to optimize your builds. For a primer on the basics of caching, see the [page on managing cache](./cache.md).

We will discuss a specific example where dependencies are cached using multiple techniques. Let's take the following build for instance

```Dockerfile
# Earthfile

VERSION 0.7
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    COPY build.gradle ./
    COPY src src
    RUN gradle build
    RUN gradle install
    SAVE ARTIFACT build/install/java-example/bin /bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib /lib AS LOCAL build/lib
```

If we ran this every time we changed the code, then the build would end up downloading all the dependencies specified in Gradle every single time. Here is how this can be improved.

### Option 1: Layer-based caching

One option is to use layer-based caching to first download all dependencies and only afterwards to copy the code and build it. Thus, when the code changes, the cache is bust after the step where we download dependencies.

```Dockerfile
# Earthfile

VERSION 0.7
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    # First, copy gradle build spec and download dependencies.
    COPY build.gradle ./
    RUN gradle build
    # Then copy code and build it.
    COPY src src
    RUN gradle build
    RUN gradle install
    SAVE ARTIFACT build/install/java-example/bin /bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib /lib AS LOCAL build/lib
```

This technique is the most simple and most robust technique, however it does not always cover all situations efficiently.

### Option 2: Mount-based caching (advanced)

Consider, for instance, a situation where the dependencies themselves change frequently. For example, in Java, if the build has `SNAPSHOT` dependencies (and/or are marked in Gradle as `changing = true`), then it needs to download the latest version available for those frequently.

For these cases, the build could use a cache mount. A cache mount is a volume that gets attached to a `RUN` command and is persisted and shared between runs.

```Dockerfile
# Earthfile

VERSION 0.7
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    COPY build.gradle ./
    COPY src src
    RUN --mount=type=cache,target=/root/.gradle/caches \
        gradle build && \
        gradle install
    SAVE ARTIFACT build/install/java-example/bin /bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib /lib AS LOCAL build/lib
```

This allows the Gradle command to reuse dependencies downloaded in a previous run and would only download new and updated dependencies, if any apply. It would not download everything from scratch again.

Although this approach works for many situations, it also has some possible downsides. For example, removing a dependency from the gradle spec, will not remove it from the cache. A mild problem can be that the cache becomes bloated over time. A worse possible problem can surface due to build inconsistencies created by the contents of the cache. Although some build tools are smart enough to ignore those cached dependencies if they have been removed from the spec, the behavior will vary from tool to tool and from language to language. In some cases, the build may continue to use a cached but removed dependency, yet if the build is executed on another system, it breaks or behaves differently.

To manage these situations you can elect to reset the cache of a build, by simply running the build again with the flag `--no-cache`: `earthly --no-cache +build`. Or, you can eliminate the entire cache altogether using [`earthly prune`](../earthly-command/earthly-command.md#earthly-prune).

Cache mounts are a leaky abstraction on top of the Earthly base principles by which nothing is shared. Although they are necessary to help optimize builds in some situations, care must be taken of possible edge cases whereby a build behaves differently because of the cache mount alone.

To help minimize such edge cases, Earthly only shares cache mounts between repeated builds of the *same target*, **and** only if *the build args are the same between invocations*. For more information see the [`RUN --mount` option reference](../earthfile/earthfile.md#run).
