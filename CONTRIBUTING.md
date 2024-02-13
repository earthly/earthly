# Contributing

## Code of Conduct

Please refer to [code-of-conduct.md](./code-of-conduct.md) for details.

## Using Earthly prerelease

To build Earthly from source, you need the same requirements as Earthly. We recommend that you use the prerelease version of Earthly for development purposes. To launch the prerelease Earthly, simply use the `./earthly` script provided in the root of the earthly repository. The prerelease Earthly tracks the version on main. You can use `./earthly --version` to identify which Git hash was used to build it.

## Building from source

To build Earthly from source for your target system, use

* Linux and WSL
    ```bash
    ./earthly +for-linux
    ```
* Mac
    ```bash
    ./earthly +for-darwin
    ```
* Mac with M1 chip
    ```bash
    ./earthly +for-darwin-m1
    ```

This builds the earthly binary in `./build/*/*/earthly`, typically one of:

* `./build/linux/amd64/earthly`
* `./build/darwin/amd64/earthly`
* `./build/darwin/arm64/earthly`

It also builds the buildkitd image.

The buildkitd image is tagged with your current branch name and also the built binary defaults to using that built image. The built binary will always check on startup whether it has the latest buildkitd running for its configured image name and will restart buildkitd automatically to update. If during your development you end up making changes to just the buildkitd image, the binary will pick up the change on its next run.

For development purposes, you may use the built `earthly` binary to rebuild itself. It's usually faster than switching between the built binary and the prerelease binary because it avoids constant buildkitd restarts. After the first initial build, you'll end up using:


* Linux and WSL
    ```bash
    ./build/linux/amd64/earthly +for-linux
    ```
* Mac
    ```bash
    ./build/darwin/amd64/earthly +for-darwin
    ```

* Mac with M1 chip
    ```bash
    ./build/darwin/amd64/earthly +for-darwin-m1
    ```

## Delve

