#!/usr/bin/env bash
set -ex

ORG="ryan-test"
PROJECT="registry-command-test-project"

oldlockvalue=""
while true; do
    lock="$(earthly secrets --org "$ORG" --project "$PROJECT" get lock || true)"
    if [ -z "$lock" ]; then
        echo "no lock exists; proceeding to lock it"
        break
    fi
    if [ "$lock" = "$oldlockvalue" ]; then
        echo "lock value hasn't changed; forcing it open"
        earthly secrets --org "$ORG" --project "$PROJECT" rm lock || true
        sleep $[ ( $RANDOM % 5 ) + 1 ]s
        continue
    fi
    oldlockvalue="$lock"
    # TODO implement a secrets ls --long, which would show a "date created/modified" column
    # then if the lock is older than 1 minute, we would consider it abandoned, delete it, and create
    # a new lock. For now, we will simply sleep for 60 seconds (which should be enough time for the test to pass)
    duration=$[ ( $RANDOM % 30 ) + 180 ]
    echo "lock exists; sleeping for $duration seconds"
    sleep "$duration"
done

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
