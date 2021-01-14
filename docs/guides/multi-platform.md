# Multi-platform builds

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/536) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#multi-platform` channel.
{% endhint %}

Earthly has the ability to perform builds for multiple platforms, in parallel. This page walks through setting up your system to support emulation as well as through a few simple examples of how to use this feature.

Currently only `linux` is supported as the build platform OS. Building with Windows containers will be available in a future version of Earthly.

By default, builds are performed on the same processor architecture as available on the host natively. Using the `--platform` flag across various Earthfile commands or as part of the `earthly` command, it is possible to override the build platform and thus be able to execute builds on non-native processor architectures. Execution of non-native binaries can be performed via QEMU emulation.

In some cases, execution of the build itself does not need to happen on the target architecture, through cross-compilation features of the compiler. Examples of languages that support cross-compilation are Go and Rust. This approach may be more beneficial in many cases, as there is no need to install QEMU and also, the build is more performant.

## Pre-requisites for emulation

In order to execute emulated build steps (usually `RUN`), QEMU needs to be installed and set up. This will allow you perform Earthly builds on non-native platforms, but also incidentally, to run Docker images on your host system through `docker run --platform=...`.

### Windows and Mac

On Mac and on Windows, the Docker Desktop app comes with QEMU readilly installed and ready to go, so no special consideration is necessary.

### Linux

On linux, QEMU needs to be installed manually. On Ubuntu, this can be achieved by running:

```
sudo apt-get install qemu binfmt-support qemu-user-static
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker stop earthly-buildkitd || true
```

However, note that the command

```
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

will need to be re-run on every system restart. For a more permanent installation of QEMU see this guide on [building multi architecture Docker images, which includes QEMU installation instructions](https://medium.com/@artur.klauser/building-multi-architecture-docker-images-with-buildx-27d80f7e2408).

### GitHub Actions

To make use of emulation in GitHub Actions, the following step needs to be included in every job that performs a multi-platform build:

```yaml
jobs:
    <job-name>:
        steps:
            -
                name: Set up QEMU
                id: qemu
                uses: docker/setup-qemu-action@v1
                with:
                    image: tonistiigi/binfmt:latest
                    platforms: all
            - uses: actions/checkout@v2
            - ...
```

## Performing multi-platform builds

In order to execute builds for multiple platforms, the execution may be parallelized through the repeated use of the `BUILD --platform` flag. For example:

```Dockerfile
build-all-platforms:
    BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build

build:
    ...
```

If the `+build` target were invoked without the use of any flag, Earthly would simply perform the build on the native architecture of the host system.

However, invoking the target `+build-all-platforms` causes `+build` to execute twice, in parallel: one time on `linux/amd64` and another time on `linux/arm/v7`.

You may also override the target platform when issuing the `earthly` build command. For example:

```bash
earthly --platform=linux/arm64 +build
```

This would cause the build to execute on the `linux/arm64` architecture.

## Saving multi-platform images

The easiest way to include platform information as part of a build is through the use of `FROM --platform`. For example:

```Dockerfile
FROM --platform=linux/arm/v7 alpine:3.11
```

If multiple targets create an image with the same name, but for different platforms, the images will be merged into a multi-platform image during export. For example:

```Dockerfile
build-all-platforms:
    BUILD +build-amd64
    BUILD +build-arm-v7

build-amd64:
    FROM --platform=linux/amd64 alpine:3.11
    ...
    SAVE IMAGE --push org/myimage:latest

build-arm-v7:
    FROM --platform=linux/arm/v7 alpine:3.11
    ...
    SAVE IMAGE --push org/myimage:latest
```

When `earthly --push +build-all-platforms` is executed, the build will push a multi-manifest image to the Docker registry. The manifest will contain two images: one for `linux/amd64` and one for `linux/arm/v7`. This works as such because both targets that save images use the exact same Docker tag for the image.

Of course, in some situations, the build steps are the same (except they run on different platform), so the two definitions can be merged like so:

```Dockerfile
build-all-platforms:
    BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build

build:
    FROM alpine:3.11
    ...
    SAVE IMAGE --push org/myimage:latest
