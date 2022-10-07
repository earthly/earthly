#!/bin/bash

# clean up old images (best effort)
docker images | grep myuser/earthly-multiplatform-wait-test | awk '{print $1 ":" $2}' | xargs -n 1 docker rmi

set -e
cd "$(dirname "$0")"
../common/test.sh
