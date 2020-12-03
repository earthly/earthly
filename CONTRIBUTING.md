# Contributing

## Code of Conduct

Please refer to [code-of-conduct.md](./code-of-conduct.md) for details.

## Using Earthly pre-release

To build Earthly from source, you need the same requirements as Earthly. We recommend that you use the pre-prelease version of Earthly for development purposes. To launch the pre-release Earthly, simply use the `./earth` script provided in the root of the earthly repository. The pre-release Earthly tracks the version on main. You can use `./earth --version` to identify which Git hash was used to build it.

## Building from source

To build Earthly from source for your target system, use

* Linux and WSL
    ```bash
    ./earth +for-linux
    ```
* Mac
    ```bash
    ./earth +for-darwin
    ```

This builds the earth binary in `./build/<platform>/amd64/earth` and also the buildkitd image.

The buildkitd image is tagged with your current branch name and also the built binary defaults to using that built image. The built binary will always check on startup whether it has the latest buildkitd running for its configured image name and will restart buildkitd automatically to update. If during your development you end up making changes to just the buildkitd image, the binary will pick up the change on its next run.

For development purposes, you may use the built `earth` binary to rebuild itself. It's usually faster than switching between the built binary and the prerelease binary because it avoids constant buildkitd restarts. After the first initial build, you'll end up using:


* Linux and WSL
    ```bash
    ./build/linux/amd64/earth +for-linux
    ```
* Mac
    ```bash
    ./build/darwin/amd64/earth +for-darwin
    ```

## Running tests

To run most tests you can issue

```bash
./build/<platform>/amd64/earth -P +test
```

To also build the examples, you can run

```bash
./build/<platform>/amd64/earth -P +test-all
```

## Gotchas

### Auth

If you have issues with git-related features or with private docker registries, make sure you have configured auth correctly. See the [auth page](https://docs.earthly.dev/guides/auth) for more details.

## CLA

### Individual

All contributions must indicate agreement to the [Earthly Contributor License Agreement](https://earthly.dev/cla) by signing off all commits. The sign-off is a simple line at the end of the explanation for the patch.

```
Signed-off-by: John Doe <john.doe@email.com>
```

Using your real name is a requirement. If you set your `user.name` and `user.email` git configs, you can sign your commit automatically with

```
git commit -s
```

### Entity

If you are an entity, please use the [Earthly Contributor License Agreement form](https://earthly.dev/cla-form) in addition to requiring individual contributors to sign off all contributions.
