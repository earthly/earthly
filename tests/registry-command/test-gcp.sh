#!/bin/sh
set -e
./lock.sh

finish() {
  status="$?"
  ./unlock.sh
  if [ "$status" = "0" ]; then
    echo "test-ecr passed"
  else
    echo "test-ecr failed with $status"
  fi
}
trap finish EXIT

# Test Google artifact registry
export GCP_SERVER="us-west1-docker.pkg.dev"
export GCP_FULL_ADDRESS="$GCP_SERVER/ci-cd-302220"
export IMAGE="integration-test/test"
./test-gcp-user.sh
./test-gcp-project.sh

# Test Google container registry
export GCP_SERVER="gcr.io"
export GCP_FULL_ADDRESS="$GCP_SERVER/ci-cd-302220"
export IMAGE="test"
./test-gcp-user.sh
./test-gcp-project.sh
