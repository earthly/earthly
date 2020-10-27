#!/bin/bash
set -xeuo pipefail

export EARTHLY_INSTALL_ID=buildkite-earthly-mac

earth --version

earth github.com/earthly/earthly-example-scala/simple+test
earth github.com/earthly/earthly-example-scala/simple+docker
