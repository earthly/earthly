# Introduction

Earthly is a build automation tool based on containers. It allows you to execute self-contained and portable builds, to create Docker images and artifacts (eg binaries, packages, arbitrary files), with the help of container technology.

It's like Dockerfile and Makefile had a baby. And they called it Earthly.

Earthly has a number of key features. You can specify builds with a Dockerfile-like syntax. Because everything runs on containers, your builds run the same on your laptop as they run in CI or on your colleague's laptop. Strong isolation also gives you easy to use parallelism, with no strings attached. You can also import dependencies from other directories or other repositories with ease, making Earthly great for large mono-repo builds that span a vast directory hierarchy; but also for multi-repo setups where builds might depend on each other across repositories.

## Installation

To install `earth` (the Earthly CLI) on your system, see [instructions on the Earthly GitHub page](https://github.com/vladaionescu/earthly#installation).

You may optionally also install the [VS Code Syntax Highlighting extension](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting).

## Getting started

If you are new to Earthly, check out the [Basics page](./guides/basics.md), to get started.

A high-level overview is available on [the Earthly GitHub page](https://github.com/vladaionescu/earthly).

## Quick Links

* [Earthly GitHub page](https://github.com/vladaionescu/earthly)
* [Earthly basics](./guides/basics.md)
* [Earthfile reference](./earthfile/earthfile.md)
* [earth command reference](./earth-command/earth-command.md)
* [Earthfile examples](./examples/examples.md)
