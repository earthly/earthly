#!/usr/bin/env bash
set -eu # don't use -x as it will leak the JWT

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
echo "running tests with $earthly"

## Test that we display web link correctly without being logged in and with no params
NO_COLOR=0 "$earthly" web 2>&1 | grep -qE "https://ci(-beta)?.earthly.dev/login"
RESULT=$?
if [ $RESULT -ne 0 ]; then
  echo "failed without login and no arguments"
  exit $RESULT
fi

# Test with arguments passed for provider and final
NO_COLOR=0 "$earthly" web --provider github --final ci 2>output-file
output=$(cat output-file)
rm output-file

RESULT=$?
if [ $RESULT -ne 0 ]; then
  echo "failed without login provider / final arguments passed"
  exit $RESULT
fi

if [[ "$output" != *"final=ci"* ]]; then
  echo "failed to find final=ci in output"
  exit 1
fi

if [[ "$output" != *"provider=github"* ]]; then
  echo "failed to find provider=github in output"
  exit 1
fi

# Login for the authed tests
test -n "$EARTHLY_TOKEN"

echo Logging in
"$earthly" account login

# Test logged in - no args
NO_COLOR=0 "$earthly" web 2>&1 | grep -qE "https://ci(-beta)?.earthly.dev/login\?token=(.*)+"
RESULT=$?
if [ $RESULT -ne 0 ]; then
  echo "failed with login and no arguments"
  exit $RESULT
fi


# Test logged in with arguments passed for provider and final
NO_COLOR=0 "$earthly" web --provider github --final ci 2>output-file
output=$(cat output-file)
rm output-file

RESULT=$?
if [ $RESULT -ne 0 ]; then
  echo "failed without login provider / final arguments passed"
  exit $RESULT
fi

if [[ "$output" != *"final=ci"* ]]; then
  echo "failed to find final=ci in output"
  exit 1
fi

if [[ "$output" != *"provider=github"* ]]; then
  echo "failed to find provider=github in output"
  exit 1
fi

if [[ "$output" != *"token="* ]]; then
  echo "failed to find provider= in output"
  exit 1
fi
