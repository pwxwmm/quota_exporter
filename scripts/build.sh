#!/bin/bash
set -e

# Get the current script directory
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PROJECT_ROOT=$(dirname "$SCRIPT_DIR")

# Navigate to the project root directory
cd "$PROJECT_ROOT"

echo "[INFO] Building quota_exporter..."

# Compile the Go project and place the binary in the bin directory
go build -o bin/quota_exporter ./cmd/main.go

echo "[INFO] Build completed! Binary is located at bin/quota_exporter"