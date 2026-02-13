#!/bin/bash
# fix_proto_type_args.sh - Patch manpasik.pb.dart to fix generic type arguments

PROTO_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart"

if [ ! -f "$PROTO_FILE_DART" ]; then
    echo "Error: $PROTO_FILE_DART not found"
    exit 1
fi

echo "Patching $PROTO_FILE_DART (type args)..."

# Use single quotes for sed to avoid shell expansion
sed -i 's/<int>/<$core.int>/g' "$PROTO_FILE_DART"
sed -i 's/<String>/<$core.String>/g' "$PROTO_FILE_DART"
sed -i 's/<double>/<$core.double>/g' "$PROTO_FILE_DART"
sed -i 's/<bool>/<$core.bool>/g' "$PROTO_FILE_DART"

echo "Patch complete."
