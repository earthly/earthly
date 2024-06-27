#!/bin/sh
set -e
test -n "$earthly_config"

# the stub server is sometimes failing with a "OSError: [Errno 98] Address in use", let's try to see if anything else is listening.
echo "== netstat before =="
netstat -nltp
echo "== done =="

# Start the stub server; it will create a file under /server-got-a-connection if any connections are received.
# If this file is created, we know that earthly attempted to connect to the server (which means the DO_NOT_TRACK setting was ignored).
/bin/api-earthly-stub-server

echo "== netstat after stub server =="
sleep 3
netstat -nltp
echo "== done =="

# prevent earthly auto-login
earthly --config $earthly_config account logout

# build a target
earthly --config $earthly_config +true

# shutdown the stub-server
set +e
echo ""
echo "== shutting down api-earthly-stub-server =="
api-earthly-stub-server-shutdown
export exit_code="$?"
set -e

# display logs
echo ""
echo "== api-earthly-stub-server log =="
cat /var/log/do-not-track-server.log
echo ""

ls -la /tmp

# check test result
acbtest "$exit_code" = "0"
