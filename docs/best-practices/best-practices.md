# Best practices

Although Earthly has been designed to be unambiguous about what command to use for the job, writing Earthfiles can sometimes still be tricky, when it comes to nuances. As you try to accomplish certain tasks, you may find that sometimes the same result can be achieved using more than one technique. Or so it might seem.

Below we list some of the best practices that we have found to be useful in writing Earthfiles, with a focus on certain commands or techniques that seem similar, but aren't really, but also on some key points that we have seen newcomers stumble into.

## `COPY` only the minimal amount of files. Avoid copying `.git`

A typical mistake is to `COPY` entire large directories into the build environment and only using a subset of the files within them. Or worse, copying the entire repository (which might also include `.git`) for no good reason.

```Dockerfile
# Avoid
COPY . .
COPY * ./
```

The problem with this is that many of the files copied are not actually used during the build, however Earthly will react to changes to them, causing it to reuse cache inefficiently. It's not an issue of file size (though sometimes that too can hurt performance). It is much of an issue of re-executing build commands that wouldn't have to be re-executed.

```Dockerfile
# Avoid
COPY . .
RUN go mod download
RUN go build ...
```

In the above example, changing the project's `README.md` or running `git fetch` might cause slow commands like `go mod download` to be re-executed.

Earthly uses `COPY` commands (among other things) to mark certain files as inputs to the build. If any file included in a `COPY` changes, then the build will continue from that `COPY` command onwards. For this reason, you want to be as specific as possible when including files in a `COPY` command. In some cases, you might even have to list files one-by-one individually.

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

## `ENV` for image env vars, `ARG` for build configurability

`ENV` variables and `ARG` variables seem similar, however they are meant for different use-cases. Here is a breakdown of the differences, as well as how they differ from the Dockerfile-specific `ARG` command:

| | `ENV` | `ARG` | Dockerfile `ARG` |
| --- | --- | --- | --- |
| Available as an env-var in the same target | ✅ | ✅ | ❌ |
| Available for expanding within non-RUN commands | ❌ | ✅ | ✅ |
| Stored in the final image as an env-var | ✅ | ❌ | ❌ |
| Inherited via `FROM` | ✅ | ❌ | ❌ |
| Can be overriden when calling a build | ❌ | ✅ | ✅ |
| Can be propagated to other targets (via `BUILD +target --<key>=<value>` or similar) | ❌ | ✅ | N/A |

As you can see, the key situation where `ENV` is needed is when you want the value to be stored as part of the final image's configuration. This causes any `FROM` or `docker run` using that image to inherit the value.

However, if the use-case is build configurability, then `ARG` is the way to achieve that.

## `IF [...]` vs `RUN if [...]`

Earthly 0.6 introduces the conditional `IF` command, which allows for complex control flow within Earthly recipes. However, there is also the possiblity of using the shell `if` command to accomplish similar behavior. Which one should you use? Here is a quick comparison:

| | `IF` | `RUN if` |
| --- | --- | --- |
| Can execute any command as the expression | ✅ | ✅ |
| Can use mounts and secrets | ✅ | ✅ |
| Can use ARGs | ✅ | ✅ |
| Expression can be cached | ✅ | ✅ |
| Body runs in the same layer as the condition expression | ❌ | ✅ |
| Body can include any Earthly command  | ✅ | ❌ |

As you can see, `IF` is more powerful in that it can include other Earthly commands within it, allowing for rich conditional behavior. Examples might include optionally saving images, using different base images depending on a set of conditions, initializing `ARG`s with varying values.

`RUN if`, however is often simpler, and it only uses one layer.

As a best practice, it is recommended to use `RUN if` whenever possible (e.g. only `RUN` commands would be involved), to encourage simplicity, and otherwise to use `IF`.

## `FOR ... IN ...` vs `RUN for ... in ...`

As is the case with `IF` vs `RUN if`, there is a similar debate for the Earthly builtin command `FOR` vs `RUN for`. Here is a quick comparison of the two for flavors:

| | `FOR` | `RUN for` |
| --- | --- | --- |
| Can execute any command as the expression | ✅ | ✅ |
| Can use mounts and secrets | ✅ | ✅ |
| Can use ARGs | ✅ | ✅ |
| Expression can be cached | ✅ | ✅ |
| Can iterate over a constant list | ✅ | ✅ |
| Can iterate over a list resulting from an expression | ✅ | ✅ |
| Body runs in the same layer as the for expression | ❌ | ✅ |
| Body can include any Earthly command  | ✅ | ❌ |

