
# Table of contents

* [👋 Introduction](README.md)
* [⬇️ Installation](https://earthly.dev/get-earthly)
* [🎓 Learn the basics](basics/basics.md)
    * [Part 1: A simple Earthfile](basics/part-1-a-simple-earthfile.md)
    * [Part 2: Detailed explanation](basics/part-2-detailed-explanation.md)
    * [Part 3: Adding dependencies in the mix](basics/part-3-adding-dependencies-in-the-mix.md)
    * [Part 4: Efficient caching of dependencies](basics/part-4-efficient-caching-of-dependencies.md)
    * [Part 5: Reduce code duplication](basics/part-5-reduce-code-duplication.md)
    * [Final words](basics/final-words.md)

## 📖 Docs

* Guides
    * [Authenticating Git and image registries](guides/auth.md)
    * [Target, artifact and command referencing](guides/target-ref.md)
    * [Build arguments and secrets](guides/build-args.md)
    * [User-defined commands (UDCs)](guides/udc.md)
    * [Managing cache](guides/cache.md)
    * [Advanced local caching](guides/advanced-local-caching.md)
    * [Shared cache](guides/shared-cache.md)
    * [Using Docker in Earthly](guides/docker-in-earthly.md)
    * [Cloud secrets](guides/cloud-secrets.md)
    * [Integration Testing](guides/integration.md)
    * [Debugging techniques](guides/debugging.md)
    * [Multi-platform builds](guides/multi-platform.md)
    * Configuring registries
        * [AWS ECR](guides/registries/aws-ecr.md)
        * [GCP Artifact Registry](guides/registries/gcp-artifact-registry.md)
        * [Azure ACR](guides/registries/azure-acr.md)
        * [Self-signed certificates](guides/registries/self-signed.md)
    * Using the Earthly Docker Images
        * [earthly/earthly](docker-images/all-in-one.md)
        * [earthly/buildkitd](docker-images/buildkit-standalone.md)
* [Earthfile reference](earthfile/earthfile.md)
    * [Builtin args](earthfile/builtin-args.md)
    * [Excluding patterns](earthfile/earthignore.md)
    * [Experimental features](earthfile/features.md)
* [The `earthly` command](earthly-command/earthly-command.md)
* [Configuration reference](earthly-config/earthly-config.md)
* [CI integration](./ci-integration.md)
* Misc
    * [Alternative installation](./alt-installation.md)
    * [Definitions](definitions/definitions.md)
    * [Data collection](data-collection/data-collection.md)

## 🧙 Examples

* [Overview](examples/examples.md)
* CI examples
    * [GitHub Actions](ci-examples/gh-actions-integration.md)
    * [Circle CI](ci-examples/circle-integration.md)
    * [AWS CodeBuild](ci-examples/codebuild-integration.md)
