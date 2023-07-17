#!/usr/bin/env bash
set -e

if ! which curl 2>/dev/null >&2; then
  echo >&2 "curl is required; please install it."
  exit 1
fi

cd "$(dirname "$0")"
curl https://raw.githubusercontent.com/alexcb/acbtest/main/acbgrep > acbgrep && chmod +x acbgrep
curl https://raw.githubusercontent.com/alexcb/acbtest/main/acbtest > acbtest && chmod +x acbtest
