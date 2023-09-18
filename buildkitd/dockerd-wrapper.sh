#!/bin/sh
set -eu

EARTHLY_DOCKER_WRAPPER_DEBUG=${EARTHLY_DOCKER_WRAPPER_DEBUG:-''}
if [ "$EARTHLY_DOCKER_WRAPPER_DEBUG" = "1" ]; then
    echo "enabling docker wrapper debug mode"
    set -x
fi

# This host is used to pull images from the embedded BuildKit Docker registry.
buildkit_docker_registry='172.30.0.1:8371'

detect_docker_compose_cmd() {
    if command -v docker-compose >/dev/null; then
        echo "docker-compose"
        return 0
    fi
    if docker help | grep -w compose >/dev/null; then
        echo "docker compose"
        return 0
    fi
    echo >&2 "failed to detect docker compose / docker-compose command"
    return 1
}

# Runs docker-compose with the right -f flags.
docker_compose_cmd() {
    compose_file_flags=""
    for f in $EARTHLY_COMPOSE_FILES; do
        compose_file_flags="$compose_file_flags -f $f"
    done
    export COMPOSE_HTTP_TIMEOUT=600
    docker_compose="$(detect_docker_compose_cmd)"
    export COMPOSE_PROJECT_NAME="default" # newer versions of docker fail if this is not set; older versions used "default" when it was not set
    # shellcheck disable=SC2086
    $docker_compose $compose_file_flags "$@"
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
    mkdir -p "$EARTHLY_DOCKERD_DATA_ROOT"

    if [ -f "/sys/fs/cgroup/cgroup.controllers" ]; then
        if [ "$EARTHLY_DOCKER_WRAPPER_DEBUG" = "1" ]; then
            echo >&2 "detected cgroups v2"
        fi

        # move script to separate cgroup, to prevent the root cgroup from becoming threaded (which will prevent systemd images (e.g. kind) from running)
        mkdir /sys/fs/cgroup/dockerd-wrapper
        echo $$ > /sys/fs/cgroup/dockerd-wrapper/cgroup.procs

       # earthly wraps dockerd-wrapper.sh with a call via /bin/sh -c '....'
       # so we also need to move the parent pid into this new group, which is weird
       # TODO: we should unwrap this so $$ is all we need to move
        echo 1 > /sys/fs/cgroup/dockerd-wrapper/cgroup.procs

        if [ "$(wc -l < /sys/fs/cgroup/cgroup.procs)" != "0" ]; then
            echo >&2 "warning: processes exist in the root cgroup; this may cause errors during cgroup initialization"
        fi

        root_cgroup_type="$(cat /sys/fs/cgroup/cgroup.type)"
        if [ "$root_cgroup_type" != "domain" ]; then
            echo >&2 "WARNING: expected cgroup type of \"domain\", but got \"$root_cgroup_type\" instead"
        fi
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

    load_file_images
    load_registry_images
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
    data_root=$(TMPDIR="$EARTHLY_DOCKERD_DATA_ROOT/" mktemp -d)
    echo "Starting dockerd with data root $data_root"

    if uname -a | grep microsoft-standard-WSL >/dev/null; then
        if iptables --version | grep nf_tables >/dev/null; then
            echo "WARNING: WSL and iptables-nft may not work; attempting to switch to iptables-legacy"
            ln -sf "/sbin/iptables-legacy" /sbin/iptables
        fi
    fi

    # Use a specific IP range to avoid collision with host dockerd (we need to also connect to host
    # docker containers for the debugger).
    if ! [ -f /etc/docker/daemon.json ]; then
        mkdir -p /etc/docker
        echo >/etc/docker/daemon.json '{}'
    fi

    # compliments of https://stackoverflow.com/a/53666584
    # this will concatenate arrays found in both the LHS and RHS; default jq will overwrite the LHS with the RHS
    cat <<'EOF' > /tmp/meld.jq
def meld(a; b):
  a as $a | b as $b
  | if ($a|type) == "object" and ($b|type) == "object"
    then reduce ([$a,$b]|add|keys_unsorted[]) as $k ({};
      .[$k] = meld( $a[$k]; $b[$k]) )
    elif ($a|type) == "array" and ($b|type) == "array"
    then $a+$b
    elif $b == null then $a
    else $b
    end;
meld($user; .)
EOF

    daemon_data="$(cat /etc/docker/daemon.json)"
    cat <<EOF | jq --argjson user "$daemon_data" -f /tmp/meld.jq > /etc/docker/daemon.json
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
    "data-root": "$data_root",
    "insecure-registries" : ["$buildkit_docker_registry"]
}
EOF

    # Start with wiping the dir to make sure a previous interrupted build did not leave its state around.
    wipe_data_root "$data_root"
    mkdir -p "$data_root"
    rm -f /var/run/docker.pid
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
  cat /var/log/docker.log || true
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

        # Wait for the PID to exit. This ensures that dockerd cannot keep any files in data root open.
        wait "$dockerd_pid" || true
    fi

    # Wipe dockerd data when done.
    wipe_data_root "$data_root"
}

