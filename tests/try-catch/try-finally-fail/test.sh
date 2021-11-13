#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

rm .testdata || true # cleanup

set +e
"$earthly" $@ +test
exit_code="$?"
set -e

test -f .testdata
test "$(cat .testdata)" = "magic"
test "$(md5sum .testdata | awk '{print $1}')" = "b6e7f4684fae52db1c9e33687b697d38"
test "$exit_code" -ne "0"
