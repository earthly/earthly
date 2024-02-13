# Builtin args

Builtin args are variables with values automatically filled-in by Earthly.

The value of a builtin arg can never be overridden. However, you can always have an additional `ARG`, which takes as the default value, the value of the builtin arg. The additional arg can be overridden. Example

```Dockerfile
ARG EARTHLY_TARGET_TAG
ARG TAG=$EARTHLY_TARGET_TAG
SAVE IMAGE --push some/name:$TAG
```

{% hint style='danger' %}
##### Important
Earthly builtin args need to be pre-declared before they can be used. For example

```Dockerfile
ARG EARTHLY_TARGET
RUN echo "The current target is $EARTHLY_TARGET"
```
{% endhint %}

### General args

| Name | Description | Example value |
| --- | --- | --- |
| `EARTHLY_CI` | Whether the build is being executed in --ci mode. | `true`, `false` |
| `EARTHLY_BUILD_SHA` | The git hash of the commit which built the currently running version of Earthly. | `1a9eda7a83af0e2ec122720e93ff6dbe9231fc0c` |
| `EARTHLY_LOCALLY` | Whether the target is being executed `LOCALLY`. | `true`, `false` |
| `EARTHLY_PUSH` | Whether `earthly` was called with the `--push` flag, or not. | `true`, `false` |
| `EARTHLY_VERSION` | The version of Earthly currently running. | `v0.8.0` |

### Target-related args

| Name | Description | Example value |
| --- | --- | --- |
| `EARTHLY_TARGET_NAME` | The name part of the canonical reference of the current target. | For the target `github.com/bar/buz/src:john/work+foo`, the name would be `foo` |
| `EARTHLY_TARGET_PROJECT_NO_TAG` | The project part of the canonical reference of the current target, but without the tag. | For the target `github.com/bar/buz/src:john/work+foo`, this would be `github.com/bar/buz/src` |
| `EARTHLY_TARGET_PROJECT` | The project part of the canonical reference of the current target. | For the target `github.com/bar/buz/src:john/work+foo`, the canonical project would be `github.com/bar/buz/src:john` |
| `EARTHLY_TARGET_TAG_DOCKER` | The tag part of the canonical reference of the current target, sanitized for safe use as a docker tag. This is guaranteed to be a valid docker tag, even if no canonical form exists, in which case, `latest` is used. | For the target `github.com/bar/buz/src:john/work+foo`, the docker tag would be `john_work` |
| `EARTHLY_TARGET_TAG` | The tag part of the canonical reference of the current target. Note that if the target has no [canonical form](../guides/importing.md#canonical-form), the value is an empty string. | For the target `github.com/bar/buz/src:john/work+foo`, the tag would be `john/work` |
| `EARTHLY_TARGET` | The canonical reference of the current target. | For example, for a target named `foo`, which exists on `john/work` branch, in a repository at `github.com/bar/buz`, in a subdirectory `src`, the canonical reference would be `github.com/bar/buz/src:john/work+foo`. For more information about canonical references, see [importing guide](../guides/importing.md). |

### Git-related args

| Name | Description | Example value |
| --- | --- | --- |
| `EARTHLY_GIT_AUTHOR` | The git author detected within the build context directory. If no git directory is detected, then the value is an empty string. | `John Doe <john@example.com>` |
| `EARTHLY_GIT_CO_AUTHORS` | The git co-authors detected within the build context directory, separated by space. If no git directory is detected, then the value is an empty string. | `Jane Doe <jane@example.com Jack Smith <jack@example.com>` |
| `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` | The author timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `1626881847` |
| `EARTHLY_GIT_BRANCH` | The git branch of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `main` |
| `EARTHLY_GIT_COMMIT_TIMESTAMP` | The committer timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `1626881847` |
| `EARTHLY_GIT_HASH` | The git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666ade67b29e42bef121144456d3977a67a` |
| `EARTHLY_GIT_ORIGIN_URL` | The git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. Please note that this may be inconsistent, depending on whether an HTTPS or SSH URL was used. | `git@github.com:bar/buz.git` or `https://github.com/bar/buz.git` |
| `EARTHLY_GIT_PROJECT_NAME` | The git project name from within the git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `bar/buz` |
| `EARTHLY_GIT_REFS` | The git references of the git commit detected within the build context directory, separated by space. If no git directory is detected, then the value is an empty string. | `issue-2735-git-ref main` |
| `EARTHLY_GIT_SHORT_HASH` | The first 8 characters of the git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666` |
| `EARTHLY_SOURCE_DATE_EPOCH` | The timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is `0` (the unix epoch) | `1626881847`, `0` |

### Platform-related args

| Name | Description | Example value |
| --- | --- | --- |
| `NATIVEARCH` | The native processor architecture of the build runner. | `arm`, `amd64`, `arm64` |
| `NATIVEOS` | The native OS of the build runner. | `linux` |
| `NATIVEPLATFORM` | The native platform of the build runner. | `linux/arm/v7`, `linux/amd64`, `darwin/arm64` |
| `NATIVEVARIANT` | The native processor architecture variant of the build runner. | `v7` |
| `TARGETARCH` | The target processor architecture the target is being built for. | `arm`, `amd64`, `arm64` |
| `TARGETOS` | The target OS the target is being built for. | `linux` |
| `TARGETPLATFORM` | The target platform the target is being built for. This defaults to the native platform. | `linux/arm/v7`, `linux/amd64`, `linux/arm64` |
| `TARGETVARIANT` | The target processor architecture variant the target is being built for. | `v7` |
| `USERARCH` | The processor architecture of the user (the environment the `earthly` binary is invoked from). | `arm`, `amd64`, `arm64` |
| `USEROS` | The OS of the user (the environment the `earthly` binary is invoked from). | `darwin` |
| `USERPLATFORM` | The platform of the user (the environment the `earthly` binary is invoked from). | `darwin/amd64`, `linux/amd64`, `darwin/arm64` |
| `USERVARIANT` | The processor architecture variant of the user (the environment the `earthly` binary is invoked from). | `v7` |

The default value of the `TARGETPLATFORM` arg is the native platform of the runner, for non-LOCALLY targets. This can be overriden by using the `--platform` flag, when using the `earthly` CLI. For example, `earthly --platform linux/amd64 +my-target` will set the `TARGETPLATFORM` arg to `linux/amd64`. You can also override the target platform in an Earthfile, when issuing `BUILD` commands. For example, `BUILD --platform linux/amd64 +my-target`. Or you can override the platform within the target definition by setting the platform in the `FROM` statement. For example `FROM --platform linux/amd64 alpine:3.13`.

Under `LOCALLY`, the `TARGETPLATFORM` arg is always set to the user platform (the environment the `earthly` binary is invoked from) and it is not overriden by the `--platform` flag.

{% hint style='info' %}
##### Note
Under `LOCALLY` targets, it is important to declare the `TARGETPLATFORM` arg **after** the `LOCALLY` command, to ensure that it gets the approriate user platform value. For example:

```Dockerfile
my-target:
    LOCALLY
    ARG TARGETPLATFORM
    RUN echo "The target platform under LOCALLY is $TARGETPLATFORM"
```
{% endhint %}
