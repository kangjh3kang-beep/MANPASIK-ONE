# Measurement Evidence Contract P2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** P1에서 계산한 assay evidence status를 measurement gRPC/Flutter 외부 계약까지 안전하게 노출한다.

**Architecture:** `MeasurementResult` proto에 evidence status, diagnostic readiness, evidence gaps를 backward-compatible field number로 추가한다. measurement-service handler는 내부 `ProcessedResult`에서 값을 복사하고, generated Go/Dart는 공식 generator gate로 재생성 또는 검증한다.

**Tech Stack:** Proto3, Go gRPC generated code, Dart gRPC generated code, Go 1.26.2 local toolchain, Flutter/Dart proto compile gates.

---

## Stage Gate Rules

- 각 Task는 RED 테스트를 먼저 작성하고 실패를 확인한다.
- proto 변경은 `backend/shared/proto/manpasik.proto`를 SSOT로 삼고 generated Go/Dart drift를 남기지 않는다.
- external contract 변경 후에는 Go handler tests와 Dart proto compile/preflight gate를 통과해야 한다.
- 단계 종료마다 `docs/audit/measurement-evidence-contract-p2.md`, `CHANGELOG.md`, `CONTEXT.md`에 코드리뷰와 품질 게이트를 기록한다.

## File Structure

- Modify: `backend/shared/proto/manpasik.proto`
  - `MeasurementResult`에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 field 7~9로 추가한다.
- Modify: `backend/shared/gen/go/v1/manpasik.pb.go`
  - protoc generator로 Go message struct/getter/raw descriptor를 갱신한다.
- Modify: `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
  - proto generator 실행 결과로 timestamp/header만 바뀌는지 확인한다.
- Modify: `backend/services/measurement-service/internal/handler/grpc.go`
  - `service.ProcessedResult` evidence fields를 `v1.MeasurementResult`로 복사한다.
- Modify: `backend/services/measurement-service/internal/handler/grpc_stream_test.go`
  - stream response가 `research_only`, `DiagnosticReady=false`, evidence gaps를 노출하는지 검증한다.
- Modify: `frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh`
  - compile smoke가 새 evidence fields를 사용해 Dart generated contract를 확인한다.
- Modify or generate: `frontend/flutter-app/lib/generated/manpasik.pb.dart`, `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`, `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`, `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`
  - Dart checked-in generated output을 proto와 맞춘다.
- Create: `docs/audit/measurement-evidence-contract-p2.md`
  - P2 구현 범위, 자체 코드리뷰, 품질 게이트, 잔여 리스크를 기록한다.

## Task 1: Go gRPC MeasurementResult Evidence Contract

**Files:**
- Modify: `backend/services/measurement-service/internal/handler/grpc_stream_test.go`
- Modify: `backend/shared/proto/manpasik.proto`
- Modify: `backend/shared/gen/go/v1/manpasik.pb.go`
- Modify: `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
- Modify: `backend/services/measurement-service/internal/handler/grpc.go`

- [x] **Step 1: Write failing handler contract test**

Add this assertion block to `TestStreamMeasurementStoresMeasurementAndFingerprint` after the existing response field checks:

```go
if response.EvidenceStatus != "research_only" {
    t.Fatalf("response EvidenceStatus = %q, want research_only", response.EvidenceStatus)
}
if response.DiagnosticReady {
    t.Fatal("research-only response must not be diagnostic ready")
}
assertContainsString(t, response.EvidenceGaps, "clinical_lock_required")
```

Add this helper near the other test helpers:

```go
func assertContainsString(t *testing.T, values []string, want string) {
    t.Helper()
    for _, value := range values {
        if value == want {
            return
        }
    }
    t.Fatalf("values = %v, want %q", values, want)
}
```

