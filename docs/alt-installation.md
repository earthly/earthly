# Alternative Installation

This page outlines alternative installation instructions for the `earthly` build tool. The main instructions that most users need are available on the [installation intructions page](https://earthly.dev/get-earthly).

## Pre-requisites

* [Docker](https://docs.docker.com/install/)
* [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* (*Windows only*) [Docker WSL 2 backend](https://docs.docker.com/docker-for-windows/wsl/)

## Install earthly

Download the binary relevant to your platform from [the releases page](https://github.com/earthly/earthly/releases), rename it to `earthly` and place it in your `bin`.

To initialize auto-completion, run

```bash
earthly bootstrap
```

and then restart your shell.

### CI

For instructions on how to install `earthly` for CI use, see the [CI integration guide](./ci-integration/ci-integration.md).

### Installing from source

To install from source, see the [contributing page](https://github.com/earthly/earthly/blob/main/CONTRIBUTING.md).

## Configuration

If you use SSH-based git authentication, then your git credentials will just work with Earthly. Read more about [git auth](./guides/auth.md).

For a full list of configuration options, see the [Configuration reference](./earthly-config/earthly-config.md)

## Verify installation

To verify that the installation works correctly, you can issue a simple build of an existing hello-world project

```bash
earthly github.com/earthly/hello-world:main+hello
```

You should see the output

```
github.com/earthly/hello-world:main+hello | --> RUN [echo 'Hello, world!']
github.com/earthly/hello-world:main+hello | Hello, world!
github.com/earthly/hello-world:main+hello | Target github.com/earthly/hello-world:main+hello built successfully
=========================== SUCCESS ===========================
```
