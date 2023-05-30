#!/usr/bin/env bash
set -eu

echo "running tests with $earthly"

earthly=${earthly:=earthly}
earthly=$(realpath "$earthly")
frontend="${frontend:-$(which docker)}"
tag=test-image:v1.2.3
testdir=/tmp/earthly-docker-build-test
mkdir -p $testdir

function reset() {
    chmod 777 -R $testdir >/dev/null 2>&1 || true
    rm -rf output output2 $testdir/Earthfile $testdir/.dockerignore $testdir/.earthlyignore
    $frontend rmi -f $tag >/dev/null 2>&1
}

function cleanup() {
    reset
    rm -rf $testdir
}
trap cleanup EXIT

function run_test_cmd() {
  cmd=$1
  if  eval NO_COLOR=0 "$cmd" > output 2>&1; then
      echo "earthly docker-build should have failed"
      exit 1
  fi
}

echo "=== test 1 - command fails when tag is not specified:"
reset
run_test_cmd "\"$earthly\" docker-build"
tail -n1 output > output2
diff output2 <(echo "Error: Required flag \"tag\" not set")

echo "=== test 2 - command fails when build context is not specified:"
reset
run_test_cmd "\"$earthly\" docker-build -t $tag"
tail -n1 output > output2
diff output2 <(echo "Error: no build context path provided. Try earthly docker-build <path>")

echo "=== test 3 - command fails when more than one build context (non arg flag) is specified:"
reset
run_test_cmd "\"$earthly\" docker-build -t $tag ctx1 ctx2"
tail -n1 output > output2
diff output2 <(echo "Error: invalid arguments ctx1 ctx2")

echo "=== test 4 - command fails when it cannot check if Earthfile already exists:"
reset
chmod 000 $testdir
run_test_cmd "\"$earthly\" docker-build -t $tag $testdir"
diff output <(echo "Error: failed to check if \"$testdir/Earthfile\" exists: unable to stat $testdir/Earthfile: stat $testdir/Earthfile: permission denied")

echo "=== test 5 - command fails when the Earthfile already exists:"
reset
touch $testdir/Earthfile

run_test_cmd "\"$earthly\" docker-build -t $tag $testdir"
diff output <(echo "Error: earthfile already exists; please delete it if you wish to continue")

echo "=== test 6 - command fails when .dockerignore is inaccessible:"
reset
touch $testdir/.dockerignore
chmod 000 $testdir/.dockerignore

run_test_cmd "\"$earthly\" docker-build -t $tag $testdir"
diff output <(echo "Error: failed to copy \"$testdir/.dockerignore\" to \"$testdir/.earthlyignore\": open $testdir/.dockerignore: permission denied")

echo "=== test 7 - command fails when it cannot create the Earthfile:"
reset
chmod -w $testdir

run_test_cmd "\"$earthly\" docker-build -t $tag $testdir"
diff output <(echo "Error: docker-build: failed to create Earthfile \"$testdir/Earthfile\": open $testdir/Earthfile: permission denied")

# happy paths:
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

echo "=== test 8 - it creates an image:"
reset

"$earthly" docker-build -t "$tag" $testdir

"$frontend" inspect "$tag" > /dev/null

echo "=== test 9 - it uses the correct target:"

# use target1:
reset

"$earthly" docker-build -t "$tag" --target target1 $testdir

"$frontend" run --rm "$tag" |grep target1 > output 2>&1

diff output <(echo "target1")

# use target2:
reset

"$earthly" docker-build -t "$tag" --target target2 $testdir

"$frontend" run --rm "$tag" |grep target2 > output 2>&1

diff output <(echo "target2")

echo "=== test 10 - it uses the correct arg value:"

# use override-value:
reset

"$earthly" docker-build -t "$tag" $testdir --arg_to_override=override-value

"$frontend" run --rm "$tag" |grep override-value > output 2>&1

diff output <(echo "override-value")

# use default-value:
reset

"$earthly" docker-build -t "$tag" $testdir

"$frontend" run --rm "$tag" |grep default-value > output 2>&1

diff output <(echo "default-value")

echo "=== test 11 - it builds the image using the correct platforms:"

reset

"$earthly" docker-build -t "$tag" --platform linux/amd64 --platform linux/arm64 $testdir

# arm64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_arm64" > output

diff output <(echo "arm64")

# amd64:
"$frontend" inspect --format='{{.Architecture}}' "${tag}_linux_amd64" > output

diff output <(echo "amd64")

echo "=== test 12 - it ignores files according to .dockerignore:"

reset

# create .dockerignore file
cat << EOF > $testdir/.dockerignore
bad.txt
EOF

"$earthly" docker-build -t "$tag" $testdir

"$frontend" run --rm "$tag" > output

grep bad.txt output > output2 || true

diff output2 <(echo -n "")

echo "=== All tests have passed ==="
