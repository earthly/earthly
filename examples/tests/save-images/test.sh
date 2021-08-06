#!/usr/bin/env bash
set -uex
set -o pipefail

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
"$earthly" --version

echo "=== ($LINENO): Test Old Copy Behaviour ==="

docker rmi -f myimage:copytest
docker rmi -f mysubimage:copytest
"$earthly" ./old-behaviour-copy+myimage
echo done
docker inspect myimage:copytest >/dev/null
docker inspect mysubimage:copytest >/dev/null

echo "=== ($LINENO): Test Old build Behaviour ==="
docker rmi -f myimage:buildtest
docker rmi -f mysubimage:buildtest
"$earthly" ./old-behaviour-build+myimage
docker inspect myimage:buildtest >/dev/null
docker inspect mysubimage:buildtest >/dev/null

echo "=== ($LINENO): Test Old from Behaviour ==="
docker rmi -f myimage:fromtest
docker rmi -f mysubimage:fromtest
"$earthly" ./old-behaviour-from+myimage
docker inspect myimage:fromtest >/dev/null
docker inspect mysubimage:fromtest >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by from ==="
docker rmi -f myimage:old-remote-from-test
docker rmi -f earthly-test-saveimage:test
"$earthly" ./old-behaviour-remote+myimage-fromtest
docker inspect myimage:old-remote-from-test >/dev/null
docker inspect earthly-test-saveimage:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by build ==="
docker rmi -f myimage:old-remote-build-test
docker rmi -f earthly-test-saveimage:test
"$earthly" ./old-behaviour-remote+myimage-buildtest
docker inspect myimage:old-remote-build-test >/dev/null
docker inspect earthly-test-saveimage:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour with remote target referenced by copy ==="
docker rmi -f myimage:old-remote-copy-test
docker rmi -f earthly-test-saveimage-and-artifact:test
"$earthly" ./old-behaviour-remote+myimage-copytest
docker inspect myimage:old-remote-copy-test >/dev/null
docker inspect earthly-test-saveimage-and-artifact:test >/dev/null

echo "=== ($LINENO): Test Old Behaviour on directly referenced remote target ==="
docker rmi -f earthly-test-saveimage:test
"$earthly" github.com/earthly/test-remote/save-image:main+saveimage
docker inspect earthly-test-saveimage:test >/dev/null


echo "=== ($LINENO): Test New Behaviour ==="

docker rmi -f myimage:test
docker rmi -f mysubimage:test
docker rmi -f earthly-test-saveimage:test
"$earthly" ./new-behaviour+myimage

docker inspect myimage:test >/dev/null

if docker inspect mysubimage:latest; then
	echo "ERROR ($LINENO): mysubimage should not have been saved1."
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via from ==="

docker rmi -f myimage:fromtest
docker rmi -f earthly-test-saveimage:test

"$earthly" -V ./new-behaviour-remote+myimage-fromtest
docker inspect myimage:fromtest >/dev/null
if docker inspect earthly-test-saveimage:test >/dev/null; then
    echo "ERROR ($LINENO): earthly-test-saveimage:test should not have been saved."   # THIS IS LEGIT BROKEN
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via copy ==="

docker rmi -f myimage:copytest
docker rmi -f earthly-test-saveimage:test

"$earthly" ./new-behaviour-remote+myimage-copytest
docker inspect myimage:copytest >/dev/null
if docker inspect earthly-test-saveimage:test >/dev/null; then
	echo "ERROR ($LINENO) : earthly-test-saveimage:test should not have been save3."
    exit 1
fi

echo "=== ($LINENO): Test New Behaviour referencing remote via build ==="

docker rmi -f myimage:buildtest
docker rmi -f earthly-test-saveimage:test

"$earthly" ./new-behaviour-remote+myimage-buildtest
docker inspect myimage:buildtest >/dev/null
docker inspect earthly-test-saveimage:test >/dev/null

echo "save-images test passed"
