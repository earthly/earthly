#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key

earthly=${earthly:=earthly}
if [ "$earthly" != "earthly" ]; then
  earthly=$(realpath "$earthly")
fi
echo "running tests with $earthly"
"$earthly" --version

PATH="$(realpath "$(dirname "$0")/../acbtest"):$PATH"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)

EARTHLY_INSTALLATION_NAME="earthly.integration"
export EARTHLY_INSTALLATION_NAME
rm -rf "$HOME/.earthly.integration/"

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# test logout has no effect when EARTHLY_TOKEN is set
if GITHUB_ACTIONS="" NO_COLOR=0 "$earthly" account logout > output 2>&1; then
    echo "earthly account logout should have failed"
    exit 1
fi
diff output <(echo "Error: account logout has no effect when --auth-token (or the EARTHLY_TOKEN environment variable) is set")

# fetch shared secret key (this step assumes your personal user has access to the /earthly-technologies/ secrets org
echo "fetching manitou-id_rsa"
ID_RSA=$("$earthly" secrets --org earthly-technologies --project core get -n secrets-integration-manitou-id_rsa)

# now that we grabbed the manitou credentials, unset our token, to ensure that we're only testing using manitou's credentials
unset EARTHLY_TOKEN
"$earthly" account logout

echo starting new instance of ssh-agent, and loading credentials
eval "$(ssh-agent)"

# grab first 6chars of md5sum of key to help sanity check that the same key is consistently used
set +x # make sure we don't print the key here
md5sum=$(echo -n "$ID_RSA" | md5sum | awk '{ print $1 }' | head -c6)

echo "Adding key (with md5sum $md5sum...) into ssh-agent"
echo "$ID_RSA" | ssh-add -

echo testing that key was correctly loaded into ssh-agent
ssh-add -l | acbgrep manitou

echo testing that the ssh-agent only contains a single key
test "$(ssh-add -l | wc -l)" = "1"

echo "testing earthly account login works (and is using the earthly-manitou account)"
"$earthly" account login 2>&1 | acbgrep 'Logged in as "other-service+earthly-manitou@earthly.dev" using ssh auth'

mkdir -p /tmp/earthtest
cat << EOF > /tmp/earthtest/Earthfile
VERSION 0.7
PROJECT manitou-org/earthly-core-integration-test
FROM alpine:3.18
test-local-secret:
    WORKDIR /test
    RUN --mount=type=secret,target=/tmp/test_file,id=my_secret test "\$(cat /tmp/test_file)" = "my-local-value"
test-server-secret:
    WORKDIR /test
    RUN --mount=type=secret,target=/tmp/test_file,id=my_test_file test "\$(cat /tmp/test_file)" = "secret-value"
EOF

# set and test get returns the correct value
"$earthly" secrets --org manitou-org --project earthly-core-integration-test set my_test_file "secret-value"
"$earthly" secrets --org manitou-org --project earthly-core-integration-test get my_test_file | acbgrep 'secret-value'

# test earthly will prompt if value is missing
/usr/bin/expect -c '
spawn '"$earthly"' secrets --org manitou-org --project earthly-core-integration-test set my_test_file
expect "secret value: "
send "its my secret value\n"
expect eof
'
"$earthly" secrets --org manitou-org --project earthly-core-integration-test get my_test_file | acbgrep 'its my secret value'

# test set --stdin works
echo -e "hello\nworld" | "$earthly" secrets --org manitou-org --project earthly-core-integration-test set --stdin my_test_file
# note "echo -e "hello\nworld" | md5sum" -> 0f723ae7f9bf07744445e93ac5595156
"$earthly" secrets --org manitou-org --project earthly-core-integration-test get -n my_test_file
"$earthly" secrets --org manitou-org --project earthly-core-integration-test get -n my_test_file | md5sum | acbgrep '0f723ae7f9bf07744445e93ac5595156'

# test set --file works
"$earthly" secrets --org manitou-org --project earthly-core-integration-test set --file <(echo -e "foo\nbar") my_test_file
# note "echo -e "foo\nbar" | md5sum" -> f47c75614087a8dd938ba4acff252494
"$earthly" secrets --org manitou-org --project earthly-core-integration-test get -n my_test_file | md5sum | acbgrep 'f47c75614087a8dd938ba4acff252494'


# restore the "secret-value", which the org selection test requires
"$earthly" secrets --org manitou-org --project earthly-core-integration-test set my_test_file "secret-value"

# test selecting org
"$earthly" org select manitou-org
"$earthly" org ls | acbgrep '^\* \+manitou-org'

# test secrets with org selected in config file
"$earthly" secrets --project earthly-core-integration-test get my_test_file | acbgrep 'secret-value'
"$earthly" secrets --project earthly-core-integration-test set my_other_file "super-secret-value"
"$earthly" secrets --project earthly-core-integration-test get my_other_file | acbgrep 'super-secret-value'
"$earthly" secrets --project earthly-core-integration-test ls | acbgrep '^my_test_file$'

# test secrets with personal org
"$earthly" org select user:other-service+earthly-manitou@earthly.dev
"$earthly" secrets set super/secret hello
"$earthly" secrets get super/secret | acbgrep 'hello'
"$earthly" secrets get /user/super/secret | acbgrep 'hello'
"$earthly" secrets ls | acbgrep '^super/secret$'
"$earthly" secrets ls /user | acbgrep '^super/secret$'

echo "=== test 1 ==="
# test RUN --mount can reference a secret from the command line
"$earthly" --no-cache --secret my_secret=my-local-value /tmp/earthtest+test-local-secret

echo "=== test 2 ==="
# test RUN --mount can reference a secret from the server that is only specified in the Earthfile
"$earthly" --no-cache /tmp/earthtest+test-server-secret

echo "=== test 3 ==="
# Test earthly will display a message containing the name of the secret that was not found
set +e
"$earthly" --no-cache /tmp/earthtest+test-local-secret > output 2>&1
exit_code="$?"
set -e
cat output
test "$exit_code" != "0"
acbgrep 'unable to lookup secret "my_secret": not found' output
acbgrep 'Help: Make sure to set the project at the top of the Earthfile' output
echo "=== All tests have passed ==="
