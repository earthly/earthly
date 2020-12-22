#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key

earthly=${earthly:=earthly}
echo "running tests with $earthly"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

# ensure earthly login works (and print out who gets logged in)
$earthly account login

# fetch shared secret key (this step assumes your personal user has access to the /earthly-technologies/ secrets org
ID_RSA=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/id_rsa)
GITHUB_PASSWORD=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/password)

# now that we grabbed the cinnamonthecat credentials, unset our token, to ensure that we're only testing using cinnamonthecat's credentials
unset EARTHLY_TOKEN

echo starting new instance of ssh-agent, and loading credentials
eval "$(ssh-agent)"

# grab first 6chars of md5sum of key to help sanity check that the same key is consistently used
md5sum=$(echo -n "$ID_RSA" | md5sum | awk '{ print $1 }' | head -c6)

echo "Adding key (with md5sum $md5sum...) into ssh-agent"
echo "$ID_RSA" | ssh-add -

echo testing that key was correctly loaded into ssh-agent
ssh-add -l | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

echo testing that the ssh-agent only contains a single key
test "$(ssh-add -l | wc -l)" = "1"

echo "testing earthly account login works (and is using the earthly-cinnamon account)"
$earthly account login | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /other-service\+earthly-cinnamon\@earthly.dev/;'

mkdir -p /tmp/earthtest
cat << EOF > /tmp/earthtest/Earthfile
FROM alpine:3.11
test-local-secret:
    WORKDIR /test
    RUN --mount=type=secret,target=/tmp/test_file,id=+secrets/my_secret test "\$(cat /tmp/test_file)" = "secret-value"
test-server-secret:
    WORKDIR /test
    RUN --mount=type=secret,target=/tmp/test_file,id=+secrets/user/earthly_integration_tests/my_test_file test "\$(cat /tmp/test_file)" = "secret-value"
EOF

# set and test get returns the correct value
$earthly secrets set /user/earthly_integration_tests/my_test_file "secret-value"
$earthly secrets get /user/earthly_integration_tests/my_test_file | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /secret-value/;'

echo === test 1 ===
# test RUN --mount can reference a secret from the command line
$earthly --no-cache --secret my_secret=secret-value /tmp/earthtest+test-local-secret

echo === test 2 ===
# test RUN --mount can reference a secret from the server that is only specified in the Earthfile
$earthly --no-cache /tmp/earthtest+test-server-secret

echo === All tests have passed ===

