# Contributing

## Using Earthly pre-release

To build Earthly from source, you need the same requirements as Earthly. We recommend that you use the pre-prelease version of Earthly for development purposes. To launch the pre-release Earthly, simply use the `./earth` script provided in the root of the earthly repository. The pre-release Earthly tracks the version on master. You can use `./earth --version` to identify which Git hash was used to build it.

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

### Running a full build stops after `SUCCESS`, but does not exit

After starting a run of `earth +all` the console displays a green line which states

```
=========================== SUCCESS ===========================
```

It does not denote the end of the build process however and some additional work might still be required. However it might take a while to be displayed.
