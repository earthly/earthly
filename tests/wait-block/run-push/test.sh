#!/usr/bin/env bash
set -uex
set -o pipefail

# Unset referenced-save-only.
export EARTHLY_VERSION_FLAG_OVERRIDES=""

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}
"$earthly" --version

echoserver="earthly-test-echoserver"

# clean up previous run
docker kill "$echoserver" || true
docker rm "$echoserver" || true


# display a pass/fail message at the end
function finish {
  status="$?"
  if [ "$status" = "0" ]; then
    echo "run-push test passed"
  else
    echo "run-push test failed with $status"
  fi
}
trap finish EXIT

# Get host IP (it must be accessible from both docker and runc)
HOST_IP=$(ip route get 8.8.8.8 | awk -F"src " 'NR==1{split($2,a," ");print a[1]}')

# pick the first free port starting at 5000 (up to 5050)
PORT_FOUND="false"
for i in {0..50}
do
    PORT="$(( 5000 + i ))"
    ACTIVE_PORT="$(netstat -lnt | awk '{print $4}' | (grep ":$PORT\$" || true) | wc -l)"
    if [ "$ACTIVE_PORT" = "0" ]; then
        PORT_FOUND="true"
        break
    fi
done
test "$PORT_FOUND" = "true"

# run a tcp server
# NOTE: there is a bug in the busybox implementation of netcat which prevents second connections
# from being established; instead, the regular BSD-version is used instead:
# nc:
#     FROM debian:11.4
#     RUN apt-get update && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*
#     SAVE IMAGE --push alexcb132/netcat
docker pull alexcb132/netcat
docker run \
    -d \
    --network=host \
    --name "$echoserver" alexcb132/netcat \
    /bin/sh -c "nc -l -k -p $PORT"
timeout 5 sh -c 'until nc -z $0 $1; do sleep 1; done' 127.0.0.1 $PORT


echo "===test1===" > "/dev/tcp/127.0.0.1/$PORT"

"$earthly" $@ +test --echoserver_ip="$HOST_IP" --echoserver_port="$PORT"

diff <(docker logs "$echoserver" | grep -A 999 '===test1===') <(cat <<EXPECTED
===test1===
no-cache-1
no-cache-2
EXPECTED
)

echo "===test2===" > "/dev/tcp/127.0.0.1/$PORT"

"$earthly" --push $@ +test --echoserver_ip="$HOST_IP" --echoserver_port="$PORT"

diff <(docker logs "$echoserver" | grep -A 999 '===test2===') <(cat <<EXPECTED
===test2===
no-cache-1
run-push
no-cache-2
EXPECTED
)
