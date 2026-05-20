# Flutter Evidence UI Gate Shards Audit (P22)

## 목적

P21에서 Flutter evidence UI gate의 전체 test file 상한을 13개로 고정했다. 현재 gate가 이미 13개에 도달했으므로, 이번 단계는 상한을 늘리기 전에 `contract`, `ui`, `data` category shard를 도입해 범위와 비용을 작게 관리할 수 있게 한다.

## 변경 사항

- `scripts/flutter_evidence_ui_gate.sh`
  - `--list-categories`를 추가했다.
  - `--category contract|ui|data`를 추가했다.
  - `--list`와 `--count`가 category option을 함께 처리하도록 확장했다.
  - 기본 실행은 기존처럼 전체 13개 evidence test를 실행한다.
- `scripts/flutter_evidence_ui_gate_shard_test.sh`
  - category 목록, 전체/카테고리별 count, shard별 핵심 test 포함 여부, unknown category 실패를 검증한다.
- `scripts/flutter_evidence_ui_gate_policy_test.sh`
  - `category_shards`와 `max_test_files_per_category` marker를 검증한다.
  - 각 category count가 shard 상한을 넘지 않는지 검증한다.
- `docs/ci/flutter-evidence-ui-gate-policy.md`
  - `category_shards: contract,ui,data`와 `max_test_files_per_category: 5`를 추가했다.
  - 새 evidence test가 category 상한을 넘길 경우 shard split 또는 정책 근거 갱신이 필요하다고 기록했다.

## TDD 기록

- RED: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`
  - 실패 이유: category support 없음. `--list-categories`가 실제 Flutter test 경로로 떨어졌고 category 목록 검증 실패.
- GREEN: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`
  - `FLUTTER_EVIDENCE_UI_GATE_SHARD_TEST_PASS`.
- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`
  - 실패 이유: `category_shards: contract,ui,data` policy marker 없음.
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`
  - `FLUTTER_EVIDENCE_UI_GATE_POLICY_PASS`.

## 자체 코드리뷰

- Category별 테스트 수는 `contract=4`, `ui=4`, `data=5`로 모두 `max_test_files_per_category: 5` 이하이다.
- Aggregate 기본 실행은 기존 CI step과 호환되며, 추후 CI shard 분리가 필요하면 같은 스크립트의 `--category` option만 사용하면 된다.
- Unknown category는 non-zero로 실패시켜 오타가 전체 gate 실행으로 조용히 대체되지 않게 했다.
- `--list`와 `--count`는 Flutter SDK를 호출하지 않아 policy/guard에서 빠르게 실행된다.

## 품질 게이트

- RED: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: FAIL
- GREEN: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: FAIL
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category contract`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category ui`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category data`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 CI 비용 최적화 단계에서는 `.github/workflows/ci.yml`에서 aggregate evidence gate를 category shard step으로 나눌지 결정한다.
- 새 evidence test 추가 시 먼저 `scripts/flutter_evidence_ui_gate.sh --count --category <name>`로 shard 상한을 확인한다.
