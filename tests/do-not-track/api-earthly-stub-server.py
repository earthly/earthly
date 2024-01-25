#!/usr/bin/env python3
import socket
import os
import sys
import fcntl
import time
from contextlib import suppress

host = '127.0.0.1'
port = 443

server_got_a_connection_path = '/server-got-a-connection'

# This starts a stub server in the background via a double-fork
pidfile='/do-not-track-tracker.pid'
stdin='/dev/null'
stdout='/dev/null'
stderr='/var/log/do-not-track-server.log'

ready_pipe_r, ready_pipe_w = os.pipe()

# first fork
pid = os.fork()
if pid > 0:
    os.close(ready_pipe_w)
    fcntl.fcntl(ready_pipe_r, fcntl.F_SETFL, os.O_NONBLOCK)
    num_attemps_remaining = 10
    while True:
        try:
            msg = os.read(ready_pipe_r, 1024).decode('utf8')
        except BlockingIOError as e:
            num_attemps_remaining -= 1
            if num_attemps_remaining <= 0:
                print('server failed to start')
                sys.exit(1)
            print('waiting for stub-server to start')
            time.sleep(1)
            continue
        msg = msg
        break
    if msg == 'ready':
        print('stub-server ready')
        sys.exit(0)
    print(f'unexpected msg "{msg}" while waiting for server to start')
    sys.exit(1)

os.close(ready_pipe_r)

try:
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
    os.write(ready_pipe_w, 'ready'.encode('utf8'))
    os.close(ready_pipe_w)

    with suppress(FileNotFoundError):
        os.remove(server_got_a_connection_path)

    print(f'creating socket', file=sys.stderr)
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        print(f'binding socket to {host}:{port}', file=sys.stderr)
        s.bind((host, port))
        print(f'listening on socket', file=sys.stderr)
        s.listen()
        conn, addr = s.accept()
        print(f'received connection from {addr}', file=sys.stderr)
        with conn:
            with open(server_got_a_connection_path, 'w') as f:
                f.write('this should not have happened')
except Exception as e:
    # log to stderr
    print(f'unexpected exception {e}', file=sys.stderr)

    # send the exception back over the ready_pipe (so the initial process can display the error if it is still running)
    os.write(ready_pipe_w, f'unexpected exception while starting server: {e}'.encode('utf8'))
    os.close(ready_pipe_w)
    raise
