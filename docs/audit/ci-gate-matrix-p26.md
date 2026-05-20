# CI Gate Matrix Audit (P26)

## 목적

P20-P25에서 evidence, web release, Wasm readiness, security/assay gate를 순차적으로 보강했다. 이번 단계는 이 gate들의 blocking 여부, workflow 위치, policy 문서, guard script를 하나의 matrix로 고정해 향후 timing/cost 정책 변경의 기준점을 만든다.

## 변경 사항

- `docs/ci/ci-gate-matrix.md`
  - `matrix_version: 1`, `release_target: js_web`, `wasm_release_target: false` marker를 추가했다.
  - SSOT, security, assay evidence, Flutter evidence policy/workflow, Flutter web release, Wasm readiness gate를 matrix로 문서화했다.
  - blocking/non-blocking, execution mode, policy doc, guard script를 명시했다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - matrix 문서와 핵심 marker를 검증한다.
  - matrix가 참조하는 policy docs와 guard scripts가 실제 파일로 존재하는지 검증한다.
  - `.github/workflows/ci.yml`에 주요 workflow step marker가 연결되어 있는지 검증한다.

## TDD 기록

- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패 이유: matrix document 없음.
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`
  - `CI_GATE_MATRIX_POLICY_PASS`.

## 자체 코드리뷰

- Matrix는 기존 policy 문서의 결정을 중복 구현하지 않고, gate 간 관계와 위치를 색인하는 역할로 제한했다.
- Wasm readiness는 `blocking: false`로 명시해 JS web release gate와 혼동하지 않게 했다.
- Evidence gate는 CI category shard 운영 상태를 matrix에 반영했다.
- Flutter web release gate는 `release_target: js_web`와 PR/release branch blocking 정책을 반영했다.

## 품질 게이트

- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash -n scripts/ci_gate_matrix_policy_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 실제 Flutter web release gate timing capture를 추가해 P25의 cost review 기준을 수치화한다.
- 새 CI gate를 추가할 때는 matrix, policy doc, guard script, workflow 연결을 함께 갱신한다.
