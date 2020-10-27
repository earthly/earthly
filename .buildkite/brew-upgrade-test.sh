#!/bin/bash
set -xeuo pipefail

export EARTHLY_INSTALL_ID=buildkite-earthly-mac

brew upgrade earthly

earth --version

earth github.com/earthly/earthly/examples/go+docker
