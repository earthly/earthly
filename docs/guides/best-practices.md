# Best practices

Although Earthly has been designed to be unambiguous about what command to use for the job, writing Earthfiles can sometimes still be tricky, when it comes to nuances. As you try to accomplish certain tasks, you may find that sometimes the same result can be achieved using more than one technique. Or so it might seem.

Below we list some of the best practices that we have found to be useful in designing Earthly builds, with a focus on certain commands or techniques that seem similar, but aren't really, but also on some key points that we have seen newcomers stumble into.

## Table of contents

* [Earthfile-specific](#earthfile-specific)
    * [`COPY` only the minimal amount of files. Avoid copying `.git`](#copy-only-the-minimal-amount-of-files.-avoid-copying-.git)
    * [`ENV` for image env vars, `ARG` for build configurability](#env-for-image-env-vars-arg-for-build-configurability)
    * [Use cross-repo references, and avoid `GIT CLONE` if possible](#use-cross-repo-references-and-avoid-git-clone-if-possible)
    * [`GIT CLONE` vs `RUN git clone`](#git-clone-vs-run-git-clone)
    * [`IF [...]` vs `RUN if [...]`](#if-...-vs-run-if-...)
    * [`FOR ... IN ...` vs `RUN for ... in ...`](#for-...-in-...-vs-run-for-...-in-...)
    * [Pattern: Optionally `LOCALLY`](#pattern-optionally-locally)
    * [Pattern: Deciding on a base image based on a condition](#pattern-deciding-on-a-base-image-based-on-a-condition)
    * [Use `RUN --push` for deployment commands](#use-run-push-for-deployment-commands)
    * [Use `--secret`, not `ARG`s to pass secrets to the build](#use-secret-not-args-to-pass-secrets-to-the-build)
    * [Avoid copying secrets to the build environment](#avoid-copying-secrets-to-the-build-environment)
    * [Do not pass Earthly dependencies from one target to another via the local file system or via an external registry](#do-not-pass-earthly-dependencies-from-one-target-to-another-via-the-local-file-system-or-via-an-external-registry)
    * [Use `WITH DOCKER --pull`](#use-with-docker-pull)
    * [Style: Define the high-level targets at the top of the Earthfile](#style-define-the-high-level-targets-at-the-top-of-the-earthfile)
    * [Use `COPY +my-target/...` to pass files to and from `LOCALLY` targets](#use-copy-+my-target-...-to-pass-files-to-and-from-locally-targets)
    * [Use `WITH DOCKER --load=+my-target` to pass images to `LOCALLY` targets](#use-with-docker-load-+my-target-to-pass-images-to-locally-targets)
    * [Avoid non-deterministic behavior](#avoid-non-deterministic-behavior)
    * [Use `COPY --dir` to copy multiple directories](#use-copy-dir-to-copy-multiple-directories)
    * [Use separate images for build and production](#use-separate-images-for-build-and-production)
    * [Use `SAVE ARTIFACT ... AS LOCAL ...` for generated code, not `LOCALLY`](#use-save-artifact-...-as-local-...-for-generated-code-not-locally)
    * [Multi-line strings](#multi-line-strings)
    * [Multi-line commands](#multi-line-commands)
    * [Copying files from outside the build context](#copying-files-from-outside-the-build-context)
    * [Repository structure: Place build logic as close to the relevant code as possible](#repository-structure-place-build-logic-as-close-to-the-relevant-code-as-possible)
    * [Repository structure: Do not place all Earthfiles in a dedicated directory](#repository-structure-do-not-place-all-earthfiles-in-a-dedicated-directory)
    * [Pattern: Pass-through artifacts or images](#pattern-pass-through-artifacts-or-images)
    * [Use `earthly/dind`](#use-earthly-dind)
    * [Pattern: Saving artifacts resulting from a `WITH DOCKER`](#pattern-saving-artifacts-resulting-from-a-with-docker)
* [Usage-specific](#usage-specific)
    * [Use `--ci` when running in CI](#use-ci-when-running-in-ci)
    * [Avoid `LOCALLY` and other non-strict commands](#avoid-locally-and-other-non-strict-commands)
    * [Pattern: Push on the `main` branch only](#pattern-push-on-the-main-branch-only)
    * [Do not expose cache image tags publicly if the cache contains private code or dependencies](#do-not-expose-cache-image-tags-publicly-if-the-cache-contains-private-code-or-dependencies)
    * [Technique: Use `earthly -i` to debug failures](#technique-use-earthly-i-to-debug-failures)
    * [Run everything in a single Earthly invocation, do not wrap Earthly](#run-everything-in-a-single-earthly-invocation-do-not-wrap-earthly)
    * [Use `RUN --ssh` for passing host SSH keys to builds](#use-run-ssh-for-passing-host-ssh-keys-to-builds)
    * [Use SAVE IMAGE to always cache](#use-save-image-to-always-cache)
    * [Future: Saving an artifact even if the build fails](#future-saving-an-artifact-even-if-the-build-fails)

## Earthfile-specific

### `COPY` only the minimal amount of files. Avoid copying `.git`

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

An additional way in which you can improve the precision of the `COPY` command is to use the [`.earthlyignore`](../earthfile/earthlyignore.md) file. Note, however, that this is best left as a last resort, as new files added to the project (that may be irrelevant to builds) would need to be manually added to `.earthlyignore`, which may be error-prone. It is much better to have to include every new file manually into the build (by adding it to a `COPY` command), than to exclude every new file manually (by adding it to the `.earthlyignore`), as whenever any such new file *must* be included, then the build would typically fail, making it harder to make a mistake compared to the opposite.

### `ENV` for image env vars, `ARG` for build configurability

`ENV` variables and `ARG` variables seem similar, however they are meant for different use-cases. Here is a breakdown of the differences, as well as how they differ from the Dockerfile-specific `ARG` command:

| | `ENV` | `ARG` | Dockerfile `ARG` |
| --- | --- | --- | --- |
| Available as an env-var in the same target | ✅ | ✅ | ❌ |
| Available for expanding within non-RUN commands | ❌ | ✅ | ✅ |
| Stored in the final image as an env-var | ✅ | ❌ | ❌ |
| Inherited via `FROM` | ✅ | ❌ | ❌ |
| Can be overridden when calling a build | ❌ | ✅ | ✅ |
| Can be propagated to other targets (via `BUILD +target --<key>=<value>` or similar) | ❌ | ✅ | N/A |

As you can see, the key situation where `ENV` is needed is when you want the value to be stored as part of the final image's configuration. This causes any `FROM` or `docker run` using that image to inherit the value.

However, if the use-case is build configurability, then `ARG` is the way to achieve that.

### Use cross-repo references, and avoid `GIT CLONE` if possible

Earthly provides rich set of features to allow working with and across Git repositories. It is recommended to use Earthly [cross-repository references](../guides/importing.md) rather than `GIT CLONE` or `RUN git clone`, whenever possible.

Here is an example.

Repo 1:

```
repo 1
├── README.md
└── my-file.txt
```

Repo 2:

```Dockerfile
# Bad
VERSION 0.8
FROM alpine:3.18
WORKDIR /work
print-file:
    GIT CLONE git@github.com:my-co/repo-1.git
    RUN echo my-file.txt
```

This might be addressed in the following way:

Repo 1:

```
repo 1
├── README.md
├── Earthfile
└── my-file.txt
```

```Dockerfile
# Repo 1 Earthfile
VERSION 0.8
FROM alpine:3.18
WORKDIR /work
file:
    COPY ./my-file.txt ./
    SAVE ARTIFACT ./my-file.txt
```

Repo 2:

```Dockerfile
# Repo 2 Earthfile
VERSION 0.8
IMPORT github.com/my-co/repo-1
FROM alpine:3.18
WORKDIR /work
print-file:
    COPY repo-1+file/my-file.txt ./
    RUN echo my-file.txt
```


There are multiple benefits to using cross-repository references in this manner:

* The build of repo 1 can evolve to more than just passing a file to another repository. It may be possible to also export generated code, artifacts, base images or full microservice images in the future, if they are needed.
* It is clearer about which files are actually needed externally, as they are declared via `SAVE ARTIFACT`. This makes the code more readable and maintainable. The fact that an artifact is saved during a build constitutes an explicit API of the repository.

Of course, the down-side is that repo 1 requires an Earthfile to be added, and that might not always be feasible. It's possible that repo 1 is controlled by another team, or that it is entirely external to the company. In such cases, `GIT CLONE` might help to provide a faster, yet imperfect solution.

Another use-case where `GIT CLONE` is better suited is when the operation needs to take place on the whole source repository. For example, performing Git operations, such as tagging, creating branches, or merging.

Finally, here is a comparison between cross-repo references and `GIT CLONE`:

| | Cross-repo reference | `GIT CLONE` |
| --- | --- | --- |
| Example | `FROM github.com/my-co/my-proj:my-branch+my-target` | `GIT CLONE --branch=my-branch git@github.com:my-co/my-proj` |
| Earthly can pass-through SSH agent access from the host | ✅ | ✅ |
| Access to HTTPS repositories can be configured in Earthly | ✅ | ✅ |
| Can specify branch or tag | ✅ - via `:<branch>` | ✅ - via `--branch` |
| Source configurable via `ARG`s | ✅ | ✅ |
| Protocol-agnostic referencing | ✅ | ❌ - can be `ssh://`, `https://`, `git@github.com` etc |
| Clear declaration of the dependency | ✅ - source repo needs to expose it in the Earthfile | ❌ |
| Can be used without modifications to the source repository | ❌ - requires Earthfile | ✅ |
| Can operate on the repository itself | ❌ - possible, but not designed for this | ✅ |

### `GIT CLONE` vs `RUN git clone`

Earthly has a built-in `GIT CLONE` instruction that can be used to clone a Git repository. It is recommended that `GIT CLONE` is used rather than `RUN git clone`, for a few reasons:

* Earthly treats `GIT CLONE` as a first-class input (BuildKit source). As such, Earthly caches the repository internally and downloading only incremental differences on changes.
* Earthly is commit hash-aware, so it'll be able to detect when the build needs to take place versus when there are no changes to be made and the cache can be reused. If a change takes place in the source repository, `RUN git clone` would not be able to detect that, as it is not recognized as an input. So it would naively reuse the cache when it shouldn't.
* `GIT CLONE` will pass-through Earthly settings for [authentication](../guides/auth.md), such as SSH agent access and/or HTTPS credentials.

`GIT CLONE` does have some limitations, however. It only performs a shallow clone, it does not have the branch information, it does not have origin information, and it does not have the tags downloaded. Even in such cases, it might be better to attempt to reintroduce the information after a `GIT CLONE`, whenever possible, in order to gain the caching benefits.

When this proves to be too difficult, or impossible, and you really need to perform a custom `RUN git clone`, consider using both in conjunction, to gain the hash awareness benefits.

```Dockerfile
# Bad
RUN git clone git@github.com/my-co/my-proj
WORKDIR my-proj
RUN ls
```

```Dockerfile
# Good
GIT CLONE git@github.com/my-co/my-proj my-proj
WORKDIR my-proj
RUN ls
```

```Dockerfile
# Ok, if you have no choice
ARG git_url="git@github.com/my-co/my-proj"
GIT CLONE "$git_url" my-proj
ARG git_hash=$(cd my-proj; git rev-parse HEAD)
RUN rm -rf my-proj &&\
    git clone "$git_url" my-proj &&\
    cd my-proj &&\
    git checkout "$git_hash"
WORKDIR my-proj
RUN ls
```

Alternatively, an optimal deep clone can be achieved by calling the `FUNCTION` [DEEP_CLONE](https://github.com/earthly/lib/blob/3.0.1/utils/git/Earthfile#L4):
```Dockerfile
ARG git_url="git@github.com/my-co/my-proj"
DO github.com/earthly/lib/utils/git:3.0.1+DEEP_CLONE --GIT_URL=$git_url
RUN ls
```

{% hint style='info' %}
##### Note
The final "if you have no choice" example contains an `ARG` shell-out, which will be set to the latest git sha; this is done so that if the git repository doesn't
contain any new commits between multiple runs, the `ARG git_hash` value will stay constant, and will allow earthly to skip over the `RUN` command; whereas a `RUN --no-cache`
command will be run each time even if the `ARG git_hash` value has never changed.
{% endhint %}

Finally, here is a comparison between `GIT CLONE` and `RUN git clone`:

| | `GIT CLONE` | `RUN git clone` |
| --- | --- | --- |
| Earthfiles can be protocol-agnostic | ❌ - can be `ssh://`, `https://`, `git@github.com` etc | ❌ - can be `ssh://`, `https://`, `git@github.com` etc |
| Can configure access in Earthly, to keep Earthfiles agnostic | ✅ | ❌ |
| Earthly can pass-through SSH agent access from the host | ✅ | ✅ - but it requires `RUN --ssh` |
| Access to HTTPS repositories can be configured in Earthly | ✅ | ❌ - but possible to pass credentials via secrets |
| Cache-aware - incremental pulls | ✅ | ❌ |
| Commit hash-aware - rebuild when there are changes in remote repository | ✅ | ❌ |

### `IF [...]` vs `RUN if [...]`

Earthly 0.6 introduces the conditional `IF` command, which allows for complex control flow within Earthly recipes. However, there is also the possibility of using the shell `if` command to accomplish similar behavior. Which one should you use? Here is a quick comparison:

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

### `FOR ... IN ...` vs `RUN for ... in ...`

As is the case with `IF` vs `RUN if`, there is a similar debate for the Earthly builtin command `FOR` vs `RUN for`. Here is a quick comparison of the two 'for' flavors:

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

### Pattern: Optionally `LOCALLY`

In certain cases, it may be desirable to execute certain targets on the host machine, rather than in the sandboxed build environment, for debugging purposes. However, we need most of the targets to execute in strict mode in CI. The solution to this is to use a target that can be optionally executed via `LOCALLY`. Here is an example:

Suppose we wanted the following target to be executed on against the host's Docker daemon:

```Dockerfile
FROM earthly/dind:alpine-3.19-docker-25.0.3-r2
WORKDIR /app
COPY docker-compose.yml ./
WITH DOCKER --compose docker-compose.yml \
        --service db \
        --load=+integration-test
    RUN docker-compose up integration
END
```

We could therefore have an equivalent `LOCALLY` target:

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
FROM alpine:3.18
ARG run_locally=false
IF [ "$run_locally" = "true" ]
    LOCALLY
ELSE
    FROM earthly/dind:alpine-3.19-docker-25.0.3-r2
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

### Pattern: Deciding on a base image based on a condition

In some cases, it is useful to switch up which base image to use depending on the result of an `IF` expression. For example, let's assume that the company provided Go image only supports the `linux/amd64` platform, and therefore, you'd like to use the official golang image when ARM (`linux/arm64`) is detected. Here's how this can be achieved:

```Dockerfile
FROM alpine:3.18
ARG TARGETPLATFORM
IF [ "$TARGETPLATFORM" = "linux/arm64" ]
    FROM golang:1.16
ELSE
    FROM my-company/golang:1.16
END
```

This will cause the execution of consecutive `FROM`s within the same target. This is completely valid. On encountering another `FROM` expression, the current build environment is reset and another fresh root is initialized, containing the specified images data.

### Use `RUN --push` for deployment commands

If the result of a build needs to be pushed to an external service (or storage provider) and the destination is not an image registry, then you will need to use a custom push command (as opposed to a `SAVE IMAGE --push`).

To execute a custom push command, you can simply use a regular `RUN` command together with the `--push` flag. The `--push` will ensure that:

* The command is only executed when Earthly is in push mode (`earthly --push`)
* No cache is reused for that specific command, causing it to execute every time
* The command is executed during the push phase of the build, ensuring that everything else (e.g. testing) has completed successfully first

Let's look at an example of using the [github-release utility](https://github.com/github-release/github-release) to perform a push to GitHub Releases.

```Dockerfile
# Bad, and dangerous
RUN --no-cache --secret GITHUB_TOKEN github-release upload ...
```

`RUN --no-cache` should be avoided for this use-case, as it has some potentially dangerous downsides:

* The upload command may be executed in parallel with any testing (meaning that tests might not pass yet the upload may still complete)
* The upload will execute even when earthly is not invoked in `--push` mode.

To address this issue, it is advisable to use `RUN --push` instead.

```Dockerfile
# Good
RUN --push --secret GITHUB_TOKEN github-release upload ...
```

### Use `--secret`, not `ARG`s to pass secrets to the build

If a build requires the usage of secrets, it is strongly recommended that you use the builtin secrets constructs, such as `earthly --secret`, [Earthly Cloud Secrets](../cloud/cloud-secrets.md), and `RUN --secret`.

Using `ARG`s for passing secrets is strongly discouraged, as the secrets will be leaked in build logs, the build cache and the possibly in published images.

### Avoid copying secrets to the build environment

Even when using the proper builtin constructs for handling secrets, it is possible to then copy secrets in the build environment, which cause secrets to be leaked to a remote build cache, or to published images.

A simple example of how this may be possible:

```Dockerfile
# Bad
RUN --secret MY_SECRET echo "secret: $MY_SECRET" > /app/secret.txt
```

While this seems innocuous and possibly uncommon, consider the following, which on the face of it might look like a good idea:

```Dockerfile
# Bad
RUN --secret AWS_ACCESS_KEY_ID --secret AWS_SECRET_ACCESS_KEY echo "[default]\naws_access_key_id=$AWS_ACCESS_KEY_ID\naws_secret_access_key=$AWS_SECRET_ACCESS_KEY" > /root/.aws/credentials
RUN aws ec2 describe-images
```

Another negative example is copying the local credentials file:

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

### Do not pass Earthly dependencies from one target to another via the local file system or via an external registry

If you are new to Earthly, you may be tempted to save an artifact locally in one target and then to retrieve it another one.

```Dockerfile
# Bad
all:
    BUILD +dep
    BUILD +build
dep:
    ...
    SAVE ARTIFACT my-artifact.jar AS LOCAL ./build/my-artifact.jar
build:
    ...
    COPY ./build/my-artifact.jar ./
    ...
```

This will not actually work, as in Earthly all output takes place only at the end of a successful build. Meaning that when `+build` starts, the artifact would not have been output yet. In fact, `+dep` and `+build` will run completely parallel anyway - as Earthly does not know of a dependency between them.

The proper way to achieve this is to use [artifact references](../guides/importing.md).

```Dockerfile
# Good
all:
    BUILD +build
dep:
    ...
    SAVE ARTIFACT my-artifact.jar
build:
    ...
    COPY +dep/my-artifact.jar ./
    ...
```

Notice that `+dep` no longer needs to save the file locally. Also, the `COPY` command no longer references the file from the local file system. It has been replaced with an artifact reference from the target `+dep`. This reference will tell Earthly that these two targets depend on each other and will therefore schedule the relevant parts to run sequentially.

Notice also that in our `+all` target, we no longer have to call both `+dep` and `+build`. The system will automatically infer that when building `+build`, `+dep` is also required.

Another example of what you should **not** do is to pass Earthly images via between targets via an external registry.

```Dockerfile
# Bad
all:
    BUILD +dep-img
    BUILD +test
dep-img:
    ...
    SAVE IMAGE --push my-co/my-image:latest
test:
    WITH DOCKER
        RUN docker run my-co/my-image:latest
    END
```

```Dockerfile
# Also bad
all:
    BUILD +test
dep-img:
    ...
    SAVE IMAGE --push my-co/my-image:latest
test:
    BUILD +dep-img # This still does not work
    WITH DOCKER
        RUN docker run my-co/my-image:latest
    END
```

Similarly, in this case, pushing of the image takes place at the end of the build, which means that when `+test` runs, it will not have the image available, unless it has been pushed in a previous execution (which means that the image may be stale).

To fix this, we need to use `WITH DOCKER --load` and a [target reference](../guides/importing.md):

```Dockerfile
# Good
all:
    BUILD +test
dep-img:
    ...
    SAVE IMAGE my-co/my-image:latest
test:
    WITH DOCKER --load=+dep-img
        RUN docker run my-co/my-image:latest
    END
```

The `--load` instruction will inform Earthly that the two targets depend on each other and will therefore build the image and load it into the Docker daemon provided by `WITH DOCKER`.

### Use `WITH DOCKER --pull`

When referencing an external image in the body of a `WITH DOCKER` block, it is important to declare it via `WITH DOCKER --pull`, for a few reasons:

* The image will be cached as part of buildkit, allowing for faster builds. This is especially important as `WITH DOCKER` wipes the state of the Docker daemon (including its cache) after every run.
* The Daemon within `WITH DOCKER` is not logged into registries. Your local Docker login config is not propagated to the daemon. This means that you may run into issues when trying to pull images from private registries, but also, DockerHub rate limiting may prevent you from pulling images consistently from public repositories.

```Dockerfile
# Bad: Image hello-world needs to be pulled every time and is not part of the Earthly-managed cache.
WITH DOCKER
   RUN docker run hello-world
END
```

```Dockerfile
# Good
WITH DOCKER --pull hello-world
   RUN docker run hello-world
END
```

If you use `WITH DOCKER --compose`, Earthly will automatically pull images declared in the compose file for you, as long as they are not already being loaded from another target via `WITH DOCKER --load`. So in this case, you do not need to declare those image with `WITH DOCKER --pull`.

### Style: Define the high-level targets at the top of the Earthfile

High-level targets are those targets that are meant to be executed directly by the user on the command-line or via the CI.

As software engineers, we read code more often than we write it. As a matter of style, it is recommended to declare the higher-level targets at the top of the Earthfile, to help with the usability of the Earthfile. This will help fellow engineers who have not worked on the Earthfile to quickly find the relevant targets to use in their day-to-day development.

It also helps a reader to consume the Earthfile starting from the top, forming a high-level picture first, then gradually going deeper and deeper to lower-level logic.

### Use `COPY +my-target/...` to pass files to and from `LOCALLY` targets

When using `LOCALLY`, it is tempting to skip on using Earthly constructs for passing files between targets. However, this can be problematic.

```Dockerfile
# Bad
all:
    BUILD +dep
    BUILD +build
dep:
    LOCALLY
    RUN echo "Hello World" > ./my-artifact.txt
build:
    COPY ./my-artifact.txt ./
    ...
```

This setup may actually work, but it has a key issue: the order of `+dep` and `+build` is not guaranteed. So in some runs, the file `./my-artifact.txt` will be created before the `+build` target is executed, and in some runs it will be created after.

To fix this race condition, you need to use an [artifact reference](../guides/importing.md), to ensure that Earthly is aware of the dependency between the two targets:

```Dockerfile
# Good
all:
    BUILD +build
dep:
    LOCALLY
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt
build:
    COPY +dep/my-artifact.txt ./
    ...
```

Here is another example of the reverse (copying a file to a `LOCALLY` target):

```Dockerfile
# Bad
all:
    BUILD +dep
    BUILD +run-locally
dep:
    FROM alpine:3.18
    WORKDIR /work
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt AS LOCAL ./build/my-artifact.txt
run-locally:
    LOCALLY
    RUN echo ./build/my-artifact.txt
```

The mistake here is relying on `SAVE ARTIFACT ... AS LOCAL ...` for the transfer of the artifact to the `LOCALLY` target. As Earthly outputs are written at the end of the build, the target `+run-locally` will not have the file in time (or it might have it from a previous run only, meaning that it might be stale).

Here is how to fix this:

```Dockerfile
# Good
all:
    BUILD +run-locally
dep:
    FROM alpine:3.18
    WORKDIR /work
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt
run-locally:
    LOCALLY
    COPY +dep/my-artifact.txt ./build/my-artifact.txt
    RUN echo ./build/my-artifact.txt
```

The `COPY` command using an artifact reference will inform Earthly of the dependency between the two targets, and will therefore cause the transfer of artifact between the two properly.

And finally, here is another common mistake, when passing files between two `LOCALLY` targets:

```Dockerfile
# Bad
all:
    BUILD +dep
    BUILD +run-locally
dep:
    LOCALLY
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt AS LOCAL ./build/my-artifact.txt
run-locally:
    LOCALLY
    RUN echo ./build/my-artifact.txt
```

```Dockerfile
# Also bad
all:
    BUILD +run-locally
dep:
    LOCALLY
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt AS LOCAL ./build/my-artifact.txt
run-locally:
    BUILD +dep # Order still not guaranteed
    LOCALLY
    RUN echo ./build/my-artifact.txt
```

Here, the mistake is that the order of operations is not guaranteed. Earthly does not know that the two targets depend on each other, and therefore might decide to run them out of order. It might work sometimes, but it is not guaranteed that it will work every time.

To address this, again, the relationship between the two targets should be declared via `COPY` and an artifact reference.

```Dockerfile
# Good
all:
    BUILD +run-locally
dep:
    LOCALLY
    RUN echo "Hello World" > ./my-artifact.txt
    SAVE ARTIFACT ./my-artifact.txt
run-locally:
    LOCALLY
    COPY +dep/my-artifact.txt ./build/my-artifact.txt
    RUN echo ./build/my-artifact.txt
```

### Use `WITH DOCKER --load=+my-target` to pass images to `LOCALLY` targets

Earthly is able to output Docker images to the local Docker daemon at the end of each build. However, when requiring an image for a `LOCALLY` target, the image needs to be output in the *middle* of the build.

```Dockerfile
# Bad
all:
    BUILD +build-img
    BUILD +run-img
build-img:
    ...
    SAVE IMAGE my-co/my-img:latest
run-img:
    LOCALLY
    RUN docker run my-co/my-img:latest
```

The above will not work as the output will take place at the end of the build only. In addition, Earthly is unaware that there is a dependency between the two targets. To address this, we need to use `WITH DOCKER --load` and a [target reference](../guides/importing.md):

```Dockerfile
# Good
all:
    BUILD +run-img
build-img:
    ...
    SAVE IMAGE my-co/my-img:latest
run-img:
    LOCALLY
    WITH DOCKER --load=+build-img
        RUN docker run my-co/my-img:latest
    END
```

The `--load` instruction will inform Earthly of the dependency and will therefore cause the image to be output right before the `WITH DOCKER` `RUN` command executes.

### Avoid non-deterministic behavior

It is generally recommended to avoid any non-deterministic behavior when designing Earthly builds. This may include:

* Introducing time-stamps in builds or in tags
* Generating unique IDs
* Initializing `ARG` with values that include randomness

The main reason to avoid non-deterministic behavior is to ensure that builds are repeatable, and to maximize the use of cache. If an intermediate step leads to the same result as a previous run, Earthly may be able to reuse further computation performed previously.

Many compilers, code generators and other tools might not be deterministic and there may be no way around it. Earthly still functions correctly in these cases, however there may be occasions where the cache is not fully utilized to its potential.

### Use `COPY --dir` to copy multiple directories

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

### Use separate images for build and production

To keep production images small, it is advisable to start from a new base image and to install only production-required dependencies and then to copy in only the final built binaries or packages. This technique may vary from language to language, depending on the ecosystem-specific tooling.

An example, for Go, you might have a development image, that contains the entire Go development tools, including the `go` binary. After the application binary has been built via `go build`, there is no longer a need for the `go` binary. So the production image should not contain it. Here is an example:

```Dockerfile
# Avoid: production image is bloated
FROM go:1.16
RUN apk add ... # development + production dependencies
build:
    COPY ...
    RUN go mod download
    COPY ...
    RUN go build ... -o /usr/bin/app
    ENTRYPOINT ["/usr/bin/app"]
    SAVE IMAGE my-production-image:latest
```

Here is a way to address this:

```Dockerfile
# Good
FROM go:1.16
RUN apk add ... # development dependencies
build:
    COPY ...
    RUN go mod download
    COPY ...
    RUN go build ... -o ./build/app
    SAVE ARTIFACT ./build/app
image:
    FROM alpine:3.18 # start afresh
    RUN apk add ... # production dependencies only
    COPY +build/app /usr/bin/app
    ENTRYPOINT ["/usr/bin/app"]
    SAVE IMAGE my-production-image:latest
```

### Use `SAVE ARTIFACT ... AS LOCAL ...` for generated code, not `LOCALLY`

Many programming tools require the generation of code. The generated code is often used in completing a build, but also it might be required for IDEs to perform code completion. For this reason, it's often preferable that generated code is also output as local files during development.

It is recommended that any generated code is saved via `SAVE ARTIFACT ... AS LOCAL ...` via regular Earthly targets, rather than via running the generation command in `LOCALLY`. There are multiple reasons for this:

* Executing commands via `LOCALLY` loses the repeatability benefits. This means that the same command could end up generating different code, depending on the system it is being run on. Differences in the environment, such as the version of code generator installed (e.g. `protoc`), or certain environment variables (e.g. `GOPATH`) could cause the generated code to be different.
* The logic to generate code via `LOCALLY` will not be usable in the CI, as the CI script would typically enable `--strict` mode.
* If the code generation workflow requires that the generated code is committed to the repository and then used in a subsequent earthly build, it is possible that due to human error, changes will be made to the input files, without the generated code to be updated correctly. If a problem or an incompatibility is introduced in this manner, it will only show up later, for other people when they try to generate the code themselves. In worse cases, it may even go unnoticed and end up in production.

### Multi-line strings

To specify a multi-line string in Earthly, you can simply start quotes on one line and end them on another.

```Dockerfile
# Bad
RUN echo "this is a" > /tmp/file
RUN echo "multi-line string" >> /tmp/file
RUN echo "that goes" >> /tmp/file
RUN echo "on" >> /tmp/file
RUN echo "and on" >> /tmp/file
ARG MULTILINE_STRING=$(cat /tmp/file)
```

```Dockerfile
# Good
ARG MULTILINE_STRING="this is a
multi-line string
that goes
on
and on"
```

### Multi-line commands

To execute commands that may span multiple lines, you can use the line continuation character (`\`). Remember to chain multiple shell commands via `&&` in order to correctly exit if one of the commands fails.

```Dockerfile
RUN go build ... && \
    if [ "$FOO" = "bar" ]; then \
        echo "spaghetti" > ./default-food.txt ;\
    fi
```

### Copying files from outside the build context

It is generally advisable to avoid copying files outside of the build context. If a file is required from a sibling directory, or from a parent directory, it is recommended that those files are exported via `SAVE ARTIFACT` and then copied over using an artifact reference.

```
├── dir1
|   └── some-file.txt
└── dir2
    ├── other-files...
    └── Earthfile
```

```Dockerfile
# ./dir2/Earthfile
# Bad: does not work
COPY ../dir1/some-file.txt ./
```

In the above example, the file `some-file.txt` is copied from the sibling directory `dir1`. This will not work in Earthly as the file is not in the build context of `./dir2/Earthfile` (the build context in this case is `./dir2`). To address this issue, we can create an Earthfile in `dir1` that exports the file `some-file.txt` as an artifact.

```
├── dir1
|   ├── some-file.txt
|   └── Earthfile
└── dir2
    ├── other-files...
    └── Earthfile
```

```Dockerfile
# ./dir1/Earthfile
VERSION 0.8
FROM alpine:3.18
WORKDIR /work
file:
    COPY some-file.txt ./
    SAVE ARTIFACT ./some-file.txt
```

```Dockerfile
# ./dir2/Earthfile
# Good
COPY ../dir1+file/some-file.txt ./
```

The passing of the file as an artifact will also help create a build API of `dir1`, where all the files required outside of it are explicitly exported.

If a file is needed and there is no good way of adding an Earthfile to the directory containing it (e.g. a common file from the user's home directory), then an option is to use `LOCALLY`.

```Dockerfile
file:
    LOCALLY
    SAVE ARTIFACT $HOME/some-file.txt
do-something:
    COPY +file/some-file.txt ./
```

Note, however, that `LOCALLY` is not allowed in `--strict` mode (or in `--ci` mode), as it introduces a dependency from the host machine, which may interfere with the repeatability property of the build.

Although performing a `COPY ../` is not possible in Earthly today, there are some rare, but valid use-cases for this functionality. This is being discussed in GitHub issue [#1221](https://github.com/earthly/earthly/issues/1221).

### Repository structure: Place build logic as close to the relevant code as possible

When designing builds, it is advisable to place lower-level build logic closer to the code that it is building. This can be achieved by splitting Earthly builds across multiple Earthfiles, and placing some of the Earthfiles deeper inside the directory structure. The lower-level Earthfiles can then export artifacts and/or images via `SAVE ARTIFACT` or `SAVE IMAGE` commands, respectively. Those artifacts can then be referenced in higher-level Earthfiles via artifact and target references (`COPY ./deep/dir+some-target/an/artifact ...`, `FROM ./some/path+my-target`).

This allows for low coupling between modules within your code and creates a "build API" for your directories, whereby all externally accessible artifacts are exposed explicitly.

As one example, you might find the [monorepo example](https://github.com/earthly/earthly/tree/main/examples/monorepo) to be a useful case-study. However, even when a repository contains a single project, you might still find it useful to split logic across multiple Earthfiles. An example might be including Protocol Buffers generation logic inside the subdirectory containing the `.proto` files, in its own Earthfile.

For a real-world example, you can also take a look at Earthly's own build, where several Earthfiles are scattered across the repository to help organize build logic across modules, very much like regular code. Here are some examples:

<!-- GitBook currently has a bug where any references to an "Earthfile" gets confused with "docs/Earthfile" and somehow appends a /README.md
https://github.com/earthly/earthly/blob/main/Earthfile was changed to https://tinyurl.com/yt3d3cx6 -->

* [`ast/parser`](https://github.com/earthly/earthly/tree/main/ast/parser) - Earthfile contains the logic for generating Go source code based on an ANTLR grammar.
* [`ast/parser/tests`](https://github.com/earthly/earthly/tree/main/ast/tests) - Earthfile contains logic for running AST-specific tests.
* [`buildkitd`](https://github.com/earthly/earthly/tree/main/buildkitd) - Earthfile contains the logic for building the Earthly buildkit image.
* [`tests`](https://github.com/earthly/earthly/tree/main/tests) - Earthfile contains logic for executing e2e tests.
* [`release/**/`](https://github.com/earthly/earthly/tree/main/release) - Multiple Earthfiles contain logic used for the release of Earthly.
* [The main Earthfile](https://tinyurl.com/yt3d3cx6) - ties everything together, referencing the various targets across the sub-directories.

### Repository structure: Do not place all Earthfiles in a dedicated directory

A common practice when using Dockerfiles is to place all Dockerfiles in a special directory of the repository.

```
# Pattern common for Dockerfiles, but should be avoided for Earthfiles
├── app1-src-dir
|   └── ...
├── app2-src-dir
|   └── ...
├── app3-src-dir
|   └── ...
└── services
    ├── app1.Dockerfile
    ├── app2.Dockerfile
    └── app3.Dockerfile
```

And then running `docker build -f ./services/app1.Dockerfile ./app1-src-dir ...` and so on.

In Earthly, however, this is an anti-pattern, for a couple reasons:

- Every repository using Earthly should have a common structure, to help the user navigate the build. The convention is that Earthfiles are as close to the code as possible, with some high-level targets exposed in the root of the repository, or the root of the directory containing the code for a specific app. Having this convention helps the users who have not written the Earthfiles to quickly be able to browse around and understand the build, at least at a high level.
- Cross-directory and cross-repository references will point to directories where the user expects an Earthfile to be present, and then to a specific target or function within that Earthfile. It is important for this discoverability to be available to anyone browsing the build code and understanding the connections between Earthfiles.

For these reasons, Earthly does not support placing all Earthfiles in a single directory, nor the equivalent of a `docker build -f` option.

### Pattern: Pass-through artifacts or images

If a target acts as a wrapper for another target and that other target produces artifacts, you may find it useful for the wrapper to also emit the same artifacts. Consider the following example of the target `+build-for-windows`:

```Dockerfile
# No pass-through artifacts
VERSION 0.8
FROM alpine:3.18
build:
    ARG some_arg=...
    ARG another_arg=...
    ARG os=linux
    RUN ...
    SAVE ARTIFACT ./output
build-for-windows:
    BUILD +build --some_arg=... --another_arg=... --os=windows
```

```Dockerfile
# With pass-through artifacts
VERSION 0.8
FROM alpine:3.18
build:
    ARG some_arg=...
    ARG another_arg=...
    ARG os=linux
    RUN ...
    SAVE ARTIFACT ./output
build-for-windows:
    COPY (+build --some_arg=... --another_arg=... --os=windows) ./
    SAVE ARTIFACT ./*
```

The fact that `+build-for-windows` itself exports the artifacts means that it can be referenced directly in other targets as `COPY +build-for-windows/output ./`.

Similarly, if a target emits an image, then that image can be also emitted by a wrapping target like so:

```Dockerfile
# No pass-through image
VERSION 0.8
FROM alpine:3.18
build:
    ARG some_arg=...
    ARG another_arg=...
    RUN ...
    SAVE IMAGE some-intermediate-image:latest
build-wrapper:
    BUILD +build --some_arg=... --another_arg=...
```

```Dockerfile
# With pass-through image
VERSION 0.8
FROM alpine:3.18
build:
    ARG some_arg=...
    ARG another_arg=...
    RUN ...
    SAVE IMAGE some-intermediate-image:latest
build-wrapper:
    FROM +build --some_arg=... --another_arg=...
    SAVE IMAGE i-can-give-this-another-name:latest
```

This allows for `+build-wrapper` to reuse the logic in `+build`, but ultimately to create an image that is saved under a different name. This can then be used in a `WITH DOCKER --load` statement directly (whereas if there was no image pass-through, then `+build-wrapper` couldn't have been used).

### Use `earthly/dind`

When using `WITH DOCKER`, it is recommended that you use the official `earthly/dind` image (preferably `:alpine`) for running Docker-in-Docker. Earthly's `WITH DOCKER` requires that the Docker engine is installed already in the image it is running in.

If Docker engine is not detected, `WITH DOCKER` will need to first install it - it usually does so automatically - however, the cache will be inefficient. Consider the following example:

```Dockerfile
# Avoid
integration-test:
    FROM some-other-image:latest
    COPY docker-compose.yml ./
    WITH DOCKER --compose docker-compose.yml
        RUN ...
    END
```

Let's assume that `some-other-image:latest` does not already have Docker engine installed. This means that on the `WITH DOCKER` line, Earthly will add a hidden installation step, to add Docker engine. This takes some time to execute, but it will work.

The problem, however, will be apparent when there is a change (no matter how small) to `docker-compose.yml`. That will cause the build to re-execute without cache from the `COPY` command onwards, meaning that the installation of Docker engine will be repeated.

A simple way to fix this is to use an earthly-provided [function](../guides/functions.md) to install Docker engine before the `COPY` command. Please note that this particular function is fastest when ran on top of an alpine-based image.

```Dockerfile
# Better
integration-test:
    FROM some-other-image:latest
    DO github.com/earthly/lib+INSTALL_DIND
    COPY docker-compose.yml ./
    WITH DOCKER --compose docker-compose.yml
        RUN ...
    END
```

The best supported option, however, is to use the `earthly/dind` image, if possible.

```Dockerfile
# Best - if possible
integration-test:
    FROM earthly/dind:alpine-3.19-docker-25.0.3-r2
    COPY docker-compose.yml ./
    WITH DOCKER --compose docker-compose.yml
        RUN ...
    END
```

### Pattern: Saving artifacts resulting from a `WITH DOCKER`

In Earthly, `WITH DOCKER` starts up a transient Docker daemon for that specific instruction and then shuts it down and completely wipes its data afterwards. That does not mean, however, that you cannot export any information from it (or from any container running within), to be used in another part of the build. Although you may not run any non-`RUN` commands within `WITH DOCKER`, you can still use `SAVE ARTIFACT` (and any other command) after the `WITH DOCKER` instruction. The Docker daemon's data is wiped - but the rest of the build environment remains intact.

```Dockerfile
WITH DOCKER ...
    RUN docker run -v ./screenshots:/screenshots ... && \
        docker logs ... >./full-logs.txt && \
        docker inspect ... >./some-docker-state.json
END
SAVE ARTIFACT ./screenshots
SAVE ARTIFACT ./full-logs.txt
SAVE ARTIFACT ./some-docker-state.json
```

## Usage-specific

### Use `--ci` when running in CI

Build scripts serve slightly different purposes when they are run in CI compared to when they are executed for local development. Most of the logic is similar (hence Earthly attempts to unify the two concepts in a single syntax), but there are some small differences.

For example, for development purposes, you may use commands such as `LOCALLY`, which cause Earthly to be less repeatable, and yet might satisfy very much needed use-cases that are typically out of scope of a CI build.

In addition, in CI it is much more likely that shared caching will be needed, while outputting artifacts and images in the CI's host environment itself would not be needed.

For these reasons, Earthly comes with the `--ci` flag, which simply expands to `--no-output --use-inline-cache --save-inline-cache --strict`. The `--ci` flag therefore, prevents the use of commands that are not repeatable, enables inline caching and disables outputting artifacts and images on the CI host.

### Avoid `LOCALLY` and other non-strict commands

Certain Earthly functionality is only meant to be used for local development only. Most such commands do not fully abide by the Earthly spirit of repeatable builds, however, for certain specific development use-cases they are needed and therefore Earthly provides them. When Earthly is used in `--strict` mode (`earthly --strict +my-target` or `earthly --ci +my-target`), the usage of these commands is not allowed.

For this reason, it is recommended to avoid using these commands as much as possible, as doing so will:

1. Cause Earthly to behave in a non-repeatable way across other platforms, as it will rely on host-specific environment configuration.
2. Disable caching.
3. Cause the specific targets to not work at all when `--ci` is passed in.

An example of a command that is not allowed in strict mode is `LOCALLY`. The `LOCALLY` command skips the sandboxing of the build and executes all commands directly on the host machine.

Examples of valid cases where `LOCALLY` may be used are:

* installing dependencies on the host machine (e.g. to help IDEs provide better suggestions)
* executing tests on the host docker daemon, to help with inspection and debugging
* executing development commands which would otherwise require copying very large amounts of files to the sandboxed build environment

Note, however, that none of these cases are needed in a CI environment, and ultimately these commands are not regularly tested by a CI, which means they may break more frequently.

### Pattern: Push on the `main` branch only

When using Earthly in a sandboxed CI, it may be useful to perform pushes on occasion, in order to populate shared caches. In addition, image pushes might also help developers to grab pre-built images for local use for various development workflows.

Pushing too often can result in slowing down builds due to the upload time, while pushing too infrequently results in stale cache or stale development images.

A good balance is often to perform pushes on the `main` branch only, and to disable any pushing on PR builds. Although the `main` build will be slower, it will allow for maximum use of cache in PR builds, without the slowdown of further pushes.

Main branch build: `earthly --ci --push +target`.
PR build: `earthly --ci +target`.

The push option can also be configured via the env var `EARTHLY_PUSH`, which may be easier to manipulate in your CI of choice.

A more extreme case of this idea can be to use explicit maximum cache: `earthly --ci --push --remote-cache=.... --max-remote-cache +target`. The idea, again is to tradeoff performance on the `main` branch, for the benefit of faster PR builds. Whether this is actually beneficial needs to be measured on a project-by-project basis, however.

### Do not expose cache image tags publicly if the cache contains private code or dependencies

When using shared caching, Earthly has the ability to save some intermediate information about the build in a Docker registry of choice. Note, however, that the intermediate cache may also contain private code or dependencies, which could be exposed via the cache in some cases. Even if the final image only contains compiled binaries, it may still be possible to access intermediate layers that lead to the fully compiled binaries - some of which may contain private code or dependencies.

As such, when working on a closed-source project, it is advisable to only use private image repositories to prevent any code leaks.

### Technique: Use `earthly -i` to debug failures

In Earthly it is possible to drop into the container of a failed step to diagnose the failure better. To turn on this setting you can use the `-i` flag: `earthly -i +my-target`. Once dropped into the container's shell, you may use the up arrow to pre-populate the previously failed command should you wish to retry it or amend it. To exit the shell, press `Ctrl+D`.

### Run everything in a single Earthly invocation, do not wrap Earthly

Historically, build scripts have been constructed by cobbling up multiple technologies together: Makefiles, Bash scripts, Dockerfiles, Python scripts, Ruby scripts, and so on. The possibilities are endless, but also the readability and maintainability of the scripts suffer.

Earthly has been designed with a few key goals in mind:

* Repeatability - the builds should just work on another system
* Readability - the builds should be understandable by any team member on the team, without much effort
* A universal CI script - a script that contains all the information needed for the CI to perform a complete build

In this spirit, Earthly has been designed to not require any wrapping around it. Here are some examples of antipatterns:

* **Antipattern**: Earthly is called repeatedly in a single script in order to process intermediate results and then pass those results to later invocations.
* **Antipattern**: Downloading dependencies outside of Earthly and then copying them in.
* **Antipattern**: Performing a build without `--push` and then repeating the same build, but with `--push` enabled.
* **Antipattern**: Computing the value of an `ARG` outside of Earthly and then calling Earthly with that value.
* **Antipattern**: Running `earthly` in a bash `for` loop, to process multiple targets as separate builds.
* **Antipattern**: Running `earthly` repeatedly, rather than using an `all` target encapsulating all the targets needed to be built.
* **Antipattern**: Running `earthly` repeatedly with `--no-cache`, to control the order of a deployment, instead of using `RUN --push` adequately.

All of the above should be avoided as they hinder repeatability and/or they are abusing Earthly features in ways Earthly was not designed for. If a wrapping script is used outside of Earthly, it means that the script is not containerized, which means that the script is susceptible to host environment nuances.

The differences can be somewhat surprising: `make` and `sed` can be different on a Mac compared to a Linux machine, for example. Various linux distributions might have different versions of `bash` installed. Environment variables could play surprising roles. Causes for inconsistencies can sip in from anywhere, making builds more difficult to maintain.

To keep your build scripts uniform across projects (and thus more readable) and to keep them repeatable, it is best if `earthly` is used directly and with minimal argument overrides.

### Use `RUN --ssh` for passing host SSH keys to builds

Earthly provides a way to pass-through access to your host's SSH keys to the build, by proxying the host's ssh-agent connection inside the build. This may be useful if you need to access private repositories where you authenticate with SSH. An example of such a case might be downloading Go dependencies from private repositories:

```Dockerfile
RUN --ssh go mod download
```

### Use SAVE IMAGE to always cache

The simplicity of Earthly comes from its ability to cache build steps that do not need to be rerun. There are times, though, when you have more insight into the build process than Earthly does and may need to manually manage the cache for certain steps. Using an image push is an advanced technique that lets you tell Earthly, **"I want to always cache this step."**

For example, consider the Earthly Blog's Earthfile, which has a base image that installs Pandoc, among other utilities:

```Dockerfile
base:
    FROM ruby:2.7
    ARG TARGETARCH
    WORKDIR /site
    ...
    RUN apt-get install pandoc  pandocfilters -y

build:
    FROM +base
    # Build blog
```

Typically, installing Pandoc from source can take up to 30 minutes. When cached, the blog build process takes about 5 minutes, but without the cache, it jumps to 35 minutes. Occasionally, due to the heavy disk usage of the blog build, the base image may be evicted from the cache, leading to a slower build.

To circumvent this, the solution is to push the base step as an image and reference it directly in subsequent builds:

```Dockerfile
base:
    FROM ruby:2.7
    ARG TARGETARCH
    WORKDIR /site
    ...
    RUN apt-get install pandoc  pandocfilters -y

    # Manually cache the base image by pushing it to a registry
    SAVE IMAGE –push earthly/blog-base-image:latest

build:
    # Use the cached base image for builds  
    FROM earthly/blog-base-image:latest
    RUN ...
```

By doing this, we ensure that the build process is consistently swift, as the base layers do not need to be rebuilt. The trade-off is that the base target needs to be updated manually. For the Earthly website, we address this by running earthly `+base` once a week via a CI job, and on-demand whenever there are changes to the base image. You can see the full example, including flags for manual rerunning the base image, on GitHub.

#### Caveats to this approach

- Caching externally in an image registry adds download time to the build. This technique is most effective when the build time saved is significantly greater than the potential time to download the layer from a registry.

- By manually caching part of the build, you risk introducing inconsistencies. It's important to rerun the step that produces the image manually whenever there are changes. Scheduling a regular job can mitigate this risk.

#### Effective use cases

- **Expensive build steps that rarely change.** Earthly's cache operates on a least-recently-used (LRU) basis, trying to retain frequently used cache steps. However, some steps are more costly to regenerate than others. An external image can serve as a reliable fallback if an expensive step is evicted from the cache.

- **Benign Nondeterminism:** Earthly’s cache expects determinism. It assumes that if an input file changes, the associated build step must be rerun. If you know a step can be cached despite changes in inputs, using SAVE IMAGE for caching allows you to manage this manually and avoid unnecessary cache invalidation.

### Future: Saving an artifact even if the build fails

We are aware of the lack of capability here. Please follow GitHub issues [#988](https://github.com/earthly/earthly/issues/988) and [#587](https://github.com/earthly/earthly/issues/587) for updates.

There are currently workarounds for this (see [this comment](https://github.com/earthly/earthly/issues/988#issuecomment-870504677) and [this comment](https://github.com/earthly/earthly/issues/988#issuecomment-981088796)), however they have significant limitations.
