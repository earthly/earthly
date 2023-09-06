# Pushing and Pulling images with GCP Artifact Registry

## Introduction

The GCP Artifact Registry is a hosted docker repository that requires extra configuration for day-to-day use. This configuration is not typical of other repositories, and there are some considerations to account for when using it with Earthly. This guide will walk you through creating an Earthfile, building an image, and pushing it to Artifact Registry.

[Artifact Registry is the successor to the GCP Container Registry (GCR)](https://cloud.google.com/artifact-registry/docs/transition/transition-from-gcr). It can accommodate more than just Docker images, but those are beyond the scope of this guide. Most of what we detail here applies to GCR as well, it will just require some [small tweaks](https://cloud.google.com/artifact-registry/docs/transition/transition-from-gcr#compare).

This guide assumes you have already installed the [`gcloud` CLI tool](https://cloud.google.com/sdk/docs/install), [enabled the Artifact Repository API](https://console.cloud.google.com/flows/enableapi?apiid=artifactregistry.googleapis.com&redirect=https://cloud.google.com/artifact-registry/docs/docker/quickstart), and [created a new repository named `hello-earthly`](https://console.cloud.google.com/artifacts).

## Create an Earthfile

No special considerations are needed in the Earthfile itself. You can use `SAVE IMAGE` just like any other repository.

```
FROM alpine:3.18

build:
    RUN echo "Hello from Earthly!" > motd
    ENTRYPOINT cat motd
    SAVE IMAGE --push <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
```

## Configure the Artifact Repository Credential Helper

Artifact Repository does not issue permanent credentials. Instead, it relies on your Google credentials to issue Docker credentials. To this end, Google has built-in a credential helper to the `gcloud` CLI tool. `gcloud` can update your `.docker/config.json` on its own by running `gcloud auth configure-docker <region>-docker.pkg.dev`. Here is a sample entry it might create:

```
 {
    "credHelpers": {
        "<region>-docker.pkg.dev": "gcloud"
  }
}

```

## IAM

Ensure that you have correct permissions to push and pull the images. Please reference the [GCP documentation](https://cloud.google.com/artifact-registry/docs/access-control#grant) to ensure you have the correct permissions set. You will need to add the `Artifact Registry Reader` and `Artifact Registry Writer` roles to complete the tasks in this guide.

If you are using GCR; keep in mind that the needed permissions are based on the GCP storage permissions. We used the `Storage Admin` permissions to complete the guide with GCR.

Service Accounts also work with Earthly. Rather than `gcloud init`, simply log in using the Google-provided key like this:

```
RUN gcloud auth activate-service-account --key-file /test/key.json
```

## Run the Target

With the helper installed, no special To build and push an image, simply execute the build target. Don't forget the `--push` flag!

```
❯ earthly --push +build
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
         alpine:3.18 | --> Load metadata linux/amd64
               +base | --> FROM alpine:3.18
               +base | [██████████] resolve docker.io/library/alpine:3.18@sha256:0bd0e9e03a022c3b0226667621da84fc9bf562a9056130424b5bfbd8bcb0397f ... 100%
              +build | --> RUN echo "Hello from Earthly!" > motd
              output | --> exporting outputs
              output | [██████████] exporting layers ... 100%
              output | [██████████] exporting manifest sha256:08f310b4520418a60f7c12b168167ea22b886bc03d43ab87058e959ef5c14cf2 ... 100%
              output | [██████████] exporting config sha256:8a54361d584a6a51f0136b9ae1526aba8f99cc0a1583954b0f206d3a472eaac9 ... 100%
              output | [██████████] pushing layers ... 100%
              output | [██████████] pushing manifest for <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love ... 100%
              output | [██████████] sending tarballs ... 100%
2bc1eb057e55: Loading layer [==================================================>]     187B/187B
=========================== SUCCESS ===========================
Loaded image: <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
              +build | Image +build as <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love (pushed)


```

## Pulling Images

Using this credential helper; you can also pull images without any special handling in an Earthfile:

```
FROM earthly/dind:alpine-main

run:
    WITH DOCKER --pull <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
        RUN docker run <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
    END
```

And here is how you would run it:

```
❯ earthly -P +run
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
  e/dind:alpine-main | --> Load metadata linux/amd64
u/e/h/hello-earthly:with-love | --> Load metadata linux/amd64
u/e/h/hello-earthly:with-love | --> DOCKER PULL <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
u/e/h/hello-earthly:with-love | [          ] resolve <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love@sha256:08f310b4520418a60f7c12b168167ea22b886bc03d43ab87058e959ef5c14cf2 ... 0%                               [██████████] resolve <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love@sha256:08f310b4520418a60f7c12b168167ea22b886bc03d43ab87058e959ef5c14cf2 ... 100%
               +base | --> FROM earthly/dind:alpine-main
               +base | [██████████] resolve docker.io/earthly/dind:alpine-main@sha256:09f497f0114de1f3ac6ce2da05568fcb50b0a4fd8b9025ed7c67dc952d092766 ... 100%
                +run | *cached* --> WITH DOCKER (install deps)
                +run | *cached* --> WITH DOCKER RUN docker run <region>-docker.pkg.dev/<project>/hello-earthly/hello-earthly:with-love
              output | --> exporting outputs
              output | [██████████] sending tarballs ... 100%
=========================== SUCCESS ===========================

```
