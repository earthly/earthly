#!/bin/bash

set -xeu

function cleanup() {
    jobs="$(jobs -p)"
    if [ -n "$jobs" ]
    then
        # shellcheck disable=SC2086 # Intended splitting of 
        kill $jobs
    fi
    wait
}
trap cleanup EXIT

export EARTHLY_CONVERSION_PARALLELISM=5

EARTHLY_VERSION_FLAG_OVERRIDES="$(tr -d '\n' < .earthly_version_flag_overrides)"
export EARTHLY_VERSION_FLAG_OVERRIDES

case "$EARTHLY_OS" in
    darwin)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-darwin-amd64"
        earthly="./build/darwin/amd64/earthly"
        ;;

    darwin-m1)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-darwin-arm64"
        earthly="./build/darwin/arm64/earthly"
        ;;

    linux)
        download_url="https://github.com/earthly/earthly/releases/latest/download/earthly-linux-amd64"
        earthly="./build/linux/amd64/earthly"
        ;;
esac

echo "The detected architecture of the runner is $(uname -m)"

echo "Add branch info back to git (Earthly uses it for tagging)"
git checkout -B "$BUILDKITE_BRANCH" || true

echo "Download latest Earthly binary"
if [ -n "$download_url" ]; then
    curl -o ./earthly-released -L "$download_url" && chmod +x ./earthly-released
    released_earthly=./earthly-released
fi

echo "Build latest earthly using released earthly"
"$released_earthly" --version
"$released_earthly" config global.disable_analytics true
"$released_earthly" +for-"$EARTHLY_OS"
chmod +x "$earthly"

# WSL2 sometimes gives a "Text file busy" when running the native binary, likely due to crossing the WSL/Windows divide.
# This should be enough retry to skip that, and fail if theres _actually_ a problem.
att_max=5
att_num=1
until "$earthly" --version || (( att_num == att_max ))
do
    echo "Attempt $att_num failed! Trying again in $att_num seconds..."
    sleep $(( att_num++ ))
done

export EARTHLY_VERSION_FLAG_OVERRIDES="referenced-save-only"
"$earthly" config global.local_registry_host 'tcp://127.0.0.1:8371'

# Yes, there is a bug in the upstream YAML parser. Sorry about the jank here.
# https://github.com/go-yaml/yaml/issues/423
"$earthly" config global.buildkit_additional_config "'[registry.\"docker.io\"]

 mirrors = [\"registry-1.docker.io.mirror.corp.earthly.dev\"]'"

echo "Execute tests"
"$earthly" --ci -P \
    --build-arg DOCKERHUB_AUTH=true \
    --build-arg DOCKERHUB_USER_SECRET=+secrets/earthly-technologies/dockerhub-mirror/user \
    --build-arg DOCKERHUB_TOKEN_SECRET=+secrets/earthly-technologies/dockerhub-mirror/pass \
    --build-arg DOCKERHUB_MIRROR=registry-1.docker.io.mirror.corp.earthly.dev \
  +test

echo "Execute fail test"
bash -c "! $earthly --ci ./tests/fail+test-fail"
