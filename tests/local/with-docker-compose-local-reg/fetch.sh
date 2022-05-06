#!/bin/sh
set -e
curl --connect-timeout 1 \
    --max-time 1 \
    --retry 5 \
    --retry-delay 0 \
    --retry-max-time 10 \
    "http://${WEBHOST}:${WEBPORT}"
