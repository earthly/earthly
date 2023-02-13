#!/usr/bin/env bash
set -e
./lock.sh

function finish {
  status="$?"
  ./unlock.sh
  if [ "$status" = "0" ]; then
    echo "test-gcp passed"
  else
    echo "test-gcp failed with $status"
  fi
}
trap finish EXIT

# Test Google Artifact registry
export ARTIFACT_SERVER="us-west1-docker.pkg.dev"
export ARTIFACT_FULL_ADDRESS="$ARTIFACT_SERVER/ci-cd-302220"
export IMAGE="integration-test/test"
./test-gcp-user.sh
./test-gcp-project.sh

# Test Google container registry
export ARTIFACT_SERVER="gcr.io"
export ARTIFACT_FULL_ADDRESS="$ARTIFACT_SERVER/ci-cd-302220"
export IMAGE="test"
./test-gcp-user.sh
./test-gcp-project.sh
