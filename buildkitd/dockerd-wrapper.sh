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

            if grep -q "^failed to start containerd: timeout waiting for containerd to start$" /var/log/docker.log; then
                # This error is the sentinel string for retrying to start dockerd.
                echo "Attempting to restart dockerd (attempt $i), since the error may be transient..."
                sleep 5
            else
                # If the logs do not contain this, then fail fast to maintain prior behavior.
                exit 1
            fi
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
    if ! [ -f /etc/docker/daemon.json ]; then
        mkdir -p /etc/docker
        echo >/etc/docker/daemon.json '{}'
    fi

    daemon_data="$(cat /etc/docker/daemon.json)"
    cat <<EOF | jq ". + $daemon_data" > /etc/docker/daemon.json
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
    ],
    "bip": "172.20.0.1/16",
    "data-root": "$EARTHLY_DOCKERD_DATA_ROOT"
}
EOF

    # Start with a rm -rf to make sure a previous interrupted build did not leave its state around.
    rm -rf "$EARTHLY_DOCKERD_DATA_ROOT"
    mkdir -p "$EARTHLY_DOCKERD_DATA_ROOT"
    dockerd >/var/log/docker.log 2>&1 &
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
            print_dockerd_logs
            echo "If you are having trouble running docker, try using the official earthly/dind image instead"
            return 1
        fi
        i=$((i+1))
    done
}

print_dockerd_logs() {
  echo "Architecture: $(uname -m)"
  echo "==== Begin dockerd logs ===="
  cat /var/log/docker.log
  echo "==== End dockerd logs ===="
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
                echo "dockerd did not exit after $timeout seconds, force-exiting"
                kill -9 "$dockerd_pid" >/dev/null 2>&1 || true
            fi
            i=$((i+1))
        done

        # Wait for the PID to exit. This ensures that dockerd cannot keep any files in EARTHLY_DOCKERD_DATA_ROOT open.
        wait "$dockerd_pid" || true
    fi

      # Wipe dockerd data when done.
    if ! rm -rf "$EARTHLY_DOCKERD_DATA_ROOT"; then
        # We have some issues about failing to delete files. If we fail, list the processes keeping it open for results.
        echo "==== Begin file info ===="
        lsof +D "$EARTHLY_DOCKERD_DATA_ROOT"
        echo "==== End file info logs ===="
        echo "" # Add space between above and docker logs
        print_dockerd_logs
    fi
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
