#!/bin/sh
set -ex

# WARNING -- RACE-CONDITION: this test is not thread-safe (since it makes use of a shared user's secrets)
# the lock.sh and unlock.sh scripts must first be run

clearusersecrets() {
    earthly secrets ls /user/std/ | xargs -r -n 1 earthly secrets rm
}

# clear out secrets from previous test
clearusersecrets

# test dockerhub credentials do not exist
earthly registry list | grep -v registry-1.docker.io

# set dockerhub credentials
earthly registry setup --username mytest --password keepitsafe

# test dockerhub credentials exist
earthly registry list | grep registry-1.docker.io

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/username)" = "mytest"
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/password)" = "keepitsafe"

earthly registry remove
earthly registry list | grep -v registry-1.docker.io

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
