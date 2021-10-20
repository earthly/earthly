# Earthfile reference

<!--

Note to person editing!!

The general order of the commands is as follows:

- Core classical Dockerfile commands (order is the same as in the Dockferfile official docs)
- Core, GA'd Earthly commands
- Other Dockerfile commands which have the exact same behavior in Earthly as in Dockerfiles
- Beta Earthly commands
- Experimental Earthly commands
- Classical Dockerfile commands that are not supported
- Deprecated Earthly commands

-->

Earthfiles are comprised of a series of target declarations and recipe definitions. Earthfiles are named `Earthfile`, regardless of their location in the codebase.

Earthfiles have the following rough structure:

```
<base-recipe>
...

<target-name>:
    <recipe>
    ...

<target-name>:
    <recipe>
    ...

<command-name>:
    <recipe>
    ...
```

Each recipe contains a series of commands, which are defined below. For an introduction into Earthfiles, see the [Basics page](../basics/basics.md).

## FROM

#### Synopsis

* `FROM <image-name>`
* `FROM [--build-arg <key>=<value>] [--platform <platform>] [--allow-privileged] <target-ref>`

#### Description

The `FROM` command initializes a new build environment and sets the base image for subsequent instructions. It works similarly to the classical [Dockerfile `FROM` instruction](https://docs.docker.com/engine/reference/builder/#from), but it has the added ability to use another target's image as the base image for the build. For example: `FROM +another-target`.

{% hint style='info' %}
##### Note

The `FROM ... AS ...` form available in the classical Dockerfile syntax is not supported in Earthfiles. Instead, define a new Earthly target. For example, the following Dockerfile
 
```Dockerfile
# Dockerfile

FROM alpine:3.13 AS build
# ... instructions for build

FROM build as another
# ... further instructions inheriting build

FROM busybox as yet-another
COPY --from=build ./a-file ./
```

can become

```Dockerfile
# Earthfile

build:
    FROM alpine:3.13
    # ... instructions for build
    SAVE ARTIFACT ./a-file

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

##### `--platform <platform>` (**beta**)

Specifies the platform to build on.

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Allows remotely-referenced targets to request privileged capabilities; this flag has no effect when referencing local targets.

Additionally, for privileged capabilities, earthly must be invoked on the command line with the `--allow-privileged` (or `-P`) flag.

For example, consider two Earthfiles, one hosted on a remote GitHub repo:

```Dockerfile
# github.com/earthly/example
FROM alpine:latest
elevated-target:
    RUN --privileged echo do something requiring privileged access.
```

and a local Earthfile:

```Dockerfile
FROM alpine:latest
my-target:
    FROM --allow-privileged github.com/earthly/example+elevated-target
    # ... further instructions inheriting remotely referenced Earthfile
```

then one can build `my-target` by invoking earthly with the `--allow-privileged` (or `-P`) flag:

```bash
earthly --allow-privileged +my-target
```

## RUN

#### Synopsis

* `RUN [--push] [--entrypoint] [--privileged] [--secret <env-var>=<secret-ref>] [--ssh] [--mount <mount-spec>] [--] <command>` (shell form)
* `RUN [[<flags>...], "<executable>", "<arg1>", "<arg2>", ...]` (exec form)

#### Description

The `RUN` command executes commands in the build environment of the current target, in a new layer. It works similarly to the [Dockerfile `RUN` command](https://docs.docker.com/engine/reference/builder/#run), with some added options.

The command allows for two possible forms. The *exec form* runs the command executable without the use of a shell. The *shell form* uses the default shell (`/bin/sh -c`) to interpret the command and execute it. In either form, you can use a `\` to continue a single `RUN` instruction onto the next line.

When the `--entrypoint` flag is used, the current image entrypoint is used to prepend the current command.

To avoid any ambiguity regarding whether an argument is a `RUN` flag option or part of the command, the delimiter `--` may be used to signal the parser that no more `RUN` flag options will follow.

#### Options

##### `--push`

Marks the command as a "push command". Push commands are only executed if all other non-push instructions succeed. In addition, push commands are never cached, thus they are executed on every applicable invocation of the build.

Push commands are not run by default. Add the `--push` flag to the `earthly` invocation to enable pushing. For example

```bash
earthly --push +deploy
```

Push commands were introduced to allow the user to define commands that have an effect external to the build. This kind of effects are only allowed to take place if the entire build succeeds. Good candidates for push commands are uploads of artifacts to artifactories, commands that make a change to an external environment, like a production or staging environment.

Note that non-push commands are not allowed to follow a push command within a recipe.

##### `--no-cache`

Force the command to run every time; ignoring any cache. Any commands following the invocation of `RUN --no-cache`, will also ignore the cache. If `--no-cache` is used as an option on the `RUN` statement within a `WITH DOCKER` statement, all commands after the `WITH DOCKER` will also ignore the cache.

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

Note that privileged mode is not enabled by default. In order to use this option, you need to additionally pass the flag `--allow-privileged` (or `-P`) to the `earthly` command. Example:

```bash
earthly --allow-privileged +some-target
```

##### `--secret <env-var>=<secret-ref> | <secret-id>`

Makes available a secret, in the form of an env var (its name is defined by `<env-var>`), to the command being executed. 
If you only specify `<secret-id>`, the name of the env var will be `<secret-id>` and its value the value of `<secret-id>`.

The `<secret-ref>` needs to be of the form `+secrets/<secret-id>`, where `<secret-id>` is the identifier passed to the `earthly` command when passing the secret: `earthly --secret <secret-id>=<value>`.

Here is an example that showcases both syntaxes:

```Dockerfile
release:
    RUN --push --secret GITHUB_TOKEN=+secrets/GH_TOKEN github-release upload
release-short:
    RUN --push --secret GITHUB_TOKEN github-release upload
```

```bash
earthly --secret GH_TOKEN="the-actual-secret-token-value" +release
earthly --secret GITHUB_TOKEN="the-actual-secret-token-value" +release-short
```

An empty string is also allowed for `<secret-ref>`, allowing for optional secrets, should it need to be disabled.

```Dockerfile
release:
    ARG SECRET_ID=+secrets/GH_TOKEN
    RUN --push --secret GITHUB_TOKEN=$SECRET_ID github-release upload
release-short:
    ARG SECRET_ID=GITHUB_TOKEN
    RUN --push --secret $SECRET_ID github-release upload
```

```bash
earthly --build-arg SECRET_ID="" +release
earthly --build-arg SECRET_ID="" +release-short
```

See also the [Cloud secrets guide](../guides/cloud-secrets.md).

##### `--ssh`

Allows a command to access the ssh authentication client running on the host via the socket which is referenced by the environment variable `SSH_AUTH_SOCK`.

Here is an example:

```Dockerfile
RUN mkdir -p ~/.ssh && \
    echo 'github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==' >> ~/.ssh/known_hosts && \
    echo 'gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9' >> ~/.ssh/known_hosts
RUN --ssh git config --global url."git@github.com:".insteadOf "https://github.com/" && \
    go mod download
```

##### `--mount <mount-spec>`

Mounts a file or directory in the context of the build environment.

The `<mount-spec>` is defined as a series of comma-separated list of key-values. The following keys are allowed

| Key | Description | Example |
| --- | --- | --- |
| `type` | The type of the mount. Currently only `cache`, `tmpfs`, and `secret` are allowed. | `type=cache` |
| `target` | The target path for the mount. | `target=/var/lib/data` |
| `id` | The secret ID for the contents of the `target` file, only applicable for `type=secret`. | `id=+secrets/password` |

Example:

```Dockerfile
ENV GOCACHE=/go-cache
RUN --mount=type=cache,target=/go-cache go build main.go
```

Note that mounts cannot be shared between targets, nor can they be shared within the same target,
if the build-args differ between invocations.

##### `--interactive` / `--interactive-keep` (**experimental**)

Opens an interactive prompt during the target build. An interactive prompt must:

1. Be the last issued command in the target, with the exception of `SAVE IMAGE` commands. This also means that you cannot `FROM` a target containing a `RUN --interactive`.
2. Be the only `--interactive` target within the run.
3. Not be within a `LOCALLY`-designated target.

###### Examples:

Start an interactive python REPL:
```Dockerfile
python:
    FROM alpine:3.13
    RUN apk add python
    RUN --interactive python
```

Start `bash` to tweak an image by hand. Changes made will be included:
```Dockerfile
build:
    FROM alpine:3.13
    RUN apk add bash
    RUN --interactive-keep bash
```

## COPY

#### Synopsis

* `COPY [options...] <src>... <dest>` (classical form)
* `COPY [options...] <src-artifact>... <dest>` (artifact form)

#### Description

The command `COPY` allows copying of files and directories between different contexts.

The command may take a couple of possible forms. In the *classical form*, `COPY` copies files and directories from the build context into the build environment - in this form, it works similarly to the [Dockerfile `COPY` command](https://docs.docker.com/engine/reference/builder/#copy). In the *artifact form*, `COPY` copies files or directories (also known as "artifacts" in this context) from the artifact environment of other build targets into the build environment of the current target. Either form allows the use of wildcards for the sources.

The parameter `<src-artifact>` is an [artifact reference](../guides/target-ref.md#artifact-reference) and is generally of the form `<target-ref>/<artifact-path>`, where `<target-ref>` is the reference to the target which needs to be built in order to yield the artifact and `<artifact-path>` is the path within the artifact environment of the target, where the file or directory is located. The `<artifact-path>` may also be a wildcard.

Note that the Dockerfile form of `COPY` whereby you can reference a source as a URL is not yet supported in Earthfiles.

{% hint style='info' %}
##### Note
To prevent Earthly from copying unwanted files, you may specify file patterns to be excluded from the build context using an [`.earthignore`](./earthignore.md) file. This file has the same syntax as a [`.dockerignore` file](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
{% endhint %}

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

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

##### `--keep-own`

Instructs Earthly to keep file ownership information. This applies only to the *artifact form* and has no effect otherwise.

##### `--if-exists`

Only copy source if it exists; if it does not exist, earthly will simply ignore the COPY command and won't treat any missing sources as failures.

##### `--from`

Although this option is present in classical Dockerfile syntax, it is not supported by Earthfiles. You may instead use a combination of `SAVE ARTIFACT` and `COPY` *artifact form* commands to achieve similar effects. For example, the following Dockerfile

```Dockerfile
# Dockerfile
COPY --from=some-image /path/to/some-file.txt ./
```

... would be equivalent to `final-target` in the following Earthfile

```Dockerfile
# Earthfile
intermediate:
    FROM some-image
    SAVE ARTIFACT /path/to/some-file.txt

final-target:
    COPY +intermediate/some-file.txt ./
```

##### `--platform <platform>` (**beta**)

In *artifact form*, it specifies the platform to build the artifact on.

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

#### Examples

Assuming the following directory tree, of a folder named `test`:

```
test
  └── file
```

Here is how the following copy commands will behave:

```
# Copies the contents of the test directory.
# To access the file, it would be found at ./file
COPY test .

# Also copies the contents of the test directory.
# To access the file, it would be found at ./file
COPY test/* .

# Copies the whole test folder.
# To access the file, it would be found at ./test/file
COPY --dir test .
```

For detailed examples demonstrating how other scenarios may function, please see our [test suite](https://github.com/earthly/earthly/blob/main/tests/copy.earth).

## ARG

#### Synopsis

* `ARG <name>[=<default-value>]`

#### Description

The command `ARG` declares a variable (or arg) with the name `<name>` and with an optional default value `<default-value>`. If no default value is provided, then empty string is used as the default value.

This command works similarly to the [Dockerfile `ARG` command](https://docs.docker.com/engine/reference/builder/#arg), with a few differences regarding the scope and the predefined args (called builtin args in Earthly). The variable's scope is always limited to the recipe of the current target or command and only from the point it is declared onward. For more information regarding builtin args, see the [builtin args page](./builtin-args.md).

If an `ARG` is defined in the `base` target of the Earthfile, then it becomes a global `ARG` and it is made available to every other target or command in that file, regardless of their base images used.

The value of an arg can be overridden either from the `earthly` command

```bash
earthly --build-arg <name>=<override-value>
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

A number of builtin args are available and are pre-filled by Earthly. For more information see [builtin args](./builtin-args.md).

## SAVE ARTIFACT

#### Synopsis

* `SAVE ARTIFACT [--keep-ts] [--keep-own] [--if-exists] [--force] <src> [<artifact-dest-path>] [AS LOCAL <local-path>]`

#### Description

The command `SAVE ARTIFACT` copies a file, a directory, or a series of files and directories represented by a wildcard, from the build environment into the target's artifact environment.

If `AS LOCAL ...` is also specified, it additionally marks the artifact to be copied to the host at the location specified by `<local-path>`, once the build is deemed as successful.

If `<artifact-dest-path>` is not specified, it is inferred as `/`.

Files within the artifact environment are also known as "artifacts". Once a file has been copied into the artifact environment, it can be referenced in other places of the build (for example in a `COPY` command), using an [artifact reference](../guides/target-ref.md#artifact-reference).

{% hint style='info' %}
##### Hint
In order to inspect the contents of an artifacts environment, you can run

```bash
earthly --artifact +<target>/* ./output/
```

This command dumps the contents of the artifact environment of the target `+<target>` into a local directory called `output`, which can be inspected directly.
{% endhint %}

{% hint style='danger' %}
##### Important
Note that there is a distinction between a *directory artifact* and *file artifact* when it comes to local output. When saving an artifact locally, a directory artifact will **replace** the destination entirely, while a file (or set of files) artifact will be copied **into** the destination directory.

```Dockerfile
# This will wipe ./destination and replace it with the contents of the ./my-directory artifact.
SAVE ARTIFACT ./my-directory AS LOCAL ./destination
# This will merge the contents of ./my-directory into ./destination.
SAVE ARTIFACT ./my-directory/* AS LOCAL ./destination
```
{% endhint %}

#### Options

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

##### `--keep-own`

Instructs Earthly to keep file ownership information.

##### `--if-exists`

Only save artifacts if they exists; if not, earthly will simply ignore the SAVE ARTIFACT command and won't treat any missing sources as failures.

##### `--force`

Force save operations which may be unsafe, such as writing to (or overwriting) a file or directory on the host filesystem located outside of the context of the directory containing the Earthfile.

#### Examples

Assuming the following directory tree, of a folder named `test`:

```
test
  └── file

```

Here is how the following `SAVE ARTIFACT ... AS LOCAL` commands will behave:

```
WORKDIR base
COPY test .

# This will copy the base folder into the output directory.
# You would find file at out-dot/base/file.
SAVE ARTIFACT . AS LOCAL out-dot/

# This will copy the contents of the base folder into the output directory.
# You would find sub-file at out-glob/file. Note the base directory is not in the output.
SAVE ARTIFACT ./* AS LOCAL out-glob/
```

For detailed examples demonstrating how other scenarios may function, please see our [test suite](https://github.com/earthly/earthly/blob/main/examples/tests/file-copying.earth).

## SAVE IMAGE

#### Synopsis

* `SAVE IMAGE [--cache-from=<cache-image>] [--push] <image-name>...` (output form)
* `SAVE IMAGE --cache-hint` (cache hint form)

#### Description

In the *output form*, the command `SAVE IMAGE` marks the current build environment as the image of the target and assigns an output image name.

In the *cache hint form*, it instructs Earthly that the current target should be included as part of the explicit cache. For more information see the [shared caching guide](../guides/shared-cache.md).

#### Options

##### `--push`

The `--push` options marks the image to be pushed to an external registry after it has been loaded within the docker daemon available on the host.

If inline caching is enabled, the `--push` option also instructs Earthly to use the specified image names as cache sources.

The actual push is not executed by default. Add the `--push` flag to the earthly invocation to enable pushing. For example

```bash
earthly --push +docker-image
```

##### `--cache-from=<cache-image>` (**experimental**)

Adds additional cache sources to be used when `--use-inline-cache` is enabled. For more information see the [shared caching guide](../guides/shared-cache.md).

##### `--cache-hint` (**experimental**)

Instructs Earthly that the current target should be included as part of the explicit cache. For more information see the [shared caching guide](../guides/shared-cache.md).

## BUILD

#### Synopsis

* `BUILD [--build-arg <key>=<value>] [--platform <platform>] [--allow-privileged] <target-ref>`

#### Description

The command `BUILD` instructs Earthly to additionally invoke the build of the target referenced by `<target-ref>`, where `<target-ref>` follows the rules defined by [target referencing](../guides/target-ref.md#target-reference).

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

##### `--platform <platform>` (**beta**)

Specifies the platform to build on.

This flag may be repeated in order to instruct the system to perform the build for multiple platforms. For example

```Dockerfile
build-all-platforms:
    BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build
```

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

## VERSION

#### Synopsis

* `VERSION [--use-copy-include-patterns] <version-number>`

#### Description

The command `VERSION` identifies which set of features to enable in Earthly while handling the corresponding Earthfile. The `VERSION` command is currently optional;
however will become mandatory in a future version of Earthly. When specified, `VERSION` must be the first command in the Earthfile.


| Version number | enabled features |
| --- | --- | --- |
| `0.5` | _initial functionality will be preserved_ |

#### Options

Individual features may be enabled by setting the corresponding feature flag.
New features start off as experimental, which is why they are disabled by default.
Once a feature reaches maturity, it will be enabled by default under a new version number.

All features are described in [a corresponding table](./features.md).

## GIT CLONE

#### Synopsis

* `GIT CLONE [--branch <git-ref>] [--keep-ts] <git-url> <dest-path>`

#### Description

The command `GIT CLONE` clones a git repository from `<git-url>`, optionally referenced by `<git-ref>`, into the build environment, within the `<dest-path>`.

In contrast to an operation like `RUN git clone <git-url> <dest-path>`, the command `GIT CLONE` is cache-aware and correctly distinguishes between different git commit IDs when deciding to reuse a previous cache or not. In addition, `GIT CLONE` can also use [Git authentication configuration](../guides/auth.md) passed on to `earthly`, whereas `RUN git clone` would require additional secrets passing, if the repository is not publicly accessible.

#### Options

##### `--branch <git-ref>`

Points the `HEAD` to the git reference specified by `<git-ref>`. If this option is not specified, then the remote `HEAD` is used instead.

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

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

The `WORKDIR` command `strs` the working directory for other commands that follow in the recipe. The working directory is also persisted as the default directory for the image. If the directory does not exist, it is automatically created. This command works the same way as the [Dockerfile `WORKDIR` command](https://docs.docker.com/engine/reference/builder/#workdir).

## HEALTHCHECK (same as Dockerfile HEALTHCHECK)

#### Synopsis

* `HEALTHCHECK NONE` (disable health checking)
* `HEALTHCHECK [--interval=DURATION] [--timeout=DURATION] [--start-period=DURATION] [--retries=N] CMD command arg1 arg2` (check container health by running command inside the container)

#### Description

The `HEALTHCHECK` command tells Docker how to test a container to check that it is still working. It works the same way as the [Dockerfile `HEALTHCHECK` command](https://docs.docker.com/engine/reference/builder/#healthcheck), with the only exception that the exec form of this command is not yet supported.

#### Options

##### `--interval=DURATION`

Sets the time interval between health checks. Defaults to `30s`.

##### `--timeout=DURATION`

Sets the timeout for a single run before it is considered as failed. Defaults to `30s`.

##### `--start-period=DURATION`

Sets an initialization time period in which failures are not counted towards the maximum number of retries. Defaults to `0s`.

##### `--retries=N`

Sets the number of retries before a container is considered `unhealthy`. Defaults to `3`.

## FROM DOCKERFILE (**beta**)

#### Synopsis

* `FROM DOCKERFILE [options...] <context-path>`

#### Description

The `FROM DOCKERFILE` command initializes a new build environment, inheriting from an existing Dockerfile. This allows the use of Dockerfiles in Earthly builds.

The `<context-path>` is the path where the Dockerfile build context exists. By default, it is assumed that a file named `Dockerfile` exists in that directory. The context path can be either a path on the host system, or an [artifact reference](../guides/target-ref.md#artifact-reference), pointing to a directory containing a `Dockerfile`.

{% hint style='info' %}
##### Note

This feature is currently in **Beta** and it has the following limitations:

* This feature only works with files named `Dockerfile`. The equivalent of the `-f` option available in `docker build` has not yet been implemented.
* `.dockerignore` is not used.
* The newer experimental features which exist in the Dockerfile syntax are not guaranteed to work correctly.
{% endhint %}

#### Options

##### `-f <dockerfile-path>`

Specify an alternative Dockerfile to use. The `<dockerfile-path>` can be either a path on the host system, relative to the current Earthfile, or an [artifact reference](../guides/target-ref.md#artifact-reference) pointing to a Dockerfile.

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the Dockerfile build arg identified by `<key>`. This option is similar to the `docker build --build-arg <key>=<value>` option.

##### `--target <target-name>`

In a multi-stage Dockerfile, sets the target to be used for the build. This option is similar to the `docker build --target <target-name>` option.

##### `--platform <platform>` (**beta**)

Specifies the platform to build on.

For more information see the [multi-platform guide](../guides/multi-platform.md).

## WITH DOCKER (**beta**)

#### Synopsis

```Dockerfile
WITH DOCKER [--pull <image-name>] [--load <image-name>=<target-ref>] [--compose <compose-file>]
            [--service <compose-service>] [--build-arg <key>=<value>] [--allow-privileged]
  <commands>
  ...
END
```

#### Description

The clause `WITH DOCKER` initializes a Docker daemon to be used in the context of a `RUN` command. The Docker daemon can be pre-loaded with a set of images using options such as `-pull` and `--load`. Once the execution of the `RUN` command has completed, the Docker daemon is stopped and all of its data is deleted, including any volumes and network configuration. Any other files that may have been created are kept, however.

The clause `WITH DOCKER` automatically implies the `RUN --privileged` flag.

The `WITH DOCKER` clause only supports the command [`RUN`](#run). Other commands (such as `COPY`) need to be run either before or after `WITH DOCKER ... END`. In addition, only one `RUN` command is permitted within `WITH DOCKER`. However, multiple shell commands may be stringed together using `;` or `&&`.

A typical example of a `WITH DOCKER` clause might be:

```Dockerfile
FROM earthly/dind:alpine
WORKDIR /test
COPY docker-compose.yml ./
WITH DOCKER \
        --compose docker-compose.yml \
        --load image-name:latest=+some-target \
        --pull some-image:latest
    RUN docker run ... && \
        docker run ... && \
        ...
END
```

For more examples, see the [Docker in Earthly guide](../guides/docker-in-earthly.md) and the [Integration testing guide](../guides/integration.md).

{% hint style='info' %}
##### Note
For performance reasons, it is recommended to use a Docker image that already contains `dockerd`. If `dockerd` is not found, Earthly will attempt to install it.

Earthly provides officially supported images such as `earthly/dind:alpine` and `earthly/dind:ubuntu` to be used together with `WITH DOCKER`.
{% endhint %}

#### Options

##### `--pull <image-name>`

Pulls the Docker image `<image-name>` from a remote registry and then loads it into the temporary Docker daemon created by `WITH DOCKER`.

This option may be repeated in order to provide multiple images to be pulled.

{% hint style='info' %}
##### Note
It is recommended that you avoid issuing `RUN docker pull ...` and use `WITH DOCKER --pull ...` instead. The classical `docker pull` command does not take into account Earthly caching and so it would redownload the image much more frequently than necessary.
{% endhint %}

##### `--load <image-name>=<target-ref>`

Builds the image referenced by `<target-ref>` and then loads it into the temporary Docker daemon created by `WITH DOCKER`. The image can be referenced as `<image-name>` within `WITH DOCKER`.

This option may be repeated in order to provide multiple images to be loaded.

##### `--compose <compose-file>`

Loads the compose definition defined in `<compose-file>`, adds all applicable images to the pull list and starts up all applicable compose services within.

This option may be repeated, thus having the same effect as repeating the `-f` flag in the `docker-compose` command.

##### `--service <compose-service>`

Specifies which compose service to pull and start up. If no services are specified and `--compose` is used, then all services are pulled and started up.

This option can only be used if `--compose` has been specified.

This option may be repeated in order to specify multiple services.

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the build arg identified by `<key>`, when building a `<target-ref>` (specified via `--load`). See also [BUILD](#build) for more details about the `--build-arg` option.

##### `--platform <platform>` (**beta**)

Specifies the platform for any referenced `--load` and `--pull` images.

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

## IF (**experimental**)

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/779) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack).
{% endhint %}

#### Synopsis

* ```
  IF [<condition-options>...] <condition>
    <if-block>
  END
  ```
* ```
  IF [<condition-options>...] <condition>
    <if-block>
  ELSE
    <else-block>
  END
  ```
* ```
  IF [<condition-options>...] <condition>
    <if-block>
  ELSE IF [<condition-options>...] <condition>
    <else-if-block>
  ...
  ELSE
    <else-block>
  END
  ```

#### Description

The `IF` clause can perform varying commands depending on the outcome of one or more conditions. The expression passed as part of `<condition>` is evaluated by running it in the build environment. If the exit code of the expression is zero, then the block of that condition is executed. Otherwise, the control continues to the next `ELSE IF` condition (if any), or if no condition returns a non-zero exit code, the control continues to executing the `<else-block>`, if one is provided.

A very common pattern is to use the POSIX shell `[ ... ]` conditions. For example the following marks port `8080` as exposed if the file `./foo` exists.

```Dockerfile
IF [ -f ./foo ]
  EXPOSE 8080
END
```

{% hint style='info' %}
##### Note
Performing a condition requires that a `FROM` (or a from-like command, such as `LOCALLY`) has been issued before the condition itself.

For example, the following is NOT a valid Earthfile.

```Dockerfile
# NOT A VALID EARTHFILE.
ARG base=alpine
IF [ "$base" = "alpine" ]
    FROM alpine:3.13
ELSE
    FROM ubuntu:20.10
END
```

The reason this is invalid is because the `IF` condition is actually running the `/usr/bin/[` executable to test if the condition is true or false, and therefore requires that a valid build environment has been initialized.

Here is how this might be fixed.

```Dockerfile
ARG base=alpine
FROM busybox
IF [ "$base" = "alpine" ]
    FROM alpine:3.13
ELSE
    FROM ubuntu:20.10
END
```

By initializing the build environment with `FROM busybox`, the `IF` condition can execute on top of the `busybox` image.
{% endhint %}

{% hint style='danger' %}
##### Important
Changes to the filesystem in any of the conditions are not preserved. If a file is created as part of a condition, then that file will not be present in the build environment for any subsequent commands.
{% endhint %}

#### Options

##### `--privileged`

Same as [`RUN --privileged`](#privileged).

##### `--ssh`

Same as [`RUN --ssh`](#ssh).

##### `--no-cache`

Same as [`RUN --no-cache`](#no-cache).

##### `--mount <mount-spec>`

Same as [`RUN --mount <mount-spec>`](#mount-less-than-mount-spec-greater-than).

##### `--secret <env-var>=<secret-ref>`

Same as [`RUN --secret <env-var>=<secret-ref>`](#secret-less-than-env-var-greater-than-less-than-secret-ref-greater-than).

## FOR (**experimental**)

Enable via `VERSION --for-in 0.5`.

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/779) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack).
{% endhint %}

#### Synopsis

* ```
  FOR [<options>...] <variable-name> IN <expression>
    <for-block>
  END
  ```

#### Description

The `FOR` clause can iterate over the items resulting from the expression `<expression>`. On each iteration, the value of `<variable-name>` is set to the current item in the iteration and the block of commands `<for-block>` is executed in the context of that variable set as a build arg.

The expression may be either a constant list of items (e.g. `foo bar buz`), or the output of a command (e.g. `$(echo foo bar buz)`), or a parameterized list of items (e.g. `foo $BARBUZ`). The result of the expression is then tokenized using the list of separators provided via the `--sep` option. If unspecified, the separator list defaults to `[tab]`, `[new line]` and `[space]` (`\t\n `).

{% hint style='danger' %}
##### Important
Changes to the filesystem in expressions are not preserved. If a file is created as part of a `FOR` expression, then that file will not be present in the build environment for any subsequent commands.
{% endhint %}

#### Examples

As an example, `FOR` may be used to iterate over a list of files for compilation

```Dockerfile
FOR file IN $(ls)
  RUN gcc "${file}" -o "${file}.o" -c
END
```

As another example, `FOR` may be used to iterate over a set of directories in a monorepo and invoking targets within them.

```Dockerfile
FOR dir IN $(ls -d */)
  BUILD "./$dir+build"
