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

# first test we can do a regular build
"$earthly" +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "1"
docker rmi myimg:623cb5fb1b8c4cff8693281095724bb0

# copy shouldn't produce an image
"$earthly" +copy
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

# --no-output should prevent outputting images
"$earthly" --no-output +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

# --image mode only ouputs image of directly-referenced image,
# in the case of +build, there is no SAVE IMAGE
"$earthly" --image +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

# the +myimg target on the otherhand contains an explicit SAVE IMAGE
"$earthly" --image +myimg
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "1"
docker rmi myimg:623cb5fb1b8c4cff8693281095724bb0
