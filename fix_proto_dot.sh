#!/bin/bash
# fix_proto_dot.sh - Patch manpasik.pb.dart to remove leading dot before $core

PROTO_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart"

if [ ! -f "$PROTO_FILE_DART" ]; then
    echo "Error: $PROTO_FILE_DART not found"
    exit 1
fi

echo "Patching $PROTO_FILE_DART (removing leading dot)..."

# Replace .$core with $core
sed -i 's/\.\$core/\$core/g' "$PROTO_FILE_DART"

echo "Patch complete."
