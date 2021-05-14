#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

echo "machine github.com
login cinnamonthecat
password $REPO_ACCESS_TOKEN" > ~/.netrc

docker image rm -f test-private:latest
$earthly -VD github.com/cinnamonthecat/test-private:main+docker
docker run --rm test-private:latest | grep "Salut Lume"
