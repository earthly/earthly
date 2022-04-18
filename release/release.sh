#!/bin/bash
set -ex

# earthly does not (yet) have fine control over controlling the order of RUN --push commands (they all happen at the end)
# Our release process requires the following commands be done in order:
# - uploading binaries to github (as a --push command)
# - pushing a new commit to the earthly/earthly-homebrew repo
# - downloading the binaries from github and existing binaries from s3 buckets
# - signing those apt and yum packages containing those binaries
# - and pushing signed up to s3 buckets (which back our apt and yum repos)
# - building an AMI for the latest version

# Args
#   Required
#     RELEASE_TAG
#   Optional
#     GITHUB_USER         override the default earthly github user
#     DOCKERHUB_USER      override the default earthly dockerhub user
#     EARTHLY_REPO        override the default earthly repo name
#     BREW_REPO           override the default homebrew-earthly repo name
#     GITHUB_SECRET_PATH  override the default +secrets/earthly-technologies/github/griswoldthecat/token secrets location where the github token is stored
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

if [[ "$earthly" == .* ]]; then
  earthly="$(pwd)/$earthly"
fi

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd $SCRIPT_DIR

test -n "$HOME" || (echo "ERROR: HOME is not set"; exit 1);

test -n "$RELEASE_TAG" || (echo "ERROR: RELEASE_TAG is not set" && exit 1);
(echo "$RELEASE_TAG" | grep '^v[0-9]\+.[0-9]\+.[0-9]\+\(-rc[0-9]\+\)\?$' > /dev/null) || (echo "ERROR: RELEASE_TAG must be formatted as v1.2.3 (or v1.2.3-RC1); instead got \"$RELEASE_TAG\""; exit 1);

# Set default values
export GITHUB_USER=${GITHUB_USER:-earthly}
export DOCKERHUB_USER=${DOCKERHUB_USER:-earthly}
export EARTHLY_REPO=${EARTHLY_REPO:-earthly}
export BREW_REPO=${BREW_REPO:-homebrew-earthly}
export GITHUB_SECRET_PATH=$GITHUB_SECRET_PATH
export PRERELEASE=${PRERELEASE:-false}
export SKIP_CHANGELOG_DATE_TEST=${SKIP_CHANGELOG_DATE_TEST:-false}


if [ "$PRERELEASE" != "false" ] && [ "$PRERELEASE" != "true" ]; then
    echo "PRERELEASE must be \"true\" or \"false\""
    exit 1
fi

if [ "$GITHUB_USER" = "earthly" ] && [ "$EARTHLY_REPO" = "earthly" ]; then
    if [ "$DOCKERHUB_USER" != "earthly" ]; then
        echo "expected DOCKERHUB_USER=earthly but got $DOCKERHUB_USER"
        exit 1
    fi
fi

if [ -z "$earthly" ]; then
  ../earthly upgrade
  earthly="../earthly"
fi

# fail-fast if release-notes do not exist (or if date is incorrect)
"$earthly" --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG --build-arg SKIP_CHANGELOG_DATE_TEST +release-notes

if [ -n "$GITHUB_SECRET_PATH" ]; then
    GITHUB_SECRET_PATH_BUILD_ARG="--build-arg GITHUB_SECRET_PATH=$GITHUB_SECRET_PATH"
else
    ("$earthly" secrets ls /earthly-technologies >/dev/null) || (echo "ERROR: current user does not have access to earthly-technologies shared secrets"; exit 1);
fi

release_apt_and_yum="false"
release_ami="false"
if [ "$GITHUB_USER" = "earthly" ] && [ "$EARTHLY_REPO" = "earthly" ]; then
    ("$earthly" secrets get /user/earthly-technologies/aws/credentials >/dev/null) || (echo "ERROR: user-secrets /user/earthly-technologies/aws/credentials does not exist"; exit 1);
    release_apt_and_yum="true"
    release_ami="true"
fi

existing_release=$(curl -s https://api.github.com/repos/earthly/earthly/releases/tags/$RELEASE_TAG | jq -r .tag_name)
if [ "$existing_release" != "null" ]; then
    test "$OVERWRITE_RELEASE" = "1" || (echo "a release for $RELEASE_TAG already exists, to proceed with overwriting this release set OVERWRITE_RELEASE=1" && exit 1);
    echo "overwriting existing release for $RELEASE_TAG"
fi

"$earthly" --push --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG +release-dockerhub
"$earthly" --push --build-arg GITHUB_USER --build-arg EARTHLY_REPO --build-arg BREW_REPO --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG --build-arg SKIP_CHANGELOG_DATE_TEST --build-arg PRERELEASE="$PRERELEASE" $GITHUB_SECRET_PATH_BUILD_ARG +release-github

if [ "$PRERELEASE" != "false" ]; then
    echo "exiting due to prerelease = true"
    exit 0
fi

echo "homebrew release with gu=$GITHUB_USER; er=$EARTHLY_REPO; br=$BREW_REPO; du=$DOCKERHUB_USER; rt=$RELEASE_TAG"
"$earthly" --push --build-arg GITHUB_USER --build-arg EARTHLY_REPO --build-arg BREW_REPO --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG $GITHUB_SECRET_PATH_BUILD_ARG +release-homebrew

# TODO pass along a RELEASE_REPO_TEST_SUFFIX which would cause us to host our yum/apt repos under https://test-pkg.earthly.dev/$RELEASE_REPO_TEST_SUFFIX/...
# and when it is empty, we would use https://pkg.earthly.dev/...
#../earthly --push --build-arg GITHUB_USER --build-arg EARTHLY_REPO --build-arg BREW_REPO --build-arg DOCKERHUB_USER --build-arg RELEASE_TAG +release-repo
# until then, we will just print this out:
echo "TODO: the apt/yum release must be triggered seperately; until we get https://test-pkg.earthly.dev/ setup"

if [ "$release_apt_and_yum" = "true" ]; then
    "$earthly" --push --build-arg RELEASE_TAG ./apt-repo+build-and-release
    "$earthly" --push --build-arg RELEASE_TAG ./yum-repo+build-and-release
else
    echo "WARNING: there is no staging environment for apt or yum repos"
fi

if [ "$release_ami" = "true" ]; then
  "$earthly" --push --build-arg RELEASE_TAG ./ami+update-pipelines
  "$earthly" --push --build-arg RELEASE_TAG ./ami+build-ami
else
  echo "WARNING: there is no staging environment for AMIs"
fi
