#!/usr/bin/env python3
import argparse
import pexpect
import time
import sys
import io
import os
import shlex

script_dir = os.path.dirname(os.path.realpath(__file__))

def test_with_exit_code(earthly_path, timeout, exit_code_to_send):
    output = io.StringIO()

    # change dir to where our test Earthfile is
    cwd = os.getcwd()
    os.chdir(script_dir)

    exit_code = 0
    try:
        c = pexpect.spawn(f'{shlex.quote(earthly_path)} +interactive-target', encoding='utf-8', timeout=timeout)
        c.logfile_read = output
        try:
            try:
                question_text = 'What exit code do you want'
                c.expect(question_text, timeout=timeout)
            except Exception as e:
                raise RuntimeError(f'failed to find {question_text} in output')

            print(f'sending "{exit_code_to_send}" as response to "What exit code do you want" prompt')
            c.sendline(str(exit_code_to_send))

            if exit_code_to_send:
                expected_text = f'did not complete successfully. Exit code {exit_code_to_send}'
                try:
                    c.expect(expected_text, timeout=10)
                except Exception as e:
                    raise RuntimeError(f'failed to find text: {expected_text}')

            status = c.wait()
            print(f'earthly exited with code={status}')
            if exit_code_to_send and not status:
                raise RuntimeError(f'earthly should have failed (due to requested exit {exit_code_to_send}), but didnt')
            if not exit_code_to_send and status:
                raise RuntimeError('earthly should not have failed')

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

def test_interactive(earthly_path, timeout):
    for code in (0, 73):
        print(f'running test_with_exit_code with code={code}')
        result = test_with_exit_code(earthly_path, timeout, code)
        if result:
            return result
    return 0
