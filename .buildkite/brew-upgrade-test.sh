#!/bin/bash
set -xeuo pipefail

earthly="earthly"
if ! command -v "$earthly"; then
    earthly="earth"
fi

brew upgrade earthly

$earthly --version

$earthly github.com/earthly/earthly/examples/go:main+docker
