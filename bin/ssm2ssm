#!/usr/bin/env bash

function _run() {
  local script=''
  script="$(basename "${BASH_SOURCE[0]}")"
  local directory=''
  directory="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
  local executable="$directory/bin/$script.build"
  local log_directory="$directory/bin/.log"
  local stdout="$log_directory/$script.stdout.log"
  local stderr="$log_directory/$script.stderr.log"

  mkdir -p "$log_directory"

  go build \
    -o "$executable" \
    "$directory/cmd/$script" \
    >"$stdout" \
    2>"$stderr"

  status=$?
  [ $status -eq 0 ] || {
    echo "[ERROR] Error while building $script"
    cat "$stderr"
    exit 255
  } >&2

  set -e

  "$executable" "${@:-}"
}

if [ "${BASH_SOURCE[0]}" = "$0" ]; then
  _run "${@:-}"
  exit $?
fi
