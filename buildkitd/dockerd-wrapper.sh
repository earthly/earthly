#!/bin/sh

set -e

if [ -z "$EARTHLY_DOCKERD_DATA_ROOT" ]; then
    echo "EARTHLY_DOCKERD_DATA_ROOT not set"
    exit 1
fi

start_dockerd() {
    mkdir -p "$EARTHLY_DOCKERD_DATA_ROOT"
    dockerd --data-root="$EARTHLY_DOCKERD_DATA_ROOT" >/var/log/docker.log 2>&1 &
    i=1
    timeout=30
    while ! docker ps >/dev/null 2>&1; do
        sleep 1
        if [ "$i" -gt "$timeout" ]; then
            # Print dockerd logs on start failure.
            cat /var/log/docker.log
            exit 1
        fi
        i=$((i+1))
    done
}

stop_dockerd() {
    dockerd_pid="$(cat /var/run/docker.pid)"
    timeout=10
    if [ -n "$dockerd_pid" ]; then
        kill "$dockerd_pid" >/dev/null 2>&1
        i=1
        while kill -0 "$dockerd_pid" >/dev/null 2>&1; do
            sleep 1
            if [ "$i" -gt "$timeout" ]; then
                kill -9 "$dockerd_pid" >/dev/null 2>&1 || true
            fi
            i=$((i+1))
        done
    fi
    # Wipe dockerd data when done.
    rm -rf "$EARTHLY_DOCKERD_DATA_ROOT"
}

load_images() {
    if [ -n "$EARTHLY_DOCKER_LOAD_IMAGES" ]; then
        echo "Loading images..."
        for img in $EARTHLY_DOCKER_LOAD_IMAGES; do
            docker load -i "$img" || (stop_dockerd; exit 1)
        done
        echo "...done"
    fi
}

export EARTHLY_WITH_DOCKER=1

# Lock the creation and destruction of the docker daemon - only one daemon can be started at a time
# (dockerd race conditions in handling networking setup).
# shellcheck disable=SC2039
(
    flock -x 200
    start_dockerd
    flock -u 200
) 200>/var/earthly/dind/lock

load_images

set +e
"$@"
exit_code="$?"
set -e

# shellcheck disable=SC2039
(
    flock -x 200
    stop_dockerd
    flock -u 200
) 200>/var/earthly/dind/lock

exit "$exit_code"
