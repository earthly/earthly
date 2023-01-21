<h1 align="center"><a href="https://earthly.dev"><img src="img/logo-banner-white-bg.png" alt="Earthly" align="center" width="700px" /></a></h1>

[![GitHub Actions CI](https://github.com/earthly/earthly/workflows/GitHub%20Actions%20CI/badge.svg)](https://github.com/earthly/earthly/actions?query=workflow%3A%22GitHub+Actions+CI%22+branch%3Amain)
[![Join the chat on Slack](https://img.shields.io/badge/slack-join%20chat-red.svg)](https://earthly.dev/slack)
[![Docs](https://img.shields.io/badge/docs-git%20book-blue)](https://docs.earthly.dev)
[![Website](https://img.shields.io/badge/website-earthly.dev-blue)](https://earthly.dev)
[![Install Earthly](https://img.shields.io/github/v/release/earthly/earthly.svg?label=install&color=1f626c)](https://earthly.dev/get-earthly)
[![Docker Hub](https://img.shields.io/badge/docker%20hub-earthly-blue)](https://hub.docker.com/u/earthly)
[![License MPL-2](https://img.shields.io/badge/license-MPL-blue.svg)](./LICENSE)

**🐳 Build anything via containers** - *build images or standalone artifacts (binaries, packages, arbitrary files)*

**🛠 Programming language agnostic** - *allows the use of language-specific build tooling*

**🔁 Repeatable builds** - *does not depend on user's local installation: runs the same locally, as in CI*

**⛓ Parallelism that just works** - *build in parallel without special considerations*

**🏘 Mono and Poly-repo friendly** - *ability to split the build definitions across vast project hierarchies*

**💾 Shared caching** - *share build cache between CI runners*

**🔀 Multi-platform** - *build for multiple platforms in parallel*

---------------------------------

[🌍 Earthly](https://earthly.dev) is a CI/CD framework that allows you to develop pipelines locally and run them anywhere. Earthly leverages containers for the execution of pipelines. This makes them self-contained, repeatable, portable and parallel.

<br/>
<div align="center"><a href="https://earthly.dev/get-earthly"><img src="docs/img/get-earthly-button.png" alt="Get Earthly" title="Get Earthly" /></a></div>
<br/>

---------------------------------

<br/>
<br/>
<h2 align="center">Table of Contents</h2>

* [Why use Earthly?](#why-use-earthly)
* [Where Does Earthly Fit?](#where-does-earthly-fit)
* [How Does It Work?](#how-does-it-work)
* [Installation](#installation)
* [Quick Start](#quick-start)
* [Features](#features)
* [FAQ](#faq)
* [Contributing](#contributing)
* [Licensing](#licensing)

<br/>
<br/>
<h2 align="center">Why Use Earthly?</h2>

### 🔁 Reproduce CI failures

Earthly builds are self-contained, isolated and repeatable. Regardless of whether Earthly runs in your CI or on your laptop, there is a degree of guarantee that the build will run the same way. This allows for faster iteration on the build scripts and easier debugging when something goes wrong. No more `git commit -m "try again"`.

### 🤲 Builds that run the same for everyone

Repeatable builds also mean that your build will run the same on your colleagues' laptop without any additional project-specific or language-specific setup. This fosters better developer collaboration and mitigates works-for-me type of issues.

### 🚀 From zero to working build in minutes

Jump from project to project with ease, regardless of the language they are written in. Running the project's test suites is simply a matter of running an Earthly target (without fiddling with project configuration to make it compile and run on your system). Contribute across teams with confidence.

### 📦 Reusability

A simple, yet powerful import system allows for reusability of builds across directories or even across repositories. Importing other builds does not have hidden environment-specific implications - it just works.

### ❤️ It's like Makefile and Dockerfile had a baby

Taking some of the best ideas from Makefiles and Dockerfiles, Earthly combines two build specifications into one.

<br/>
<br/>
<h2 align="center">Where Does Earthly Fit?</h2>

<div align="center"><img src="docs/img/integration-diagram-v2.png" alt="Earthly fits between language-specific tooling and the CI" width="700px" /></div>
<br/>

Earthly is meant to be used both on your development machine and in CI. It can run on top of popular CI systems (like Jenkins, [Circle](https://docs.earthly.dev/examples/circle-integration), [GitHub Actions](https://docs.earthly.dev/examples/gh-actions-integration)). It is typically the layer between language-specific tooling (like maven, gradle, npm, pip, go build) and the CI build spec.

<br/>
<br/>
<h2 align="center">How Does It Work?</h2>

In short: **containers**, **layer caching** and **complex build graphs**!

Earthly executes builds in containers, where execution is isolated. The dependencies of the build are explicitly specified in the build definition, thus making the build self-sufficient.

We use a target-based system to help users break up complex builds into reusable parts. Nothing is shared between targets other than clearly declared dependencies. Nothing shared means no unexpected race conditions. In fact, the build is executed in parallel whenever possible, without any need for the user to take care of any locking or unexpected environment interactions.

| ℹ️ Note <br/><br/> Earthfiles might seem very similar to Dockerfile multi-stage builds. In fact, the [same technology](https://github.com/moby/buildkit) is used underneath. However, a key difference is that Earthly is designed to be a general-purpose build system, not just a Docker image specification. Read more about [how Earthly is different from Dockerfiles](#how-is-earthly-different-from-dockerfiles). |
| :--- |

<br/>
<br/>
<h2 align="center">Installation</h2>

See [installation instructions](https://earthly.dev/get-earthly).

To build from source, check the [contributing page](./CONTRIBUTING.md).

<br/>
<br/>
<h2 align="center">Quick Start</h2>

Here are some resources to get you started with Earthly

* 🏁 [Getting started guide](https://docs.earthly.dev/guides/basics)
* 👀 [Examples](./examples)
  * [C](./examples/c)
  * [C++](./examples/cpp)
  * [COBOL](./examples/cobol)
  * [Go](./examples/go)
  * [Java](./examples/java)
  * [JS](./examples/js)
  * [Python](./examples/python)
  * [Ruby](./examples/ruby)
  * [Rust](./examples/rust)
  * [Scala](./examples/scala)
  * [Mono-repo](./examples/monorepo)
  * [Multi-repo](./examples/multirepo)
* 🔍 Explore [Earthly's own build](https://docs.earthly.dev/examples/examples#earthlys-own-build)
* ✔️ [Best practices](https://docs.earthly.dev/best-practices)

See also the [full documentation](https://docs.earthly.dev).

Reference pages

* 📑 [Earthfile reference](https://docs.earthly.dev/earthfile)
* #️⃣ [Earthly command reference](https://docs.earthly.dev/earthly-command)
* ⚙️ [Configuration reference](https://docs.earthly.dev/earthly-config)

### A simple example (for Go)

```earthly
# Earthfile
VERSION 0.6
FROM golang:1.15-alpine3.13
RUN apk --update --no-cache add git
WORKDIR /go-example

all:
  BUILD +lint
  BUILD +docker

build:
  COPY main.go .
  RUN go build -o build/go-example main.go
  SAVE ARTIFACT build/go-example AS LOCAL build/go-example

lint:
  RUN go get golang.org/x/lint/golint
  COPY main.go .
  RUN golint -set_exit_status ./...

docker:
  COPY +build/go-example .
  ENTRYPOINT ["/go-example/go-example"]
  SAVE IMAGE go-example:latest
```

```go
// main.go
package main

import "fmt"

func main() {
  fmt.Println("hello world")
}
```

Invoke the build using `earthly +all`.

<div align="center"><a href="https://asciinema.org/a/351683?speed=2"><img src="img/demo-351683.gif" alt="Demonstration of a simple Earthly build" title="View on asciinema.org" width="600px" /></a></div>

Examples for other languages are available in the [examples dir](./examples).

<br/>
<br/>
<h2 align="center">Features</h2>

### 📦 Modern import system

Earthly can be used to reference and build targets from other directories or even other repositories. For example, if we wanted to build [an example target from the `github.com/earthly/earthly` repository](./examples/go/Earthfile#L17-L20), we could issue

```bash
# Try it yourself! No need to clone.
earthly github.com/earthly/earthly/examples/go:main+docker
# Run the resulting image.
docker run --rm earthly/examples:go
```

### 🔨 Reference other targets using +

Use `+` to reference other targets and create complex build inter-dependencies.

<div align="center"><a href="https://docs.earthly.dev/guides/target-ref"><img src="docs/guides/img/ref-infographic-v2.png" alt="Target and artifact reference syntax" title="Reference targets using +" width="600px" /></a></div>

Examples

* Same directory (same Earthfile)
  
  ```earthly
  BUILD +some-target
  FROM +some-target
  COPY +some-target/my-artifact ./
  ```

* Other directories

  ```earthly
  BUILD ./some/local/path+some-target
  FROM ./some/local/path+some-target
  COPY ./some/local/path+some-target/my-artifact ./
  ```

* Other repositories

  ```earthly
  BUILD github.com/someone/someproject:v1.2.3+some-target
  FROM github.com/someone/someproject:v1.2.3+some-target
  COPY github.com/someone/someproject:v1.2.3+some-target/my-artifact ./
  ```

### 💾 Caching that works the same as Docker builds

<div align="center"><a href="https://asciinema.org/a/351674?speed=2"><img src="img/demo-351674.gif" alt="Demonstration of Earthly's caching" title="View on asciinema.org" width="600px" /></a></div>

Cut down build times in CI through [Shared Caching](https://docs.earthly.dev/guides/shared-cache).

### 🛠 Multi-platform support

Build for multiple platforms in parallel.

```earthly
VERSION 0.6
all:
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        --platform=linux/arm/v7 \
        --platform=linux/arm/v6 \
        +build

build:
    FROM alpine:3.15
    CMD ["uname", "-m"]
    SAVE IMAGE multiplatform-image
```

### ⛓ Parallelization that just works

Whenever possible, Earthly automatically executes targets in parallel.

<div align="center"><a href="https://asciinema.org/a/351678?speed=2"><img src="img/demo-351678.gif" alt="Demonstration of Earthly's parallelization" title="View on asciinema.org" width="600px" /></a></div>

### 🤲 Make use of build tools that work everywhere

No need to ask your team to install `protoc`, a specific version of Python, Java 1.6 or the .NET Core ecosystem. You only install once, in your Earthfile, and it works for everyone. Or even better, you can just make use of the rich Docker Hub ecosystem.

```earthly
VERSION 0.6
FROM golang:1.15-alpine3.13
WORKDIR /proto-example

proto:
  FROM namely/protoc-all:1.29_4
  COPY api.proto /defs
  RUN --entrypoint -- -f api.proto -l go
  SAVE ARTIFACT ./gen/pb-go /pb AS LOCAL pb

build:
  COPY go.mod go.sum .
  RUN go mod download
  COPY +proto/pb pb
  COPY main.go ./
  RUN go build -o build/proto-example main.go
  SAVE ARTIFACT build/proto-example
```

See full [example code](./examples/readme/proto).

### 🔑 Cloud secrets support built-in

Secrets are never stored within an image's layers and they are only available to the commands that need them.

```bash
earthly set /user/github/token 'shhh...'
```

```earthly
release:
  RUN --push --secret GITHUB_TOKEN=user/github/token github-release upload file.bin
```

<br/>
<br/>
<h2 align="center">FAQ</h2>

### How is Earthly different from Dockerfiles?

[Dockerfiles](https://docs.docker.com/engine/reference/builder/) were designed for specifying the make-up of Docker images and that's where Dockerfiles stop. Earthly takes some key principles of Dockerfiles (like layer caching), but expands on the use-cases. For example, Earthly can output regular artifacts, run unit and integration tests, and create several Docker images at a time - all outside the scope of Dockerfiles.

It is possible to use Dockerfiles in combination with other technologies (e.g., Makefiles or bash files) to solve such use-cases. However, these combinations are difficult to parallelize, challenging to scale across repositories as they lack a robust import system and also they often vary in style from one team to another. Earthly does not have these limitations as it was designed as a general-purpose build system.

For example, Earthly introduces a richer target, artifact and image [referencing system](https://docs.earthly.dev/guides/target-ref), allowing for better reuse in complex builds spanning a single large repository or multiple repositories. Because Dockerfiles are only meant to describe one image at a time, such features are outside the scope of applicability of Dockerfiles.

### How do I tell apart classical Dockerfile commands from Earthly commands?

Check out the [Earthfile reference doc page](https://docs.earthly.dev/earthfile). It has all the commands there and specifies which commands are the same as Dockerfile commands and which are new.

### Can Earthly build Dockerfiles?

Yes! You can use the command `FROM DOCKERFILE` to inherit the commands in an existing Dockerfile.

```earthly
build:
  FROM DOCKERFILE .
  SAVE IMAGE some-image:latest
```

You may also optionally port your Dockerfiles to Earthly entirely. Translating Dockerfiles to Earthfiles is usually a matter of copy-pasting and making minor adjustments. See the [getting started page](https://docs.earthly.dev/guides/basics) for some Earthfile examples.

### How is Earthly different from Bazel?

[Bazel](https://bazel.build) is a build tool developed by Google to optimize the speed, correctness, and reproducibility of their internal monorepo codebase. The main difference between Bazel and Earthly is that Bazel is a **build system**, whereas Earthly is a **general-purpose CI/CD framework**. For a more in-depth explanation see [our FAQ](https://earthly.dev/faq#bazel).

<br/>
<br/>
<h2 align="center">Contributing</h2>

* Please report bugs as [GitHub issues](https://github.com/earthly/earthly/issues).
* Join us on [Slack](https://earthly.dev/slack)!
* Questions via GitHub issues are welcome!
* PRs welcome! But please give a heads-up in a GitHub issue before starting work. If there is no GitHub issue for what you want to do, please create one.
* To build from source, check the [contributing page](./CONTRIBUTING.md).

<br/>
<br/>
<h2 align="center">Licensing</h2>

Earthly is licensed under the Mozilla Public License Version 2.0. See [LICENSE](./LICENSE).
