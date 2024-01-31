#!/bin/sh
set -eo pipefail # DONT add a set -x or you will leak the key

acbtest -n "$USER1_TOKEN"
acbtest -n "$USER2_TOKEN"
acbtest -n "$USER2_SSH_KEY"
acbtest -z "$SSH_AUTH_SOCK"

echo "== it should login to user1 with token =="
EARTHLY_TOKEN="$USER1_TOKEN" earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user1@earthly.dev" using token auth'

echo "== it should stay logged in as user1 even though EARTHLY_TOKEN is no longer set =="
earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user1@earthly.dev" using cached jwt auth'

echo "== it should stay logged in as user1 since the cached jwt is used (even though user2's ssh key is available via ssh keys) =="
eval "$(ssh-agent)"
echo "$USER2_SSH_KEY" | ssh-add -
ssh-add -l | acbgrep '(ED25519)'

earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user1@earthly.dev" using cached jwt auth'

ssh-add -D # remove the key

echo "== forcing a logout should allow us to change users =="
earthly account logout
EARTHLY_TOKEN="$USER2_TOKEN" earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user2@earthly.dev" using token auth'

echo "== it should stay logged in as user2 =="
earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user2@earthly.dev" using cached jwt auth'

echo "== it should be able to login as user2 with ssh =="
earthly account logout
echo "$USER2_SSH_KEY" | ssh-add -
earthly account login 2>&1 | acbgrep 'Logged in as "other-service.earthly-user2@earthly.dev" using ssh auth'

echo "== using token param should behave similarly to EARTHLY_TOKEN env =="
earthly account login --token "$USER2_TOKEN" 2>&1 | acbgrep 'Logged in as "other-service.earthly-user2@earthly.dev" using token auth'

echo "== same as above but first ensure we're logged out =="
earthly account logout
rm -vf ~/.earthly/auth.*
earthly account login --token "$USER2_TOKEN" 2>&1 | acbgrep 'Logged in as "other-service.earthly-user2@earthly.dev" using token auth'

echo "== ensure auth files are recreated =="
acbtest -f ~/.earthly/auth.credentials
acbtest -f ~/.earthly/auth.jwt
