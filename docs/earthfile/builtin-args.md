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

In contrast to Dockerfile predefined args, Earthly builtin args need to be pre-declared before they can be used. For example

```Dockerfile
ARG EARTHLY_TARGET
RUN echo "The current target is $EARTHLY_TARGET"
```
{% endhint %}

The following builtin args are available

| Name | Description | Example value |
| --- | --- | --- |
| `EARTHLY_BUILD_SHA` | The git hash of the commit which built the currently running version of Earthly. | `1a9eda7a83af0e2ec122720e93ff6dbe9231fc0c` |
| `EARTHLY_CI` | Whether the build is being executed in --ci mode. | `true`, `false` |
| `EARTHLY_CI_RUNNER` | Whether the build is being executed within Earthly CI. | `true`, `false` |
| `EARTHLY_GIT_AUTHOR` | The git author detected within the build context directory. If no git directory is detected, then the value is an empty string. | `John Doe <john@example.com>` |
| `EARTHLY_GIT_CO_AUTHORS` | The git co-authors detected within the build context directory, separated by space. If no git directory is detected, then the value is an empty string. | `Jane Doe <jane@example.com Jack Smith <jack@example.com>` |
| `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` | The author timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `1626881847` |
| `EARTHLY_GIT_BRANCH` | The git branch of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. Must be enabled with the `VERSION --git-branch 0.7` feature flag. | `main` |
| `EARTHLY_GIT_COMMIT_TIMESTAMP` | The committer timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `1626881847` |
| `EARTHLY_GIT_HASH` | The git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666ade67b29e42bef121144456d3977a67a` |
| `EARTHLY_GIT_ORIGIN_URL` | The git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. Please note that this may be inconsistent, depending on whether an HTTPS or SSH URL was used. | `git@github.com:bar/buz.git` or `https://github.com/bar/buz.git` |
| `EARTHLY_GIT_PROJECT_NAME` | The git project name from within the git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `bar/buz` |
| `EARTHLY_GIT_SHORT_HASH` | The first 8 characters of the git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666` |
| `EARTHLY_LOCALLY` | Whether the target is being executed `LOCALLY`. | `true`, `false` |
| `EARTHLY_PUSH` | Whether `earthly` was called with the `--push` flag, or not. | `true`, `false` |
| `EARTHLY_SOURCE_DATE_EPOCH` | The timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is `0` (the unix epoch) | `1626881847`, `0` |
| `EARTHLY_TARGET_NAME` | The name part of the canonical reference of the current target. | For the target `github.com/bar/buz/src:john/work+foo`, the name would be `foo` |
| `EARTHLY_TARGET_PROJECT_NO_TAG` | The project part of the canonical reference of the current target, but without the tag. | For the target `github.com/bar/buz/src:john/work+foo`, this would be `github.com/bar/buz/src` |
| `EARTHLY_TARGET_PROJECT` | The project part of the canonical reference of the current target. | For the target `github.com/bar/buz/src:john/work+foo`, the canonical project would be `github.com/bar/buz/src:john` |
| `EARTHLY_TARGET_TAG_DOCKER` | The tag part of the canonical reference of the current target, sanitized for safe use as a docker tag. This is guaranteed to be a valid docker tag, even if no canonical form exists, in which case, `latest` is used. | For the target `github.com/bar/buz/src:john/work+foo`, the docker tag would be `john_work` |
| `EARTHLY_TARGET_TAG` | The tag part of the canonical reference of the current target. Note that if the target has no [canonical form](../guides/target-ref.md#canonical-form), the value is an empty string. | For the target `github.com/bar/buz/src:john/work+foo`, the tag would be `john/work` |
| `EARTHLY_TARGET` | The canonical reference of the current target. | For example, for a target named `foo`, which exists on `john/work` branch, in a repository at `github.com/bar/buz`, in a subdirectory `src`, the canonical reference would be `github.com/bar/buz/src:john/work+foo`. For more information about canonical references, see [target referencing](../guides/target-ref.md). |
| `EARTHLY_VERSION` | The version of Earthly currently running. | `v0.6.2` |
| `TARGETARCH` | The target processor architecture the target is being built for. | `arm`, `amd64`, `arm64` |
| `TARGETOS` | The target OS the target is being built for. | `linux` |
| `TARGETPLATFORM` | The target platform the target is being built for. | `linux/arm/v7`, `linux/amd64`, `linux/arm64` |
| `TARGETVARIANT` | The target processor architecture variant the target is being built for. | `v7` |
| `USERARCH` | The processor architecture the target is being built from. | `arm`, `amd64`, `arm64` |
| `USEROS` | The OS the target is being built from. | `linux`, `darwin` |
| `USERPLATFORM` | The platform the target is being built from. | `linux/arm/v7`, `linux/amd64`, `darwin/arm64` |
| `USERVARIANT` | The processor architecture variant the target is being built from. | `v7` |

{% hint style='info' %}
##### Note

The classical Dockerfile predefined args are currently not available in Earthly.
{% endhint %}
