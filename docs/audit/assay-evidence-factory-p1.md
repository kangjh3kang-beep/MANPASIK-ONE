# Assay Evidence Factory P1

## 목적

P0에서 제거한 측정 결과 하드코딩 위에 카트리지별 근거 manifest를 추가했다. 이번 단계는 실제 임상 성능을 주장하지 않고, 각 assay가 연구용인지 분석잠금인지 임상잠금인지 코드에서 명확히 구분하는 기반이다.

## Task 1 완료 범위

- `backend/shared/assay/registry.go`
  - `EvidenceStatusResearchOnly`, `EvidenceStatusAnalyticalLocked`, `EvidenceStatusClinicalLocked`를 추가했다.
  - `EvidenceManifest`, `AnalyticalPerformance`, `AcceptanceCriteria`를 추가했다.
  - 건강 바이오마커 기본 registry에 LOINC, UCUM, 필수 comparator method 힌트를 붙였다.
  - 실제 만파식 assay는 모두 `research_only`로 유지해 diagnostic-ready claim을 하지 않게 했다.
  - `Definition.IsDiagnosticReady()`와 `Definition.EvidenceGaps()`를 추가했다.
- `backend/shared/assay/registry_test.go`
  - glucose의 LOINC/UCUM/evidence status를 검증한다.
  - research-only assay가 diagnostic-ready가 아님을 검증한다.
  - clinical-locked definition은 LOINC, UCUM, reference method, LoD, LoQ, linearity, bias, CV, sensitivity, specificity, sample count, comparator study를 모두 갖춰야 통과한다.

## TDD 기록

- RED 1: `Evidence`, `EvidenceStatusResearchOnly`, `IsDiagnosticReady`, `EvidenceGaps`가 없어 `go test ./shared/assay`가 build fail.
- GREEN 1: evidence manifest schema와 readiness/gap method를 추가해 `go test ./shared/assay` PASS.
- REVIEW RED: 자체 리뷰에서 sensitivity/specificity acceptance threshold 누락을 발견했고, 누락 시 diagnostic-ready가 되면 안 된다는 테스트를 추가해 FAIL 확인.
- GREEN 2: `min_sensitivity_required`, `min_specificity_required` gap을 추가해 target tests PASS.

## 자체 코드리뷰

- 실제 registry의 `glucose`, `crp` 등은 LOINC/UCUM이 있어도 `research_only`이므로 `IsDiagnosticReady()`가 false다. 검증되지 않은 임상 claim이 생기지 않는다.
- `EvidenceGaps()`는 clinical lock 기준의 빈 칸을 문자열로 반환하므로 다음 CI gate나 관리자 화면에서 그대로 사용할 수 있다.
- `Evaluate()`는 기존처럼 corrected signal과 channel completeness confidence를 반환하므로 measurement-service 경로의 동작은 유지된다.
- `RequiredReference`는 연구용 assay에서 “어떤 comparator가 필요할지”를 남기는 필드이고, clinical lock의 실제 증거 필드는 `ReferenceMethod`다.

## 품질 게이트

- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS

## 다음 단계 지침

1. Task 2에서는 `MeasurementData`와 `ProcessedResult`에 evidence status를 노출하되, 기존 proto/gRPC 계약을 흔들지 않는 내부 service field부터 추가한다.
2. 외부 API 응답 확장은 proto 변경이 필요하므로 별도 Task에서 generated Go/Dart 재생성 게이트와 함께 진행한다.
3. clinical lock을 실제 assay에 부여하는 변경은 금지한다. 실측 데이터셋, comparator study, acceptance criteria가 문서와 DB에 들어온 뒤 별도 PR로 처리한다.

## Task 2 완료 범위

- `backend/services/measurement-service/internal/service/measurement.go`
  - `MeasurementData`에 `EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps`를 추가했다.
  - `ProcessedResult`에도 같은 evidence surface를 추가했다.
  - `applyAssaySemantics()`가 `assay.Definition`에서 evidence status와 diagnostic readiness를 채운다.
  - research-only assay는 처리와 저장은 가능하지만 `DiagnosticReady=false`와 evidence gaps를 유지한다.
- `backend/services/measurement-service/internal/service/measurement_test.go`
  - `TestProcessMeasurement_연구용_EvidenceStatus_노출`을 추가했다.
  - 저장된 measurement와 반환 result 모두 `research_only`, `DiagnosticReady=false`, evidence gaps 존재를 검증한다.