END
```

#### Options

##### `--sep <separator-list>`

The list of separators to use when tokenizing the output of the expression. If unspecified, the separator list defaults to `[tab]`, `[new line]` and `[space]` (`\t\n `).

##### `--privileged`

Same as [`RUN --privileged`](#privileged).

##### `--ssh`

Same as [`RUN --ssh`](#ssh).

##### `--no-cache`

Same as [`RUN --no-cache`](#no-cache).

##### `--mount <mount-spec>`

Same as [`RUN --mount <mount-spec>`](#mount-less-than-mount-spec-greater-than).

##### `--secret <env-var>=<secret-ref>`

Same as [`RUN --secret <env-var>=<secret-ref>`](#secret-less-than-env-var-greater-than-less-than-secret-ref-greater-than).

## LOCALLY (**experimental**)

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/580) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#locally` channel.
{% endhint %}

#### Synopsis

* `LOCALLY`

#### Description

The `LOCALLY` command can be used in place of a `FROM` command, which will cause earthly to execute all commands under the target directly
on the host system, rather than inside a container. Commands within a `LOCALLY` target will never be cached.
This feature should be used with caution as locally run commands have no guarantee they will behave the same on different systems.

Only `RUN` commands are supported under a `LOCALLY` defined target; furthermore only `RUN`'s `--push` flag is supported.

