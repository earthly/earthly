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
    echo "no-cache test passed"
  else
    echo "no-cache test failed with $status"
  fi
}
trap finish EXIT

"$earthly" +test 2>&1 | tee output

alphaline=$(grep -n alpha output | awk -F : '{print $1}')
bravoline=$(grep -n bravo output | awk -F : '{print $1}')

test "$alphaline" -lt "$bravoline"
