#!/bin/bash

set -euxo pipefail

FRONTEND=${FRONTEND:-docker}
EARTHLY_IMAGE=${EARTHLY_IMAGE:-earthly/earthly:dev-main}

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
"$FRONTEND" run --rm --privileged "${EARTHLY_IMAGE}" --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
grep "Hello World" output.txt
grep "Earthly installation is working correctly" output.txt

if [ "$FRONTEND" = "docker" ]; then
    echo "Test use /var/run/docker.sock, but not privileged."
    "$FRONTEND" run --rm -e NO_BUILDKIT=1 -v /var/run/docker.sock:/var/run/docker.sock "${EARTHLY_IMAGE}" --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
    grep "Hello World" output.txt
    grep "Earthly installation is working correctly" output.txt
fi

echo "Test satellite (not privileged, no buildkit)."
"$FRONTEND" run --rm -e EARTHLY_TOKEN="${EARTHLY_TOKEN}" -e NO_BUILDKIT=1 "${EARTHLY_IMAGE}" --org earthly-technologies --sat core-test --no-cache github.com/earthly/hello-world+hello 2>&1 | tee output.txt
grep "Hello World" output.txt
grep "Earthly installation is working correctly" output.txt

rm output.txt
echo "Success"
