#!/usr/bin/env bash
# Compile full Dart gRPC output in an isolated temporary package.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_PROTO="${ROOT}/../../backend/shared/proto"
INCLUDE="${ROOT}/proto_include"
PUB_BIN="${PUB_CACHE:-$HOME/.pub-cache}/bin"
TMP_DIR="$(mktemp -d)"

cleanup() {
  rm -rf "$TMP_DIR"
}
trap cleanup EXIT

export PATH="${PUB_BIN}:$PATH"

if ! command -v protoc >/dev/null 2>&1; then
  echo "PROTO_COMPILE_GATE_STATUS=blocked_missing_protoc"
  exit 1
fi

if ! command -v protoc-gen-dart >/dev/null 2>&1; then
  echo "PROTO_COMPILE_GATE_STATUS=blocked_missing_protoc_gen_dart"
  echo "HINT=dart pub global activate protoc_plugin"
  exit 1
fi

if ! command -v dart >/dev/null 2>&1; then
  echo "PROTO_COMPILE_GATE_STATUS=blocked_missing_dart"
  exit 1
fi

mkdir -p "$TMP_DIR/lib/generated" "$TMP_DIR/bin"

protoc -I="$BACKEND_PROTO" -I="$INCLUDE" \
  --dart_out=grpc:"$TMP_DIR/lib/generated" \
  "$BACKEND_PROTO/manpasik.proto"

PB_DART="$TMP_DIR/lib/generated/manpasik.pb.dart"
PB_GRPC_DART="$TMP_DIR/lib/generated/manpasik.pbgrpc.dart"

if ! grep -q "class MeasurementData" "$PB_DART"; then
  echo "PROTO_COMPILE_GATE_STATUS=failed_missing_measurement_data"
  exit 1
fi

if ! grep -q "streamMeasurement" "$PB_GRPC_DART"; then
  echo "PROTO_COMPILE_GATE_STATUS=failed_missing_streammeasurement"
  exit 1
fi

cat > "$TMP_DIR/pubspec.yaml" <<'YAML'
name: manpasik_proto_compile_gate
publish_to: none
environment:
  sdk: ">=3.8.0 <4.0.0"
dependencies:
  fixnum: ^1.1.1
  grpc: ^5.1.0
  protobuf: ^6.0.0
YAML

cat > "$TMP_DIR/bin/compile_smoke.dart" <<'DART'
import 'package:manpasik_proto_compile_gate/generated/manpasik.pbgrpc.dart';

void main() {
  final data = MeasurementData(
    sessionId: 'session-compile-1',
    rawChannels: [1, 2, 3],
    differential: DifferentialCorrection(
      sDet: 10,
      sRef: 2,
      alpha: 0.95,
      sCorrected: 8.1,
    ),
    envMeta: EnvironmentMeta(tempC: 25, humidityPct: 40),
  );
  final result = MeasurementResult(
    sessionId: data.sessionId,
    primaryValue: data.differential.sCorrected,
    unit: 'mg/dL',
    confidence: 0.95,
    fingerprintVector: data.rawChannels.map((value) => value.toDouble()).toList(),
    evidenceStatus: 'research_only',
    diagnosticReady: false,
    evidenceGaps: ['clinical_lock_required'],
  );
  final clientType = MeasurementServiceClient;
  if (data.sessionId != 'session-compile-1' ||
      result.primaryValue != 8.1 ||
      result.evidenceStatus != 'research_only' ||
      result.diagnosticReady ||
      !result.evidenceGaps.contains('clinical_lock_required') ||
      clientType.toString().isEmpty) {
    throw StateError('generated proto compile smoke failed');
  }
}
DART

(
  cd "$TMP_DIR"
  dart pub get >/dev/null
  dart analyze bin lib
  dart run bin/compile_smoke.dart
)

generated_lines="$(wc -l "$TMP_DIR"/lib/generated/*.dart | tail -n 1 | awk '{print $1}')"
grpc_version="$(awk '
  /^  grpc:/ { in_pkg=1; next }
  in_pkg && /^  [^ ]/ { in_pkg=0 }
  in_pkg && /version:/ {
    gsub(/"/, "", $2)
    print $2
    exit
  }
' "$TMP_DIR/pubspec.lock")"
protobuf_version="$(awk '
  /^  protobuf:/ { in_pkg=1; next }
  in_pkg && /^  [^ ]/ { in_pkg=0 }
  in_pkg && /version:/ {
    gsub(/"/, "", $2)
    print $2
    exit
  }
' "$TMP_DIR/pubspec.lock")"

echo "PROTO_COMPILE_GATE_STATUS=passed"
echo "PROTO_COMPILE_GATE_OUTPUT_DIR=$TMP_DIR"
echo "PROTO_COMPILE_GATE_GENERATED_DART_LINES=$generated_lines"
echo "PROTO_COMPILE_GATE_GRPC_VERSION=${grpc_version:-unknown}"
echo "PROTO_COMPILE_GATE_PROTOBUF_VERSION=${protobuf_version:-unknown}"
