#!/usr/bin/env bash

set -ue
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

# run a build twice. Validate missing, then present, then no bootstrap. Remove, then bootstrap realsies and validate that next doesnt.

echo "=== Test 1: Hand Bootstrapped ==="

"$earthly" bootstrap

if [[ ! -d "$HOME/.earthly" ]]; then
  echo ".earthly directory was missing after bootstrap"
  exit 1
fi

if [[ ! -f "$HOME/.earthly/install_id" ]]; then
  echo "install_id was missing after bootstrap"
  exit 1
fi

echo "----"
"$earthly" +test | tee hand_boot_output # Hand boots are gloves ;)

if  cat hand_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.earthly/"

echo "=== Test 2: Implied Bootstrap ==="

"$earthly" +test

if [[ ! -d "$HOME/.earthly" ]]; then
  echo ".earthly directory was missing after bootstrap"
  exit 1
fi

if [[ ! -f "$HOME/.earthly/install_id" ]]; then
  echo "install_id was missing after bootstrap"
  exit 1
fi

echo "----"
"$earthly" +test | tee imp_boot_output # Hand boots are gloves ;)

if  cat imp_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.earthly/"