#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT/backend"

GO_BIN="${MANPASIK_GO_BINARY:-/home/kangjh3kang/sdk/go-go1.26.2/bin/go}"
if [[ ! -x "$GO_BIN" ]]; then
  GO_BIN="${GO_BIN_FALLBACK:-go}"
fi

"$GO_BIN" test -count=1 ./shared/assay -run 'TestEvidenceGate'
echo "ASSAY_EVIDENCE_GATE_PASS"
