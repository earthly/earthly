#!/bin/sh
set -eo pipefail # DONT add a set -x or you will leak the key

acbtest -n "$OIDC_USER_TOKEN"
acbtest -n "$earthly_config" # set by earthly-entrypoint.sh


echo "== it should login to user with token =="
# todo: remove address env vars when production is ready
EARTHLY_SERVER_ADDRESS=https://ci.staging.earthly.dev EARTHLY_GRPC_ADDRESS=ci.staging.earthly.dev:443 \
EARTHLY_TOKEN="$OIDC_USER_TOKEN" earthly account login 2>&1 | acbgrep 'Logged in as "ido+testemail@earthly.dev" using token auth'

echo "== it should access aws via oidc =="
# todo: remove address env vars when production is ready
EARTHLY_SERVER_ADDRESS=https://ci.staging.earthly.dev EARTHLY_GRPC_ADDRESS=ci.staging.earthly.dev:443 \
earthly --config "$earthly_config" +oidc
