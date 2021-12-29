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

An additional way in which you can improve the precision of the `COPY` comamnd is to use the [`.earthlyignore`](../earthfile/earthignore.md) file. Note, however, that this is best left as a last resort, as new files added to the project (that may be irrelevant to builds) would need to be manually added to `.earthlyignore`, which may be error-prone. It is much better to have to include every new file manually into the build (by adding it to a `COPY` command), than to exclude every new file manually (by adding it to the `.earthlyignore`), as whenever any such new file *must* be included, then the build would typically fail, making it harder to make a mistake compared to the reverse.

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

Build scripts serve slightly different purposes when they are run in CI compared to when they are executed for local development. Most of the logic is similar (hence Earthly attempts to unify the two concepts in a unified syntax), but there are some small differences.

For example, for development purposes, you may use commands such as `LOCALLY`, which cause Earthly to be less repeatable, and yet might satisfy very much needed use-cases that are typically out of scope of a CI build.

In addition, in CI it is much more likely that shared caching will be needed, while outputting artifacts and images locally would not be needed.

For these reasons, Earthly comes with the `--ci` flag, which simply expands to `--no-output --use-inline-cache --save-inline-cache --strict`. The `--ci` flag therefore, prevents the use of commands that are not repeatable, enables inline caching and disables outputting artifacts and images locally (as in CI the output is typically pushed or uploaded).

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

When using Earthly in a sandboxed CI, it may be useful to perform pushes on occasion, in order to populate shared caches. In addition, image pushes might also help developers to grab pre-built images for local use for various development workflows.

Pushing too often can result in slowing down builds due to the push time, while pushing too infrequently results in stale cache or development images.

A good balance is often to perform pushes on the `main` branch only, and to disable any pushing on PR builds. Although the `main` build will be slower, it will allow for maximum use of cache in PR builds, without the slowdown of further pushes.

Main branch build: `earthly --ci --push +target`.
PR build: `earthly --ci +target`.

The push option can also be configured via the env var `EARTHLY_PUSH`, which may be easier to manipulate in your CI of choice.

A more extreme case of this idea can be to use explicit maximum cache: `earthly --ci --push --remote-cache=.... --max-remote-cache +target`. The idea, again is to tradeoff performance on the `main` branch, for the benefit of faster PR builds. Whether this is actually beneficial needs to be measured on a project-by-project basis, however.

## Use `RUN --push` for deployment commands

If the result of a build needs to be pushed to an external service (or storage provider) and the destination is not an image registry, then you will need to use a custom push command (as opposed to a `SAVE IMAGE --push`).

To execute a custom push command, you can simply use a regular `RUN` command together with the `--push` flag. The `--push` will ensure that:

* The command is only executed when Earthly is run in push mode (`earthly --push`)
* No cache is reused for that specific command, causing it to execute every time
* The command is executed during the push phase of the build, ensuring that everything else (e.g. testing) has completed successfully

