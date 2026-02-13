#!/bin/bash
# fix_proto.sh - Patch manpasik.pb.dart to use correct protobuf types

PROTO_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart"
PBGRPC_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart"

if [ ! -f "$PROTO_FILE_DART" ]; then
    echo "Error: $PROTO_FILE_DART not found"
    exit 1
fi

echo "Patching $PROTO_FILE_DART..."

# Replace basic types with $core prefix (handle nullable ? and space)
# Use loop to handle multiple occurrences if needed or just specific patterns
sed -i 's/String? /\$core.String? /g' "$PROTO_FILE_DART"
sed -i 's/int? /\$core.int? /g' "$PROTO_FILE_DART"
sed -i 's/double? /\$core.double? /g' "$PROTO_FILE_DART"
sed -i 's/bool? /\$core.bool? /g' "$PROTO_FILE_DART"

# Also handle non-nullable return types or parameters if missed previously
# (Note: previously ran script handled 'String ' but missed 'String?')
sed -i 's/String /\$core.String /g' "$PROTO_FILE_DART"
sed -i 's/int /\$core.int /g' "$PROTO_FILE_DART"
sed -i 's/double /\$core.double /g' "$PROTO_FILE_DART"
sed -i 's/bool /\$core.bool /g' "$PROTO_FILE_DART"

# Handle List<int> -> List<$core.int>
sed -i 's/List<int>/List<\$core.int>/g' "$PROTO_FILE_DART"
sed -i 's/List<String>/List<\$core.String>/g' "$PROTO_FILE_DART"
sed -i 's/List<double>/List<\$core.double>/g' "$PROTO_FILE_DART"
sed -i 's/List<bool>/List<\$core.bool>/g' "$PROTO_FILE_DART"

# Fix List itself if not already handled
sed -i 's/List</\$core.List</g' "$PROTO_FILE_DART"
sed -i 's/Map</\$core.Map</g' "$PROTO_FILE_DART"

# Remove potential double prefixes if run multiple times (e.g. $core.$core.String)
sed -i 's/\$core\.\$core/\$core/g' "$PROTO_FILE_DART"

# Fix method names
sed -i 's/\$_getND/\$_getN/g' "$PROTO_FILE_DART"

echo "Patch complete."
