#!/usr/bin/env bash

set -ue
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../build/linux/amd64/earthly"}
# docker or podman
frontend=${frontend-"docker"}

echo "=== Test Single --push ==="

"$earthly" --allow-privileged --no-cache +only-push 2>&1 | tee single_output

if ! cat single_output | grep "Did not push image earthly/sap:only-push"; then
    echo "Invalid push text (expected: Did not push image earthly/sap:only-push)"
    cat single_output
    exit 1
fi

"$frontend" images --format "{{.Repository}}:{{.Tag}}" | grep "earthly/sap:" > single_images

if  ! cat single_images | grep -q "earthly/sap:only-push"; then
    echo "Missed only valid image"
    cat single_images
    exit 1
fi

echo "=== Test All Phase Transitions ==="

"$earthly" --allow-privileged --no-cache +test 2>&1 | tee multi_output

if ! cat multi_output | grep "Did not push image earthly/sap:after-push"; then
    echo "Invalid push text (expected: Did not push image earthly/sap:after-push)"
    cat multi_output
    exit 1
fi

"$frontend" images --format "{{.Repository}}:{{.Tag}}" | grep "earthly/sap:" > multi_images

if echo "$EARTHLY_VERSION_FLAG_OVERRIDES" | grep "wait-block" >/dev/null; then
    echo "skipping non-output after push test; --wait-block feature does not impose such a limitation"
else
    if ! cat multi_output | grep "Did not output image earthly/sap:after-push"; then
        echo "Invalid push text (expected: Did not output image earthly/sap:after-push)"
        cat multi_output
        exit 1
    fi
    if cat multi_images | grep "earthly/sap:after-push"; then
        echo "Saved invalid image"
        "$frontend" images --format "{{.Repository}}:{{.Tag}}" | grep "sap:"
        exit 1
    fi
fi

if  ! cat multi_images | grep -q "earthly/sap:empty" || \
    ! cat multi_images | grep -q "earthly/sap:before-push" || \
    ! cat multi_images | grep -q "earthly/sap:first-push"
then
    echo "Missed valid image"
    cat multi_images
    exit 1
fi

echo "Success!"
