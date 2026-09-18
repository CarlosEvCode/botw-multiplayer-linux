#!/usr/bin/env bash
# Backward-compatibility wrapper for install.sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$SCRIPT_DIR/install.sh" "$@"
