#!/bin/sh

set -eu

# Start buildkitd.
rm -f "/run/buildkit/buildkitd.sock"
/usr/bin/entrypoint.sh \
    buildkitd \
    --allow-insecure-entitlement=security.insecure \
    --config=/etc/buildkitd.toml \
    >/var/log/buildkitd.log 2>&1 \
    &
buildkitd_pid="$!"
export EARTHLY_BUILDKIT_HOST="unix:///run/buildkit/buildkitd.sock"

# Poll for buildkitd readiness.
i=1
timeout=120
while [ ! -S "/run/buildkit/buildkitd.sock" ]; do
    sleep 1
    i=$((i+1))
    if [ "$i" -gt "$timeout" ]; then
        kill -9 "$buildkitd_pid" >/dev/null 2>&1 || true
        echo "Buildkitd did not start within $timeout seconds"
        echo "Buildkitd log"
        echo "=============="
        cat /var/log/buildkitd.log
        echo "=============="
        exit 1
    fi
done

# Run earthly with given args.
set +e
earthly "$@"
exit_code="$?"
set -e

# Shut down buildkitd.
kill "$buildkitd_pid" >/dev/null 2>&1 || true
i=1
timeout=10
while kill -0 "$buildkitd_pid" >/dev/null 2>&1 ; do
    sleep 1
    i=$((i+1))
    if [ "$i" -gt "$timeout" ]; then
        kill -9 "$buildkitd_pid" >/dev/null 2>&1 || true
    fi
done

exit "$exit_code"
