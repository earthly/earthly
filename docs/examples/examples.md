
# Examples

## Examples of CI integration

Examples of integrating Earthly into various CI systems can be found on the following pages:

* [Circle CI](../ci-integration/guides/circle-integration.md)
* [GitHub Actions](../ci-integration/guides/gh-actions-integration.md)
* [AWS CodeBuild](../ci-integration/guides/codebuild-integration.md)
* [Jenkins](../ci-integration/guides/jenkins.md)
* [Kubernetes](../ci-integration/guides/kubernetes.md)

For more general information on CI systems not listed above, see the [CI integration guide](../ci-integration/overview.md).

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
<!-- vale HouseStyle.Spelling = YES -->

### Examples by use-cases

* [integration-test](https://github.com/earthly/earthly/tree/main/examples/integration-test) - shows how `WITH DOCKER` and `docker-compose` can be used to start up services and then run an integration test suite.
* [monorepo](https://github.com/earthly/earthly/tree/main/examples/monorepo) - shows how multiple sub-projects can be co-located in a single repository and how the build can be fragmented across these.
* [multirepo](https://github.com/earthly/earthly/tree/main/examples/multirepo) - shows how artifacts from multiple repositories can be referenced in a single build. See also the `grpc` example for a more extensive use-case.

### Examples by Earthly features

* [import](https://github.com/earthly/earthly/tree/main/examples/import) - shows how to use the `IMPORT` command to alias project references.
* [cutoff-optimization](https://github.com/earthly/earthly/tree/main/examples/cutoff-optimization) - shows that if an intermediate artifact does not change, then the rest of the build will use the cache, even if the source has changed.
* [multiplatform](https://github.com/earthly/earthly/tree/main/examples/multiplatform) - shows how Earthly can execute builds and create images for multiple platforms, using QEMU emulation.
* [multiplatform-cross-compile](https://github.com/earthly/earthly/tree/main/examples/multiplatform-cross-compile) - shows has through the use of cross-compilation, you can create images for multiple platforms, without using QEMU emulation.

### Examples by use of other technologies

* [grpc](https://github.com/earthly/earthly/tree/main/examples/grpc) - shows how to use Earthly to compile a protobuf grpc definition into protobuf code for both a Go-based server, and a python-based client, in a multirepo setup.
* [terraform](https://github.com/earthly/earthly/tree/main/examples/terraform) - shows how Terraform could be used from Earthly.

### Other

* [readme](https://github.com/earthly/earthly/tree/main/examples/readme) - some sample code we used in our README.
* [tests](https://github.com/earthly/earthly/tree/main/tests) - a suite of tests Earthly uses to ensure that its features are working correctly.

### Earthly's own build

As a distinct example of a complete build, you can take a look at Earthly's own build. Earthly builds itself, and the build files are available on GitHub:

* [Earthfile](https://github.com/earthly/earthly/blob/main/Earthfile) - the root build file
* [buildkitd/Earthfile](https://github.com/earthly/earthly/blob/main/buildkitd/Earthfile) - the build of the Buildkit daemon
* [AST/parser/Earthfile](https://github.com/earthly/earthly/blob/main/ast/parser/Earthfile) - the build of the parser, which generates .go files
* [examples/tests/Earthfile](https://github.com/earthly/earthly/blob/main/examples/tests/Earthfile) - system and smoke tests
* [contrib/earthfile-syntax-highlighting/Earthfile](https://github.com/earthly/earthly/blob/main/contrib/earthfile-syntax-highlighting/Earthfile) - the build of the VS Code extension

To invoke Earthly's build, check out the code and then run the following in the root of the repository

```bash
earthly +all
```

[![asciicast](https://asciinema.org/a/313845.svg)](https://asciinema.org/a/313845)
