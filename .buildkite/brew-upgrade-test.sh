#!/bin/bash

set -xeuo pipefail

brew upgrade earthly
earth bootstrap

earth --version

earth github.com/earthly/earthly/examples/go+docker
