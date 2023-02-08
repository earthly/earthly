#!/bin/sh
set -e

ORG="ryan-test"
PROJECT="registry-command-test-project"

id="$(cat /tmp/registry-command-lock)"

lock="$(earthly secrets --org "$ORG" --project "$PROJECT" get lock || true)"
if [ "$lock" = "$id" ]; then
  earthly secrets --org "$ORG" --project "$PROJECT" rm lock
else
  echo "unlock failed: unexpected lock contents (expected $id; got $lock)"
fi
