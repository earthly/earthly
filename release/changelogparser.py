#!/usr/bin/env python3
import argparse
import re
import sys
from collections import OrderedDict

class ChangeLogParseError(Exception):
    def __init__(self, message, line):
        super().__init__(message)
        self.line = line

class UnexpectedHeaderError(ChangeLogParseError):
    pass

class MissingTitleError(ChangeLogParseError):
    pass

class MalformedVersionHeaderError(ChangeLogParseError):
    pass

class MalformedHeaderError(ChangeLogParseError):
    pass

class MalformedUnorderedItemError(ChangeLogParseError):
    pass

class DuplicateVersionError(ChangeLogParseError):
    pass


def parse_line(line, line_num):
    '''
    parses lines of the form "# <title>", "## <sub title>", etc.
    if line is not a header, a regular string is returned.
    headers must contain exactly one space between the '#' and title, and may not contain trailing spaces.
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
        raise MalformedHeaderError(line, line_num)
    if line[0] != " ":
        raise MalformedHeaderError(line, line_num)
    line = line[1:]
    if line.startswith(" ") or line.endswith(" "):
        raise MalformedHeaderError(line, line_num)

    return num_headers, line

version_line_re = re.compile(r'^(v[0-9]+\.[0-9]+\.[0-9]+(-rc[0-9]+)?) - ([0-9]{4}-[0-9]{2}-[0-9]{2})( \(aborted release/not recommended\))?$')

def parse_changelog(changelog_data):
    versions = OrderedDict()
    def save_version(version, release_date, body):
        if version in versions:
            raise DuplicateVersionError(version, line_num)
        versions[version] = {
            'date': release_date,
            'body': '\n'.join(body),
        }

    line_num = 1
    version = None
    is_title_body = False
    dash_found = False
    body = []
    ignore = False
    prev_header_num = None
    prev_header_title = None
    for line_num, line in enumerate(changelog_data.splitlines()):
        num_headers, title = parse_line(line, line_num)

        if line_num == 0:
            if num_headers != 1:
                raise MissingTitleError(f'expected title main `# <project-name> Changelog` title; got {line}', line_num)
            if not title.endswith(' Changelog'):
                raise MissingTitleError("expected title ending with Changelog", line_num)
            is_title_body = True
            prev_header_num = num_headers
            prev_header_title = title
            continue

        if num_headers == 0:
            if line == '<!--changelog-parser-ignore-start-->':
                ignore = True
                continue
            if line == '<!--changelog-parser-ignore-end-->':
                ignore = False
                continue
            if ignore:
                pass
            elif is_title_body:
                pass # no linting of title body
            elif is_intro:
                pass # no linting of intro text
            elif line == '':
                dash_found = False
            elif line.startswith('-'):
                if not line.startswith('- '):
                    raise MalformedUnorderedItemError(f'expected unordered item of the form `- <text>`; got {line}', line_num)
                dash_found = True
            elif not line.startswith(' '):
                raise MalformedUnorderedItemError(f'expected unordered item of the form `- <text>` (or `- <text>\\n  <more text>`); got {line}', line_num)
            elif line.startswith(' ') and dash_found is False:
                raise MalformedUnorderedItemError(f'expected unordered item of the form `- <text>` (or `- <text>\\n  <more text>`); got {line}', line_num)
            body.append(line)
        elif num_headers == 1:
            raise UnexpectedHeaderError(line, line_num)
        elif num_headers == 2:
            is_intro = True
            ignore = False
            if is_title_body:
                if title != 'Unreleased':
                    raise MissingTitleError(f'expected `## Unreleased` title; got {line}', line_num)
                is_title_body = False
                assert version is None
                version = title
                release_date = None
            else:
                if version:
                    save_version(version, release_date, body)
                m = version_line_re.match(title)
                if not m:
                    raise MalformedVersionHeaderError(line, line_num)
                version = m.group(1)
                release_date = m.group(3)
            body = []
        elif num_headers == 3:
            ignore = False
            is_intro = False
            allowed_titles = ('Added', 'Changed', 'Removed', 'Fixed')
            if title not in allowed_titles:
                raise UnexpectedHeaderError(f'expected header of {allowed_titles}; but got "{title}"', line_num)
            if prev_header_num not in (2, 3):
                raise UnexpectedHeaderError(f'expected header "{title}" to be under a "vX.Y.Z" or "Unreleased" section, but instead it was located after "{prev_header_title}"', line_num)

            body.append(line)
        else:
            raise UnexpectedHeaderError(f'unsupported header {line}')

        if num_headers > 0:
            prev_header_num = num_headers
            prev_header_title = title

    if version:
        save_version(version, release_date, body)

    return versions

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--changelog', help='path to change log; if omitted changelog is read from stdin')
    parser.add_argument('--version', help='version to display; if omitted, changelog is still parsed and any errors are displayed', default=None)
    parser.add_argument('--date', help='display release date of specified version', action='store_true')
    args = parser.parse_args()

    path_str = args.changelog
    if path_str is None:
        path_str = 'stdin'
        changelog_str = sys.stdin.read()
    else:
        with open(args.changelog, 'rb') as fp:
            changelog_str = fp.read().decode('utf8')

    try:
        changelog = parse_changelog(changelog_str)
    except MalformedVersionHeaderError as e:
        print(f'failed to parse {path_str}:{e.line+1}: unable to parse "{e}"; should be of the form "v1.2.3 - YYYY-MM-DD" (or "v1.2.3-rc4 - YYYY-MM-DD")', file=sys.stderr)
        sys.exit(1)
    except MalformedHeaderError as e:
        print(f'failed to parse {path_str}:{e.line+1}: malformed header found ({e}); should be "#[#[...]] <title>"', file=sys.stderr)
        sys.exit(1)
    except DuplicateVersionError as e:
        print(f'failed to parse {path_str}:{e.line+1}: duplicate titles ({e}) detected', file=sys.stderr)
        sys.exit(1)
    except ChangeLogParseError as e:
        print(f'failed to parse {path_str}:{e.line+1}: unhandled exception {e.__class__.__name__}: {e}', file=sys.stderr)
        sys.exit(1)

    if args.version is None:
        # running under linting mode
        sys.exit(0)

    try:
        details = changelog[args.version]
    except KeyError:
        print(f'No changelog entry exists for {args.version}', file=sys.stderr)
        sys.exit(1)
    if args.date:
        print(details['date'])
    else:
        print(details['body'].strip())
