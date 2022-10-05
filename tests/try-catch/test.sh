#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly)
set -ue
set -o pipefail

initialwd="$(pwd)"
cd "$(dirname "$0")"

# display a pass/fail message at the end
function finish {
  status="$?"
  RED='\033[0;31m'
  GREEN='\033[0;32m'
  NC='\033[0m' # No Color
  if [ "$status" = "0" ]; then
    printf "${GREEN}try-catch tests passed${NC}\n"
  else
    printf "${RED}try-catch tests failed with ${status}${NC}\n"
  fi
}
trap finish EXIT

# TODO: add back docker-try-finally-fail
for test_path in try-catch-not-currently-implemented try-finally-fail try-finally-pass try-finally-if-exists try-finally-two-files
do
    printf "=== running $test_path ===\n\n"
    "${test_path}/test.sh"
    printf "${test_path} passed\n\n"
done
