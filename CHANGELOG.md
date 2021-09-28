# Earthly Changelog

All notable changes to [Earthly](https://github.com/earthly/earthly) will be documented in this file.

## Unreleased

### Added

- New ./release+all target that combines existing release commands into a single command
  that can be run as a single command.

## v0.5.23 - 2021-08-24

- introduced `COPY --if-exists` which allows users to ignore errors which would have occurred [if the file did not exist](https://docs.earthly.dev/docs/earthfile#if-exists).
- introduced new `ip_tables` config option for controlling which iptables binary is used; fixes #1160
- introduced warning message when saving to paths above the current directory of the current Earthfile; these warnings will eventually become errors unless the `--force` [flag](https://docs.earthly.dev/docs/earthfile#force) is passed to `SAVE ARTIFACT`.
- bugfix for remote buildkit configuration options being ignored; fixes #1177
- suppressed erroneous internal-term error messages which occurred when running under non-interactive ( e.g. ci ) modes; fixes #1108
- changed help text for `--artifact` mode
- deb and yum packages no longer clear the earthly cache on upgrades


## v0.5.22 - 2021-08-11

- when running under --ci mode, earthly now raises an error if a user attempts to use the interactive debugger
- updated underlying buildkit version
- print all request and responses to buildkit when running under --debug mode
- support for specifying files to ignore under .earthlyignore in addition to .earthignore; an error is raised if both exist
- new ARG `EARTHLY_GIT_SHORT_HASH` will contain an 8 char representation of the current git commit hash
- new ARG `EARTHLY_GIT_COMMIT_TIMESTAMP` will contain the timestamp of the current git commit
- new ARG `EARTHLY_SOURCE_DATE_EPOCH` will contain the same value as `EARTHLY_GIT_COMMIT_TIMESTAMP` or 0 when the timestamp is not available
- only directly referenced artifacts or images will be saved when the VERSION's --referenced-save-only feature flag is defined #896
- experimental support for FOR statements, when the VERSION's --for-in feature flag is defined #1142
- fixes bug where error was not being repeated as the final output
- fixes bug where https-based git credentials were leaked to stdout

## v0.5.20 - 2021-07-22

- Support for passing true/false values to boolean flags #1109
- fixes error that stated `http is insecure` when configuring a https git source. #1115

## v0.5.19 - 2021-07-21

- Improved selective file-transferring via buildkit's include patterns; this feature is currently disabled by default, but can be enabled by including the `--use-copy-include-patterns` feature-flag in the `VERSION` definition (e.g. add `VERSION --use-copy-include-patterns 0.5` to the top of your Earthfiles). This will become enabled by default in a later version.
- Support for host systems that have `nf_tables` rather than `ip_tables`.
- Show hidden dev flags when `EARTHLY_AUTOCOMPLETE_HIDDEN="1"` is set (or when running a custom-built version).
- Improved crash logs.

## v0.5.18 - 2021-07-08

* Added a `--symlink-no-follow` flag to allow copying invalid symbolic links (https://github.com/earthly/earthly/issues/1067)
* Updated buildkit, which contains a fix for "failed to get edge" panic errors (https://github.com/earthly/earthly/issues/1016)
* Fix bug that prevented using an absolute path to reference targets which contained relative imports
* Added option to disable analytics data collection when environment variables `EARTHLY_DISABLE_ANALYTICS` or `DO_NOT_TRACK` are set.
* Include version and help flags in autocompletion output.

## v0.5.17 - 2021-06-15

* Begin experimental official support for `earthly/earthly` and `earthly/buildkitd` images; including a new `entrypoint` for `earthly/earthly` (https://github.com/earthly/earthly/pull/1050)
* When running in `verbose` mode, log all files sent to buildkit (https://github.com/earthly/earthly/pull/1051, https://github.com/earthly/earthly/pull/1056)
* Adjust `deb` and `rpm` packages to auto-install the shell completions though post-installation mechanisms (https://github.com/earthly/earthly/pull/1019, https://github.com/earthly/earthly/pull/1057)

## v0.5.16 - 2021-06-03

- fixes handling of `Error getting earthly dir` lookup failures which prevents earthly from running (https://github.com/earthly/earthly/issues/1026)
- implements ability to perform local exports via buildkit-hosted local registery in order to speed up exports; the feature is currently disabled by default but can be enabled with `earthly config global.local_registry_host 'tcp://127.0.0.1:8371'` (https://github.com/earthly/earthly/issues/500)

## v0.5.15 - 2021-05-27

* `earthly config` is no longer experimental. (https://github.com/earthly/earthly/pull/979)
* Running a target, will now `bootstrap` automatically, if it looks like `earthly bootstrap` has not been run yet. (https://github.com/earthly/earthly/pull/989)
* `earthly bootstrap` ensures the permissions on the `.earthly` folder are correct (belonging to the user) ( https://github.com/earthly/earthly/pull/993)
* Cache mount ID now depends on a target input hash which does not include inactive variables (https://github.com/earthly/earthly/pull/1000)
* Added `EARTHLY_TARGET_PROJECT_NO_TAG` built-in argument (https://github.com/earthly/earthly/pull/1011)
* When `~` is used as the path to a secret file, it now expands as expected. (https://github.com/earthly/earthly/pull/977)
* Use the environment-specified `$HOME`, unless `$SUDO_USER` is set. If it is, use the users home directory. (https://github.com/earthly/earthly/pull/1015)


## v0.5.14 - 2021-05-27

* `earthly config` is no longer experimental. (https://github.com/earthly/earthly/pull/979)
* Running a target, will now `bootstrap` automatically, if it looks like `earthly bootstrap` has not been run yet. (https://github.com/earthly/earthly/pull/989)
* `earthly bootstrap` ensures the permissions on the `.earthly` folder are correct (belonging to the user) ( https://github.com/earthly/earthly/pull/993)
* Cache mount ID now depends on a target input hash which does not include inactive variables (https://github.com/earthly/earthly/pull/1000)
* Added `EARTHLY_TARGET_PROJECT_NO_TAG` built-in argument (https://github.com/earthly/earthly/pull/1011)
* When `~` is used as the path to a secret file, it now expands as expected. (https://github.com/earthly/earthly/pull/977)


## v0.5.13 - 2021-05-13

- fixes panic on invalid (or incomplete) `~/.netrc` file (https://github.com/earthly/earthly/issues/980)

## v0.5.12 - 2021-05-07

* Adds a retry for remote Buildkit hosts when using the `EARTHLY_BUILDKIT_HOST` configuration option. (#952)
* Re-fetch credentials when they expire (#957)
* Make use of `~/.netrc` credentials when no config is set under `~/.earthly/config.yml` (#964)
* Make use of auth credentials when performing a GIT CLONE command within an Earthfile. (#964)
* Improved error output when desired secret does not exist, including the name of the missing secret. (#972)
* Warn if `build-arg` appears after the target in CLI invocations.(#959)

## v0.5.11 - 2021-04-27

Changes:

* Support for `FROM DOCKERFILE -f` (https://github.com/earthly/earthly/pull/950)
* Fixes missing access to global arguments in user defined commands (https://github.com/earthly/earthly/pull/947)
* Users's `~/.earthly` directory is now referenced when earthly is invoked with sudo


## v0.5.10 - 2021-04-19

Changes:

* Added ability to run `WITH DOCKER` under `LOCALLY` (https://github.com/earthly/earthly/pull/840)
* Fix `FROM DOCKERFILE` `--build-arg`s not being passed correctly (https://github.com/earthly/earthly/issues/932)
* Docs: Add uninstall instructions
* Docs: Improve onboarding tutorial based on user feedback


## v0.5.9 - 2021-04-05

Changes:

* [**experimental**] Improved parallelization when using commands such as `IF`, `WITH DOCKER`, `FROM DOCKERFILE`, `ARG X=$(...)` and others. To enable this feature, pass `--conversion-parallelism=5` or set `EARTHLY_CONVERSION_PARALLELISM=5`. (https://github.com/earthly/earthly/issues/888)
* Auto-detect MTU (https://github.com/earthly/earthly/issues/847)
* MTU may set via config `earthly config global.cni_mtu 12345` (https://github.com/earthly/earthly/pull/906)
* Hide `--debug` flag since it is only used for development on Earthly itself
* Download and start buildkitd as part of the earthly bootstrap command
* Improved buildkitd startup logic (https://github.com/earthly/earthly/pull/892)
* Check for reserved target names and disallow them (e.g. `+base`) (https://github.com/earthly/earthly/pull/898)
* Fix use of self-hosted repos when a subdirectory is used (https://github.com/earthly/earthly/pull/897)


## v0.5.8 - 2021-03-23


Changes:

* [**experimental**] Support for ARGs in user-defined commands (UDCs). UDCs are templates (much like functions in regular programming languages), which can be used to define a series of steps to be executed in sequence. In other words, it is a way to reuse common build steps in multiple contexts. This completes the implementation of UDCs and the feature is now in **experimental** phase (https://github.com/earthly/earthly/issues/581). For more information see the [UDC guide](https://docs.earthly.dev/guides/udc).
* [**experimental**] New command: `IMPORT` (https://github.com/earthly/earthly/pull/868)
  ```
  IMPORT github.com/foo/bar:v1.2.3
  IMPORT github.com/foo/buz:main AS zulu
  
  ...
  
  FROM bar+target
  BUILD zulu+something
  ```
* Fix handling of some escaped quotes (https://github.com/earthly/earthly/issues/859)
* Fix: empty targets are now valid (https://github.com/earthly/earthly/pull/872)
* Fix some line continuation issues (https://github.com/earthly/earthly/pull/873 & https://github.com/earthly/earthly/pull/874)
* Earthly now limits parallelism to `20`. This fixes some very large builds attempting to use resources all at the same time
* Automatically retry TLS handshake timeout errors


## v0.5.7 - 2021-03-13

- raise error when duplicate target names exists in Earthfile
- fix zsh autocompletion issue for mac users
- basic user defined commands (experimental)
- cleans up console output for saving artifacts (#848)
- implement support for WORKDIR under LOCALLY targets

If the autocompletion bug persists for anyone (e.g. seeing an error like `command not found: __earthly__`), and the issues persists after upgrading to v0.5.7; it might be necessary to delete the _earthly autocompletion file before re-running earthly bootstrap (or alternatively manually replace `__earthly__` with the full path to the earthly binary).

## v0.5.6 - 2021-03-09

This release removes the `ongoing` updates "Provide intermittent updates on long-running targets (#844)" from the previous release, as it has issues in the interactive mode.

## v0.5.5 - 2021-03-08

* Keep `.git` directory in build context. (#815 )
* Wait extra time for buildkitd to start if the cache is larger than 30 GB  (#827)
* *Experimental:* Allow RUN commands to open an interactive session (`RUN --interactive`), with the option to save the manual changes into the final image (`RUN --interactive-keep`) (#833)
* Provide intermittent updates on long-running targets (#844)
* Fix ZSH autocompletion in some instances (#838)

## v0.5.4 - 2021-02-26

* New experimental `--strict` flag, which doesn't allow the use of `LOCALLY`. `--strict` is now implied when using `--ci`. (https://github.com/earthly/earthly/pull/801)
* Add help text when issuing `earthly config <item> --help`. Improved user experience. (https://github.com/earthly/earthly/pull/814)
* Detect if the build doesn't start with a FROM-like command and return a meaningful error. Previously `FROM scratch` was assumed automatically. (https://github.com/earthly/earthly/issues/807)
* Fix an issue where `.tmpXXXXX` directories were created in the current directory (https://github.com/earthly/earthly/pull/821)
* Fix auto-complete in zsh (https://github.com/earthly/earthly/pull/811)
* Improved startup logic for buildkit daemon, which speeds up some rare edge cases (https://github.com/earthly/earthly/pull/808)
* Print buildkit logs if it crashes or times out on startup (https://github.com/earthly/earthly/pull/819)
* Create config path if it's missing (https://github.com/earthly/earthly/pull/812)


## v0.5.3 - 2021-02-24

Changes
 - Support for conditional logic using new `IF`, `ELSE IF`, and `ELSE` keywords (required for #779)
 - Support for copying artifacts to `LOCALLY` targets (required for #580)
 - bugfix: segfault when no output or error is displayed (fixes #798)
 - bugfix: unable to run earthly in docker container with mounted host-docker socket (fixes #791)
 - bugfix: `./.tmp-earthly-outXXXXXX` temp files are now stored under `./.tmp-earthly-out/tmpXXXXXX` and are correctly excluded from the builld context 


## v0.5.2 - 2021-02-18

Changes:

* New experimental command for editing the Earthly config (https://github.com/earthly/earthly/issues/675)
* `SAVE IMAGE --push` after a `RUN --push` now includes the effects of the `RUN --push` too (https://github.com/earthly/earthly/pull/754)
* Improved syntax errors when parsing Earthfiles
* Improved error message when QEMU is missing
* Fix `earthly-linux-arm64` binary - was a Mac binary by mistake (https://github.com/earthly/earthly/issues/789)
* Fix override of build arg not being detected properly (https://github.com/earthly/earthly/pull/790)
* Fix image export error when it doesn't contain any `RUN` commands (https://github.com/earthly/earthly/issues/782)


## v0.5.1 - 2021-02-08

Changes
 - Support for SAVE ARTIFACT under LOCALLY contexts; this allows one to run a command locally and save the output to a different container.
 - Support for build arg matrix; supplying multiple `--build-args` for the same value will cause the `BUILD` target to be built for each different build arg value.
 - Improvements for Apple M1 support
 - Improved errors when parsing invalid Earthfiles (to enable the new experimental code, set the `EARTHLY_ENABLE_AST` variable to `true`)

## v0.5.0 - 2021-02-01

Changes:
* Switch to BSL license. For [more information about this decision, take a look at our blog post](https://blog.earthly.dev/every-open-core-company-should-be-a-source-available-company/).
* `--platform` setting is now automatically propagated between Earthfiles. In addition, you can now specify the empty string `--platform=` to automatically detect your system's architecture.
* `earthly/dind` images now available for `linux/arm/v7` and `linux/arm64`
* Improved visibility of platform used for each build step, as well as for any build args that have been overridden.
* Allow saving an artifact after a `RUN --push` (https://github.com/earthly/earthly/pull/735)
* Allow specifying `--no-cache` for a single `RUN` command (https://github.com/earthly/earthly/issues/585)
* There are now separate `SUCCESS` lines for each of the two possible phases of an earthly run: `main` and `push`.
* [Support of popular cloud registries for the experimental shared cache feature is now properly documented](https://docs.earthly.dev/guides/shared-cache#compatibility-with-major-registry-providers)
* Fix `SAVE IMAGE --cache-hint` not working correctly (https://github.com/earthly/earthly/issues/744)
* Fix `i/o timeout` errors being cached and requiring BuildKit restart
* Experimental support for running commands on the host system via `LOCALLY` (https://github.com/earthly/earthly/issues/580)
* Bug fixes for Apple Silicon. `earthly-darwin-arm64` binary is now available. Please treat this version as highly experimental for now. (https://github.com/earthly/earthly/issues/722)

## v0.5.0-rc2 - 2021-02-01

No details provided

## v0.5.0-rc1 - 2021-02-01

No details provided

## v0.4.6 - 2021-01-29

No details provided

## v0.4.5 - 2021-01-13

 * Fix inconsistent `COPY --dir` behavior [https://github.com/earthly/earthly/issues/705](#705)
 * Fix `AS LOCAL` behavior with directories [https://github.com/earthly/earthly/issues/703](#703)

## v0.4.4 - 2021-01-06

Changes:
* Improved experimental support for arm-based platforms, including Apple M1. Builds run natively on arm platforms now. (For Apple M1, you need to use darwin-amd64 download and have Rosetta 2 installed - the build steps themselves will run natively, however, via the buildkit daemon).
* Add `SAVE ARTIFACT --if-exists` (https://github.com/earthly/earthly/issues/588)
* Fix an issue where comments at the end of the Earthfile were not allowed (https://github.com/earthly/earthly/issues/681)
* Fix an issue where multiple `WITH DOCKER --load` with the same target, but different image tag were not working (https://github.com/earthly/earthly/issues/685)
* Fix an issue where `SAVE ARTIFACT ./* AS LOCAL` was flattening subdirectories (https://github.com/earthly/earthly/issues/689)
* Binaries for `arm5` and `arm6` are no longer supported

## v0.4.3 - 2020-12-23

Changes:
* Fix regression for `WITH DOCKER --compose=... --load=...` (https://github.com/earthly/earthly/issues/676)
* Improvements to the multiplatform experimental support. See the [multiplatform example](https://github.com/earthly/earthly/blob/main/examples/multiplatform/Earthfile).

## v0.4.2 - 2020-12-22

Changes:
 - bugfix: EARTHLY_GIT_PROJECT_NAME contained the raw git url when https-based auth was used (fixes #671)
 - feature: support for mounting secrets as files rather than environment variables
 - feature: experimental support for multi-platform builds
 - misc: sending anonymized usage metrics to earthly
