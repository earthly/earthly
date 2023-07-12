#!/bin/bash
set -euo pipefail # don't use -x as it will leak the mirror credentials

# to run this locally; in the root of the repo:
#   ./earthly +earthly-docker && EARTHLY_IMAGE="earthly/earthly:dev-$(git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')" scripts/tests/satellites-integration.sh

FRONTEND=${FRONTEND:-docker}
EARTHLY_IMAGE=${EARTHLY_IMAGE:-earthly/earthly:dev-main}

if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)

"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" --entrypoint="" "${EARTHLY_IMAGE}" /bin/sh -c "earthly org select earthly-technologies"

"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" --entrypoint="" "${EARTHLY_IMAGE}" /bin/sh -c "earthly sat ls"

"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" --entrypoint="" "${EARTHLY_IMAGE}" /bin/sh -c "earthly sat select core-test"

"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" --entrypoint="" "${EARTHLY_IMAGE}" /bin/sh -c "earthly sat inspect core-test"

"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" --entrypoint="" "${EARTHLY_IMAGE}" /bin/sh -c "earthly github.com/earthly/hello+hello-world"
