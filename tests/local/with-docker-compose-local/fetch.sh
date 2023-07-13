#!/bin/sh
set -e
sleep 10
output="$(mktemp)"
curl --connect-timeout 1 \
    --max-time 1 \
    --retry 5 \
    --retry-delay 0 \
    --retry-max-time 10 \
    "http://${WEBHOST}:${WEBPORT}" | tee "$output"
echo "curl command worked"
echo "checking $output contains the magic string"
grep "Hello World" "$output"
echo "grep command worked"
