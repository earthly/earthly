#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key

# This script can be run locally with:
# earthly=./build/linux/amd64/earthly ./scripts/tests/private-repo.sh

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

# prevent the self-update of earthly from running (this ensures no bogus data is printed to stdout,
# which would mess with the secrets data being fetched)
date +%s > /tmp/last-earthly-prerelease-check

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# fetch shared secret key (this step assumes your personal user has access to the /earthly-technologies/ secrets org
ID_RSA=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/id_rsa)
GITHUB_PASSWORD=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/password)

# now that we grabbed the cinnamonthecat credentials, unset our token, to ensure that we're only testing using cinnamonthecat's credentials
unset EARTHLY_TOKEN

echo starting new instance of ssh-agent, and loading credentials
eval "$(ssh-agent)"

# make sure that no keys are initially loaded
ssh-add -l | grep 'The agent has no identities'


# test we can still clone a public repo
echo === Test 1 ===
"$earthly" github.com/earthly/hello-world:main+hello


# test we can still clone a public repo even if we have an unauthorized key-loaded
echo === Test 2 ===

# first create a new random key which is NOT authorized for any services
ssh-keygen -b 3072 -t rsa -f /tmp/sshkey -q -N "" -C "this-key-is-not-authorized"
ssh-add /tmp/sshkey

# then test that we still have access to this public repo
"$earthly" github.com/earthly/hello-world:main+hello


# Test a private repo can be cloned
echo === Test 3 ===

# grab first 6chars of md5sum of key to help sanity check that the same key is consistently used
md5sum=$(echo -n "$ID_RSA" | md5sum | awk '{ print $1 }' | head -c6)

echo "Adding key (with md5sum $md5sum...) into ssh-agent"
echo "$ID_RSA" | ssh-add -

echo testing that key was correctly loaded into ssh-agent
ssh-add -l | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

echo testing github connection
ssh git@github.com 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

docker image rm -f test-private:latest
"$earthly" -VD github.com/cinnamonthecat/test-private:main+docker
docker run --rm test-private:latest | grep "Salut Lume"


# Test public repo can be cloned without ssh, when GIT_URL_INSTEAD_OF is set as recommended by our CI docs
echo === Test 4 ===
docker image rm -f earthly/examples:cpp
SSH_AUTH_SOCK="" GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:" $earthly -VD github.com/earthly/earthly/examples/cpp:main+docker
docker run --rm earthly/examples:cpp | grep fib


# Test a private repo can be cloned using https
echo === Test 5 ===

cat << EOF > /tmp/earthconfig.https
git:
  github.com:
    auth: https
    user: cinnamonthecat
    password: "$GITHUB_PASSWORD"
EOF
cat /tmp/earthconfig.https

docker image rm -f other-test-private:latest
SSH_AUTH_SOCK="" $earthly -VD --config /tmp/earthconfig.https github.com/cinnamonthecat/other-test-private:main+docker
docker run --rm other-test-private:latest | grep "Salut le monde"


# Test a private repo can be cloned using https, as setup via args
echo === Test 6 ===

docker image rm -f other-test-private:latest
SSH_AUTH_SOCK="" $earthly -VD --git-username=cinnamonthecat --git-password="$GITHUB_PASSWORD" github.com/cinnamonthecat/another-test-private:main+docker
docker run --rm another-test-private:latest | grep "Hola Mundo"


echo === All private-repo.sh tests have passed ===
