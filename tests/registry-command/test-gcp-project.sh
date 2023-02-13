#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

ARTIFACT_SERVER="us-west1-docker.pkg.dev"
ARTIFACT_FULL_ADDRESS="$ARTIFACT_SERVER/ci-cd-302220"
IMAGE="integration-test/test"

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

echo "here we go in test-gcp-project.sh"
test -n "$earthly_config" # set by earthly-entrypoint.sh
which earthly

# clear out secrets from previous test
clearusersecrets

# test dockerhub credentials do not exist
earthly registry list | grep -v $ARTIFACT_SERVER

# set dockerhub credentials

# TODO implement registry login command for gcr, then switch this test over

ORG="ryan-test"
PROJECT="registry-command-test-project"

echo "setting up cred helper manually"
earthly secrets --org "$ORG" --project "$PROJECT" set std/registry/$ARTIFACT_SERVER/cred_helper gcp-login
set +x # don't remove, or keys will be leaked
test -n "$GCP_KEY" || (echo "GCP_KEY is empty" && exit 1)
echo $GCP_KEY | earthly secrets --org "$ORG" --project "$PROJECT" set --stdin std/registry/$ARTIFACT_SERVER/GCP_KEY
set -x
echo "done setting up cred helper (and secrets)"

# test dockerhub credentials exist
earthly registry list # TODO validate this works

uuid="$(uuidgen)"

cat > Earthfile <<EOF
VERSION 0.7
PROJECT ryan-test/registry-command-test-project
pull:
  FROM $ARTIFACT_FULL_ADDRESS/$IMAGE:latest
  RUN test -f /etc/passwd

push:
  FROM alpine
  RUN echo $uuid > /some-data
  SAVE IMAGE --push $ARTIFACT_FULL_ADDRESS/$IMAGE:latest
EOF

# --no-output is required for earthly-in-earthly; however a --push to gcloud artifact registry will still occur
earthly --config "$earthly_config" --verbose +pull
earthly --config "$earthly_config" --no-output --push --verbose +push

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
