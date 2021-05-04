#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

# start ssh-agent and make sure no keys are loaded
eval "$(ssh-agent)"
ssh-add -l | grep 'The agent has no identities'

# display the first 6 characters of the md5sum of the key
md5sum=$(echo -n "$ID_RSA" | md5sum | awk '{ print $1 }' | head -c6)
echo "Adding key (with md5sum $md5sum...) into ssh-agent"

# load the key into the ssh auth agent
echo "$ID_RSA" | ssh-add -

echo testing that key was correctly loaded into ssh-agent
ssh-add -l | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

echo testing github connection
ssh git@github.com 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /cinnamonthecat/;'

docker image rm -f test-private:latest
"$earthly" -VD github.com/cinnamonthecat/test-private:main+docker
docker run --rm test-private:latest | grep "Salut Lume"
