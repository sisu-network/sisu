#!/usr/bin/env bash
# This scripts provide some helpful method and predefined shell variables. This
# is not intended to be used by its own. Use should source this script into
# other scripts.

SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPTS_DIR/.." && pwd)"

echo_info() {
  echo "[INFO] $*"
}
