#!/usr/bin/env bash
set -euo pipefail

earthly=${earthly:=earthly}
if [ "$earthly" != "earthly" ]; then
  earthly=$(realpath "$earthly")
fi
"$earthly" --version
PATH="$(realpath "$(dirname "$0")/../acbtest"):$PATH"

echo "running tests with $earthly"
"$earthly" --version

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)

set -x # don't move this to the top; or we'll leak the token

EARTHLY_INSTALLATION_NAME="earthly-integration"
export EARTHLY_INSTALLATION_NAME
rm -rf "$HOME/.earthly.integration/"

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# test --org
"$earthly" sat --org earthly-technologies inspect core-test

# test EARTHLY_ORG env
EARTHLY_ORG=earthly-technologies "$earthly" sat inspect core-test

# test inspect can use the org name from the config file
"$earthly" org select earthly-technologies
"$earthly" sat inspect core-test

# test the sat select works correctly uses the config's org
"$earthly" sat select core-test

# Print ls for debugging
"$earthly" satellite ls

# Check core-test is selected
NO_COLOR=1 "$earthly" satellite ls | acbgrep '^\* \+core-test'

echo "=== All tests have passed ==="
