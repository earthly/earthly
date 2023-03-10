#!/bin/bash
set -e

export CHECK_TAG_WAS_PUSHED=true

cd "$(dirname "$0")"
../common/test.sh --push
