#!/bin/sh

# Start buildkitd.
/usr/bin/entrypoint.sh \
    buildkitd \
    --allow-insecure-entitlement=security.insecure \
    --config=/etc/buildkitd.toml \
    &>/var/log/buildkitd.log \
    &
buildkitd_pid="$!"
export EARTHLY_BUILDKIT_HOST="unix:///run/buildkit/buildkitd.sock"

# Poll for buildkitd readiness.
let i=1
while [ ! -S "/run/buildkit/buildkitd.sock" ]; do
    sleep 1
    let i+=1
    if [ "$i" -gt "10" ]; then
        kill -9 "$buildkitd_pid" &>/dev/null
        echo "Buildkitd did not start within 10 seconds"
        exit 1
    fi
done

# Run earth with given args.
earth "$@"
exit_code="$?"

# Shut down buildkitd.
kill "$buildkitd_pid" &>/dev/null
let i=1
while kill -0 "$buildkitd_pid" &>/dev/null ; do
    sleep 1
    let i+=1
    if [ "$i" -gt "10" ]; then
        kill -9 "$buildkitd_pid" &>/dev/null
    fi
done
rm -f "/run/buildkit/buildkitd.sock"

exit "$exit_code"
