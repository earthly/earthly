#!/usr/bin/env bash
set -uex
set -o pipefail

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
"$earthly" --version

docker rmi myimg:623cb5fb1b8c4cff8693281095724bb0 || true
"$earthly" +copy
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

"$earthly" --no-output +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "0"

"$earthly" +build
test "$(docker images -q myimg:623cb5fb1b8c4cff8693281095724bb0 | wc -l)" = "1"
