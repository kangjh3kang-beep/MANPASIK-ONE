#!/bin/bash
TARGET_DIR="/mnt/d/우리집/flutter_cache/flutter"
echo "Fixing CRLF in $TARGET_DIR recursively..."

find "$TARGET_DIR" -type f -name "*.sh" -exec sed -i 's/\r$//' {} +
find "$TARGET_DIR/bin" -type f -name "flutter" -exec sed -i 's/\r$//' {} +

echo "Triggering Flutter SDK download..."
# Ensure we are using the flutter binary we just fixed
"$TARGET_DIR/bin/flutter" --version
