#!/bin/bash
set -uxe
set -o pipefail

cd "$(dirname "$0")"
earthly=${earthly-"../../build/linux/amd64/earthly"}

cp ~/.earthly/config.yml ~/.earthly/config.yml.bkup

function finish {
  mv ~/.earthly/config.yml.bkup ~/.earthly/config.yml
}
trap finish EXIT


"$earthly" config global.tls_enabled true

# FIXME bootstrap is failing with "open /home/runner/.earthly/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap || (echo "ignoring bootstrap failure")

# bootstrapping should generate six pem files
test $(ls ~/.earthly/certs/*.pem | wc -l) = "6"

"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'
