#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

rm .testdata || true # cleanup

"$earthly" $@ +test
! test -f .testdata
test -f .otherdata

! "$earthly" $@ +test --fail=yes
! test -f .testdata
test -f .otherdata
