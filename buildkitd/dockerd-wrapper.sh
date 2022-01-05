#!/bin/sh

set -eu

# Runs docker-compose with the right -f flags.
docker_compose_cmd() {
    compose_file_flags=""
    for f in $EARTHLY_COMPOSE_FILES; do
        compose_file_flags="$compose_file_flags -f $f"
    done
    export COMPOSE_HTTP_TIMEOUT=600
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

    # Sometimes, when dockerd starts containerd, it doesn't come up in time. This timeout is not configurable from
    # dockerd, therefore we retry... since most instances of this timeout seem to be related to networking or scheduling
    # when many WITH DOCKER commands are also running. Logs are printed for each failure.
    for i in 1 2 3 4 5; do
        if start_dockerd; then
            break
        else
            if [ "$i" = 5 ]; then
                # Exiting here on the last retry maintains prior behavior of exiting when this cant start.
                exit 1
            fi
            sleep 5
        fi
    done

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
    # Use a specific IP range to avoid collision with host dockerd (we need to also connect to host
    # docker containers for the debugger).
    mkdir -p /etc/docker
    cat <<'EOF' >/etc/docker/daemon.json
{
    "default-address-pools" : [
        {
            "base" : "172.21.0.0/16",
            "size" : 24
        },
        {
            "base" : "172.22.0.0/16",
            "size" : 24
        }
    ]
}
EOF

    # Start with a rm -rf to make sure a previous interrupted build did not leave its state around.
    rm -rf "$EARTHLY_DOCKERD_DATA_ROOT"
    mkdir -p "$EARTHLY_DOCKERD_DATA_ROOT"
    dockerd --data-root="$EARTHLY_DOCKERD_DATA_ROOT" --bip=172.20.0.1/16 >/var/log/docker.log 2>&1 &
    dockerd_pid="$!"
    i=1
    timeout=300
    while ! docker ps >/dev/null 2>&1; do
        sleep 1
        fail=false
        if [ "$i" -gt "$timeout" ]; then
            echo "ERROR: dockerd start timeout (${timeout}s)"
            fail=true
        fi
        if ! kill -0 "$dockerd_pid" >/dev/null 2>&1; then
            echo "ERROR: dockerd crashed on startup"
            fail=true
        fi
        if [ "$fail" = "true" ]; then
            # Print dockerd logs on start failure.
            echo "==== Begin dockerd logs ===="
            cat /var/log/docker.log
            echo "==== End dockerd logs ===="
            echo "If you are having trouble running docker, try using the official earthly/dind image instead"
            return 1
        fi
        i=$((i+1))
    done
}

stop_dockerd() {
    dockerd_pid="$(cat /var/run/docker.pid)"
    timeout=30
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
