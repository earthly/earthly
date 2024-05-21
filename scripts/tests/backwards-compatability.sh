#!/usr/bin/env bash
set -xeu

# used to start earthly-integration-buildkitd
earthly="${earthly:=earthly}"
if [ -f "$earthly" ]; then
  earthly=$(realpath "$earthly")
fi

# used for testing backwards compatability issues
crustly="${crustly:=earthly-v0.8.0}"
if [ -f "$crustly" ]; then
  crustly=$(realpath "$crustly")
fi

# change directory to script location
cd -- "$( dirname -- "${BASH_SOURCE[0]}" )"

current_git_sha="$(git rev-parse HEAD)"

if "$("$crustly" --version)" | grep "$current_git_sha" >/dev/null; then
  echo "ERROR: $crustly was built using the current git sha $current_git_sha"
  exit 1
fi

echo "running tests with earthly=$earthly for bootstrapping and crustly=$crustly for cli"
echo "earthly=$("$earthly" --version)"
echo "crustly=$("$crustly" --version)"
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)
echo "using frontend=$frontend"

PATH="$(realpath ../acbtest):$PATH"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

set +x # dont remove or the token will be leaked
if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)
set -x

EARTHLY_INSTALLATION_NAME="earthly-integration"
export EARTHLY_INSTALLATION_NAME
rm -rf "$HOME/.earthly.integration/"

echo "$earthly"
# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# start buildkitd container
"$earthly" bootstrap

# start a build using an older version of the earthly cli
"$crustly" --no-buildkit-update -P ../../tests/with-docker+all

# validate buildkitd container was compiled using the current branch
buildkitd_earthly_version="$(docker logs earthly-integration-buildkitd |& grep -o 'EARTHLY_GIT_HASH=[a-z0-9]*')"
acbtest "$buildkitd_earthly_version" = "EARTHLY_GIT_HASH=$current_git_sha"

echo "=== All tests have passed ==="
