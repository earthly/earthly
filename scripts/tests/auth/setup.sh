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

# fetch shared secret key (this step assumes your personal user has access to the /earthly-technologies/ secrets org
ID_RSA=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/id_rsa)
GITHUB_PASSWORD=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/password)
REPO_ACCESS_TOKEN=$($earthly secrets get -n /earthly-technologies/github/other-service+github-cinnamon@earthly.dev/repo-access-token)

# now that we grabbed the cinnamonthecat credentials, unset our token, to ensure that we're only testing using cinnamonthecat's credentials
unset EARTHLY_TOKEN

# make sure ssh-agent is not running
test -z "${SSH_AUTH_SOCK:-}"

# make sure tests start without a config
rm -f ~/.earthly/config.yml
