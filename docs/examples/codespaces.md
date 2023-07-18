# Github Codespacse

## Overview

Earthly can be run from within configured Codespaces as part of your customized dev environments. It can be added either by:

* Adding the Earthly feature to your devcontainer configuration
* Installing Earthly directly onto your devcontainer Dockerfile

### Compatibility

Earthly has been tested with Codespaces on Github.com and in Visual Stuidio Code.

### Resources

 * [Codespaces Deep Dive](https://docs.github.com/en/codespaces/getting-started/deep-dive)
 * [Codespaces Devcontainers](https://docs.github.com/en/codespaces/setting-up-your-project-for-codespaces/adding-a-dev-container-configuration/introduction-to-dev-containers)
 * [Install Earthly](https://earthly.dev/get-earthly)


## Setup with Satellites (Recommended)

With Satellites you can run Earthly in your devcontainer without needing to install docker or podman! This can allow you to have a more simplified setup as well as faster Earthly build times.

If you don't already have an Earthly account or satellites read more about them here:
- https://docs.earthly.dev/earthly-cloud/satellites

Login here: 
- https://cloud.earthly.dev/login


How To:
1. Update your devcontainer configuration in .devcontainer/devcontainer.json to use a Dockerfile

    ```
    {
        "name": "My Devcontainer",

        "build": {
            "dockerfile": "Dockerfile"
        },
        ...
        // YOUR OTHER CONFIG HERE

        // You can also add on the Earthly syntax highlighting extension!
        "customizations": {
            "vscode": {
                "extensions": [
                    "earthly.earthfile-syntax-highlighting"
                ]
            }
        },
        ...
    }
    ```

1. Update your Dockerfile to include installing Earthly and any other setup you may need
    
    `sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly`
    
    _Note_: You don't need to install docker or podman, nor bootstrap earthly, when using satellites! Allowing you to further simplify your codespaces configuration
1. Rebuild your container to include the installation
1. Login `earthly account login --token {YOUR_TOKEN_HERE}`
    
    For how to handle secrets view Codespaces secret guide here: https://docs.github.com/en/codespaces/managing-your-codespaces/managing-encrypted-secrets-for-your-codespaces
1. Select your org `earthly org select {YOUR_ORG_NAME_HERE}`
1. Launch your satellite: `earthly sat launch {YOUR_SAT_NAME_HERE}`
1. Your Earthly commands will now run on your satellite!
    
    _Note_: It is recommended you login to your docker registry account to prevent rate limiting issues.


## Setup with Docker (Using Devcontainer Feature)

The Feature used in this example includes installing docker so there is no need to install it separately

This example was run using the `mcr.microsoft.com/devcontainers/universal:2` image.

1. Update your devcontainer configuration in .devcontainer/devcontainer.json in the root of your project 
    ````
    // YOUR OTHER CONFIG HERE
    ...
    ...
    "features": {
        "ghcr.io/shepherdjerred/devcontainers-features/earthly:1": {
            "bootstrap": true
        }
    },

    // You can also add on the Earthly syntax highlighting extension!
    "customizations": {
        "vscode": {
            "extensions": [
                "earthly.earthfile-syntax-highlighting"
            ]
        }
    },
    ...
    ...
    ````
1. Run `Rebuild Container` on your codespace to have it install the feature
1. Your Earthly commands will now run on docker!

## Setup with Podman (Manually customizing your Devcontainer)

Currently there are known issues with running Earthly with Podman on an environment that uses cgroup v2 on Codespaces (which runs as a Docker container). You can check what version of cgroup your podman installation is using by running:

```podman info | grep cgroup```

These issues are still being investigated at this time and there is no fix available yet. You can confirm that you are being affected by these issues by checking the logs from the earthly-buildkit container if you see `sh: write error: Invalid argument` before the container exits.  

Related to: https://github.com/containers/podman/issues/12559

For more information on Podman with Devcontainers:
    
- https://code.visualstudio.com/remote/advancedcontainers/docker-options#_podman 
    
    This refers to setting up Podman as a remote container instead of as part of the devcontainer. This is for using devcontainers on your local machine rather than within Codespaces.

