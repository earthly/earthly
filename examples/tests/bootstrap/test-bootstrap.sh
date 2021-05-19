#!/usr/bin/env bash

set -ue
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../../build/linux/amd64/earthly"}

echo "=== Test 1: Hand Bootstrapped ==="

"$earthly" bootstrap

if [[ ! -d "$HOME/.earthly" ]]; then
  echo ".earthly directory was missing after bootstrap"
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

echo "----"
"$earthly" +test | tee imp_boot_output

if  cat imp_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.earthly/"

echo "=== Test 3: CI ==="

"$earthly" --ci +test

if [[ ! -d "$HOME/.earthly" ]]; then
  echo ".earthly directory was missing after bootstrap"
  exit 1
fi

echo "----"
"$earthly" --ci +test | tee ci_boot_output

if  cat ci_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.earthly/"

echo "=== Test 4: With Autocomplete ==="

"$earthly" bootstrap

if [[ -f "/usr/share/bash-completion/completions/earthly" ]]; then
  echo "autocompletions were present when they should not have been"
  exit 1
fi

echo "----"
sudo "$earthly" bootstrap --with-autocomplete

if [[ ! -f "/usr/share/bash-completion/completions/earthly" ]]; then
  echo "autocompletions were missing when they should have been present"
  exit 1
fi

rm -rf "$HOME/.earthly/"
sudo rm -rf "/usr/share/bash-completion/completions/earthly"

echo "=== Test 5: Permissions ==="

"$earthly" bootstrap

if [[ $(stat --format '%U' "$HOME/.earthly") != "$USER" ]]; then
  echo "earthly directory is not owned by the user"
  stat "$HOME/.earthly"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.earthly") != "$USER" ]]; then
  echo "earthly directory is not owned by the users group"
  stat "$HOME/.earthly"
  exit 1
fi

echo "----"
touch $HOME/.earthly/config.yml
sudo chown -R 12345:12345 $HOME/.earthly

sudo "$earthly" bootstrap

if [[ $(stat --format '%U' "$HOME/.earthly") != "$USER" ]]; then
  echo "earthly directory is not owned by the user"
  stat "$HOME/.earthly"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.earthly") != "$USER" ]]; then
  echo "earthly directory is not owned by the users group"
  stat "$HOME/.earthly"
  exit 1
fi

if [[ $(stat --format '%U' "$HOME/.earthly/config.yml") != "$USER" ]]; then
  echo "earthly config is not owned by the user"
  stat "$HOME/.earthly/config.yml"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.earthly/config.yml") != "$USER" ]]; then
  echo "earthly config is not owned by the users group"
  stat "$HOME/.earthly/config.yml"
  exit 1
fi