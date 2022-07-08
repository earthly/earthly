# Earthly Changelog

All notable changes to [Earthly](https://github.com/earthly/earthly) will be documented in this file.

## Unreleased

### Changed

- Updated buildkit to include changes up to 12cfc87450c8d4fc31c8c0a09981e4c3fb3e4d9f

### Added

- Adding support for saving artifact from `--interactive-keep`. [#1980](https://github.com/earthly/earthly/issues/1980)

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
- Git config options for non-standard port and path prefix; these options are incompatible with a custom git substition regex.
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

- The feature flag `--exec-after-build` has been enabled retroactively for `VERSION 0.5`. This speeds up largs builds by 15-20%.
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

- Duplicate execution occuring when using ARGs. [#1572](https://github.com/earthly/earthly/issues/1572), [#1582](https://github.com/earthly/earthly/issues/1582)
- Overriding builtin ARG value now displays an error (rather than silently ignoring it).

## v0.6.3 - 2022-01-12

### Changed

- Updated buildkit to contain changes up to `15fb1145afa48bf81fbce41634bdd36c02454f99` from `moby/master`.

### Added

- Expirmental `CACHE` command can be used in Earthfiles to optimize the cache in projects that perform better with incremental changes. For example, a Maven
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
