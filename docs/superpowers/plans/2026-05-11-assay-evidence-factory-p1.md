# Assay Evidence Factory P1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 카트리지별 분석 성능 근거를 코드로 추적해 만파식 범용분석시스템이 연구용/분석잠금/임상잠금 상태를 명확히 구분하게 한다.

**Architecture:** `backend/shared/assay`를 assay manifest의 중심으로 삼고, measurement-service는 이 registry가 계산한 의미론만 사용한다. P1에서는 실제 성능 수치를 주장하지 않고, LOINC/UCUM/필수 reference method/validation status/acceptance criteria 구조를 먼저 도입한다.

**Tech Stack:** Go 1.24, `backend/shared/assay`, measurement-service Go tests, shell quality gates.

---

## Stage Gate Rules

- 각 Task는 RED 테스트를 먼저 작성하고 실패를 확인한다.
- GREEN 구현 후 `gofmt`, 대상 `go test`, 관련 `git diff --check`를 통과해야 한다.
- 단계 종료마다 자체 코드리뷰를 수행하고 `docs/audit/*`에 리뷰 결과와 잔여 리스크를 기록한다.
- 한 단계가 통과하기 전 다음 단계 구현을 시작하지 않는다.

## File Structure

- Modify: `backend/shared/assay/registry.go`
  - `EvidenceStatus`, `EvidenceManifest`, `AcceptanceCriteria`, `ReferenceMethodRequirement`를 추가한다.
  - `Definition`이 evidence manifest를 포함하게 한다.
  - `IsDiagnosticReady()`와 `EvidenceGaps()`로 운영 claim 가능 여부를 판단한다.
- Modify: `backend/shared/assay/registry_test.go`
  - LOINC/UCUM/validation status, research-only diagnostic gate, clinical lock 필수 필드 검증 테스트를 추가한다.
- Create: `docs/audit/assay-evidence-factory-p1.md`
  - P1 구현 범위, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 기록한다.
- Modify: `CHANGELOG.md`, `CONTEXT.md`
  - 작업 기록 프로토콜에 따라 완료 단계와 다음 단계 지침을 반영한다.

## Task 1: Evidence Manifest Schema

**Files:**
- Modify: `backend/shared/assay/registry.go`
- Modify: `backend/shared/assay/registry_test.go`
- Create: `docs/audit/assay-evidence-factory-p1.md`

- [x] **Step 1: Write failing tests**

Add tests that require:

```go
definition, err := Resolve("glucose")
if err != nil {
    t.Fatal(err)
}
if definition.Evidence.LOINCCode != "15074-8" {
    t.Fatalf("LOINCCode = %q, want 15074-8", definition.Evidence.LOINCCode)
}
if definition.Evidence.UCUMUnit != "mg/dL" {
    t.Fatalf("UCUMUnit = %q, want mg/dL", definition.Evidence.UCUMUnit)
}
if definition.IsDiagnosticReady() {
    t.Fatal("research-only assay must not be diagnostic ready")
}
if len(definition.EvidenceGaps()) == 0 {
    t.Fatal("research-only assay must expose evidence gaps")
}
```

Also add a clinical-lock test that constructs a complete `Definition` and expects `IsDiagnosticReady() == true`.

- [x] **Step 2: Verify RED**

Run: `cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay`

Expected: FAIL because `Evidence`, `IsDiagnosticReady`, and `EvidenceGaps` are not defined.

- [x] **Step 3: Implement minimal schema**

Add:

```go
type EvidenceStatus string

const (
    EvidenceStatusResearchOnly EvidenceStatus = "research_only"
    EvidenceStatusAnalyticalLocked EvidenceStatus = "analytical_locked"
    EvidenceStatusClinicalLocked EvidenceStatus = "clinical_locked"
)
```

Then add evidence structs and methods to satisfy the tests.

- [x] **Step 4: Verify GREEN**

Run:

```bash
cd backend
/home/kangjh3kang/sdk/go-go1.26.2/bin/gofmt -w shared/assay/registry.go shared/assay/registry_test.go
/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay
```

Expected: PASS.

- [x] **Step 5: Self Review**

Review the diff for:

- no clinical-ready claim on real ManPaSik assays without evidence;
- no hardcoded false certainty;
- no empty evidence fields on clinical locked definitions;
- no breakage to existing `Evaluate()` behavior.

Record the review in `docs/audit/assay-evidence-factory-p1.md`.

## Task 2: Measurement Evidence Surface

**Files:**
- Modify: `backend/services/measurement-service/internal/service/measurement.go`
- Modify: `backend/services/measurement-service/internal/service/measurement_test.go`
- Modify: `docs/audit/assay-evidence-factory-p1.md`

- [x] **Step 1: Write failing tests**

Add a service test that verifies a research-only assay can still process measurements but never marks the measurement as diagnostically validated.

- [x] **Step 2: Implement minimal surface**

Add evidence status to `MeasurementData` and `ProcessedResult`, filled from `assay.Definition.Evidence.Status`.

- [x] **Step 3: Verify**

Run: `cd backend && /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/service ./services/measurement-service/internal/handler`

Expected: PASS.

## Task 3: Evidence Gate CI Hook

**Files:**
- Modify: `scripts/security_release_gate.sh` or create `scripts/assay_evidence_gate.sh`
- Modify: `.github/workflows/ci.yml`
- Modify: `docs/audit/assay-evidence-factory-p1.md`

- [x] **Step 1: Add gate**

Add a gate that fails if any clinical-locked assay lacks LOINC, UCUM, reference method, LoD, LoQ, or acceptance criteria.

- [x] **Step 2: Verify**

Run the gate and target Go tests.

## Current Execution Choice

The user requested sequential implementation, so this plan will be executed inline with a review and verification checkpoint after each task.
