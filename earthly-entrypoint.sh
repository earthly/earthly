#!/bin/sh
set -e

earthly_config="/etc/.earthly/config.yml"
if [ ! -f "$earthly_config" ]; then
  # Missing config, generate it and use the env vars
  # Do not do both, since that would write to the mounted config
  mkdir -p "$(dirname $earthly_config)" && touch $earthly_config

  # Apply global configuration
  if [ -n "$GLOBAL_CONFIG" ]; then
    earthly --config "$earthly_config" config global "$GLOBAL_CONFIG"
  fi

  # Apply git configuration
  if [ -n "$GIT_CONFIG" ]; then
    earthly --config $earthly_config config git "$GIT_CONFIG"
  fi
fi

# Skip docker if you are not exporting any images.
if [ -z "$NO_DOCKER" ]; then
  # Right now, this container is BYOD - Bring Your Own Docker.
  if [ -z "$DOCKER_HOST" ]; then
    echo "DOCKER_HOST is not set"
  fi

  # Light check if docker is functional
  if ! docker images > /dev/null 2>&1; then
    echo "Docker appears not to be connected. Please check your DOCKER_HOST variable, and try again."
    exit 1
  fi
fi

# If no host specified, start an internal buildkit. If it is specified, rely on external setup
if [ -z "$BUILDKIT_HOST" ]; then

  # Is container running as privileged? This is currently required when starting up and using buildkit
  if ! captest --text | grep sys_admin > /dev/null; then
    echo "Container appears to be running unprivileged. Currently, privileged mode is required when buildkit runs inside the container."
    exit 1
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
else
  export EARTHLY_BUILDKIT_HOST="$BUILDKIT_HOST"
fi

echo "Using $EARTHLY_BUILDKIT_HOST as buildkit daemon"

if [ -n "$SRC_DIR" ]; then
  echo 'Please note that SRC_DIR is deprecated. This script will no longer automatically switch to it in the future.'
  echo 'Please change the container'"'"'s working directory instead (e.g. via docker run -w)'
  cd "$SRC_DIR"
fi

if [ -n "$EARTHLY_EXEC_CMD" ]; then
    export earthly_config
    exec "$EARTHLY_EXEC_CMD"
    exit 1 # this should never be reached
fi

# Run earthly with given args.
# Exec so we don't have to trap and manage signal propagation
exec earthly --config $earthly_config "$@"
