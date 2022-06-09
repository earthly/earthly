#!/bin/bash
# This test is designed to be run directly by github actions or on your host (i.e. not earthly-in-earthly);
# as a result, you may run into issues if you have a firewall which prevents access to the registry -- make sure access to it's port is open
set -uxe
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../build/linux/amd64/earthly"}
echo "using earthly=$(realpath "$earthly")"

registry_name="wait-block-registry"

# Cleanup previous run.
docker stop "$registry_name" || true
docker rm "$registry_name" || true
rm -rf ./certs || true

# Get host IP (it must be accessible from both docker and runc)
HOST_IP=$(ip route get 8.8.8.8 | awk -F"src " 'NR==1{split($2,a," ");print a[1]}')

# pick the first free port starting at 5000 (up to 5050)
PORT_FOUND="false"
for i in {0..50}
do
    PORT="$(( 5000 + i ))"
    ACTIVE_PORT="$(netstat -lnt | awk '{print $4}' | (grep ":$PORT\$" || true) | wc -l)"
    if [ "$ACTIVE_PORT" = "0" ]; then
        PORT_FOUND="true"
        break
    fi
done
test "$PORT_FOUND" = "true"

REGISTRY="$HOST_IP:$PORT"
echo "running registry on $REGISTRY"

# Generate certs.
"$earthly" \
    --build-arg REGISTRY_IP="$HOST_IP" \
    --artifact ../registry-certs+certs/certs "$(pwd)/certs"

CRT_PATH="$(pwd)/certs/domain.crt"

# A random tmp file which shouldn't conflict with anything else
config_path="/tmp/earthly-34a5d7b5-903e-40d8-ade3-260ff9794f93.yml"

cat > "$config_path" <<EOF
global:
  buildkit_additional_args: ["-v", "$CRT_PATH:/etc/config/wait-block-test.ca"]
  buildkit_additional_config: |
    [registry."$REGISTRY"]
      ca=["/etc/config/wait-block-test.ca"]
EOF

# Run registry. This will use the same IP address as allocated above.
docker run -d \
    -p "$PORT:$PORT" \
    -v "$(pwd)"/certs:/certs \
    -e REGISTRY_HTTP_ADDR="0.0.0.0:$PORT" \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    --name "$registry_name" registry:2

# First make sure all deps get cached, to increase the likelihood of a race-condition
"$earthly" --config="$config_path" +deps

# Test.
tag="$(uuidgen)"
test -n "$tag"
echo "using tag=$tag"
pwd
set +e
"$earthly" --config="$config_path" --push $@ +test --tag="$tag" --REGISTRY="$REGISTRY"
exit_code="$?"
set -e

# Cleanup.
docker stop "$registry_name" || true

exit "$exit_code"
