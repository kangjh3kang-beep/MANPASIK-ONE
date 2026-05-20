# Flutter Web Timing Collection Policy Audit (P28)

## 목적

P27에서 Flutter web release gate가 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 출력하게 되었다. 이번 단계는 이 marker를 CI 로그에서 어떻게 수집하고, 언제 PR-required gate 완화를 검토할 수 있는지 문서/guard/CI matrix로 고정한다.

## 변경 사항

- `scripts/ci_web_gate_timing_collection_policy_test.sh`
  - timing collection policy 문서와 web release policy, CI matrix, workflow 연결을 검증한다.
  - 최소 sample 수, advisory threshold 상태, nightly split 전환 조건을 marker로 검증한다.
- `docs/ci/flutter-web-gate-timing-collection.md`
  - `source_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS`와 `log_query_pattern`을 기록했다.
  - `minimum_samples_before_relaxing: 5`를 명시했다.
  - `threshold_policy_status: advisory`, `blocking_threshold_seconds: unset`, `current_nightly_split: false`를 기록했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Flutter web timing collection policy` step을 추가했다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - timing collection policy 문서와 guard script, matrix marker, workflow step을 검증한다.
- `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing collection policy` row를 추가했다.

## TDD 기록

- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`
  - 실패 이유: `docs/ci/flutter-web-gate-timing-collection.md` 없음.
- GREEN: 같은 guard 재실행 PASS.
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패 이유: matrix에 `Flutter web timing collection policy` marker 없음.
- GREEN: matrix row 추가 후 PASS.

## 자체 코드리뷰

- Timing threshold는 아직 blocking으로 만들지 않고 advisory로 유지했다.
- PR gate 완화는 최소 5개 successful CI sample과 measured cost review 이후로 제한했다.
- Timing collection guard를 `ssot-governance`에 연결해 문서 drift가 CI에서 잡히게 했다.
- Matrix도 새 policy와 guard를 참조하므로 gate 관계가 한 곳에서 추적된다.

## 품질 게이트

- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 CI duration samples를 저장하는 간단한 artifact/report 포맷을 정하거나, threshold policy를 실제 숫자로 승격할 조건을 작성한다.
- `blocking_threshold_seconds`를 숫자로 바꾸려면 sample dataset과 release risk review를 먼저 남긴다.
