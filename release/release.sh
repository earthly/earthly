#!/bin/bash
set -ex

# earthly does not (yet) have fine control over controlling the order of RUN --push commands (they all happen at the end)
# Our release process requires the following commands be done in order:
# - uploading binaries to github (as a --push command)
# - pushing a new commit to the earthly/earthly-homebrew repo
# - downloading the binaries from github and existing binaries from s3 buckets
# - signing those apt and yum packages containing those binaries
# - and pushing signed up to s3 buckets (which back our apt and yum repos)

# Args
#   Required
#     RELEASE_TAG
#   Optional
#     GITHUB_USER         override the default earthly github user
#     DOCKERHUB_USER      override the default earthly dockerhub user
#     DOCKERHUB_IMG       override the default earthly dockerhub image name
#     EARTHLY_REPO        override the default earthly repo name
#     BREW_REPO           override the default homebrew-earthly repo name
#     GITHUB_SECRET_PATH  override the default +secrets/earthly-technologies/github/littleredcorvette/token secrets location where the github token is stored
#     PRERELEASE          override the default false value (must be false or true)
#
# examples
#
#  performing a test release to a non-earthly location
#    env -i HOME="$HOME" PATH="$PATH" SSH_AUTH_SOCK="$SSH_AUTH_SOCK" RELEASE_TAG=v0.5.10 GITHUB_USER=alexcb DOCKERHUB_USER=alexcb132 EARTHLY_REPO=earthly BREW_REPO=homebrew-earthly-1 ./release.sh
#
#  performing a release candidate
#    env -i HOME="$HOME" PATH="$PATH" SSH_AUTH_SOCK="$SSH_AUTH_SOCK" RELEASE_TAG=v0.6.0-rc1 PRERELEASE=true ./release.sh
#
#  performing a regular release
#    env -i HOME="$HOME" PATH="$PATH" SSH_AUTH_SOCK="$SSH_AUTH_SOCK" RELEASE_TAG=v0.6.0 ./release.sh
#

# must happen before we change dirs
if [[ "$earthly" == .* ]]; then
  earthly="$(pwd)/$earthly"
fi

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "$SCRIPT_DIR"

test -n "$HOME" || (echo "ERROR: HOME is not set"; exit 1);
test -n "$RELEASE_TAG" || (echo "ERROR: RELEASE_TAG is not set" && exit 1);
(echo "$RELEASE_TAG" | grep '^v[0-9]\+.[0-9]\+.[0-9]\+\(-rc[0-9]\+\)\?$' > /dev/null) || (echo "ERROR: RELEASE_TAG must be formatted as v1.2.3 (or v1.2.3-RC1); instead got \"$RELEASE_TAG\""; exit 1);
command -v jq || (echo "ERROR: jq is not installed"; exit 1);

if [ -z "$earthly" ]; then
  ../earthly upgrade
  earthly="../earthly"
fi

test -f "$earthly" || (echo "ERROR: earthly is set to $earthly which does not exist" && exit 1)

# TODO once v 0.7 is fully released, we can remove this
if ! "$earthly" secrets --help 2>&1 | grep migrate > /dev/null; then
    echo "you are using an older version of earthly, please upgrade to v0.7.X (or build it from main)"
    exit 1
fi

# Set default values
export GITHUB_USER=${GITHUB_USER:-earthly}
export DOCKERHUB_HOST=${DOCKERHUB_HOST:-docker.io}
export DOCKERHUB_USER=${DOCKERHUB_USER:-earthly}
export DOCKERHUB_IMG=${DOCKERHUB_IMG:-earthly}
export DOCKERHUB_BUILDKIT_IMG=${DOCKERHUB_BUILDKIT_IMG:-buildkitd}
export EARTHLY_REPO=${EARTHLY_REPO:-earthly}
export BREW_REPO=${BREW_REPO:-homebrew-earthly}
export GITHUB_SECRET_PATH=$GITHUB_SECRET_PATH
export PRERELEASE=${PRERELEASE:-true}
export SKIP_CHANGELOG_DATE_TEST=${SKIP_CHANGELOG_DATE_TEST:-false}
export S3_BUCKET=${S3_BUCKET:-production-pkg}
export EARTHLY_STAGING=${EARTHLY_STAGING:-false}


if [ "$PRERELEASE" != "false" ] && [ "$PRERELEASE" != "true" ]; then
    echo "PRERELEASE must be \"true\" or \"false\""
    exit 1
fi

if [[ "$RELEASE_TAG" =~ "rc" ]] && [ "$PRERELEASE" != "true" ]; then
    echo "unable to have release candidate \"$RELEASE_TAG\" when PRERELEASE is not true (got \"$PRERELEASE\" instead)"
    exit 1
fi

PRODUCTION_RELEASE="false"
if [ "$GITHUB_USER" = "earthly" ] && [ "$EARTHLY_REPO" = "earthly" ]; then
    PRODUCTION_RELEASE="true"
    if [ "$DOCKERHUB_HOST" != "docker.io" ]; then
        echo "expected DOCKERHUB_HOST=docker.io but got $DOCKERHUB_HOST"
        exit 1
    fi
    if [ "$DOCKERHUB_USER" != "earthly" ]; then
        echo "expected DOCKERHUB_USER=earthly but got $DOCKERHUB_USER"
        exit 1
    fi
    if [ "$DOCKERHUB_IMG" != "earthly" ]; then
        echo "expected DOCKERHUB_IMG=earthly but got $DOCKERHUB_IMG"
        exit 1
    fi
    if [ "$DOCKERHUB_BUILDKIT_IMG" != "buildkitd" ]; then
        echo "expected DOCKERHUB_BUILDKIT_IMG=buildkitd but got $DOCKERHUB_BUILDKIT_IMG"
        exit 1
    fi
    if [ "$S3_BUCKET" != "production-pkg" ]; then
        echo "expected S3_BUCKET=production-pkg but got $S3_BUCKET"
        exit 1
    fi
