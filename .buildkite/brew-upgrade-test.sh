#!/bin/bash

set -xeuo pipefail

brew upgrade earthly

earthly --version

earthly github.com/earthly/earthly/examples/go:main+docker
