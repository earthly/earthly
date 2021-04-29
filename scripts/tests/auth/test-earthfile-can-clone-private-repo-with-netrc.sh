#!/usr/bin/env bash
set -eu # don't use -x as it will leak the private key
# shellcheck source=./setup.sh
source "$(dirname "$0")/setup.sh"

echo "machine github.com
login cinnamonthecat
password $REPO_ACCESS_TOKEN" > ~/.netrc

mkdir /tmp/earthly-68686b7e-ff3e-4c68-86d9-cc84458b1bab
cat << EOF > /tmp/earthly-68686b7e-ff3e-4c68-86d9-cc84458b1bab/Earthfile
FROM alpine:3.13
test-clone-https:
    GIT CLONE --branch main https://github.com/cinnamonthecat/test-private.git .
    RUN cat README.md
test-clone-ssh:
    GIT CLONE --branch main git@github.com:cinnamonthecat/test-private.git .
    RUN cat README.md
EOF

"$earthly" -VD /tmp/earthly-68686b7e-ff3e-4c68-86d9-cc84458b1bab+test-clone-https
"$earthly" -VD /tmp/earthly-68686b7e-ff3e-4c68-86d9-cc84458b1bab+test-clone-ssh
