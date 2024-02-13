# Go Monorepo Example

The following example demonstrates some ways to use Earthly effectively to build a multi-module monorepo.

The example contains two Microservices in a `services` directory, and a single shared library in a `libs` directory.

Each sub-module of this monorepo manages it's own build, test, and release steps, with parent Earthfile acting as an orchestrator.

The `services` import code from `libs` during their Earthly build by utilizing [artifacts](https://docs.earthly.dev/docs/earthfile?q=save+artifact).

Some other noteworthy features of this demo include releases microservices based on a `.semver.yaml` file, starting the entire stack locally, and running unit-tests across the entire monorepo in parallel.

For a more detailed explanation, you can read about this example in [our blog post](https://earthly.dev/blog/golang-monorepo/).
