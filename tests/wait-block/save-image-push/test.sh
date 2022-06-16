#!/bin/bash
set -e

cd "$(dirname "$0")"
../common/test.sh --push
