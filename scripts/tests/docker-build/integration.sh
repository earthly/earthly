#!/usr/bin/env bash
set -eu

earthly=${earthly:=earthly}
if [ "$earthly" != "earthly" ]; then
  earthly=$(realpath "$earthly")
fi
echo "running tests with $earthly"
"$earthly" --version
frontend="${frontend:-$(which docker)}"
tag=test-image:v1.2.3
tag2="${tag}.4"

PATH="$(realpath "$(dirname "$0")/../../acbtest"):$PATH"

testdir=/tmp/earthly-docker-build-test
mkdir -p $testdir

function reset() {
    chmod 777 -R $testdir >/dev/null 2>&1 || true
    rm -rf output output2 $testdir/Earthfile $testdir/.dockerignore $testdir/.earthlyignore
    $frontend rmi -f $tag $tag2 ${tag}_linux_amd64 ${tag}_linux_arm64 >/dev/null 2>&1
}

function cleanup() {
    reset
    rm -rf $testdir
}
trap cleanup EXIT

function run_test_cmd() {
  cmd=$1
  if  eval GITHUB_ACTIONS="" NO_COLOR=0 "$cmd" > output 2>&1; then
      echo "earthly docker-build should have failed"
      exit 1
  fi
}

echo "=== test 1 - command fails when build context is not specified:"
reset
run_test_cmd "\"$earthly\" docker-build"
tail -n1 output > output2
diff output2 <(echo "Error: no build context path provided. Try earthly docker-build <path>")

echo "=== test 2 - command fails when more than one build context (non arg flag) is specified:"
reset
run_test_cmd "\"$earthly\" docker-build ctx1 ctx2"
tail -n1 output > output2
diff output2 <(echo "Error: invalid arguments ctx1 ctx2")

## happy paths:
reset
# create Dockerfile for the happy paths
cat << EOF > $testdir/Dockerfile
FROM alpine as target1
ARG arg_to_override=default-value
WORKDIR /app
COPY good.txt bad.txt* .
RUN ls -1 *.txt > output
RUN echo \$arg_to_override >> output
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["cat output;echo target1"]

FROM target1 as target2
CMD ["cat output;echo target2"]
EOF

touch $testdir/good.txt $testdir/bad.txt

echo "=== test 3 - it creates an image:"
reset

"$earthly" docker-build -t "$tag" $testdir

"$frontend" inspect "$tag" > /dev/null

echo "=== test 4 - it creates an image with multiple tags:"
reset

"$earthly" docker-build -t "$tag" -t "$tag2" $testdir

"$frontend" inspect "$tag" > /dev/null
"$frontend" inspect "$tag2" > /dev/null

echo "=== test 5 - it creates an image without tags"
reset

"$earthly" docker-build $testdir

echo "=== test 6 - it uses the correct target:"

# use target1:
reset

"$earthly" docker-build -t "$tag" --target target1 $testdir

"$frontend" run --rm "$tag" | acbgrep target1 > output 2>&1

diff output <(echo "target1")

# use target2:
reset

"$earthly" docker-build -t "$tag" --target target2 $testdir

"$frontend" run --rm "$tag" | acbgrep target2 > output 2>&1

diff output <(echo "target2")

echo "=== test 7 - it uses the correct arg value:"

# use override-value:
reset

"$earthly" docker-build -t "$tag" $testdir --arg_to_override=override-value

"$frontend" run --rm "$tag" | acbgrep override-value > output 2>&1

diff output <(echo "override-value")

# use default-value:
reset

"$earthly" docker-build -t "$tag" $testdir

"$frontend" run --rm "$tag" | acbgrep default-value > output 2>&1

diff output <(echo "default-value")

echo "=== test 8 - it builds the image using the correct platforms (multiple flags):"

reset

"$earthly" docker-build -t "$tag" --platform linux/amd64 --platform linux/arm64 $testdir

# arm64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_arm64" > output

diff output <(echo "arm64")

# amd64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_amd64" > output

diff output <(echo "amd64")

echo "=== test 9 - it builds the image using the correct platforms (one flag, docker command style):"

reset

"$earthly" docker-build -t "$tag" --platform linux/amd64,linux/arm64 $testdir

# arm64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_arm64" > output

diff output <(echo "arm64")

# amd64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_amd64" > output

diff output <(echo "amd64")

echo "=== test 10 - it ignores files according to .dockerignore:"

reset

# create .dockerignore file
cat << EOF > $testdir/.dockerignore
bad.txt
EOF

"$earthly" docker-build -t "$tag" $testdir

"$frontend" run --rm "$tag" > output

if grep bad.txt output > /dev/null; then
    >&2 echo "failure: bad.txt was found in output"
    exit 1
fi

# verify file is removed
if [ -f $testdir/.earthlyignore ]; then
    >&2 echo "$testdir/.earthlyignore was not removed"
    exit 1
fi

echo "=== All tests have passed ==="
