#!/bin/bash

set -xeuo pipefail

earthly --version

earthly github.com/earthly/earthly-example-scala/simple:main+test
earthly github.com/earthly/earthly-example-scala/simple:main+docker
