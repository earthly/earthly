#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

ORG="ryan-test"
PROJECT="registry-command-test-project"

clearprojectsecrets() {
    earthly secrets --org "$ORG" --project "$PROJECT" ls std/ | xargs -r -n 1 earthly secrets --org "$ORG" --project "$PROJECT" rm
}

test -n "$earthly_config" # set by earthly-entrypoint.sh

# clear out secrets from previous test
clearprojectsecrets

# test credentials do not exist
earthly registry list | grep -v "$GCP_SERVER"

# set credentials
set +x # don't remove, or keys will be leaked
test -n "$GCP_KEY" || (echo "GCP_KEY is empty" && exit 1)
set -x
earthly registry setup --org "$ORG" --project "$PROJECT" --cred-helper=gcloud "$GCP_SERVER"

# test credentials exist
earthly registry list --org "$ORG" --project "$PROJECT" | grep "$GCP_SERVER"

uuid="$(uuidgen)"

cat > Earthfile <<EOF
VERSION 0.7
PROJECT $ORG/$PROJECT
pull:
  FROM $GCP_FULL_ADDRESS/$IMAGE:latest
  RUN test -f /etc/passwd

push:
  FROM alpine
  RUN echo $uuid > /some-data
  SAVE IMAGE --push $GCP_FULL_ADDRESS/$IMAGE:latest
EOF

# --no-output is required for earthly-in-earthly; however a --push to gcp will still occur
earthly --config "$earthly_config" --verbose +pull
earthly --config "$earthly_config" --no-output --push --verbose +push

earthly registry remove --org "$ORG" --project "$PROJECT" "$GCP_SERVER"
earthly registry list --org "$ORG" --project "$PROJECT" | grep -v $GCP_SERVER

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearprojectsecrets
