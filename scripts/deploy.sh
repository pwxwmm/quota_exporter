#!/bin/bash
set -e

# Define the target server and deployment directory
TARGET_SERVER="user@remote-server"
TARGET_DIR="/opt/quota_exporter"

# Build the project before deployment
./scripts/build.sh

# Create necessary directories on the remote server
ssh "$TARGET_SERVER" "mkdir -p $TARGET_DIR/bin $TARGET_DIR/config $TARGET_DIR/logs"

# Transfer the necessary files to the remote server
scp bin/quota_exporter "$TARGET_SERVER:$TARGET_DIR/bin/"
scp config/config.yaml "$TARGET_SERVER:$TARGET_DIR/config/"
scp scripts/start.sh "$TARGET_SERVER:$TARGET_DIR/"

# Start the exporter on the remote server
echo "[INFO] Deploying and starting quota_exporter on $TARGET_SERVER..."
ssh "$TARGET_SERVER" "bash $TARGET_DIR/start.sh"
echo "[INFO] Deployment completed!"
