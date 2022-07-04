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
    echo "no-output test passed"
  else
    echo "no-output test failed with $status"
  fi
}
trap finish EXIT

# Cleanup from previous tests
docker rmi myimg:623cb5fb1b8c4cff8693281095724bb0 || true

# Do the tests

"$earthly" +copy
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

"$earthly" --no-output +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

# TODO this test is currently failing --image mode does not prevent a BUILD in a WAIT/END from saving the image
#"$earthly" --image +build
#test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

"$earthly" +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "1"
