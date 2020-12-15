#!/bin/bash
set -xeuo pipefail

earthly="earthly"
if ! command -v "$earthly"; then
    earthly="earth"
fi

$earthly --version

$earthly github.com/earthly/earthly-example-scala/simple:main+test
$earthly github.com/earthly/earthly-example-scala/simple:main+docker
