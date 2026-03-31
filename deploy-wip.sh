#!/bin/bash
# Deploy local changes to the GeekBudget.WIP Portainer stack.
# Works with uncommitted and unpushed changes. No git push required.
#
# Usage: ./deploy-wip.sh [--dry-run]

set -e
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
python3 /data/portainer.py deploy "$SCRIPT_DIR/deploy-wip.json" "$@"
