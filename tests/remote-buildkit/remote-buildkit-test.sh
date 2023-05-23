#!/bin/bash
set -uxe
set -o pipefail

cd "$(dirname "$0")"
earthly=${earthly-"../../build/linux/amd64/earthly"}
host=$(hostname)

mkdir -p ~/.earthly-dev
touch ~/.earthly-dev/config.yml
cp ~/.earthly-dev/config.yml ~/.earthly-dev/config.yml.bkup

function finish {
  mv ~/.earthly-dev/config.yml.bkup ~/.earthly-dev/config.yml
}
trap finish EXIT

echo "=== Test 1: TLS Enabled ==="
# FIXME bootstrap is failing with "open /home/runner/.earthly-dev/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap || (echo "ignoring bootstrap failure")

# bootstrapping should generate six pem files
test $(ls ~/.earthly-dev/certs/*.pem | wc -l) = "6"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'

rm -rf ~/.earthly-dev/certs

# force buildkit restart before next test
"$earthly" bootstrap || (echo "ignoring bootstrap failure")

echo "=== Test 2: TLS Enabled with different hostname ==="
# FIXME bootstrap is failing with "open /home/runner/.earthly-dev/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap --certs-hostname "$host" || (echo "ignoring bootstrap failure")

# bootstrapping should generate six pem files
test $(ls ~/.earthly-dev/certs/*.pem | wc -l) = "6"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'

rm -rf ~/.earthly-dev/certs

# force buildkit restart before next test
"$earthly" bootstrap || (echo "ignoring bootstrap failure")

echo "=== Test 3: TLS Disabled ==="
"$earthly" config global.tls_enabled false
# FIXME bootstrap is failing with "open /home/runner/.earthly-dev/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap || (echo "ignoring bootstrap failure")

# bootstrapping should not generate any pem files
test $(ls ~/.earthly-dev/certs/*.pem | wc -l) = "0"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'
