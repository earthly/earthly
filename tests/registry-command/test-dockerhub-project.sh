#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared project's secrets)
# the lock.sh and unlock.sh scripts must first be run

ORG="ryan-test"
PROJECT="registry-command-test-project"

clearprojectsecrets() {
    earthly secrets --org "$ORG" --project "$PROJECT" ls /user/std/registry | xargs -r -n 1 earthly secrets --org "$ORG" --project "$PROJECT" rm
}

# clear out secrets from previous test
clearprojectsecrets

# test dockerhub credentials do not exist
earthly registry --org "$ORG" --project "$PROJECT" list | grep -v registry-1.docker.io

# set dockerhub credentials
earthly registry --org "$ORG" --project "$PROJECT" setup --username myprojecttest --password keepitsecret

# test dockerhub credentials exist
earthly registry --org "$ORG" --project "$PROJECT" list | grep registry-1.docker.io

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/registry-1.docker.io/username)" = "myprojecttest"
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/registry-1.docker.io/password)" = "keepitsecret"

# test a different host
echo -n keepitsecret2  | earthly registry --org "$ORG" --project "$PROJECT" setup --username myprojecttest2 --password-stdin corp-registry.earthly.dev

# both dockerhub and corp-registry should exist
earthly registry --org "$ORG" --project "$PROJECT" list | grep registry-1.docker.io
earthly registry --org "$ORG" --project "$PROJECT" list | grep corp-registry.earthly.dev

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/registry-1.docker.io/username)" = "myprojecttest"
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/registry-1.docker.io/password)" = "keepitsecret"
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/corp-registry.earthly.dev/username)" = "myprojecttest2"
test "$(earthly secrets --org "$ORG" --project "$PROJECT" get std/registry/corp-registry.earthly.dev/password)" = "keepitsecret2"

earthly registry --org "$ORG" --project "$PROJECT" remove
earthly registry --org "$ORG" --project "$PROJECT" list | grep -v registry-1.docker.io

clearprojectsecrets
