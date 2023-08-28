#!/bin/sh
set -e

# this script is run within a container running under a WITH DOCKER, which validates the container receives a correctly setup cgroup

if ! [ -f /sys/fs/cgroup/cgroup.type ]; then
  echo "skip: test requires cgroup v2" && exit 0;
fi
if [ "$(cat /sys/fs/cgroup/cgroup.type)" = "domain" ]; then
  echo "pass: cgroup type is domain"
else
  echo "fail: expected domain cgroup type but got $(cat  /sys/fs/cgroup/cgroup.type)"
  exit 1
fi

if grep cpu /sys/fs/cgroup/cgroup.controllers >/dev/null; then
  echo "pass: cpu controller is available in cgroup.controllers"
else
  echo "fail: cpu controller is NOT available in cgroup.controllers"
  exit 1
fi
