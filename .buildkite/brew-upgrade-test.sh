#!/bin/bash

set -xeuo pipefail

brew upgrade earth
earth bootstrap

earth --version

earth github.com/earthly/earthly/examples/go+docker
