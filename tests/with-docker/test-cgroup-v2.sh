#!/bin/sh
set -e
insidecode="$(base64 -w0 test-cgroup-v2-inside-container.sh)"
test -n "$insidecode"
docker run --privileged -t --name foo ubuntu:23.04 /bin/sh -c "echo $insidecode | base64 -d > test-cgroup && chmod +x test-cgroup && ./test-cgroup"
