#!/bin/bash

set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# Test.
"$earthly" --allow-privileged \
    --ci \
    +all
exit "$?"