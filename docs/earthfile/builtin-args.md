# Builtin args

Builtin args are variables with values automatically filled-in by Earthly.

The value of a builtin arg can never be overriden. However, you can always have an additional `ARG`, which takes as the default value, the value of the builtin arg. The additional arg can be overriden. Example

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
| `EARTHLY_TARGET_PROJECT` | The project part of the canonical reference of the current target. | For the example above, the canonical project would be `github.com/bar/buz/src` |
| `EARTHLY_TARGET_NAME` | The name part of the canonical reference of the current target. | For the example above, the name would be `foo` |
| `EARTHLY_TARGET_TAG` | The tag part of the canonical reference of the current target. Note that if the target has no [canonical form](../guides/target-ref.md#canonical-form)), the value is an empty string. | For the example above, the tag would be `john/work` |
| `EARTHLY_TARGET_TAG_DOCKER` | The tag part of the canonical reference of the current target, sanitized for safe use as a docker tag. This is guaranteed to be a valid docker tag, even if no canonical form exists, in which case, `latest` is used. | For the example above, the docker tag would be `john_work` |
| `EARTHLY_GIT_HASH` | The git hash detected within the build context directory. If no git directory is detected, then the value is an empty string. Take care when using this arg, as the frequently changing git hash may be cause for not using the cache. | `41cb5666ade67b29e42bef121144456d3977a67a` |
| `EARTHLY_GIT_ORIGIN_URL` | The git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `git@github.com:earthly/earthly.git` |
| `EARTHLY_GIT_PROJECT_NAME` | The git project name from within the git URL detected within the build context directory. If no git directory is detected, then the value is an empty string. | `earthly/earthly` |

{% hint style='info' %}
##### Note

The classical Dockerfile predefined args are currently not available in Earthly.
{% endhint %}
