# Flutter Web Release Gate Timing Audit (P27)

## 목적

P25에서 Flutter web release gate를 PR/release branch blocking으로 유지하되, 완화 전에는 측정 기반 cost review가 필요하다고 정책화했다. 이번 단계는 release gate가 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 출력하게 해 향후 CI 비용 검토의 수치 근거를 남긴다.

## 변경 사항

- `scripts/flutter_web_release_gate.sh`
  - 실제 build 경로에서 epoch seconds로 elapsed time을 측정한다.
  - `policy_check`가 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` env를 numeric으로 검증하고 marker를 출력한다.
  - invalid duration value는 실패한다.
- `scripts/flutter_web_release_gate_test.sh`
  - policy-check 성공 fixture에서 duration env를 주입하고 timing marker 출력을 검증한다.
- `scripts/ci_web_gate_policy_test.sh`
  - `timing_capture: true`와 `timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` marker를 요구한다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - CI gate matrix도 timing marker를 포함하는지 검증한다.
- `docs/ci/flutter-web-release-gate-policy.md`
  - web gate run이 duration seconds marker를 출력해야 한다고 문서화했다.
- `docs/ci/ci-gate-matrix.md`
  - timing capture marker를 matrix marker block에 추가했다.

## TDD 기록

- RED: `bash scripts/flutter_web_release_gate_test.sh`
  - 실패 이유: duration seconds marker 없음.
- GREEN: `bash scripts/flutter_web_release_gate_test.sh`
  - `FLUTTER_WEB_RELEASE_GATE_TEST_PASS`.
- RED: `bash scripts/ci_web_gate_policy_test.sh`
  - 실패 이유: timing capture marker 없음.
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패 이유: matrix timing marker 없음.
- GREEN: 두 policy guard 재실행 PASS.

## 자체 코드리뷰

- `--policy-check` 경로는 실제 build를 돌리지 않고 marker format을 검증할 수 있게 유지했다.
- 실제 build path는 build command 직전/직후 시간을 측정해 `policy_check`에 전달한다.
- duration marker는 숫자만 허용해 CI log parser가 안정적으로 읽을 수 있게 했다.
- 이번 단계에서는 실제 full web build를 다시 실행하지 않고 parser/policy/common gate를 검증했다.

## 품질 게이트

- `bash scripts/flutter_web_release_gate_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 CI logs에서 duration marker를 수집하는 문서 또는 timing threshold policy를 추가한다.
- `nightly_split: true`를 검토하려면 여러 CI run의 duration marker를 먼저 수집한다.