`RUN` commands have access to the environment variables which are exposed to the `earthly` command; however, the commands 
are executed within a working directory which is set to the location of the referenced Earthfile and not where the `earthly` command is run from.

For example, the following Earthfile will display the current user, hostname, and directory where the Earthfile is stored:

```Dockerfile
whoami:
    LOCALLY
    RUN echo "I am currently running under $USER on $(hostname) under $(pwd)"
```

{% hint style='info' %}
##### Note
In Earthly, outputting images and artifacts locally takes place only at the end of a successful build. In order to use such images or artifacts in `LOCALLY` targets, they need to be referenced correctly.

For images, use the `--load` option under `WITH DOCKER`:

```Dockerfile
my-image:
    FROM alpine 3.13
    ...
    SAVE IMAGE my-example-image

a-locally-example:
    LOCALLY
    WITH DOCKER --load=+my-image
        RUN docker run --rm my-example-image
    END
```

Do NOT use `BUILD` for using images in `LOCALLY` targets:

```Dockerfile
# INCORRECT - do not use!
my-image:
    FROM alpine 3.13
    ...
    SAVE IMAGE my-example-image

a-locally-example:
    LOCALLY
    BUILD +my-image
    # The image will not be available here because the local export of the
    # image only takes place at the end of an entire successful build.
    RUN docker run --rm my-example-image
```

