# Measurement Evidence Safe UI P4

## 목적

P3에서 Flutter 도메인까지 전달된 measurement evidence fields를 화면 표시용 안전 문구로 변환했다. 이번 범위는 `research_only`를 진단, 정상, 위험, 확정 판정처럼 보이지 않게 만드는 UI 표면 보강이며, 저장소 영속화나 의료 상담 문구 확장은 포함하지 않는다.

## Task 1 완료 범위: Evidence Presentation Copy

- `frontend/flutter-app/lib/features/measurement/domain/measurement_evidence_presentation.dart`
  - `MeasurementEvidencePresentation`을 추가해 `badgeLabel`과 `detailText`를 제공한다.
  - `research_only`는 `연구용` 배지와 참고용 상세 문구로 변환한다.
  - `unknown`은 `검증 확인 중`으로 안전하게 실패한다.
- `frontend/flutter-app/test/features/measurement/domain/measurement_evidence_presentation_test.dart`
  - `research_only` 상세 문구가 "정상", "위험", "확정" 표현을 포함하지 않는지 검증한다.

## Task 1 TDD 기록

- RED: `MeasurementEvidencePresentation` 미정의로 copy test compile fail.
- GREEN: presentation helper를 추가해 copy test PASS.

## Task 2 완료 범위: Orchestrator Snapshot Evidence Propagation

- `frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart`
  - `MeasurementGoldenPathSnapshot`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - `serverProcessed`, `sessionEnded` snapshot이 `ProcessMeasurementResult` evidence fields를 그대로 전달한다.
- `frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`
  - 골든패스 서버 처리 snapshot이 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 보존하는지 검증한다.

## Task 2 TDD 기록

- RED: snapshot에 evidence getters가 없어 orchestrator test compile fail.
- GREEN: snapshot 필드와 서버 결과 매핑을 추가해 orchestrator test PASS.

## Task 3 완료 범위: Measure Screen Compact Badge

- `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
  - 서버 처리 및 세션 종료 단계에서 evidence status가 있으면 안전한 badge label을 표시한다.
  - evidence status가 없으면 기존 `SERVER`, `DONE` fallback을 유지한다.

## 자체 코드리뷰

- UI에는 짧은 배지 라벨만 연결했고 상세 의료 문구, 상담 유도, 위험 판정은 추가하지 않았다.
- `research_only`는 `연구용`으로 표시되며 진단 가능 상태로 승격하지 않는다.
- snapshot 기본값은 `diagnosticReady=false`, `evidenceGaps=[]`라 legacy 경로가 보수적으로 동작한다.
- `MeasureScreen` fallback은 기존 상태 텍스트를 유지해 evidence fields가 없는 환경에서도 UI 회귀가 작다.

## 품질 게이트

- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart`: FAIL, helper missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: FAIL, snapshot evidence fields missing
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos` targeted P4 Flutter files: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: PASS
- `git diff --check` targeted tracked P4 files: PASS
- P4 targeted trailing whitespace check: PASS

## 다음 단계 지침

1. P5에서 저장 계층을 확장하려면 Timescale/Postgres migration과 history response contract test를 먼저 작성한다.
2. UI 상세 설명을 늘릴 때는 규정 검토 문구를 별도 copy test로 잠그고, `research_only`에 진단/정상/위험/확정 표현을 금지한다.
3. 실제 서비스 smoke는 gateway REST와 measurement-service gRPC를 함께 띄우는 경로에서 evidence fields가 유지되는지 검증한다.
