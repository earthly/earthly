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
| `EARTHLY_TARGET` | The canonical reference of the current target. | For example, for a target named `foo`, which exists on `john/work` branch, in a repository at `github.com/bar/buz`, in a subdirectory `src`, the canonical reference would be `github.com/bar/buz/src:john/work+foo`. For more information about canonical references, see [target referencing](../guides/target-ref.md). |
| `EARTHLY_TARGET_PROJECT` | The project part of the canonical reference of the current target. | For the example above, the canonical project would be `github.com/bar/buz/src:john` |
| `EARTHLY_TARGET_PROJECT_NO_TAG` | The project part of the canonical reference of the current target, but without the tag. | For the example above, this would be `github.com/bar/buz/src` |
| `EARTHLY_TARGET_NAME` | The name part of the canonical reference of the current target. | For the example above, the name would be `foo` |
| `EARTHLY_TARGET_TAG` | The tag part of the canonical reference of the current target. Note that if the target has no [canonical form](../guides/target-ref.md#canonical-form), the value is an empty string. | For the example above, the tag would be `john/work` |
| `EARTHLY_TARGET_TAG_DOCKER` | The tag part of the canonical reference of the current target, sanitized for safe use as a docker tag. This is guaranteed to be a valid docker tag, even if no canonical form exists, in which case, `latest` is used. | For the example above, the docker tag would be `john_work` |
| `EARTHLY_GIT_HASH` | The git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666ade67b29e42bef121144456d3977a67a` |
| `EARTHLY_GIT_SHORT_HASH` | The first 8 characters of the git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666` |
| `EARTHLY_GIT_ORIGIN_URL` | The git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. Please note that this may be inconsistent, depending on whether an HTTPS or SSH URL was used. | `git@github.com:bar/buz.git` or `https://github.com/bar/buz.git` |
| `EARTHLY_GIT_PROJECT_NAME` | The git project name from within the git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `bar/buz` |
| `EARTHLY_GIT_COMMIT_TIMESTAMP` | The timestamp, as unix seconds, of the git commit detected within the build context directory. If no git directory is detected, then the value is an empty string. | `1626881847` |
| `TARGETPLATFORM` | (**experimental**) The target platform the target is being built for. | `linux/arm/v7`, `linux/amd64`, `linux/arm64` |
| `TARGETOS` | (**experimental**) The target OS the target is being built for. | `linux` |
| `TARGETARCH` | (**experimental**) The target processor architecture the target is being built for. | `arm`, `amd64`, `arm64` |
| `TARGETVARIANT` | (**experimental**) The target processor architecture variant the target is being built for. | `v7` |

{% hint style='info' %}
##### Note

The classical Dockerfile predefined args are currently not available in Earthly.
{% endhint %}
