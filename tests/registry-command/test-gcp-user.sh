#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

test -n "$earthly_config" # set by earthly-entrypoint.sh
which earthly

# clear out secrets from previous test
clearusersecrets

# test dockerhub credentials do not exist
earthly registry list | grep -v $ARTIFACT_SERVER

# set dockerhub credentials

# TODO implement registry login command for gcloud artifact registry, then switch this test over

echo "setting up cred helper manually"
earthly secrets set /user/std/registry/$ARTIFACT_SERVER/cred_helper gcp-login
set +x # don't remove, or keys will be leaked
test -n "$GCP_KEY" || (echo "GCP_KEY is empty" && exit 1)
echo $GCP_KEY | earthly secrets set --stdin /user/std/registry/$ARTIFACT_SERVER/GCP_KEY
set -x
echo "done setting up cred helper (and secrets)"

# test dockerhub credentials exist
earthly registry list # TODO validate this works

uuid="$(uuidgen)"

cat > Earthfile <<EOF
VERSION 0.7
pull:
  FROM $ARTIFACT_FULL_ADDRESS/$IMAGE:latest
  RUN test -f /etc/passwd

push:
  FROM alpine
  RUN echo $uuid > /some-data
  SAVE IMAGE --push $ARTIFACT_FULL_ADDRESS/$IMAGE:latest
EOF

# --no-output is required for earthly-in-earthly; however a --push to gcp will still occur
earthly --config "$earthly_config" --verbose +pull
earthly --config "$earthly_config" --no-output --push --verbose +push

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
