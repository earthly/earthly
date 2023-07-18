# Earthly examples

This folder contains a series of examples to help you familiarize yourself with Earthly.

<!-- NOTE: If you change this, please also change docs/examples/examples.md -->

## Examples from the Basics tutorial

If you are new to Earthly, you may find the [Basics tutorial](https://docs.earthly.dev/basics) helpful.

<!-- vale HouseStyle.Spelling = NO -->
* [tutorial](./tutorial)
    * [go](./tutorial/go)
    * [js](./tutorial/js)
    * [java](./tutorial/java)
    * [python](./tutorial/python)
<!-- vale HouseStyle.Spelling = YES -->

## Examples by language

Please note that these examples, although similar, are distinct from the ones used in the [tutorial](./tutorial).

<!-- vale HouseStyle.Spelling = NO -->
* [clojure](./clojure)
* [cobol](./cobol)
* [c](./c)
* [cpp](./cpp)
* [dotnet](./dotnet)
* [elixir](./elixir)
* [go](./go)
* [java](./java)
* [js](./js)
* [next-js-netlify](./next-js-netlify)
* [python](./python)
* [ruby](./ruby)
* [ruby-on-rails](./ruby-on-rails)
* [rust](./rust)
* [scala](./scala)
* [typescript-node](./typescript-node)
<!-- vale HouseStyle.Spelling = YES -->

## Examples by use-cases

* [integration-test](./integration-test) - shows how `WITH DOCKER` and `docker-compose` can be used to start up services and then run an integration test suite.
* [monorepo](./monorepo) - shows how multiple sub-projects can be co-located in a single repository and how the build can be fragmented across these.
* [multirepo](./multirepo) - shows how artifacts from multiple repositories can be referenced in a single build. See also the `grpc` example for a more extensive use-case.

## Examples by Earthly features

* [import](./import) - shows how to use the `IMPORT` command to alias project references.
* [cutoff-optimization](./cutoff-optimization) - shows that if an intermediate artifact does not change, then the rest of the build will use the cache, even if the source has changed.
* [multiplatform](./multiplatform) - shows how Earthly can execute builds and create images for multiple platforms, using QEMU emulation.
* [multiplatform-cross-compile](./multiplatform-cross-compile) - shows has through the use of cross-compilation, you can create images for multiple platforms, without using QEMU emulation.

## Examples by use of other technologies

* [grpc](./grpc) - shows how to use Earthly to compile a protobuf grpc definition into protobuf code for both a Go-based server, and a python-based client, in a multirepo setup.
* [terraform](./terraform) - shows how Terraform could be used from Earthly.

## Other

* [readme](./readme) - some sample code we used in our README.
* [tests](./tests) - a suite of tests Earthly uses to ensure that its features are working correctly.