```

A more complete version of this example is available in [examples/multiplatform](https://github.com/earthly/earthly/tree/main/examples/multiplatform) in GitHub. You may try out this example without cloning by running

```bash
earthly github.com/earthly/earthly/examples/multiplatform:main+all
docker run --rm earthly/examples:multiplatform
docker run --rm earthly/examples:multiplatform_linux_amd64
docker run --rm earthly/examples:multiplatform_linux_arm_v7
```

{% hint style='info' %}
##### Note
As of the time of writing this article, the `docker` CLI has limited support for working with multi-manifest images locally. For this reason, when exporting an image to the local Docker daemon, Earthly provides the different architectures as different Docker tags.

For example, the above build would yield locally:

* `org/myimage:latest`
* `org/myimage:latest_linux_amd64` (the same as `org/myimage:latest` if running on a `linux/amd64` host)
* `org/myimage:latest_linux_arm_v7`

The additional Docker tags are only available for use on the local system. When pushing an image to a Docker registry, it is pushed as a single multi-manifest image.
{% endhint %}

## Creating multi-platform images without emulation

Building multi-platform images does not necessarily require that execution of the build itself takes place on the target platform. Through the use of cross-compilation, it is possible to obtain target-platform binaries compiled on the host-native platform. At the end, these binaries may be placed in a final image which is marked for a specific platform.

Note, however, that not all programming languages have support for cross-compilation. The applicability of this approach may be limited as a result. Examples of languages that *can* cross-compile for other platforms are Go and Rust.

Here is an example where a multi-platform image can be created without actually executing any `RUN` on the target platform (and therefore emulation is not necessary):

```Dockerfile
build-all-platforms:
    BUILD +build-amd64
    BUILD +build-arm-v7

build:
    FROM golang:1.13-alpine3.11
    WORKDIR /example
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG GOARM
    COPY main.go ./
    RUN go build -o main main.go
    SAVE ARTIFACT ./main

build-amd64:
    FROM --platform=linux/amd64 alpine:3.11
    COPY +build/main ./example/main
    ENTRYPOINT ["/example/main"]
    SAVE IMAGE --push org/myimage:latest

build-arm-v7:
    FROM --platform=linux/arm/v7 alpine:3.11
    COPY \
        --platform=linux/amd64 \
        --build-arg GOARCH=arm \
        --build-arg GOARM=v7 \
        +build/main ./example/main
    ENTRYPOINT ["/example/main"]
    SAVE IMAGE --push org/myimage:latest
```

The key here is the use of the `COPY` commands. The execution of the target `+build` may take place on the host platform (in this case, `linux/amd64`) and yet produce binaries for either `amd64` or `arm/v7`. Since there is no `RUN` command as part of the `+build-arm-v7` target, no emulation is necessary.

## Making use of builtin platform args

A number of [builtin build args](../earthfile/builtin-args.md) are made available to be used in conjunction with multi-platform builds:

* `TARGETPLATFORM` (eg `linux/arm/v7`)
* `TARGETOS` (eg `linux`)
* `TARGETARCH` (eg `arm`)
* `TARGETVARIANT` (eg `v7`)

Here is an example of how the build described above could be simplified through the use of these build args:

```Dockerfile
build-all-platforms:
    BUILD --platform=linux/amd64 --platform=linux/arm/v7 +build-image

build:
    FROM golang:1.13-alpine3.11
    WORKDIR /example
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG VARIANT
    COPY main.go ./
    RUN GOARM=${VARIANT#"v"} go build -o main main.go
    SAVE ARTIFACT ./main

build-image:
    ARG TARGETPLATFORM
    ARG TARGETARCH
    ARG TARGETVARIANT
    FROM --platform=$TARGETPLATFORM alpine:3.11
    COPY \
        --platform=linux/amd64 \
        --build-arg GOARCH=$TARGETARCH \
        --build-arg VARIANT=$TARGETVARIANT \
        +build/main ./example/main
    ENTRYPOINT ["/example/main"]
    SAVE IMAGE --push org/myimage:latest
```

The code of this example is available in [examples/multiplatform-cross-compile](https://github.com/earthly/earthly/tree/main/examples/multiplatform-cross-compile) in GitHub. You may try out this example without cloning by running

```bash
earthly github.com/earthly/earthly/examples/multiplatform-cross-compile:main+build-all-platforms
```