Here is an example of using the [github-release utility](https://github.com/github-release/github-release) to perform a push to GitHub Releases:

```Dockerfile
RUN --push --secret GITHUB_TOKEN github-release upload ...
```

## Use `--secret`, not `ARG`s to pass secrets to the build

If a build requires the usage of secrets, it is strongly recommended that you use the builtin secrets constructs, such as `earthly --secret`, [Earthly Cloud Secrets](../guides/cloud-secrets.md), and `RUN --secret`.

Using `ARG`s for passing secrets is strongly discouraged, as the secrets will be leaked in build logs, the build cache and the possibly in published images.

## Avoid copying secrets to the build environment

Even when using the proper builtin constructs for handling secrets, it is possible to then copy secrets in the build environment, which cause secrets to be leaked to a remote build cache, or to published images.

An simple example of how this may be possible:

```Dockerfile
# Bad
RUN --secret MY_SECRET echo "secret: $MY_SECRET" > /app/secret.txt
```

While this seems inoccuous and possibly uncommon, consider the following, which on the face of it might look like a good idea:

```Dockerfile
# Bad
RUN --secret AWS_ACCESS_KEY_ID --secret AWS_SECRET_ACCESS_KEY echo "[default]\naws_access_key_id=$AWS_ACCESS_KEY_ID\naws_secret_access_key=$AWS_SECRET_ACCESS_KEY" > /root/.aws/credentials
RUN aws ec2 describe-images
```

Another negative example is `COPY`ing the local credentials file:

```Dockerfile
# Bad
aws-creds:
    LOCALLY
    RUN cp "$HOME"/.aws/credentials ./.aws-creds
    SAVE ARTIFACT ./.aws-creds

do-something-with-aws:
    FROM ...
    COPY +aws-creds/.aws-creds /root/.aws/credentials
    RUN aws ec2 describe-images
```

The correct way to handle secrets that need to exist as files is to either mount them as secret files in the first place:

```Dockerfile
# Best
RUN --mount=type=secret,target=/root/.aws/credentials,id=AWS_CREDENTIALS \
    aws ec2 describe-images
```

This way, the credentials are never stored in the stored environment - they are only mounted during the execution of the `RUN` command.

Or, if you really have no choice, you may copy the secrets temporarily, but you **have** to remove them in the same layer:

```Dockerfile
# Ok, but error prone
RUN --secret AWS_ACCESS_KEY_ID --secret AWS_SECRET_ACCESS_KEY echo "[default]\naws_access_key_id=$AWS_ACCESS_KEY_ID\naws_secret_access_key=$AWS_SECRET_ACCESS_KEY" > /root/.aws/credentials ;\
    aws ec2 describe-images ;\
    rm /root/.aws/credentials
```

This should be avoided if possible, as it is error prone and might get secrets leaked if the `rm` is forgotten, or if the removal is performed under a separate `RUN` command.

```Dockerfile
# Bad: removal takes place in a separate layer, which means that the secrets will be leaked to the cache
RUN --secret AWS_ACCESS_KEY_ID --secret AWS_SECRET_ACCESS_KEY echo "[default]\naws_access_key_id=$AWS_ACCESS_KEY_ID\naws_secret_access_key=$AWS_SECRET_ACCESS_KEY" > /root/.aws/credentials
RUN aws ec2 describe-images
RUN rm /root/.aws/credentials
```

## Avoid exposing cache tags publicly if the cache contains private code or dependencies

TODO ...

## Do not pass Earthly dependencies from one target to another via the local file system or via the local Docker daemon

TODO ...

## Use `COPY +my-target/...` to pass files to and from `LOCALLY` targets

TODO ...

## Use `WITH DOCKER --load=+my-target` to pass images to `LOCALLY` targets

TODO ...

## Avoid non-deterministic behavior (such as randomness)

TODO ....

## Use cross-repo references, avoid `GIT CLONE` if possible

TODO ...

## Use `COPY --dir` to copy multiple directories

The classical Dockerfile `COPY` command differs from the unix `cp` in that it will copy directory *contents*, not the directories themselves. This requires that copying multiple directories to be split across multiple lines:

```Dockerfile
# Avoid: too verbose
COPY dir-1 dir-1
COPY dir-2 dir-2
COPY dir-3 dir-3
```

This is repetitive and uses more cache layers than should be necessary.

Earthly introduces a setting, `COPY --dir`, which makes `COPY` behave more like `cp` and less like the Dockerfile `COPY`. The `--dir` flag can be used therefore to copy multiple directories in a single command:

```Dockerfile
# Good
COPY --dir dir-1 dir-2 dir-3 ./
```

## Technique: Use `earthly -i` to debug failures

TODO ...

## Use separate images for build and production

TODO ...

(Build in one image, copy only the necessary files for the the final production image)

## Use `SAVE ARTIFACT ... AS LOCAL ...` for generated code, not `LOCALLY`

Many programming tools require the generation of code. The generated code is often used in completing a build, but also it might be required for IDEs to perform code completion. For this reason, it's often preferable that generated code is also output as local files during development.

It is recommended that generated code is saved via `SAVE ARTIFACT ... AS LOCAL ...` via regular Earthly targets, rather than via running the generation command in `LOCALLY`. There are multiple reasons for this:

* Executing commands via `LOCALLY` loses the repeatability benefits. This means that the same command could end up generating different code, depending on the system it is being run on. Differences in the environment, such as the version of code generator installed (e.g. `protoc`), or certain environment variables (e.g. `GOPATH`) could cause the generated code to be different.
* The logic to generate code via `LOCALLY` will not be usable in the CI, as the CI script would typically enable `--strict` mode.
* If the code generation workflow requires that the generated code is committed to the repository and then used in a subsequent earthly build, it is possible that due to human error, changes will be made to the input files, without the generated code to be updated correctly. If a problem or an incompatibility is introduced in this manner, it will show up for other people when they try to generate the code themselves. In worse cases, it may even go unnoticed and end up in production.

## Run everything in a single Earthly invocation, do not wrap Earthly

TODO ...
