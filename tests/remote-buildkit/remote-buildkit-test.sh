#!/bin/bash
set -uxe
set -o pipefail

cd "$(dirname "$0")"
earthly=${earthly-"../../build/linux/amd64/earthly"}
host=$(hostname)

# so that we can test production/staging binaries
default_install_name="${default_install_name:-"earthly-dev"}"

mkdir -p ~/.${default_install_name}
touch ~/.${default_install_name}/config.yml
cp ~/.${default_install_name}/config.yml ~/.${default_install_name}/config.yml.bkup

function finish {
  mv ~/.${default_install_name}/config.yml.bkup ~/.${default_install_name}/config.yml
}
trap finish EXIT

echo "=== Test 1: TLS Enabled ==="
# FIXME bootstrap is failing with "open /home/runner/.${default_install_name}/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap || (echo "ignoring bootstrap failure")

# bootstrapping should generate six pem files
test $(ls ~/.${default_install_name}/certs/*.pem | wc -l) = "6"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'

rm -rf ~/.${default_install_name}/certs

# force buildkit restart before next test
"$earthly" bootstrap || (echo "ignoring bootstrap failure")

echo "=== Test 2: TLS Enabled with different hostname ==="
# FIXME bootstrap is failing with "open /home/runner/.${default_install_name}/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap --certs-hostname "$host" || (echo "ignoring bootstrap failure")

# bootstrapping should generate six pem files
test $(ls ~/.${default_install_name}/certs/*.pem | wc -l) = "6"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'

rm -rf ~/.${default_install_name}/certs

# force buildkit restart before next test
"$earthly" bootstrap || (echo "ignoring bootstrap failure")

echo "=== Test 3: TLS Disabled ==="
"$earthly" config global.tls_enabled false
# FIXME bootstrap is failing with "open /home/runner/.${default_install_name}/certs/ca_cert.pem: permission denied", but generates them nonetheless.
"$earthly" --verbose --buildkit-host tcp://127.0.0.1:8372 bootstrap || (echo "ignoring bootstrap failure")

# bootstrapping should not generate any pem files
test $(ls ~/.${default_install_name}/certs/*.pem | wc -l) = "0"

"$earthly" --no-cache --verbose --buildkit-host tcp://127.0.0.1:8372 +target 2>&1 | perl -pe 'BEGIN {$status=1} END {exit $status} $status=0 if /running under remote-buildkit test/;'
