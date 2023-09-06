#!/bin/sh
set -e

configpath="/etc/.earthly/config.yml"

if [ -f "$configpath" ]; then
  echo >&2 "error: $configpath already exists, unable to setup registry"
  exit 1
fi

mkdir -p "$(dirname "$configpath")"
cat>"$configpath"<<EOF
global:
  disable_analytics: true
EOF

if [ -n "$DOCKERHUB_MIRROR" ]; then
  # create a mirror entry for dockerhub (aka docker.io)
  cat>>"$configpath"<<EOF
  buildkit_additional_config: |
    [registry."docker.io"]
      mirrors = ["$DOCKERHUB_MIRROR"]
EOF
  # create a second registry config for the mirror if either insecure or http flags must be set
  if [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ] || [ "$DOCKERHUB_MIRROR_HTTP" = "true" ]; then
    cat>>"$configpath"<<EOF
    [registry."$DOCKERHUB_MIRROR"]
EOF
    if [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ]; then
    cat>>"$configpath"<<EOF
      insecure = true
EOF
    fi
    if [ "$DOCKERHUB_MIRROR_HTTP" = "true" ]; then
    cat>>"$configpath"<<EOF
      http = true
EOF
    fi
  fi
  if [ "$DOCKERHUB_MIRROR_AUTH" = "true" ]; then
    if [ -z "$DOCKERHUB_MIRROR" ]; then
      echo >&2 "error: expected value for DOCKERHUB_MIRROR, but got none"
      exit 1
    fi
    if [ -z "$DOCKERHUB_MIRROR_USER" ]; then
      echo >&2 "error: expected value for DOCKERHUB_MIRROR_USER, but got none"
      exit 1
    fi
    if [ -z "$DOCKERHUB_MIRROR_PASS" ]; then
      echo >&2 "error: expected value for DOCKERHUB_MIRROR_PASS, but got none"
      exit 1
    fi
    docker login "$DOCKERHUB_MIRROR" --username="$DOCKERHUB_MIRROR_USER" --password="$DOCKERHUB_MIRROR_PASS"
  fi
else
  echo >&2 "WARNING: no dockerhub mirror has been setup; you may get rate limited"
fi
