#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

rm .testdata || true # cleanup

set +e
"$earthly" -P $@ +test
exit_code="$?"
set -e

test -f .testdata
cat .testdata | grep '^CONTAINER ID'
test "$(md5sum .testdata | awk '{print $1}')" = "8abf6fc0ebbfaac48c70d23d2dce27fd"
test "$exit_code" -ne "0"