To use the [delve debugger](https://github.com/go-delve/delve) with the earthly binary, you need to disable optimizations in the 'go build' command. This is done using the -GO_GCFLAGS arg:

```
./earthly +for-own -GO_GCFLAGS='all=-N -l'
```

From there, you may use `dlv exec` against the binary, using `--` to separate dlv args from earthly args:

```
dlv exec ./build/own/earthly -- +base
Type 'help' for list of commands.

(dlv) break /earthly/earthfile2llb/interpreter.go:670
Breakpoint 1 set at 0x182866a for github.com/earthly/earthly/earthfile2llb.(*Interpreter).handleRun() /earthly/earthfile2llb/interpreter.go:670
(dlv) continue
 Init 🚀
————————————————————————————————————————————————————————————————————————————————

           buildkitd | Found buildkit daemon as podman container (earthly-dev-buildkitd)

 Build 🔧
————————————————————————————————————————————————————————————————————————————————

golang:1.20-alpine3.17 | --> Load metadata golang:1.20-alpine3.17 linux/amd64
> github.com/earthly/earthly/earthfile2llb.(*Interpreter).handleRun() /earthly/earthfile2llb/interpreter.go:670 (hits goroutine(295):1 total:1) (PC: 0x182866a)
(dlv)
```

## Running tests

To run most tests you can issue

```bash
./build/*/*/earthly -P \
  --secret DOCKERHUB_USER=<my-docker-username> \
  --secret DOCKERHUB_PASS=<my-docker-password-or-token> \
  +test --DOCKERHUB_AUTH=true
```

To also build the examples, you can run

```bash
./build/*/*/earthly -P  \
  --secret DOCKERHUB_USER=<my-docker-username> \
  --secret DOCKERHUB_PASS=<my-docker-password-or-token> \
  +test-all --DOCKERHUB_AUTH=true
```

It is also possible to run tests without credentials. But running all of them, or running too frequently may incur rate limits. You could run a single test, without credentials like this:

```bash
./build/*/*/earthly -P ./tests+env-test
```

If you don't want to specify these directly on the CLI, or don't want to type these each time, it's possible to store them in [.arg and .secret files](https://docs.earthly.dev/docs/earthly-command#build-args) instead.
Here is a template to get you started:

```shell
# .arg file
DOCKERHUB_AUTH=true
```

```shell
# .secret file
DOCKERHUB_USER=<my-docker-username>
DOCKERHUB_PASS=<my-docker-password-or-token>
```

### Running tests with an insecure pull through cache

The [Insecure Docker Hub Cache Example](https://docs.earthly.dev/ci-integration/pull-through-cache#insecure-docker-hub-cache-example), provides a guide on both
running an insecure docker hub pull through cache as well as configuring Earthly to use that cache.

Since Earthly uses itself for running the tests (Earthly-in-Earthly), simply configuring `~/.earthly/config.yml`[^dir] is insufficient -- one must also configure the
embedded version of Earthly to use the cache via build-args:

```bash
./build/*/*/earthly -P ./tests+all --DOCKERHUB_MIRROR=<ip-address-or-hostname>:<port> --DOCKERHUB_MIRROR_INSECURE=true
```

or if you are using a plain http cache, use:

```bash
./build/*/*/earthly -P ./tests+all --DOCKERHUB_MIRROR=<ip-address-or-hostname>:<port> --DOCKERHUB_MIRROR_HTTP=true
```

### Running tests with a mirror that requires authentication

To use a mirror that requires authentication, you can run:

```bash
./build/*/*/earthly -P \
  --secret DOCKERHUB_MIRROR_USER=<my-mirror-username> \
  --secret DOCKERHUB_MIRROR_PASS=<my-mirror-password> \
  ./tests+all --DOCKERHUB_MIRROR=<ip-address-or-hostname>:<port> --DOCKERHUB_MIRROR_AUTH=true
```

You can alternatively store these settings in the `.arg` and `.secret` files:

```shell
# .arg file
DOCKERHUB_MIRROR=<ip-address-or-hostname>:<port>
DOCKERHUB_MIRROR_AUTH=true
```

```shell
# .secret file
DOCKERHUB_MIRROR_USER=<my-mirror-username>
DOCKERHUB_MIRROR_PASS=<my-mirror-password>
```

### Running tests using earthly's internal mirror (only for members of the earthly org)

If you have access to `earthly-technologies/core`, you can make use of the internal mirror by running:

```bash
./build/*/*/earthly -P \
  ./tests+all --DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS=true
```

which will use the credentials which are stored in earthly's [cloud-hosted secrets](https://docs.earthly.dev/earthly-cloud/cloud-secrets).


## Updates to buildkit or fsutil

Earthly is built against a fork of [buildkit](https://github.com/earthly/buildkit) and [fsutil](https://github.com/earthly/fsutil).

To work with changes to this fork, you can use `earthly +for-linux --BUILDKIT_PROJECT=../buildkit`. This will use the local directory `../buildkit` for the buildkit code, when using buildkit in both `go.mod` and when building the buildkitd image.

For contributions that require updates to these forks, a PR must be opened in in the earthly-fork of the repository, and a corresponding PR should
be opened in the earthly repository -- please link the two PRs together, in order to show that earthly's tests will continue to pass with the changes to buildkit or fsutil.

The linked-PRs should be merged at the same time, in order to prevent earthly's main branch from pointing to a non-main branch of buildkit.
This is because the buildkit tests in the earthly fork of buildkit may not all pass -- this is a tech-debt trade-off -- instead of fixing (and extending these tests),
we instead rely on the earthly integration tests to pass before merging in changes to our fork.

The earthly-fork of the buildkit repository does not automatically squash commits; if you are submitting a PR for a new feature, it must be squashed manually.
The buildkit github repo is only setup to explicitly create a new merge commit for all merges -- this means you will have to go update your PR in earthly with a reference to
the new merge-commit (and wait for tests to run again); however, you may do a `git merge --ff-only <branch> && git push` using git on the command line, which will speed up the merge
process, since you will not have to edit your linked earthly PR (since a fast-forward change will not create a new merge commit, therefore allowing you to keep the existing referenced git sha in the
earthly PR).

If on the otherhand, you are pulling in upstream changes from moby, they should never be squashed.

To update earthly's reference to buildkit, you may run `earthly +update-buildkit --BUILDKIT_GIT_ORG=<git-user-or-org> --BUILDKIT_GIT_SHA=<40-char-git-reference-here>`.

Updates to fsutil must first be vendored into buildkit, then updated under `go.mod`; additional docs and scripts exist in the buildkit repo.


## Running buildkit under debug mode

Buildkit's scheduler has a debug mode, which can be enabled with the following `~/.earthly/config.yml`[^dir] config:

```yml
global:
  buildkit_additional_args: [ '-e', 'BUILDKIT_SCHEDULER_DEBUG=1' ]
```

then run `earthly --debug +target`.

This will produce scheduler debug messages such as

```
time="2022-10-27T18:18:06Z" level=debug msg="<< unpark [eyJzbCI6eyJmaWxlIjoiRWFydGhmaWxlIiwic3RhcnRMaW5lIjoyMSwic3RhcnRDb2x1bW4iOjQsImVuZExpbmUiOjIxLCJlbmRDb2x1bW4iOjI1fSwidGlkIjoiOTEyMWZkNzYtYjI5MS00YmQyLTg2MGUtNTZhYzJjZDVhMmY3IiwidG5tIjoiK3NsZWVwIiwicGx0IjoibGludXgvYW1kNjQifQ==] RUN --no-cache sleep 123\n"
time="2022-10-27T18:18:06Z" level=debug msg="> creating jzaxegge8eh5hjqe33jybv2ml [/bin/sh -c EARTHLY_LOCALLY=false PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin /usr/bin/earth_debugger /bin/sh -c 'sleep 123']" span="[eyJzbCI6eyJmaWxlIjoiRWFydGhmaWxlIiwic3RhcnRMaW5lIjoyMSwic3RhcnRDb2x1bW4iOjQsImVuZExpbmUiOjIxLCJlbmRDb2x1bW4iOjI1fSwidGlkIjoiOTEyMWZkNzYtYjI5MS00YmQyLTg2MGUtNTZhYzJjZDVhMmY3IiwidG5tIjoiK3NsZWVwIiwicGx0IjoibGludXgvYW1kNjQifQ==] RUN --no-cache sleep 123"
```

## Gotchas

### Auth

If you have issues with git-related features or with private docker registries, make sure you have configured auth correctly. See the [auth page](https://docs.earthly.dev/guides/auth) for more details.

You may need to adjust the docker login command in the `earthly-integration-test-base:` target by removing the Earthly repository and adjusting for your login credentials provider.

### Documentation

We maintain three different branches for [0.6](https://github.com/earthly/earthly/tree/docs-0.6), [0.7](https://github.com/earthly/earthly/tree/docs-0.7) and [0.8](https://github.com/earthly/earthly/tree/docs-0.8) docs, which are automatically propagated to [docs.earthly.dev](https://docs.earthly.dev/) (which has a dropdown options to switch between versions).

Documentation related to new unreleased features should be submitted in a PR to `main`, which we will merge into the `docs-0.8` branch when we perform a release.

To contribute improvements to documentation related to currently released features, please open a PR against the `docs-0.8` branch; we will cherry-pick these changes to the older documentation branches if necessary.

### Config

Starting with [v0.6.30](CHANGELOG.md#v0630---2022-11-22), the default location of the built binary's config file has
changed to `~/.earthly-dev/config.yml`. The standard location is not used as a fallback; it is possible to `export EARTHLY_CONFIG=~/.earthly/config.yml`, or create a symlink if required.

## Prereleases

In addition to the `./earthly` prerelease script, we maintain a repository dedicated to [prereleases versions](https://github.com/earthly/earthly-staging/releases) of earthly.

The prerelease versions follow a pseudo-semantic versioning scheme: `0.<epoch>.<decimal-git-sha>`; which is described in greater detail in the repository's [README](https://github.com/earthly/earthly-staging).

Additionally, prerelease docker images are pushed to [earthly/earthly-staging](https://hub.docker.com/r/earthly/earthly-staging/tags) and [earthly/buildkitd-staging](https://hub.docker.com/r/earthly/buildkitd-staging/tags).

## CLA

### Individual

All contributions must indicate agreement to the [Earthly Contributor License Agreement](https://gist.github.com/vladaionescu/ed990fa149a38a53ac74b64155bc6766) by logging into GitHub via the CLA assistant and signing the provided CLA. The CLA assistant will automatically notify the PRs that require CLA signing.

### Entity

If you are an entity, please use the [Earthly Contributor License Agreement form](https://earthly.dev/cla-form) in addition to requiring your individual contributors to sign all contributions.

[^dir]: Depending on the point in time earthly is being built the actual location [may be different](#config)
