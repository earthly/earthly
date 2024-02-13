#!/usr/bin/env bash

source "./assert.sh"
set -v
results=$(docker run --network=host earthly/examples:integration)
expected="The first 5 countries alphabetically are: Afghanistan, Albania, Algeria, American Samoa, Andorra"

assert_eq "$expected" "$results"
