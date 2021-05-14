#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

"$earthly" config git "{github.com: {auth: https, user: 'cinnamonthecat', password: '$GITHUB_PASSWORD'}}"

docker image rm -f other-test-private:latest
"$earthly" -VD github.com/cinnamonthecat/other-test-private:main+docker
docker run --rm other-test-private:latest | grep "Salut le monde"
