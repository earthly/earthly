#!/bin/bash
set -euo pipefail # don't use -x as it will leak the mirror credentials

FRONTEND=${FRONTEND:-docker}
EARTHLY_IMAGE=${EARTHLY_IMAGE:-earthly/earthly:dev-main}

test -n "$DOCKERHUB_MIRROR_USERNAME" || (echo "error: DOCKERHUB_MIRROR_USERNAME is not set" && exit 1)
test -n "$DOCKERHUB_MIRROR_PASSWORD" || (echo "error: DOCKERHUB_MIRROR_PASSWORD is not set" && exit 1)

ENCODED_AUTH="$(echo -n "$DOCKERHUB_MIRROR_USERNAME:$DOCKERHUB_MIRROR_PASSWORD" | base64 -w 0)"

dockerconfig="$(mktemp /tmp/earthly-image-test-docker-config.XXXXXX)"
chmod 600 "$dockerconfig"
cat > "$dockerconfig" <<EOF
{
	"auths": {
		"registry-1.docker.io.mirror.corp.earthly.dev": {
			"auth": "$ENCODED_AUTH"
		}
	}
}
EOF

# Note that it is not possible to use GLOBAL_CONFIG for this, due to the fact
# earthly-entrypoint.sh starts buildkit instead of the earthly binary,
# as a result the buildkit_additional_config value in ~/.earthly/config.yml is ignored.
export EARTHLY_ADDITIONAL_BUILDKIT_CONFIG='[registry."docker.io"]
  mirrors = ["registry-1.docker.io.mirror.corp.earthly.dev"]'

function finish {
  status="$?"
  if [ "$status" = "0" ]; then
    echo "earthly-image.sh test passed"
  else
    echo "earthly-image.sh failed with $status"
  fi
  rm "$dockerconfig"
}
trap finish EXIT

echo "Test no --privileged and no NO_BUILDKIT=1 -> fail."
if "$FRONTEND" run --rm "${EARTHLY_IMAGE}" 2>&1 | tee output.txt; then
    echo "expected failure"
    exit 1
fi
grep "Container appears to be running unprivileged" output.txt

echo "Test no target provided -> fail."
if "$FRONTEND" run --rm --privileged "${EARTHLY_IMAGE}" 2>&1 | tee output.txt; then
    echo "expected failure"
    exit 1
fi
grep "Executes Earthly builds" output.txt # Display help
grep "Error: no target reference provided" output.txt # Show error
if "$FRONTEND" run --rm -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" 2>&1 | tee output.txt; then
    echo "expected failure"
    exit 1
fi
grep "Executes Earthly builds" output.txt # Display help
grep "Error: no target reference provided" output.txt # Show error

echo "Test --version (smoke test)."
"$FRONTEND" run --rm --privileged "${EARTHLY_IMAGE}" --version 2>&1
"$FRONTEND" run --rm -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" --version 2>&1

echo "Test --help."
"$FRONTEND" run --rm --privileged "${EARTHLY_IMAGE}" --help 2>&1 | tee output.txt
grep "Executes Earthly builds" output.txt # Display help
"$FRONTEND" run --rm -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" --help 2>&1 | tee output.txt
grep "Executes Earthly builds" output.txt # Display help

echo "Test earthly sat inspect."
"$FRONTEND" run --rm --privileged -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" "${EARTHLY_IMAGE}" sat --org earthly-technologies inspect core-test 2>&1 | tee output.txt
grep "core-test" output.txt
"$FRONTEND" run --rm -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" sat --org earthly-technologies inspect core-test 2>&1 | tee output.txt
grep "core-test" output.txt

echo "Test hello world with embedded buildkit."
"$FRONTEND" run --rm --privileged -e EARTHLY_ADDITIONAL_BUILDKIT_CONFIG -v "$dockerconfig:/root/.docker/config.json" "${EARTHLY_IMAGE}" --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
grep "Hello World" output.txt
grep "Earthly installation is working correctly" output.txt

if [ "$FRONTEND" = "docker" ]; then
    echo "Test use /var/run/docker.sock, but not privileged."
    "$FRONTEND" run --rm -e EARTHLY_ADDITIONAL_BUILDKIT_CONFIG -v "$dockerconfig:/root/.docker/config.json" -e NO_BUILDKIT=1 -e EARTHLY_NO_BUILDKIT_UPDATE=1 -v /var/run/docker.sock:/var/run/docker.sock "${EARTHLY_IMAGE}" --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
    grep "Hello World" output.txt
    grep "Earthly installation is working correctly" output.txt
fi

echo "Test satellite (not privileged, no buildkit)."
"$FRONTEND" run --rm -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" --org earthly-technologies --sat core-test --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
grep "Hello World" output.txt
grep "Earthly installation is working correctly" output.txt

rm output.txt
echo "Success"
