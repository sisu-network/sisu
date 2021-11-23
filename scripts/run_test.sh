#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

cd "$PROJECT_DIR" || exit 1

echo_info "Test all packages"
go test -race $(go list ./...)

EXIT_CODE=$?
exit $EXIT_CODE
