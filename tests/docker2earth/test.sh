#!/bin/bash
set -eu

earthly=$(pwd)/build/linux/amd64/earthly
dockerfiles="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "=== Testing Dockerfile1 ==="
cd "$(mktemp -d)"
echo "working out of $(pwd)"
cp "$dockerfiles"/Dockerfile1 Dockerfile
$earthly docker2earthly --tag=myimage:latest
$earthly +build
docker run --rm myimage:latest say-hi | grep hello

echo "=== Testing Dockerfile2 ==="
cd "$(mktemp -d)"
echo "working out of $(pwd)"
cat "$dockerfiles"/Dockerfile2 | $earthly docker2earthly --dockerfile - --tag myotherimage:test
cp "$dockerfiles"/app.go .
$earthly +build
docker run --rm myotherimage:test | grep greetings

echo "=== Testing args-before-from.Dockerfile ==="
cd "$(mktemp -d)"
echo "working out of $(pwd)"
$earthly docker2earthly --dockerfile - --tag onemoreimage:test < "$dockerfiles"/args-before-from.Dockerfile
cp "$dockerfiles"/app.go .
$earthly +build --BASE=golang --GO_MAJOR=1
docker run --rm onemoreimage:test | grep greetings