- [x] **Step 2: Verify RED**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/handler -run TestStreamMeasurementStoresMeasurementAndFingerprint
```

Expected: FAIL because generated `v1.MeasurementResult` has no `EvidenceStatus`, `DiagnosticReady`, or `EvidenceGaps`.

- [x] **Step 3: Extend proto contract**

Modify `MeasurementResult` in `backend/shared/proto/manpasik.proto`:

```proto
message MeasurementResult {
  string session_id = 1;
  double primary_value = 2;
  string unit = 3;
  double confidence = 4;
  repeated float fingerprint_vector = 5;
  google.protobuf.Timestamp processed_at = 6;
  string evidence_status = 7;
  bool diagnostic_ready = 8;
  repeated string evidence_gaps = 9;
}
```

- [x] **Step 4: Regenerate Go proto**

Run:

```bash
cd backend/shared/proto
PATH="/home/kangjh3kang/go/bin:$PATH" protoc -I=. \
  --go_out=../gen/go/v1 --go_opt=paths=source_relative \
  --go-grpc_out=../gen/go/v1 --go-grpc_opt=paths=source_relative \
  manpasik.proto
```

Expected: `backend/shared/gen/go/v1/manpasik.pb.go` contains `EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps` fields and getters.

- [x] **Step 5: Map service result to proto response**

Update the `stream.Send(&v1.MeasurementResult{...})` literal:

```go
EvidenceStatus:     string(result.EvidenceStatus),
DiagnosticReady:    result.DiagnosticReady,
EvidenceGaps:       append([]string(nil), result.EvidenceGaps...),
```

- [x] **Step 6: Verify GREEN**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/gofmt -w services/measurement-service/internal/handler/grpc.go services/measurement-service/internal/handler/grpc_stream_test.go
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service
```

Expected: PASS.

## Task 2: Dart Generated Contract Gate

**Files:**
- Modify: `frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh`
- Modify or generate: `frontend/flutter-app/lib/generated/manpasik.pb.dart`
- Modify or generate: `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
- Modify or generate: `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`
- Modify or generate: `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`

- [x] **Step 1: Write failing Dart compile smoke**

In `check_proto_generation_compile_gate.sh`, update the temporary `MeasurementResult` constructor:

```dart
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
```

Update the assertion:

```dart
      result.evidenceStatus != 'research_only' ||
      result.diagnosticReady ||
      !result.evidenceGaps.contains('clinical_lock_required') ||
```

- [x] **Step 2: Verify Dart generation gate**

Run:

```bash
cd frontend/flutter-app
bash scripts/check_proto_generation_preflight.sh
bash scripts/check_proto_generation_compile_gate.sh
```

Expected: PASS if Dart generator and dependencies are available; otherwise record the exact BLOCKED status.

- [x] **Step 3: Regenerate checked-in Dart output**

Run:

```bash
cd frontend/flutter-app
bash scripts/generate_proto.sh
```

Expected: checked-in generated Dart files expose `evidenceStatus`, `diagnosticReady`, and `evidenceGaps`.

## Task 3: P2 Audit and Shared Context

**Files:**
- Create: `docs/audit/measurement-evidence-contract-p2.md`
- Modify: `CHANGELOG.md`
- Modify: `CONTEXT.md`

- [x] **Step 1: Record implementation scope**

Create `docs/audit/measurement-evidence-contract-p2.md` with Task 1/2 TDD records, code review, quality gates, and residual risks.

- [x] **Step 2: Update shared logs**

Add a new top entry to `CHANGELOG.md` and a new top block to `CONTEXT.md` that records proto fields, generated outputs, tests, and next-stage instructions.

- [x] **Step 3: Final verification**

Run:

```bash
python3 scripts/validate_ssot_constants.py
bash scripts/security_release_gate.sh
bash scripts/assay_evidence_gate.sh
cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service
git diff --check -- backend/shared/proto/manpasik.proto backend/shared/gen/go/v1/manpasik.pb.go backend/shared/gen/go/v1/manpasik_grpc.pb.go backend/services/measurement-service/internal/handler/grpc.go backend/services/measurement-service/internal/handler/grpc_stream_test.go frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh CHANGELOG.md CONTEXT.md
```

Expected: PASS, except Dart generator gates may be BLOCKED only if the required local generator/toolchain is unavailable.
