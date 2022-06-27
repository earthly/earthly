#!/bin/bash
set -ex

cd "$(dirname "$0")"

# run each test.sh script from each sub directory
for d in *; do
    if [ -d "$d" ] && [ "$d" != "common" ]; then
        echo "=== Running $d/test.sh ==="
        "$d/test.sh"
    fi
done
