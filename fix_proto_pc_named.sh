#!/bin/bash
# fix_proto_pc_named.sh - Patch manpasik.pb.dart to use named subBuilder for pc

PROTO_FILE_DART="$HOME/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart"

if [ ! -f "$PROTO_FILE_DART" ]; then
    echo "Error: $PROTO_FILE_DART not found"
    exit 1
fi

echo "Patching $PROTO_FILE_DART (pc named args)..."

# Case 1: pc<DeviceInfo>(1, 'devices', DeviceInfo.CreateMessage)
# We want: pc<DeviceInfo>(1, 'devices', $pb.PbFieldType.PM, subBuilder: DeviceInfo.CreateMessage)

# Use sed with capture groups or just specific patterns
# Since we already removed PbFieldType.PM, current state is:
# ..pc<DeviceInfo>(1, 'devices', DeviceInfo.CreateMessage);

sed -i "s/\.\.pc<DeviceInfo>(1, 'devices', DeviceInfo\.CreateMessage)/\.\.pc<DeviceInfo>(1, 'devices', \$pb.PbFieldType.PM, subBuilder: DeviceInfo.CreateMessage)/g" "$PROTO_FILE_DART"

# Case 2: ..pc<MeasurementSummary>(1, 'measurements', MeasurementSummary.CreateMessage)

sed -i "s/\.\.pc<MeasurementSummary>(1, 'measurements', MeasurementSummary\.CreateMessage)/\.\.pc<MeasurementSummary>(1, 'measurements', \$pb.PbFieldType.PM, subBuilder: MeasurementSummary.CreateMessage)/g" "$PROTO_FILE_DART"

echo "Patch complete."
