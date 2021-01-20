#!/bin/bash
set -xeuo pipefail

earthly="earthly"
if ! command -v "$earthly"; then
    earthly="earth"
fi

$earthly --version

$earthly examples/tests/cloud-push-pull+all
