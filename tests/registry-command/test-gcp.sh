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

./test-gcp-user.sh
./test-gcp-project.sh
