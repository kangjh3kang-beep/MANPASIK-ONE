# Flutter Evidence UI Gate Audit (P20)

## 목적

P3-P15에서 측정 evidence 계약, safe copy, home/result/DataHub badge, REST/gRPC mapper 호환성을 순차 보강했다. 이번 단계는 해당 회귀 테스트를 하나의 Flutter evidence UI gate로 묶고 CI에서 실행해, evidence 표시와 매핑 계약이 일반 테스트 집합 속에서 흐려지지 않도록 한다.

## 변경 사항

- `scripts/flutter_evidence_ui_gate.sh`
  - evidence 관련 Flutter test file 목록을 한 번에 실행한다.
  - `--list` 모드로 게이트 포함 범위를 빠르게 검증할 수 있게 했다.
  - 테스트 통과 후 `FLUTTER_EVIDENCE_UI_GATE_PASS` marker를 출력한다.
- `scripts/flutter_evidence_ui_gate_test.sh`
  - gate list가 핵심 evidence presentation, badge, DataHub, REST mapper, generated contract tests를 포함하는지 검증한다.
- `scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`
  - CI workflow가 `Flutter evidence UI gate` step과 `scripts/flutter_evidence_ui_gate.sh` 호출을 포함하는지 검증한다.
- `.github/workflows/ci.yml`
  - Flutter job의 `Run tests` 이후, `Flutter web release gate` 이전에 evidence UI gate를 실행한다.

## TDD 기록

- RED: `bash scripts/flutter_evidence_ui_gate_test.sh`
  - 실패 이유: `scripts/flutter_evidence_ui_gate.sh` 없음.
- GREEN: `bash scripts/flutter_evidence_ui_gate_test.sh`
  - `FLUTTER_EVIDENCE_UI_GATE_TEST_PASS`.
- RED: `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`
  - 실패 이유: CI workflow에 evidence gate step 없음.
- GREEN: `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`
  - `CI_FLUTTER_EVIDENCE_UI_GATE_WORKFLOW_PASS`.

## 자체 코드리뷰

- gate는 새 테스트를 중복 작성하지 않고 기존 evidence 회귀 테스트를 묶는 방식으로 유지했다.
- CI working directory가 `frontend/flutter-app`이므로 workflow에서는 `../../scripts/flutter_evidence_ui_gate.sh`를 호출한다.
- `--list` 모드는 Flutter SDK 없이도 gate 범위 drift를 빠르게 탐지하도록 남겼다.
- P16-P19의 web release gate 정책과 충돌하지 않도록 evidence gate를 web release gate 앞에 배치했다.

## 품질 게이트

- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계에서는 evidence gate 비용을 줄이기 위해 test shard/분류 정책을 검토한다.
- 또는 DataHub/measurement 외 화면에 evidence badge가 추가될 때 `scripts/flutter_evidence_ui_gate.sh --list` 범위를 함께 갱신한다.
