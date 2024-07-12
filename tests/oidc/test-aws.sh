#!/bin/sh
set -eo pipefail # DONT add a set -x or you will leak the key

acbtest -n "$OIDC_USER_TOKEN"
acbtest -n "$ROLE_ARN"
acbtest -n "$earthly_config" # set by earthly-entrypoint.sh


echo "== it should login to user with token =="
EARTHLY_TOKEN="$OIDC_USER_TOKEN" earthly account login 2>&1 | acbgrep 'Logged in as "other-service+oidc-ci-test@earthly.dev" using token auth'

echo "== it should access aws via oidc =="
earthly --config "$earthly_config" +oidc --ROLE_ARN="$ROLE_ARN"

echo "== it should access aws via oidc-with-docker =="
earthly --config "$earthly_config" --allow-privileged +oidc-with-docker --ROLE_ARN="$ROLE_ARN"
