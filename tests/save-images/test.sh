#!/usr/bin/env bash
set -uex
set -o pipefail

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

cd "$(dirname "$0")"

earthly=${earthly-"../../build/linux/amd64/earthly"}
"$earthly" --version
# docker / podman
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)

# NOTE: the "old behaviour" tests were removed when earthly v0.8 was released
# which obsoleted VERSION 0.5 functionality. TODO: rename "new behaviour to current behaviour"

echo "=== ($LINENO): Test New Behaviour ==="

"$frontend" rmi -f myimage:test
"$frontend" rmi -f mysubimage:test
"$frontend" rmi -f earthly-test-saveimage:test
"$earthly" ./new-behaviour+myimage

"$frontend" inspect myimage:test >/dev/null

if "$frontend" inspect mysubimage:latest; then
	echo "ERROR ($LINENO): mysubimage should not have been saved1."
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via from ==="

"$frontend" rmi -f myimage:fromtest
"$frontend" rmi -f earthly-test-saveimage:test

"$earthly" -V ./new-behaviour-remote+myimage-fromtest
"$frontend" inspect myimage:fromtest >/dev/null
if "$frontend" inspect earthly-test-saveimage:test >/dev/null; then
    echo "ERROR ($LINENO): earthly-test-saveimage:test should not have been saved."   # THIS IS LEGIT BROKEN
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via copy ==="

"$frontend" rmi -f myimage:copytest
"$frontend" rmi -f earthly-test-saveimage:test

"$earthly" ./new-behaviour-remote+myimage-copytest
"$frontend" inspect myimage:copytest >/dev/null
if "$frontend" inspect earthly-test-saveimage:test >/dev/null; then
	echo "ERROR ($LINENO) : earthly-test-saveimage:test should not have been save3."
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via build ==="

"$frontend" rmi -f myimage:buildtest
"$frontend" rmi -f earthly-test-saveimage:test

"$earthly" ./new-behaviour-remote+myimage-buildtest
"$frontend" inspect myimage:buildtest >/dev/null
"$frontend" inspect earthly-test-saveimage:test >/dev/null

echo "=== ($LINENO): Test Disable Earthly Labels ==="

"$frontend" rmi -f myimage:test
"$frontend" rmi -f myimagewithlabels:test

"$earthly" ./disable-earthly-labels+myimage
"$frontend" inspect myimage:test >/dev/null
labels=$("$frontend" inspect myimage:test | jq -r '.[].Config.Labels')

if [ "${labels}" != "null" ] ; then
    echo "ERROR ($LINENO): myimage: should not have labels."
    exit 1
fi

"$earthly" ./disable-earthly-labels+myimagewithlabels
"$frontend" inspect myimagewithlabels:test >/dev/null
labels=$("$frontend" inspect myimagewithlabels:test | jq -r '.[].Config.Labels')

if [ "${labels}" == "null" ] ; then
    echo "ERROR ($LINENO): myimagewithlabels: should *have* labels."
    exit 1
fi

echo "save-images test passed"
