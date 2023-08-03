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

# set dockerhub credentials via stdin
echo -n "fromstdin" | earthly registry setup --username mytest2 --password-stdin

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/username)" = "mytest2"
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/password)" = "fromstdin"

# test no extra newline was stored; note that "echo -n fromstdin | md5sum" = 4b1fb3bf88ee25da648fefd5af81c921
earthly secrets get -n /user/std/registry/registry-1.docker.io/password | md5sum | grep 4b1fb3bf88ee25da648fefd5af81c921

# set dockerhub credentials via tty
/usr/bin/expect -c '
spawn earthly registry setup
expect "username: "
send "mytest3\n"
expect "password: "
send "fromexpect\n"
expect eof
'

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/username)" = "mytest3"
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/password)" = "fromexpect"

# test no extra newline was stored; note that "echo -n fromexpect | md5sum" = bd62328338f2f6a8cb8adf2e3712afad
earthly secrets get -n /user/std/registry/registry-1.docker.io/password | md5sum | grep bd62328338f2f6a8cb8adf2e3712afad

# set dockerhub credentials via tty
/usr/bin/expect -c '
spawn earthly registry setup --username mytest4
expect "password: "
send "fromexpect2\n"
expect eof
'

# test username and password were correctly stored in underlying std secret
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/username)" = "mytest4"
test "$(earthly secrets get /user/std/registry/registry-1.docker.io/password)" = "fromexpect2"

# test no extra newline was stored; note that "echo -n fromexpect2 | md5sum" = d581f3b642ece7e7b559b8a73c60aeae
earthly secrets get -n /user/std/registry/registry-1.docker.io/password | md5sum | grep d581f3b642ece7e7b559b8a73c60aeae

earthly registry remove
earthly registry list | grep -v registry-1.docker.io

# clear out secrets (just in case project-based registry accidentally uses user-based)
clearusersecrets
