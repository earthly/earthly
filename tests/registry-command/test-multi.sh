#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

./lock.sh

finish() {
  status="$?"
  ./unlock.sh
  if [ "$status" = "0" ]; then
    echo "$0 passed"
  else
    echo "$0 failed with $status"
  fi
}
trap finish EXIT


# ECR details
export ECR_REGISTRY_HOST="404851345508.dkr.ecr.us-west-2.amazonaws.com"

# Google artifact registry details
export GCP_SERVER="us-west1-docker.pkg.dev"
export GCP_FULL_ADDRESS="$GCP_SERVER/ci-cd-302220"
export IMAGE="integration-test/test"

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

# clear out secrets from previous test
clearusersecrets


echo "Setting up ECR credentials"
set +x # don't remove, or keys will be leaked
test -n "$AWS_ACCESS_KEY_ID" || (echo "AWS_ACCESS_KEY_ID is empty" && exit 1)
test -n "$AWS_SECRET_ACCESS_KEY" || (echo "AWS_SECRET_ACCESS_KEY is empty" && exit 1)
set -x
earthly registry setup --cred-helper=ecr-login "$ECR_REGISTRY_HOST"

echo "Setting up GCP credentials"
set +x # don't remove, or keys will be leaked
test -n "$GCP_KEY" || (echo "GCP_KEY is empty" && exit 1)
export GCP_SERVICE_ACCOUNT_KEY="$GCP_KEY" # registry setup reads from this env
set -x
earthly registry setup --cred-helper=gcloud "$GCP_SERVER"


echo "done setting up cred helper (and secrets)"

earthly registry list | grep "$ECR_REGISTRY_HOST"
earthly registry list | grep "$GCP_SERVER"

cat > Earthfile <<EOF
VERSION 0.7
pull-dockerhub:
  FROM earthly/rot13
  RUN which ncat # installed on earthly/rot13

pull-ecr:
  FROM $ECR_REGISTRY_HOST/integration-test:latest
  RUN test -f /etc/passwd

pull-gcp:
  FROM $GCP_FULL_ADDRESS/$IMAGE:latest
  RUN test -f /etc/passwd

pull:
  BUILD +pull-dockerhub
  BUILD +pull-ecr
  BUILD +pull-gcp
EOF

earthly --config "$earthly_config" --verbose +pull

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
