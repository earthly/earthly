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
    print('A connection was made, when it should not have')
    sys.exit(1)

# next make sure the server is still working

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((host, port))
    s.sendall(b"Hello, world")

time.sleep(0.1)

if not os.path.exists('/server-got-a-connection'):
    print('stub server malfunction; it should have created a file but didnt. The results of the DO_NOT_TRACK test can not be trusted')
    sys.exit(1)


pid = int(open(pidfile, 'r').read())
os.kill(pid, signal.SIGKILL)
