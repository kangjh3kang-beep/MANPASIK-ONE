#!/bin/bash
echo "Killing locked processes..."
/mnt/c/Windows/System32/taskkill.exe /F /IM dart.exe /T 2>/dev/null
/mnt/c/Windows/System32/taskkill.exe /F /IM flutter.exe /T 2>/dev/null

echo "Removing cache..."
rm -rf /mnt/d/우리집/flutter_cache/flutter/bin/cache

echo "Triggering Flutter SDK download..."
flutter --version
