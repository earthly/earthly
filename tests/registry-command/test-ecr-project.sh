#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

echo "here we go in test-ecr-user.sh"
test -n "$earthly_config" # set by earthly-entrypoint.sh
which earthly

# clear out secrets from previous test
clearusersecrets

# test dockerhub credentials do not exist
earthly registry list | grep -v 404851345508.dkr.ecr.us-west-2.amazonaws.com

# set dockerhub credentials

# TODO implement registry login command for ecr, then switch this test over

ORG="ryan-test"
PROJECT="registry-command-test-project"

echo "setting up cred helper manually"
earthly secrets --org "$ORG" --project "$PROJECT" set std/registry/404851345508.dkr.ecr.us-west-2.amazonaws.com/cred_helper ecr-login
set +x # don't remove, or keys will be leaked
test -n "$AWS_ACCESS_KEY_ID" || (echo "AWS_ACCESS_KEY_ID is empty" && exit 1)
test -n "$AWS_SECRET_ACCESS_KEY" || (echo "AWS_SECRET_ACCESS_KEY is empty" && exit 1)
echo $AWS_ACCESS_KEY_ID | earthly secrets --org "$ORG" --project "$PROJECT" set --stdin std/registry/404851345508.dkr.ecr.us-west-2.amazonaws.com/AWS_ACCESS_KEY_ID
echo $AWS_SECRET_ACCESS_KEY | earthly secrets --org "$ORG" --project "$PROJECT" set --stdin std/registry/404851345508.dkr.ecr.us-west-2.amazonaws.com/AWS_SECRET_ACCESS_KEY
set -x
echo "done setting up cred helper (and secrets)"

# test dockerhub credentials exist
earthly registry list # TODO validate this works

uuid="$(uuidgen)"

cat > Earthfile <<EOF
VERSION 0.7
PROJECT ryan-test/registry-command-test-project
pull:
  FROM 404851345508.dkr.ecr.us-west-2.amazonaws.com/integration-test:latest
  RUN test -f /etc/passwd

push:
  FROM alpine
  RUN echo $uuid > /some-data
  SAVE IMAGE --push 404851345508.dkr.ecr.us-west-2.amazonaws.com/integration-test:latest
EOF

# --no-output is required for earthly-in-earthly; however a --push to ecr will still occur
earthly --config "$earthly_config" --verbose +pull
earthly --config "$earthly_config" --no-output --push --verbose +push

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