wipe_data_root() {
    if ! rm -rf "$1" 2>/dev/null >&2 && [ -n "$(ls -A "$1")" ]; then
        # We have some issues about failing to delete files.
        # If we fail, list the processes keeping it open for results.
        rm -rf "$1" || true # Do it again, but now print the error.
        echo "==== Begin file lsof info ===="
        if ! lsof +D "$1" ; then
            echo "Failed to run lsof +D $1. Trying lsof $1"
            if ! lsof "$1"; then
                echo "Failed to run lsof $1"
            fi
        fi
        echo "==== End file lsof info ===="
        echo "==== Begin file ls info ===="
        if ! ls -Ral "$1"; then
            echo "Failed to run ls -Ral $1"
        fi
        echo "==== End file ls info ===="
        echo "" # Add space between above and docker logs
        print_dockerd_logs
    fi
}

load_file_images() {
    if [ -n "$EARTHLY_DOCKER_LOAD_FILES" ]; then
        echo "Loading images from BuildKit via tar files..."
        for img in $EARTHLY_DOCKER_LOAD_FILES; do
            docker load -i "$img" || (stop_dockerd; exit 1)
        done
        echo "...done"
    fi
}

load_registry_images() {
    EARTHLY_DOCKER_LOAD_REGISTRY=${EARTHLY_DOCKER_LOAD_REGISTRY:-''}
    if [ -n "$EARTHLY_DOCKER_LOAD_REGISTRY" ]; then
        echo "Loading images from BuildKit via embedded registry..."
        for img in $EARTHLY_DOCKER_LOAD_REGISTRY; do
            case "$img" in
                *'|'*)
                    with_reg="$buildkit_docker_registry/$(printf '%s' "$img" | cut -d'|' -f1)"
                    user_tag="$(printf '%s' "$img" | cut -d'|' -f2-)"
                    ;;
                *)
                    # Old format before v0.6.21.
                    with_reg="$buildkit_docker_registry/$img"
                    user_tag="$(printf '%s' "$img" | cut -d'/' -f2-)"
                    echo "Detected old format"
                    ;;
            esac
            echo "Pulling $with_reg and retagging as $user_tag"
            (docker pull "$with_reg" && docker tag "$with_reg" "$user_tag" && docker rmi "$with_reg") || (stop_dockerd; exit 1)
        done
        echo "...done"
    fi
}

EARTHLY_DOCKER_WRAPPER_DEBUG_CMD=${EARTHLY_DOCKER_WRAPPER_DEBUG_CMD:-''}
if [ -n "$EARTHLY_DOCKER_WRAPPER_DEBUG_CMD" ]; then
    echo "Running debug command: $EARTHLY_DOCKER_WRAPPER_DEBUG_CMD"
    eval "$EARTHLY_DOCKER_WRAPPER_DEBUG_CMD"
    echo "debug command exited with $?; forcing exit 1 to prevent saving RUN snapshot"
    exit 1
fi

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
