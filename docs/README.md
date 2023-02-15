# Introduction

Earthly is a versatile, approachable CI/CD framework that gives you repeatable builds by running each build inside containers *– write builds once, and they will run anywhere*. Earthly has a super simple, instantly recognizable syntax that is easy to write and understand *– like Dockerfile and Makefile had a baby*. And it leverages and augments popular build tools instead of replacing them, *so you don’t have to rewrite all your builds no matter what languages you use*. With Earthly, you can create Docker images and build artifacts (e.g. binaries, packages, and arbitrary files).

Earthly can run locally or on top of popular CI systems – such as [Jenkins](./ci-integration/guides/jenkins.md), [CircleCI](./ci-integration/guides/circle-integration.md), [GitHub Actions](./ci-integration/guides/gh-actions-integration.md), [AWS CodeBuild](./ci-integration/guides/codebuild-integration.md), [Google Cloud Build](./ci-integration/guides/google-cloud-build.md), and [GitLab CI/CD](./ci-integration/guides/gitlab-integration.md)). It typically acts as the layer between language-specific tooling (such as maven, gradle, npm, pip, and go build) and the CI build spec.

![Earthly fits between language-specific tooling and the CI](img/integration-diagram-v2.png)

Earthly's key features/benefits are:
  * **🔁 Repeatable Builds**
    Earthly runs all builds in containers, making them self-contained, isolated, repeatable, and portable. When you write a build, you know it will execute correctly no matter where it runs – your laptop, a colleague’s laptop, or any CI. You don’t have to configure language-specific tooling, install additional dependencies, or complicate your build scripts to ensure they are compatible with different OSs. Earthly gives you consistent, repeatable builds regardless of where they run.
  * **❤️ Super Simple**  
    Earthly’s syntax is easy to write and understand. Most engineers can read an Earthfile instantly, without prior knowledge of Earthly. We combined some of the best ideas from Dockerfiles and Makefiles into one specification *– like Dockerfile and Makefile had a baby*.
  * **🛠 Compatible with Every Language, Framework, and Build Tool**  
    One of the key principles of Earthly is that the best build tooling for a specific language is built by the community of that language itself. Earthly does not intend to replace any language-specific build tooling, but rather to leverage and augment them. Earthly works with the compilers and build tools you use. If it runs on Linux, it runs on Earthly. And you don’t have to rewrite your existing builds or replace your `package.json`, `go.mod`, `build.gradle`, or `Cargo.toml` files. You can use Earthly as a wrapper around your existing tooling and still get Earthly’s repeatable builds, parallel execution, and build caching.
  * **🏘 Great for Monorepos and Polyrepos**  
    Earthly is great for both [monorepos](https://github.com/earthly/earthly/tree/main/examples/monorepo) and [polyrepos](https://github.com/earthly/earthly/tree/main/examples/multirepo). You can split your build logic across multiple Earthfiles, placing some deeper inside the directory structure or even in other repositories. Referencing targets from other Earthfiles is easy regardless of where they are stored. So you can organize your build logic however makes the most sense for your project.
  * **💨 Fast Builds**  
    Earthly automatically executes build targets in parallel and makes maximum use of cache. This makes builds fast. Earthly also has powerful shared caching capabilities that speed up builds frequently run across a team or in sandboxed environments, such as Earthly Satellites, GitHub Actions, or your CI.  
    &nbsp;  
    If your build has multiple steps, Earthly will:
    1. Build a directed acyclic graph (DAG).
    2. Isolate execution of each step.
    3. Run independent steps in parallel.
    4. Cache results for future use.
  * **♻️ Reuse, Don't Repeat**  
    Never have to write the same code in multiple builds again. With Earthly, you can reuse targets, artifacts, and images across multiple Earthfiles, even ones in other repositories, in a single line. Earthly is cache-aware, based on the individual hashes of each file, and has shared caching capabilities. So you can create a vast and efficient build hierarchy that only executes the minimum required steps.

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
