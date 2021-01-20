#!/bin/bash
# Note: Most of this test runs as Earthly-in-Earthly so that we can easily send local cache to a tmpfs.

set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# Test.
set +e
"$earthly" --allow-privileged \
    --ci \
    +all
exit_code="$?"
set -e


exit "$exit_code"
