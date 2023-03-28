#!/bin/sh

set -e

EARTHLY_DEBUG=${EARTHLY_DEBUG:-false}
if [ "$EARTHLY_DEBUG" = "true" ]; then
    set -x
    export EARTHLY_DEBUG
fi

earthly_config="/etc/.earthly/config.yml"
if [ ! -f "$earthly_config" ]; then
  # Missing config, generate it and use the env vars
  # Do not do both, since that would write to the mounted config
  mkdir -p "$(dirname $earthly_config)" && touch "$earthly_config"

  # Apply global configuration
  if [ -n "$GLOBAL_CONFIG" ]; then
    earthly --config "$earthly_config" config global "$GLOBAL_CONFIG"
  fi

  # Apply git configuration
  if [ -n "$GIT_CONFIG" ]; then
    earthly --config $earthly_config config git "$GIT_CONFIG"
  fi
fi

# If no host specified, start an internal buildkit. If it is specified, rely on external setup
if [ -z "$NO_BUILDKIT" ]; then
  if [ -z "$BUILDKIT_HOST" ]; then
    if ! captest --text | grep sys_admin > /dev/null; then
      echo 1>&2 "Container appears to be running unprivileged. Currently, privileged mode is required when buildkit runs inside the container."
      echo 1>&2 "To run this image without buildkit, set the environment variable NO_BUILDKIT=1"
      exit 1
    fi

    if [ -f "/sys/fs/cgroup/cgroup.controllers" ]; then
        ! "$EARTHLY_DEBUG" || echo 1>&2 "detected cgroups v2; earthly-entrypoint.sh pid=$$"

        # move the process under a new cgroup to prevent buildkitd/entrypoint.sh
        # from getting a "h: write error: Resource busy" error while enabling controllers
        # via echo +pids > /sys/fs/cgroup/cgroup.subtree_control
        ( \
          mkdir -p /sys/fs/cgroup/earthly-entrypoint && \
          echo "$$" > /sys/fs/cgroup/earthly-entrypoint/cgroup.procs \
        ) || true
    fi

    export BUILDKIT_TCP_TRANSPORT_ENABLED=true

    /usr/bin/entrypoint.sh \
      buildkitd \
        --config=/etc/buildkitd.toml \
        >/var/log/buildkitd.log 2>&1 \
        &

    if [ "$BUILDKIT_DEBUG" = "true" ]; then
        tail -f /var/log/buildkitd.log &
    fi

    EARTHLY_BUILDKIT_HOST="tcp://$(hostname):8372" # hostname is not recognized as local for this reason
    export EARTHLY_BUILDKIT_HOST
    earthly --config "$earthly_config" config global.tls_enabled false
  else
    export EARTHLY_BUILDKIT_HOST="$BUILDKIT_HOST"
  fi
  ! "$EARTHLY_DEBUG" || echo 1>&2 "Using $EARTHLY_BUILDKIT_HOST as buildkit daemon"
fi

if [ -n "$SRC_DIR" ]; then
  echo 1>&2 'Please note that SRC_DIR is deprecated. This script will no longer automatically switch to it in the future.'
  echo 1>&2 'Please change the container'"'"'s working directory instead (e.g. via docker run -w)'
  cd "$SRC_DIR"
fi

if [ -n "$EARTHLY_EXEC_CMD" ]; then
    export earthly_config
    exec "$EARTHLY_EXEC_CMD"
    exit 1 # this should never be reached
fi

# Run earthly with given args.
# Exec so we don't have to trap and manage signal propagation
exec earthly --config "$earthly_config" "$@"
