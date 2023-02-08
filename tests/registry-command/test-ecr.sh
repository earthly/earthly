#!/usr/bin/env bash
set -e
./lock.sh

function finish {
  status="$?"
  ./unlock.sh
  if [ "$status" = "0" ]; then
    echo "test-ecr passed"
  else
    echo "test-ecr failed with $status"
  fi
}
trap finish EXIT

./test-ecr-user.sh
./test-ecr-project.sh
