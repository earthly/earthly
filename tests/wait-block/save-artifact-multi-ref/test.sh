#!/usr/bin/env bash
set -uex
set -o pipefail

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
"$earthly" --version

# display a pass/fail message at the end
function finish {
  status="$?"
  if [ "$status" = "0" ]; then
    echo "save-artifact-multi-ref test passed"
  else
    echo "save-artifact-multi-ref test failed with $status"
  fi
}
trap finish EXIT

# Cleanup from previous tests
rm -f data

"$earthly" $@ +test
test "$(cat data)" = "162948536bc2"
