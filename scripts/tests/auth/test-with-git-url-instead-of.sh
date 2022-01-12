#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

if [ "$GITHUB_EVENT_NAME" == "pull_request" ]; then
    # GHA does a pre-merge before running tests, this per-merge is not pushed to the repo
    # until after the PR is merged. Therefore running it against this pre-merged, non-pushed
    # commit will fail and should be avoided.
    gitsha=$(jq -r .pull_request.head.sha < "$GITHUB_EVENT_PATH")
else
    gitsha=$(git rev-parse HEAD)
fi
test -n "$gitsha"

GITHUB_SERVER_URL="${VARIABLE:-https://github.com}"
GITHUB_REPOSITORY="${GITHUB_REPOSITORY:-earthly/earthly}"
github_server_without_protocol=$(echo "$GITHUB_SERVER_URL" | sed -e 's#^https\?://##; s#//$##')
earthly_ref="$github_server_without_protocol/$GITHUB_REPOSITORY/examples/cpp:$gitsha+docker"

docker image rm -f earthly/examples:cpp
echo "running $earthly -VD $earthly_ref"
GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:" "$earthly" -VD "$earthly_ref"
docker run --rm earthly/examples:cpp | grep fib
