# Wasm Readiness Audit (P19)

## 목적

P16-P18에서 Flutter JS web build gate와 CI 정책을 고정했다. 이번 단계는 Wasm을 아직 릴리스 타깃으로 보지 않는다는 상태와 현재 dry-run blocker를 별도 readiness 문서로 남겨, 향후 Wasm 지원 전환 시 필요한 작업 범위를 명확히 한다.

## 결정 사항

- `wasm_release_target: false`
- `current_release_target: js_web`
- `wasm_dry_run_blocking: false`

## 변경 사항

- `docs/ci/flutter-wasm-readiness.md`
  - 현재 Wasm release target 상태와 dry-run blocker 목록을 기록했다.
  - blocker package: `flutter_secure_storage_web`, `share_plus`, `connectivity_plus`, `package:js`.
  - Wasm target 승격 기준을 문서화했다.
- `scripts/ci_wasm_readiness_policy_test.sh`
  - readiness 문서가 필수 marker와 blocker package 이름을 포함하는지 검증한다.

## TDD 기록

- RED: `bash scripts/ci_wasm_readiness_policy_test.sh`
  - 실패 이유: readiness document 없음.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- Wasm blocker를 JS web release gate와 혼동하지 않도록 별도 readiness 문서로 분리했다.
- 보안/PHI 측면에서 `flutter_secure_storage_web` 대체는 별도 설계 없이 진행하지 않도록 승격 기준에 명시했다.
- 현재 CI release policy와 모순되지 않도록 `wasm_dry_run_blocking: false`를 유지했다.

## 품질 게이트

- `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash -n scripts/ci_wasm_readiness_policy_test.sh scripts/ci_web_gate_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS

## 다음 단계 지침

- Wasm을 실제 제품 타깃으로 전환한다는 결정이 나오기 전까지는 P16-P18의 JS web release gate 정책을 유지한다.
- 다음 일반 강화 단계는 전체 CI gate 비용과 병렬화 최적화 또는 DataHub 외 화면의 evidence badge 회귀 확장이다.
