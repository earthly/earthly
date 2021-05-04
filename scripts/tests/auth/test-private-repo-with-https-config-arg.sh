#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

docker image rm -f other-test-private:latest
SSH_AUTH_SOCK="" $earthly -VD --git-username=cinnamonthecat --git-password="$GITHUB_PASSWORD" github.com/cinnamonthecat/another-test-private:main+docker
docker run --rm another-test-private:latest | grep "Hola Mundo"
