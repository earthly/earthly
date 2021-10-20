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

This builds the earthly binary in `./build/<platform>/amd64/earthly` and also the buildkitd image.

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

## Running tests

To run most tests you can issue

```bash
./build/<platform>/amd64/earthly -P \
  --secret DOCKERHUB_USER=<my-docker-username> \
  --secret DOCKERHUB_TOKEN=<my-docker-token> \
  +test
```

To also build the examples, you can run

```bash
./build/<platform>/amd64/earthly -P  \
  --secret DOCKERHUB_USER=<my-docker-username> \
  --secret DOCKERHUB_TOKEN=<my-docker-token> \
  +test-all
```

The token should be the same token you use to login with Docker Hub. Other repositories are not supported. It is also possible to run tests without credentials. But running all of them, or running too frequently may incur rate limits. You could run a single test, without credentials like this:

```bash
./build/<platform>/amd64/earthly -P --build-arg DOCKERHUB_AUTH=false ./tests+env-test
```

If you don't want to specify these directly on the CLI, or don't want to type these each time, its possible [to use an .env file instead](https://docs.earthly.dev/docs/earthly-command#environment-variables-and-.env-file). Here is a template to get you started:

```shell
DOCKERHUB_USER=<my-docker-username>
DOCKERHUB_TOKEN=<my-docker-token>
```

## Gotchas

### Auth

If you have issues with git-related features or with private docker registries, make sure you have configured auth correctly. See the [auth page](https://docs.earthly.dev/guides/auth) for more details.

## CLA

### Individual

All contributions must indicate agreement to the [Earthly Contributor License Agreement](https://gist.github.com/vladaionescu/ed990fa149a38a53ac74b64155bc6766) by logging into GitHub via the CLA assistant and signing the provided CLA. The CLA assistant will automatically notify the PRs that require CLA signing.

### Entity

If you are an entity, please use the [Earthly Contributor License Agreement form](https://earthly.dev/cla-form) in addition to requiring your individual contributors to sign all contributions.
