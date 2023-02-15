#!/bin/sh
set -ex
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

export ECR_REGISTRY_HOST="404851345508.dkr.ecr.us-west-2.amazonaws.com"

./test-ecr-user.sh
./test-ecr-project.sh
