#!/usr/bin/env bash

#####################################################################
##
## title: Assert Extension
##
## description:
## Assert extension of shell (bash, ...)
##   with the common assert functions
## Function list based on:
##   http://junit.sourceforge.net/javadoc/org/junit/Assert.html
## Log methods : inspired by
##	- https://natelandau.com/bash-scripting-utilities/
## author: Mark Torok
##
## date: 07. Dec. 2016
##
## license: MIT
##
#####################################################################

if command -v tput &>/dev/null && tty -s; then
  RED=$(tput setaf 1)
  GREEN=$(tput setaf 2)
  MAGENTA=$(tput setaf 5)
  NORMAL=$(tput sgr0)
  BOLD=$(tput bold)
else
  RED=$(echo -en "\e[31m")
  GREEN=$(echo -en "\e[32m")
  MAGENTA=$(echo -en "\e[35m")
  NORMAL=$(echo -en "\e[00m")
  BOLD=$(echo -en "\e[01m")
fi

log_header() {
  printf "\n${BOLD}${MAGENTA}==========  %s  ==========${NORMAL}\n" "$@" >&2
}

log_success() {
  printf "${GREEN}✔ %s${NORMAL}\n" "$@" >&2
}

log_failure() {
  printf "${RED}✖ %s${NORMAL}\n" "$@" >&2
}


assert_eq() {
  local expected="$1"
  local actual="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [ "$expected" == "$actual" ]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$expected == $actual :: $msg" || true
    return 1
  fi
}

assert_not_eq() {
  local expected="$1"
  local actual="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [ ! "$expected" == "$actual" ]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$expected != $actual :: $msg" || true
    return 1
  fi
}

assert_true() {
  local actual
  local msg

  actual="$1"

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  assert_eq true "$actual" "$msg"
  return "$?"
}

assert_false() {
  local actual
  local msg

  actual="$1"

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  assert_eq false "$actual" "$msg"
  return "$?"
}

assert_array_eq() {

  declare -a expected=("${!1}")
  # echo "AAE ${expected[@]}"

  declare -a actual=("${!2}")
  # echo "AAE ${actual[@]}"

  local msg
  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  local return_code
  return_code=0
  if [ ! "${#expected[@]}" == "${#actual[@]}" ]; then
    return_code=1
  fi

  local i
  for (( i=1; i < ${#expected[@]} + 1; i+=1 )); do
    if [ ! "${expected[$i-1]}" == "${actual[$i-1]}" ]; then
      return_code=1
      break
    fi
  done

  if [ "$return_code" == 1 ]; then
    [ "${#msg}" -gt 0 ] && log_failure "(${expected[*]}) != (${actual[*]}) :: $msg" || true
  fi

  return "$return_code"
}

assert_array_not_eq() {

  declare -a expected=("${!1}")
  declare -a actual=("${!2}")

  local msg
  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  local return_code
  return_code=1
  if [ ! "${#expected[@]}" == "${#actual[@]}" ]; then
    return_code=0
  fi

  local i
  for (( i=1; i < ${#expected[@]} + 1; i+=1 )); do
    if [ ! "${expected[$i-1]}" == "${actual[$i-1]}" ]; then
      return_code=0
      break
    fi
  done

  if [ "$return_code" == 1 ]; then
    [ "${#msg}" -gt 0 ] && log_failure "(${expected[*]}) == (${actual[*]}) :: $msg" || true
  fi

  return "$return_code"
}

assert_empty() {
  local actual
  local msg

  actual="$1"

  if [ "$#" -ge 2 ]; then
    msg="$2"
  fi

  assert_eq "" "$actual" "$msg"
  return "$?"
}

assert_not_empty() {
  local actual
  local msg

  actual="$1"

  if [ "$#" -ge 2 ]; then
    msg="$2"
  fi

  assert_not_eq "" "$actual" "$msg"
  return "$?"
}

assert_contain() {
  local haystack="$1"
  local needle="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [ -z "${needle:+x}" ]; then
    return 0;
  fi

  if [ -z "${haystack##*$needle*}" ]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$haystack doesn't contain $needle :: $msg" || true
    return 1
  fi
}

assert_not_contain() {
  local haystack="$1"
  local needle="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [ -z "${needle:+x}" ]; then
    return 0;
  fi

  if [ "${haystack##*$needle*}" ]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$haystack contains $needle :: $msg" || true
    return 1
  fi
}

assert_gt() {
  local first="$1"
  local second="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [[ "$first" -gt  "$second" ]]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$first > $second :: $msg" || true
    return 1
  fi
}

assert_ge() {
  local first="$1"
  local second="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [[ "$first" -ge  "$second" ]]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$first >= $second :: $msg" || true
    return 1
  fi
}

assert_lt() {
  local first="$1"
  local second="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [[ "$first" -lt  "$second" ]]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$first < $second :: $msg" || true
    return 1
  fi
}

assert_le() {
  local first="$1"
  local second="$2"
  local msg

  if [ "$#" -ge 3 ]; then
    msg="$3"
  fi

  if [[ "$first" -le  "$second" ]]; then
    return 0
  else
    [ "${#msg}" -gt 0 ] && log_failure "$first <= $second :: $msg" || true
    return 1
  fi
}
