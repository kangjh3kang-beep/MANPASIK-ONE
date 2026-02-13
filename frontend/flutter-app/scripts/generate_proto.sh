#!/usr/bin/env bash
# Generate Dart gRPC code from manpasik.proto
# Requires: protoc, dart pub global activate protoc_plugin
# Then: export PATH="$PATH:$HOME/.pub-cache/bin"
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BACKEND_PROTO="${ROOT}/../../backend/shared/proto"
INCLUDE="${ROOT}/proto_include"
OUT="${ROOT}/lib/generated"
# Use local protoc_plugin from pub cache if available
PUB_BIN="${PUB_CACHE:-$HOME/.pub-cache}/bin"
export PATH="${PUB_BIN}:$PATH"
mkdir -p "$OUT"
protoc -I="$BACKEND_PROTO" -I="$INCLUDE" \
  --dart_out=grpc:"$OUT" \
  "$BACKEND_PROTO/manpasik.proto"
echo "Generated in $OUT"
