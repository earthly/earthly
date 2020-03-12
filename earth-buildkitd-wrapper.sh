#!/bin/sh

# Start buildkitd.
/usr/bin/entrypoint.sh \
    buildkitd --allow-insecure-entitlement=security.insecure \
    &>/var/log/buildkitd.log \
    &
buildkitd_pid="$!"
export EARTHLY_BUILDKIT_HOST="unix:///run/buildkit/buildkitd.sock"

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
