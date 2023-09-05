# Version-specific features

Earthly makes use of feature flags to release new and experimental features.
Some features must be explicitly enabled to use them.

Earthly uses [semantic versioning](http://semver.org/); once a new feature
has reached stability, a new major or minor version of Earthly will be released with
the feature enabled by default.

## Specifying Version and features

Each Earthfile should list the current earthly version it depends on using the [`VERSION`](../earthfile/earthfile.md#version) command.
The `VERSION` command was first introduced under `0.5` and is required as of `0.7`.

```Dockerfile
VERSION [<flags>...] <version-number>
```

### Example

To Future-proof Earthfiles it is recommended to add a `VERSION` command. Consider a case where an Earthfile is developed
against earthly `v0.5.23` and makes use of the experimental `FOOBAR` command, the first line of the Earthfile should be:

```Dockerfile
VERSION --foobar 0.5
```

This will ensure that backwards-breaking features that are introduced in a later version will not change how this Earthfile is interpreted.

In a future release (e.g. `0.X`), the `FOOBAR` command **might** be promoted from the _experimental_ stage to _stable_ stage,
at that point, version `0.X` would automatically set the `--foobar` flag to `true`, and the Earthfile could be updated
to require version `0.X` (or later), and could be rewritten as `VERSION 0.X`.

## Feature flags

| Feature flag | status | description |
| --- | --- | --- |
| `--use-registry-for-with-docker` | 0.5 | Makes use of the embedded BuildKit Docker registry (instead of tar files) for `WITH DOCKER` loads and pulls |
| `--use-copy-include-patterns` | 0.6 | Speeds up COPY transfers |
| `--referenced-save-only` | 0.6 | Changes the behavior of SAVE commands in a significant way |
| `--for-in` | 0.6 | Enables support for `FOR ... IN ...` commands |
| `--require-force-for-unsafe-saves` | 0.6 | Requires `--force` for saving artifacts locally outside the Earthfile's directory  |
| `--no-implicit-ignore` | 0.6 | Eliminates implicit `.earthlyignore` entries, such as `Earthfile` and `.tmp-earthly-out` |
| `--earthly-version-arg` | 0.7 | Enables builtin ARGs: `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA` |
| `--shell-out-anywhere` | 0.7 | Allows shelling-out in any earthly command (including in the middle of `ARG`) |
| `--explicit-global` | 0.7 | Base target args must have a `--global` flag in order to be considered global args |
| `--check-duplicate-images` | 0.7 | Check for duplicate images during output |
| `--use-cache-command` | 0.7 | Allow use of `CACHE` command in Earthfiles |
| `--use-host-command` | 0.7 | Allow use of `HOST` command in Earthfiles |
| `--use-copy-link` | 0.7 | Use the equivalent of `COPY --link` for all copy-like operations |
| `--new-platform` | 0.7 | Enable new platform behavior |
| `--no-tar-build-output` | 0.7 | Do not print output when creating a tarball to load into `WITH DOCKER` |
| `--use-no-manifest-list` | 0.7 | Enable the `SAVE IMAGE --no-manifest-list` option |
| `--use-chmod` | 0.7 | Enable the `COPY --chmod` option |
| `--earthly-locally-arg` | 0.7 | Enable the `EARTHLY_LOCALLY` arg |
| `--use-project-secrets` | 0.7 | Enable project-based secret resolution |
| `--use-pipelines` | 0.7 | Enable the `PIPELINE` and `TRIGGER` commands |
| `--earthly-git-author-args` | 0.7 | Enable the `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` args |
| `--wait-block` | 0.7 | Enable the `WAIT` / `END` block commands |
| `--earthly-ci-runner-arg` | 0.7 | Enable the `EARTHLY_CI_RUNNER` builtin ARG |
| `--use-docker-ignore` | 0.7 | Enable the use of `.dockerignore` files in `FROM DOCKERFILE` targets |
| `--try` | Experimental | Enable the `TRY` / `FINALLY` / `END` block commands |
| `--arg-scope-and-set` | Experimental | Enable the `LET` / `SET` commands and nested `ARG` scoping |
| `--pass-args` | Experimental | Enable the optional `--pass-args` flag for the `BUILD`, `FROM`, `COPY`, `WITH DOCKER --load` commands |


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
