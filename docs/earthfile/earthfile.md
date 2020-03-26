# Earthfile reference

## FROM

#### Synopsis

* `FROM <image-name>`
* `FROM [--build-arg <key>=<value>] <target-ref>`

#### Description

The `FROM` command initializes a new build environment and sets the base image for subsequent instructions. It works similarly to the classical [Dockerfile `FROM` instruction](https://docs.docker.com/engine/reference/builder/#from), but it has the added ability to use another target's image as the base image for the build. For example: `FROM +another-target`.

> ##### Note
> The `FROM ... AS ...` form available in the classical Dockerfile syntax is not supported in Earthfiles. Instead, define a new Earthly target. For example, the following Dockerfile
> 
> ```Dockerfile
>     FROM alpine:3.11 AS build
>     # ... instructions for build
>     FROM build as another
>     # ... further instructions inheriting build
>     FROM busybox as yet-another
>     COPY --from=build ./a-file ./
> ```
>
> can become
>
> ```Dockerfile
>     build:
>         FROM alpine:3.11
>         # ... instructions for build
>         SAVE ARTIFACT ./a-file
>         SAVE IMAGE
>     another:
>         FROM +build
>         # ... further instructions inheriting build
>     yet-another:
>         FROM busybox
>         COPY +build/a-file ./
> ```

#### Options

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`. The value may be a constant string

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

## RUN

#### Synopsis

* `RUN [--push] [--entrypoint] [--privileged] [--with-docker] [--secret <env-var>=<secret-ref>] [--mount <mount-spec>] [--] <command>` (shell form)
* `RUN [[<flags>...], "<executable>", "<arg1>", "<arg2>", ...]` (exec form)

#### Description

The `RUN` command executes commands in the build environment of the current target, in a new layer.

The command allows for two possible forms. The *exec form* runs the command executable without the use of a shell. The *shell form* uses the default shell (`/bin/sh -c`) to interpret the command and execute it. In either form, you can use a `\` to continue a single `RUN` instruction onto the next line.

When the `--entrypoint` flag is used, the current image entrypoint is used to prepend the current command.

To avoid any abiguity regarding whether an argument is a `RUN` flag option or part of the command, the delimiter `--` may be used to signal the parser that no more `RUN` flag options will follow.

#### Options

##### `--push`

Marks the command as a "push command". Push commands are only executed if all other non-push instructions succeed. In addition, push commands are never cached, thus they are executed on every applicable invocation of the build.

Push commands are not run by default. Add the `--push` flag to the `earth` invocation to enable pushing. Example:

```bash
earth --push +deploy
```

Push commands were introduced to allow the user to define commands that have an effect external to the build. This kind of effects are only allowed to take place if the entire build succeeds. Good candidates for push commands are uploads of artifacts to artifactories, commands that make a change to an external environment, like a production or staging environment.

Note that no other non-push commands are allowed to follow a push command within a recipe.

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

Makes available a docker daemon within the build environment, which the command can leverage. This option automatically implies `--privileged`.

Example:

```Dockerfile
RUN --with-docker docker run hello-world
```

or 

```Dockerfile
RUN --with-docker docker-compose up -d ; docker run some-test-image
```

> ##### Note
> This feature is experimental and is subject to change.

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

The command may take a couple of possible forms. In the *classical form*, `COPY` copies files and directories from the build context into the build environment. In the *artifact form*, `COPY` copies files or directories (also known as "artifacts" in this context) from the artifact environment of other build targets into the build environment of the current target. Either form allows the use of wildcards for the sources.

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

Sets a value override of `<value>` for the build arg identified by `<key>`, when building the target containing the mentioned artifact. The value may be a constant string

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

## SAVE ARTIFACT

#### Synopsis

* `SAVE ARTIFACT <src> [<artifact-dest-path>] [AS LOCAL <local-path>]`

#### Description

The command `SAVE ARTIFACT` copies a file, a directory, or a series of files and directories represented by a wildcard, from the build environment into the target's artifact environment.

If `AS LOCAL ...` is also specified, it additionally marks the artifact to be copied to the host at the location specified by `<local-path>`, once the build is deemed as successful.

If `<artifact-dest-path>` is not specified, it is inferred as `/`.

File within the artifact environment are also known as "artifacts". Once a file has been copied into the artifact environment, it can be referenced in other places of the build (for example in a `COPY` command), using an [artifact reference](../guides/target-ref.md).

## SAVE IMAGE

#### Synopsis

* `SAVE IMAGE [[--push] <image-name>...]`

#### Description

The command `SAVE IMAGE` marks the current build environment as the image of the target. The image can then be referenced using an [image reference](../guides/target-ref.md) in other parts of the build (for example in a `FROM` command).

If one ore more `<image-name>`>'s are specified, the command also marks the image to be loaded within the docker daemon available on the host.

> ##### Note
> It is an error to issue the command `SAVE IMAGE` twice within the same recipe. The `SAVE IMAGE` command is always implied at the end of the `base` target, thus issuing `SAVE IMAGE` within the recipe of the `base` target is also an error.

#### Options

##### `--push`

The `--push` options marks the image to be pushed to an external registry after it has been loaded within the docker daemon available on the host.
