#!/usr/bin/env python3
import socket
import os
import sys
from contextlib import suppress

host = '127.0.0.1'
port = 443

server_got_a_connection_path = '/server-got-a-connection'

# This starts a stub server in the background via a double-fork
pidfile='/do-not-track-tracker.pid'
stdin='/dev/null'
stdout='/dev/null'
stderr='/dev/null'

# first fork
pid = os.fork()
if pid > 0:
    sys.exit(0)

os.chdir('/')
os.setsid()
os.umask(0)

# second fork
pid = os.fork()
if pid > 0:
    sys.exit(0)

# redirect stdio
sys.stdout.flush()
sys.stderr.flush()
si = os.open(stdin, os.O_RDWR)
so = os.open(stdout, os.O_WRONLY|os.O_TRUNC|os.O_CREAT)
se = os.open(stderr, os.O_WRONLY|os.O_TRUNC|os.O_CREAT)
os.dup2(si, sys.stdin.fileno())
os.dup2(so, sys.stdout.fileno())
os.dup2(se, sys.stderr.fileno())

# write pid to disk
pid = str(os.getpid())
with open(pidfile, 'w') as f:
    f.write(pid)

# daemon ready to go

with suppress(FileNotFoundError):
    os.remove(server_got_a_connection_path)

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((host, port))
    s.listen()
    conn, addr = s.accept()
    with conn:
        with open(server_got_a_connection_path, 'w') as f:
            f.write('this should not have happened')
