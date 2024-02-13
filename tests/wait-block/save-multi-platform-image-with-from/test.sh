#!/usr/bin/env bash
set -uex

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

# clean up old images (best effort)
docker images | grep earthly-multiplatform-wait-test-with-from | awk '{print $1 ":" $2}' | xargs -r -n 1 docker rmi

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
"$earthly" +test
