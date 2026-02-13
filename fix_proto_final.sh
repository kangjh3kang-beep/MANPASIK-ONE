#!/bin/bash
# fix_proto_final.sh - Patch manpasik.pb.dart to handle pc args and types

PROTO_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart"

if [ ! -f "$PROTO_FILE_DART" ]; then
    echo "Error: $PROTO_FILE_DART not found"
    exit 1
fi

echo "Patching $PROTO_FILE_DART (final fixes)..."

# Fix pc arguments: remove PbFieldType.PM
sed -i 's/\$pb.PbFieldType.PM, //g' "$PROTO_FILE_DART"

# Fix create -> CreateMessage
# We only replace .create) or .create, to avoid replacing other words ending in create
sed -i 's/\.create)/\.CreateMessage)/g' "$PROTO_FILE_DART"
sed -i 's/\.create,/\.CreateMessage,/g' "$PROTO_FILE_DART"

# Fix remaining types in generics
# Handle List<int>, List<String>, etc within $core.List<...>
# Previously we tried replacing List<int>, but it might be $core.List<int>
sed -i 's/<int>/<\$core.int>/g' "$PROTO_FILE_DART"
sed -i 's/<String>/<\$core.String>/g' "$PROTO_FILE_DART"
sed -i 's/<double>/<\$core.double>/g' "$PROTO_FILE_DART"
sed -i 's/<bool>/<\$core.bool>/g' "$PROTO_FILE_DART"

echo "Patch complete."
