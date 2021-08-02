
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
* [Examples](examples/examples.md)
* Misc
    * [Alternative installation](./alt-installation.md)
    * [Definitions](definitions/definitions.md)
    * [Data collection](data-collection/data-collection.md)

## 🔧 CI Integration
* [Overview](ci-integration/overview.md)
* [Build An Earthly CI Image](ci-integration/build-an-earthly-ci-image.md)
* [Pull-Through Cache](ci-integration/pull-through-cache.md)
* [Remote Buildkit](ci-integration/remote-buildkit.md)
* Vendor-Specific Guides
  * [Jenkins](ci-integration/guides/jenkins.md)
  * [Circle CI](ci-integration/guides/circle-integration.md)
  * [GitHub Actions](ci-integration/guides/gh-actions-integration.md)
  * [AWS CodeBuild](ci-integration/guides/codebuild-integration.md)
  * [Kubernetes](ci-integration/guides/kubernetes.md)
