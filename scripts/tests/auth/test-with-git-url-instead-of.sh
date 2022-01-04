#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

gitsha=$(git rev-parse HEAD)
test -n "$gitsha"

docker image rm -f earthly/examples:cpp
GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:" "$earthly" -VD "github.com/earthly/earthly/examples/cpp:$gitsha+docker"
docker run --rm earthly/examples:cpp | grep fib
