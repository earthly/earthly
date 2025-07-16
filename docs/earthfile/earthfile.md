# Earthfile reference

<!--

Note to person editing!!

The general order of the commands is as follows:

- Core classical Dockerfile commands (order is the same as in the Dockerfile official docs)
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
* `FROM [--platform <platform>] [--allow-privileged] <target-ref> [--<build-arg-key>=<build-arg-value>...]`

#### Description

The `FROM` command initializes a new build environment and sets the base image for subsequent instructions. It works similarly to the classical [Dockerfile `FROM` instruction](https://docs.docker.com/engine/reference/builder/#from), but it has the added ability to use another [target](https://docs.earthly.dev/docs/guides/target-ref#target-reference)'s image as the base image.

Examples:

* Classical reference: `FROM alpine:latest`
* Local reference: `FROM +another-target`
* Relative reference: `FROM ./subdirectory+some-target` or `FROM ../otherdirectory+some-target`
* Absolute reference: `FROM /absolute/path+some-target`
* Remote reference from a public or [private](https://docs.earthly.dev/docs/guides/auth) git repository: `FROM github.com/example/project+remote-target`

The `FROM` command does not mark any saved images or artifacts of the referenced target for output, nor does it mark any push commands of the referenced target for pushing. For that, please use [`BUILD`](#build).

{% hint style='info' %}
##### Note

The `FROM ... AS ...` form available in the classical Dockerfile syntax is not supported in Earthfiles. Instead, define a new Earthly target. For example, the following Dockerfile

```Dockerfile
# Dockerfile

FROM alpine:3.18 AS build
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
    FROM alpine:3.18
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

##### `--<build-arg-key>=<build-arg-value>`

Sets a value override of `<build-arg-value>` for the build arg identified by `<build-arg-key>`. See also [BUILD](#build) for more details about build args.

##### `--platform <platform>`

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

##### `--pass-args`

Earthly automatically passes all current arguments to referenced targets in the _same_ Earthfile.
However, when the `--pass-args` flag is set, Earthly will also propagate all arguments to an externally referenced target.

##### `--build-arg <key>=<value>` (**deprecated**)

This option is deprecated. Use `--<build-arg-key>=<build-arg-value>` instead.

## RUN

#### Synopsis

* `RUN [options...] [--] <command>` (shell form)
* `RUN [[options...], "<executable>", "<arg1>", "<arg2>", ...]` (exec form)

#### Description

The `RUN` command executes commands in the build environment of the current target, in a new layer. It works similarly to the [Dockerfile `RUN` command](https://docs.docker.com/engine/reference/builder/#run), with some added options.

The command allows for two possible forms. The *exec form* runs the command executable without the use of a shell. The *shell form* uses the default shell (`/bin/sh -c`) to interpret the command and execute it. In either form, you can use a `\` to continue a single `RUN` instruction onto the next line.

When the `--entrypoint` flag is used, the current image entrypoint is used to prepend the current command.

To avoid any ambiguity regarding whether an argument is a `RUN` flag option or part of the command, the delimiter `--` may be used to signal the parser that no more `RUN` flag options will follow.

#### Options

##### `--push`

Marks the command as a "push command". Push commands are never cached, thus they are executed on every applicable invocation of the build.

Push commands are not run by default. Add the `--push` flag to the `earthly` invocation to enable pushing. For example

```bash
earthly --push +deploy
```

Push commands were introduced to allow the user to define commands that have an effect external to the build. Good candidates for push commands are uploads of artifacts to artifactories, commands that make a change to an external environment, like a production or staging environment.

##### `--no-cache`

Force the command to run every time; ignoring the layer cache. Any commands following the invocation of `RUN --no-cache`, will also ignore the cache. If `--no-cache` is used as an option on the `RUN` statement within a `WITH DOCKER` statement, all commands after the `WITH DOCKER` will also ignore the cache.

{% hint style='danger' %}
##### Auto-skip
Note that `RUN --no-cache` commands may still be skipped by auto-skip. For more information see the [Caching in Earthfiles guide](../caching/caching-in-earthfiles.md#auto-skip).
{% endhint %}

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

##### `--secret <env-var>=<secret-id> | <secret-id>`

Makes available a secret, in the form of an env var (its name is defined by `<env-var>`), to the command being executed.
If you only specify `<secret-id>`, the name of the env var will be `<secret-id>` and its value the value of `<secret-id>`.

Here is an example that showcases both syntaxes:

```Dockerfile
release:
    RUN --push --secret GITHUB_TOKEN=GH_TOKEN github-release upload
release-short:
    RUN --push --secret GITHUB_TOKEN github-release upload
```

```bash
earthly --secret GH_TOKEN="the-actual-secret-token-value" +release
earthly --secret GITHUB_TOKEN="the-actual-secret-token-value" +release-short
```

An empty string is also allowed for `<secret-id>`, allowing for optional secrets, should it need to be disabled.

```Dockerfile
release:
    ARG SECRET_ID=GH_TOKEN
    RUN --push --secret GITHUB_TOKEN=$SECRET_ID github-release upload
release-short:
    ARG SECRET_ID=GITHUB_TOKEN
    RUN --push --secret $SECRET_ID github-release upload
```

```bash
earthly +release --SECRET_ID=""
earthly +release-short --SECRET_ID=""
```

It is also possible to mount a secret as a file with `RUN --mount type=secret,id=secret-id,target=/path/of/secret,chmod=0400`. See `--mount` below.

For more information on how to use secrets see the [Secrets guide](../guides/secrets.md). See also the [Cloud secrets guide](../cloud/cloud-secrets.md).

##### `--network=none`

Isolate the networking stack (and internet access) from the command.

##### `--ssh`

Allows a command to access the ssh authentication client running on the host via the socket which is referenced by the environment variable `SSH_AUTH_SOCK`.

Here is an example:

```Dockerfile
RUN mkdir -p ~/.ssh && \
    echo 'github.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCj7ndNxQowgcQnjshcLrqPEiiphnt+VTTvDP6mHBL9j1aNUkY4Ue1gvwnGLVlOhGeYrnZaMgRK6+PKCUXaDbC7qtbW8gIkhL7aGCsOr/C56SJMy/BCZfxd1nWzAOxSDPgVsmerOBYfNqltV9/hWCqBywINIR+5dIg6JTJ72pcEpEjcYgXkE2YEFXV1JHnsKgbLWNlhScqb2UmyRkQyytRLtL+38TGxkxCflmO+5Z8CSSNY7GidjMIZ7Q4zMjA2n1nGrlTDkzwDCsw+wqFPGQA179cnfGWOWRVruj16z6XyvxvjJwbz0wQZ75XK5tKSb7FNyeIEs4TT4jk+S4dhPeAUC5y+bDYirYgM4GC7uEnztnZyaVWQ7B381AK4Qdrwt51ZqExKbQpTUNn+EjqoTwvqNj4kqx5QUCI0ThS/YkOxJCXmPUWZbhjpCg56i+2aB6CmK2JGhn57K5mj0MNdBXA4/WnwH6XoPWJzK5Nyu2zB3nAZp+S5hpQs+p1vN1/wsjk=' >> ~/.ssh/known_hosts && \
    echo 'gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9' >> ~/.ssh/known_hosts
RUN --ssh git config --global url."git@github.com:".insteadOf "https://github.com/" && \
    go mod download
```

{% hint style='warning' %}
Note that `RUN --ssh` option is only used for creating a tunnel to the host's ssh-agent's socket (set via `$SSH_AUTH_SOCK`); it is **not** related to the git section of the earthly [configuration file](../earthly-config/earthly-config.md).
{% endhint %}

##### `--mount <mount-spec>`

Mounts a file or directory in the context of the build environment.

The `<mount-spec>` is defined as a series of comma-separated list of key-values. The following keys are allowed:

| Key             | Description                                                                                                                                                                                                                  | Example                                 |
|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------|
| `type`          | The type of the mount. Currently only `cache`, `tmpfs`, and `secret` are allowed.                                                                                                                                            | `type=cache`                            |
| `target`        | The target path for the mount.                                                                                                                                                                                               | `target=/var/lib/data`                  |
| `mode`, `chmod` | The permission of the mounted file, in octal format (the same format the chmod unix command line expects).                                                                                                                   | `chmod=0400`                            |
| `id`            | The cache ID for a global cache mount to be used across other targets or Earthfiles, when `type=cache`. The secret ID for the contents of the `target` file, when `type=secret`. | `id=my-shared-cache`, `id=my-password`  |
| `sharing`       | The sharing mode (`locked`, `shared`, `private`) for the cache mount, only applicable for `type=cache`.                                                                                                                      | `sharing=shared`                        |

For cache mounts, the sharing mode can be one of the following:

* `locked` (default) - the cache mount is locked for the duration of the execution, other concurrent builds will wait for the lock to be released.
* `shared` - the cache mount is shared between all concurrent builds.
* `private` - if another concurrent build attempts to use the cache, a new (empty) cache will be created for the concurrent build.

###### Examples

Persisting cache for a single `RUN` command, even when its dependencies change:

```Dockerfile
ENV GOCACHE=/go-cache
RUN --mount=type=cache,target=/go-cache go build main.go
```

{% hint style='warning' %}
Note that mounts cannot be shared between targets, nor can they be shared within the same target, if the build-args differ between invocations.
{% endhint %}

Mounting a secret as a file:

```Dockerfile
RUN --mount=type=secret,id=netrc,target=/root/.netrc curl https://example.earthly.dev/restricted/example-file-that-requires-auth > data
```

The contents of the secret `/root/.netrc` file can then be specified from the command line as:

```bash
earthly --secret netrc="machine example.earthly.dev login myusername password mypassword" +base
```

or by passing the contents of an existing file from the host filesystem:

```bash
earthly --secret-file netrc="$HOME/.netrc" +base
```


##### `--interactive` / `--interactive-keep`

Opens an interactive prompt during the target build. An interactive prompt must:

1. Be the last issued command in the target, with the exception of `SAVE IMAGE` commands. This also means that you cannot `FROM` a target containing a `RUN --interactive`.
2. Be the only `--interactive` target within the run.
3. Not be within a `LOCALLY`-designated target.

###### Examples:

Start an interactive python REPL:
```Dockerfile
python:
    FROM alpine:3.18
    RUN apk add python
    RUN --interactive python
```

Start `bash` to tweak an image by hand. Changes made will be included:
```Dockerfile
build:
    FROM alpine:3.18
    RUN apk add bash
    RUN --interactive-keep bash
```

##### `--aws` (experimental)

{% hint style='info' %}
##### Note
The `--aws` flag has experimental status. To use this feature, it must be enabled via `VERSION --run-with-aws 0.8`.
{% endhint %}

Makes AWS credentials available to the executed command via the host's environment variables or ~/.aws directory.

##### `--oidc <oidc-spec>` (experimental)

{% hint style='info' %}
##### Note
The `--oidc` flag has experimental status and can only be used conjointly with the `--aws` flag. To use this feature, it must be enabled via `VERSION --run-with-aws --run-with-aws-oidc 0.8`.
{% endhint %}

Makes AWS credentials available to the executed command via AWS OIDC provider.

The `<oidc-spec>` is defined as a series of comma-separated list of key-values. The following keys are allowed:

| Key                | Description                                                                                                                                                                                                                                            | Example                                             |
|--------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------|
| `session-name`     | The session name to identify in AWS's logs. If any `RUN ... --oidc` commands use the same `session-name`, they will share the same temporary token                                                                                                     | `session-name=my-session`                           |
| `role-arn`         | The AWS arn of the role for which to get credentials.                                                                                                                                                                                                  | `role-arn=arn:aws:iam::123456789012:role/some-role` |
| `region`           | The AWS region to connect to in order to get the credentials. This will also be the region used by the executed AWS command (though the region may be overridden in the command). If the region is not specified, the global AWS endpoint will be used | `region=us-east-1`                                  |
| `session-duration` | The time the credentials will be valid for before they expire. Default (AWS minimum): 15 minutes.                                                                                                                                                      | `session-duration=20m`                              |

Click [here](../cloud/oidc.md#openid-connect-oidc-authentication) for more information on how to configure OIDC in AWS for Earthly.

##### `--raw-output` (experimental)

{% hint style='info' %}
##### Note
The `--raw-output` flag has experimental status. To use this feature, it must be enabled via `VERSION --raw-output `.
{% endhint %}

Outputs line without target name.

###### Examples:

Given this target:
```Dockerfile
raw:
    RUN --raw-output echo "::group::"
    RUN echo "should have prefix"
    RUN --raw-output echo "::endgroup::"
```

The following is output:
```bash
 ./+gha | --> RUN --raw-output echo "::group::"
::group::
 ./+gha | --> RUN echo "should have prefix"
 ./+gha | should have prefix
 ./+gha | --> RUN --raw-output echo "::endgroup::"
::endgroup::
```

## COPY

#### Synopsis

* `COPY [options...] <src>... <dest>` (classical form)
* `COPY [options...] <src-artifact>... <dest>` (artifact form)
* `COPY [options...] (<src-artifact> --<build-arg-key>=<build-arg-value>...) <dest>` (artifact form with build args)

#### Description

The command `COPY` allows copying of files and directories between different contexts.

The command may take a couple of possible forms. In the *classical form*, `COPY` copies files and directories from the build context into the build environment - in this form, it works similarly to the [Dockerfile `COPY` command](https://docs.docker.com/engine/reference/builder/#copy). In the *artifact form*, `COPY` copies files or directories (also known as "artifacts" in this context) from the artifact environment of other build targets into the build environment of the current target. Either form allows the use of wildcards for the sources.

The parameter `<src-artifact>` is an [artifact reference](../guides/importing.md#artifact-reference) and is generally of the form `<target-ref>/<artifact-path>`, where `<target-ref>` is the reference to the target which needs to be built in order to yield the artifact and `<artifact-path>` is the path within the artifact environment of the target, where the file or directory is located. The `<artifact-path>` may also be a wildcard.

{% hint style='info' %}
##### Globbing
A target reference in a <src-artifact> may also include a glob expression.
This is useful in order to invoke multiple targets that may exist in different Earthfiles in the filesystem, in a single `COPY` command.
For example, consider the following filesystem:
```bash
services
├── Earthfile
├── service1
│    └── Earthfile
├── service2
│   ├── Earthfile
├── service3
│   ├── Earthfile
```

where a `+mocks` target is defined in services1/Earthfile, services2/Earthfile and services3/Earthfile.
The command `COPY ./services/*+mocks .` is equivalent to:
```Earthfile
    COPY ./services/service1+mocks .
    COPY ./services/service2+mocks .
    COPY ./services/service3+mocks .
```

A glob match occurs when an Earthfile in the glob expression path exists, and the named target is defined in the Earthfile.
At least one match must be found for the command to succeed.

This feature has experimental status. To use it, it must be enabled via `VERSION --wildcard-copy 0.8`.
(This is not to be confused with the usage of wildcards in the artifact name, which is fully supported, e.g. `COPY ./services/service1+mocks/* .`)

{% endhint %}

The `COPY` command does not mark any saved images or artifacts of the referenced target for output, nor does it mark any push commands of the referenced target for pushing. For that, please use [`BUILD`](#build).

Multiple `COPY` commands issued one after the other will build the referenced targets in parallel, if the targets don't depend on each other. The resulting artifacts will then be copied sequentially in the order in which the `COPY` commands were issued.

The classical form of the `COPY` command differs from Dockerfiles in three cases:

* URL sources are not yet supported.
* Absolute paths are not supported - sources in the current directory cannot be referenced with a leading `/`
* The Earthly `COPY` is a classical `COPY --link`. It uses layer merging for the copy operations.

{% hint style='info' %}
##### Note
To prevent Earthly from copying unwanted files, you may specify file patterns to be excluded from the build context using an [`.earthlyignore`](./earthlyignore.md) file. This file has the same syntax as a [`.dockerignore` file](https://docs.docker.com/engine/reference/builder/#dockerignore-file).
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

If the directories were copied without the use of `--dir`, then their contents would be merged into the destination.

##### `--<build-arg-key>=<build-arg-value>`

Sets a value override of `<build-arg-value>` for the build arg identified by `<build-arg-key>`, when building the target containing the mentioned artifact. See also [BUILD](#build) for more details about the build arg options.

Note that build args and the artifact references they apply to need to be surrounded by parenthesis:

```Dockerfile
COPY (+target1/artifact --arg1=foo --arg2=bar) ./dest/path
```

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

##### `--keep-own`

Instructs Earthly to keep file ownership information. This applies only to the *artifact form* and has no effect otherwise.

##### `--chmod <octal-format>`

Instructs Earthly to change the file permissions of the copied files. The `<chmod>` needs to be in octal format, e.g. `--chmod 0755` or `--chmod 755`.

{% hint style='info' %}
Note that you must include the flag in the corresponding `SAVE ARTIFACT --keep-own ...` command, if using *artifact form*.
{% endhint %}

##### `--if-exists`

Only copy source if it exists; if it does not exist, earthly will simply ignore the COPY command and won't treat any missing sources as failures.

##### `--symlink-no-follow`

Allows copying a symbolic link from another target; it has no effect when copying files from the host.
The option must be used in both the `COPY` and `SAVE ARTIFACT` commands; for example:

```Dockerfile
producer:
    RUN ln -s nonexistentfile symlink
    SAVE ARTIFACT --symlink-no-follow symlink

consumer:
    COPY --symlink-no-follow +producer/symlink
```

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

##### `--platform <platform>`

In *artifact form*, it specifies the platform to build the artifact on.

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

##### `--pass-args`

Same as [`FROM --pass-args`](#pass-args).

##### `--build-arg <key>=<value>` (**deprecated**)

The option `--build-arg` is deprecated. Use `--<build-arg-key>=<build-arg-value>` instead.

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

One can also copy from other Earthfile targets:

```
FROM alpine:3.18
dummy-target:
    RUN echo aGVsbG8= > encoded-data
    SAVE ARTIFACT encoded-data
example:
    COPY +dummy-target/encoded-data .
    RUN cat encoded-data | base64 -d
```

Parentheses are required when passing build-args:

```
FROM alpine:3.18
RUN apk add coreutils # required for base32 binary
dummy-target:
    ARG encoder="base64"
    RUN echo hello | $encoder > encoded-data
    SAVE ARTIFACT encoded-data
example:
    COPY ( +dummy-target/encoded-data --encoder=base32 ) .
    RUN cat encoded-data | base32 -d
```

For detailed examples demonstrating how other scenarios may function, please see our [test suite](https://github.com/earthly/earthly/blob/main/tests/copy.earth).

## ARG

#### Synopsis

* `ARG [--required] <name>[=<default-value>]` (constant form)
* `ARG [--required] <name>=$(<default-value-expr>)` (dynamic form)

#### Description

The command `ARG` declares a build argument (or arg) with the name `<name>` and with an optional default value `<default-value>`. If no default value is provided, then empty string is used as the default value.

This command works similarly to the [Dockerfile `ARG` command](https://docs.docker.com/engine/reference/builder/#arg), with a few differences regarding the scope and the predefined args (called builtin args in Earthly). The arg's scope is always limited to the recipe of the current target or command and only from the point it is declared onward. For more information regarding builtin args, see the [builtin args page](./builtin-args.md).

In its *constant form*, the arg takes a default value defined as a constant string. If the `<default-value>` is not provided, then the default value is an empty string. In its *dynamic form*, the arg takes a default value defined as an expression. The expression is evaluated at run time and its result is used as the default value. The expression is interpreted via the default shell (`/bin/sh -c`) within the build environment.

The value of an arg can be overridden either from the `earthly` command

```bash
earthly <target-ref> --<name>=<override-value>
```

or from a command from another target, when implicitly or explicitly invoking the target containing the `ARG`

```Dockerfile
BUILD <target-ref> --<name>=<override-value>
COPY (<target-ref>/<artifact-path> --<name>=<override-value>) <dest-path>
FROM <target-ref> --<name>=<override-value>
```

for example

```Dockerfile
BUILD +binary --NAME=john
COPY (+binary/bin --NAME=john) ./
FROM +docker-image --NAME=john
```

For more information on how to use build args see the [build arguments and variables guide](../guides/build-args.md). A number of builtin args are available and are pre-filled by Earthly. For more information see [builtin args](./builtin-args.md).

#### Options

##### `--required`

A required `ARG` must be provided at build time and can never have a default value. Required args can help eliminate cases where the user has unexpectedly set an `ARG` to `""`.

```
target-required:
    # user must supply build arg for target
    ARG --required NAME

build-linux:
    # or explicitly supply in build command
    BUILD +target-required --NAME=john
```

#### `--global`

A global `ARG` is an arg that is made available to all targets in the Earthfile. This is useful for setting a default value for an arg that is used in many targets.

Global args may only be declared in base targets.

{% hint style='danger' %}
##### Important
Avoid using `ARG --global` for args that change frequently (e.g. git sha, branch name, PR number, etc). Any change to the value of this arg would typically cause all targets in the Earthfile to re-execute with no cache.

It's always best to declare args as deep and late as possible within the specific target where they are needed, to get the most performance, even if this may require more verbose passing of args from one target to another. See also [`BUILD --pass-args`](#build).
{% endhint %}

## SAVE ARTIFACT

#### Synopsis

* `SAVE ARTIFACT [--keep-ts] [--keep-own] [--if-exists] [--force] <src> [<artifact-dest-path>] [AS LOCAL <local-path>]`

#### Description

The command `SAVE ARTIFACT` copies a file, a directory, or a series of files and directories represented by a wildcard, from the build environment into the target's artifact environment.

If `AS LOCAL ...` is also specified, it additionally marks the artifact to be copied to the host at the location specified by `<local-path>`, once the build is deemed as successful. Note that local artifacts are only produced by targets that are run directly with `earthly`, or when invoked using [`BUILD`](#build).

If `<artifact-dest-path>` is not specified, it is inferred as `/`.

Files within the artifact environment are also known as "artifacts". Once a file has been copied into the artifact environment, it can be referenced in other places of the build (for example in a `COPY` command), using an [artifact reference](../guides/importing.md#artifact-reference).

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

{% hint style='danger' %}
##### Important

As of [`VERSION 0.6`](#version), local artifacts are only saved [if they are connected to the initial target through a chain of `BUILD` commands](#what-is-being-output-and-pushed).

{% endhint %}

#### Options

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

##### `--keep-own`

Instructs Earthly to keep file ownership information.

##### `--if-exists`

Only save artifacts if they exists; if not, earthly will simply ignore the SAVE ARTIFACT command and won't treat any missing sources as failures.

##### `--symlink-no-follow`

Save the symbolic link rather than the contents of the symbolically linked file. Note that the same flag must also be used in the corresponding `COPY` command. For example:

```Dockerfile
producer:
    RUN ln -s nonexistentfile symlink
    SAVE ARTIFACT --symlink-no-follow symlink

consumer:
    COPY --symlink-no-follow +producer/symlink
```


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

For detailed examples demonstrating how other scenarios may function, please see our [test suite](https://github.com/earthly/earthly/blob/main/tests/file-copying.earth).

## SAVE IMAGE

#### Synopsis

* `SAVE IMAGE [--push] <image-name>...`

#### Description

The command `SAVE IMAGE` marks the current build environment as the image of the target and assigns one or more output image names.

{% hint style='info' %}
##### Assigning multiple image names

The `SAVE IMAGE` command allows you to assign more than one image name:

```Dockerfile
SAVE IMAGE my-image:latest my-image:1.0.0 my-example-registry.com/another-image:latest
```

Or

```Dockerfile
SAVE IMAGE my-image:latest
SAVE IMAGE my-image:1.0.0
SAVE IMAGE my-example-registry.com/another-image:latest
```
{% endhint %}

{% hint style='danger' %}
##### Important

As of [`VERSION 0.6`](#version), images are only saved [if they are connected to the initial target through a chain of `BUILD` commands](#what-is-being-output-and-pushed).

{% endhint %}

#### Options

##### `--push`

The `--push` options marks the image to be pushed to an external registry after it has been loaded within the docker daemon available on the host.

If inline caching is enabled, the `--push` option also instructs Earthly to use the specified image names as cache sources.

The actual push is not executed by default. Add the `--push` flag to the earthly invocation to enable pushing. For example

```bash
earthly --push +docker-image
```

##### `--no-manifest-list`

Instructs Earthly to not create a manifest list for the image. This may be useful on platforms that do not support multi-platform images (for example, AWS Lambda), and the image produced needs to be of a different platform than the default one.

## BUILD

#### Synopsis

* `BUILD [options...] <target-ref> [--<build-arg-name>=<build-arg-value>...]`

#### Description

The command `BUILD` instructs Earthly to additionally invoke the build of the target referenced by `<target-ref>`, where `<target-ref>` follows the rules defined by [target referencing](../guides/importing.md#target-reference). The invocation will mark any images, or artifacts saved by the referenced target for local output (assuming local output is enabled), and any push commands issued by the referenced target for pushing (assuming pushing is enabled).

Multiple `BUILD` commands issued one after the other will be executed in parallel if the referenced targets don't depend on each other.

{% hint style='info' %}
##### What is being output and pushed

In Earthly v0.6+, what is being output and pushed is determined either by the main target being invoked on the command-line directly, or by targets directly connected to it via a chain of `BUILD` calls. Other ways to reference a target, such as `FROM`, `COPY`, `WITH DOCKER --load` etc, do not contribute to the final set of outputs or pushes.

If you are referencing a target via some other command, such as `COPY` and you would like for the outputs or pushes to be included, you can issue an equivalent `BUILD` command in addition to the `COPY`. For example

{% hint style='info' %}
##### Globbing
A <target-ref> may also include a glob expression.
This is useful in order to invoke multiple targets that may exist in different Earthfiles in the filesystem, in a single `BUILD` command.
For example, consider the following filesystem:
```bash
services
├── Earthfile
├── service1
│    └── Earthfile
├── service2
│   ├── Earthfile
├── service3
│   ├── Earthfile
```

where a `+compile` target is defined in services1/Earthfile, services2/Earthfile and services3/Earthfile.
The command `BUILD ./services/*+compile .` is equivalent to:
```Earthfile
    BUILD ./services/service1+compile
    BUILD ./services/service2+compile
    BUILD ./services/service3+compile
```

A glob match occurs when an Earthfile in the glob expression path exists, and the named target is defined in the Earthfile.
At least one match must be found for the command to succeed.

This feature has experimental status. To use it, it must be enabled via `VERSION --wildcard-builds 0.8`.

{% endhint %}

```Dockerfile
my-target:
    COPY --platform=linux/amd64 (+some-target/some-file.txt --FOO=bar) ./
```

Should be amended with the following additional `BUILD` call:

```Dockerfile
my-target:
    BUILD --platform=linux/amd64 +some-target --FOO=bar
    COPY --platform=linux/amd64 (+some-target/some-file.txt --FOO=bar) ./
```

This, however, assumes that the target `+my-target` is itself connected via a `BUILD` chain to the main target being built. If that is not the case, additional `BUILD` commands should be issued higher up the hierarchy.
{% endhint %}

#### Options

##### `--<build-arg-key>=<build-arg-value>`

Sets a value override of `<build-arg-value>` for the build arg identified by `<build-arg-key>`.

The override value of a build arg may be a constant string

```
--SOME_ARG="a constant value"
```

or an expression involving other build args

```
--SOME_ARG="a value based on other args, like $ANOTHER_ARG and $YET_ANOTHER_ARG"
```

or a dynamic expression, based on the output of a command executed in the context of the build environment.

```
--SOME_ARG=$(find /app -type f -name '*.php')
```

Dynamic expressions are delimited by `$(...)`.

##### `--platform <platform>`

Specifies the platform to build on.

This flag may be repeated in order to instruct the system to perform the build for multiple platforms. For example

```Dockerfile
build-all-platforms:
    BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build
```

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--auto-skip` (*beta*)

Instructs Earthly to skip the build of the target if the target's dependencies have not changed from a previous successful build. For more information on how to use this feature, see the [auto-skip section of the caching in Earthfiles guide](../caching/caching-in-earthfiles.md#auto-skip).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

##### `--pass-args`

Same as [`FROM --pass-args`](#pass-args).

##### `--build-arg <build-arg-key>=<build-arg-value>` (**deprecated**)

This option is deprecated. Please use `--<build-arg-key>=<build-arg-value>` instead.

## LET

#### Synopsis

* `LET <name>=<value>`

#### Description

The command `LET` declares a variable with the name `<name>` and with a value `<value>`. This command works similarly to `ARG` except that it cannot be overridden.

`LET` variables are allowed to shadow `ARG` build arguments, which allows you to promote an `ARG` to a local variable so that it may be used with `SET`.

##### Example

```
VERSION 0.8

# mode defines the build mode. Valid values are 'dev' and 'prod'.
ARG --global mode = dev

foo:
    LET buildArgs = --mode development
    IF [ "$mode" = "prod" ]
        SET buildArgs = --mode production --optimize
    END
```

## SET

#### Synopsis

* `SET <name>=<value>`

#### Description

The command `SET` may be used to change the value of a previously declared variable, so long as the variable was declared with `LET`.

`ARG` variables may *not* be changed by `SET`, since `ARG` is intended to accept overrides from the CLI. If you want to change the value of an `ARG` variable, redeclare it with `LET someVar = "$someVar"` first.

See [the `LET` docs for more info](#let).

## VERSION

#### Synopsis

* `VERSION [options...] <version-number>`

#### Description

The command `VERSION` identifies which set of features to enable in Earthly while handling the corresponding Earthfile. Different `VERSION`s can be mixed together across different Earthfiles in the same project. Earthly handles a mix of versions gracefully, enabling or disabling features accordingly. This allows for gradual updates of `VERSION`s across large projects, without sacrificing build consistency.

The `VERSION` command is mandatory starting with Earthly 0.7. The `VERSION` command must be the first command in the Earthfile.

#### Options

Individual features may be enabled by setting the corresponding feature flag.
New features start off as experimental, which is why they are disabled by default.
Once a feature reaches maturity, it will be enabled by default under a new version number.

{% hint style='danger' %}
##### Important
Avoid using feature flags for critical workflows. You should only use feature flags for testing new experimental features. By using feature flags you are opting out of forwards/backwards compatibility guarantees. This means that running the same script in a different environment, with a different version of Earthly may result in a different behavior (i.e. it'll work on your machine, but may break the build for your colleagues or for the CI).
{% endhint %}

All features are described in [the version-specific features reference](./features.md).

## PROJECT

#### Synopsis

* `PROJECT <org-name>/<project-name>`

#### Description

The command `PROJECT` marks the current Earthfile as being part of the project belonging to the [Earthly organization](https://docs.earthly.dev/earthly-cloud/overview) `<org-name>` and the project `<project-name>`. The project is used by Earthly to retrieve [cloud-based secrets](../cloud/cloud-secrets.md) and build logs belonging to the project.

The `PROJECT` command can only be used in the `base` recipe and it applies to the entire Earthfile. The `PROJECT` command can never contain any `ARG`s that need expanding.

## GIT CLONE

#### Synopsis

* `GIT CLONE [--branch <git-ref>] [--keep-ts] <git-url> <dest-path>`

#### Description

The command `GIT CLONE` clones a git repository from `<git-url>`, optionally referenced by `<git-ref>`, into the build environment, within the `<dest-path>`.

In contrast to an operation like `RUN git clone <git-url> <dest-path>`, the command `GIT CLONE` is cache-aware and correctly distinguishes between different git commit IDs when deciding to reuse a previous cache or not. In addition, `GIT CLONE` can also use [Git authentication configuration](../guides/auth.md) passed on to `earthly`, whereas `RUN git clone` would require additional secrets passing, if the repository is not publicly accessible.

Note that the repository is cloned via a shallow-clone opperation (i.e. a single-depth clone).

{% hint style='info' %}

If you need to perform a full-depth clone of a repository, you can use the following pattern:

```Dockerfile
GIT CLONE <git-url> <dest-path>
WORKDIR <dest-path>
ARG git_hash=$(git rev-parse HEAD)
RUN git remote set-url origin <git-url> # only required if using authentication
RUN git fetch --unshallow
```
{% endhint %}

{% hint style='warning' %}
As of Earthly v0.7.21, git credentials are no longer stored in the `.git/config` file; this includes the username.
This means any ssh-based or https-based fetches or pushes will no longer work unless you restore the configured url,
which can be done with:
```Dockerfile
RUN git remote set-url origin <git-url>
```
{% endhint %}

See the "GIT CLONE vs RUN git clone" section under the [best practices guide](../guides/best-practices.md#git-clone-vs-run-git-clone) for more details.

#### Options

##### `--branch <git-ref>`

Points the `HEAD` to the git reference specified by `<git-ref>`. If this option is not specified, then the remote `HEAD` is used instead.

##### `--keep-ts`

Instructs Earthly to not overwrite the file creation timestamps with a constant.

## FROM DOCKERFILE

#### Synopsis

* `FROM DOCKERFILE [options...] <context-path>`

#### Description

The `FROM DOCKERFILE` command initializes a new build environment, inheriting from an existing Dockerfile. This allows the use of Dockerfiles in Earthly builds.

The `<context-path>` is the path where the Dockerfile build context exists. By default, it is assumed that a file named `Dockerfile` exists in that directory. The context path can be either a path on the host system, or an [artifact reference](../guides/importing.md#artifact-reference), pointing to a directory containing a `Dockerfile`.
Additionally, when using a `<context-path>` from the host system, a `.dockerignore` in the directory root will be used to exclude files (unless `.earthlyignore` or `.earthignore` are present).

#### Options

##### `-f <dockerfile-path>`

Specify an alternative Dockerfile to use. The `<dockerfile-path>` can be either a path on the host system, relative to the current Earthfile, or an [artifact reference](../guides/importing.md#artifact-reference) pointing to a Dockerfile.

{% hint style='info' %}
It is possible to split the `Dockerfile` and the build context across two separate [artifact references](../guides/importing.md#artifact-reference):

```Dockerfile
FROM alpine

mybuildcontext:
    WORKDIR /mydata
    RUN echo mydata > myfile
    SAVE ARTIFACT /mydata

mydockerfile:
    RUN echo "
FROM busybox
COPY myfile .
RUN cat myfile" > Dockerfile
    SAVE ARTIFACT Dockerfile

docker:
    FROM DOCKERFILE -f +mydockerfile/Dockerfile +mybuildcontext/mydata/*
    SAVE IMAGE testimg:latest
```

Note that `+mybuildcontext/mydata` on its own would copy the directory _and_ its contents; where as `+mybuildcontext/mydata/*` is required to copy all of the contents from within the `mydata` directory (
without copying the wrapping `mydata` directory).

If both the `Dockerfile` and build context are inside the same target, one must reference the same target twice, e.g. `FROM DOCKERFILE -f +target/dir/Dockerfile +target/dir`.
{% endhint %}

##### `--build-arg <key>=<value>`

Sets a value override of `<value>` for the Dockerfile build arg identified by `<key>`. This option is similar to the `docker build --build-arg <key>=<value>` option.

##### `--target <target-name>`

In a multi-stage Dockerfile, sets the target to be used for the build. This option is similar to the `docker build --target <target-name>` option.

##### `--platform <platform>`

Specifies the platform to build on.

For more information see the [multi-platform guide](../guides/multi-platform.md).

#### `--allow-privileged` (experimental)

{% hint style='info' %}
##### Note
The `--allow-privileged` flag has experimental status. To use this feature, it must be enabled via `VERSION --allow-privileged-from-dockerfile 0.8`.
{% endhint %}

When the Dockerfile build context points to an earthly artifact reference (e.g. `+mybuildcontext/mydata/*`), the `allow-privileged` flag will allow `RUN` commands under the referenced earthly target to make use of the `RUN --privileged` option.
This does not apply to Dockerfile's [RUN --security](https://docs.docker.com/reference/dockerfile/#run---security) flag.

## WITH DOCKER

#### Synopsis

```Dockerfile
WITH DOCKER [--pull <image-name>] [--load [<image-name>=]<target-ref>] [--compose <compose-file>]
            [--service <compose-service>] [--allow-privileged]
  <commands>
  ...
END
```

#### Description

The clause `WITH DOCKER` initializes a Docker daemon to be used in the context of a `RUN` command. The Docker daemon can be pre-loaded with a set of images using options such as `-pull` and `--load`. Once the execution of the `RUN` command has completed, the Docker daemon is stopped and all of its data is deleted, including any volumes and network configuration. Any other files that may have been created are kept, however.

If multiple targets are referenced via `--load`, the images are built in parallel. Similarly, multiple images referenced with `--pull` will be downloaded in parallel.

The clause `WITH DOCKER` automatically implies the `RUN --privileged` flag.

The `WITH DOCKER` clause only supports the command [`RUN`](#run). Other commands (such as `COPY`) need to be run either before or after `WITH DOCKER ... END`. In addition, only one `RUN` command is permitted within `WITH DOCKER`. However, multiple shell commands may be stringed together using `;` or `&&`.

A typical example of a `WITH DOCKER` clause might be:

```Dockerfile
FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
WORKDIR /test
COPY docker-compose.yml ./
WITH DOCKER \
        --compose docker-compose.yml \
        --load image-name:latest=(+some-target --SOME_BUILD_ARG=value) \
        --load another-image-name:latest=+another-target \
        --pull some-image:latest
    RUN docker run ... && \
        docker run ... && \
        ...
END
```

For more examples, see the [Docker in Earthly guide](../guides/docker-in-earthly.md) and the [Integration testing guide](../guides/integration.md).

For information on using `WITH DOCKER` with podman see the [Podman guide](../guides/podman.md)

{% hint style='info' %}
##### Note
For performance reasons, it is recommended to use a Docker image that already contains `dockerd`. If `dockerd` is not found, Earthly will attempt to install it.

Earthly provides officially supported images such as `earthly/dind:alpine-3.19-docker-25.0.5-r0` and `earthly/dind:ubuntu-23.04-docker-25.0.2-1` to be used together with `WITH DOCKER`.
{% endhint %}

{% hint style='info' %}
##### Note
Note that the cleanup phase (after the `RUN` command has finished), does not occur when using a `LOCALLY` target, users should use `RUN docker run --rm ...` to have docker remove the image after execution.
{% endhint %}

#### Options

##### `--pull <image-name>`

Pulls the Docker image `<image-name>` from a remote registry and then loads it into the temporary Docker daemon created by `WITH DOCKER`.

This option may be repeated in order to provide multiple images to be pulled.

{% hint style='info' %}
##### Note
It is recommended that you avoid issuing `RUN docker pull ...` and use `WITH DOCKER --pull ...` instead. The classical `docker pull` command does not take into account Earthly caching and so it would redownload the image much more frequently than necessary.
{% endhint %}

##### `--load [<image-name>=]<target-ref>`

Builds the image referenced by `<target-ref>` and then loads it into the temporary Docker daemon created by `WITH DOCKER`. Within `WITH DOCKER`, the image can be referenced as `<image-name>`, if specified, or otherwise by the name of the image specified in the referenced target's `SAVE IMAGE` command.

`<target-ref>` may be a simple target reference (`+some-target`), or a target reference with a build arg `(+some-target --SOME_BUILD_ARG=value)`.

This option may be repeated in order to provide multiple images to be loaded.

The `WITH DOCKER --load` option does not mark any saved images or artifacts of the referenced target for local output, nor does it mark any push commands of the referenced target for pushing. For that, please use [`BUILD`](#build).

##### `--compose <compose-file>`

Loads the compose definition defined in `<compose-file>`, adds all applicable images to the pull list and starts up all applicable compose services within.

This option may be repeated, thus having the same effect as repeating the `-f` flag in the `docker-compose` command.

##### `--service <compose-service>`

Specifies which compose service to pull and start up. If no services are specified and `--compose` is used, then all services are pulled and started up.

This option can only be used if `--compose` has been specified.

This option may be repeated in order to specify multiple services.

##### `--platform <platform>`

Specifies the platform for any referenced `--load` and `--pull` images.

For more information see the [multi-platform guide](../guides/multi-platform.md).

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

##### `--pass-args`

Same as [`FROM --pass-args`](#pass-args).

##### `--build-arg <key>=<value>` (**deprecated**)

This option is deprecated. Please use `--load <image-name>=(<target-ref> --<build-arg-key>=<build-arg-value>)` instead.

## IF

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

#### Examples

A very common pattern is to use the POSIX shell `[ ... ]` conditions. For example the following marks port `8080` as exposed if the file `./foo` exists.

```Dockerfile
IF [ -f ./foo ]
  EXPOSE 8080
END
```

It is also possible to call other commands, which can be useful for more comparisons such as semantic versioning. For example:

```Dockerfile
VERSION 0.8

test:
  FROM python:3
  RUN pip3 install semver

  # The following python script requires two arguments (v1 and v2)
  # and will return an exit code of 0 when v1 is semantically greater than v2
  # or an exit code of 1 in all other cases.
  RUN echo "#!/usr/bin/env python3
import sys
import semver
v1 = sys.argv[1]
v2 = sys.argv[2]
if semver.compare(v1, v2) > 0:
  sys.exit(0)
sys.exit(1)
  " > ./semver-gt && chmod +x semver-gt

  # Define two different versions
  ARG A="0.3.2"
  ARG B="0.10.1"

  # and compare them
  IF ./semver-gt "$A" "$B"
    RUN echo "A ($A) is semantically greater than B ($B)"
  ELSE
    RUN echo "A ($A) is NOT semantically greater than B ($B)"
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
    FROM alpine:3.18
ELSE
    FROM ubuntu:20.04
END
```

The reason this is invalid is because the `IF` condition is actually running the `/usr/bin/[` executable to test if the condition is true or false, and therefore requires that a valid build environment has been initialized.

Here is how this might be fixed.

```Dockerfile
ARG base=alpine
FROM busybox
IF [ "$base" = "alpine" ]
    FROM alpine:3.18
ELSE
    FROM ubuntu:20.04
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

##### `--secret <env-var>=<secret-id>`

Same as [`RUN --secret <env-var>=<secret-id>`](#secret-less-than-env-var-greater-than-less-than-secret-id-greater-than).

## FOR

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

##### `--secret <env-var>=<secret-id>`

Same as [`RUN --secret <env-var>=<secret-id>`](#secret-less-than-env-var-greater-than-less-than-secret-id-greater-than).

## WAIT

#### Synopsis

* ```
  WAIT
    <wait-block>
  END
  ```

#### Description

The `WAIT` clause executes the encapsulated commands and waits for them to complete.
This includes pushing and outputting local artifacts -- a feature which can be used to control the order of interactions with the outside world.

Even though the `WAIT` clause limits parallelism by forcing everything within it to finish executing before continuing, the commands **within** a `WAIT` block execute in parallel.

#### Examples

As an example, a `WAIT` block can be used to build and push to a remote registry (in parallel), then, after that execute a script which requires those images to exist in the remote registry:

```Dockerfile
myimage:
  ...
  SAVE IMAGE --push user/img:tag

myotherimage:
  ...
  SAVE IMAGE --push user/otherimg:tag

WAIT
  BUILD +myimg
  BUILD +myotherimg
END
RUN --push ./deploy ...
```

One can also use a `WAIT` block to control the order in which a `SAVE ARTIFACT ... AS LOCAL` command is executed:

```Dockerfile
RUN ./generate > data
WAIT
  SAVE ARTIFACT data AS LOCAL output/data
END
RUN ./test data # even if this fails, data will have been output
```

## TRY (experimental)

{% hint style='info' %}
##### Note
The `TRY` command is currently incomplete and has experimental status. To use this feature, it must be enabled via `VERSION --try 0.8`.
{% endhint %}

#### Synopsis

* ```
  TRY
    <try-block>
  FINALLY
    <finally-block>
  END
  ```

#### Description

The `TRY` clause executes commands within the `<try-block>`, while ensuring that the `<finally-block>` is always executed, even if the `<try-block>` fails.

This clause is still under active development. For now, only a single `RUN` command is permitted within the `<try-block>`, and only one or more `SAVE ARTIFACT` commands are permitted in the `<finally-block>`. The clause is thus useful for outputting coverage information in unit testing, outputting screenshots in UI integration tests, or outputting `junit.xml`, or similar.

#### Example

```Dockerfile
VERSION --try 0.8

example:
    FROM ...
    TRY
        # only a single RUN command is currently supported
        RUN ./test.sh
    FINALLY
        # only SAVE ARTIFACT commands are supported here
        SAVE ARTIFACT junit.xml AS LOCAL ./
    END
```

## CACHE

#### Synopsis

* ```
  CACHE [--sharing <sharing-mode>] [--chmod <octal-format>] [--id <cache-id>] [--persist] <mountpoint>
  ```

#### Description

The `CACHE` command creates a cache mountpoint at `<mountpoint>` in the build environment. The cache mountpoint is a directory which is shared between the instances of the same build target. The contents of the cache mountpoint are preserved between builds, and can be used to share data across builds.

#### Options

##### `--sharing <sharing-mode>`

The sharing mode for the cache mount, from one of the following:

* `locked` (default) - the cache mount is locked for the duration of the execution, other concurrent builds will wait for the lock to be released.
* `shared` - the cache mount is shared between all concurrent builds.
* `private` - if another concurrent build attempts to use the cache, a new (empty) cache will be created for the concurrent build.

##### `--chmod <octal-format>`

The permission of the mounted folder, in octal format (the same format the chmod unix command line expects).
Default `--chmod 0644`


##### `--id <cache-id>`

The cache ID for a global cache volume to be used across other targets or Earthfiles.

##### `--persist`

Make a copy of the cache available to any children that inherit from this target, by copying the contents of the cache to the child image.

{% hint style='warning' %}
Caches were persisted by default in version 0.7, which led to bloated images being pushed to registries. Version 0.8 changed the default behavior
to prevent copying the contents to children targets unless explicitly enabled by the newly added `--persist` flag.
{% endhint %}

## LOCALLY

#### Synopsis

* `LOCALLY`

#### Description

The `LOCALLY` command can be used in place of a `FROM` command, which will cause earthly to execute all commands under the target directly
on the host system, rather than inside a container. Commands within a `LOCALLY` target will never be cached.
This feature should be used with caution as locally run commands have no guarantee they will behave the same on different systems.

`LOCALLY` defined targets only support a subset of commands (along with a subset of their flags): `RUN`, `RUN --push`, `SAVE ARTIFACT`, and `COPY`.

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

## FUNCTION

#### Synopsis

* `FUNCTION`

#### Description

{% hint style='hint' %}
#### UDCs have been renamed to Functions

Functions used to be called UDCs (User Defined Commands). Earthly 0.7 uses `COMMAND` instead of `FUNCTION`.
{% endhint %}

The command `FUNCTION` marks the beginning of a function definition. Functions are reusable sets of instructions that can be inserted in targets or other functions. In order to reference and execute a function, you may use the command [`DO`](#do).

Unlike performing a `BUILD +target`, functions inherit the build context and the build environment from the caller.

Functions create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +MY_FUNCTION --<build-arg-key>=<build-arg-value>`.

Global imports and global args are inherited from the `base` target of the same Earthfile where the command is defined in (this may be distinct from the `base` target of the caller).

For more information see the [Functions Guide](../guides/functions.md).

## DO

#### Synopsis

* `DO [--allow-privileged] <function-ref> [--<build-arg-key>=<build-arg-value>...]`

#### Description

The command `DO` expands and executes the series of commands contained within a function [referenced by `<function-ref>`](../guides/importing.md#function-reference).

Unlike performing a `BUILD +target`, functions inherit the build context and the build environment from the caller.

Functions create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +MY_FUNCTION --<build-arg-key>=<build-arg-value>`.

For more information see the [Functions Guide](../guides/functions.md).

#### Options

##### `--allow-privileged`

Same as [`FROM --allow-privileged`](#allow-privileged).

##### `--pass-args`

Same as [`FROM --pass-args`](#pass-args).

## IMPORT

#### Synopsis

* `IMPORT [--allow-privileged] <earthfile-ref> [AS <alias>]`

#### Description

The command `IMPORT` aliases an Earthfile reference (`<earthfile-ref>`) that can be used in subsequent [target, artifact or command references](../guides/importing.md).

If not provided, the `<alias>` is inferred automatically as the last element of the path provided in `<earthfile-ref>`. For example, if `<earthfile-ref>` is `github.com/foo/bar/buz:v1.2.3`, then the alias is inferred as `buz`.

The `<earthfile-ref>` can be a reference to any directory other than `.`. If the reference ends in `..`, then mentioning `AS <alias>` is mandatory.

If an `IMPORT` is defined in the `base` target of the Earthfile, then it becomes a global `IMPORT` and it is made available to every other target or command in that file, regardless of their base images used.

For more information see the [importing guide](../guides/importing.md).

#### Options

##### `--allow-privileged`

Similar to [`FROM --allow-privileged`](#allow-privileged), extend the ability to request privileged capabilities to all invocations of the imported alias.

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

The `WORKDIR` command sets the working directory for following commands in the recipe. The working directory is also persisted as the default directory for the image. If the directory does not exist, it is automatically created. This command works the same way as the [Dockerfile `WORKDIR` command](https://docs.docker.com/engine/reference/builder/#workdir).

## HEALTHCHECK (same as Dockerfile HEALTHCHECK)

#### Synopsis

* `HEALTHCHECK NONE` (disable health checking)
* `HEALTHCHECK [--interval=DURATION] [--timeout=DURATION] [--start-period=DURATION] [--retries=N] [--start-interval=DURATION] CMD command arg1 arg2` (check container health by running command inside the container)

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

##### `--start-interval=DURATION`

Sets the time interval between health checks during the start period. Defaults to `5s`.

## HOST

#### Synopsis

* `HOST <hostname> <ip>`

#### Description

The `HOST` command creates a hostname entry (under `/etc/hosts`) that causes `<hostname>` to resolve to the specified `<ip>` address.

## SHELL (not supported)

The classical [`SHELL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#shell) is not yet supported. Use the *exec form* of `RUN`, `ENTRYPOINT` and `CMD` instead and prepend a different shell.

## ADD (not supported)

The classical [`ADD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#add) is not yet supported. Use [COPY](#copy) instead.

## ONBUILD (not supported)

The classical [`ONBUILD` Dockerfile command](https://docs.docker.com/engine/reference/builder/#onbuild) is not supported.

## STOPSIGNAL (not supported)

The classical [`STOPSIGNAL` Dockerfile command](https://docs.docker.com/engine/reference/builder/#stopsignal) is not yet supported.
