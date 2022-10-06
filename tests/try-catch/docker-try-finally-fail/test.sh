#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)
echo "using frontend=$frontend"

rm .testdata || true # cleanup

set +e
"$earthly" -P $@ +test --FRONTEND=$frontend
exit_code="$?"
set -e

test -f .testdata
cat .testdata | grep '^CONTAINER ID'
test "$(md5sum .testdata | awk '{print $1}')" = "8abf6fc0ebbfaac48c70d23d2dce27fd"
test "$exit_code" -ne "0"
