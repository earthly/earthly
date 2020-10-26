#!/bin/sh

set -eu

# Runs docker-compose with the right -f flags.
docker_compose_cmd() {
    compose_file_flags=""
    for f in $EARTHLY_COMPOSE_FILES; do
        compose_file_flags="$compose_file_flags -f $f"
    done
    # shellcheck disable=SC2086
    docker-compose $compose_file_flags "$@"
}

write_compose_config() {
    mkdir -p /tmp/earthly
    docker_compose_cmd config >/tmp/earthly/compose-config.yml
}

execute() {
    if [ -z "$EARTHLY_DOCKERD_DATA_ROOT" ]; then
        echo "EARTHLY_DOCKERD_DATA_ROOT not set"
        exit 1
    fi

    # Lock entire execution of a docker daemon - only one daemon can be used at a time
    # (dockerd race conditions in handling networking setup).
    (
        flock -x 8
        start_dockerd
        # Note that the lock will continue to be held after this subshell finishes,
        # becasue it spawns the dockerd background process. This is intentional.
        # The lock is meant to be held until dockerd exits.
    ) 8>/var/earthly/dind/lock
    load_images
    if [ "$EARTHLY_START_COMPOSE" = "true" ]; then
        # shellcheck disable=SC2086
        docker_compose_cmd up -d $EARTHLY_COMPOSE_SERVICES
    fi

    shift
    export EARTHLY_WITH_DOCKER=1
    set +e
    "$@"
    exit_code="$?"
    set -e

    if [ "$EARTHLY_START_COMPOSE" = "true" ]; then
        docker_compose_cmd down --remove-orphans
    fi
    stop_dockerd
    return "$exit_code"
}

start_dockerd() {
    mkdir -p "$EARTHLY_DOCKERD_DATA_ROOT"
    dockerd --data-root="$EARTHLY_DOCKERD_DATA_ROOT" >/var/log/docker.log 2>&1 &
    i=1
    timeout=30
    while ! docker ps >/dev/null 2>&1; do
        sleep 1
        if [ "$i" -gt "$timeout" ]; then
            # Print dockerd logs on start failure.
            echo "==== Begin dockerd logs ===="
            cat /var/log/docker.log
            echo "==== End dockerd logs ===="
            echo "If you are having trouble running docker, try using the official earthly/dind image instead"
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
    if [ -n "$EARTHLY_DOCKER_LOAD_FILES" ]; then
        echo "Loading images..."
        for img in $EARTHLY_DOCKER_LOAD_FILES; do
            docker load -i "$img" || (stop_dockerd; exit 1)
        done
        echo "...done"
    fi
}

case "$1" in
    get-compose-config)
        write_compose_config
        exit 0
        ;;
    
    execute)
        execute "$@"
        exit "$?"
        ;;
    
    *)
        echo "Invalid command $1"
        exit 1
        ;;
esac
