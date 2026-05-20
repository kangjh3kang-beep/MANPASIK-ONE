# Flutter Web Timing Report Audit (P29)

## 목적

P28에서 web release gate duration marker 수집 정책을 만들었다. 이번 단계는 CI log의 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 표준 key-value report artifact로 변환하는 스크립트와 검증을 추가해, measured cost review에 사용할 데이터 포맷을 고정한다.

## 변경 사항

- `scripts/flutter_web_timing_report.sh`
  - `--log`, `--branch-type`, `--runner-context`, `--output` option을 지원한다.
  - log에서 numeric duration marker들을 추출한다.
  - `sample_count`, `duration_seconds_values`, `latest/min/max_duration_seconds`를 key-value report로 쓴다.
  - marker가 없으면 실패한다.
- `scripts/flutter_web_timing_report_test.sh`
  - fixture log의 duration 값 `11`, `17`을 report로 변환하는지 검증한다.
  - marker 없는 log는 실패하는지 검증한다.
- `scripts/ci_web_gate_timing_collection_policy_test.sh`
  - report artifact format, report script/test, required fields marker를 검증한다.
- `docs/ci/flutter-web-gate-timing-collection.md`
  - `report_artifact_format: key_value_v1`과 report schema를 문서화했다.

## TDD 기록

- RED: `bash scripts/flutter_web_timing_report_test.sh`
  - 실패 이유: `scripts/flutter_web_timing_report.sh` 없음.
- GREEN: 같은 테스트 재실행 PASS.
- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`
  - 실패 이유: `report_artifact_format: key_value_v1` marker 없음.
- GREEN: policy 문서 갱신 후 PASS.

## 자체 코드리뷰

- Report format은 shell에서 쉽게 생성/검증 가능한 key-value v1로 제한했다.
- 여러 duration marker가 있는 log도 처리하며, latest/min/max를 함께 기록한다.
- Branch type과 runner context를 필수 인자로 받아 sample 해석에 필요한 최소 메타데이터를 포함했다.
- JSON dependency 없이 CI shell에서 동작하도록 유지했다.

## 품질 게이트

- `bash scripts/flutter_web_timing_report_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash -n scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_web_gate_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 report artifact를 CI에서 업로드할지, 아니면 local cost-review runbook으로 먼저 운영할지 결정한다.
- `blocking_threshold_seconds`를 숫자로 승격하기 전에는 최소 5개 report artifact를 수집한다.
