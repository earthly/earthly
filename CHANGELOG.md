# Earthly Changelog

All notable changes to [Earthly](https://github.com/earthly/earthly) will be documented in this file.

## Unreleased

## v0.8.12 - 2024-05-23

### Added
- An experimental modification of the buildkit scheduler, which attempts to solve the `inconsistent graph state` error, which can be enabled locally with `earthly --ticktock ...`.

### Changed
- The BYOC (bring your own cloud) commands have been updated to reflect server-side API changes.

### Fixed
- The `--buildkit-container-name` flag was incorrectly being ignored when `--no-buildkit-update` was set.

### Additional Info
- This release includes changes to buildkit

## v0.8.11 - 2024-05-16

### Added
- Support for using HTTP(S) proxies when connecting to satellites.

### Fixed
- Backwards compatability issue where `WITH DOCKER` would fail with `EARTHLY_DOCKERD_CACHE_DATA: parameter not set` when using an older version of the earthly in combination with a satellite running v0.8.10.

### Additional Info
- This release includes changes to buildkit

## v0.8.10 - 2024-05-14

### Added
- New Github Actions Workflow commands integration `--github-annotations` flag or GITHUB_ACTIONS=true env. [#2189](https://github.com/earthly/earthly/issues/2189)
- Added a new `--oidc` flag to `RUN` command which allows authentication to AWS via OIDC. Enable with the `VERSION --run-with-aws-oidc` feature flag. [#3804](https://github.com/earthly/earthly/issues/3804)
- Experimental `WITH DOCKER --cache-id=<key>` feature, which will cache the contents of the docker data root, resulting in faster `--load` and `--pull` execution. Enabled with the `VERSION --docker-cache` feature flag. [#3399](https://github.com/earthly/earthly/issues/3399)
- New `SAVE IMAGE --without-earthly-labels` feature, which will prevent any `dev.earthly.*` labels from being saved to the image. Enable with the `VERSION --allow-without-earthly-labels` feature flag. Thanks to [@3manuek](https://github.com/3manuek) for the contribution!

### Fixed
- `WITH DOCKER` load time calculation. [#3485](https://github.com/earthly/earthly/issues/3485)
- The earthly cli was not correctly setting the exit status on failures when executing a `RUN` on a satellite which reached the max execution time limit.
- Self-hosted satellite connection issue.

### Changed
- Earthly will now use source link format when displaying errors, e.g. `<path>:<line>:<col>` rather than `<path> line <line>:<col>`.
- Improved error messages for cases where a shell is required to run a command such as `IF`, `FOR`, etc.
- Earthly will now show a warning when earthly anonymously connects to a registry (which increases the chance of being rate-limited).

### Additional Info
- This release includes changes to buildkit

## v0.8.9 - 2024-04-24

### Fixed

- `BUILD --auto-skip` was recording failed steps as complete, which would lead to them being skipped on subsequent runs. [#4054](https://github.com/earthly/earthly/issues/4054)

### Additional Info
- This release has no changes to buildkit

## v0.8.8 - 2024-04-17

### Added

- New experimental wildcard-based copy, e.g. `COPY ./services/*+artifact/* .` which would invoke `COPY` for `./services/foo+artifact`, and `./services/bar+artifact` (assuming two services foo and bar, both having a `artifact` target in their respective Earthfile). Enable with the `VERSION --wildcard-copy` feature flag. [#3966](https://github.com/earthly/earthly/issues/3966).
- New built-in `ARG`s - `EARTHLY_GIT_AUTHOR_EMAIL` and `EARTHLY_GIT_AUTHOR_NAME` will contain the author email and author name respectively. Enable with the `VERSION --git-author-email-name-args` feature flag.
- New `--raw-output` flag available on `RUN` that outputs line without target name. Enable with `VERSION --raw-output`. [#3713](https://github.com/earthly/earthly/issues/3713)

### Changed

- `EARTHLY_GIT_AUTHOR` built-in `ARG` will now contain both name and email, when enabled with the `VERSION --git-author-email-name-args` feature flag. Previously it only contained the email. [#3822](https://github.com/earthly/earthly/issues/3822)

### Fixed

- Make `LET`/`SET` commands block parallel commands such as `BUILD` until the former are processed, similar to the behavior of `ARG`. [#3997](https://github.com/earthly/earthly/issues/3997)
- `LET`/`SET` commands were not properly handled with the use of Auto-skip. [#3996](https://github.com/earthly/earthly/issues/3996)

### Additional Info
- This release has no changes to buildkit

## v0.8.7 - 2024-04-03

### Added

- Warning log when resolving remote references using a git image that doesn't match Buildkit's architecture.
- New experimental `--exec-stats-summary=<path>` cli flag, which will display a summary of memory and cpu stats when earthly exits.
- A notice is now displayed when unnecessary feature flags are set (but already enabled by default by the VERSION number). Thanks to [@danqixu](https://github.com/danqixu) for the contribution! [#3641](https://github.com/earthly/earthly/issues/3641)
- A warning is displayed if the local buildkit image architecture does not match the host architecture. [#3937](https://github.com/earthly/earthly/issues/3937)

### Fixed

- Warning logs during HTTP retries are only displayed in `--debug` mode.
- The HOST command will now expand variables. Thanks to [@pbecotte](https://github.com/pbecotte) for the contribution! [#1743](https://github.com/earthly/earthly/issues/1743)
- runc has been updated to 1.1.12 in the buildkit fork

### Additional Info
- This release includes changes to buildkit

## v0.8.6 - 2024-03-18

### Added

- Ability to set arbitrary attributes which certain registries require to support explicit remote caching (via the `earthly --remote-cache` flag). [#3714](https://github.com/earthly/earthly/issues/3714) and [#3868](https://github.com/earthly/earthly/issues/3868)

### Fixed

- Fixed an issue in Auto-skip where a `+base` target's ARGs were not accounted for when calculating the cache. [#3895](https://github.com/earthly/earthly/issues/3895)

### Additional Info
- This release has no changes to buildkit

## v0.8.5 - 2024-03-11

### Added

- Added `--aws` flag to `RUN` command which makes AWS environment variables or ~/.aws available. Enable with the `VERSION --run-with-aws` feature flag. [#3803](https://github.com/earthly/earthly/issues/3803)
- Added `--allow-privileged` flag to `FROM DOCKERFILE` command. Enable with the `VERSION --allow-privileged-from-dockerfile` feature flag. Thanks to [@dustyhorizon](https://github.com/dustyhorizon) for the contribution! [#3706](https://github.com/earthly/earthly/issues/3706)

### Fixed

- Fixes an issue where wildcard `BUILD`'s are invoked from a relative directory (e.g., an `Earthfile` containing `BUILD ./*+test` invoked with `earthly ./rel-dir+target`). [#3840](https://github.com/earthly/earthly/issues/3840)
- `--pass-args` will no longer pass builtin args, which would result in `value cannot be specified for built-in build arge errors. [#3775](https://github.com/earthly/earthly/issues/3775)
- Fixes a parsing issue with `BUILD` flag arguments and wildcard targets [#3862](https://github.com/earthly/earthly/issues/3862)
- `BUILD --auto-skip` was silently ignored when the feature flag (`VERSION --build-auto-skip`) was missing [#3870](https://github.com/earthly/earthly/issues/3870)
- Fix an issue where `COPY --if-exists` would fail if the non-existing directory includes a wildcard. [#3875](https://github.com/earthly/earthly/issues/3875)
- Fixes an issue with passing the correct org value to Logstream which resulted in missing logs in the web builds view (https://cloud.earthly.dev/your-org/builds).
- Rename `UDC` to `FUNCTION` in hint when a secret is not found.

### Additional Info
- This release includes changes to buildkit

## v0.8.4 - 2024-02-21

### Added

- The internal `dockerd-wrapper.sh` script, which is used to implement `WITH DOCKER`, will execute `/usr/share/earthly/dockerd-wrapper-pre-script`, if present, prior to starting the
  inner dockerd process. This can be used to configure options that depend on the host's kernel at run-time.
- Auto-skip can now be used directly on `BUILD` commands with `BUILD --auto-skip`. [#3581](https://github.com/earthly/earthly/issues/3581)

### Changed

- Satellite `rm` requires a `--force` flag if it's running. This should help protect users from accidental deletes.

### Fixed

- Fixes an issue with the registry proxy (used for faster image & artifact exporting) on Docker Desktop for Windows/WSL. [#3769](https://github.com/earthly/earthly/issues/3769)
- Fixes a problem with cache IDs not being expanded. For example: `CACHE --id $MY_ARG` was not using the assigned value of `$MY_ARG`.

### Additional Info
- This release includes changes to buildkit

## v0.8.3 - 2024-01-31

### Fixed

- `EARTHLY_GIT_REFS` was incorrectly returning all references which contained the commit rather than pointed to the current commit. This also increases performance of looking up the branches. [#3752](https://github.com/earthly/earthly/issues/3752)
- Fixes an issue where `earthly account login --token` was leading to partially created auth config files. [#3761](https://github.com/earthly/earthly/issues/3761)

### Additional Info
- This release includes changes to buildkit

## v0.8.2 - 2024-01-25

### Added

- Added a `--force` flag to the `satellite update` command, which forces a satellite to sleep before starting the update process. This may forcibly kill ongoing builds currently running on the satellite.

### Changed

- Changed the default buildkit cache size to be adaptively set to 20GB, which is then clamped between the range of 10%-55% of the disk size.
  This logic can expressed as `min(55%, max(10%, 20GB))`.
- Satellites are now put to sleep before updating via `earthly sat update <satellite-name>`.

### Fixed

- Fixed an intermittent issue with the registry proxy support container failing immediately on Mac. [#3740](https://github.com/earthly/earthly/issues/3740)
- Fixed a problem with parsing empty results when cleaning up old registry proxy support containers on Mac.
- Fixed a case where a suggested command would incorrectly contain both `--interative` and `--ci`. [#3746](https://github.com/earthly/earthly/issues/3746)
- Disabled the registry proxy server when Earthly is run from within a container. [#3736](https://github.com/earthly/earthly/issues/3736)

### Additional Info
- This release has no changes to buildkit

## v0.8.1 - 2024-01-23

### Added

- Added a new `--disable-remote-registry-proxy` cli flag, which can be used to disable the remote registry proxy, which is used by earthly when performing a `SAVE IMAGE`
  command with a satellite / remote buildkit instance. This will cause earthly to use the slower tar-based loading of docker images. [#3736](https://github.com/earthly/earthly/issues/3736)
- A new warning if Earthly is configured with a cache size less than 10GB; running with a small cache size may lead to unexpected cache misses.

### Additional Info
- This release has no changes to buildkit

## v0.8.0 - 2024-01-22

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.8` at the top of your Earthfile.

**Migrating from 0.7**

If you are using Earthly 0.7, follow the following steps to migrate:

1. If you are still using `VERSION 0.5`, upgrade those Earthfiles to `VERSION 0.6` or `VERSION 0.7`.
2. Upgrade your Earthly binary to 0.8 in CI and across your team. The Earthly 0.8 binary can run both `VERSION 0.6` and `VERSION 0.7` Earthfiles (but `VERSION 0.5` support has been dropped).
3. Once everyone is using the Earthly 0.8 binary, upgrade your Earthfiles one by one to `VERSION 0.8`. It is ok to have a mix of `VERSION 0.6`, `VERSION 0.7` and `VERSION 0.8` Earthfiles in the same project. Earthly handles that gracefully. See changes below for information on backwards incompatible changes when migrating from `VERSION 0.7` to `VERSION 0.8`.

This process helps manage the backward breaking changes with minimal disruption.

**Summary**

Declaring `VERSION 0.8` is equivalent to

```
VERSION \
  --arg-scope-and-set \
  --cache-persist-option \
  --git-refs \
  --global-cache \
  --no-network \
  --pass-args \
  --use-docker-ignore \
  --use-function-keyword \
  --use-visited-upfront-hash-collection \
  0.7
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

It should be noted that some of these features break backwards compatibility. See below.

### Changed

- Redeclaring an `ARG` in the same scope as a previous declaration is now an error.
- `ARG`s inside of targets will no longer have their default value overridden by global `ARG`s.
- Declaring a `CACHE ...` in a target will no longer be copied to children targets when referenced via a `FROM +...`; to persist the contents of the cache, it is now required to use the `CACHE --persist ...` flag.
- The `COMMAND` keyword has been renamed to `FUNCTION`.

### Added

- `LET` - Allows declaring a local variable. This command works similarly to `ARG` except that it cannot be overridden from the CLI. `LET` variables are allowed to shadow `ARG` variables, which allows you to promote an `ARG` to a local variable so that it may be used with `SET`.
- `SET` - a new command that allows changing the value of variables declared with `LET`.
- Outputting images from a remote runner has improved performance as it no longer transfers layers that are already present locally.
- [Auto-skip](https://docs.earthly.dev/v/earthly-0.8/docs/caching/caching-in-earthfiles#3.-auto-skip) has been promoted to *beta* status.
- `RUN --network=none` allows running a command without network access.
- `.dockerignore` files are now used in `FROM DOCKERFILE` targets.
- `DO --pass-args`, `BUILD --pass-args` etc allow passing all build arguments to external Earthfiles.
- `CACHE --id=...` and `RUN --mount type=cache,id=...` allows setting a custom cache mount ID, thus allowing sharing cache mounts globally across different targets.
- New satellite sizes: 2xlarge, 3xlarge, 4xlarge
- New experimental wildcard-based builds, e.g. `BUILD ./services/*+test` which would call `./services/foo+test`, and `./services/bar+test` (assuming two services foo and bar, both having a `test` target in their respective Earthfile) [#3582](https://github.com/earthly/earthly/issues/3582).

### Removed

- `VERSION 0.5` is now obsolete. Declaring `VERSION 0.5` is no longer supported, and will now raise an error.

### Fixed

- Parallelism is improved when running the same target with different arguments in certain cases (e.g. the target uses `WITH DOCKER`).
- Fixed a log sharing upload-resumption bug
- Fixed multiple issues with the lexer failing to parse certain characters in shell command substitution (`$()`) and single quoted strings.
  - Some escaped characters, like `\#`, were failing to parse when used inside shell expressions. Example: `$(echo "a#b#c" | cut -f2 -d\#)` [#3475](https://github.com/earthly/earthly/issues/3475)
  - Some characters, like `#`, were failing to parse when used inside single-quoted strings: Example: `'this is a # string'` [#1280](https://github.com/earthly/earthly/issues/1280)
- Fixed an issue where some escaped `ARG` shell expressions were being incorrectly preprocessed. Example: `$(echo "\"")` became `$(echo """)` [#3131](https://github.com/earthly/earthly/issues/3131)
- The `--pass-args` feature was not passing active arguments which were set via a default value.
- `SAVE ARTIFACT --if-exists` was not saving files based on a wildcard glob pattern. [#1679](https://github.com/earthly/earthly/issues/1679)
- `BUILD` was not expanding `--platform` argument values.

### Additional Info
- This release includes changes to buildkit

## v0.8.0-rc2 - 2024-01-09

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.8` at the top of your Earthfile.

**Migrating from 0.7**

If you are using Earthly 0.7, follow the following steps to migrate:

1. If you are still using `VERSION 0.5`, upgrade those Earthfiles to `VERSION 0.6` or `VERSION 0.7`.
2. Upgrade your Earthly binary to 0.8 in CI and across your team. The Earthly 0.8 binary can run both `VERSION 0.6` and `VERSION 0.7` Earthfiles (but `VERSION 0.5` support has been dropped).
3. Once everyone is using the Earthly 0.8 binary, upgrade your Earthfiles one by one to `VERSION 0.8`. It is ok to have a mix of `VERSION 0.6`, `VERSION 0.7` and `VERSION 0.8` Earthfiles in the same project. Earthly handles that gracefully. See changes below for information on backwards incompatible changes when migrating from `VERSION 0.7` to `VERSION 0.8`.

This process helps manage the backward breaking changes with minimal disruption.

**Summary**

Declaring `VERSION 0.8` is equivalent to

```
VERSION \
  --arg-scope-and-set \
  --cache-persist-option \
  --git-refs \
  --global-cache \
  --no-network \
  --pass-args \
  --use-docker-ignore \
  --use-function-keyword \
  --use-visited-upfront-hash-collection \
  0.7
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

It should be noted that some of these features break backwards compatibility. See below.

### Changed

- Redeclaring an `ARG` in the same scope as a previous declaration is now an error.
- `ARG`s inside of targets will no longer have their default value overridden by global `ARG`s.
- It is no longer possible to override a global ARG when calling a target.
- Declaring a `CACHE ...` in a target will no longer be copied to children targets when referenced via a `FROM +...`; to persist the contents of the cache, it is now required to use the `CACHE --persist ...` flag.
- The `COMMAND` keyword has been renamed to `FUNCTION`.

### Added

- `LET` - Allows declaring a local variable. This command works similarly to `ARG` except that it cannot be overridden from the CLI. `LET` variables are allowed to shadow `ARG` variables, which allows you to promote an `ARG` to a local variable so that it may be used with `SET`.
- `SET` - a new command that allows changing the value of variables declared with `LET`.
- Outputting images from a remote runner has improved performance as it no longer transfers layers that are already present locally.
- [Auto-skip](https://docs.earthly.dev/v/earthly-0.8/docs/caching/caching-in-earthfiles#3.-auto-skip) has been promoted to *beta* status.
- `RUN --network=none` allows running a command without network access.
- `.dockerignore` files are now used in `FROM DOCKERFILE` targets.
- `DO --pass-args`, `BUILD --pass-args` etc allow passing all build arguments to external Earthfiles.
- `CACHE --id=...` and `RUN --mount type=cache,id=...` allows setting a custom cache mount ID, thus allowing sharing cache mounts globally across different targets.
- New satellite sizes: 2xlarge, 3xlarge, 4xlarge
- New experimental wildcard-based builds, e.g. `BUILD ./services/*+test` which would call `./services/foo+test`, and `./services/bar+test` (assuming two services foo and bar, both having a `test` target in their respective Earthfile) [#3582](https://github.com/earthly/earthly/issues/3582).

### Removed

- `VERSION 0.5` is now obsolete. Declaring `VERSION 0.5` is no longer supported, and will now raise an error.

### Fixed

- Parallelism is improved when running the same target with different arguments in certain cases (e.g. the target uses `WITH DOCKER`).
- Fixed a log sharing upload-resumption bug

### Additional Info
- This release includes changes to buildkit

## v0.8.0-rc1 - 2024-01-03

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.8` at the top of your Earthfile.

**Migrating from 0.7**

If you are using Earthly 0.7, follow the following steps to migrate:

1. If you are still using `VERSION 0.5`, upgrade those Earthfiles to `VERSION 0.6` or `VERSION 0.7`.
2. Upgrade your Earthly binary to 0.8 in CI and across your team. The Earthly 0.8 binary can run both `VERSION 0.6` and `VERSION 0.7` Earthfiles (but `VERSION 0.5` support has been dropped).
3. Once everyone is using the Earthly 0.8 binary, upgrade your Earthfiles one by one to `VERSION 0.8`. It is ok to have a mix of `VERSION 0.6`, `VERSION 0.7` and `VERSION 0.8` Earthfiles in the same project. Earthly handles that gracefully. See changes below for information on backwards incompatible changes when migrating from `VERSION 0.7` to `VERSION 0.8`.

This process helps manage the backward breaking changes with minimal disruption.

**Summary**

Declaring `VERSION 0.8` is equivalent to

```
VERSION \
  --arg-scope-and-set \
  --cache-persist-option \
  --git-refs \
  --global-cache \
  --no-network \
  --pass-args \
  --use-docker-ignore \
  --use-function-keyword \
  --use-visited-upfront-hash-collection \
  0.7
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

It should be noted that some of these features break backwards compatibility. See below.

### Changed

- Redeclaring an `ARG` in the same scope as a previous declaration is now an error.
- `ARG`s inside of targets will no longer have their default value overridden by global `ARG`s.
- It is no longer possible to override a global ARG when calling a target.
- Declaring a `CACHE ...` in a target will no longer be copied to children targets when referenced via a `FROM +...`; to persist the contents of the cache, it is now required to use the `CACHE --persist ...` flag.
- The `COMMAND` keyword has been renamed to `FUNCTION`.

### Added

- `LET` - Allows declaring a local variable. This command works similarly to `ARG` except that it cannot be overridden from the CLI. `LET` variables are allowed to shadow `ARG` variables, which allows you to promote an `ARG` to a local variable so that it may be used with `SET`.
- `SET` - a new command that allows changing the value of variables declared with `LET`.
- Outputting images from a remote runner has improved performance as it no longer transfers layers that are already present locally.
- [Auto-skip](https://docs.earthly.dev/v/earthly-0.8/docs/caching/caching-in-earthfiles#3.-auto-skip) has been promoted to *beta* status.
- `RUN --network=none` allows running a command without network access.
- `.dockerignore` files are now used in `FROM DOCKERFILE` targets.
- `DO --pass-args`, `BUILD --pass-args` etc allow passing all build arguments to external Earthfiles.
- `CACHE --id=...` and `RUN --mount type=cache,id=...` allows setting a custom cache mount ID, thus allowing sharing cache mounts globally across different targets.

### Removed

- `VERSION 0.5` is now obsolete. Declaring `VERSION 0.5` is no longer supported, and will now raise an error.

### Fixed

- Parallelism is improved when running the same target with different arguments in certain cases (e.g. the target uses `WITH DOCKER`).

### Additional Info
- This release includes changes to buildkit

## v0.7.23 - 2023-12-18

### Added
- Auto-skip (*experimental*) - a feature that allows you to skip large parts of a build in certain situations, especially suited for monorepos. For more information see [the auto-skip section from Caching in Earthfiles](https://docs.earthly.dev/docs/caching/caching-in-earthfiles#auto-skip).
- A warning when a `COPY` destination includes a tilde (~). Related to [#1789](https://github.com/earthly/earthly/issues/1789).
- A hint message to suggest the usage of `-i` flag to debug the build when a RUN command fails.
- `start-interval` flag to `HEALTHCHECK` command for dockerfile parity [#3409](https://github.com/earthly/earthly/issues/3409).
- A verbose message indicating which authentication providers are used during a build.
- `ssh_command` config option which can be used to override the ssh command that is used by `git` when connecting to an ssh-based repository. Thanks to [@weaversam8](https://github.com/weaversam8) for the contribution!

### Fixed
- Limit the number of deprecation warnings when using `COMMAND` instead of `FUNCTION` keyword.
- Fixed an error which stated `VERSION 0.0` is a valid Earthfile version.

### Changed
- Changed the color used to print metadata values (such as ARGs values) in the build log to Faint Blue.
- Updated default alpine/git image to v2.40.1.
- When creating an auth token, an existing token will no longer be overwritten by default. To overwrite, the `--overwrite` flag should be used.

### Additional Info
- This release includes changes to buildkit

## v0.7.22 - 2023-11-27

### Added
- A new experimental `earthly --exec-stats` flag, which displays per-target execution stats such as total CPU and memory usage.
- A new experimental `earthly billing view` command to get information about the organization billing plan.
- Messages informing used build minutes during a build.
- Help message when a build fails due to a missing referenced cloud secret.

### Fixed
- Remove redundant verbose error messages that were not different from messages that were already being printed.
- Fixed `failed to sign challenge` errors when attempting to login using an ed25519 key with the 1Password ssh-agent. [#3366](https://github.com/earthly/earthly/issues/3366)

### Changed
- Final error messages for executions without a known target will be displayed without `_unknown *failed* |` prefix. and instead use `Error: ` as prefix more consistently.
- Failing `RUN` commands under `LOCALLY` will display the same format of error message for `RUN` without `LOCALLY` [#3356](https://github.com/earthly/earthly/issues/3356).
- Log sharing link will be printed last, even in case of a build error.
- Help message after a build error will be printed in color.
- Use dedicated logstream failure category for param related error.
- An authentication attempt with an expired auth token will result in a `auth token expired` error instead of `unauthorized`.
- A successful authentication with an auth token will display a warning with time left before token expires if it's 14 days or under.
- The command `earthly registry` will attempt to use the selected org if no org is specified.
- Clarify error messages when failing to pass secrets to a build.
- Provide information on how to get more build minutes when a build fails due to missing minutes.
- Provide information on how to increase the max number of allowed satellites when failing to launch a satellite.
- `CACHE` mounts will no longer depend on the contents of `ARG`s, and instead will be limited to the target name.
- Child targets will no longer receive the contents of mounted `CACHE` volumes defined in the parent target; this change can be enabled with `VERSION --cache-persist-option`. [#3509](https://github.com/earthly/earthly/issues/3509)
- Improved memory usage related to log messages by no longer pre-allocating log buffers; this is most noticeable for really large Earthfiles with lots of different targets.
- Updated buildkit with upstream changes up to 3d50b97793391d81d7bc191d7c5dd5361d5dadca.
- Improved speed of `SAVE IMAGE` exports when using a remote buildkit instance (e.g. satellite) from a MacOS host; this can be enabled with the `--use-remote-registry` option.

### Additional Info
- This release includes changes to buildkit

## v0.7.21 - 2023-10-24

### Added
- The new ARG `EARTHLY_GIT_REFS` will contain the references to the current git commit, this ARG must be enabled with the `VERSION --git-refs` feature flag. [#2735](https://github.com/earthly/earthly/issues/2735)
- A new `--force-certificate-generation` flag for bootstrapping, which will force the generation of self signed TLS certificates even when the `--no-buildkit` flag is set.

### Fixed
- Fixed reduced parallelism regression which occurred when the target is the same but has different args -- can be enabled with `VERSION --use-visited-upfront-hash-collection` [#2377](https://github.com/earthly/earthly/issues/2377)
- `prune --age` did not support `d` (for days) suffix, even thought `earthly --help` said it did [#3401](https://github.com/earthly/earthly/issues/3401)
- `buildkit scheduler error: return leaving incoming open` which occured during deduplication of opperations within buildkit; cherry-picked 100d3cb6b6903be50f7a3e5dba193515aa9530fa from upstream buildkit repo. [#2957](https://github.com/earthly/earthly/issues/2957)
- Changed `WITH DOCKER` to pull images in parallel [#2351](https://github.com/earthly/earthly/issues/2351)

### Changed
- Registry proxy: Use lower-level TCP streaming [#2351](https://github.com/earthly/earthly/pull/3317)

### Additional Info
- This release includes changes to buildkit

## v0.7.20 - 2023-10-03

### Added
- Support for `mode` in mount cache [#3278](https://github.com/earthly/earthly/issues/3278).
- Support for `mode` in CACHE commands [#3290](https://github.com/earthly/earthly/pull/3290).
- Experimental support for shared/global caches (cache `id` is no longer scoped per Earthfile) [#1129](https://github.com/earthly/earthly/issues/1129). Note that this is feature-flagged, and only changed when `VERSION --global-cache 0.7` is defined.

### Fixed
- A regression where URLs will not always get shorter when used as a prefix. Partially addresses [#3200](https://github.com/earthly/earthly/issues/3200).
- If a build fails because of `qemu` missing, earthly will display a proper hint to install it [#3200](https://github.com/earthly/earthly/issues/3200).
- Removed erroneous error-message which said error: 0 errors occured [#3306](https://github.com/earthly/earthly/pull/3306).
- A race condition when exiting interactive debugger mode resulting in confusing errors [#3200](https://github.com/earthly/earthly/issues/3200).
- Docker auto-install script failures related to attempts to read from tty, while verifying docker's pgp key [#3324](https://github.com/earthly/earthly/pull/3324).
- Issue affecting pulling images in Podman [#2471](https://github.com/earthly/earthly/issues/2471).
- A `panic: send on closed channel` error would sometimes occur during shutdown of the logstream [#3325](https://github.com/earthly/earthly/pull/3325).

### Changed
- Some error messages at the end of an execution will only be displayed in verbose mode (`earthly -V ...`), e.g. `Error: build target: build main: failed to solve:`... [#3200](https://github.com/earthly/earthly/issues/3200)
- `GIT CLONE` URLs will only be printed once as part of a prefix, e.g. `+my-clone-target(https://g/e/earthly) | --> GIT CLONE (--branch ) https://github.com/earthly/earthly`
- Clarify errors in interactive debugger so that they won't be confused with the build errors [#3200](https://github.com/earthly/earthly/issues/3200).
- The `WITH DOCKER` auto-install script will now pass the `--no-tty` option to `gpg` [#3288](https://github.com/earthly/earthly/issues/3288).

### Additional Info
- This release includes changes to buildkit

## v0.7.19 - 2023-09-20

### Added
- Added "dev.earthly.*" LABELS to saved images, for example `dev.earthly.version` will be set to `v0.7.19` (or whatever version of earthly is used) [#3247](https://github.com/earthly/earthly/issues/3247).
- Added option to verbose print known_hosts to make it easier to debug git related commands [#3234](https://github.com/earthly/earthly/issues/3234).

### Fixed
- When a project based secret is not found, the name of the secret will now be displayed along with the "not found" error.

### Changed
- Log sharing will now stream logs as your build is running (rather than uploading logs when build execution completes).
- Satellite reserve calls will now retry on error [#3255](https://github.com/earthly/earthly/issues/3255).
- Display warning when TLS is disabled.

### Additional Info
- This release has no changes to buildkit

## v0.7.18 - 2023-09-18 (aborted release/not recommended)
<!--changelog-parser-ignore-start-->
Note: This release was aborted due to a regression in the log sharing functionality
<!--changelog-parser-ignore-end-->

### Added
- Added "dev.earthly.*" LABELS to saved images, for example `dev.earthly.version` will be set to `v0.7.18` (or whatever version of earthly is used) [#3247](https://github.com/earthly/earthly/issues/3247).
- Added option to verbose print known_hosts to make it easier to debug git related commands [#3234](https://github.com/earthly/earthly/issues/3234).

### Fixed
- When a project based secret is not found, the name of the secret will now be displayed along with the "not found" error.

### Changed
- Refactor console output code (e.g. removed redundant output, prepared code for a future streaming log uploads... coming soon).
- Display warning when TLS is disabled.

## v0.7.17 - 2023-08-30

### Added
- Added a `--pass-arg` flag that can be used with `BUILD`, `FROM`, `COPY`, `WITH DOCKER --load`, or `DO`, which will pass all build arguments to external Earthfiles. [#1891](https://github.com/earthly/earthly/issues/1891)

## v0.7.16 - 2023-08-28

### Fixed
- Fixed a cgroup v2 related bug that affected systemd-based images (such as kind) from being run via `WITH DOCKER`. [#3159](https://github.com/earthly/earthly/issues/3159)

### Changed
- Removed redundant output when parts of builds are re-used; the `--verbose` flag will still display the output.
- Calling `earthly secret set <path>` (when run interactively) will now prompt for a single-line secret if no other flags are given.
- fixed bug in `earthly registry setup` which was waiting for an end of file (eof) rather than newline, when prompting for a password.

### Added
- Added additional error message output when buildkit scheduller errors occur (in order to help debug the ongoing [2957](https://github.com/earthly/earthly/issues/2957) issue).

## v0.7.15 - 2023-08-04

### Fixed
- Fixed a bug in `WITH DOCKER` which prevented the use of newer versions of docker. [#3164](https://github.com/earthly/earthly/issues/3164)

## v0.7.14 - 2023-07-31

### Changed
- Update buildkit (contains upstream changes up to 687091bb6c8aaa0185cdc570c4db3db533f329d0).
- Use `HTTPS_PROXY` env when connecting to earhly cloud API.

## v0.7.13 - 2023-07-26

### Added
- `earthly account list-tokens` now shows the last time a token was used
- Experimental command `earthly init` to initialize an Earthfile in a project (currently supporting only golang projects)

### Fixed
- Fixed a bug, where the command to create tokens with a set expiration failed.
- Long pauses at the end of builds, which were characterized by apparent freezes or delays with the message `Waiting on Buildkit...`.
- `earthly account create-token` no longer panics when parsing expiration date
- `earthly account login` could change the active user when the JWT expired and an SSH key existed for a different user; now earthly will either refresh the JWT or error

### Changed
- Setting env vars like  `FORCE_COLOR`, or `EARTHLY_FULL_TARGET` to `0`, `false`, `FALSE`, or `` (an empty-string) will no longer force the color, use any other value like `1`, `true`, or `yesplease`.
- `earthly org list` now shows the currently selected org

## v0.7.12 - 2023-07-17

### Added
- warning if acquiring file-lock takes longer than 3 seconds.

### Changed
- improved error message when a 429 too many requests rate limit error occurs.
- `earthly sat ls -a` shows last accessed time
- improved output for listing auth tokens

### Fixed
- make use of org from earthly config when using satellite commands.

## v0.7.12-rc1 - 2023-07-13

### Added
- warning if acquiring file-lock takes longer than 3 seconds.

### Changed
- improved error message when a 429 too many requests rate limit error occurs.
- `earthly sat ls -a` shows last accessed time

### Fixed
- make use of org from earthly config when using satellite commands.

## v0.7.11 - 2023-07-06

### Added
- `global.org` configuration value to set a default org for all `earthly` commands that require it.
- `earthly org select` and `earthly org unselect` commands, as shortcuts to set a default organization in the `earthly` config file.

### Changed
- Removed the default size in satellite launch (the default size is now determined by the backend when not provided) [#3057](https://github.com/earthly/earthly/issues/3057)
- Deprecated the satellite org configuration value. It uses the new global configuration value.

## v0.7.10 - 2023-07-05

### Changed
- Removed the default size in satellite launch (the default size is now determined by the backend when not provided) [#3057](https://github.com/earthly/earthly/issues/3057)
- Earthly cloud organization auto-detection has been deprecated and should now be explicitly set with the `--org` flag or with the `EARTHLY_ORG` environment variable.
- Buildkit has been updated to include upstream changes up to cdf28d6fff9583a0b173c62ac9a28d1626599d3b.

### Fixed
- Updated the podman auth provider to better understand podman `auth.json` locations. [#3038](https://github.com/earthly/earthly/issues/3038)
- Fixed our aggregated authprovider ignoring the cloud authprovider when a project is set after the first creds lookup [#3058](https://github.com/earthly/earthly/issues/3058)

## v0.7.9 - 2023-06-22

### Changed
- The command `docker-build` now also supports passing multiple platforms using a comma (e.g `--platform linux/amd64,linux/arm64`)
- Increased temporary lease duration of buildkit's history queue to prevent unknown history in blob errors under high cpu load. [#3000](https://github.com/earthly/earthly/issues/3000)
- Performing an `earthly account logout` will keep you logged out -- earthly will no longer attempt an auto-login via ssh-agent (use `earthly account login` to log back in).

### Fixed
- Fixed a bug in satellite update command which was incorrectly changing satellites to medium size.
- Fixed support for being authenticated with multiple registries when using the cloud-based `earthly registry` feature. [#3010](https://github.com/earthly/earthly/issues/3010)
- Fixed `WITH DOCKER` auto install script when using latest (bookworm) version.

### Added
- Buildkit logs now include version and revision.
- Satellite name autocompletion

## v0.7.8 - 2023-06-07

### Added
- Add a new command `docker-build` to build a docker image using a Dockerfile without using an Earthfile, locally or on a satellite.

### Changed
- `FROM DOCKERFILE` will use a `.dockerignore` file when using a build context from the host system and both `.earthlyignore` and `.earthignore` do not exist. Enable with `VERSION --use-docker-ignore 0.7`.

### Fixed
- Fixed upstream race condition bug in buildkit, which resulted in `failed to solve: unknown blob sha256:<...> in history` errors. [#3000](https://github.com/earthly/earthly/issues/3000)

## v0.7.7 - 2023-06-01

### Added
- The new ARG `EARTHLY_CI_RUNNER` indicates whether the current build is executed in Earthly CI. Enable with `VERSION --earthly-ci-runner-arg 0.7`.

### Changed
- Updated buildkit up to 60d134bf7 and fsutil up to 9e7a6df48576; this includes a buildkit fix for 401 Unauthorized errors. [#2973](https://github.com/earthly/earthly/issues/2973)
- Enabled `GIT_LFS_SKIP_SMUDGE=1` when pulling git repos (to avoid pulling in large files initially).

### Fixed
- The earthly docker image incorrectly showed `dev-main` under the embedded buildkit version.

## v0.7.6 - 2023-05-23

### Added
- Better error messages when git opperations fail.
- Added a `runc-ps` script under the earthly-buildkitd container to make it easier to see what processes are running.

### Fixed
- The builtin 'docker compose' (rather than `docker-compose` script) is now used when using the `WITH DOCKER` command under alpine 3.18 or greater.
- Fixed context timeout value overflow when connecting to a remote buildkit instance.

## v0.7.5 - 2023-05-10

### Changed
- Remote BuildKit will use TLS by default.
- Deprecation warning: Secret IDs naming scheme should follow the ARG naming scheme; i.e. a letter followed by alphanumeric characters or underscores. [#2883](https://github.com/earthly/earthly/issues/2883)
- Secrets take precedence over ARGs of the same name. [#2931](https://github.com/earthly/earthly/issues/2931)

### Added
- Experimental support for performing a `git lfs pull --include=<path>` when referencing a remote target on the cli, when used with the new `--git-lfs-pull-include` flag. [#2992](https://github.com/earthly/earthly/pull/2922)

### Fixed
- `SAVE IMAGE <img>` was incorrectly pushed when earthly was run with the `--push` cli flag (this restores the requirement that images that are pushed must be defined with `SAVE IMAGE --push <img>`). [#2923](https://github.com/earthly/earthly/issues/2923)
- Incorrect global ARG values when chaining multiple DO commands together. [#2920](https://github.com/earthly/earthly/issues/2920)
- Build args autocompletion under artifact mode.

## v0.7.4 - 2023-04-12

### Changed
- Updated the github ssh-rsa public key in the pre-populated buildkitd known_hosts entries.

## v0.7.3 - 2023-04-12

### Added
- A host of changes to variables under the `--arg-scope-and-set` feature flag:
  - Redeclaring an `ARG` in the same scope as a previous declaration is now an error.
  - `ARG`s inside of targets will no longer have their default value overridden by global `ARG`s.
  - A new command, `LET`, is available for declaring non-argument variables.
    - `LET` takes precedence over `ARG`, just like `ARG` takes precedence over `ARG --global`.
  - A new command, `SET`, is available for changing the value of variables declared with `LET`.
- Introduced `--size` and `--age` flags to the prune command, to allow better control.

### Changed

- Updated buildkit with changes up to 3187d2d056de7e3f976ef62cd548499dc3472a7e.
- The `VERSION --git-branch` feature flag has been removed (`EARTHLY_GIT_BRANCH` was always available in the previous version).
- Improved earthly API connection timeout logic.
- `earthly doc` now includes `ARG`s in both summary and detail output, and `ARTIFACT`s and `IMAGE`s in its detail output.

### Fixed

- Fixed `Could not detect digest for image` warnings for when using `WITH DOCKER --load` which referenced an earthly target that
  included a `FROM` referencing an image following the `docker.io/<user>/<img>` naming scheme (rather than the `docker.io/library/<user>/<img>` scheme).
- Fixed `COPY --if-exists` to work with earthly targets.  [#2541](https://github.com/earthly/earthly/issues/2541)
- Intentional-indentation of comments is no longer removed by the doc command. [#2747](https://github.com/earthly/earthly/issues/2747)
- `SAVE ARTIFACT ... AS LOCAL ...` could not write to non-current directories upon failure of a TRY/FINALLY block. [#2800](https://github.com/earthly/earthly/issues/2800)

## v0.7.2 - 2023-03-14

### Added

- Support for [Rosetta](https://developer.apple.com/documentation/apple-silicon/about-the-rosetta-translation-environment) translation environment (emulator) in buildkit as an alternative to QEMU. To enable, go to Docker Desktop -> Settings -> Features in development -> Check `Use Rosetta for x86/amd64 emulation on Apple Silicon`.
- New ARG `EARTHLY_GIT_BRANCH` will contain the branch of the current git commit, this ARG must be enabled with the `VERSION --git-branch` feature flag. [#2735](https://github.com/earthly/earthly/pull/2735)
- Verbose logging when git configurations perform a regex substitution.

### Fixed

- SAVE IMAGE --push did not always work under `VERSION 0.7`, when image was refrenced by a `FROM` or `COPY`, followed by a `BUILD`. [#2762](https://github.com/earthly/earthly/issues/2762)

### Changed

- Simplified error message when a RUN command fails with an exit code. [#2742](https://github.com/earthly/earthly/issues/2742)
- Improved warning messages when earthly cloud-based registry auth fails. [#2783](https://github.com/earthly/earthly/issues/2783)
- Deleting a project will prompt for confirmation, unless --force is specified.
- Updated buildkit with changes up to 4451e1be0e6889ffc56225e54f7e26bd6fdada54.

## v0.7.1 - 2023-03-01

### Added

- Support for `RUN --network=none`, which prevents programs from using any network resources. [#834](https://github.com/earthly/earthly/issues/834)

### Changed

- The `unexpected env` warning can now be silenced by creating a `.arg` or `.secret` file. [#2696](https://github.com/earthly/earthly/issues/2696)

### Fixed

- Unindented comments in the middle of recipe blocks no longer cause parser errors. [#2697](https://github.com/earthly/earthly/issues/2697)

## v0.7.0 - 2023-02-21

The documentation for this version is available at the [Earthly 0.7 documentation page](https://docs.earthly.dev/v/earthly-0.7/).

**Earthly CI**

Earthly 0.7 is the first version compatible with Earthly CI.

Earthly 0.7 introduces the new keywords `PIPELINE` and `TRIGGER` to help define Earthly CI pipelines.

```
my-pipeline:
    PIPELINE --push
    TRIGGER push main
    TRIGGER pr main
    BUILD +my-target
```

For more information on how to use `PIPELINE` and `TRIGGER`, please see the [reference documentation](https://docs.earthly.dev/v/earthly-0.7/docs/earthfile#pipeline-beta).

**Podman support**

Podman support has now been promoted out of *beta* status and is generally available in 0.7. Earthly will automatically detect the container frontend, whether that's `docker` or `podman` and use it automatically for running Buildkit locally, or for outputting images locally resulting from the build.

Please note that rootful podman is required. Rootless podman is not supported.

**VERSION is now mandatory**

The `VERSION` command
is now required for all Earthfiles, and an error will occur if it is missing. If you are not ready to update your
Earthfiles to use 0.7 (or 0.6), you can declare `VERSION 0.5` to continue to use your Earthfiles.

**.env file is no longer used for `ARG` or secrets**

The `.env` file will only be used to automatically export environment variables, which can be used to configure earthly command line flags.
As a result, values will no longer be propagated to Earthfile `ARG`s or `RUN --secret=...` commands.

Instead if you want build arguments or secrets automatically passed into earthly, they must be placed in `.arg` or `.secret` files respectively.

Note that this is a **backwards incompatible** change and will apply to all Earthfiles (regardless of the defined `VERSION` value).

**Pushing no longer requires everything else to succeed**

The behavior of the `--push` mode has changed in `VERSION 0.7` and is backwards incompatible with `VERSION 0.6`. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, to allow for more flexible push ordering via the new `WAIT` clause. To achieve the behavior of the previous `--push` mode, you need to wrap any pre-required commands in a `WAIT` clause. For example, to push an image only if tests have passed, you would do the following:

```Earthfile
test-and-push:
  WAIT
    BUILD +test
  END
  BUILD +my-image
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

This type of behavior is useful in order to have better control over the order of push operations. For example, you may want to push an image to a registry, followed by a deployment that uses the newly pushed image. Here is how this might look like:

```Earthfile
push-and-deploy:
  ...
  WAIT
    BUILD +my-image
  END
  RUN --push ./deploy.sh my-org/my-image:latest
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

Where `./deploy.sh` is custom deployment script instructing a production environment to start using the image that was just pushed.

**Promoting experimental features**

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.7` at the top of your Earthfile.

Declaring `VERSION 0.7` is equivalent to

```
VERSION \
  --check-duplicate-images \
  --earthly-git-author-args \
  --earthly-locally-arg \
  --earthly-version-arg \
  --explicit-global \
  --new-platform \
  --no-tar-build-output \
  --save-artifact-keep-own \
  --shell-out-anywhere \
  --use-cache-command \
  --use-chmod \
  --use-copy-link \
  --use-host-command \
  --use-no-manifest-list \
  --use-pipelines \
  --use-project-secrets \
  --wait-block \
  0.6
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Changed

- The behavior of the `--push` mode has changed in a backwards incompatible manner. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, allowing push commands to execute in the middle of the build now. Previously under `VERSION --wait-block 0.6`.
- `ARG`s declared in the base target do not automatically become global unless explicitly declared as such via `ARG --global`. Previously under `VERSION --explicit-global 0.6`.
- The Cloud-based secrets model is now project-based; it is not compatible with the older global secrets model. Earthfiles which are defined as `VERSION 0.5` or `VERSION 0.6` will continue to use the old global secrets namespace; however
  the earthly command line no longer supports accessing or modifying the global secrets. A new `earthly secrets migrate` command has been added to help transition the global-based secrets to the new project-based secrets. If you need to manage secrets from Earthly 0.6 without migrating to the new 0.7 secrets, please use an older Earthly binary.
- All `COPY` and `SAVE ARTIFACT` operations now use union filesystem merging for performing the `COPY`. This is similar to `COPY --link` in Dockerfiles, however in Earthly it is automatically enabled for all such operations. Previously under `VERSION --use-copy-link 0.6`.
- The platform logic has been improved to allow overriding the platform in situations where previously it was not possible. Additionally, the default platform is now the native platform of the runner, and not of the host running Earthly. This makes platforms work better in remote runner settings. Previously under `VERSION --new-platform 0.6`.
- Earthly will automatically shellout to determine the `$HOME` value when referenced [#2469](https://github.com/earthly/earthly/issues/2469)
- Improved error message when invalid shell variable name is configured for a secret. [#2478](https://github.com/earthly/earthly/issues/2478)
- The `--ci` flag no longer implies `--save-inline-cache` and `--use-inline-cache` since they were 100% CPU usage in some edge cases. These flags may still be explicitly enabled with `--ci`, but earthly will print a warning.
- `earthly ls` has been promoted from *experimental* to *beta* status.
- Setting a `VERSION` feature flag boolean to false (or any other value) will now raise an error; previously it was syntactically valid but had no effect.
- `SAVE ARTIFACT <path> AS LOCAL ...` when used under a `TRY` / `FINALLY` can fail to be fully transferred to the host when the `TRY` command fails (resulting in an partially transferred file); an underflow can still occur, and is now detected and will not export the partial file. [2452](https://github.com/earthly/earthly/issues/2452)
- The `--keep-own` flag for `SAVE ARTIFACT` is now applied by default; note that `COPY --keep-own` must still be used in order to keep ownership
- Values from the `.env` file will no longer be propagated to Earthfile `ARG`s or `RUN --secret=...` commands; instead values must be placed in `.arg` or `.secret` files respectively. Note that this is a backwards incompatible change and will apply to all Earthfiles (regardless of the defined `VERSION` value). [#1736](https://github.com/earthly/earthly/issues/1736)
- Some particularly obtuse syntax errors now have hints added to help clarify what the expected syntax might be. [#2656](https://github.com/earthly/earthly/issues/2656)
- The default size when launching a new satellite is now medium instead of large.
- Satellites can be launched with a weekend-only mode for receiving auto-updates.

### Added

- The commands `PIPELINE` and `TRIGGER` have been introduced for defining Earthly CI pipelines. Previously under `VERSION --use-pipelines 0.6`.
- The clause `WAIT` is now generally available. The `WAIT` clause allows controlling of build order for operations that require it. This allows use-cases such as pushing images to a registry, followed by infrastructure changes that use the newly pushed images. Previously under `VERSION --wait-block 0.6`.
- The command `CACHE` is now generally available. The `CACHE` command allows declaring a cache mount that can be used by any `RUN` command in the target, and also persists in the final image of the target (contents available when used via `FROM`). Previously under `VERSION --use-cache-command 0.6`.
- The command `HOST` is now generally available. The `HOST` command allows declaring an `/etc/hosts` entry. Previously under `VERSION --use-host-command 0.6`.
- New ARG `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` will contain the author timestamp of the current git commit. [#2462](https://github.com/earthly/earthly/pull/2462)
- New ARGs `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA` contain the version of Earthly and the git sha of Earthly itself, respectively.
- It is now possible to execute shell commands as part of any command that allows using variables. For example `VOLUME $(cat /volume-name.txt)`. Previously under `VERSION --shell-out-anywhere 0.6`.
- Allow custom image to be used for git operations. [#2027](https://github.com/earthly/earthly/issues/2027)
- Earthly now checks for duplicate image names when performing image outputs. Previously under `VERSION --check-duplicate-images 0.6`.
- `SAVE IMAGE --no-manifest-list` allows outputting images of a different platform than the default one, but without the manifest list. This is useful for outputting images for platforms that do not support manifest lists, such as AWS Lambda. Previously under `VERSION --use-no-manifest-list 0.6`.
- `COPY --chmod <mode>` allows setting the permissions of the copied files. Previously under `VERSION --use-chmod 0.6`.
- The new ARG `EARTHLY_LOCALLY` indicates whether the current target is executed in a `LOCALLY` context. Previously under `VERSION --earthly-locally-arg 0.6`.
- The new ARGs `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` contain the author and co-authors of the current git commit, respectively. Previously under `VERSION --earthly-git-author-args 0.6`.
- `earthly doc [projectRef[+targetRef]]` is a new subcommand in *beta* status.  It will parse and output documentation comments on targets.
- Ability to store docker registry credentials in cloud secrets and corresponding `earthly registry setup|list|remove` commands; credentials can be associated with either your user or project.
- New satellite commands for enabling auto-upgrades and forcing a manual upgrade.

### Fixed

- Support for saving files larger than 64kB on failure within a `TRY/FINALLY` block. [#2452](https://github.com/earthly/earthly/issues/2452)
- Fixed race condition where `SAVE IMAGE` or `SAVE ARTIFACT AS LOCAL` commands were not always performed when contained in a target that was referenced by both a `FROM` (or `COPY`) and a `BUILD` command within the context of a `WAIT`/`END` block. [#2237](https://github.com/earthly/earthly/issues/2218)
- `WORKDIR` is lost when `--use-copy-link` feature is enabled with `GIT CLONE` or `COPY --keep-own` commands. Note that `--use-copy-link` is enabled by default in `VERSION 0.7`. [#2544](https://github.com/earthly/earthly/issues/2544)
- The `CACHE` command did not work when used inside a `WITH DOCKER` block. [#2549](https://github.com/earthly/earthly/issues/2549)
- The `--platform` argument is no longer passed to docker or podman, which caused podman to always pull the buildkit image even when it already existed locally. [#2511](https://github.com/earthly/earthly/issues/2511), [#2566](https://github.com/earthly/earthly/issues/2566)
- Fixed missing inline cache export; note that inline cache exports **do not** work when used within a `WAIT` / `END` block, this is a known current limitation. [#2178](https://github.com/earthly/earthly/issues/2178)
- Indentation in the base Earthfile target would cause a panic (when no other targets existed); now a syntax error is returned. [#2603](https://github.com/earthly/earthly/issues/2603)
- Added tighter registry read timeout, to prevent 15min stuck "ongoing" image manifest fetching.

## v0.7.0-rc3 - 2023-02-15

The documentation for this version is available at the [Earthly 0.7 documentation page](https://docs.earthly.dev/v/earthly-0.7/).

**Earthly CI**

Earthly 0.7 is the first version compatible with Earthly CI.

Earthly 0.7 introduces the new keywords `PIPELINE` and `TRIGGER` to help define Earthly CI pipelines.

```
my-pipeline:
    PIPELINE --push
    TRIGGER push main
    TRIGGER pr main
    BUILD +my-target
```

For more information on how to use `PIPELINE` and `TRIGGER`, please see the [reference documentation](https://docs.earthly.dev/v/earthly-0.7/docs/earthfile#pipeline-beta).

**Podman support**

Podman support has now been promoted out of *beta* status and is generally available in 0.7. Earthly will automatically detect the container frontend, whether that's `docker` or `podman` and use it automatically for running Buildkit locally, or for outputting images locally resulting from the build.

Please note that rootful podman is required. Rootless podman is not supported.

**VERSION is now mandatory**

The `VERSION` command
is now required for all Earthfiles, and an error will occur if it is missing. If you are not ready to update your
Earthfiles to use 0.7 (or 0.6), you can declare `VERSION 0.5` to continue to use your Earthfiles.

**.env file is no longer used for `ARG` or secrets**

The `.env` file will only be used to automatically export environment variables, which can be used to configure earthly command line flags.
As a result, values will no longer be propagated to Earthfile `ARG`s or `RUN --secret=...` commands.

Instead if you want build arguments or secrets automatically passed into earthly, they must be placed in `.arg` or `.secret` files respectively.

Note that this is a **backwards incompatible** change and will apply to all Earthfiles (regardless of the defined `VERSION` value).

**Pushing no longer requires everything else to succeed**

The behavior of the `--push` mode has changed in `VERSION 0.7` and is backwards incompatible with `VERSION 0.6`. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, to allow for more flexible push ordering via the new `WAIT` clause. To achieve the behavior of the previous `--push` mode, you need to wrap any pre-required commands in a `WAIT` clause. For example, to push an image only if tests have passed, you would do the following:

```Earthfile
test-and-push:
  WAIT
    BUILD +test
  END
  BUILD +my-image
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

This type of behavior is useful in order to have better control over the order of push operations. For example, you may want to push an image to a registry, followed by a deployment that uses the newly pushed image. Here is how this might look like:

```Earthfile
push-and-deploy:
  ...
  WAIT
    BUILD +my-image
  END
  RUN --push ./deploy.sh my-org/my-image:latest
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

Where `./deploy.sh` is custom deployment script instructing a production environment to start using the image that was just pushed.

**Promoting experimental features**

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.7` at the top of your Earthfile.

Declaring `VERSION 0.7` is equivalent to

```
VERSION \
  --check-duplicate-images \
  --earthly-git-author-args \
  --earthly-locally-arg \
  --earthly-version-arg \
  --explicit-global \
  --new-platform \
  --no-tar-build-output \
  --save-artifact-keep-own \
  --shell-out-anywhere \
  --use-cache-command \
  --use-chmod \
  --use-copy-link \
  --use-host-command \
  --use-no-manifest-list \
  --use-pipelines \
  --use-project-secrets \
  --wait-block \
  0.6
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Changed

- The behavior of the `--push` mode has changed in a backwards incompatible manner. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, allowing push commands to execute in the middle of the build now. Previously under `VERSION --wait-block 0.6`.
- `ARG`s declared in the base target do not automatically become global unless explicitly declared as such via `ARG --global`. Previously under `VERSION --explicit-global 0.6`.
- The Cloud-based secrets model is now project-based; it is not compatible with the older global secrets model. Earthfiles which are defined as `VERSION 0.5` or `VERSION 0.6` will continue to use the old global secrets namespace; however
  the earthly command line no longer supports accessing or modifying the global secrets. A new `earthly secrets migrate` command has been added to help transition the global-based secrets to the new project-based secrets. If you need to manage secrets from Earthly 0.6 without migrating to the new 0.7 secrets, please use an older Earthly binary.
- All `COPY` and `SAVE ARTIFACT` operations now use union filesystem merging for performing the `COPY`. This is similar to `COPY --link` in Dockerfiles, however in Earthly it is automatically enabled for all such operations. Previously under `VERSION --use-copy-link 0.6`.
- The platform logic has been improved to allow overriding the platform in situations where previously it was not possible. Additionally, the default platform is now the native platform of the runner, and not of the host running Earthly. This makes platforms work better in remote runner settings. Previously under `VERSION --new-platform 0.6`.
- Earthly will automatically shellout to determine the `$HOME` value when referenced [#2469](https://github.com/earthly/earthly/issues/2469)
- Improved error message when invalid shell variable name is configured for a secret. [#2478](https://github.com/earthly/earthly/issues/2478)
- The `--ci` flag no longer implies `--save-inline-cache` and `--use-inline-cache` since they were 100% CPU usage in some edge cases. These flags may still be explicitly enabled with `--ci`, but earthly will print a warning.
- `earthly ls` has been promoted from *experimental* to *beta* status.
- Setting a `VERSION` feature flag boolean to false (or any other value) will now raise an error; previously it was syntactically valid but had no effect.
- `SAVE ARTIFACT <path> AS LOCAL ...` when used under a `TRY` / `FINALLY` can fail to be fully transferred to the host when the `TRY` command fails (resulting in an partially transferred file); an underflow can still occur, and is now detected and will not export the partial file. [2452](https://github.com/earthly/earthly/issues/2452)
- The `--keep-own` flag for `SAVE ARTIFACT` is now applied by default; note that `COPY --keep-own` must still be used in order to keep ownership
- Values from the `.env` file will no longer be propagated to Earthfile `ARG`s or `RUN --secret=...` commands; instead values must be placed in `.arg` or `.secret` files respectively. Note that this is a backwards incompatible change and will apply to all Earthfiles (regardless of the defined `VERSION` value). [#1736](https://github.com/earthly/earthly/issues/1736)
- Some particularly obtuse syntax errors now have hints added to help clarify what the expected syntax might be. [#2656](https://github.com/earthly/earthly/issues/2656)


### Added

- The commands `PIPELINE` and `TRIGGER` have been introduced for defining Earthly CI pipelines. Previously under `VERSION --use-pipelines 0.6`.
- The clause `WAIT` is now generally available. The `WAIT` clause allows controlling of build order for operations that require it. This allows use-cases such as pushing images to a registry, followed by infrastructure changes that use the newly pushed images. Previously under `VERSION --wait-block 0.6`.
- The command `CACHE` is now generally available. The `CACHE` command allows declaring a cache mount that can be used by any `RUN` command in the target, and also persists in the final image of the target (contents available when used via `FROM`). Previously under `VERSION --use-cache-command 0.6`.
- The command `HOST` is now generally available. The `HOST` command allows declaring an `/etc/hosts` entry. Previously under `VERSION --use-host-command 0.6`.
- New ARG `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` will contain the author timestamp of the current git commit. [#2462](https://github.com/earthly/earthly/pull/2462)
- New ARGs `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA` contain the version of Earthly and the git sha of Earthly itself, respectively.
- It is now possible to execute shell commands as part of any command that allows using variables. For example `VOLUME $(cat /volume-name.txt)`. Previously under `VERSION --shell-out-anywhere 0.6`.
- Allow custom image to be used for git operations. [#2027](https://github.com/earthly/earthly/issues/2027)
- Earthly now checks for duplicate image names when performing image outputs. Previously under `VERSION --check-duplicate-images 0.6`.
- `SAVE IMAGE --no-manifest-list` allows outputting images of a different platform than the default one, but without the manifest list. This is useful for outputting images for platforms that do not support manifest lists, such as AWS Lambda. Previously under `VERSION --use-no-manifest-list 0.6`.
- `COPY --chmod <mode>` allows setting the permissions of the copied files. Previously under `VERSION --use-chmod 0.6`.
- The new ARG `EARTHLY_LOCALLY` indicates whether the current target is executed in a `LOCALLY` context. Previously under `VERSION --earthly-locally-arg 0.6`.
- The new ARGs `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` contain the author and co-authors of the current git commit, respectively. Previously under `VERSION --earthly-git-author-args 0.6`.
- `earthly doc [projectRef[+targetRef]]` is a new subcommand in *beta* status.  It will parse and output documentation comments on targets.
- Ability to store docker registry credentials in cloud secrets and corresponding `earthly registry setup|list|remove` commands; credentials can be associated with either your user or project.
- New satellite commands for enabling auto-upgrades and forcing a manual upgrade.

### Fixed

- Support for saving files larger than 64kB on failure within a `TRY/FINALLY` block. [#2452](https://github.com/earthly/earthly/issues/2452)
- Fixed race condition where `SAVE IMAGE` or `SAVE ARTIFACT AS LOCAL` commands were not always performed when contained in a target that was referenced by both a `FROM` (or `COPY`) and a `BUILD` command within the context of a `WAIT`/`END` block. [#2237](https://github.com/earthly/earthly/issues/2218)
- `WORKDIR` is lost when `--use-copy-link` feature is enabled with `GIT CLONE` or `COPY --keep-own` commands. Note that `--use-copy-link` is enabled by default in `VERSION 0.7`. [#2544](https://github.com/earthly/earthly/issues/2544)
- The `CACHE` command did not work when used inside a `WITH DOCKER` block. [#2549](https://github.com/earthly/earthly/issues/2549)
- The `--platform` argument is no longer passed to docker or podman, which caused podman to always pull the buildkit image even when it already existed locally. [#2511](https://github.com/earthly/earthly/issues/2511), [#2566](https://github.com/earthly/earthly/issues/2566)
- Fixed missing inline cache export; note that inline cache exports **do not** work when used within a `WAIT` / `END` block, this is a known current limitation. [#2178](https://github.com/earthly/earthly/issues/2178)
- Indentation in the base Earthfile target would cause a panic (when no other targets existed); now a syntax error is returned. [#2603](https://github.com/earthly/earthly/issues/2603)

## v0.7.0-rc2 - 2023-02-01

The documentation for this version is available at the [Earthly 0.7 documentation page](https://docs.earthly.dev/v/earthly-0.7/).

**Earthly CI**

Earthly 0.7 is the first version compatible with Earthly CI.

Earthly 0.7 introduces the new keywords `PIPELINE` and `TRIGGER` to help define Earthly CI pipelines.

```
my-pipeline:
    PIPELINE --push
    TRIGGER push main
    TRIGGER pr main
    BUILD +my-target
```

For more information on how to use `PIPELINE` and `TRIGGER`, please see the [reference documentation](https://docs.earthly.dev/v/earthly-0.7/docs/earthfile#pipeline-beta).

**Podman support**

Podman support has now been promoted out of *beta* status and is generally available in 0.7. Earthly will automatically detect the container frontend, whether that's `docker` or `podman` and use it automatically for running Buildkit locally, or for outputting images locally resulting from the build.

Please note that rootful podman is required. Rootless podman is not supported.

**VERSION is now mandatory**

The `VERSION` command
is now required for all Earthfiles, and an error will occur if it is missing. If you are not ready to update your
Earthfiles to use 0.7 (or 0.6), you can declare `VERSION 0.5` to continue to use your Earthfiles.

**Pushing no longer requires everything else to succeed**

The behavior of the `--push` mode has changed in `VERSION 0.7` and is backwards incompatible with `VERSION 0.6`. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, to allow for more flexible push ordering via the new `WAIT` clause. To achieve the behavior of the previous `--push` mode, you need to wrap any pre-required commands in a `WAIT` clause. For example, to push an image only if tests have passed, you would do the following:

```Earthfile
test-and-push:
  WAIT
    BUILD +test
  END
  BUILD +my-image
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

This type of behavior is useful in order to have better control over the order of push operations. For example, you may want to push an image to a registry, followed by a deployment that uses the newly pushed image. Here is how this might look like:

```Earthfile
push-and-deploy:
  ...
  WAIT
    BUILD +my-image
  END
  RUN --push ./deploy.sh my-org/my-image:latest
my-image:
  ...
  SAVE IMAGE --push my-org/my-image:latest
```

Where `./deploy.sh` is custom deployment script instructing a production environment to start using the image that was just pushed.

**Promoting experimental features**

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.7` at the top of your Earthfile.

Declaring `VERSION 0.7` is equivalent to

```
VERSION \
  --check-duplicate-images \
  --earthly-git-author-args \
  --earthly-locally-arg \
  --earthly-version-arg \
  --explicit-global \
  --new-platform \
  --no-tar-build-output \
  --save-artifact-keep-own \
  --shell-out-anywhere \
  --use-cache-command \
  --use-chmod \
  --use-copy-link \
  --use-host-command \
  --use-no-manifest-list \
  --use-pipelines \
  --use-project-secrets \
  --wait-block \
  0.6
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Changed

- The behavior of the `--push` mode has changed in a backwards incompatible manner. Previously, `--push` commands would only execute if all other commands had succeeded. This precondition is no longer enforced, allowing push commands to execute in the middle of the build now. Previously under `VERSION --wait-block 0.6`.
- `ARG`s declared in the base target do not automatically become global unless explicitly declared as such via `ARG --global`. Previously under `VERSION --explicit-global 0.6`.
- The Cloud-based secrets model is now project-based; it is not compatible with the older global secrets model. Earthfiles which are defined as `VERSION 0.5` or `VERSION 0.6` will continue to use the old global secrets namespace; however
  the earthly command line no longer supports accessing or modifying the global secrets. A new `earthly secrets migrate` command has been added to help transition the global-based secrets to the new project-based secrets. If you need to manage secrets from Earthly 0.6 without migrating to the new 0.7 secrets, please use an older Earthly binary.
- All `COPY` and `SAVE ARTIFACT` operations now use union filesystem merging for performing the `COPY`. This is similar to `COPY --link` in Dockerfiles, however in Earthly it is automatically enabled for all such operations. Previously under `VERSION --use-copy-link 0.6`.
- The platform logic has been improved to allow overriding the platform in situations where previously it was not possible. Additionally, the default platform is now the native platform of the runner, and not of the host running Earthly. This makes platforms work better in remote runner settings. Previously under `VERSION --new-platform 0.6`.
- Earthly will automatically shellout to determine the `$HOME` value when referenced [#2469](https://github.com/earthly/earthly/issues/2469)
- Improved error message when invalid shell variable name is configured for a secret. [#2478](https://github.com/earthly/earthly/issues/2478)
- The `--ci` flag no longer implies `--save-inline-cache` and `--use-inline-cache` since they were 100% CPU usage in some edge cases. These flags may still be explicitly enabled with `--ci`, but earthly will print a warning.
- `earthly ls` has been promoted from *experimental* to *beta* status.
- Setting a `VERSION` feature flag boolean to false (or any other value) will now raise an error; previously it was syntactically valid but had no effect.
- `SAVE ARTIFACT <path> AS LOCAL ...` when used under a `TRY` / `FINALLY` can fail to be fully transferred to the host when the `TRY` command fails (resulting in an partially transferred file); an underflow can still occur, and is now detected and will not export the partial file. [2452](https://github.com/earthly/earthly/issues/2452)
- The `--keep-own` flag for `SAVE ARTIFACT` is now applied by default; note that `COPY --keep-own` must still be used in order to keep ownership

### Added

- The commands `PIPELINE` and `TRIGGER` have been introduced for defining Earthly CI pipelines. Previously under `VERSION --use-pipelines 0.6`.
- The clause `WAIT` is now generally available. The `WAIT` clause allows controlling of build order for operations that require it. This allows use-cases such as pushing images to a registry, followed by infrastructure changes that use the newly pushed images. Previously under `VERSION --wait-block 0.6`.
- The command `CACHE` is now generally available. The `CACHE` command allows declaring a cache mount that can be used by any `RUN` command in the target, and also persists in the final image of the target (contents available when used via `FROM`). Previously under `VERSION --use-cache-command 0.6`.
- The command `HOST` is now generally available. The `HOST` command allows declaring an `/etc/hosts` entry. Previously under `VERSION --use-host-command 0.6`.
- New ARG `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` will contain the author timestamp of the current git commit. [#2462](https://github.com/earthly/earthly/pull/2462)
- New ARGs `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA` contain the version of Earthly and the git sha of Earthly itself, respectively.
- It is now possible to execute shell commands as part of any command that allows using variables. For example `VOLUME $(cat /volume-name.txt)`. Previously under `VERSION --shell-out-anywhere 0.6`.
- Allow custom image to be used for git operations. [#2027](https://github.com/earthly/earthly/issues/2027)
- Earthly now checks for duplicate image names when performing image outputs. Previously under `VERSION --check-duplicate-images 0.6`.
- `SAVE IMAGE --no-manifest-list` allows outputting images of a different platform than the default one, but without the manifest list. This is useful for outputting images for platforms that do not support manifest lists, such as AWS Lambda. Previously under `VERSION --use-no-manifest-list 0.6`.
- `COPY --chmod <mode>` allows setting the permissions of the copied files. Previously under `VERSION --use-chmod 0.6`.
- The new ARG `EARTHLY_LOCALLY` indicates whether the current target is executed in a `LOCALLY` context. Previously under `VERSION --earthly-locally-arg 0.6`.
- The new ARGs `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` contain the author and co-authors of the current git commit, respectively. Previously under `VERSION --earthly-git-author-args 0.6`.
- `earthly doc [projectRef[+targetRef]]` is a new subcommand in *beta* status.  It will parse and output documentation comments on targets.
- Ability to store docker registry credentials in cloud secrets and corresponding `earthly registry login|list|logout` commands; credentials can be associated with either your user or project.
- New satellite commands for enabling auto-upgrades and forcing a manual upgrade.

### Fixed

- Support for saving files larger than 64kB on failure within a `TRY/FINALLY` block. [#2452](https://github.com/earthly/earthly/issues/2452)
- Fixed race condition where `SAVE IMAGE` or `SAVE ARTIFACT AS LOCAL` commands were not always performed when contained in a target that was referenced by both a `FROM` (or `COPY`) and a `BUILD` command within the context of a `WAIT`/`END` block. [#2237](https://github.com/earthly/earthly/issues/2218)
- `WORKDIR` is lost when `--use-copy-link` feature is enabled with `GIT CLONE` or `COPY --keep-own` commands. Note that `--use-copy-link` is enabled by default in `VERSION 0.7`. [#2544](https://github.com/earthly/earthly/issues/2544)
- The `CACHE` command did not work when used inside a `WITH DOCKER` block. [#2549](https://github.com/earthly/earthly/issues/2549)
- The `--platform` argument is no longer passed to docker or podman, which caused podman to always pull the buildkit image even when it already existed locally. [#2511](https://github.com/earthly/earthly/issues/2511), [#2566](https://github.com/earthly/earthly/issues/2566)
- Fixed missing inline cache export; note that inline cache exports **do not** work when used within a `WAIT` / `END` block, this is a known current limitation. [#2178](https://github.com/earthly/earthly/issues/2178)
- Indentation in the base Earthfile target would cause a panic (when no other targets existed); now a syntax error is returned. [#2603](https://github.com/earthly/earthly/issues/2603)


## v0.7.0-rc1 - 2023-01-18

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of
the features in this version you need to declare `VERSION 0.7` at the top of your Earthfile. The `VERSION` command
is now required for all Earthfiles, and an error will occur if it is missing. If you are not ready to update your
Earthfiles to use 0.7 (or 0.6), you can declare `VERSION 0.5` to continue to use your Earthfiles.

Declaring `VERSION 0.7` is equivalent to

```
VERSION \
  --explicit-global \
  --check-duplicate-images \
  --earthly-version-arg \
  --use-cache-command \
  --use-host-command \
  --use-copy-link \
  --new-platform \
  --no-tar-build-output \
  --use-no-manifest-list \
  --use-chmod \
  --shell-out-anywhere \
  --earthly-locally-arg \
  --use-project-secrets \
  --use-pipelines \
  --earthly-git-author-args \
  0.6
```

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Changed

- The Cloud-based secrets model is now project-based; it is not compatible with the older global secrets model. Earthfiles which are defined as `VERSION 0.5` or `VERSION 0.6` will continue to use the old global secrets namespace; however
  the earthly command line no longer supports accessing or modifying the global secrets. A new `earthly secrets migrate` command has been added to help transition the global-based secrets to the new project-based secrets.
- Earthly will automatically shellout to determine the `$HOME` value when referenced; this requires the `--shell-out-anywhere` feature flag. [#2469](https://github.com/earthly/earthly/issues/2469)
- Improved error message when invalid shell variable name is configured for a secret. [#2478](https://github.com/earthly/earthly/issues/2478)
- The `--ci` flag no longer implies `--save-inline-cache` and `--use-inline-cache` since they were 100% CPU usage in some edge cases. These flags may still be explicitly enabled with `--ci`, but earthly will print a warning.

### Added

- New ARG `EARTHLY_GIT_COMMIT_AUTHOR_TIMESTAMP` will contain the author timestamp of the current git commit, this ARG must be enabled with the `VERSION --git-commit-author-timestamp` feature flag. [#2462](https://github.com/earthly/earthly/pull/2462)
- Allow custom image to be used for git opperations. [#2027](https://github.com/earthly/earthly/issues/2027)

### Fixed

- Support for saving files larger than 64kB on failure within a `TRY/FINALLY` block. [#2452](https://github.com/earthly/earthly/issues/2452)
- Fixed race condition where `SAVE IMAGE` or `SAVE ARTIFACT AS LOCAL` commands were not always performed when contained in a target that was referenced by both a `FROM` (or `COPY`) and a `BUILD` command within the context of a `WAIT`/`END` block. [#2237](https://github.com/earthly/earthly/issues/2218)
- `WORKDIR` is lost when `--use-copy-link` feature is enabled with `GIT CLONE` or `COPY --keep-own` commands. [#2544](https://github.com/earthly/earthly/issues/2544)
- The `CACHE` command did not work when used inside a `WITH DOCKER` block. [#2549](https://github.com/earthly/earthly/issues/2549)
- The `--platform` argument is no longer passed to docker or podman, which caused podman to always pull the buildkit image even when it already existed locally. [#2511](https://github.com/earthly/earthly/issues/2511), [#2566](https://github.com/earthly/earthly/issues/2566)

## v0.6.30 - 2022-11-22

### Added

- Added support for a custom `.netrc` file path using the standard `NETRC` environmental variable. [#2426](https://github.com/earthly/earthly/pull/2426)
- Ability to run multiple Earthly installations at a time via `EARTHLY_INSTALLATION_NAME` environment variable, or the `--installation-name` CLI flag. The installation name defaults to `earthly` if not specified. Different installations use different configurations, different buildkit Daemons, different cache volumes, and different ports.
- New `EARTHLY_CI` builtin arg, which is set to `true` when earthly is run with the `--ci` flag, this ARG must be enabled with the `VERSION --ci-arg` feature flag. [#2398](https://github.com/earthly/earthly/pull/2398)

### Changed

- Updated buildkit to include changes up to [a5263dd0f990a3fe17b67e0002b76bfd1f5b433d](https://github.com/moby/buildkit/commit/a5263dd0f990a3fe17b67e0002b76bfd1f5b433d), which includes a change to speed-up buildkit startup time.
- The Earthly Docker image works better for cases where a buildkit instance is not needed. The image now works without `--privileged` when using `NO_BUILDKIT=1`, and additionally, the image can also use `/var/run/docker.sock` or `DOCKER_HOST` for the buildkit daemon.

### Fixed

- Fixed Earthly on Mac would randomly hang on `1. Init` if Earthly was installed from Homebrew or the Earthly homebrew tap. [#2247](https://github.com/earthly/earthly/issues/2247)
- Only referenced ARGs from .env are displayed on failures, this prevents secrets contained in .env from being displayed. [#1736](https://github.com/earthly/earthly/issues/1736)
- Earthly now correctly detects if Podman is running but is under the disguise of the Docker CLI.
- Improved performance when copying files. Fully-cached builds are now dramatically faster as a result. [#2049](https://github.com/earthly/earthly/issues/2049)
- Fixed `--shell-out-anywhere` bug where inner quotes were incorrectly removed. [#2340](https://github.com/earthly/earthly/issues/2340)

## v0.6.29 - 2022-11-07

### Added

- Cache mounts sharing mode can now be specified via `RUN --mount type=cache,sharing=shared` via `CACHE --sharing=shared`. Allowed values are `locked` (default - lock concurrent acccess to the cache), `shared` (allow concurrent access) and `private` (create a new empty cache on concurrent access).

### Changed

- Increases the cache limit for local and git sources from 10% to 50% to support copying large files (e.g. binary assets).
- The default cache mount sharing mode is now `locked` instead of `shared`. This means that if you have multiple builds running concurrently, they will block on each other to gain access to the cache mount. If you want to share the cache as it was shared in previous version of Earthly, you can use `RUN --mount type=cache,sharing=shared` or `CACHE --sharing=shared`.

### Fixed

- `CACHE` command was not being correctly used in `IF`, `FOR`, `ARG` and other commands. [#2330](https://github.com/earthly/earthly/issues/2330)
- Fixed buildkit gckeepstorage config value which was was set to 1000 times larger than the cache size, now it is set to the cache size.
- Fixed Earthly not detecting the correct image digest for some images loaded in `WITH DOCKER --load` and causing cache not to be bust correctly. [#2337](https://github.com/earthly/earthly/issues/2337) and [#2288](https://github.com/earthly/earthly/issues/2288)

## v0.6.28 - 2022-10-26

### Added
- A summary of context file transfers is now displayed every 15 seconds.
- Satellite wake command, which can force a satellite to wake up (useful for calling inspect or other non-build related commands).

### Changed
- `WITH DOCKER` merging of user specific `/etc/docker/daemon.json` settings data now applies to arrays (previously only dictionaries were supported).
- A final warning will be displayed if earthly is terminated due to a interrupt signal (ctrl-c).

### Changed
- Updated buildkit to include changes up to [c717d6aa7543d4b83395e0552ef2eb311f563aab](https://github.com/moby/buildkit/commit/c717d6aa7543d4b83395e0552ef2eb311f563aab)

## v0.6.27 - 2022-10-17

### Changed
- Support for all ssh-based key types (e.g. ssh-ed25519), and not only ssh-rsa. [#1783](https://github.com/earthly/earthly/issues/1783)

### Fixed
- Unable to specify public key to add via the command-line, e.g. running `earthly account add-key <key>` ignored the key and fell back to an interactive prompt.
- `GIT CLONE` command was ignoring the `WORK DIR` command when `--use-copy-link` feature was set.

## v0.6.26 - 2022-10-13

### Added

- Build failures now show the file and line number of the failing command
- Introduced `EARTHLY_GIT_AUTHOR` and `EARTHLY_GIT_CO_AUTHORS` ARGS

### Fixed

- Some network operations were being incorrectly executed with a timeout of 0.
- Upon `earthly ls` failure it will display the failure reason

### Changed

- Loading Docker images as part of `WITH DOCKER` is now faster through the use of an embedded registry in Buildkit. This functionality was previously hidden (`VERSION --use-registry-for-with-docker`) and was only auto-enabled for Earthly Satellite users. It is now enabled by default for all builds. [#1268](https://github.com/earthly/earthly/issues/1268)

### Changed

- `VERSION` command is now required.

## v0.6.25 - 2022-10-04

### Fixed

- Fixed outputting images with long names [#2053](https://github.com/earthly/earthly/issues/2053)
- Fixed buildkit connection timing out occasionally [#2229](https://github.com/earthly/earthly/issues/2229)
- Cache size was incorrectly displayed (magnitude of 1024 higher)

## v0.6.24 - 2022-09-22

### Added

- The `earthly org invite` command now has the ability to invite multiple email addresses at once.
- Experimental support for `TRY/FINALLY`, which allows saving artifacts upon failure. [#988](https://github.com/earthly/earthly/issues/988), [#587](https://github.com/earthly/earthly/issues/587).
  Not that this is only a partial implementation, and only accepts a *single* RUN command in the `TRY`, and only `SAVE ARTIFACT` commands in the `FINALLY` block.
- Ability to enable specific satellite features via cli flags, e.g. the new experimental sleep feature can be enabled with
  `earthly satellite launch --feature-flags satellite-sleep my-satellite`.

### Changed

- Bootstrapping zsh autocompletion will first attempt to install under `/usr/local/share/zsh/site-functions`, and will now
  fallback to `/usr/share/zsh/site-functions`.
- The `earthly preview org` command has been promoted to GA, and is now available under `earthly org`.
- `earthly sat select` with no arguments now prints the current satellite and the usage text.
- The interactive debugger now connects over the buildkit session connection rather than an unencrypted tcp connection; this makes it possible
  to use the interactive debugger with remote buildkit instances.

### Fixed

- Fixed Earthly failing when using a remote docker host from a machine with an incompatible architecture. [#1895](https://github.com/earthly/earthly/issues/1895)
- Earthly will no longer race with itself when starting up buildkit. [#2194](https://github.com/earthly/earthly/issues/2194)
- The error reported when failing to initiate a connection to buildkit has been reworded to account for the remote buildkit/satellite case too.
- Errors related to parsing `VERSION` feature flags will no longer be displayed during auto-completion.

## v0.6.23 - 2022-09-06

### Fixed

- Using `--remote-cache` on a target that contains only `BUILD` instructions caused a hang. [#1945](https://github.com/earthly/earthly/issues/1945)
- Fixed WAIT/END related bug which prevent `WITH DOCKER --load` from building referenced target.
- Images and artifacts which are output (or pushed), are now displayed in the final earthly output.
- `ssh: parse error in message type 27` error when using OpenSSH 8.9; fixed by upstream in [golang/go#51689](https://github.com/golang/go/issues/51689).

### Changed

- Removed warning stating that `WAIT/END code is experimental and may be incomplete` -- it is still experimental; however, it now has a higher degree
  of test-coverage. It can be enabled with `VERSION --wait-block 0.6`.
- A warning is now displayed during exporting a multi-platform image to the local host if no platform is found that matches the host's platform type.
- Reduced verbosity of `To enable pushing use earthly --push` message.

## v0.6.22 - 2022-08-19

### Added

- `--cache-from` earthly flag, which allows defining multiple ordered caches. [#1693](https://github.com/earthly/earthly/issues/1693)
- WAIT/END support for saving artifacts to local host.
- WAIT/END support for `RUN --push` commands.

### Fixed

- Updated `EXPOSE` parsing to accept (and ignore) host IP prefix, as well as expose udp ports; this should be fully-compatible with dockerfile's format. [#1986](https://github.com/earthly/earthly/issues/1986)
- The earthly-buildkit container is now only initialized when required.

### Changed

- The earthly-buildkit container is now only initialized when required.

## v0.6.21 - 2022-08-04

### Added

- `EARTHLY_LOCALLY` builtin arg which is set to `true` or `false` when executing locally or within a container, respectively. This ARG must be enabled with
  the `VERSION --earthly-locally-arg` feature flag.

### Fixed

- Fixed an incompatibility with older versions of remote BuildKits and Satellites, which was resulting in Earthly crashing.
- Fixed `WITH DOCKER` not loading correctly when the image name contained a port number under `VERSION --use-registry-for-with-docker`. [#2071](https://github.com/earthly/earthly/issues/2071)
- Race condition in WAIT / END block, which prevented waiting on some BUILD commands.

### Changed

- Added a deprecation warning for secrets using a `+secrets/` prefix. Support for this prefix will be removed in a future release.
- per-file stat transfers are now logged when running under `--debug` mode.

## v0.6.20 - 2022-07-18

### Changed

- Updated buildkit to include changes up to 12cfc87450c8d4fc31c8c0a09981e4c3fb3e4d9f

### Added

- Adding support for saving artifact from `--interactive-keep`. [#1980](https://github.com/earthly/earthly/issues/1980)
- New `EARTHLY_PUSH` builtin arg, which is set to `true` when earthly is run with the `--push` flag, and the argument
  is referenced under the direct target, or a target which is indirectly referenced via a `BUILD` command; otherwise
  it will be set to `false`. The value mimics when a `RUN --push` command is executed. This feature must be enabled with
  `VERSION --wait-block 0.6`.

### Fixed

- Fixed `context.Canceled` being reported as the error in some builds instead of the root cause. [#1991](https://github.com/earthly/earthly/issues/1991)
- Improved cache use of `WITH DOCKER` command.
- The `earthly/earthly` docker image is now also built for arm64 (in addition to amd64).

## v0.6.19 - 2022-06-29

### Fixed

- Fixed retagging of images that are made available via the `WITH DOCKER` command when the `--use-registry-for-with-docker` feature is enabled.
- Fixed a bug where `earthly --version` would display unknown on some versions of Windows.

## v0.6.18 - 2022-06-27

### Fixed

- `sh: write error: Resource busy` error caused by running the earthly/earthly docker image on a cgroups2-enabled host. [#1934](https://github.com/earthly/earthly/issues/1934)

## v0.6.17 - 2022-06-20

### Added

- Additional debug information for failure during dind cleanup.

## v0.6.16 - 2022-06-17

### Changed

- Custom `secret_provider` is now called with user's env variables.
- Additional args can be passed to `secret_provider`, e.g. `secret_provider: my-password-manager --db=$HOME/path/to/secrets.db`
- Local registry is enabled by default in the earthly-buildkit container.

## v0.6.15 - 2022-06-02

### Changed

- Switch to MPL-2.0 license. [Announcement](https://earthly.dev/blog/earthly-open-source)

### Added

- Experimental support for Docker registry based image creation and transfer `WITH DOCKER` loads and pulls. Enable with the `VERSION --use-registry-for-with-docker` flag.
- Git config options for non-standard port and path prefix; these options are incompatible with a custom git substitution regex.
- Experimental WAIT / END blocks, to allow for finer grain of control between pushing images and running commands.
- Improved ARG error messages to include the ARG name associated with the error.

### Fixed

- Panic when running earthly --version under some versions of Windows
- Removed duplicate git commit hash from earthly --version output string (when running dev versions of earthly)
- Garbled auto-completion when using Earthfiles without a VERSION command (or with other warnings) [#1837](https://github.com/earthly/earthly/issues/1837).
- Masking of cgroups for podman support.

## v0.6.14 - 2022-04-11

### Added

- Experimental support for `SAVE IMAGE --no-manifest-list`. This option disables creating a multi-platform manifest list for the image, even if the image is created with a non-default platform. This allows the user to create non-native images (e.g. amd64 image on an M1 laptop) that are still compatible with AWS lambda. To enable this feature, please use `VERSION --use-no-manifest-list 0.6`. [#1802](https://github.com/earthly/earthly/pull/1802)
- Introduced Experimental support for `--chmod` flag in `COPY`. To enable this feature, please use `VERSION --use-chmod 0.6`. [#1817](https://github.com/earthly/earthly/pull/1817)
- Experimental `secret_provider` config option allows users to provide a script which returns secrets. [#1808](https://github.com/earthly/earthly/issues/1808)
- `/etc/ssh/ssh_known_hosts` are now passed to buildkit. [#1769](https://github.com/earthly/earthly/issues/1769)

### Fixed

- Targets with the same `CACHE` commands incorrectly shared cached contents. [#1805](https://github.com/earthly/earthly/issues/1805)
- Sometimes local outputs and pushes are skipped mistakenly when a target is referenced both via `FROM` and via `BUILD` [#1823](https://github.com/earthly/earthly/issues/1823)
- `GIT CLONE` failure (`makeCloneURL does not support gitMatcher substitution`) when used with a self-hosted git repo that was configured under `~/.earthly/config.yml`  [#1757](https://github.com/earthly/earthly/issues/1757)

## v0.6.13 - 2022-03-30

### Added

- Earthly now warns when encountering Earthfiles with no `VERSION` specified. In the future, the `VERSION` command will be mandatory. [#1775](https://github.com/earthly/earthly/pull/1775)

### Changed

- `WITH DOCKER` now merges changes into `/etc/docker/daemon.json` rather than overwriting the entire file; this change introduces `jq` as a dependency, which will
  be auto-installed if missing.

### Fixed

- The `COPY` command, when used with `LOCALLY` was incorrectly ignoring the `WORKDIR` value. [#1792](https://github.com/earthly/earthly/issues/1792)
- The `--shell-out-anywhere` feature introduced a bug which interfered with asynchronous builds. [#1785](https://github.com/earthly/earthly/issues/1785)
- `EARTHLY_GIT_SHORT_HASH` was not set when building a remotely-referenced target. [#1787](https://github.com/earthly/earthly/issues/1787)

## v0.6.12 - 2022-03-23

### Changed

- A more obvious error is printed if `WITH DOCKER` starts non-natively. This is not supported and it wasn't obvious before.
- `WITH DOCKER` will keep any settings pre-applied in `/etc/docker/daemon.json` rather than overwriting them.

### Added

- The feature flag `--exec-after-build` has been enabled retroactively for `VERSION 0.5`. This speeds up large builds by 15-20%.
- The feature flag `--parallel-load` has been enabled for every `VERSION`. This speeds up by parallelizing targets built for loading via `WITH DOCKER --load`.
- `VERSION 0.0` is now permitted, however it is only meant for Earthly internal debugging purposes. `VERSION 0.0` disables all feature flags.
- A new experimental mode in which `--platform` operates. To enable these features in your builds, set `VERSION --new-platform 0.6`:
  - There is now a distinction between **user** platform and **native** platform. The user platform is the platform of the user running Earthly, while the native platform is the platform of the build worker (these can be different when using a remote buildkit)
  - New platform shorthands are provided: `--platform=native`, `--platform=user`.
  - New builtin args are available: `NATIVEPLATFORM`, `NATIVEOS`, `NATIVEARCH`, `NATIVEVARIANT` (these are the equivalent of the `USER*` and `TARGET*` platform args).
  - When no platform is provided, earthly will default to the **native** platform
  - Additionally, earthly now default to native platform for internal operations too (copy operations, git clones etc)
  - Earthly now allows changing the platform in the middle of a target (`FROM --platform` is not a contradiction anymore). There is a distinction between the "input" platform of a target (the platform the caller passes in) vs the "output" platform (the platform that ends up being the final platform of the image). These can be different if the caller passes `BUILD --platform=something +target`, but the target starts with `FROM --platform=otherthing ...`.
- Ability to shell-out in any Earthly command, (e.g. `SAVE IMAGE myimage:$(cat version)`), as well as in the middle of ARG strings. To enable this feature, use `VERSION --shell-out-anywhere 0.6`.

### Fixed

- An experimental fix for duplicate output when building images that are loaded via `WITH DOCKER --load`. This can be enabled via `VERSION --no-tar-build-output 0.6`.

## v0.6.11 - 2022-03-17

### Added

- An experimental feature whereby `WITH DOCKER` parallelizes building of the
  images to be loaded has been added. To enable this feature use
  `VERSION --parallel-load 0.6`. [#1725](https://github.com/earthly/earthly/pull/1725)
- Added `cache_size_pct` config option to allow specifying cache size as a percentage of disk space.

### Fixed

- Fixed a duplicate build issue when using `IF` together with `WITH DOCKER` [#1724](https://github.com/earthly/earthly/issues/1724)
- Fixed a bug where `BUILD --platform=$ARG` did not expand correctly
- Fixed issue preventing use of `WITH DOCKER` with docker systemd-based images such as `kind`, when used under hosts with cgroups v2.

## v0.6.10 - 2022-03-03

### Changed

- reverted zeroing of mtime change that was introduced in v0.6.9; this restores the behavior of setting modification time to `2020-04-16T12:00`. [#1712](https://github.com/earthly/earthly/issues/1712)

## v0.6.9 - 2022-03-02

### Changed

- Log sharing is enabled by default for logged in users, it can be disabled with `earthly config global.disable_log_sharing true`.
- `SAVE ARTIFACT ... AS LOCAL` now sets mtime of output artifacts to the current time.

### Added

- Earthly is now 15-30% faster when executing large builds [#1589](https://github.com/earthly/earthly/issues/1589)
- Experimental `HOST` command, which can be used like this: `HOST <domain> <ip>` to add additional hosts during the execution of your build. To enable this feature, use `VERSION --use-host-command 0.6`. [#1168](https://github.com/earthly/earthly/issues/1168)

### Fixed

- Errors when using inline caching indicating `invalid layer index` [#1635](https://github.com/earthly/earthly/issues/1635)
- Podman can now use credentials from the default location [#1644](https://github.com/earthly/earthly/issues/1644)
- Podman can now use the local registry cache without modifying `registries.conf` [#1675](https://github.com/earthly/earthly/pull/1675)
- Podman can now use `WITH DOCKER --load` inside a target marked as `LOCALLY` [#1675](https://github.com/earthly/earthly/pull/1675)
- Interactive sessions should now work with rootless configurations that have no apparent external IP address [#1573](https://github.com/earthly/earthly/issues/1573), [#1689](https://github.com/earthly/earthly/pull/1689)
- On native Windows installations, Earthly properly detects the local git path when it's available [#1663](https://github.com/earthly/earthly/issues/1663)
- On native Windows installations, Earthly will properly identify targets in Earthfiles outside of the current directory using the `\` file separator  [#1663](https://github.com/earthly/earthly/issues/1663)
- On native Windows installations, Earthly will save local artifacts to directories using the `\` file separator [#1663](https://github.com/earthly/earthly/issues/1663)
- A parsing error, when using `WITH DOCKER --load` in conjunction with new-style
  build args. [#1696](https://github.com/earthly/earthly/issues/1696)
- `ENTRYPOINT` and `CMD` were not properly expanding args when used in shell mode.
- A race condition sometimes caused a `Canceled` error to be reported, instead of the real error that caused the build to fail

## v0.6.8 - 2022-02-16

### Fixed

- `RUN --interactive` command exit codes were being ignored.
- `RUN --ssh` command were failing to create `SSH_AUTH_SOCK` when run inside a `WITH DOCKER`. [#1672](https://github.com/earthly/earthly/issues/1672)

### Changed

- expanded help text for `earthly account register --help`.

## v0.6.7 - 2022-02-09

Log Sharing (experimental)

This version of Earthly includes an experimental log-sharing feature which will
upload build-logs to the cloud when enabled.

To enable this experimental feature, you must first sign up for an earthly account
by using the [`earthly account register`](https://docs.earthly.dev/docs/earthly-command#earthly-account-register)
command, or by visiting [https://ci.earthly.dev/](https://ci.earthly.dev/)

Once logged in, you must explicitly enable log-sharing by running:

    earthly config global.disable_log_sharing false

In a future version, log-sharing will be enabled by default for logged-in users; however, you will still be able to disable it, if needed.

When log-sharing is enabled, you will see a message such as

    Share your build log with this link: https://ci.earthly.dev/logs?logId=dc622821-9fe4-4a13-a1db-12680d73c442

as the last line of `earthly` output.

### Fixed

- `GIT CLONE` now works with annotated git tags. [#1571](https://github.com/earthly/earthly/issues/1571)
- `CACHE` command was not working for versions of earthly installed via homebrew.
- Autocompletion bug when directory has both an Earthfile and subdir containing an earthfile.
- Autocompletion bug when directory has two subdirectories where one is a prefix of the other.

### Changed

- `earthly account logout` raises an error when `EARTHLY_TOKEN` is set.

## v0.6.6 - 2022-01-26

### Added

- Ability to change mounted secret file mode. fixes [#1434](https://github.com/earthly/earthly/issues/1434)

### Changed

- Permission errors related to reading `~/.earthly/config.yml` and `.env` files are now treated as errors rather than silently ignored (and assuming the file does not exist).
- Speedup from pre-emptive execution of build steps prior to them being referenced in the build graph.

### Fixed

- earthly panic when running with `SUDO_USER` pointing to a user the current user did not have read/write permission; notably encountered when running under circleci.

### Removed

- Removed `--git-url-instead-of` flag, which has been replaced by `earthly config git ...`

## v0.6.5 - 2022-01-24

### Added

- Ability to load a different `.env` file via the `--env-file` flag.
- Added experimental feature than changes the ARGs defined in the `+base` target to be local, unless defined with a `--global` flag;
  To enable this feature use `VERSION --explicit-global 0.6`.

### Changed

- Updated buildkit to include changes up to 17c237d69a46d61653746c03bcbe6953014b41a5

### Fixed

- `failed to solve: image  is defined multiple times for the same default platform` errors. [#1594](https://github.com/earthly/earthly/issues/1594), [#1582](https://github.com/earthly/earthly/issues/1582)
- `failed to solve: image rmi after pull and retag: command failed: docker image rm ...: exit status 1: Error: No such image` errors. [#1590](https://github.com/earthly/earthly/issues/1590)

## v0.6.4 - 2022-01-17

### Fixed

- Duplicate execution occurring when using ARGs. [#1572](https://github.com/earthly/earthly/issues/1572), [#1582](https://github.com/earthly/earthly/issues/1582)
- Overriding builtin ARG value now displays an error (rather than silently ignoring it).

## v0.6.3 - 2022-01-12

### Changed

- Updated buildkit to contain changes up to `15fb1145afa48bf81fbce41634bdd36c02454f99` from `moby/master`.

### Added

- Experimental `CACHE` command can be used in Earthfiles to optimize the cache in projects that perform better with incremental changes. For example, a Maven
  project where `SNAPSHOT` dependencies are added frequently, an NPM project where `node_modules` change frequently, or programming languages using
  incremental compilers. [#1399](https://github.com/earthly/earthly/issues/1399)
- Config file entries can be deleted using a `--delete` flag (for example `earthly config global.conversion_parallelism --delete`). [#1449](https://github.com/earthly/earthly/issues/1449)
- Earthly now provides the following [builtin ARGs](https://docs.earthly.dev/docs/earthfile/builtin-args): `EARTHLY_VERSION` and `EARTHLY_BUILD_SHA`. These
  will be generally available in Earthly version 0.7+, however, they can be enabled earlier by using the `--earthly-version-arg`. [feature flag](https://docs.earthly.dev/docs/earthfile/features#feature-flags) [#1452](https://github.com/earthly/earthly/issues/1452)
- Config option to disable `known_host` checking for specific git hosts by setting `strict_host_key_checking ` to `false` under the `git` section of `earthly/config.yml` (defaults to `true`).
- Error check for using both `--interactive` and `--buildkit-host` (which are not currently supported together). [#1492](https://github.com/earthly/earthly/issues/1492)
- `earthly ls [<project-ref>]` to list Earthfile targets.

### Fixed

- Gracefully handle empty string `""` being provided as a value to `earthly config` commands. [#1449](https://github.com/earthly/earthly/issues/1449)
- `known_host` entries were being ignored when custom `pattern` and `substituted` git config options were used (commonly used for [self-hosted git repos](https://docs.earthly.dev/docs/guides/auth#self-hosted-and-private-git-repositories))
- Unable to connect to ssh server when `known_hosts` doesn't contain ssh-rsa host scan, but contains a different key-scan (e.g. `ecdsa-sha2-nistp256`, `ssh-ed25519`, etc).
- When git auth is set to ssh but no user is given, default to current user (similar to calling `ssh example.com` vs `ssh user@example.com`).

## v0.6.2 - 2021-12-01

### Fixed

- `unexpected non-relative path within git dir` bug when using case insensitive file systems [#1426](https://github.com/earthly/earthly/issues/1426)
- Unable to access private GitHub repos [#1421](https://github.com/earthly/earthly/issues/1421)

## v0.6.1 - 2021-11-29

### Fixed

- `BUILD` arguments containing a subshell (`$(...)`) were executed twice, and when `+base` target was empty would result errors such as `the first command has to be FROM, FROM DOCKERFILE, LOCALLY, ARG, BUILD or IMPORT` [#1448](https://github.com/earthly/earthly/issues/1448)
- TLS error (`transport: authentication handshake failed: remote error: tls: no application protocol`) when enabling buildkit mTLS  [#1439](https://github.com/earthly/earthly/issues/1439)
- Unable to save artifacts to local directory (`.`) [#1422](https://github.com/earthly/earthly/issues/1422)

## v0.6.0 - 2021-11-24

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of the features in this version you need to declare `VERSION 0.6` at the top of your Earthfile. If a version is not declared, then Earthly's interpreter will assume `VERSION 0.5`.

If you are not ready to update your scripts to take advantage of `VERSION 0.6`, then you may upgrade Earthly anyway and your scripts should continue to work as before, provided that they either declare `VERSION 0.5` or they don't declare a version at all.

Declaring `VERSION 0.6` is equivalent to

```
VERSION \
  --use-copy-include-patterns \
  --referenced-save-only \
  --for-in \
  --require-force-for-unsafe-saves \
  --no-implicit-ignore \
  0.5
```

It is recommended to use `VERSION 0.6` instead as individual feature flags don't guarantee proper forwards-backwards compatibility. Note, however, that Earthly `0.5.*` is not able to run a `VERSION 0.6` Earthfile and will return an error.

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Changed

<!--changelog-parser-ignore-start-->
- What Earthly outputs locally has changed in a way that is not backwards compatible. For an artifact or an image to be produced locally it needs to be part of a `BUILD` chain (or be part of the target being directly built). Artifacts and images introduced through `FROM` or `COPY` are no longer output locally.

  To update existing scripts, you may issue a duplicate `BUILD` in addition to a `FROM` (or a `COPY`), should you wish for the referenced target to perform output.

  For example, the following script

  ```
  FROM +some-target
  COPY +another-target/my-artifact ./
  ```

  could become

  ```
  FROM +some-target
  BUILD +some-target
  COPY +another-target/my-artifact ./
  BUILD +another-target
  ```

  in order to produce the same outputs.

  For more details see [#896](https://github.com/earthly/earthly/issues/896).
- The syntax for passing build args has been changed.

  Earthly v0.5 (old way)

  ```
  FROM --build-arg NAME=john +some-target
  COPY --build-arg NAME=john +something/my-artifact ./
  WITH DOCKER --build-arg NAME=john --load +another-target
    ...
  END
  ```

  Earthly v0.6 (new way)

  ```
  FROM +some-target --NAME=john
  COPY (+something/my-artifact --NAME=john) ./
  WITH DOCKER --load (+another-target --NAME=john)
    ...
  END
  ```

  Passing build args on the command-line has also changed similarly:

  Earthly v0.5 (old way)

  ```
  earthly --build-arg NAME=john +some-target
  ```

  Earthly v0.6 (new way)

  ```
  earthly +some-target --NAME=john
  ```

  This change is part of the [UDC proposal #581](https://github.com/earthly/earthly/issues/581). The old way of passing args is deprecated and will be removed in a future version (however, it still works in 0.6).
<!--changelog-parser-ignore-end-->
- If a `SAVE ARTIFACT` is unsafe (writing to a directory outside of the Earthfile directory), it'll require the `--force` flag.
- `.earthlyignore` no longer includes any implicit entries like `Earthfile` or `.earthlyignore`. These will need to be specified explicitly. [#1294](https://github.com/earthly/earthly/issues/1294)
- Buildkit was updated to `d429b0b32606b5ea52e6be4a99b69d67b7c722b2`. This includes a number of bug fixes, including eliminating crashes due to `panic failed to get edge`.

### Added

- Earthly now performs local image outputs to the local Docker daemon through a built-in registry. This speeds up the process drastically as common layers no longer need to be transferred over [#500](https://github.com/earthly/earthly/issues/500).
- Earthly now enables additional parallelism to speed up certain operations that were previously serialized [#888](https://github.com/earthly/earthly/issues/888). Note that this setting was previously controlled by `--conversion-parallelism` flag or the `EARTHLY_CONVERSION_PARALLELISM` environment variable while in experimental stage. It has now been moved as part of the Earthly config and has been promoted to GA.
- `COPY` transfers are sped up as only the necessary files are sent over to BuildKit [#1062](https://github.com/earthly/earthly/issues/1062).
- [`WITH DOCKER`](https://docs.earthly.dev/docs/earthfile#with-docker) has been promoted to GA [#576](https://github.com/earthly/earthly/issues/576).
- [`FROM DOCKERFILE`](https://docs.earthly.dev/docs/earthfile#from-dockerfile) has been promoted to GA.
- [`LOCALLY`](https://docs.earthly.dev/docs/earthfile#locally) has been promoted to GA [#580](https://github.com/earthly/earthly/issues/580).
- [`RUN --interactive` and `RUN --interactive-keep`](https://docs.earthly.dev/docs/earthfile#run) have been promoted to GA [#693](https://github.com/earthly/earthly/issues/693).
- [`IF`](https://docs.earthly.dev/docs/earthfile#if) and [`FOR`](https://docs.earthly.dev/docs/earthfile#for) have been promoted to GA [#779](https://github.com/earthly/earthly/issues/779).
- Support for Apple Silicon M1 has been promoted to GA [#722](https://github.com/earthly/earthly/issues/722).
- [Multi-platform builds](https://docs.earthly.dev/docs/guides/multi-platform) have been promoted to GA [#536](https://github.com/earthly/earthly/issues/536).
- Mounting secrets as files have been promoted as GA [#579](https://github.com/earthly/earthly/issues/579).
- [`VERSION`](https://docs.earthly.dev/docs/earthfile#version) has been promoted to GA [#991](https://github.com/earthly/earthly/issues/991)
- [User-defined commands (UDCs)](https://docs.earthly.dev/docs/guides/udc) have been promoted to GA [#581](https://github.com/earthly/earthly/issues/581).
- Allow running `SAVE ARTIFACT` after `RUN --push` is now GA [#586](https://github.com/earthly/earthly/issues/586).
- `SAVE ARTIFACT --if-exists` and `COPY --if-exists` have been promoted to GA [#588](https://github.com/earthly/earthly/issues/588).
- [Shared cache](https://docs.earthly.dev/docs/guides/shared-cache) and `--ci` mode are now GA [#11](https://github.com/earthly/earthly/issues/11).
- New builtin args `USERPLATFORM`, `USEROS`, `USERARCH`, and `USERVARIANT` which represent the platform, OS, architecture, and processor variant of the system Earthly is being called from [#1251](https://github.com/earthly/earthly/pull/1251). Thanks to @akrantz01 for the contribution!
- Config option for buildkit's `max_parallelism` configuration. Use this to increase parallelism for faster builds or decrease parallelism when resources are constraint. The default is 20. [#1308](https://github.com/earthly/earthly/issues/1308)
- Support for required ARGs (`ARG --required foo`) [#904](https://github.com/earthly/earthly/issues/904). Thanks to @camerondurham for the contribution!
- Extended auto-completion to be build-arg aware. Typing `earthly +my-target --<tab><tab>` now prints possible build-args specific to `+my-target`. [#1330](https://github.com/earthly/earthly/pull/1330).
- The console output now has an improved structure [#1226](https://github.com/earthly/earthly/pull/1226).

### Fixed

- Eliminated some spurious warnings (`ReadDataPacket failed`, `Failed to connect to terminal`, `failed to read from stdin` and others) [#1241](https://github.com/earthly/earthly/pull/1241).
- Minor fixes related to the experimental Podman support [#1239](https://github.com/earthly/earthly/pull/1239).
- Improved some error messages related to frontend detection [#1250](https://github.com/earthly/earthly/pull/1250).
- Fixed Podman's ability to load OCI images [#1287](https://github.com/earthly/earthly/pull/1287).
- Fixed homebrew installation on macOS 12. [#1370](https://github.com/earthly/earthly/pull/1370), [homebrew/earthly#13](https://github.com/earthly/homebrew-earthly/pull/13)
- `failed due to failed to autodetect a supported frontend` errors will now include underlying reason for failure
- Cache export was not honoring `EARTHLY_MAX_REMOTE_CACHE` setting.
- Buildkit logs were not being sent to `earthly-buildkitd` container's output.
- kind required permissions were not available in earthly-buildkitd.

## v0.6.0-rc3 - 2021-11-15

### Fixed

- cache export was not honoring `EARTHLY_MAX_REMOTE_CACHE` setting
- buildkit logs were not being sent to `earthly-buildkitd` container's output.
- kind required permissions were not available in earthly-buildkitd.

### Changed

- docker and fsutils versions were set to match versions defined in earthly's buildkit fork.

## v0.6.0-rc2 - 2021-11-01

### Fixed

- `failed due to failed to autodetect a supported frontend` errors will now include underlying reason for failure

### Changed

- Buildkit was updated to `d47b46cf2a16ca80a958384282e8028285b1866d`.

## v0.6.0-rc1 - 2021-10-28

This version promotes a number of features that have been previously in Experimental and Beta status. To make use of the features in this version you need to declare `VERSION 0.6` at the top of your Earthfile. If a version is not declared, then Earthly's interpreter will assume `VERSION 0.5`.

If you are not ready to update your scripts to take advantage of `VERSION 0.6`, then you may upgrade Earthly anyway and your scripts should continue to work as before, provided that they either declare `VERSION 0.5` or they don't declare a version at all.

Declaring `VERSION 0.6` is equivalent to

```
VERSION \
  --use-copy-include-patterns \
  --referenced-save-only \
  --for-in \
  --require-force-for-unsafe-saves \
  --no-implicit-ignore \
  0.5
```

It is recommended to use `VERSION 0.6` instead as individual feature flags don't guarantee proper forwards-backwards compatibility. Note, however, that Earthly `0.5.*` is not able to run a `VERSION 0.6` Earthfile and will return an error.

For more information on the individual Earthfile feature flags see the [Earthfile version-specific features page](https://docs.earthly.dev/docs/earthfile/features).

### Added

- Earthly now performs local image outputs to the local Docker daemon through a built-in registry. This speeds up the process drastically as common layers no longer need to be transferred over [#500](https://github.com/earthly/earthly/issues/500).
- Earthly now enables additional parallelism to speed up certain operations that were previously serialized [#888](https://github.com/earthly/earthly/issues/888). Note that this setting was previously controlled by `--conversion-parallelism` flag or the `EARTHLY_CONVERSION_PARALLELISM` environment variable while in experimental stage. It has now been moved as part of the Earthly config and has been promoted to GA.
- `COPY` transfers are sped up as only the necessary files are sent over to BuildKit [#1062](https://github.com/earthly/earthly/issues/1062).
- [`WITH DOCKER`](https://docs.earthly.dev/docs/earthfile#with-docker) has been promoted to GA [#576](https://github.com/earthly/earthly/issues/576).
- [`FROM DOCKERFILE`](https://docs.earthly.dev/docs/earthfile#from-dockerfile) has been promoted to GA.
- Support for Apple Silicon M1 has been promoted to GA [#722](https://github.com/earthly/earthly/issues/722).
- [Multi-platform builds](https://docs.earthly.dev/docs/guides/multi-platform) have been promoted to GA [#536](https://github.com/earthly/earthly/issues/536).
- Mounting secrets as files have been promoted as GA [#579](https://github.com/earthly/earthly/issues/579).
- [`VERSION`](https://docs.earthly.dev/docs/earthfile#version) has been promoted to GA [#991](https://github.com/earthly/earthly/issues/991)
- [User-defined commands (UDCs)](https://docs.earthly.dev/docs/guides/udc) have been promoted to GA [#581](https://github.com/earthly/earthly/issues/581).
- Allow running `SAVE ARTIFACT` after `RUN --push` is now GA [#586](https://github.com/earthly/earthly/issues/586).
- `SAVE ARTIFACT --if-exists` and `COPY --if-exists` have been promoted to GA [#588](https://github.com/earthly/earthly/issues/588).
- [Shared cache](https://docs.earthly.dev/docs/guides/shared-cache) and `--ci` mode are now GA [#11](https://github.com/earthly/earthly/issues/11).
- [`LOCALLY`](https://docs.earthly.dev/docs/earthfile#locally) has been promoted to GA [#580](https://github.com/earthly/earthly/issues/580).
- [`RUN --interactive` and `RUN --interactive-keep`](https://docs.earthly.dev/docs/earthfile#run) have been promoted to GA [#693](https://github.com/earthly/earthly/issues/693).
- [`IF`](https://docs.earthly.dev/docs/earthfile#if) and [`FOR`](https://docs.earthly.dev/docs/earthfile#for) have been promoted to GA [#779](https://github.com/earthly/earthly/issues/779).
- If a `SAVE ARTIFACT` is unsafe (writing to a directory outside of the Earthfile directory), it'll require the `--force` flag.
- `.earthlyignore` no longer includes any implicit entries like `Earthfile` or `.earthlyignore`. These will need to be specified explicitly. [#1294](https://github.com/earthly/earthly/issues/1294)
- The console output now has an improved structure [#1226](https://github.com/earthly/earthly/pull/1226).
- Fixed homebrew installation on macOS 12. [#1370](https://github.com/earthly/earthly/pull/1370), [homebrew/earthly#13](https://github.com/earthly/homebrew-earthly/pull/13)
### Changed

<!--changelog-parser-ignore-start-->
- What Earthly outputs locally has changed in a way that is not backwards compatible. For an artifact or an image to be produced locally it needs to be part of a `BUILD` chain (or be part of the target being directly built). Artifacts and images introduced through `FROM` or `COPY` are no longer output locally.

  To update existing scripts, you may issue a duplicate `BUILD` in addition to a `FROM` (or a `COPY`), should you wish for the referenced target to perform output.

  For example, the following script

  ```
  FROM +some-target
  COPY +another-target/my-artifact ./
  ```

  could become

  ```
  FROM +some-target
  BUILD +some-target
  COPY +another-target/my-artifact ./
  BUILD +another-target
  ```

  in order to produce the same outputs.

  For more details see [#896](https://github.com/earthly/earthly/issues/896).
- The syntax for passing build args has been changed.

  Earthly v0.5 (old way)

  ```
  FROM --build-arg NAME=john +some-target
  COPY --build-arg NAME=john +something/my-artifact ./
  WITH DOCKER --build-arg NAME=john --load +another-target
    ...
  END
  ```

  Earthly v0.6 (new way)

  ```
  FROM +some-target --NAME=john
  COPY (+something/my-artifact --NAME=john) ./
  WITH DOCKER --load (+another-target --NAME=john)
    ...
  END
  ```

  Passing build args on the command-line has also changed similarly:

  Earthly v0.5 (old way)

  ```
  earthly --build-arg NAME=john +some-target
  ```

  Earthly v0.6 (new way)

  ```
  earthly +some-target --NAME=john
  ```

  This change is part of the [UDC proposal #581](https://github.com/earthly/earthly/issues/581). The old way of passing args is deprecated and will be removed in a future version (however, it still works in 0.6).
<!--changelog-parser-ignore-end-->
- Add builtin args `USERPLATFORM`, `USEROS`, `USERARCH`, and `USERVARIANT` which represent the platform, OS, architecture, and processor variant of the system Earthly is being called from [#1251](https://github.com/earthly/earthly/pull/1251). Thanks to @akrantz01 for the contribution!
- Support for required ARGs (`ARG --required foo`) [#904](https://github.com/earthly/earthly/issues/904). Thanks to @camerondurham for the contribution!
- Add a config item for buildkit's `max_parallelism` configuration. Use this to increase parallelism for faster builds or decrease parallelism when resources are constraint. The default is 20. [#1308](https://github.com/earthly/earthly/issues/1308)
- Extend auto-completion to be build-arg aware. Typing `earthly +my-target --<tab><tab>` now prints possible build-args specific to `+my-target`. [#1330](https://github.com/earthly/earthly/pull/1330).
- Buildkit was updated to `d429b0b32606b5ea52e6be4a99b69d67b7c722b2`. This includes a number of bug fixes, including eliminating crashes due to `panic failed to get edge`.

### Fixed

- Eliminated some spurious warnings (`ReadDataPacket failed`, `Failed to connect to terminal`, `failed to read from stdin` and others) [#1241](https://github.com/earthly/earthly/pull/1241).
- Minor fixes related to the experimental Podman support [#1239](https://github.com/earthly/earthly/pull/1239).
- Improved some error messages related to frontend detection [#1250](https://github.com/earthly/earthly/pull/1250).
- Fixed Podman's ability to load OCI images [#1287](https://github.com/earthly/earthly/pull/1287).

## v0.5.24 - 2021-09-30

### Added

- New `--output` flag, which forces earthly to enable outputs, even when running under `--ci` mode [#1200](https://github.com/earthly/earthly/issues/1200).
- Experimental support for Podman [#760](https://github.com/earthly/earthly/issues/760).
- Automatically adds compatibility arguments for cases where docker is running under user namespaces.

### Fixed

- Removed spurious `BuildKit and Local Registry URLs are pointed at different hosts (earthly-buildkitd vs. 127.0.0.1)` warning.
- Scrub git credentials when running under --debug mode.
- "FROM DOCKERFILE" command was ignoring the path (when run on a remote target), which prevented including dockerfiles which were named something else.
- Removed the creation of a temporary output directory when run in `--no-output` mode, or when building targets that don't output artifacts,
  the temporary directory is now created just before it is needed.
- Fixed race condition involving `WITH DOCKER` and `IF` statements, which resulted in `failed to solve: NotFound: no access allowed to dir` errors.

## v0.5.23 - 2021-08-24

- introduced `COPY --if-exists` which allows users to ignore errors which would have occurred [if the file did not exist](https://docs.earthly.dev/docs/earthfile#if-exists).
- introduced new `ip_tables` config option for controlling which iptables binary is used; fixes #1160
- introduced warning message when saving to paths above the current directory of the current Earthfile; these warnings will eventually become errors unless the `--force` [flag](https://docs.earthly.dev/docs/earthfile#force) is passed to `SAVE ARTIFACT`.
- fixed remote BuildKit configuration options being ignored; fixes #1177
- suppressed erroneous internal-term error messages which occurred when running under non-interactive ( e.g. `--ci` ) modes; fixes #1108
- changed help text for `--artifact` mode
- deb and yum packages no longer clear the earthly cache on upgrades


## v0.5.22 - 2021-08-11

- when running under `--ci` mode, earthly now raises an error if a user attempts to use the interactive debugger
- updated underlying BuildKit version
- print all request and responses to BuildKit when running under --debug mode
- support for specifying files to ignore under `.earthlyignore` in addition to `.earthignore`; an error is raised if both exist
- new ARG `EARTHLY_GIT_SHORT_HASH` will contain an 8 char representation of the current git commit hash
- new ARG `EARTHLY_GIT_COMMIT_TIMESTAMP` will contain the timestamp of the current git commit
- new ARG `EARTHLY_SOURCE_DATE_EPOCH` will contain the same value as `EARTHLY_GIT_COMMIT_TIMESTAMP` or 0 when the timestamp is not available
- only directly referenced artifacts or images will be saved when the VERSION's --referenced-save-only feature flag is defined #896
- experimental support for FOR statements, when the VERSION's --for-in feature flag is defined #1142
- fixes bug where error was not being repeated as the final output
- fixes bug where HTTPS-based git credentials were leaked to stdout

## v0.5.20 - 2021-07-22

- Support for passing true/false values to boolean flags #1109
- fixes error that stated `http is insecure` when configuring a HTTPS git source. #1115

## v0.5.19 - 2021-07-21

- Improved selective file-transferring via BuildKit's include patterns; this feature is currently disabled by default, but can be enabled by including the `--use-copy-include-patterns` feature-flag in the `VERSION` definition (e.g. add `VERSION --use-copy-include-patterns 0.5` to the top of your Earthfiles). This will become enabled by default in a later version.
- Support for host systems that have `nf_tables` rather than `ip_tables`.
- Show hidden dev flags when `EARTHLY_AUTOCOMPLETE_HIDDEN="1"` is set (or when running a custom-built version).
- Improved crash logs.

## v0.5.18 - 2021-07-08

- Added a `--symlink-no-follow` flag to allow copying invalid symbolic links (https://github.com/earthly/earthly/issues/1067)
- Updated BuildKit, which contains a fix for "failed to get edge" panic errors (https://github.com/earthly/earthly/issues/1016)
- Fix bug that prevented using an absolute path to reference targets which contained relative imports
- Added option to disable analytics data collection when environment variables `EARTHLY_DISABLE_ANALYTICS` or `DO_NOT_TRACK` are set.
- Include version and help flags in autocompletion output.

## v0.5.17 - 2021-06-15

- Begin experimental official support for `earthly/earthly` and `earthly/buildkitd` images; including a new `entrypoint` for `earthly/earthly` (https://github.com/earthly/earthly/pull/1050)
- When running in `verbose` mode, log all files sent to BuildKit (https://github.com/earthly/earthly/pull/1051, https://github.com/earthly/earthly/pull/1056)
- Adjust `deb` and `rpm` packages to auto-install the shell completions though post-installation mechanisms (https://github.com/earthly/earthly/pull/1019, https://github.com/earthly/earthly/pull/1057)

## v0.5.16 - 2021-06-03

- fixes handling of `Error getting earthly dir` lookup failures which prevents earthly from running (https://github.com/earthly/earthly/issues/1026)
- implements ability to perform local exports via buildkit-hosted local registry in order to speed up exports; the feature is currently disabled by default but can be enabled with `earthly config global.local_registry_host 'tcp://127.0.0.1:8371'` (https://github.com/earthly/earthly/issues/500)

## v0.5.15 - 2021-05-27

- `earthly config` is no longer experimental. (https://github.com/earthly/earthly/pull/979)
- Running a target, will now `bootstrap` automatically, if it looks like `earthly bootstrap` has not been run yet. (https://github.com/earthly/earthly/pull/989)
- `earthly bootstrap` ensures the permissions on the `.earthly` folder are correct (belonging to the user) ( https://github.com/earthly/earthly/pull/993)
- Cache mount ID now depends on a target input hash which does not include inactive variables (https://github.com/earthly/earthly/pull/1000)
- Added `EARTHLY_TARGET_PROJECT_NO_TAG` built-in argument (https://github.com/earthly/earthly/pull/1011)
- When `~` is used as the path to a secret file, it now expands as expected. (https://github.com/earthly/earthly/pull/977)
- Use the environment-specified `$HOME`, unless `$SUDO_USER` is set. If it is, use the users home directory. (https://github.com/earthly/earthly/pull/1015)


## v0.5.14 - 2021-05-27

- `earthly config` is no longer experimental. (https://github.com/earthly/earthly/pull/979)
- Running a target, will now `bootstrap` automatically, if it looks like `earthly bootstrap` has not been run yet. (https://github.com/earthly/earthly/pull/989)
- `earthly bootstrap` ensures the permissions on the `.earthly` folder are correct (belonging to the user) ( https://github.com/earthly/earthly/pull/993)
- Cache mount ID now depends on a target input hash which does not include inactive variables (https://github.com/earthly/earthly/pull/1000)
- Added `EARTHLY_TARGET_PROJECT_NO_TAG` built-in argument (https://github.com/earthly/earthly/pull/1011)
- When `~` is used as the path to a secret file, it now expands as expected. (https://github.com/earthly/earthly/pull/977)


## v0.5.13 - 2021-05-13

- fixes panic on invalid (or incomplete) `~/.netrc` file (https://github.com/earthly/earthly/issues/980)

## v0.5.12 - 2021-05-07

- Adds a retry for remote BuildKit hosts when using the `EARTHLY_BUILDKIT_HOST` configuration option. (#952)
- Re-fetch credentials when they expire (#957)
- Make use of `~/.netrc` credentials when no config is set under `~/.earthly/config.yml` (#964)
- Make use of auth credentials when performing a GIT CLONE command within an Earthfile. (#964)
- Improved error output when desired secret does not exist, including the name of the missing secret. (#972)
- Warn if `build-arg` appears after the target in CLI invocations.(#959)

## v0.5.11 - 2021-04-27

- Support for `FROM DOCKERFILE -f` (https://github.com/earthly/earthly/pull/950)
- Fixes missing access to global arguments in user defined commands (https://github.com/earthly/earthly/pull/947)
- Users's `~/.earthly` directory is now referenced when earthly is invoked with sudo


## v0.5.10 - 2021-04-19

- Added ability to run `WITH DOCKER` under `LOCALLY` (https://github.com/earthly/earthly/pull/840)
- Fix `FROM DOCKERFILE` `--build-arg`s not being passed correctly (https://github.com/earthly/earthly/issues/932)
- Docs: Add uninstall instructions
- Docs: Improve onboarding tutorial based on user feedback


## v0.5.9 - 2021-04-05

- [**experimental**] Improved parallelization when using commands such as `IF`, `WITH DOCKER`, `FROM DOCKERFILE`, `ARG X=$(...)` and others. To enable this feature, pass `--conversion-parallelism=5` or set `EARTHLY_CONVERSION_PARALLELISM=5`. (https://github.com/earthly/earthly/issues/888)
- Auto-detect MTU (https://github.com/earthly/earthly/issues/847)
- MTU may set via config `earthly config global.cni_mtu 12345` (https://github.com/earthly/earthly/pull/906)
- Hide `--debug` flag since it is only used for development on Earthly itself
- Download and start buildkitd as part of the earthly bootstrap command
- Improved buildkitd startup logic (https://github.com/earthly/earthly/pull/892)
- Check for reserved target names and disallow them (e.g. `+base`) (https://github.com/earthly/earthly/pull/898)
- Fix use of self-hosted repositories when a subdirectory is used (https://github.com/earthly/earthly/pull/897)


## v0.5.8 - 2021-03-23

- [**experimental**] Support for ARGs in user-defined commands (UDCs). UDCs are templates (much like functions in regular programming languages), which can be used to define a series of steps to be executed in sequence. In other words, it is a way to reuse common build steps in multiple contexts. This completes the implementation of UDCs and the feature is now in **experimental** phase (https://github.com/earthly/earthly/issues/581). For more information see the [UDC guide](https://docs.earthly.dev/guides/udc).
- [**experimental**] New command: `IMPORT` (https://github.com/earthly/earthly/pull/868)
  ```
  IMPORT github.com/foo/bar:v1.2.3
  IMPORT github.com/foo/buz:main AS zulu

  ...

  FROM bar+target
  BUILD zulu+something
  ```
- Fix handling of some escaped quotes (https://github.com/earthly/earthly/issues/859)
- Fix: empty targets are now valid (https://github.com/earthly/earthly/pull/872)
- Fix some line continuation issues (https://github.com/earthly/earthly/pull/873 & https://github.com/earthly/earthly/pull/874)
- Earthly now limits parallelism to `20`. This fixes some very large builds attempting to use resources all at the same time
- Automatically retry TLS handshake timeout errors

## v0.5.7 - 2021-03-13

- raise error when duplicate target names exists in Earthfile
- basic user defined commands (experimental)
- cleans up console output for saving artifacts (#848)
- implement support for WORKDIR under LOCALLY targets
- fix zsh autocompletion issue for mac users
  If the autocompletion bug persists for anyone (e.g. seeing an error like `command not found: __earthly__`), and the issues persists after upgrading to v0.5.7; it might be necessary to delete the _earthly autocompletion file before re-running earthly bootstrap (or alternatively manually replace `__earthly__` with the full path to the earthly binary).

## v0.5.6 - 2021-03-09

- This release removes the `ongoing` updates "Provide intermittent updates on long-running targets (#844)" from the previous release, as it has issues in the interactive mode.

## v0.5.5 - 2021-03-08

- Keep `.git` directory in build context. (#815 )
- Wait extra time for buildkitd to start if the cache is larger than 30 GB  (#827)
- *Experimental:* Allow RUN commands to open an interactive session (`RUN --interactive`), with the option to save the manual changes into the final image (`RUN --interactive-keep`) (#833)
- Provide intermittent updates on long-running targets (#844)
- Fix ZSH autocompletion in some instances (#838)

## v0.5.4 - 2021-02-26

- New experimental `--strict` flag, which doesn't allow the use of `LOCALLY`. `--strict` is now implied when using `--ci`. (https://github.com/earthly/earthly/pull/801)
- Add help text when issuing `earthly config <item> --help`. Improved user experience. (https://github.com/earthly/earthly/pull/814)
- Detect if the build doesn't start with a FROM-like command and return a meaningful error. Previously `FROM scratch` was assumed automatically. (https://github.com/earthly/earthly/issues/807)
- Fix an issue where `.tmpXXXXX` directories were created in the current directory (https://github.com/earthly/earthly/pull/821)
- Fix auto-complete in zsh (https://github.com/earthly/earthly/pull/811)
- Improved startup logic for BuildKit daemon, which speeds up some rare edge cases (https://github.com/earthly/earthly/pull/808)
- Print BuildKit logs if it crashes or times out on startup (https://github.com/earthly/earthly/pull/819)
- Create config path if it's missing (https://github.com/earthly/earthly/pull/812)


## v0.5.3 - 2021-02-24

- Support for conditional logic using new `IF`, `ELSE IF`, and `ELSE` keywords (required for #779)
- Support for copying artifacts to `LOCALLY` targets (required for #580)

### Fixed
- segfault when no output or error is displayed (fixes #798)
- unable to run earthly in docker container with mounted host-docker socket (fixes #791)
- `./.tmp-earthly-outXXXXXX` temp files are now stored under `./.tmp-earthly-out/tmpXXXXXX` and are correctly excluded from the build context


## v0.5.2 - 2021-02-18

- New experimental command for editing the Earthly config (https://github.com/earthly/earthly/issues/675)
- `SAVE IMAGE --push` after a `RUN --push` now includes the effects of the `RUN --push` too (https://github.com/earthly/earthly/pull/754)
- Improved syntax errors when parsing Earthfiles
- Improved error message when QEMU is missing
- Fix `earthly-linux-arm64` binary - was a Mac binary by mistake (https://github.com/earthly/earthly/issues/789)
- Fix override of build arg not being detected properly (https://github.com/earthly/earthly/pull/790)
- Fix image export error when it doesn't contain any `RUN` commands (https://github.com/earthly/earthly/issues/782)


## v0.5.1 - 2021-02-08

- Support for SAVE ARTIFACT under LOCALLY contexts; this allows one to run a command locally and save the output to a different container.
- Support for build arg matrix; supplying multiple `--build-args` for the same value will cause the `BUILD` target to be built for each different build arg value.
- Improvements for Apple M1 support
- Improved errors when parsing invalid Earthfiles (to enable the new experimental code, set the `EARTHLY_ENABLE_AST` variable to `true`)

## v0.5.0 - 2021-02-01

- Switch to BSL license. For [more information about this decision, take a look at our blog post](https://blog.earthly.dev/every-open-core-company-should-be-a-source-available-company/).
- `--platform` setting is now automatically propagated between Earthfiles. In addition, you can now specify the empty string `--platform=` to automatically detect your system's architecture.
- `earthly/dind` images now available for `linux/arm/v7` and `linux/arm64`
- Improved visibility of platform used for each build step, as well as for any build args that have been overridden.
- Allow saving an artifact after a `RUN --push` (https://github.com/earthly/earthly/pull/735)
- Allow specifying `--no-cache` for a single `RUN` command (https://github.com/earthly/earthly/issues/585)
- There are now separate `SUCCESS` lines for each of the two possible phases of an earthly run: `main` and `push`.
- [Support of popular cloud registries for the experimental shared cache feature is now properly documented](https://docs.earthly.dev/guides/shared-cache#compatibility-with-major-registry-providers)
- Fix `SAVE IMAGE --cache-hint` not working correctly (https://github.com/earthly/earthly/issues/744)
- Fix `i/o timeout` errors being cached and requiring BuildKit restart
- Experimental support for running commands on the host system via `LOCALLY` (https://github.com/earthly/earthly/issues/580)
- Bug fixes for Apple Silicon. `earthly-darwin-arm64` binary is now available. Please treat this version as highly experimental for now. (https://github.com/earthly/earthly/issues/722)

## v0.5.0-rc2 - 2021-02-01

- No details provided

## v0.5.0-rc1 - 2021-02-01

- No details provided

## v0.4.6 - 2021-01-29

- No details provided

## v0.4.5 - 2021-01-13

- Fix inconsistent `COPY --dir` behavior [#705](https://github.com/earthly/earthly/issues/705)
- Fix `AS LOCAL` behavior with directories [#703](https://github.com/earthly/earthly/issues/703)

## v0.4.4 - 2021-01-06

- Improved experimental support for arm-based platforms, including Apple M1. Builds run natively on arm platforms now. (For Apple M1, you need to use darwin-amd64 download and have Rosetta 2 installed - the build steps themselves will run natively, however, via the BuildKit daemon).
- Add `SAVE ARTIFACT --if-exists` (https://github.com/earthly/earthly/issues/588)
- Fix an issue where comments at the end of the Earthfile were not allowed (https://github.com/earthly/earthly/issues/681)
- Fix an issue where multiple `WITH DOCKER --load` with the same target, but different image tag were not working (https://github.com/earthly/earthly/issues/685)
- Fix an issue where `SAVE ARTIFACT ./* AS LOCAL` was flattening subdirectories (https://github.com/earthly/earthly/issues/689)
- Binaries for `arm5` and `arm6` are no longer supported

## v0.4.3 - 2020-12-23

- Fix regression for `WITH DOCKER --compose=... --load=...` (https://github.com/earthly/earthly/issues/676)
- Improvements to the multiplatform experimental support. See the [multiplatform example](https://github.com/earthly/earthly/blob/main/examples/multiplatform/Earthfile).

## v0.4.2 - 2020-12-22

- fixed: `EARTHLY_GIT_PROJECT_NAME` contained the raw git URL when HTTPS-based auth was used (fixes #671)
- feature: support for mounting secrets as files rather than environment variables
- feature: experimental support for multi-platform builds
- misc: sending anonymized usage metrics to earthly
