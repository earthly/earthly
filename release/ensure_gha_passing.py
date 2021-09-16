#!/usr/bin/env python3
import argparse
import sys

import requests

def get_checks(org, repo, sha):
    url = f'https://api.github.com/repos/{org}/{repo}/commits/{sha}/check-runs'
    print(f'querying {url} for check statuses')
    r = requests.get(url)
    return r.json()

def display_checks(checks):
    num_non_completed = 0
    check_runs = checks['check_runs']
    for check_run in sorted(check_runs, key=lambda x: x['name']):
        name = check_run['name']
        status = check_run['status']
        details_url = check_run['details_url']
        print(f'{name}: {status}')
        if status != 'completed':
            print(f'  details: {details_url}')
            num_non_completed += 1
    return num_non_completed, len(check_runs)

def get_argparser():
    parser = argparse.ArgumentParser(description='check statuses')
    parser.add_argument('--org', default='earthly', help='github org')
    parser.add_argument('--repo', default='earthly', help='github repo')
    parser.add_argument('--sha', required=True, help='commit sha to check')
    return parser

if __name__ == '__main__':
    args = get_argparser().parse_args()
    num_non_completed, num_total = display_checks(get_checks(args.org, args.repo, args.sha))
    if num_total == 0:
        print(f'error: no checks were detected')
        sys.exit(2)
    if num_non_completed:
        print(f'error: {num_non_completed} check(s) reported non-completed status')
        sys.exit(1)
    print(f'success: {num_total} check(s) reported completed status')
