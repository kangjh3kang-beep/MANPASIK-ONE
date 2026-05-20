# Measurement Evidence Contract P2

## 목적

P1에서 내부적으로 계산하던 assay evidence status를 measurement gRPC 응답 계약으로 확장했다. Task 1은 Go gRPC 경로를 닫았고, Task 2는 Flutter/Dart checked-in generated output까지 같은 계약으로 정렬했다.

## Task 1 완료 범위

- `backend/shared/proto/manpasik.proto`
  - `MeasurementResult`에 `evidence_status = 7`, `diagnostic_ready = 8`, `evidence_gaps = 9`를 추가했다.
  - 기존 field number 1~6은 유지해 backward-compatible 확장으로 처리했다.
- `backend/shared/gen/go/v1/manpasik.pb.go`
  - protoc로 Go generated message struct/getters/raw descriptor를 갱신했다.
- `backend/shared/gen/go/v1/manpasik_grpc.pb.go`
  - 같은 protoc 실행 결과를 반영했다.
- `backend/services/measurement-service/internal/handler/grpc.go`
  - `ProcessedResult.EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps`를 `v1.MeasurementResult`로 복사한다.
- `backend/services/measurement-service/internal/handler/grpc_stream_test.go`
  - stream response가 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 노출하는지 검증한다.

## TDD 기록

- RED: `MeasurementResult` generated type에 `EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps`가 없어 handler target test가 build fail.
- GREEN: proto field 7~9 추가, Go generated 재생성, handler mapping 추가 후 handler target test와 measurement target suite가 PASS.

## 자체 코드리뷰

- 새 proto field는 7~9만 사용해 기존 wire field 1~6을 변경하지 않았다.
- handler는 service layer가 계산한 registry 기반 evidence만 복사하므로, gRPC 경로에서 별도 임상 판정을 만들지 않는다.
- `EvidenceGaps`는 응답 생성 시 slice copy로 전달해 service result mutation 부작용을 줄였다.
- 현재 실제 assay는 `research_only`이므로 외부 응답에도 임상 준비 완료 claim이 노출되지 않는다.
- Dart checked-in generated output은 Task 2에서 갱신되므로 Go/Dart contract drift를 남기지 않는다.

## 품질 게이트

- RED: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/handler -run TestStreamMeasurementStoresMeasurementAndFingerprint`: FAIL, generated field missing
- GREEN: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/handler -run TestStreamMeasurementStoresMeasurementAndFingerprint`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `git diff --check` targeted tracked files: PASS
- untracked P1/P2 file trailing whitespace check: PASS

## Task 2 완료 범위

- `frontend/flutter-app/test/generated/measurement_result_evidence_contract_test.dart`
  - checked-in Dart `MeasurementResult`가 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 직렬화/역직렬화하는지 검증한다.
- `frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh`
  - 임시 generated Dart compile smoke가 새 evidence fields를 생성자와 runtime assertion에서 사용하게 했다.
- `frontend/flutter-app/lib/generated/manpasik.pb.dart`
  - 공식 Dart generator 출력으로 `MeasurementResult` evidence fields/getters/setters를 갱신했다.
- `frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbenum.dart`
- `frontend/flutter-app/lib/generated/manpasik.pbjson.dart`
  - 동일 proto 기준의 checked-in generated output을 갱신했다.

## Task 2 TDD 기록

- RED: checked-in Dart generated `MeasurementResult`에 evidence fields가 없어 `flutter test --no-pub test/generated/measurement_result_evidence_contract_test.dart`가 compile fail.
- GREEN: `bash scripts/generate_proto.sh`로 checked-in Dart generated files를 재생성한 뒤 같은 Flutter test가 PASS.

## Task 2 자체 코드리뷰

- Dart compile gate는 임시 fresh generated output을 사용하므로, proto generator가 새 fields를 실제로 생성하는지 확인한다.
- checked-in Dart contract test는 저장소에 커밋된 `lib/generated/manpasik.pb.dart`를 직접 사용하므로, proto와 checked-in generated drift를 잡는다.
- 새 field 이름은 Dart generator의 lowerCamelCase 규칙인 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`로 고정했다.
- 임상 claim은 여전히 `research_only`와 `diagnosticReady=false` 테스트 데이터만 사용한다.

## Task 2 품질 게이트

- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/generated/measurement_result_evidence_contract_test.dart`: FAIL, checked-in generated fields missing
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
- `bash scripts/check_proto_generation_compile_gate.sh`: PASS, `PROTO_COMPILE_GATE_STATUS=passed`
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/generated/measurement_result_evidence_contract_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos test/generated/measurement_result_evidence_contract_test.dart lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/generated/manpasik.pbenum.dart lib/generated/manpasik.pbjson.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `git diff --check` targeted tracked files: PASS
- untracked P1/P2 file trailing whitespace check: PASS

## 다음 단계 지침

1. 다음 단계는 REST/gateway 응답 또는 Flutter UI 표시 계층에 evidence status를 노출하는 별도 P3로 분리한다.
2. UI에 표시할 때는 `research_only`를 진단 가능 또는 의료적 확정 표현으로 번역하지 않는다.
3. 저장 계층에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 영속화하려면 schema migration과 repository contract test를 먼저 작성한다.
