# Earthly Debian repository

We host a Debian repository which Debian and ubuntu users can use to install earthly.

## Setup for Ubuntu

TODO: move these notes elsewhere, this readme should only be notes on how to release to our repo, and is only intended for those with
access to earthly credentials.

Ubuntu users can use this guide to set up our repo:

First install the following tools:

    sudo apt-get update
    sudo apt-get install \
       apt-transport-https \
       ca-certificates \
       curl \
       gnupg \
       lsb-release

Second, add earthly's official GPG key:

    curl -fsSL https://pkg.earthly.dev/earthly.pgp | sudo gpg --dearmor -o /usr/share/keyrings/earthly-archive-keyring.gpg


Finally, set up the stable repository:

    echo \
      "deb [arch=amd64 signed-by=/usr/share/keyrings/earthly-archive-keyring.gpg] https://pkg.earthly.dev/deb \
      stable main" | sudo tee /etc/apt/sources.list.d/earthly.list > /dev/null

## Requirements

To package a new version of earthly, ensure the following requirements are met:

1. you have aws credentials configured in the earthly secret store under `/user/earthly-technologies/aws/credentials`, and have access to the developer role

    # you can upload them via
    earthly secrets set --file ~/.aws/credentials /user/earthly-technologies/aws/credentials

2. you have access to the earthly-technologies secrets; specifically the following two commands should work:

    earthly secrets ls /earthly-technologies/apt/keys/earthly-apt-public.pgp
    earthly secrets ls /earthly-technologies/apt/keys/earthly-apt-private.pgp

## Release steps

Once earthly has been released to GitHub, visit https://github.com/earthly/earthly/releases to determine the latest version:

    export RELEASE_TAG="v0.0.0"

Then run

    earthly +build-and-release --RELEASE_TAG="$RELEASE_TAG"

### Running steps independently

It is also possible to run steps independently:

#### Building deb packages

To package all platforms

    earthly +deb-all --RELEASE_TAG="$RELEASE_TAG"

To package a specific platform

    earthly +deb --RELEASE_TAG="$RELEASE_TAG" --EARTHLY_PLATFORM=arm7

#### Cloning the s3 repo to your local disk

    earthly +download

#### Indexing and signing the repo

    earthly +index-and-sign

#### Uploading the repo to s3

    earthly +upload
