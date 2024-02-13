#!/usr/bin/bash
set -e
cd /home/alex/gh/earthly/buildkit
sha="$(git rev-parse HEAD)"
cd /home/alex/gh/earthly/earthly
./earthly +update-buildkit --BUILDKIT_GIT_SHA="$sha"
./earthly +for-linux

rm /tmp/earth.log || true
./build/linux/amd64/earthly -P --no-cache +test-no-qemu-group1 |& tee /tmp/earth.log
