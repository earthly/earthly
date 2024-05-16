
# Examples

## Examples of CI integration

Examples of integrating Earthly into various CI systems can be found on the following pages:

* [Circle CI](../ci-integration/guides/circle-integration.md)
* [GitHub Actions](../ci-integration/guides/gh-actions-integration.md)
* [AWS CodeBuild](../ci-integration/guides/codebuild-integration.md)
* [Jenkins](../ci-integration/guides/jenkins.md)
* [Kubernetes](../ci-integration/guides/kubernetes.md)

For more general information on CI systems not listed above, see the [CI integration guide](../ci-integration/overview.md).

## Examples by dev environments

Examples of how Earthly can be integrated into different dev environments

* [GitHub Codespaces](https://github.com/earthly/codespaces-example) - shows how Earthly can be used in GitHub Codespaces and Devcontainers

## Example Earthfiles

In this section, you will find some examples of Earthfiles to familiarize yourself with Earthly.

The code for all the examples is available in the [examples GitHub directory](https://github.com/earthly/earthly/tree/main/examples).

<!-- NOTE: If you change this, please also change examples/README.md -->

### Examples from the Basics tutorial

If you are new to Earthly, you may find the [Basics tutorial](../basics/basics.md) helpful.

* [tutorial](https://github.com/earthly/earthly/tree/main/examples/tutorial)
    * [go](https://github.com/earthly/earthly/tree/main/examples/tutorial/go)
    * [js](https://github.com/earthly/earthly/tree/main/examples/tutorial/js)
    * [java](https://github.com/earthly/earthly/tree/main/examples/tutorial/java)
    * [python](https://github.com/earthly/earthly/tree/main/examples/tutorial/python)

### Examples by language

Please note that these examples, although similar, are distinct from the ones used in the [tutorial](https://github.com/earthly/earthly/tree/main/examples/tutorial).

<!-- vale HouseStyle.Spelling = NO -->
* [c](https://github.com/earthly/earthly/tree/main/examples/c)
* [clojure](https://github.com/earthly/earthly/tree/main/examples/clojure)
* [cobol](https://github.com/earthly/earthly/tree/main/examples/cobol)
* [cpp](https://github.com/earthly/earthly/tree/main/examples/cpp)
* [dotnet](https://github.com/earthly/earthly/tree/main/examples/dotnet)
* [elixir](https://github.com/earthly/earthly/tree/main/examples/elixir)
* [go](https://github.com/earthly/earthly/tree/main/examples/go)
* [java](https://github.com/earthly/earthly/tree/main/examples/java)
* [js](https://github.com/earthly/earthly/tree/main/examples/js)
* [python](https://github.com/earthly/earthly/tree/main/examples/python)
* [ruby](https://github.com/earthly/earthly/tree/main/examples/ruby)
* [ruby-on-rails](https://github.com/earthly/earthly/tree/main/examples/ruby-on-rails)
* [rust](https://github.com/earthly/earthly/tree/main/examples/rust)
* [scala](https://github.com/earthly/earthly/tree/main/examples/scala)
* [typescript-node](https://github.com/earthly/earthly/tree/main/examples/typescript-node)
<!-- vale HouseStyle.Spelling = YES -->

### Examples by use-cases

* [integration-test](https://github.com/earthly/earthly/tree/main/examples/integration-test) - shows how `WITH DOCKER` and `docker-compose` can be used to start up services and then run an integration test suite.
* [monorepo](https://github.com/earthly/earthly/tree/main/examples/monorepo) - shows how multiple sub-projects can be co-located in a single repository and how the build can be fragmented across these.
* [multirepo](https://github.com/earthly/earthly/tree/main/examples/multirepo) - shows how artifacts from multiple repositories can be referenced in a single build. See also the `grpc` example for a more extensive use-case.

### Examples by Earthly features

* [import](https://github.com/earthly/earthly/tree/main/examples/import) - shows how to use the `IMPORT` command to alias Earthfile references.
* [cutoff-optimization](https://github.com/earthly/earthly/tree/main/examples/cutoff-optimization) - shows that if an intermediate artifact does not change, then the rest of the build will use the cache, even if the source has changed.
* [multiplatform](https://github.com/earthly/earthly/tree/main/examples/multiplatform) - shows how Earthly can execute builds and create images for multiple platforms, using QEMU emulation.
* [multiplatform-cross-compile](https://github.com/earthly/earthly/tree/main/examples/multiplatform-cross-compile) - shows has through the use of cross-compilation, you can create images for multiple platforms, without using QEMU emulation.

### Examples by use of other technologies

* [grpc](https://github.com/earthly/earthly/tree/main/examples/grpc) - shows how to use Earthly to compile a protobuf grpc definition into protobuf code for both a Go-based server, and a python-based client, in a multirepo setup.
* [terraform](https://github.com/earthly/earthly/tree/main/examples/terraform) - shows how Terraform could be used from Earthly.

### Other

* [readme](https://github.com/earthly/earthly/tree/main/examples/readme) - some sample code we used in our README.
* [tests](https://github.com/earthly/earthly/tree/main/tests) - a suite of tests Earthly uses to ensure that its features are working correctly.

### Larger Examples And Community Examples

* [Earthly, Rust, GoLang, NodeJS and GitHub Actions Example](https://github.com/earthly/earthly-vs-gha)
* [Cloud Services In GoLang](https://github.com/earthly/cloud-services-example)
* [Earthfile workshop Repo](https://github.com/earthly/workshop-2023-09-18)
* [Python & C Example](https://github.com/earthly/pymerge)
* [Python Docker Example](https://github.com/earthly/build-transpose/blob/main/Earthfile)
* [Awesome Earthly - Community Examples](https://github.com/earthly/awesome-earthly)

### Earthly's own build

As a distinct example of a complete build, you can take a look at Earthly's own build. Earthly builds itself, and the build files are available on GitHub:

<!--

GitBook currently has a bug where any references to an "Earthfile" gets confused with "docs/Earthfile" and somehow appends a /README.md

e.g. https://github.com/earthly/earthly/blob/main/Earthfile is changed to https://github.com/earthly/earthly/blob/main/Earthfile/README.md

Here's a snip from an support request with gitbook:

    On Thu, Dec 23, 2021 at 7:15:12 UTC, GitBook Support <support@gitbook.com> wrote:

    There is a file:

    https://github.com/earthly/earthly/blob/main/Earthfile

    And you want to reference it directly in your GitBook space as a link.

    The problem here is that GitBook is thrown off by the fact it has a folder under the docs root. Remember you documentation root is set to /docs.

    So when it sees that reference, it assumes you are referencing a default README.md file under that folder. The folder I am talking about is this one:

    https://github.com/earthly/earthly/tree/main/docs/earthfile

    Now, the question is, if there's an easy way out of this.

    On Thu, Dec 23, 2021 at 11:41:41 UTC, GitBook Support <support@gitbook.com> wrote:

    I can't confirm it yet, but this might be an edge case that we could patch.

    One not very ideal workaround I thought of is to temporarily switch to shortened URLs for those that fail because of this scenario.


* [Earthfile](https://github.com/earthly/earthly/blob/main/Earthfile) - the root build file
* [buildkitd/Earthfile](https://github.com/earthly/earthly/blob/main/buildkitd/Earthfile) - the build of the BuildKit daemon
* [AST/parser/Earthfile](https://github.com/earthly/earthly/blob/main/ast/parser/Earthfile) - the build of the parser, which generates .go files
* [tests/Earthfile](https://github.com/earthly/earthly/blob/main/tests/Earthfile) - system and smoke tests
* [earthfile-grammar/Earthfile](https://github.com/earthly/earthfile-grammar/blob/main/Earthfile) - the build of the VS Code extension
-->

* [Earthfile](https://tinyurl.com/yt3d3cx6) - the root build file
* [buildkitd/Earthfile](https://tinyurl.com/yvnpuru7) - the build of the BuildKit daemon
* [AST/parser/Earthfile](https://tinyurl.com/2k3u4vty) - the build of the parser, which generates .go files
* [tests/Earthfile](https://tinyurl.com/2p8ws579) - system and smoke tests
* [earthfile-grammar/Earthfile](https://tinyurl.com/2vyjprt6) - the build of the VS Code extension

To invoke Earthly's build, check out the code and then run the following in the root of the repository

```bash
earthly +all
```

[![asciicast](https://asciinema.org/a/313845.svg)](https://asciinema.org/a/313845)