For artifacts, use `COPY`, the same way you would in a regular target:

```Dockerfile
my-artifact:
    FROM alpine 3.13
    ...
    SAVE ARTIFACT ./my-example-artifact

a-locally-example:
    LOCALLY
    COPY +my-artifact/my-example-artifact ./
    RUN cat ./my-example-artifact
```

Do NOT use `SAVE ARTIFACT ... AS LOCAL` and `BUILD` for referencing artifacts in `LOCALLY` targets:

```Dockerfile
# INCORRECT - do not use!
my-artifact:
    FROM alpine 3.13
    ...
    SAVE ARTIFACT ./my-example-artifact AS LOCAL ./my-example-artifact

a-locally-example:
    LOCALLY
    BUILD +my-artifact
    # The artifact will not be available here because the local export of the
    # artifact only takes place at the end of an entire successful build.
    RUN cat ./my-example-artifact
```
{% endhint %}

## COMMAND (**experimental**)

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/581) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#udc` channel.
{% endhint %}

#### Synopsis

* `COMMAND`

#### Description

The command `COMMAND` marks the beginning of a user-defined command (UDC) definition. UDCs are templates (much like functions in regular programming languages), which can be used to define a series of steps to be executed in sequence. In order to reference and execute a UDC, you may use the command [`DO`](#do-experimental).

Unlike performing a `BUILD +target`, UDCs inherit the build context and the build environment from the caller.

UDCs create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +COMMAND --<build-arg-key>=<build-arg-value>`.

Global imports and global args are inherited from the `base` target of the same Earthfile where the command is defined in (this may be distinct from the `base` target of the caller).

For more information see the [User-defined commands guide](../guides/udc.md).

## DO (**experimental**)

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/581) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#udc` channel.
{% endhint %}

#### Synopsis

* `DO [--allow-privileged] <command-ref> [--<build-arg-key>=<build-arg-value>...]`

#### Description

The command `DO` expands and executes the series of commands contained within a user-defined command (UDC) [referenced by `<command-ref>`](../guides/target-ref.md#command-reference).

Unlike performing a `BUILD +target`, UDCs inherit the build context and the build environment from the caller.

UDCs create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +COMMAND --<build-arg-key>=<build-arg-value>`.

For more information see the [User-defined commands guide](../guides/udc.md).

#### Options

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

## IMPORT (**experimental**)

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/581) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#udc` channel.
{% endhint %}

#### Synopsis

* `IMPORT [--allow-privileged] <project-ref> [AS <alias>]`

#### Description

The command `IMPORT` aliases a project reference (`<project-ref>`) that can be used in subsequent [target, artifact or command references](../guides/target-ref.md).

If not provided, the `<alias>` is inferred automatically as the last element of the path provided in `<project-ref>`. For example, if `<project-ref>` is `github.com/foo/bar/buz:v1.2.3`, then the alias is inferred as `buz`.

The `<project-ref>` can be a reference to any directory other than `.`. If the reference ends in `..`, then mentioning `AS <alias>` is mandatory.

If an `IMPORT` is defined in the `base` target of the Earthfile, then it becomes a global `IMPORT` and it is made available to every other target or command in that file, regardless of their base images used.

For more information see the [target, artifact and command references guide](../guides/target-ref.md).

#### Options

##### `--allow-privileged`

Similar to [`FROM --allow-privileged`](#allow-privileged), extend the ability to request privileged capabilities to all invocations of the imported alias.

## SHELL (not supported)

The classical [`SHELL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#add) is not yet supported. Use the *exec form* of `RUN`, `ENTRYPOINT` and `CMD` instead and prepend a different shell.

## ADD (not supported)

The classical [`ADD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#add) is not yet supported. Use [COPY](#copy) instead.

## ONBUILD (not supported)

The classical [`ONBUILD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#onbuild) is not supported.

## STOPSIGNAL (not supported)

The classical [`STOPSIGNAL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#stopsignal) is not yet supported.

## DOCKER PULL (**deprecated**)

#### Synopsis

* `DOCKER PULL <image-name>`

#### Description

{% hint style='danger' %}
`DOCKER PULL` is now deprecated and will not be supported in future versions of Earthly. Please use `WITH DOCKER --pull <image-name>` instead.
{% endhint %}

## DOCKER LOAD (**deprecated**)

#### Synopsis

* `DOCKER LOAD [--build-arg <name>=<override-value>] <target-ref> <image-name>`

#### Description

{% hint style='danger' %}
`DOCKER LOAD` is now deprecated and will not be supported in future versions of Earthly. Please use `WITH DOCKER --load <image-name>=<target-ref>` instead.
{% endhint %}
