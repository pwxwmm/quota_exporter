#!/bin/bash
set -e

# Get the current script directory
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(dirname "$SCRIPT_DIR")

# Navigate to the project root directory
cd "$PROJECT_ROOT"

# Start the quota_exporter process in the background
echo "[INFO] Starting quota_exporter..."
nohup bin/quota_exporter > logs/quota_exporter.log 2>&1 &
echo "[INFO] quota_exporter started! Logs: logs/quota_exporter.log"
