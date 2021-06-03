#!/bin/sh

set -e

# Apply global configuration
if [ -n "$GLOBAL_CONFIG" ]; then
  earthly config global "$GLOBAL_CONFIG"
fi

# Apply git configuration
if [ -n "$GIT_CONFIG" ]; then
  earthly config git "$GIT_CONFIG"
fi

# We prefer you use EARTHLY_TOKEN directly, but this is here just in case
if [ -n "$EARTHLY_EMAIL" ] && [ -n "$EARTHLY_PASSWORD" ]; then
  earthly account login --email "$EARTHLY_EMAIL" --password "$EARTHLY_PASSWORD"
fi

# Right now, this container is BYOD - Bring Your Own Docker.
if [ -z "$DOCKER_HOST" ]; then
  echo "DOCKER_HOST is not set"
  sleep 1
fi

# Light check if docker is functional
if ! docker images > /dev/null 2>&1; then
  echo "Docker appears not to be connected. Please check your DOCKER_HOST variable, and try again."
  exit 1
fi

# If no host specified, start an internal buildkit. If it is specified, rely on external setup
buildkit_pid=
if [ -z "$BUILDKIT_HOST" ]; then
  export BUILDKIT_TCP_TRANSPORT_ENABLED=true

  /usr/bin/entrypoint.sh \
    buildkitd \
    --config=/etc/buildkitd.toml \
    &
  buildkitd_pid="$!"

  EARTHLY_BUILDKIT_HOST="tcp://$(hostname):8372" # hostname is not recognized as local for this reason
  export EARTHLY_BUILDKIT_HOST
else
  export EARTHLY_BUILDKIT_HOST="$BUILDKIT_HOST"
fi

echo "Using $EARTHLY_BUILDKIT_HOST as buildkit daemon"

# Use the desired target dir for running a target, saves typing if you use the convention
BASE_DIR="/src"
if [ -n "$SRC_DIR" ]; then
  BASE_DIR="$SRC_DIR"
fi

cd "$BASE_DIR"

# Run earthly with given args.
set +e
earthly "$@"
exit_code="$?"
set -e

# Shut down buildkitd for a graceful exit.
if [ -n "$buildkit_pid" ]; then
  kill "$buildkitd_pid" >/dev/null 2>&1 || true
  i=1
  timeout=10
  while kill -0 "$buildkitd_pid" >/dev/null 2>&1 ; do
      sleep 1
      i=$((i+1))
      if [ "$i" -gt "$timeout" ]; then
          kill -9 "$buildkitd_pid" >/dev/null 2>&1 || true
      fi
  done
fi

exit "$exit_code"