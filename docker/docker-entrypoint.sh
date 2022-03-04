#!/bin/sh

set -e

if [ -n "$1" ]; then
  hash "$1" 2>/dev/null || {
    echo "[DEBUG] Passing all arguments to ssm2ssm" >&2
    ssm2ssm "${@:-}"
    exit $?
  }
fi

echo "[DEBUG] Running command as is"
exec "$@"