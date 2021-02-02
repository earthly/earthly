#!/usr/bin/env bash

set -ue
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

"$earthly" --allow-privileged --no-cache +test | tee output

if ! cat output | grep "Did not push, OR save earthly/sap:after-push locally"; then
    echo "Invalid push text"
    cat output
    exit 1
fi

docker images --format "{{.Repository}}:{{.Tag}}" | grep "earthly/sap:" > images

if cat images | grep "earthly/sap:after-push"; then
    echo "Saved invalid image"
    docker images --format "{{.Repository}}:{{.Tag}}" | grep "sap:"
    exit 1
fi

if  ! cat images | grep -q "earthly/sap:empty" || \
    ! cat images | grep -q "earthly/sap:before-push" || \
    ! cat images | grep -q "earthly/sap:first-push"
then
    echo "Missed valid image"
    docker images --format "{{.Repository}}:{{.Tag}}" | grep "sap:"
    exit 1
fi

echo "Success!"
