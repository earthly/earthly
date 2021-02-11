#!/bin/bash
set -xeuo pipefail

cpubrand=$(sysctl -n machdep.cpu.brand_string)
echo "macOS test running on $cpubrand"

earthly="earthly"
if ! command -v "$earthly"; then
    earthly="earth"
fi

brew upgrade earthly

$earthly --version

$earthly github.com/earthly/earthly/examples/go:main+docker
