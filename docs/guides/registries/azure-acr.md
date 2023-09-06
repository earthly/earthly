# Pushing and Pulling images with Azure ACR

## Introduction

The Azure Container Registry (ACR) is a hosted docker repository that requires extra configuration for day-to-day use. This configuration is not typical of other repositories, and there are some considerations to account for when using it with Earthly. This guide will walk you through creating an Earthfile, building an image, and pushing it to ACR.


This guide assumes you have already installed the [Azure CLI tool](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli), and [created a new repository named `helloearthly`](https://portal.azure.com/?quickstart=true#create/Microsoft.ContainerRegistry).

## Create an Earthfile

No special considerations are needed in the Earthfile itself. You can use `SAVE IMAGE` just like any other repository.

```
FROM alpine:3.18

build:
    RUN echo "Hello from Earthly!" > motd
    ENTRYPOINT cat motd
    SAVE IMAGE --push helloearthly.azurecr.io/hello-earthly:with-love
```

## Login and Configure the ACR Credential Helper

ACR does not issue permanent credentials. Instead, it relies on your Azure AD credentials to issue Docker credentials. As an individual user, you will need to log into your repository first:

```
❯ az acr login --name helloearthly
Login Succeeded
```

After logging in, the [ACR Credential Helper](https://github.com/Azure/acr-docker-credential-helper) will help keep your credentials up to date, as long as it is invoked again before your already issued credentials expire. When all this is complete, your `.docker/config.json` might look like this:
```
{
	"auths": {
		"helloearthly.azurecr.io": {
			"auth": "...",
			"identitytoken": "..."
		}
	},
	"credsStore": "acr-linux"
}
```

ACR boasts many other methods of logging in, including [Service Principals](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-auth-service-principal) and [admin accounts](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-authentication#admin-account). Note that the admin account method is not recommended for production usage. Please follow the relevant guides to authenticate if you wish to use one of these other methods.

## RBAC

Ensure that you have correct permissions to push and pull the images. Please reference the [ACR RBAC documentation](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-roles) to ensure you have the correct permissions set. To complete all the activities in this guide, you will need to have at least the `AcrPush` role.

Earthly also works with Service Principals; and these do not require `az acr login`. You can simply login directly with `docker` like this: 

```
RUN --secret AZ_USERNAME=earthly-technologies/azure/ci-cd-username \
    --secret AZ_PASSWORD=earthly-technologies/azure/ci-cd-password \
    docker login helloearthly.azurecr.io --username $AZ_USERNAME --password $AZ_PASSWORD
```

## Run the Target

Once you are logged in, and have the optional credential helper installed, then you are ready to use Earthly to access images in ACR. To build and push an image, simply execute the build target. Don't forget the `--push` flag!

```
❯ ../earthly/earthly --push --no-cache +build
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
         alpine:3.18 | --> Load metadata linux/amd64
               +base | --> FROM alpine:3.18
               +base | [██████████] resolve docker.io/library/alpine:3.18@sha256:0bd0e9e03a022c3b0226667621da84fc9bf562a9056130424b5bfbd8bcb0397f ... 100%
              +build | --> RUN echo "Hello from Earthly!" > motd
              output | --> exporting outputs
              output | [██████████] exporting layers ... 100%
              output | [██████████] exporting manifest sha256:02df2d4600094d5550f7475b868ce9bb17d6c3a529e9669a453bbba7b2cdb659 ... 100%
              output | [██████████] exporting config sha256:722368416f5de51291ce937feac2c246d66dff351678968b1b6ebc533ceaaa0c ... 100%
              output | [██████████] pushing layers ... 100%
              output | [██████████] pushing manifest for helloearthly.azurecr.io/hello-earthly:with-love ... 100%
              output | [██████████] sending tarballs ... 100%
824d26cf8432: Loading layer [==================================================>]     192B/192B
=========================== SUCCESS ===========================
Loaded image: helloearthly.azurecr.io/hello-earthly:with-love
              +build | Image +build as helloearthly.azurecr.io/hello-earthly:with-love (pushed)
```

## Pulling Images

By logging in and optionally installing the credential helper; you can also pull images without any special handling in an Earthfile:

```
FROM earthly/dind:alpine-main

run:
    WITH DOCKER --pull helloearthly.azurecr.io/hello-earthly:with-love
        RUN docker run helloearthly.azurecr.io/hello-earthly:with-love
    END
```

And here is how you would run it:

```
❯ earthly -P +run
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
  e/dind:alpine-main | --> Load metadata linux/amd64
h/hello-earthly:with-love | --> Load metadata linux/amd64
h/hello-earthly:with-love | --> DOCKER PULL helloearthly.azurecr.io/hello-earthly:with-love
h/hello-earthly:with-love | [██████████] resolve helloearthly.azurecr.io/hello-earthly:with-love@sha256:02df2d4600094d5550f7475b868ce9bb17d6c3a529e9669a453bbba7b2cdb659 ... 100%
               +base | --> FROM earthly/dind:alpine-main
               +base | [██████████] resolve docker.io/earthly/dind:alpine-main@sha256:09f497f0114de1f3ac6ce2da05568fcb50b0a4fd8b9025ed7c67dc952d092766 ... 100%
                +run | *cached* --> WITH DOCKER (install deps)
                +run | --> WITH DOCKER RUN docker run helloearthly.azurecr.io/hello-earthly:with-love
                +run | Loading images...
                +run | Loaded image: helloearthly.azurecr.io/hello-earthly:with-love
                +run | ...done
                +run | Hello from Earthly!
              output | --> exporting outputs
              output | [██████████] sending tarballs ... 100%
=========================== SUCCESS ===========================
```

## Troubleshooting

### 401 (authentication required)

Re-run `az acr login --name` to log in again and refresh your credentials. Azure recommends that you run this at the beginning o each automated script; keep this in mind for your CI runs.
