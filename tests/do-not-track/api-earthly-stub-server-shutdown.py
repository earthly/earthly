#!/usr/bin/env python3
import socket
import os
import time
import signal
import sys

host = '127.0.0.1'
port = 443

pidfile='/do-not-track-tracker.pid'
server_got_a_connection_path = '/server-got-a-connection'

if os.path.exists(server_got_a_connection_path):
    print('ERROR: A connection was made, when it should not have')
    sys.exit(1)

# next make sure the server is still working
# if it is *not* working, then this test is invalid, since we wouldn't have
# detected if earthly ever attempted to connect to it

try:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((host, port))
        s.sendall(b"Hello, world")
except Exception as e:
    print('ERROR: stub server malfunction; failed to connect to stub server: {e}')
    sys.exit(1)

time.sleep(0.1)

if not os.path.exists('/server-got-a-connection'):
    print('ERROR: stub server malfunction; it should have created a file but didnt. The results of the DO_NOT_TRACK test can not be trusted')
    sys.exit(1)

try:
    pid = int(open(pidfile, 'r').read())
    os.kill(pid, signal.SIGKILL)
except Exception as e:
    print('WARN: failed to shutdown api-earthly-stub-server: {e}')
