#!/bin/bash
set -euo pipefail # don't use -x as it will leak the mirror credentials

# to run this locally; in the root of the repo:
#   ./earthly +earthly-docker && EARTHLY_IMAGE="earthly/earthly:dev-$(git rev-parse --abbrev-ref HEAD | sed 's/\//_/g')" scripts/tests/earthly-image-sat.sh

FRONTEND=${FRONTEND:-docker}
EARTHLY_IMAGE=${EARTHLY_IMAGE:-earthly/earthly:dev-main}
PATH="$(realpath "$(dirname "$0")/../acbtest"):$PATH"

if [ -z "${DOCKERHUB_MIRROR_USERNAME:-}" ] && [ -z "${DOCKERHUB_MIRROR_PASSWORD:-}" ]; then
  echo "using dockerhub credentials from earthly secrets"
  DOCKERHUB_MIRROR_USERNAME="$(earthly secrets --org earthly-technologies --project core get dockerhub-mirror/user)"
  export DOCKERHUB_MIRROR_USERNAME
  DOCKERHUB_MIRROR_PASSWORD="$(earthly secrets --org earthly-technologies --project core get dockerhub-mirror/pass)"
  export DOCKERHUB_MIRROR_PASSWORD
fi
test -n "$DOCKERHUB_MIRROR_USERNAME" || (echo "error: DOCKERHUB_MIRROR_USERNAME is not set" && exit 1)
test -n "$DOCKERHUB_MIRROR_PASSWORD" || (echo "error: DOCKERHUB_MIRROR_PASSWORD is not set" && exit 1)

if [ -z "${EARTHLY_TOKEN:-}" ]; then
  echo "using EARTHLY_TOKEN from earthly secrets"
  EARTHLY_TOKEN="$(earthly secrets --org earthly-technologies --project core get earthly-token-for-satellite-tests)"
  export EARTHLY_TOKEN
fi
test -n "$EARTHLY_TOKEN" || (echo "error: EARTHLY_TOKEN is not set" && exit 1)

function finish {
  status="$?"
  if [ "$status" = "0" ]; then
    echo "earthly-image-sat.sh test passed"
  else
    echo "earthly-image-sat.sh failed with $status"
  fi
}
trap finish EXIT

# Wait for a stopped satellite to finish shutting down. Other states may continue.
while true; do
  state=$("$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" "${EARTHLY_IMAGE}" sat --org earthly-technologies inspect core-test 2>&1 | grep State | awk '{print $2}')
  echo "Current state: $state"
  case $state in
    Stopping)
      echo "Waiting for shutdown to finish"
      ;;
    *)
      break
      ;;
  esac
  sleep 10
done

# This will catch the case where the satellite was previously stopping. No effect when already awake.
"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" "${EARTHLY_IMAGE}" sat --org earthly-technologies wake core-test

echo "Test earthly sat inspect."
"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" "${EARTHLY_IMAGE}" sat --org earthly-technologies inspect core-test 2>&1 | tee output.txt
acbgrep "core-test" output.txt
"$FRONTEND" run --rm -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" sat --org earthly-technologies inspect core-test 2>&1 | tee output.txt
acbgrep "core-test" output.txt

echo "Test satellite (not privileged, no buildkit)."
satconfig="$(mktemp /tmp/earthly-image-test-satellite-config.XXXXXX)"
chmod 600 "$satconfig"
if [ -n "${DOCKERHUB_USERNAME:-}" ] && [ -n "${DOCKERHUB_PASSWORD:-}" ]; then
    ENCODED_SAT_AUTH="$(echo -n "$DOCKERHUB_USERNAME:$DOCKERHUB_PASSWORD" | base64 -w 0)"

    cat > "$satconfig" <<EOF
{
	"auths": {
		"docker.io": {
			"auth": "$ENCODED_SAT_AUTH"
		}
	}
}
EOF
fi

# TODO FIXME test that the above credentials are actually used by the satellite

"$FRONTEND" run --rm -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" -v "$satconfig:/root/.docker/config.json" -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" --org earthly-technologies --sat core-test --no-cache github.com/earthly/hello-world:4d466d524f768a379374c785fdef30470e87721d+hello 2>&1 | tee output.txt
acbgrep "Hello World" output.txt
acbgrep "Earthly installation is working correctly" output.txt

rm output.txt
echo "=== All tests have passed ==="
