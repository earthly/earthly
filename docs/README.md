# Introduction

Earthly is a build automation tool from the same era as your code. It allows you to execute all your builds in containers. This makes them self-contained, repeatable, portable and parallel. You can use Earthly to create Docker images and artifacts (eg binaries, packages, arbitrary files).

Earthly can run on top of popular CI systems (like Jenkins, [Circle](./examples/circle-integration.md), [GitHub Actions](./examples/gh-actions-integration.md), [AWS CodeBuild](./examples/codebuild-integration.md)). It is typically the layer between language-specific tooling (like maven, gradle, npm, pip, go build) and the CI build spec.

![Earthly fits between language-specific tooling and the CI](img/integration-diagram.png)

Earthly has a number of key features. It has a familiar syntax (it's like Dockerfile and Makefile had a baby). Everything runs on containers, so your builds run the same on your laptop as they run in CI or on your colleague's laptop. Strong isolation also gives you easy to use parallelism, with no strings attached. You can also import dependencies from other directories or other repositories with ease, making Earthly great for large [mono-repo builds](./examples/monorepo.md) that span a vast directory hierarchy; but also for [multi-repo setups](./examples/multirepo.md) where builds might depend on each other across repositories.

One of the key principles of Earthly is that the best build tooling of a specific language is built by the community of that language itself. Earthly does not intend to replace that tooling, but rather to leverage and augment it.

## Installation

For a full list of installation options see the [Installation page](./installation/installation.md).

### Linux

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly && /usr/local/bin/earthly bootstrap'
```

### Mac

```bash
brew install earthly
```

### Windows via WSL (**beta**)

Earthly on Windows requires [Docker Desktop WSL2 backend](https://docs.docker.com/docker-for-windows/wsl/). Under `wsl`, run the following to install `earthly`.

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly && /usr/local/bin/earthly bootstrap'
```

### Syntax highlighting

* [VS Code](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting)
* [Vim](https://github.com/earthly/earthly.vim)
* [Sublime Text](https://packagecontrol.io/packages/Earthly%20Earthfile)

## Getting started

If you are new to Earthly, check out the [Basics page](./guides/basics.md), to get started.

A high-level overview is available on [the Earthly GitHub page](https://github.com/earthly/earthly).

## Quick Links

* [Earthly GitHub page](https://github.com/earthly/earthly)
* [Full installation instructions](./installation/installation.md)
* [Earthly basics](./guides/basics.md)
* [Earthfile reference](./earthfile/earthfile.md)
* [Earthly command reference](./earthly-command/earthly-command.md)
* [Configuration reference](./earthly-config/earthly-config.md)
* [Earthfile examples](./examples/examples.md)
