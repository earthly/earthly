#!/usr/bin/env python3
import argparse
import re
import sys
from collections import OrderedDict

class UnexpectedHeaderError(Exception):
    pass

class MissingTitleError(Exception):
    pass

class MalformedVersionHeaderError(Exception):
    pass

class MalformedHeaderError(Exception):
    pass

class DuplicateVersionError(KeyError):
    pass


def parse_line(line):
    '''
    parses lines of the form "# <title>", "## <sub title>", etc.
    if line is not a header, a regular string is returned.
    headers must contain exactly once space between the '#' and title, and may not contain trailling space.
    tabs are not friends.
    '''
    num_headers = 0
    for c in line:
        if c == '#':
            num_headers += 1
        else:
            break
    if num_headers == 0:
        return 0, line
    line = line[num_headers:]
    if line == "":
        raise MalformedHeaderError(line)
    if line[0] != " ":
        raise MalformedHeaderError(line)
    line = line[1:]
    if line.startswith(" ") or line.endswith(" "):
        raise MalformedHeaderError(line)

    return num_headers, line

version_line_re = re.compile(r'^(v[0-9]+\.[0-9]+\.[0-9]+(-rc[0-9]+)?) - ([0-9]{4}-[0-9]{2}-[0-9]{2})$')

def parse_changelog(changelog_data):
    lines = changelog_data.splitlines()
    num_headers, title = parse_line(lines[0])
    if num_headers != 1:
        raise MissingTitleError
    if not title.endswith(' Changelog'):
        raise MissingTitleError

    versions = OrderedDict()
    def save_version(version, release_date, body):
        if version in versions:
            raise DuplicateVersionError(version)
        versions[version] = {
            'date': release_date,
            'body': '\n'.join(body),
        }

    lines = lines[1:]
    version = None
    for line in lines:
        num_headers, line = parse_line(line)
        if num_headers == 1:
            raise UnexpectedHeaderError(line)
        elif num_headers == 2:
            if version:
                save_version(version, release_date, body)
            if line == 'Unreleased':
                version = line
                release_date = None
            else:
                m = version_line_re.match(line)
                if not m:
                    raise MalformedVersionHeaderError(line)
                version = m.group(1)
                release_date = m.group(2)
            body = []
        elif version:
            body.append(line)

    if version:
        save_version(version, release_date, body)

    return versions

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('changelog', help='path to change log')
    parser.add_argument('version', help='version to display')
    args = parser.parse_args()

    try:
        with open(args.changelog, 'rb') as fp:
            changelog = parse_changelog(fp.read().decode('utf8'))
    except MalformedVersionHeaderError as e:
        print(f'failed to parse changelog: unable to parse "{e}"; should be of the form "v1.2.3 - YYYY-MM-DD"', file=sys.stderr)
        sys.exit(1)
    except MalformedHeaderError as e:
        print(f'failed to parse changelog: malformed header found ({e}); should be "#[#[...]] <title>"', file=sys.stderr)
        sys.exit(1)
    except DuplicateVersionError as e:
        print(f'failed to parse changelog: duplicate titles ({e}) detected', file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f'failed to parse changelog: unhandled exception {e.__class__.__name__}: {e}', file=sys.stderr)
        sys.exit(1)

    try:
        details = changelog[args.version]
    except KeyError:
        print('No changelog entry exists for {args.version}', file=sys.stderr)
        sys.exit(1)
    print(details['body'].strip())
