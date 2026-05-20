# Flutter Evidence CI Shard Policy Audit (P24)

## 목적

P23에서 CI의 evidence gate를 `contract`, `ui`, `data` shard step으로 분리했다. 이번 단계는 그 운영 상태를 policy marker로 고정하고, policy guard가 문서와 workflow의 일치성을 직접 검증하게 한다.

## 변경 사항

- `docs/ci/flutter-evidence-ui-gate-policy.md`
  - `ci_execution_mode: category_shards`를 추가했다.
  - `ci_shard_steps: contract,ui,data`를 추가했다.
  - `aggregate_ci_step: false`를 추가했다.
  - `aggregate_local_gate_supported: true`를 추가했다.
  - CI는 shard step을 사용하고 로컬 aggregate 실행은 개발자 검증용으로 남긴다는 설명을 추가했다.
- `scripts/flutter_evidence_ui_gate_policy_test.sh`
  - CI shard 운영 marker를 검증한다.
  - `ci_shard_steps`의 각 category가 workflow에 `--category <name>` command로 연결되어 있는지 검증한다.
  - `aggregate_ci_step: false`일 때 aggregate-only CI command가 남아 있으면 실패한다.

## TDD 기록

- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`
  - 실패 이유: `ci_execution_mode: category_shards` marker 없음.
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`
  - `FLUTTER_EVIDENCE_UI_GATE_POLICY_PASS`.

## 자체 코드리뷰

- Policy 문서의 CI 실행 상태와 실제 workflow command가 같은 guard 안에서 검증된다.
- 로컬 aggregate gate는 삭제하지 않았다. 개발자 로컬 검증과 CI 실행 모드를 분리해 둘 다 유지한다.
- Aggregate CI step은 `aggregate_ci_step: false` marker에 의해 금지된다.

## 품질 게이트

- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: FAIL
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate_policy_test.sh scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## 다음 단계 지침

- 다음 단계는 Flutter web release gate의 PR 필수 정책을 유지할지, nightly/release branch 전용으로 분리할지 비용 정책을 검토한다.
- Evidence gate policy를 변경할 때는 `scripts/flutter_evidence_ui_gate_policy_test.sh`를 먼저 확장한다.
