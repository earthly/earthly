#!/usr/bin/env python3
import argparse
import pexpect
import time
import sys
import io
import os
import shlex
import importlib.util

script_dir = os.path.dirname(os.path.realpath(__file__))

def import_test_func(path, func_name='test_interactive'):
    module_name = os.path.basename(path)
    spec = importlib.util.spec_from_file_location(module_name, path)
    foo = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(foo)
    return getattr(foo, func_name)

def get_earthly_binary(earthly):
    if os.path.isfile(earthly):
        return earthly
    if os.path.sep not in earthly:
        for path in os.environ['PATH'].split(':'):
            fullpath = os.path.join(path, earthly)
            if os.path.isfile(fullpath):
                return fullpath
    raise RuntimeError(f'failed to find earthly binary: {earthly}')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-e', '--earthly', help="earthly binary to run test against", default='earthly')
    parser.add_argument('-t', '--timeout', help="fail test if it takes longer than this many seconds", type=float, default=30.0)
    args = parser.parse_args()

    earthly_path = os.path.realpath(get_earthly_binary(args.earthly))
    print(f'Running interactive tests against "{earthly_path}"')

    exit_code = 0
    for test_name, test in (
            ('test-simple', import_test_func(os.path.join(script_dir, 'simple', 'test-simple.py'))),
            ('test-interactive-run', import_test_func(os.path.join(script_dir, 'interactive-run', 'test-interactive-run.py'))),
            ('test-docker-compose', import_test_func(os.path.join(script_dir, 'docker-compose', 'test-docker-compose.py'))),
            ):
        print(f'Running {test_name}')
        test_exit_code = test(earthly_path, args.timeout)
        if test_exit_code == 2:
            print(f'{test_name} timedout')
            exit_code = test_exit_code
        elif test_exit_code:
            print(f'{test_name} failed with exit code={test_exit_code}')
            exit_code = test_exit_code
        else:
            print(f'{test_name} passed')
    sys.exit(exit_code)
