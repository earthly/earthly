<h1 align="center"><a href="https://earthly.dev"><img src="img/logo-banner-white-bg.png" alt="Earthly" align="center" width="700px" /></a></h1>

![CI](https://github.com/earthly/earthly/workflows/CI/badge.svg)
[![Join the chat at https://gitter.im/earthly-room/community](https://badges.gitter.im/earthly-room.svg)](https://gitter.im/earthly-room/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Docs](https://img.shields.io/badge/docs-git%20book-blue)](https://docs.earthly.dev)
[![Website](https://img.shields.io/badge/website-earthly.dev-blue)](https://earthly.dev)
[![Docker Hub](https://img.shields.io/badge/docker%20hub-earthly-blue)](https://hub.docker.com/u/earthly)
[![License](https://img.shields.io/badge/license-MPL--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
<a href="https://www.producthunt.com/posts/earthly-2?utm_source=badge-featured&utm_medium=badge&utm_souce=badge-earthly-2" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=266617&theme=dark" alt="Earthly - Build anything via containers | Product Hunt" width="93" height="20" /></a>

**🐳 Build anything via containers** - *build images or standalone artifacts (binaries, packages, arbitrary files)*

**🛠 Programming language agnostic** - *allows use of language-specific build tooling*

**🔁 Reproducible builds** - *does not depend on user's local installation. Runs the same locally, as in CI*

**⛓ Parallelism that just works** - *builds in parallel without special considerations the user has to make*

**🏠 Mono-repo friendly** - *ability to split the build definitions across a vast directory hierarchy*

**🏘 Multi-repo friendly** - *ability to import builds or artifacts from other repositories*

---------------------------------

[🌍 Earthly](https://earthly.dev) is a build automation tool for the post-container era. It allows you to execute all your builds in containers. This makes them self-contained, reproducible, portable and parallel. You can use Earthly to create Docker images and artifacts (eg binaries, packages, arbitrary files).

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

Earthly builds are self-contained, isolated and reproducible. Regardless of whether Earthly runs in your CI or on your laptop, there is a degree of guarantee that the build will run the same way. This allows for faster iteration on the build scripts and easier debugging when something goes wrong. No more `git commit -m "try again"`.

### 🤲 Builds that run the same for everyone

Reproducible builds also means that your build will run the same on your colleagues laptop without any additional project-specific or language-specific setup. This fosters better developer collaboration and mitigates works-for-me type of issues.

### 🚀 From zero to working build in minutes

Jump from project to project with ease, regardless of the language they are written in. Running the project's test suites is simply a matter of running an Earthly target (without fiddling with project configuration to make it compile and run on your system). Contribute across teams with confidence.

### 📦 Reusability

A simple, yet powerful import system allows for reusability of builds across directories or even across repositories. Importing other builds does not have hidden environment-specific implications - it just works.

### ❤️ It's like Makefile and Dockerfile had a baby

Taking some of the best ideas from Makefiles and Dockerfiles, Earthly combines two build specifications into one.

<br/>
<br/>
<h2 align="center">Where Does Earthly Fit?</h2>

<div align="center"><img src="docs/img/integration-diagram.png" alt="Earthly fits between language-specific tooling and the CI" width="700px" /></div>
<br/>

Earthly is meant to be used both on your development machine and in CI. It can run on top of popular CI systems (like Jenkins, [Circle](https://docs.earthly.dev/examples/circle-integration), [GitHub Actions](https://docs.earthly.dev/examples/gh-actions-integration)). It is typically the layer between language-specific tooling (like maven, gradle, npm, pip, go build) and the CI build spec.

<br/>
<br/>
<h2 align="center">How Does It Work?</h2>

In short: **containers**, **layer caching** and **complex build graphs**!

Earthly executes builds in containers, where execution is isolated. The dependencies of the build are explicitly specified in the build definition, thus making the build self-sufficient.

We use a target-based system to help users break-up complex builds into reusable parts. Nothing is shared between targets, other than clearly declared dependencies. Nothing shared means no unexpected race conditions. In fact, the build is executed in parallel whenever possible, without any need for the user to take care of any locking or unexpected environment interactions.

| ℹ️ Note <br/><br/> Earthfiles might seem very similar to Dockerfile multi-stage builds. In fact, the [same technology](https://github.com/moby/buildkit) is used underneath. However, a key difference is that Earthly is designed to be a general purpose build system, not just a Docker image specification. Read more about [how Earthly is different from Dockerfiles](#how-is-earthly-different-from-dockerfiles). |
| :--- |

<br/>
<br/>
<h2 align="center">Installation</h2>

For a full list of installation options see the [Installation page](https://docs.earthly.dev/installation).

### Linux

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth && /usr/local/bin/earth bootstrap'
```

### Mac

```bash
brew install earthly
earth bootstrap
```

### Windows via WSL (**beta**)
  
Earthly on Windows requires [Docker Desktop WSL2 backend](https://docs.docker.com/docker-for-windows/wsl/). Under `wsl`, run the following to install `earth`.

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth && /usr/local/bin/earth bootstrap'
```

### Your CI

See the [CI integration guide](https://docs.earthly.dev/guides/ci-integration)

### Syntax highlighting

[Install](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting) for VS code.

```
ext install earthly.earthfile-syntax-highlighting
```

<br/>
<br/>
<h2 align="center">Quick Start</h2>

Here are some resources to get you started with Earthly

* 🏁 [Getting started guide](https://docs.earthly.dev/guides/basics)
* 👀 [Examples](https://docs.earthly.dev/examples)
  * [Go](https://docs.earthly.dev/examples/go)
  * [Java](https://docs.earthly.dev/examples/java)
  * [JS](https://docs.earthly.dev/examples/js)
  * [C++](https://docs.earthly.dev/examples/cpp)
  * [Mono-repo](https://docs.earthly.dev/examples/monorepo)
  * [Multi-repo](https://docs.earthly.dev/examples/multirepo)
  * The [examples](./examples) dir
* 🔍 Explore [Earthly's own build](https://docs.earthly.dev/examples/earthly)

See also the [full documentation](https://docs.earthly.dev).

Reference pages

* 📑 [Earthfile reference](https://docs.earthly.dev/earthfile)
* #️⃣ [Earth command reference](https://docs.earthly.dev/earth-command)
* ⚙️ [Configuration reference](https://docs.earthly.dev/earth-config)

### A simple example (for Go)

```Dockerfile
# Earthfile
FROM golang:1.13-alpine3.11
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

Invoke the build using `earth +all`.

<div align="center"><a href="https://asciinema.org/a/351683?speed=2"><img src="img/demo-351683.gif" alt="Demonstration of a simple Earthly build" title="View on asciinema.org" width="600px" /></a></div>

Examples for other languages are available on the [examples page](https://docs.earthly.dev/examples).

<br/>
<br/>
<h2 align="center">Features</h2>

### 📦 Modern import system

Earthly can be used to reference and build targets from other directories or even other repositories. For example, if we wanted to build [an example target from the `github.com/earthly/earthly` repository](./examples/go/Earthfile#L17-L20), we could issue

```bash
# Try it yourself! No need to clone.
earth github.com/earthly/earthly/examples/go+docker
# Run the resulting image.
docker run --rm go-example:latest
```

### 🔨 Reference other targets using +

Use `+` to reference other targets and create complex build inter-dependencies.

<div align="center"><a href="https://docs.earthly.dev/guides/target-ref"><img src="docs/guides/img/ref-infographic.png" alt="Target and artifact reference syntax" title="Reference targets using +" width="600px" /></a></div>

Examples

* Same directory (same Earthfile)
  
  ```Dockerfile
  BUILD +some-target
  FROM +some-target
  COPY +some-target/my-artifact ./
  ```

* Other directories

  ```Dockerfile
  BUILD ./some/local/path+some-target
  FROM ./some/local/path+some-target
  COPY ./some/local/path+some-target/my-artifact ./
  ```

* Other repositories

  ```Dockerfile
  BUILD github.com/someone/someproject:v1.2.3+some-target
  FROM github.com/someone/someproject:v1.2.3+some-target
  COPY github.com/someone/someproject:v1.2.3+some-target/my-artifact ./
  ```

### 💾 Caching that works the same as docker builds

<div align="center"><a href="https://asciinema.org/a/351674?speed=2"><img src="img/demo-351674.gif" alt="Demonstration of Earthly's caching" title="View on asciinema.org" width="600px" /></a></div>

### 🛠 Reusability with build args

Here is an example where building for multiple platforms can leverage build args.

```Dockerfile
FROM golang:1.13-alpine3.11
RUN apk add --update --no-cache g++
WORKDIR /go-example

all:
  BUILD \
    --build-arg GOOS=linux \
    --build-arg GOARCH=amd64 \
    --build-arg GO_LDFLAGS="-linkmode external -extldflags -static" \
    +build
  BUILD \
    --build-arg GOOS=darwin \
    --build-arg GOARCH=amd64 \
    +build
  BUILD \
    --build-arg GOOS=windows \
    --build-arg GOARCH=amd64 \
    +build

build:
  COPY main.go .
  ARG GOOS
  ARG GOARCH
  ARG GO_LDFLAGS
  RUN go build -ldflags "$GO_LDFLAGS" -o build/go-example main.go && \
      echo "Build for $GOOS/$GOARCH was successful"
  SAVE ARTIFACT build/go-example AS LOCAL "build/$GOOS/$GOARCH/go-example"
```

### ⛓ Parallelization that just works

Whenever possible, Earthly automatically executes targets in parallel.

<div align="center"><a href="https://asciinema.org/a/351678?speed=2"><img src="img/demo-351678.gif" alt="Demonstration of Earthly's parallelization" title="View on asciinema.org" width="600px" /></a></div>

### 🤲 Make use of build tools that work everywhere

No need to ask your team to install `protoc`, a specific version of Python, Java 1.6 or the .NET Core ecosystem. You only install once, in your Earthfile, and it works for everyone. Or even better, you can just make use of the rich Docker Hub ecosystem.

```Dockerfile
FROM golang:1.13-alpine3.11
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

### 🔑 Secrets support built-in

Secrets are never stored within an image's layers and they are only available to the commands that need them.

```Dockerfile
release:
  RUN --push --secret GITHUB_TOKEN=+secrets/GITHUB_TOKEN github-release upload file.bin
```

```bash
earth --secret GITHUB_TOKEN --push +release
```

<br/>
<br/>
<h2 align="center">FAQ</h2>

### How is Earthly different from Dockerfiles?

[Dockerfiles](https://docs.docker.com/engine/reference/builder/) were designed for specifying the make-up of Docker images and that's where Dockerfiles stop. Earthly takes some key principles of Dockerfiles (like layer caching), but expands on the use-cases. For example, Earthly can output regular artifacts, run unit and integration tests and also create several Docker images at a time - all of which are outside the scope of Dockerfiles.

It is possible to use Dockerfiles in combination with other technologies (eg Makefiles or bash files) in order to solve for such use-cases. However, these combinations are difficult to parallelize, difficult to scale across repositories as they lack a robust import system and also they often vary in style from one team to another. Earthly does not have these limitations as it was designed as a general purpose build system.

As an example, Earthly introduces a richer target, artifact and image [referencing system](https://docs.earthly.dev/guides/target-ref), which allows for better reuse in complex builds spanning a single large repository or multiple repositories. Because Dockerfiles are only meant to describe one image at a time, such features are outside the scope of applicability of Dockerfiles.

### How do I tell apart classical Dockerfile commands from Earthly commands?

Check out the [Earthfile reference doc page](https://docs.earthly.dev/earthfile). It has all the commands there and it specifies which commands are the same as Dockerfile commands and which are new.

### Can Earthly build Dockerfiles?

Yes! You can use the command `FROM DOCKERFILE` to inherit the commands in an existing Dockerfile.

```Dockerfile
build:
  FROM DOCKERFILE .
  SAVE IMAGE some-image:latest
```

You may also optionally port your Dockerfiles to Earthly entirely. Translating Dockerfiles to Earthfiles is usually a matter of copy-pasting and making small adjustments. See the [getting started page](https://docs.earthly.dev/guides/basics) for some Earthfile examples.

### How is Earthly different from Bazel?

[Bazel](https://bazel.build) is a build tool developed by Google for the purpose of optimizing speed, correctness and reproducibility of their internal monorepo codebase. Earthly draws inspiration from some of the principles of Bazel (mainly reproducibility), but it is different in a few key ways:

* Earthly does not replace language-specific tools, like Maven, Gradle, Webpack etc. Instead, it leverages and integrates with them. Adopting Bazel usually means that all build files need to be completely rewritten. This is not the case with Earthly as it mainly acts as the glue between builds.
* The learning curve of Earthly is more accessible, especially if the user already has experience with Dockerfiles. Bazel, on the other hand, introduces some completely new concepts.
* Bazel has a purely descriptive specification language. Earthly is a mix of descriptive and imperative language.
* Bazel uses tight control of compiler tool chain to achieve consistent builds, whereas Earthly uses containers and well-defined inputs.

Overall, compared to Bazel, Earthly sacrifices a little correctness and reproducibility in favor of significantly better usability and composability with existing open-source technologies.

<br/>
<br/>
<h2 align="center">Contributing</h2>

* Please report bugs as [GitHub issues](https://github.com/earthly/earthly/issues).
* Join us on [Gitter](https://gitter.im/earthly-room/community)!
* Questions via GitHub issues are welcome!
* PRs welcome! But please give a heads-up in GitHub issue before starting work. If there is no GitHub issue for what you want to do, please create one.
* To build from source, you will need the `earth` binary ([Earthly builds itself](https://docs.earthly.dev/examples/earthly)). Git clone the code and run `earth +all`. To run the tests, run `earth -P +test`.

<br/>
<br/>
<h2 align="center">Licensing</h2>

Earthly is licensed under the Mozilla Public License Version 2.0. See [LICENSE](./LICENSE) for the full license text.
