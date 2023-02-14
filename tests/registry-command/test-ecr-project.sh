#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

test -n "$earthly_config" # set by earthly-entrypoint.sh
test -n "$ECR_REGISTRY_HOST"

# clear out secrets from previous test
clearusersecrets

# test dockerhub credentials do not exist
earthly registry list | grep -v "$ECR_REGISTRY_HOST"

# set dockerhub credentials

# TODO implement registry login command for ecr, then switch this test over

ORG="ryan-test"
PROJECT="registry-command-test-project"

echo "setting up cred helper manually"
earthly secrets --org "$ORG" --project "$PROJECT" set "std/registry/$ECR_REGISTRY_HOST/cred_helper" ecr-login
set +x # don't remove, or keys will be leaked
test -n "$AWS_ACCESS_KEY_ID" || (echo "AWS_ACCESS_KEY_ID is empty" && exit 1)
test -n "$AWS_SECRET_ACCESS_KEY" || (echo "AWS_SECRET_ACCESS_KEY is empty" && exit 1)
echo $AWS_ACCESS_KEY_ID | earthly secrets --org "$ORG" --project "$PROJECT" set --stdin "std/registry/$ECR_REGISTRY_HOST/AWS_ACCESS_KEY_ID"
echo $AWS_SECRET_ACCESS_KEY | earthly secrets --org "$ORG" --project "$PROJECT" set --stdin "std/registry/$ECR_REGISTRY_HOST/AWS_SECRET_ACCESS_KEY"
set -x
echo "done setting up cred helper (and secrets)"

# test dockerhub credentials exist
earthly registry list # TODO validate this works

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

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
