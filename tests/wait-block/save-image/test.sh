#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "first running without the --push flag"
../common/test.sh

echo "re-running test with --push flag"
../common/test.sh --push
