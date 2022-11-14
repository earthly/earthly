#!/bin/bash
set -eu

if [ -z ${GITHUB_ACTIONS+x} ]; then
    echo "this script should only be run from GHA; if run locally it will modify your ssh settings and earthly config"
    exit 1
fi

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

# ensure earthly login works (and print out who gets logged in)
"$earthly" account login

# these tests require the EARTHLY_TOKEN not be set
unset EARTHLY_TOKEN

# make sure ssh-agent is not running
test -z "${SSH_AUTH_SOCK:-}"

# make sure tests start without a config
rm -f ~/.earthly-dev/config.yml
