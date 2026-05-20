# CI Flutter Evidence UI Shards Audit (P23)

## 목적

P22에서 Flutter evidence UI gate가 `contract`, `ui`, `data` shard를 지원하게 되었다. 이번 단계는 CI의 aggregate evidence gate를 category shard step으로 분리해 실패 위치를 더 빠르게 파악하고, 향후 CI 비용 최적화 기준을 명확히 한다.

## 변경 사항

- `.github/workflows/ci.yml`
  - 기존 `Flutter evidence UI gate` aggregate step을 제거했다.
  - `Flutter evidence UI gate: contract` step을 추가했다.
  - `Flutter evidence UI gate: ui` step을 추가했다.
  - `Flutter evidence UI gate: data` step을 추가했다.
- `scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`
  - 세 shard step 이름과 `--category` 호출을 검증한다.
  - aggregate-only `run: bash ../../scripts/flutter_evidence_ui_gate.sh`가 남아 있으면 실패한다.

## TDD 기록

- RED: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`
  - 실패 이유: `Flutter evidence UI gate: contract` marker 없음.
- GREEN: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`
  - `CI_FLUTTER_EVIDENCE_UI_GATE_SHARDS_WORKFLOW_PASS`.

## 자체 코드리뷰

- 기존 workflow guard는 step prefix와 script path를 계속 검증하므로 P20 호환성이 유지된다.
- 신규 shard workflow guard는 aggregate-only command 제거까지 확인해 중복 CI 실행 비용을 막는다.
- `Flutter web release gate`는 shard 세 step 이후에 유지해 evidence regression이 먼저 드러나게 했다.
- 실제 shard 실행은 `contract`, `ui`, `data` 모두 통과했다.

## 품질 게이트

- RED: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: FAIL
- GREEN: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash -n scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_shard_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category contract`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category ui`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category data`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계에서는 CI policy 문서에 aggregate-to-shard 운영 상태를 marker로 고정하거나, Flutter web release gate도 비용 정책에 따라 PR/nightly 분리를 검토한다.
- 새 evidence test 추가 시 CI shard step이 어느 category에 속하는지 함께 검토한다.
