# CI Web Gate Policy Audit (P18)

## 목적

P17에서 Flutter web release gate를 CI에 연결했다. 이번 단계는 해당 gate를 PR 필수로 유지할지 여부와 Wasm dry-run warning의 blocking 여부를 정책 문서로 고정한다.

## 결정 사항

- `release_target: js_web`
- `required_on_pull_request: true`
- `wasm_dry_run_blocking: false`

즉, JS web artifact build failure는 PR/release blocking이고, Wasm dry-run warning은 marker로 추적하되 현재는 blocking하지 않는다.

## 변경 사항

- `docs/ci/flutter-web-release-gate-policy.md`
  - Flutter web release target과 PR 필수 gate 여부를 machine-readable marker로 기록했다.
  - Wasm을 제품 릴리스 타깃으로 승격할 때 필요한 후속 계획 기준을 기록했다.
- `scripts/ci_web_gate_policy_test.sh`
  - 정책 문서 marker와 CI workflow 연결 상태를 검증한다.

## TDD 기록

- RED: `bash scripts/ci_web_gate_policy_test.sh`
  - 실패 이유: policy document 없음.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- 정책 문서와 CI workflow guard를 분리해, 문서만 있고 CI가 빠지는 상태를 방지했다.
- Wasm warning은 P16 marker와 같은 정책을 문서화했다.
- 의료/헬스케어 제품 특성상 PR에서 JS web build gate를 유지하는 보수적 선택을 했다.

## 품질 게이트

- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS

## 다음 단계 지침

- P19에서는 Wasm strict compatibility가 실제 제품 요구사항인지 검토하는 별도 readiness plan을 작성한다.
- Wasm을 지원하기 전까지는 JS web release gate를 CI 필수로 유지한다.
