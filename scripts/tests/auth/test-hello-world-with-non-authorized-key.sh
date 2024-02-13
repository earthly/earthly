#!/bin/bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

# start ssh-agent and make sure no keys are loaded
eval "$(ssh-agent)"
ssh-add -l | grep 'The agent has no identities'

# generate a new random key which is NOT authorized for any services
ssh-keygen -b 3072 -t rsa -f /tmp/non-authorized-ssh-key -q -N "" -C "this-key-is-not-authorized"
ssh-add /tmp/non-authorized-ssh-key

# test that only the above key is loaded
test "$(ssh-add -l | wc -l)" = "1"

# test earthly can access a public repo
"$earthly" github.com/earthly/hello-world:main+hello
