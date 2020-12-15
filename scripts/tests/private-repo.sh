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

echo testing github connection
ssh git@github.com 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

# Test a private repo can be cloned
echo === Test 1 ===
docker image rm -f test-private:latest
$earthly -VD github.com/cinnamonthecat/test-private:main+docker
docker run --rm test-private:latest | grep "Salut Lume"

# Test public repo can be cloned without ssh, when GIT_URL_INSTEAD_OF is set as recommended by our CI docs
echo === Test 2 ===
docker image rm -f cpp-example:latest
SSH_AUTH_SOCK="" GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:" $earthly -VD github.com/earthly/earthly/examples/cpp:main+docker
docker run --rm cpp-example:latest | grep fib

# Test a private repo can be cloned using https
echo === Test 3 ===

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
echo === Test 4 ===

docker image rm -f other-test-private:latest
SSH_AUTH_SOCK="" $earthly -VD --git-username=cinnamonthecat --git-password="$GITHUB_PASSWORD" github.com/cinnamonthecat/another-test-private:main+docker
docker run --rm another-test-private:latest | grep "Hola Mundo"


echo === All private-repo.sh tests have passed ===