else
    if [ "$S3_BUCKET" == "production-pkg" ]; then
        echo "unable to perform non-production release with S3_BUCKET set to production-pkg"
        exit 1
    fi
fi

# fail-fast if release-notes do not exist (or if date is incorrect)
"$earthly" --build-arg RELEASE_TAG --build-arg SKIP_CHANGELOG_DATE_TEST +release-notes

if [ -n "$GITHUB_SECRET_PATH" ]; then
    GITHUB_SECRET_PATH_BUILD_ARG="--build-arg GITHUB_SECRET_PATH=$GITHUB_SECRET_PATH"
else
    ("$earthly" secrets --org earthly-technologies --project core ls >/dev/null) || (echo "ERROR: current user does not have access to the earthly-technologies core project"; exit 1);
fi

existing_release=$(curl -s https://api.github.com/repos/earthly/earthly/releases/tags/$RELEASE_TAG | jq -r .tag_name)
if [ "$existing_release" != "null" ]; then
    test "$OVERWRITE_RELEASE" = "1" || (echo "a release for $RELEASE_TAG already exists, to proceed with overwriting this release set OVERWRITE_RELEASE=1" && exit 1);
    echo "overwriting existing release for $RELEASE_TAG"
fi

PUSH_LATEST_TAG="true"
if [ "$PRERELEASE" = "true" ] || [ "$PRODUCTION_RELEASE" != "true" ]; then
    PUSH_LATEST_TAG="false"
fi

GITHUB_PRERELEASE="$PRERELEASE"
if [ "$EARTHLY_STAGING" = "true" ]; then
    # special case to ensure https://github.com/earthly/earthly-staging/releases/latest/download/earthly-linux-amd64 is kept up to date
    GITHUB_PRERELEASE="false"

    # make sure we aren't accidentally doing a regular release
    if [ "$PUSH_LATEST_TAG" = "true" ]; then
        echo "something is wrong; PUSH_LATEST_TAG should be false"
        exit 1
    fi
fi

earthlynext="$(cat ../earthly-next)"
if [[ ! "$earthlynext" =~ ^[a-zA-Z0-9]{40}$ ]]; then
    echo "../earthly-next does not contain a valid git commit hash; got $earthlynext"
    exit 1
fi
echo "earthlynext is $earthlynext"

"$earthly" --push --build-arg DOCKERHUB_USER --build-arg DOCKERHUB_IMG --build-arg DOCKERHUB_BUILDKIT_IMG +release-dockerhub --PUSH_PRERELEASE_TAG="$PRERELEASE" --PUSH_LATEST_TAG="$PUSH_LATEST_TAG" --RELEASE_TAG="$RELEASE_TAG"
"$earthly" --push --build-arg DOCKERHUB_USER --build-arg DOCKERHUB_IMG --build-arg DOCKERHUB_BUILDKIT_IMG +perform-release-buildkitd-dockerhub --PUSH_PRERELEASE_TAG="false" --PUSH_LATEST_TAG="false" --RELEASE_TAG="$RELEASE_TAG-ticktock" --BUILDKIT_PROJECT=github.com/earthly/buildkit:$earthlynext
"$earthly" --push --build-arg GITHUB_USER --build-arg EARTHLY_REPO --build-arg BREW_REPO --build-arg DOCKERHUB_USER --build-arg DOCKERHUB_BUILDKIT_IMG --build-arg RELEASE_TAG --build-arg SKIP_CHANGELOG_DATE_TEST $GITHUB_SECRET_PATH_BUILD_ARG +release-github --PRERELEASE="$GITHUB_PRERELEASE"

if [ "$PRERELEASE" != "false" ]; then
    echo "exiting due to PRERELEASE=$PRERELEASE"
    exit 0
fi

if [ "$EARTHLY_STAGING" = "true" ]; then
    echo "exiting due to EARTHLY_STAGING=$EARTHLY_STAGING"
    exit 0
fi

echo "homebrew release with gu=$GITHUB_USER; er=$EARTHLY_REPO; br=$BREW_REPO; du=$DOCKERHUB_USER; rt=$RELEASE_TAG"
"$earthly" --push --build-arg GITHUB_USER --build-arg EARTHLY_REPO --build-arg BREW_REPO --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG $GITHUB_SECRET_PATH_BUILD_ARG +release-homebrew

if [ "$PRODUCTION_RELEASE" = "true" ]; then
    echo "releasing to apt under s3://$S3_BUCKET/deb"
    "$earthly" --push --build-arg RELEASE_TAG --build-arg S3_BUCKET ./apt-repo+build-and-release

    echo "releasing to yum under s3://$S3_BUCKET/yum"
    "$earthly" --push --build-arg RELEASE_TAG --build-arg S3_BUCKET ./yum-repo+build-and-release
else
    echo "staging apt and yum repos are currently disabled" // TODO s3 related release needs to be updated to work with MFA
fi
