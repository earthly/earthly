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
test -n "$ECR_REGISTRY_HOST"

# clear out secrets from previous test
clearprojectsecrets

# test credentials do not exist
earthly registry list | grep -v "$ECR_REGISTRY_HOST" # just in case
earthly registry --org "$ORG" --project "$PROJECT" list | grep -v "$ECR_REGISTRY_HOST"

# set credentials
set +x # don't remove, or keys will be leaked
test -n "$AWS_ACCESS_KEY_ID" || (echo "AWS_ACCESS_KEY_ID is empty" && exit 1)
test -n "$AWS_SECRET_ACCESS_KEY" || (echo "AWS_SECRET_ACCESS_KEY is empty" && exit 1)
set -x
earthly registry --org "$ORG" --project "$PROJECT" setup --cred-helper=ecr-login "$ECR_REGISTRY_HOST"

# test credentials exist
earthly registry --org "$ORG" --project "$PROJECT" list | grep "$ECR_REGISTRY_HOST"

uuid="$(uuidgen)"

cat > Earthfile <<EOF
VERSION 0.7
PROJECT ryan-test/registry-command-test-project
pull:
  FROM $ECR_REGISTRY_HOST/integration-test:latest
  RUN test -f /etc/passwd

push:
  FROM alpine
  RUN echo $uuid > /some-data
  SAVE IMAGE --push $ECR_REGISTRY_HOST/integration-test:latest
EOF

# --no-output is required for earthly-in-earthly; however a --push to ecr will still occur
earthly --config "$earthly_config" --verbose +pull
earthly --config "$earthly_config" --no-output --push --verbose +push

earthly registry --org "$ORG" --project "$PROJECT" remove "$ECR_REGISTRY_HOST"
earthly registry --org "$ORG" --project "$PROJECT" list | grep -v $ECR_REGISTRY_HOST

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearprojectsecrets
