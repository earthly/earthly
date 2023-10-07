# Earthly Internals Developer Documentation

This documents the internals of the Earthly build tool; it is written for developers who wish to [contribute](../CONTRIBUTING.md) to the Earthly codebase.

Users of Earthly can find general user-documentation at [https://docs.earthly.dev/](https://docs.earthly.dev/).

## Warning

This guide may become outdated as the Earthly codebase evolves. It should be used to supplement the code-base, which is the single source of truth.

## BuildKit

Earthly use BuildKit for executing commands, caching, and generating images. BuildKit contains similar contributor-focused documentation under
the [buildkit/docs/dev](https://github.com/moby/buildkit/tree/master/docs/dev) section of their codebase.

## Jargon

| Name | Description |
| :--- | :---------- |
| **BuildKit** | BuildKit is a toolkit for converting source code to build artifacts in an efficient, expressive and repeatable manner. [^1] |
| **LLB** | LLB is a BuildKit concept, which stands for "low-level build" definition[^2]; Earthly converts Earthfiles into LLB definitions, which is sent to BuildKit via the BuildKit client |
| **LLB State** | LLB State, or simply State, is another BuildKit concept, which is used to produce low-level build definitions (LLB) from higher-level concepts like images, shell executions, mounts, etc [^2] |
| **pllb** | pllb is an Earthly thread-safe "parallel" wrapper around the LLB State functions |
| **ast** | Abstract syntax tree; a custom Earthfile grammar is defined under the ast package, which [ANTLR](https://www.antlr.org/) uses to parse the initial Earthfile |
| **buildcontext** | Borrowed from the `docker build --build-context` option, the buildcontext package ties the locations `COPY` reference to be relative to the corresponding Earthfile  |
| **resolver** | The resolver takes an Earthly target (e.g. `./sub/dir+target`, or `github.com/...+target`), and constructs a llb state that points to the buildcontext   |
| **builder** | The builder is the initial entrypoint, for the Earthly `build` cli command, it contains the `buildFunc` which is passed to the BuildKit client |
| **interpreter** | The Earthly interpreter walks the ast, performing additional parsing and validation, and makes appropriate calls to the converter |
| **converter** | The converter produces LLB, which is sent to BuildKit via the BuildKit gateway client (gwclient) |
| **build function** | A function which is passed to the initial BuildKit client's `Build(...)` function; the client performs a callback along with a newly created gateway client, which accepts LLB definitions |
| **mts** | Multi-target states; which holds multiple LLB States, in the order they should be built |
| **pullping** | Once the build function returns (passing a set of LLB references back to buildkit), the BuildKit server will execute the commands, and call the earthlyoutputs exporter, which will call back to the client (Earthly), which will be received by the pullping handler. This will cause earthly to perform a `docker pull ...` against the embedded registry |
| **dockertar** | The legacy approach for exporting images from BuildKit to the host via a `tar` file; we try to use pullping instead, since it only pulls the needed layers |
| **logbus** | An interface for writing output to both stdout and the web-based log viewer under cloud.earthly.dev |
| **earthlyoutputs** | A custom buildkit exporter (within the [earthly/buildkit fork](https://github.com/earthly/buildkit/tree/earthly-main/exporter/earthlyoutputs)), which is used to send images back to earthly |
| **embedded registry** | A [docker registry](https://github.com/distribution/distribution) which runs within the earthly-buildkitd container, used in combination with earthlyoutputs and the pullping callback; also referred to as "local registry" |

## Guides

[build action lifecycle](build-steps.md)

## Citations

[^1]: https://github.com/moby/buildkit/blob/2677a22857c917168730fe69ad617a50e0d85202/README.md
[^2]: https://github.com/moby/buildkit/blob/2677a22857c917168730fe69ad617a50e0d85202/docs/dev/README.md