Similar to the `IF` vs `RUN if` comparison, `FOR` is more powerful in that it can include other Earthly commands within it, allowing for rich iteration behavior. Examples might include iterating over a list of directories in a monorepo and calling Earthly targets within them, performing `SAVE IMAGE` over a list of container image tags.

`RUN for`, however is often simpler, and it only uses one layer.

As a best practice, it is recommended to use `RUN for` whenever possible (e.g. only `RUN` commands would be involved), to encourage simplicity, and otherwise to use `FOR`.

## Use `--ci` when running in CI

...

## Avoid `LOCALLY` and other non-strict commands

Certain Earthly functionality is only meant to be used for local development only. Most such commands do not fully abide by the Earthly spirit of repeatable builds, however, for certain specific development use-cases they are needed and therefore Earthly provides them. When Earthly is used in `--strict` mode (`earthly --strict +my-target` or `earthly --ci +my-target`), the usage of these commands is not allowed.

For this reason, it is recommended to avoid using these commands as much as possible, as doing so will:

1. Cause Earthly to behave in a non-repeatable way across other platforms, as it will realy on host-specific environment configuration.
2. Disable caching.
3. Cause the specific targets to not work at all when `--ci` is passed in.

An example of a command that is not allowed in strict mode is `LOCALLY`. The `LOCALLY` command skips the sandboxing of the build and executes all commands directly on the host machine.

Examples of valid cases where `LOCALLY` may be used are:

* installing dependencies on the host machine (e.g. to help IDEs provide better suggestions)
* executing tests on the host docker daemon, to help with inspection and debugging
* executing development commands which would otherwise require copying very large amounts of files to the sandboxed build environment

Note, however, that none of these cases are needed in a CI environment, and ultimately these commands are not regularly tested by a CI, which means they may break more frequently.

## Pattern: Optionally `LOCALLY`

In certain cases, it may be desirable to execute certain targets on the host machine, rather than in the sandboxed build environment, for debugging purposes. However, we need most of the target to execute in strict mode in CI. The solution to this is to use a target that can be optionally executed via `LOCALLY`. Here is an example:

Suppose we wanted the following target to be executed on against the host's Docker daemon:

```Dockerfile
FROM earthly/dind:alpine
WORKDIR /app
COPY docker-compose.yml ./
WITH DOCKER --compose docker-compose.yml \
        --service db \
        --load=+integration-test
    RUN docker-compose up integration
END
```

We could have an equivalent `LOCALLY` target:

```Dockerfile
LOCALLY
WITH DOCKER --compose docker-compose.yml \
        --service db \
        --load=+integration-test
    RUN docker-compose up integration
END
```

However, the code duplication is not ideal and will result in the two recipes to drift apart over time.

It is possible to use an `ARG` to decide on whether to execute the target on the host or not:

```Dockerfile
FROM alpine:3.13
ARG run_locally=false
IF [ "$run_locally" = "true" ]
    LOCALLY
ELSE
    FROM earthly/dind:alpine
    WORKDIR /app
    COPY docker-compose.yml ./
END
WITH DOCKER --compose docker-compose.yml \
        --service db \
        --load=+integration-test
    RUN docker-compose up integration
END
```

Now, to run locally, you can execute `earthly +my-target --run_locally=true`, otherwise `earthly +my-target` will execute in the sandboxed environment (the same way it executes in CI).

## Pattern: Deciding on a base image based on a condition

In some cases, it is useful to switch up which base image to use depending on the result of an `IF` expression. For example, let's assume that the company provided Go image only supports the `linux/amd64` platform, and therefore, you'd like to use the official golang image when ARM (`linux/arm64`) is detected. Here's how this can be achieved:

```Dockerfile
FROM alpine:3.13
ARG TARGETPLATFORM
IF [ "$TARGETPLATFORM" = "linux/arm64" ]
    FROM golang:1.16
ELSE
    FROM my-company/golang:1.16
END
```

This will cause the execution of consecutive `FROM`s within the same target. This is completely valid. On encountering another `FROM` expression, the current build environment is reset and another fresh root is initialized, containing the specified images data.

## Pattern: Push on the `main` branch only

TODO ...

## Use `RUN --push` for deployment commands

TODO ...

## Use `--secret`, not ARGs to pass secrets to the build

TODO ...

## Use `COPY +my-target/...` to pass files to `LOCALLY` targets

TODO ...

## Use `WITH DOCKER --load=+my-target` to pass files to `LOCALLY` targets

TODO ...

## Avoid non-deterministic behavior (such as randomness)

TODO ....

## Use cross-repo references, avoid `GIT CLONE` if possible

TODO ...
