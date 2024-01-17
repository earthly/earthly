#!/bin/sh
set -e
test -n "$earthly_config"

# Start the stub server; it will create a file under /server-got-a-connection if any connections are received.
# If this file is created, we know that earthly attempted to connect to the server (which means the DO_NOT_TRACK setting was ignored).
/bin/api-earthly-stub-server

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

# check test result
acbtest "$exit_code" = "0"
