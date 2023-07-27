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
        echo >&2 "detected cgroups v2; earthly-entrypoint.sh running under pid=$$ with controllers \"$(cat /sys/fs/cgroup/cgroup.controllers)\" in group $(cat /proc/self/cgroup)"
        test "$(cat /sys/fs/cgroup/cgroup.type)" = "domain" || (echo >&2 "WARNING: invalid root cgroup type: $(cat /sys/fs/cgroup/cgroup.type)")
    fi

    # generate certificates
    earthly --config "$earthly_config" --buildkit-host=tcp://127.0.0.1:8372 bootstrap --certs-hostname="$(hostname)"

    if [ ! -f /etc/ca.pem ]; then
      ln -s /root/.earthly/certs/ca_cert.pem /etc/ca.pem
    fi

    if [ ! -f /etc/cert.pem ]; then
      ln -s /root/.earthly/certs/buildkit_cert.pem /etc/cert.pem
    fi

    if [ ! -f /etc/key.pem ]; then
      ln -s /root/.earthly/certs/buildkit_key.pem /etc/key.pem
    fi


    export BUILDKIT_TCP_TRANSPORT_ENABLED=true
    export BUILDKIT_TLS_ENABLED=true

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
