#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

cleanup() {
  rm -r .testdata data out/ || true
}

execute() {
  local target="$1"
  local expected_path="$2"
  shift 2

  cleanup

  set +e
  "$earthly" $@ "$target"
  exit_code="$?"
  set -e

  test -f "$expected_path"
  test "$(cat "$expected_path")" = "magic"
  test "$(md5sum "$expected_path" | awk '{print $1}')" = "b6e7f4684fae52db1c9e33687b697d38"
  test "$exit_code" -ne "0"
}

execute '+test' .testdata "$@"
execute '+test-save-to-curdir' data "$@"
execute '+test-save-to-child-dir' out/data "$@"
execute '+test-save-to-child-file' out/.testdata "$@"
cleanup
