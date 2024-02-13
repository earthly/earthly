#!/usr/bin/env python3
import argparse
import pexpect
import time
import sys
import io
import os
import shlex

script_dir = os.path.dirname(os.path.realpath(__file__))

def test_interactive(earthly_path, timeout):
    output = io.StringIO()

    # change dir to where our test Earthfile is
    cwd = os.getcwd()
    os.chdir(script_dir)

    exit_code = 0
    try:
        c = pexpect.spawn(f'{shlex.quote(earthly_path)} -Pi +fail-with-docker-compose', encoding='utf-8', timeout=timeout)
        c.logfile_read = output
        try:
            c.expect('Entering interactive debugger')

            # give the shell time to startup, otherwise stdin might be lost
            time.sleep(5)

            # test we can obtain decoded text from the rot13 echo server that is created via docker-compose
            c.sendline('docker exec rot13 sh -c \'(echo guvf vf zl Frpe3g Z3ff4tr; sleep 1) | ncat localhost 5432\'')
            try:
                c.expect('this is my Secr3t M3ss4ge', timeout=timeout)
            except Exception as e:
                raise RuntimeError('failed to find "this is my Secr3t M3ss4ge" in output (indicating we were unable to obtain decoded text from the rot13 echo server that was started via docker compose)')

            # decode the data.txt file to ensure the debugger is still running
            c.sendline("cat /data.txt | base64 -d")
            try:
                c.expect('e88cc2b8-c179-4ed7-8eae-207a0ef7546c', timeout=1)
            except Exception as e:
                raise RuntimeError('failed to find uuid in output (indicating the cat + base64 decode failed to run)')

            c.sendline("exit")
            c.wait()

            assert not c.isalive()
        except pexpect.exceptions.TIMEOUT as e:
            print('ERROR: interactive test timed out')
            exit_code = 2
        except Exception as e:
            print(f'ERROR: interactive test failed with {e}')
            exit_code = 1
        finally:
            print('--------------')
            print('earthly output')
            s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
            print(s)
            if exit_code:
                print('additional pexpect debug information')
                print(str(c)+'\n')
    finally:
        os.chdir(cwd)

    return exit_code
