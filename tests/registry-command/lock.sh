#!/bin/sh
set -ex

ORG="ryan-test"
PROJECT="registry-command-test-project"

lock="$(earthly secrets --org "$ORG" --project "$PROJECT" get lock || true)"
if [ -n "$lock" ]; then
    # TODO implement a secrets ls --long, which would show a "date created/modified" column
    # then if the lock is older than 1 minute, we would consider it abandoned, delete it, and create
    # a new lock. For now, we will simply sleep for 30 seconds (which should be enough time for the test to pass)
    echo "lock exists; sleeping for 30 seconds"
    sleep 30
fi

echo "no lock exists; proceeding to lock it"

id="$(uuidgen)"
test -n "$id"

earthly secrets --org "$ORG" --project "$PROJECT" set lock "$id"

sleep 1

lock="$(earthly secrets --org "$ORG" --project "$PROJECT" get lock)"
if [ "$lock" != "$id" ]; then
  echo "failed to lock"
  exec ./lock.sh # try again
fi

echo "$id" > /tmp/registry-command-lock
echo locked
