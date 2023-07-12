#!/usr/bin/env bash
set -euo pipefail

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# test --org
"$earthly" sat --org earthly-technologies inspect core-test

# test EARTHLY_ORG env
EARTHLY_ORG=earthly-technologies "$earthly" sat inspect core-test

# test org select
"$earthly" org select earthly-technologies
"$earthly" sat select core-test
