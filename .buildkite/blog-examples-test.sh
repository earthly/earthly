#!/bin/bash

set -xeuo pipefail

earth --version

earth github.com/earthly/earthly-example-scala:main/simple+test
earth github.com/earthly/earthly-example-scala:main/simple+docker
