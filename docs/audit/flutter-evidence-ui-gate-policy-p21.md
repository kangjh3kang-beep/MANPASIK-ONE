# Flutter Evidence UI Gate Policy Audit (P21)

## 목적

P20에서 Flutter evidence UI gate를 CI에 연결했다. 이번 단계는 해당 gate가 시간이 지나며 범위와 비용이 무제한으로 커지지 않도록 정책 문서와 기계 검증을 추가한다.

## 변경 사항

- `scripts/flutter_evidence_ui_gate.sh`
  - `--count` 모드를 추가해 현재 gate test file 수를 출력한다.
- `scripts/flutter_evidence_ui_gate_policy_test.sh`
  - evidence UI gate 정책 문서 marker를 검증한다.
  - gate count가 `max_test_files` 상한을 넘지 않는지 검증한다.
  - CI workflow가 policy guard를 실행하는지 검증한다.
- `docs/ci/flutter-evidence-ui-gate-policy.md`
  - gate scope, PR 필수 여부, full Flutter test와의 중복 상태, 최대 test file 수, 변경 기준을 문서화했다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job에 `Flutter evidence UI gate policy` step을 추가했다.

## TDD 기록

- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`
  - 실패 이유: policy document 없음.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- `--count` 모드는 Flutter SDK를 호출하지 않아 SSOT governance job에서 빠르게 실행된다.
- 정책 guard는 문서 marker와 CI 연결, 실제 gate 목록 수를 함께 확인해 문서와 실행 스크립트의 drift를 줄인다.
- 현재는 full Flutter test와 evidence UI gate가 일부 중복되지만, 명명된 evidence 회귀 gate의 가시성을 우선한다는 정책을 명시했다.
- test file 수 상한을 `13`으로 고정해 신규 evidence 화면 추가 시 범위 확장이나 shard 분리를 의식적으로 결정하게 했다.

## 품질 게이트

- `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- evidence gate가 13개 test file을 넘으면 category shard를 먼저 설계한다.
- full Flutter test가 추후 shard 기반으로 바뀌면 evidence UI gate를 독립 shard로 승격하고 `duplicates_full_flutter_test` 정책을 갱신한다.
