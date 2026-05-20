# Flutter Web Gate Cost Policy Audit (P25)

## 목적

P16-P19에서 Flutter web release gate와 Wasm readiness 정책을 만들었다. 이번 단계는 web build gate의 CI 비용 정책을 명시해, JS web release target이 유지되는 동안 PR과 release branch에서 blocking gate를 계속 실행한다는 결정을 기계 검증으로 고정한다.

## 변경 사항

- `scripts/ci_web_gate_policy_test.sh`
  - `ci_execution_mode: pull_request_and_release_branch` marker를 검증한다.
  - `release_branch_required: true` marker를 검증한다.
  - `nightly_split: false` marker를 검증한다.
  - `cost_review_required_before_relaxing: true` marker를 검증한다.
  - `blocking_build_command: flutter build web --no-pub` marker를 검증한다.
- `docs/ci/flutter-web-release-gate-policy.md`
  - JS web이 현재 release artifact이므로 PR blocking을 유지한다고 기록했다.
  - nightly-only 또는 release-branch-only로 완화하려면 측정 기반 CI cost review와 별도 policy update가 필요하다고 명시했다.

## TDD 기록

- RED: `bash scripts/ci_web_gate_policy_test.sh`
  - 실패 이유: `ci_execution_mode: pull_request_and_release_branch` marker 없음.
- GREEN: `bash scripts/ci_web_gate_policy_test.sh`
  - `CI_WEB_GATE_POLICY_PASS`.

## 자체 코드리뷰

- 기존 `required_on_pull_request: true`, `release_target: js_web`, `wasm_dry_run_blocking: false` 정책과 모순되지 않게 marker를 추가했다.
- 실제 workflow 연결은 기존 `Flutter web release gate` step guard로 유지했다.
- 비용 완화는 구현하지 않고, 완화 전 측정 기반 검토를 요구하도록 정책화했다.

## 품질 게이트

- RED: `bash scripts/ci_web_gate_policy_test.sh`: FAIL
- GREEN: `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash -n scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 실제 Flutter web release gate 시간을 기록하는 timing capture를 추가하거나, CI 전체 gate matrix 문서를 만든다.
- `nightly_split: true`로 바꾸기 전에는 timing evidence와 release risk review를 먼저 남긴다.
