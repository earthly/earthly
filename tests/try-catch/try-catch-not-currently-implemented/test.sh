#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

rm .testdata || true # cleanup

! "$earthly" $@ +test 2>&1 | tee .earthlyoutput
! test -f .testdata

grep "TRY/FINALLY doesn't (currently) support CATCH statements" .earthlyoutput
