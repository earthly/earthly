#!/bin/bash

set -xeuo pipefail

earth --version

earth github.com/earthly/earthly-example-scala/simple:main+test
earth github.com/earthly/earthly-example-scala/simple:main+docker
