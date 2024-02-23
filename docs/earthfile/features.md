# Version-specific features

Earthly makes use of feature flags to release new and experimental features.
Some features must be explicitly enabled to use them.

{% hint style='danger' %}
##### Important
Avoid using feature flags for critical workflows. You should only use feature flags for testing new experimental features. By using feature flags you are opting out of forwards/backwards semver compatibility guarantees. This means that running the same script in a different environment, with a different version of Earthly may result in a different behavior (i.e. it'll work on your machine, but may break the build for your colleagues or for the CI).
{% endhint %}

Earthly uses [semantic versioning](http://semver.org/); once a new feature
has reached stability, a new `VERSION` release will include the feature enabled by default.

## Difference between the Earthly binary version and the Earthfile version

Earthly binary versions and Earthfile versions (declared via `VERSION`) follow the same minor versioning milestones, but are not the same.

The Earthly binary is able to run some older Earthfiles, but newer Earthfiles are not able to run on older Earthly binaries. The table below shows the compatibility matrix:

| Earthly binary version | Supported Earthfile VERSIONs |
|------------------------|--------------------------------|
| 0.8.x | `VERSION 0.6`, `VERSION 0.7`, `VERSION 0.8` |
| 0.7.x | `VERSION 0.5`, `VERSION 0.6`, `VERSION 0.7` |
| 0.6.x | No version specified (0.5 implied), `VERSION 0.5`, `VERSION 0.6` |
| 0.5.x | No version specified (0.5 implied), `VERSION 0.5` |
| <0.5.x | `VERSION` not supported |

## Upgrading to a newer version

In order to upgrade to `VERSION 0.8` safely, follow these steps:

1. If you are still using `VERSION 0.5`, upgrade those Earthfiles to `VERSION 0.6` or `VERSION 0.7`.
2. Upgrade your Earthly binary to 0.8 in CI and across your team. The Earthly 0.8 binary can run both `VERSION 0.6` and `VERSION 0.7` Earthfiles.
3. Once everyone is using the Earthly 0.8 binary, upgrade your Earthfiles one by one to `VERSION 0.8`. It is ok to have a mix of `VERSION 0.6`, `VERSION 0.7` and `VERSION 0.8` Earthfiles in the same project. Earthly handles that gracefully.

When upgrading between `VERSION`s, keep in mind that you will encounter backwards-incompatible changes. Check out the change log of each version for more information.

* [0.8](https://github.com/earthly/earthly/releases/tag/v0.8.0-rc1)
* [0.7](https://github.com/earthly/earthly/releases/tag/v0.7.0)
* [0.6](https://github.com/earthly/earthly/releases/tag/v0.6.0)
* [0.5](https://github.com/earthly/earthly/releases/tag/v0.5.0)

## Specifying Version and features

Each Earthfile should list the current earthly version it depends on using the [`VERSION`](../earthfile/earthfile.md#version) command.
The `VERSION` command was first introduced under `0.5` and is required as of `0.7`.

```Dockerfile
VERSION [<flags>...] <version-number>
```

## Feature flags

| Feature flag                        | Status       | Description                                                                                                        |
|-------------------------------------|--------------|--------------------------------------------------------------------------------------------------------------------|
| `--use-registry-for-with-docker`    | 0.5          | Makes use of the embedded BuildKit Docker registry (instead of tar files) for `WITH DOCKER` loads and pulls        |
| `--use-copy-include-patterns`       | 0.6          | Speeds up COPY transfers                                                                                           |
| `--referenced-save-only`            | 0.6          | Changes the behavior of SAVE commands in a significant way                                                         |
| `--for-in`                          | 0.6          | Enables support for `FOR ... IN ...` commands                                                                      |
| `--require-force-for-unsafe-saves`  | 0.6          | Requires `--force` for saving artifacts locally outside the Earthfile's directory                                  |
| `--no-implicit-ignore`              | 0.6          | Eliminates implicit `.earthlyignore` entries, such as `Earthfile` and `.tmp-earthly-out`                           |
| `--earthly-version-arg`             | 0.7          | Enables builtin ARGs: `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA`                                                    |
| `--shell-out-anywhere`              | 0.7          | Allows shelling-out in any earthly command (including in the middle of `ARG`)                                      |
| `--explicit-global`                 | 0.7          | Base target args must have a `--global` flag in order to be considered global args                                 |
| `--check-duplicate-images`          | 0.7          | Check for duplicate images during output                                                                           |
| `--use-cache-command`               | 0.7          | Allow use of `CACHE` command in Earthfiles                                                                         |
| `--use-host-command`                | 0.7          | Allow use of `HOST` command in Earthfiles                                                                          |
| `--use-copy-link`                   | 0.7          | Use the equivalent of `COPY --link` for all copy-like operations                                                   |
| `--new-platform`                    | 0.7          | Enable new platform behavior                                                                                       |
| `--no-tar-build-output`             | 0.7          | Do not print output when creating a tarball to load into `WITH DOCKER`                                             |
| `--use-no-manifest-list`            | 0.7          | Enable the `SAVE IMAGE --no-manifest-list` option                                                                  |
| `--use-chmod`                       | 0.7          | Enable the `COPY --chmod` option                                                                                   |
| `--earthly-locally-arg`             | 0.7          | Enable the `EARTHLY_LOCALLY` arg                                                                                   |
| `--use-project-secrets`             | 0.7          | Enable project-based secret resolution                                                                             |
| `--use-pipelines`                   | 0.7          | Enable the `PIPELINE` and `TRIGGER` commands                                                                       |
| `--earthly-git-author-args`         | 0.7          | Enable the `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` args                                                  |
| `--wait-block`                      | 0.7          | Enable the `WAIT` / `END` block commands                                                                           |
| `--no-network`                      | 0.8 | Allow the use of `RUN --network=none` commands                                                                     |
| `--arg-scope-and-set`               | 0.8 | Enable the `LET` / `SET` commands and nested `ARG` scoping                                                         |
| `--use-docker-ignore`               | 0.8 | Enable the use of `.dockerignore` files in `FROM DOCKERFILE` targets                                               |
| `--pass-args`                       | 0.8 | Enable the optional `--pass-args` flag for the `BUILD`, `FROM`, `COPY`, `WITH DOCKER --load` commands              |
| `--global-cache`                    | 0.8 | Enable global caches (shared across different Earthfiles), for cache mounts and `CACHE` commands having an ID      |
| `--cache-persist-option`            | 0.8 | Adds `CACHE --persist` option to persist cache content in images, Changes default `CACHE` behaviour to not persist |
| `--use-function-keyword`            | 0.8 | Enable using `FUNCTION` instead of `COMMAND` when declaring a function |
| `--use-visited-upfront-hash-collection` | 0.8 | Switches to a newer target parallelization algorithm |
| `--no-use-registry-for-with-docker` | Experimental | Disable `use-registry-for-with-docker`                                                                             |
| `--try`                             | Experimental | Enable the `TRY` / `FINALLY` / `END` block commands                                                                |
| `--earthly-ci-runner-arg`           | Experimental | Enable the `EARTHLY_CI_RUNNER` builtin ARG                                                                         |

Note that the features flags are disabled by default in Earthly versions lower than the version listed in the "status" column above.

##### `--use-copy-include-patterns`

*Speeds up COPY transfers.*

When enabled, Earthly will only send the files listed for the specific [`COPY`](../earthfile/earthfile.md#copy) command.
Without this feature, Earthly sends the entire directory of files excluding files listed in the [`.earthlyignore` file](../earthfile/earthlyignore.md).

##### `--referenced-save-only`

*Changes the behavior of SAVE commands in a significant way*

When enabled, Earthly will output artifacts resulting from `SAVE ARTIFACT ... AS LOCAL ...` and images resulting from `SAVE IMAGE` and also execute `RUN --push` commands only if they are connected to the main target through a chain of `BUILD` commands.

For example, chains like these will produce outputs (and possibly push, if enabled):

* main target -> `SAVE`
* main target -> `BUILD -> SAVE`
* main target -> `BUILD -> BUILD -> SAVE`
* main target -> `BUILD -> BUILD -> BUILD -> SAVE`

While chains like these will NOT produce outputs nor would they push:

* main target -> `FROM -> SAVE`
* main target -> `COPY -> SAVE`
* main target -> `FROM -> BUILD -> SAVE`
* main target -> `BUILD -> FROM -> SAVE`
* main target -> `BUILD -> BUILD -> COPY -> SAVE`

This works the same regardless of whether the targets in the chain are remote or local.

When this feature is **disabled**, Earthly will output artifacts and images regardless of whether they are connected to the main target through a chain of `BUILD` commands, however the outputs will be subject to the following rules:

* All `SAVE ARTIFACT ... AS LOCAL ...`, with local Earthfiles will be output
* `SAVE ARTIFACT ... AS LOCAL ...` produced in remote targets will not be output
* All images with tag names (both local and remote Earthfiles) will be output
* No image will be pushed or `RUN --push` command will be executed if the target is remote

##### `--for-in`

*Enables support for `FOR ... IN ...` commands*

When enabled, Earthly will allow the use of `FOR ... IN ...` commands.