## Task 2 TDD 기록

- RED: `ProcessedResult`와 `MeasurementData`에 `EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps`가 없어 target test가 build fail.
- GREEN: 내부 service field와 `applyAssaySemantics()` 매핑을 추가해 target test PASS.

## Task 2 자체 코드리뷰

- 외부 proto/gRPC 계약은 변경하지 않았다. `MeasurementResult` generated schema에는 새 필드를 노출하지 않아 Dart/Go proto 재생성 리스크를 피했다.
- evidence status는 저장 전 `MeasurementData`에 채워지므로 Timescale/PostgreSQL repository 확장 시 같은 내부 값을 사용할 수 있다.
- upstream이 `DiagnosticReady=true`를 넣더라도 `applyAssaySemantics()`가 registry 결과로 덮어쓰므로 research-only assay가 임상 준비 완료로 상승하지 않는다.
- `EvidenceGaps`는 slice copy로 `ProcessedResult`에 전달해 후속 mutation 부작용을 줄였다.

## Task 2 품질 게이트

- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/service -run TestProcessMeasurement_연구용_EvidenceStatus_노출`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `git diff --check` targeted service files: PASS

## Task 3 지침

1. `scripts/assay_evidence_gate.sh`를 새로 만들고 `backend/shared/assay` package tests를 실행하는 방식으로 시작한다.
2. gate는 clinical-locked assay가 evidence gaps를 갖는 순간 실패해야 한다.
3. 현 registry의 실제 assay가 `research_only`인 상태는 실패로 보지 않는다.
4. CI에는 `ssot-governance` job 아래에 추가하되, 전체 Go test 중복 실행을 피하기 위해 shared assay gate만 좁게 호출한다.

## Task 3 완료 범위

- `backend/shared/assay/registry.go`
  - `ClinicalLockEvidenceGaps()`를 추가했다.
  - `clinical_locked` assay만 gate 대상으로 삼고, `research_only` assay는 실패 대상으로 보지 않는다.
- `backend/shared/assay/registry_test.go`
  - broken clinical lock fixture가 evidence gaps로 보고되는지 검증한다.
  - 현재 registry의 실제 assay들이 broken clinical lock 상태가 아님을 검증한다.
- `scripts/assay_evidence_gate.sh`
  - CI에서 좁게 실행 가능한 assay evidence gate를 추가했다.
  - `MANPASIK_GO_BINARY`로 Go binary override를 받을 수 있게 했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Assay evidence gate` step을 추가했다.

## Task 3 TDD 기록

- RED: `ClinicalLockEvidenceGaps`가 없어 `go test -count=1 ./shared/assay -run TestEvidenceGate`가 build fail.
- GREEN: `ClinicalLockEvidenceGaps()`와 gate script를 추가해 target tests와 script가 PASS.

## Task 3 자체 코드리뷰

- gate는 `EvidenceStatusClinicalLocked`만 검사하므로 연구용 assay가 LOINC/UCUM 또는 comparator hint만 가진 상태로도 개발과 실험을 계속할 수 있다.
- clinical lock으로 승격된 assay는 `EvidenceGaps()`가 요구하는 LOINC, UCUM, reference method, LoD, LoQ, linearity, bias, CV, sensitivity, specificity, sample count, comparator study를 모두 채워야 한다.
- 현재 registry에는 실제 clinical-locked assay가 없으므로 임상 성능 claim이 새로 생기지 않았다.
- CI hook은 전체 backend test를 중복 실행하지 않고 shared assay gate만 실행해 빠르게 실패 원인을 좁힌다.

## Task 3 품질 게이트

- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay -run TestEvidenceGate`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `git diff --check` targeted tracked files: PASS
- untracked P1 file trailing whitespace check: PASS

## P1 종료 및 다음 단계 지침

1. 외부 API 응답에 evidence status를 노출하려면 proto 변경, generated Go/Dart 재생성, Flutter contract test를 같은 단계에서 묶어야 한다.
2. clinical lock 승격은 실제 reference method dataset, analytical validation report, acceptance criteria 승인 기록이 들어온 뒤 별도 PR로만 허용한다.
3. 다음 구현 단계는 `ProcessedResult` evidence surface를 저장 계층과 gRPC/REST 응답 계약으로 확장하는 P2 계획으로 시작한다.
