#!/bin/bash
set -eu

earth=`pwd`/build/linux/amd64/earth
dockerfiles="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo === Testing Dockerfile1 ===
cd $(mktemp -d)
echo "working out of $(pwd)"
cp $dockerfiles/Dockerfile1 Dockerfile
$earth docker2earth
$earth +build
docker run --rm myimage:latest say-hi | grep hello

echo === Testing Dockerfile2 ===
cd $(mktemp -d)
echo "working out of $(pwd)"
cat $dockerfiles/Dockerfile2 | $earth docker2earth --dockerfile - --tag myotherimage:test
cp $dockerfiles/app.go .
$earth +build
docker run --rm myotherimage:test | grep greetings
