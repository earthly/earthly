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
        c = pexpect.spawn(f'{shlex.quote(earthly_path)} -i +fail', encoding='utf-8', timeout=timeout)
        c.logfile_read = output
        try:
            c.expect('Entering interactive debugger')

            # give the shell time to startup, otherwise stdin might be lost
            time.sleep(0.5)

            # /data.txt contains bWFnaWNzdHJpbmcK which is the base64 encoded value of "magicstring\n"
            c.sendline("cat /data.txt | base64 -d")

            try:
                c.expect('magicstring', timeout=1)
            except Exception as e:
                raise RuntimeError('failed to find magicstring in output (indicating the cat + base64 decode failed to run)')

            assert c.isalive()

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
