# Introduction

Earthly is a build automation tool from the same era as your code. It allows you to execute all your builds in containers. This makes them self-contained, repeatable, portable and parallel. You can use Earthly to create Docker images and artifacts (e.g. binaries, packages, arbitrary files).

Earthly can run on top of popular CI systems (like [Jenkins](./ci-integration/guides/jenkins.md), [CircleCI](./ci-integration/guides/circle-integration.md), [GitHub Actions](./ci-integration/guides/gh-actions-integration.md), [AWS CodeBuild](./ci-integration/guides/codebuild-integration.md)). It is typically the layer between language-specific tooling (like maven, gradle, npm, pip, go build) and the CI build spec.

![Earthly fits between language-specific tooling and the CI](img/integration-diagram-v2.png)

Earthly has a number of key features. It has a familiar syntax (it's like Dockerfile and Makefile had a baby). Everything runs on containers, so your builds run the same on your laptop as they run in CI or on your colleague's laptop. Strong isolation also gives you easy to use parallelism, with no strings attached. You can also import dependencies from other directories or other repositories with ease, making Earthly great for large [mono-repo builds](https://github.com/earthly/earthly/tree/main/examples/monorepo) that span a vast directory hierarchy; but also for [multi-repo setups](https://github.com/earthly/earthly/tree/main/examples/multirepo) where builds might depend on each other across repositories.

One of the key principles of Earthly is that the best build tooling of a specific language is built by the community of that language itself. Earthly does not intend to replace that tooling, but rather to leverage and augment it.

## Installation

See [installation instructions](https://earthly.dev/get-earthly).

For a full list of installation options see the [alternative installation page](./alt-installation/alt-installation.md).

## Getting started

If you are new to Earthly, check out the [Basics page](./basics/basics.md), to get started.

A high-level overview is available on [the Earthly GitHub page](https://github.com/earthly/earthly).

## Quick Links

* [Earthly GitHub page](https://github.com/earthly/earthly)
* [Installation instructions](https://earthly.dev/get-earthly)
* [Earthly basics](./basics/basics.md)
* [Earthfile reference](./earthfile/earthfile.md)
* [Earthly command reference](./earthly-command/earthly-command.md)
* [Configuration reference](./earthly-config/earthly-config.md)
* [Earthfile examples](./examples/examples.md)
* [Best practices](./best-practices/best-practices.md)
