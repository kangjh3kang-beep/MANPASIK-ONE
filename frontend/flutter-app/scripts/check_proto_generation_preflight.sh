#!/usr/bin/env bash
# Preflight Dart gRPC generation without overwriting checked-in manual stubs.
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
  echo "PROTO_PREFLIGHT_STATUS=blocked_missing_protoc"
  exit 1
fi

if ! command -v protoc-gen-dart >/dev/null 2>&1; then
  echo "PROTO_PREFLIGHT_STATUS=blocked_missing_protoc_gen_dart"
  echo "HINT=dart pub global activate protoc_plugin"
  exit 1
fi

protoc -I="$BACKEND_PROTO" -I="$INCLUDE" \
  --dart_out=grpc:"$TMP_DIR" \
  "$BACKEND_PROTO/manpasik.proto"

PB_DART="$TMP_DIR/manpasik.pb.dart"
PB_GRPC_DART="$TMP_DIR/manpasik.pbgrpc.dart"

if ! grep -q "class MeasurementData" "$PB_DART"; then
  echo "PROTO_PREFLIGHT_STATUS=failed_missing_measurement_data"
  exit 1
fi

if ! grep -q "streamMeasurement" "$PB_GRPC_DART"; then
  echo "PROTO_PREFLIGHT_STATUS=failed_missing_streammeasurement"
  exit 1
fi

generated_lines="$(wc -l "$TMP_DIR"/*.dart | tail -n 1 | awk '{print $1}')"
echo "PROTO_PREFLIGHT_STATUS=generated_streammeasurement_available"
echo "PROTO_PREFLIGHT_OUTPUT_DIR=$TMP_DIR"
echo "PROTO_PREFLIGHT_GENERATED_DART_LINES=$generated_lines"

timestamp_import="$(grep -E "protobuf/.+timestamp.pb.dart" "$PB_DART" || true)"
if [[ -n "$timestamp_import" ]]; then
  timestamp_path="$(echo "$timestamp_import" | sed -E "s/.*package:protobuf\\/(.*)'.*/\\1/" | head -n 1)"
  pub_cache="${PUB_CACHE:-$HOME/.pub-cache}/hosted/pub.dev"
  locked_protobuf_version="$(awk '
    /^  protobuf:/ { in_protobuf=1; next }
    in_protobuf && /^  [^ ]/ { in_protobuf=0 }
    in_protobuf && /version:/ {
      gsub(/"/, "", $2)
      print $2
      exit
    }
  ' "$ROOT/pubspec.lock")"
  if [[ -n "$locked_protobuf_version" ]]; then
    locked_timestamp_file="$pub_cache/protobuf-${locked_protobuf_version}/lib/$timestamp_path"
  else
    locked_timestamp_file=""
  fi
  if [[ -z "$locked_timestamp_file" || ! -f "$locked_timestamp_file" ]]; then
    echo "PROTO_PREFLIGHT_FULL_REPLACEMENT=blocked_timestamp_import_incompatible_with_current_pub_cache"
    echo "PROTO_PREFLIGHT_LOCKED_PROTOBUF_VERSION=${locked_protobuf_version:-unknown}"
    echo "PROTO_PREFLIGHT_TIMESTAMP_IMPORT=$timestamp_import"
    exit 0
  fi
fi

echo "PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate"
