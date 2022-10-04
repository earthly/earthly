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

echo "=== ($LINENO): Test Old Copy Behaviour ==="

"$frontend" rmi -f myimage:copytest
"$frontend" rmi -f mysubimage:copytest
"$earthly" ./old-behaviour-copy+myimage
echo done
"$frontend" inspect myimage:copytest >/dev/null
"$frontend" inspect mysubimage:copytest >/dev/null

echo "=== ($LINENO): Test Old build Behaviour ==="
"$frontend" rmi -f myimage:buildtest
"$frontend" rmi -f mysubimage:buildtest
"$earthly" ./old-behaviour-build+myimage
"$frontend" inspect myimage:buildtest >/dev/null
"$frontend" inspect mysubimage:buildtest >/dev/null

echo "=== ($LINENO): Test Old from Behaviour ==="
"$frontend" rmi -f myimage:fromtest
"$frontend" rmi -f mysubimage:fromtest
"$earthly" ./old-behaviour-from+myimage
"$frontend" inspect myimage:fromtest >/dev/null
"$frontend" inspect mysubimage:fromtest >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by from ==="
"$frontend" rmi -f myimage:old-remote-from-test
"$frontend" rmi -f earthly-test-saveimage:test
"$earthly" ./old-behaviour-remote+myimage-fromtest
"$frontend" inspect myimage:old-remote-from-test >/dev/null
"$frontend" inspect earthly-test-saveimage:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by build ==="
"$frontend" rmi -f myimage:old-remote-build-test
"$frontend" rmi -f earthly-test-saveimage:test
"$earthly" ./old-behaviour-remote+myimage-buildtest
"$frontend" inspect myimage:old-remote-build-test >/dev/null
"$frontend" inspect earthly-test-saveimage:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by copy ==="
"$frontend" rmi -f myimage:old-remote-copy-test
"$frontend" rmi -f earthly-test-saveimage-and-artifact:test
"$earthly" ./old-behaviour-remote+myimage-copytest
"$frontend" inspect myimage:old-remote-copy-test >/dev/null
"$frontend" inspect earthly-test-saveimage-and-artifact:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour on directly referenced remote target ==="
"$frontend" rmi -f earthly-test-saveimage:test
"$earthly" github.com/earthly/test-remote/save-image:main+saveimage
"$frontend" inspect earthly-test-saveimage:test >/dev/null


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

echo "save-images test passed"
