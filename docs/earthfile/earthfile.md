# Earthfile reference

Earthfiles are comprised of a series of target declarations and recipe definitions. Each recipe contains a series of commands, which are defined below. For an introduction into Earthfiles, see the [Basics page](../guides/basics.md).

## FROM

#### Synopsis

* `FROM <image-name>`
* `FROM [--build-arg <key>=<value>] <target-ref>`

#### Description

The `FROM` command initializes a new build environment and sets the base image for subsequent instructions. It works similarly to the classical [Dockerfile `FROM` instruction](https://docs.docker.com/engine/reference/builder/#from), but it has the added ability to use another target's image as the base image for the build. For example: `FROM +another-target`.

{% hint style='info' %}
##### Note

The `FROM ... AS ...` form available in the classical Dockerfile syntax is not supported in Earthfiles. Instead, define a new Earthly target. For example, the following Dockerfile
 
```Dockerfile
# Dockerfile

FROM alpine:3.11 AS build
# ... instructions for build

FROM build as another
# ... further instructions inheriting build

FROM busybox as yet-another
COPY --from=build ./a-file ./
```

can become

```Dockerfile
# build.earth

build:
    FROM alpine:3.11
    # ... instructions for build
    SAVE ARTIFACT ./a-file
    SAVE IMAGE

another:
    FROM +build
    # ... further instructions inheriting build

yet-another:
    FROM busybox
    COPY +build/a-file ./
```
{% endhint %}

#### Options

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`. See also [BUILD](#build) for more details about the `--build-arg` option.

## RUN

#### Synopsis

* `RUN [--push] [--entrypoint] [--privileged] [--with-docker] [--secret <env-var>=<secret-ref>] [--mount <mount-spec>] [--] <command>` (shell form)
* `RUN [[<flags>...], "<executable>", "<arg1>", "<arg2>", ...]` (exec form)

#### Description

The `RUN` command executes commands in the build environment of the current target, in a new layer. It works similarly to the [Dockerfile `RUN` command](https://docs.docker.com/engine/reference/builder/#run), with some added options.

The command allows for two possible forms. The *exec form* runs the command executable without the use of a shell. The *shell form* uses the default shell (`/bin/sh -c`) to interpret the command and execute it. In either form, you can use a `\` to continue a single `RUN` instruction onto the next line.

When the `--entrypoint` flag is used, the current image entrypoint is used to prepend the current command.

To avoid any abiguity regarding whether an argument is a `RUN` flag option or part of the command, the delimiter `--` may be used to signal the parser that no more `RUN` flag options will follow.

#### Options

##### `--push`

Marks the command as a "push command". Push commands are only executed if all other non-push instructions succeed. In addition, push commands are never cached, thus they are executed on every applicable invocation of the build.

Push commands are not run by default. Add the `--push` flag to the `earth` invocation to enable pushing. For example

```bash
earth --push +deploy
```

Push commands were introduced to allow the user to define commands that have an effect external to the build. This kind of effects are only allowed to take place if the entire build succeeds. Good candidates for push commands are uploads of artifacts to artifactories, commands that make a change to an external environment, like a production or staging environment.

Note that non-push commands are not allowed to follow a push command within a recipe.

##### `--entrypoint`

Prepends the currently defined entrypoint to the command.

This option is useful for replacing `docker run` in a traditional build environment. For example, a command like

```bash
docker run --rm -v "$(pwd):/data" cytopia/golint .
```

Might become the following in an Earthfile

```Dockerfile
FROM cytopia/goling
COPY . /data
RUN --entrypoint .
```

##### `--privileged`

Allows the command to use privileged capabilities.

Note that privileged mode is not enabled by default. In order to use this option, you need to additionally pass the flag `--allow-privileged` (or `-P`) to the `earth` command. Example:

```bash
earth --allow-privileged +some-target
```

##### `--with-docker` [**experimental**]

{% hint style='danger' %}
`RUN --with-docker` is experimental and is subject to change.
{% endhint %}

Makes available a docker daemon within the build environment, which the command can leverage. This option automatically implies `--privileged`.

Example:

```Dockerfile
RUN --with-docker docker run hello-world
```

or 

```Dockerfile
RUN --with-docker docker-compose up -d ; docker run some-test-image
```

##### `--secret <env-var>=<secret-ref>`

Makes available a secret, in the form of an env var (its name is defined by `<env-var>`), to the command being executed.

The `<secret-ref>` needs to be of the form `+secrets/<secret-id>`, where `<secret-id>` is the identifier passed to the `earth` command when passing the secret: `earth --secret <secret-id>=<value>`.

Here is an example:

```Dockerfile
release:
    RUN --push --secret GITHUB_TOKEN=+secrets/GH_TOKEN github-release upload
```

```bash
earth --secret GH_TOKEN="the-actual-secret-token-value" +release
```

##### `--mount <mount-spec>`

Mounts a file or directory in the context of the build environment.

The `<mount-spec>` is defined as a series of comma-separated list of key-values. The following keys are allowed

| Key | Description | Example |
| --- | --- | --- |
| `type` | The type of the mount. Currently only `cache` is allowed. | `type=cache` |
| `target` | The target path for the mount. | `target=/var/lib/data` |

Example:

```Dockerfile
ENV GOCACHE=/go-cache
RUN --mount=type=cache,target=/go-cache go build main.go
```

Note that mounts cannot be shared between targets, nor can it be shared within the same target,
if the build-args differ between invocations.

## COPY

#### Synopsis

* `COPY [--dir] <src>... <dest>` (classical form)
* `COPY [--dir] [--build-arg <key>=<value>] <src-artifact>... <dest>` (artifact form)

#### Description

The command `COPY` allows copying of files and directories between different contexts.

The command may take a couple of possible forms. In the *classical form*, `COPY` copies files and directories from the build context into the build environment - in this form, it works similarly to the [Dockerfile `COPY` command](https://docs.docker.com/engine/reference/builder/#copy). In the *artifact form*, `COPY` copies files or directories (also known as "artifacts" in this context) from the artifact environment of other build targets into the build environment of the current target. Either form allows the use of wildcards for the sources.

The parameter `<src-artifact>` is an artifact reference and is generally of the form `<target-ref>/<artifact-path>`, where `<target-ref>` is the reference to the target which needs to be built in order to yield the artifact and `<artifact-path>` is the path within the artifact environment of the target, where the file or directory is located. The `<artifact-path>` may also be a wildcard.

Note that the Dockerfile form of `COPY` whereby you can reference a source as a URL is not yet supported in Earthfiles.

#### Options

##### `--dir`

The option `--dir` changes the behavior of the `COPY` command to copy the directories themselves, rather than the contents of the directories. It allows the command to behave similarly to a `cp -r` operation on a unix system. This allows the enumeration of several directories to be copied over on a single line (and thus, within a single layer). For example, the following two are equivalent with respect to what is being copied in the end (but not equivalent with respect to the number of layers used).

```Dockerfile
COPY dir1 dir1
COPY dir2 dir2
COPY dir3 dir3
```

```Dockerfile
COPY --dir dir1 dir2 dir3 ./
```

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`, when building the target containing the mentioned artifact. See also [BUILD](#build) for more details about the `--build-arg` option.

##### `--from`

Although this option is present in classical Dockerfile syntax, it is not supported by Earthfiles. You may instead use a combination of `SAVE ARTIFACT` and `COPY` *artifact form* commands to achieve similar effects. For example, the following Dockerfile

```Dockerfile
# Dockerfile
COPY --from=some-image /path/to/some-file.txt ./
```

... would be equivalent to `final-target` in the following Earthfile

```Dockerfile
# build.earth
intermediate:
    FROM some-image
    SAVE ARTIFACT /path/to/some-file.txt

final-target:
    COPY +intermediate/some-file.txt ./
```

## GIT CLONE

#### Synopsis

* `GIT CLONE [--branch <git-ref>] <git-url> <dest-path>`

#### Description

The command `GIT CLONE` clones a git repository from `<git-url>`, optionally referenced by `<git-ref>`, into the build environment, within the `<dest-path>`. If the `--branch` option is not specified, then `HEAD` is inferred.

In contrast to an operation like `RUN git clone <git-url> <dest-path>`, the command `GIT CLONE` is cache-aware and correctly distinguishes between different git commit IDs when deciding to reuse a previous cache or not. In addition, `GIT CLONE` can also use [Git authentication configuration](../guides/auth.md) passed on to `earth`, whereas `RUN git clone` would require additional secrets passing, if the repository is not publicly accessible.

#### Options

##### `--branch <git-ref>`

Points the `HEAD` to the git reference specified by `<git-ref>`. If this option is not specified, then the remote `HEAD` is used instead.

## SAVE ARTIFACT

#### Synopsis

* `SAVE ARTIFACT <src> [<artifact-dest-path>] [AS LOCAL <local-path>]`

#### Description

The command `SAVE ARTIFACT` copies a file, a directory, or a series of files and directories represented by a wildcard, from the build environment into the target's artifact environment.

If `AS LOCAL ...` is also specified, it additionally marks the artifact to be copied to the host at the location specified by `<local-path>`, once the build is deemed as successful.

If `<artifact-dest-path>` is not specified, it is inferred as `/`.

Files within the artifact environment are also known as "artifacts". Once a file has been copied into the artifact environment, it can be referenced in other places of the build (for example in a `COPY` command), using an [artifact reference](../guides/target-ref.md).

## SAVE IMAGE

#### Synopsis

* `SAVE IMAGE [[--push] <image-name>...]`

#### Description

The command `SAVE IMAGE` marks the current build environment as the image of the target. The image can then be referenced using an [image reference](../guides/target-ref.md) in other parts of the build (for example in a `FROM` command).

If one ore more `<image-name>`>'s are specified, the command also marks the image to be loaded within the docker daemon available on the host.

{% hint style='info' %}
##### Note
It is an error to issue the command `SAVE IMAGE` twice within the same recipe. In addition, the `SAVE IMAGE` command is always implied at the end of the `base` target, thus issuing `SAVE IMAGE` within the recipe of the `base` target is also an error.
{% endhint %}

#### Options

##### `--push`

The `--push` options marks the image to be pushed to an external registry after it has been loaded within the docker daemon available on the host.

Push commands are not run by default. Add the --push flag to the earth invocation to enable pushing. For example

```bash
earth --push +docker-image
```

## BUILD

#### Synopsis

* `BUILD [--build-arg <key>=<value>] <target-ref>`

#### Description

The command `BUILD` instructs Earthly to additionally invoke the build of the target referenced by `<target-ref>`, where `<target-ref>` follows the rules defined by [target referencing](../guides/target-ref.md).

#### Options

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`.

The override value of a build arg may be a constant string

```
--build-arg SOME_ARG="a constant value"
```

or an expression involving other build args

```
--build-arg SOME_ARG="a value based on other args, like $ANOTHER_ARG and $YET_ANOTHER_ARG"
```

or a dynamic expression, based on the output of a command executed in the context of the build environment. In this case, the build arg becomes a "variable build arg".

```
--build-arg SOME_ARG=$(find /app -type f -name '*.php')
```

## ARG

#### Synopsis

* `ARG <name>[=<default-value>]`

#### Description

The command `ARG` declares a variable (or arg) with the name `<name>` and with an optional default value `<default-value>`. If no default value is provided, then empty string is used as the default value.

This command works similarly to the [Dockerfile `ARG` command](https://docs.docker.com/engine/reference/builder/#arg), with a few differences regarding the scope and the predefined args (called builtin args in Earthly). For more information see [builtin args](../guides/builtin-args.md). The variable's scope is always limited to the current target's recipe and only from the point it is declared onwards.

The value of an arg can be overridden either from the `earth` command

```bash
earth --build-arg <name>=<override-value>
```

or from a command from another target, when implicitly or explicitly invoking the target containing the `ARG`

```Dockerfile
BUILD --build-arg <name>=<override-value> <target-ref>
COPY --build-arg <name>=<override-value> <target-ref>/<artifact-path>... <dest-path>
FROM --build-arg <name>=<override-value> <target-ref>
```

for example

```Dockerfile
BUILD --build-arg PLATFORM=linux +binary
COPY --build-arg PLATFORM=linux +binary/bin ./
FROM --build-arg NAME=john +docker-image
```

{% hint style='danger' %}
##### Important

In contrast to Dockerfile predefined args, Earthly builtin args need to be pre-declared before they can be used. For example

```Dockerfile
ARG EARTHLY_TARGET
RUN echo "The current target is $EARTHLY_TARGET"
```
{% endhint %}

The value of a builtin arg can never be overriden. However, you can always have an additional `ARG`, which takes as the default value, the value of the builtin arg. The additional arg can be overriden. Example

```Dockerfile
ARG EARTHLY_TARGET_TAG
ARG TAG=$EARTHLY_TARGET_TAG
SAVE IMAGE --push some/name:$TAG
```

The following builtin args are available

| Name | Description | Example value |
| --- | --- | --- |
| `EARTHLY_TARGET` | The canonical reference of the current target. | For example, for a target named `foo`, which exists on `master` branch, in a repository at `github.com/bar/buz`, in a subdirectory `src`, the canonical reference would be `github.com/bar/buz/src:master+foo`. For more information about canonical references, see [target referencing](../guides/target-ref.md). |
| `EARTHLY_TARGET_PROJECT` | The project part of the canonical reference of the current target. | For the example above, the canonical project would be `github.com/bar/buz/src` |
| `EARTHLY_TARGET_NAME` | The name part of the canonical reference of the current target. | For the example above, the name would be `foo` |
| `EARTHLY_TARGET_TAG` | The tag part of the canonical reference of the current target. Note that in some cases, no tag is detected, and as such, the value is an empty string | For the example above, the tag would be `master` |
| `EARTHLY_GIT_HASH` | The git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666ade67b29e42bef121144456d3977a67a` |
| `EARTHLY_GIT_ORIGIN_URL` | The git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `git@github.com:vladaionescu/earthly.git` |
| `EARTHLY_GIT_PROJECT_NAME` | The git project name from within the git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `vladaionescu/earthly` |

## DOCKER PULL [**experimental**]

{% hint style='danger' %}
`DOCKER PULL` is experimental and is subject to change.
{% endhint %}

#### Synopsis

* `DOCKER PULL <image-name>`

#### Description

The command `DOCKER PULL` pulls a docker image from a remote registry into the docker daemon available within the build envionment. It can be used in conjunction with `RUN --with-docker docker run ...` to execute docker images in the context of the build environment.

## DOCKER LOAD [**experimental**]

{% hint style='danger' %}
`DOCKER LOAD` is experimental and is subject to change.
{% endhint %}

#### Synopsis

* `DOCKER LOAD [--build-arg <name>=<override-value>] <target-ref> AS <image-name>`

#### Description

The command `DOCKER LOAD` builds the image referenced by `<target-ref>` and then loads it into the docker daemon available within the build environment, as a docker image `<image-name>`. It can be used in conjunction with `RUN --with-docker docker run ...` to execute docker images that are produced by other targets of the build.

#### Options

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`, when invoking the build referenced by `<target-ref>`. See also [BUILD](#build) for more details about the `--build-arg` option.

## CMD (same as Dockerfile CMD)

#### Synopsis

* `CMD ["executable", "arg1", "arg2"]` (exec form)
* `CMD ["arg1, "arg2"]` (as default arguments to the entrypoint)
* `CMD command arg1 arg2` (shell form)

#### Description

The command `CMD` sets default arguments for an image, when executing as a container. It works the same way as the [Dockerfile `CMD` command](https://docs.docker.com/engine/reference/builder/#cmd).

## LABEL (same as Dockerfile LABEL)

#### Synopsis

* `LABEL <key>=<value> <key>=<value> ...`

#### Description

The `LABEL` command adds label metadata to an image. It works the same way as the [Dockerfile `LABEL` command](https://docs.docker.com/engine/reference/builder/#label).

## EXPOSE (same as Dockerfile EXPOSE)

#### Synopsis

* `EXPOSE <port> <port> ...`
* `EXPOSE <port>/<protocol> <port>/<protocol> ...`

#### Description

The `EXPOSE` command marks a series of ports as listening ports within the image. It works the same way as the [Dockerfile `EXPOSE` command](https://docs.docker.com/engine/reference/builder/#expose).

## ENV (same as Dockerfile ENV)

#### Synopsis

* `ENV <key> <value>`
* `ENV <key>=<value>`

#### Description

The `ENV` command sets the environment variable `<key>` to the value `<value>`. It works the same way as the [Dockerfile `ENV` command](https://docs.docker.com/engine/reference/builder/#env).

{% hint style='info' %}
##### Note
Do not use the `ENV` command for secrets used during the build. All `ENV` values used during the build are persisted within the image itself. See the [`RUN --secret` option](#run) to pass secrets to build instructions.
{% endhint %}

## ENTRYPOINT (same as Dockerfile ENTRYPOINT)

#### Synopsis

* `ENTRYPOINT ["executable", "arg1", "arg2"]` (exec form)
* `ENTRYPOINT command arg1 arg2` (shell form)

#### Description

The `ENTRYPOINT` command sets the default command or executable to be run when the image is executed as a container. It works the same way as the [Dockerfile `ENTRYPOINT` command](https://docs.docker.com/engine/reference/builder/#entrypoint).

## VOLUME (same as Dockerfile VOLUME)

#### Synopsis

* `VOLUME <path-to-target-mount> <path-to-target-mount> ...`
* `VOLUME ["<path-to-target-mount>", <path-to-target-mount> ...]`

#### Description

The `VOLUME` command creates a mount point at the specified path and marks it as holding externally mounted volumes. It works the same way as the [Dockerfile `VOLUME` command](https://docs.docker.com/engine/reference/builder/#volume).

## USER (same as Dockerfile USER)

#### Synopsis

* `USER <user>[:<group>]`
* `USER <UID>[:<GID>]`

#### Description

The `USER` command sets the user name (or UID) and optionally the user group (or GID) to use when running the image and also for any subsequent instructions in the build recipe. It works the same way as the [Dockerfile `USER` command](https://docs.docker.com/engine/reference/builder/#user).

## WORKDIR (same as Dockerfile WORKDIR)

#### Synopsis

* `WORKDIR <path-to-dir>`

#### Description

The `WORKDIR` command strs the working directory for other commands that follow in the recipe. The working directory is also persisted as the default directory for the image. If the directory does not exist, it is automatically created. This command works the same way as the [Dockerfile `WORKDIR` command](https://docs.docker.com/engine/reference/builder/#workdir).

## SHELL (not supported)

The classical [`SHELL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#add) is not yet supported. Use the *exec form* of `RUN`, `ENTRYPOINT` and `CMD` instead and prepend a different shell.

## ADD (not supported)

The classical [`ADD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#add) is not yet supported. Use [COPY](#copy) instead.

## ONBUILD (not supported)

The classical [`ONBUILD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#onbuild) is not supported.

## STOPSIGNAL (not supported)

The classical [`STOPSIGNAL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#stopsignal) is not yet supported.

## HEALTHCHECK (not supported)

The classical [`HEALTHCHECK` Dockerfile command](https://docs.docker.com/engine/reference/builder/#healthcheck) is not yet supported.
