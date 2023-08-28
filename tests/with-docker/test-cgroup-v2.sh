#!/bin/sh
set -e
insidecode="$(base64 -w0 test-cgroup-v2-inside-container.sh)"
test -n "$insidecode"
test -n "$ubuntu_img_tag"
docker run --privileged -t --name foo "ubuntu:$ubuntu_img_tag" /bin/sh -c "echo $insidecode | base64 -d > test-cgroup && chmod +x test-cgroup && ./test-cgroup"
