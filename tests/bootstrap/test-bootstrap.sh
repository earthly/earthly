#!/usr/bin/env bash

set -uex
set -o pipefail

cd "$(dirname "$0")"

earthly=${earthly-"../../build/linux/amd64/earthly"}
earthly=$(realpath "$earthly")

# docker / podman
frontend="${frontend:-$(which docker || which podman)}"
test -n "$frontend" || (>&2 echo "Error: frontend is empty" && exit 1)

# so that we can test production/staging binaries
default_install_name="${default_install_name:-"earthly-dev"}"

echo "=== Test 1: Hand Bootstrapped ==="

"$earthly" bootstrap

if [[ ! -d "$HOME/.${default_install_name}" ]]; then
  echo ".${default_install_name} directory was missing after bootstrap"
  exit 1
fi

EARTHLY_INSTALLATION_NAME=earthly-test "$earthly" bootstrap

if [[ ! -d "$HOME/.earthly-test" ]]; then
  echo ".earthly-test directory was missing after bootstrap"
  exit 1
fi

echo "----"
"$earthly" +test | tee hand_boot_output # Hand boots are gloves ;)

if  cat hand_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.${default_install_name}" "$HOME/.earthly-test"

echo "=== Test 2: Implied Bootstrap ==="

"$earthly" +test

if [[ ! -d "$HOME/.${default_install_name}" ]]; then
  echo ".${default_install_name} directory was missing after bootstrap"
  exit 1
fi

EARTHLY_INSTALLATION_NAME=earthly-test "$earthly" +test

if [[ ! -d "$HOME/.earthly-test" ]]; then
  echo ".earthly-test directory was missing after bootstrap"
  exit 1
fi

echo "----"
"$earthly" +test | tee imp_boot_output

if  cat imp_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.${default_install_name}" "$HOME/.earthly-test"

echo "=== Test 3: CI ==="

"$earthly" --ci +test

if [[ ! -d "$HOME/.${default_install_name}" ]]; then
  echo ".${default_install_name} directory was missing after bootstrap"
  exit 1
fi

EARTHLY_INSTALLATION_NAME=earthly-test "$earthly" --ci +test

if [[ ! -d "$HOME/.earthly-test" ]]; then
 echo ".earthly-test directory was missing after bootstrap"
 exit 1
fi

echo "----"
"$earthly" --ci +test | tee ci_boot_output

if  cat ci_boot_output | grep -q "bootstrap |"; then
    echo "build did extra bootstrap"
    exit 1
fi

rm -rf "$HOME/.${default_install_name}" "$HOME/.earthly-test"

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

rm -rf "$HOME/.${default_install_name}"
sudo rm -rf "/usr/share/bash-completion/completions/earthly"

echo "=== Test 5: Permissions ==="

touch testfile
USR=$(stat --format '%U' testfile)
GRP=$(stat --format '%G' testfile)

echo "Current defaults:"
echo "User : $USR"
echo "Group: $GRP"

"$earthly" bootstrap

if [[ $(stat --format '%U' "$HOME/.${default_install_name}") != "$USR" ]]; then
  echo "earthly directory is not owned by the user"
  stat "$HOME/.${default_install_name}"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.${default_install_name}") != "$GRP" ]]; then
  echo "earthly directory is not owned by the users group"
  stat "$HOME/.${default_install_name}"
  exit 1
fi

echo "----"

touch $HOME/.${default_install_name}/config.yml
sudo chown -R 12345:12345 $HOME/.${default_install_name}

sudo "$earthly" bootstrap

if [[ $(stat --format '%U' "$HOME/.${default_install_name}") != "$USR" ]]; then
  echo "earthly directory is not owned by the user"
  stat "$HOME/.${default_install_name}"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.${default_install_name}") != "$GRP" ]]; then
  echo "earthly directory is not owned by the users group"
  stat "$HOME/.${default_install_name}"
  exit 1
fi

if [[ $(stat --format '%U' "$HOME/.${default_install_name}/config.yml") != "$USR" ]]; then
  echo "earthly config is not owned by the user"
  stat "$HOME/.${default_install_name}/config.yml"
  exit 1
fi

if [[ $(stat --format '%G' "$HOME/.${default_install_name}/config.yml") != "$GRP" ]]; then
  echo "earthly config is not owned by the users group"
  stat "$HOME/.${default_install_name}/config.yml"
  exit 1
fi

echo "=== Test 6: works in read-only directory ==="

sudo mkdir /tmp/earthly-read-only-test
sudo cp Earthfile /tmp/earthly-read-only-test/.
sudo chmod 0755 /tmp/earthly-read-only-test/.

prevdir=$(pwd)
cd /tmp/earthly-read-only-test/.

if touch this-should-fail 2>/dev/null; then
  echo "this directory should have been read-only; something is wrong with this test"
  exit 1
fi

"$earthly" +test

cd "$prevdir"

echo "=== Test 7: Homebrew Source ==="

if which "$frontend" > /dev/null; then
  "$frontend" rm -f earthly-buildkitd
fi

bash=$("$earthly" bootstrap --source bash)
if [[ "$bash" != *"complete -o nospace"* ]]; then
  echo "bash autocompletion appeared to be incorrect"
  echo "$bash"
  exit 1
fi

zsh=$("$earthly" bootstrap --source zsh)
if [[ "$zsh" != *"complete -o nospace"* ]]; then
  echo "zsh autocompletion appeared to be incorrect"
  echo "$zsh"
  exit 1
fi

if "$frontend" container ls | grep earthly-buildkitd; then
  echo "--source created a $frontend container"
  exit 1
fi

if [[ -f ../../build/linux/amd64/earth ]]; then
  echo "--source symlinked earthly to earth"
fi

if ! DOCKER_HOST="$frontend is missing" "$earthly" bootstrap --source zsh > /dev/null 2>&1; then
  echo "--source failed when $frontend was missing"
  exit 1
fi

rm -rf "$HOME/.${default_install_name}"

echo "=== Test 8: No Buildkit ==="

"$earthly" bootstrap --no-buildkit
if "$frontend" container ls | grep earthly-buildkitd; then
  echo "--no-buildkit created a $frontend container"
  exit 1
fi

if ! DOCKER_HOST="$frontend is missing" "$earthly" bootstrap --no-buildkit; then
  echo "--no-buildkit fails when $frontend is missing"
  exit 1
fi

rm -rf "$HOME/.${default_install_name}"
