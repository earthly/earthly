
# Table of contents

* [👋 Introduction](README.md)
* [💻 Install Earthly](install/install.md)
* [🎓 Learn the basics](basics/basics.md)
    * [Part 1: A simple Earthfile](basics/part-1-a-simple-earthfile.md)
    * [Part 2: Outputs](basics/part-2-outputs.md)
    * [Part 3: Adding dependencies With Caching](basics/part-3-adding-dependencies-with-caching.md)
    * [Part 4: Args](basics/part-4-args.md)
    * [Part 5: Importing](basics/part-5-importing.md)
    * [Part 6: Using Docker In Earthly](basics/part-6-using-docker-with-earthly.md)
    * [Part 7: Using remote runners](basics/part-7-using-remote-runners.md)
    * [Part 8a: Using Earthly in your current CI](basics/part-8a-using-earthly-in-your-current-ci.md)
    * [Final words](basics/final-words.md)
* [⭐ Featured guides](lang-guides/featured-guides.md)
    * [Rust](lang-guides/rust.md)

## 📖 Docs

* Guides
    * [Importing](guides/importing.md)
    * [Build arguments and variables](guides/build-args.md)
    * [Secrets](guides/secrets.md)
    * [Functions](guides/functions.md)
    * [Using Docker in Earthly](guides/docker-in-earthly.md)
    * [Multi-platform builds](guides/multi-platform.md)
    * [Authenticating Git and image registries](guides/auth.md)
    * [Integration Testing](guides/integration.md)
    * [Debugging techniques](guides/debugging.md)
    * [Podman](guides/podman.md)
    * Configuring registries
        * [AWS ECR](guides/registries/aws-ecr.md)
        * [GCP Artifact Registry](guides/registries/gcp-artifact-registry.md)
        * [Azure ACR](guides/registries/azure-acr.md)
        * [Self-signed certificates](guides/registries/self-signed.md)
    * Using the Earthly Docker Images
        * [earthly/earthly](docker-images/all-in-one.md)
        * [earthly/buildkitd](docker-images/buildkit-standalone.md)
    * [Best practices](guides/best-practices.md)
* [Caching](./caching/caching.md)
    * [Caching in Earthfiles](./caching/caching-in-earthfiles.md)
    * [Managing cache](./caching/managing-cache.md)
    * [Caching via remote runners](./caching/caching-via-remote-runners.md)
    * [Caching via a registry (advanced)](./caching/caching-via-registry.md)
* [Remote runners](remote-runners.md)
* [Earthfile reference](earthfile/earthfile.md)
    * [Builtin args](earthfile/builtin-args.md)
    * [Excluding patterns](earthfile/earthlyignore.md)
    * [Version-specific features](earthfile/features.md)
* [The `earthly` command](earthly-command/earthly-command.md)
* [Earthly lib](earthly-lib/earthly-lib.md)
* [Configuration reference](earthly-config/earthly-config.md)
* [Examples](examples/examples.md)
* Misc
    * [Alternative installation](./alt-installation/alt-installation.md)
    * [Data collection](data-collection/data-collection.md)
    * [Definitions](definitions/definitions.md)
    * [Public key authentication](public-key-auth/public-key-auth.md)

## 🔧 CI Integration

* [Overview](ci-integration/overview.md)
* [Use the Earthly CI Image](ci-integration/use-earthly-ci-image.md)
* [Build your own Earthly CI Image](ci-integration/build-an-earthly-ci-image.md)
* [Pull-Through Cache](ci-integration/pull-through-cache.md)
* [Remote BuildKit](ci-integration/remote-buildkit.md)
* Vendor-Specific Guides
    * [GitHub Actions](ci-integration/guides/gh-actions-integration.md)
    * [Circle CI](ci-integration/guides/circle-integration.md)
    * [GitLab CI/CD](ci-integration/guides/gitlab-integration.md)
    * [Jenkins](ci-integration/guides/jenkins.md)
    * [AWS CodeBuild](ci-integration/guides/codebuild-integration.md)
    * [Google Cloud Build](ci-integration/guides/google-cloud-build.md)
    * [Bitbucket Pipelines](ci-integration/guides/bitbucket-pipelines-integration.md)
    * [Woodpecker CI](ci-integration/guides/woodpecker-integration.md)
    * [Kubernetes](ci-integration/guides/kubernetes.md)

## ☁️ Earthly Cloud

* [Overview](cloud/overview.md)
* [Managing permissions](cloud/managing-permissions.md)
* [Cloud secrets](cloud/cloud-secrets.md)
* [Earthly Satellites](cloud/satellites.md)
    * [Managing Satellites](cloud/satellites/managing.md)
    * [Using Satellites](cloud/satellites/using.md)
    * [Self-Hosted Satellites](cloud/satellites/self-hosted.md)
    * [GitHub runners](cloud/satellites/gha-runners.md)
    * [Best Practices](cloud/satellites/best-practices.md)
    * [Bring Your Own Cloud (BYOC)](cloud/satellites/byoc/byoc.md)
      * AWS
        * [Requirements](cloud/satellites/byoc/aws/requirements.md)
        * [CloudFormation](cloud/satellites/byoc/aws/cloudformation.md)
        * [Terraform](cloud/satellites/byoc/aws/terraform.md)
        * [Manual](cloud/satellites/byoc/aws/manual.md)
      * VPN
        * [Tailscale](cloud/satellites/byoc/vpn/tailscale.md)
