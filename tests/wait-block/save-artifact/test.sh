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
    echo "save-artifact test passed"
  else
    echo "save-artifact test failed with $status"
  fi
}
trap finish EXIT

# Cleanup from previous tests
rm -f data

"$earthly" $@ +test
test "$(cat data)" = "foo"

# next, check for an expected failure
set +e
("$earthly" $@ +test-fail; echo $? > earthly.exitcode) 2>&1 | tee earthly.log
set -e
test "$(cat earthly.exitcode)" = "1"
grep 'unable to copy file data, which has is outputted elsewhere' earthly.log

if grep "this magic string should never appear" earthly.log >/dev/null; then
  echo "magic string command should never have run, but did"
  exit 1
fi
