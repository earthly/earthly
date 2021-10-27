# Earthly rpm repository

We host a rpm repository which fedora and CentOS users can use to install earthly.

## Setup for Fedora

TODO: move these notes elsewhere, this readme should only be notes on how to release to our repo, and is only intended for those with
access to earthly credentials.

fedora users can use this guide to set up our repo:



First install the following tools:

    sudo dnf -y install dnf-plugins-core

Second, add earthly's repo

    dnf config-manager \
        --add-repo \
        https://pkg.earthly.dev/earthly.repo

Finally, install earthly:

    dnf -y install earthly

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

    earthly --push +build-and-release --RELEASE_TAG="$RELEASE_TAG"

### Running steps independently

It is also possible to run steps independently:

#### Building rpm packages

To package all platforms

    earthly +rpm-all --RELEASE_TAG="$RELEASE_TAG"

To package a specific platform

    earthly +rpm --RELEASE_TAG --EARTHLY_PLATFORM=arm7

#### Cloning the s3 repo to your local disk

    earthly +download

#### Indexing and signing the repo

    earthly +index-and-sign

#### Uploading the repo to s3

    earthly --push +upload
