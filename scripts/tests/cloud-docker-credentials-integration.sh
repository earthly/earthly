#!/usr/bin/env bash
set -eu

earthly=${earthly:=earthly}
if [ "$earthly" != "earthly" ]; then
  earthly=$(realpath "$earthly")
fi
echo "running tests with $earthly"
"$earthly" --version

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

set +x # dont remove or the token will be leaked
if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)
set -x

EARTHLY_INSTALLATION_NAME="earthly-integration"
export EARTHLY_INSTALLATION_NAME
rm -rf "$HOME/.earthly.integration/"

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# A username / password has been stored in the cloud to a docker hub user (that is not part of earthly) via:
#   earthly secret --org earthly-technologies --project core-test-cloud-docker-credentials-test set std/registry/registry-1.docker.io/username verygoodusername
#   earthly secret --org earthly-technologies --project core-test-cloud-docker-credentials-test set std/registry/registry-1.docker.io/password verygoodpassword # just kidding
#
# And our earthly-user that GHA uses has been configured to have access:
#   earthly projects --org earthly-technologies --project core-test-cloud-docker-credentials-test members add other-service+earthly@earthly.dev read+secrets

# check that secrets is working, and we are running with the correct user
test "$("$earthly" secret --org earthly-technologies --project core-test-cloud-docker-credentials-test get std/registry/registry-1.docker.io/username)" = "verygoodusername"

echo ==== test that the private verygoodimage can be fetched using the credentials from secrets ====
rm -rf /tmp/earthly-cloud-docker-credentials-test-1

mkdir /tmp/earthly-cloud-docker-credentials-test-1
cd /tmp/earthly-cloud-docker-credentials-test-1
cat >> Earthfile <<EOF
VERSION 0.7
PROJECT earthly-technologies/core-test-cloud-docker-credentials-test
test1:
    FROM verygoodusername/verygoodimage:verygoodtag
    RUN base64 -d /my-test-data | grep verygooddata
EOF

# first make sure docker can't access the verygoodimage
test "$(docker images -q verygoodusername/verygoodimage | wc -l)" = "0"
if docker pull verygoodusername/verygoodimage:verygoodtag 2> docker-pull.log; then
  cat docker-pull.log
  echo "error: this test requires that docker does not have access to pull the verygoodimage"
  exit 1
fi
if ! grep 'requested access to the resource is denied' docker-pull.log >/dev/null; then
  cat docker-pull.log
  echo expected denied failed, but got somthing else
  exit 1
fi

# then test that earthly can access the verygoodimage (by using the cloud-hosted registry credentials)
"$earthly" --no-cache +test1

echo "=== All tests have passed ==="
