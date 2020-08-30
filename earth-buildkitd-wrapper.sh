#!/bin/sh

# Start buildkitd.
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
while [ ! -S "/run/buildkit/buildkitd.sock" ]; do
    sleep 1
    i=$((i+1))
    if [ "$i" -gt "30" ]; then
        kill -9 "$buildkitd_pid" >/dev/null 2>&1
        echo "Buildkitd did not start within 10 seconds"
        exit 1
    fi
done

# Run earth with given args.
earth "$@"
exit_code="$?"

# Shut down buildkitd.
kill "$buildkitd_pid" >/dev/null 2>&1
i=1
while kill -0 "$buildkitd_pid" >/dev/null 2>&1 ; do
    sleep 1
    i=$((i+1))
    if [ "$i" -gt "10" ]; then
        kill -9 "$buildkitd_pid" >/dev/null 2>&1
    fi
done
rm -f "/run/buildkit/buildkitd.sock"

exit "$exit_code"
