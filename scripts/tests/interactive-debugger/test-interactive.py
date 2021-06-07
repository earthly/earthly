#!/usr/bin/env python3
import argparse
import pexpect
import time
import sys
import io
import os
import shlex

script_dir = os.path.dirname(os.path.realpath(__file__))

def print_box(s):
    n = len(s)
    print('-'*n)
    print(s)
    print('-'*n)

def test_interactive(earthly_path, timeout):
    earthly_path = os.path.realpath(earthly_path)

    print(f'Running interactive-debugger test against "{earthly_path}"')
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
                raise RuntimeError('failed to find magicstring in output (indicating the echo + base64 decode failed to run)')

            assert c.isalive()

            c.sendline("exit")
            c.wait()

            assert not c.isalive()
            print_box('interactive debugger test passed')
        except Exception as e:
            print_box('interactive test failed')
            print(f'{e}')

            print_box('pexpect debug information')
            print(str(c))
            exit_code = 1
        finally:
            print_box('earthly output')
            s = ''.join(ch for ch in output.getvalue() if ch.isprintable() or ch == '\n')
            print(s)
    finally:
        os.chdir(cwd)

    return exit_code


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-e', '--earthly', help="earthly binary to run test against", default='earthly')
    parser.add_argument('-t', '--timeout', help="fail test if it takes longer than this many seconds", type=float, default=30.0)
    args = parser.parse_args()
    exit_code = test_interactive(args.earthly, args.timeout)
    print(f'exiting with status code {exit_code}')
    sys.exit(exit_code)
