##_2026-05-21_Codex_go_work_assay_gate_alignment
**status**: merged_followup_done
**plan**:
- main fast-forward 병합 후 clean worktree에서 `assay_evidence_gate`가 `go.work`의 낮은 Go 버전 때문에 실패하는 문제를 보정한다.
**changes**:
- go.work
  - workspace Go 버전을 1.25.0으로 정렬해 backend/shared/assay evidence gate가 clean main에서도 실행되게 했다.
**quality_gate**:
- RED on clean main: `bash scripts/assay_evidence_gate.sh`: FAIL, `go.work lists go 1.21`
- GREEN after go.work alignment: `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --cached --check`: PASS
**notes**:
- 작업트리에 남아 있는 `./go-backend` workspace 추가 변경은 이번 커밋에 포함하지 않았다.
---

##_2026-05-19_Codex_flutter_web_timing_governance_index_p38
**status**: task1_task5_integration_done
**plan**:
- P38 Flutter Web Timing Governance Index 계획을 추가하고, P30-P37 timing CI governance 문서/스크립트를 canonical entrypoint로 묶는다.
**changes**:
- docs/superpowers/plans/2026-05-19-flutter-web-timing-governance-index-p38.md
  - P38 Task 1~5 구현계획과 기존 정책/감사 체인 통합 보강 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_web_timing_governance_index_policy_test.sh
  - governance index marker, 필수 문서/스크립트 참조, 기존 정책 가드 체인, P27-P38 감사 기록 체인, workflow, matrix 연결을 검증한다.
- docs/ci/flutter-web-timing-governance-index.md
  - Flutter web timing CI governance의 canonical entrypoint, policy chain, guard chain, audit trail을 문서화했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter web timing governance index policy` step을 추가했다.
- scripts/ci_gate_matrix_policy_test.sh
  - governance index policy doc/guard/workflow 연결을 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web timing governance index policy` row와 governance marker를 추가했다.
- docs/audit/flutter-web-timing-governance-index-p38.md
  - P38 TDD 기록, 기존 작업 통합 보강, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_timing_governance_index_policy_test.sh`: FAIL, governance index missing
- RED: `bash scripts/ci_web_timing_governance_index_policy_test.sh`: FAIL, matrix marker missing
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, matrix marker missing
- GREEN: `bash scripts/ci_web_timing_governance_index_policy_test.sh`: PASS
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- RED integration supplement: `bash scripts/ci_web_timing_governance_index_policy_test.sh`: FAIL, guard_chain marker missing
- GREEN integration supplement: `bash scripts/ci_web_timing_governance_index_policy_test.sh`: PASS
- `bash scripts/ci_web_release_review_checklist_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_synthetic_review_example_test.sh`: PASS
- `bash -n` targeted governance index scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P38/P37 integration files: PASS
- trailing whitespace check on targeted P38/P37 integration files: PASS
**notes**:
- 실제 production timing artifact는 추가하지 않았고, governance entrypoint를 문서/guard로 고정했다.
**next_steps**:
- 다음 단계는 실제 GitHub Actions timing artifact가 준비되면 P39 production-sample intake checklist를 추가하는 것이다.
---

##_2026-05-19_Codex_flutter_web_timing_synthetic_review_example_p37
**status**: task1_task4_done
**plan**:
- P37 Flutter Web Timing Synthetic Review Example 계획을 추가하고, release-review checklist를 비운영 synthetic fixture/audit example로 재현 검증한다.
**changes**:
- docs/superpowers/plans/2026-05-19-flutter-web-timing-synthetic-review-example-p37.md
  - P37 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_web_timing_synthetic_review_example_test.sh
  - synthetic review fixture와 audit file marker를 확인하고 validator/wrapper 경로를 실행한다.
- docs/ci/fixtures/flutter-web-timing/2026-05-19-synthetic-review/
  - `manifest.env`, five synthetic sample files, `aggregate.env`를 추가했다.
- docs/audit/flutter-web-timing-synthetic-review-example.md
  - wrapper가 검증할 수 있는 synthetic audit example을 추가했다.
- scripts/ci_web_release_review_checklist_policy_test.sh
  - synthetic example script, fixture path, audit path를 요구하고 synthetic test를 실행하도록 확장했다.
- docs/ci/flutter-web-release-review-checklist.md
  - synthetic dry-run command와 fixture/audit example path를 문서화했다.
- docs/audit/flutter-web-timing-synthetic-review-example-p37.md
  - P37 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`: FAIL, synthetic review dir missing
- GREEN: `bash scripts/ci_web_timing_synthetic_review_example_test.sh`: PASS
- RED: `bash scripts/ci_web_release_review_checklist_policy_test.sh`: FAIL, synthetic checklist marker missing
- GREEN: `bash scripts/ci_web_release_review_checklist_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n` targeted synthetic review scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Synthetic review example is marked as non-production and has `decision: rejected`.
**next_steps**:
- 다음 단계는 P38 release-governance index page를 추가해 Flutter web timing CI policies를 한 entrypoint에서 찾을 수 있게 하는 것이다.
---

##_2026-05-19_Codex_flutter_web_release_review_checklist_p36
**status**: task1_task4_done
**plan**:
- P36 Flutter Web Release Review Checklist 계획을 추가하고, timing threshold/PR gate 정책 변경 전 실행해야 할 release-review command order를 문서와 CI guard로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-19-flutter-web-release-review-checklist-p36.md
  - P36 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_web_release_review_checklist_policy_test.sh
  - checklist 문서 marker, 필수 command, 참조 파일, workflow, matrix 연결을 검증한다.
- docs/ci/flutter-web-release-review-checklist.md
  - policy gates, artifact export, offline aggregation, fixture validation, threshold review check, audit decision 순서를 문서화했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter web release review checklist policy` step을 추가했다.
- scripts/ci_gate_matrix_policy_test.sh
  - checklist policy doc/guard/workflow 연결을 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web release review checklist policy` row와 checklist marker를 추가했다.
- docs/audit/flutter-web-release-review-checklist-p36.md
  - P36 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_release_review_checklist_policy_test.sh`: FAIL, checklist document missing
- RED: `bash scripts/ci_web_release_review_checklist_policy_test.sh`: FAIL, matrix marker missing
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, matrix marker missing
- GREEN: `bash scripts/ci_web_release_review_checklist_policy_test.sh`: PASS
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash -n` targeted release review checklist scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- 실제 timing artifact fixture는 추가하지 않았고, release-review command order만 문서/guard로 고정했다.
**next_steps**:
- 다음 단계는 P37 synthetic timing review example을 추가해 checklist를 실제 비운영 예제로 재현 가능하게 검증하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_threshold_change_review_check_p35
**status**: task1_task4_done
**plan**:
- P35 Flutter Web Timing Threshold Change Review Check 계획을 추가하고, threshold 변경 PR에서 reviewer가 validator와 audit template 값을 빠짐없이 실행/작성했는지 검증하는 wrapper command를 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-threshold-change-review-check-p35.md
  - P35 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/flutter_web_timing_threshold_change_review_check_test.sh
  - 임시 review fixture와 audit file의 성공 경로 및 `p95_seconds` mismatch 실패 경로를 검증한다.
- scripts/flutter_web_timing_threshold_change_review_check.sh
  - fixture validator를 실행하고 manifest/aggregate/audit Markdown field consistency를 검증한다.
- scripts/ci_web_timing_review_fixture_policy_test.sh
  - fixture policy가 wrapper script/test, required audit fields, template `validator_command` field를 문서화했는지 검증한다.
- docs/ci/flutter-web-timing-review-fixtures.md
  - threshold change review check command와 required audit fields를 문서화했다.
- docs/audit/templates/flutter-web-timing-review-template.md
  - wrapper가 읽을 수 있는 machine-check field block을 추가했다.
- docs/audit/flutter-web-timing-threshold-change-review-check-p35.md
  - P35 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`: FAIL, wrapper script missing
- GREEN: `bash scripts/flutter_web_timing_threshold_change_review_check_test.sh`: PASS
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: FAIL, wrapper marker missing
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n` targeted timing threshold review scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- 실제 timing review fixture는 추가하지 않았고, reviewer wrapper/check contract만 구현했다.
**next_steps**:
- 다음 단계는 P36 release-review checklist page를 추가해 timing gate command 실행 순서를 한 문서에서 안내하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_review_validator_p34
**status**: task1_task4_done
**plan**:
- P34 Flutter Web Timing Review Validator 계획을 추가하고, concrete review fixture directory를 threshold change PR 전에 검증하는 validator를 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-review-validator-p34.md
  - P34 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/flutter_web_timing_review_fixture_validate_test.sh
  - 임시 review fixture 성공 경로와 `no_phi_attestation=false` 실패 경로를 검증한다.
- scripts/flutter_web_timing_review_fixture_validate.sh
  - manifest, five samples, aggregate output, review status, no-PHI attestation, forbidden fields, aggregate reproducibility를 검증한다.
- scripts/ci_web_timing_review_fixture_policy_test.sh
  - fixture policy가 validator script/test, manifest required fields, forbidden fixture fields를 문서화했는지 검증한다.
- docs/ci/flutter-web-timing-review-fixtures.md
  - validator 사용법과 concrete review directory 검증 범위를 문서화했다.
- docs/audit/flutter-web-timing-review-validator-p34.md
  - P34 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`: FAIL, validator script missing
- GREEN: `bash scripts/flutter_web_timing_review_fixture_validate_test.sh`: PASS
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: FAIL, validator marker missing
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n` targeted timing review validator scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- 실제 timing fixture directory는 추가하지 않았고, validator와 policy contract만 구현했다.
**next_steps**:
- 다음 단계는 P35 threshold change PR checklist 또는 wrapper command를 추가해 reviewer가 validator와 audit template을 빠짐없이 실행하도록 하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_review_fixture_p33
**status**: task1_task4_done
**plan**:
- P33 Flutter Web Timing Review Fixture 계획을 추가하고, 실제 5-sample timing promotion review를 fixture 구조와 audit template으로 재현 가능하게 기록하는 규약을 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-review-fixture-p33.md
  - P33 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_web_timing_review_fixture_policy_test.sh
  - fixture convention 문서, fixture root README, audit template, promotion policy cross-reference, workflow, matrix 연결을 검증한다.
- docs/ci/flutter-web-timing-review-fixtures.md
  - fixture root, manifest, sample naming, aggregate output, GitHub Actions artifact source, no-PHI 규칙을 문서화했다.
- docs/ci/fixtures/flutter-web-timing/README.md
  - 실제 review directory 구조와 반입 제한을 기록했다.
- docs/audit/templates/flutter-web-timing-review-template.md
  - aggregate result, decision, sign-off, rollback notes를 기록하는 audit template을 추가했다.
- docs/ci/flutter-web-timing-threshold-promotion.md
  - review fixture policy, fixture root, audit template을 promotion policy에 교차 참조했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter web timing review fixture policy` step을 추가했다.
- scripts/ci_gate_matrix_policy_test.sh
  - review fixture policy doc/template/guard/workflow 연결을 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web timing review fixture policy` row와 fixture marker를 추가했다.
- docs/audit/flutter-web-timing-review-fixture-p33.md
  - P33 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: FAIL, fixture policy document missing
- RED: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: FAIL, matrix marker missing
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, matrix marker missing
- GREEN: `bash scripts/ci_web_timing_review_fixture_policy_test.sh`: PASS
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/flutter_web_timing_sample_aggregate_test.sh`: PASS
- `bash -n` targeted timing review fixture scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- 실제 timing artifact는 추가하지 않았고, fixture convention 및 audit template만 고정했다.
**next_steps**:
- 다음 단계는 P34 fixture manifest validator를 추가해 concrete review directory가 정책 변경 PR 전에 검증되도록 하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_sample_aggregator_p32
**status**: task1_task4_done
**plan**:
- P32 Flutter Web Timing Sample Aggregator 계획을 추가하고, exported timing artifact 5개 이상에서 median/p95/worst-case를 산출하는 offline aggregator를 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-sample-aggregator-p32.md
  - P32 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/flutter_web_timing_sample_aggregate_test.sh
  - 5개 fixture artifact의 sample count, branch types, median, p95, worst-case와 최소 sample 부족 실패를 검증한다.
- scripts/flutter_web_timing_sample_aggregate.sh
  - exported `.env` artifact directory를 읽어 sorted durations, median, nearest-rank p95, worst-case, branch types, runner contexts를 출력한다.
- scripts/ci_web_timing_threshold_promotion_policy_test.sh
  - promotion policy가 aggregator script/test와 aggregate output schema를 문서화했는지 검증한다.
- docs/ci/flutter-web-timing-threshold-promotion.md
  - offline aggregator 사용법과 aggregate output fields를 문서화했다.
- docs/audit/flutter-web-timing-sample-aggregator-p32.md
  - P32 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`: FAIL, aggregator script missing
- GREEN: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`: PASS
- RED: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: FAIL, aggregator marker missing
- GREEN: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n` targeted timing aggregator scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Aggregator is offline and deterministic; actual GitHub artifact download automation is intentionally outside this step.
**next_steps**:
- 다음 단계는 P33 real-sample review fixture convention과 audit template을 추가해 실제 5-sample promotion review를 재현 가능하게 기록하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_threshold_promotion_p31
**status**: task1_task4_done
**plan**:
- P31 Flutter Web Timing Threshold Promotion 계획을 추가하고, web release gate 완화 조건을 최소 artifact sample, measured cost review, 자동 완화 금지 정책으로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-threshold-promotion-p31.md
  - P31 Task 1~4 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_web_timing_threshold_promotion_policy_test.sh
  - promotion policy 문서, timing collection 교차 참조, workflow step, CI matrix 연결을 검증한다.
- docs/ci/flutter-web-timing-threshold-promotion.md
  - 최소 5개 성공 artifact sample, branch type, runner context, median/p95/worst-case review, 자동 완화 금지를 문서화했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter web timing threshold promotion policy` step을 추가했다.
- docs/ci/flutter-web-gate-timing-collection.md
  - threshold promotion policy 문서와 `flutter-web-timing-report` sample source를 교차 참조한다.
- scripts/ci_gate_matrix_policy_test.sh
  - promotion policy marker와 guard/doc/workflow 연결을 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web timing threshold promotion policy` row와 promotion marker를 추가했다.
- docs/audit/flutter-web-timing-threshold-promotion-p31.md
  - P31 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: FAIL, promotion policy document missing
- RED: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: FAIL, matrix marker missing
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, matrix marker missing
- GREEN: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`: PASS
- `bash -n` targeted timing promotion scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Gate relaxation remains manual and evidence-driven while `automated_gate_relaxation: false`.
**next_steps**:
- 다음 단계는 P32 offline timing artifact sample aggregator를 추가해 5개 이상 artifact export가 있을 때 median/p95/worst-case를 산출하는 것이다.
---

##_2026-05-14_Codex_flutter_web_timing_artifact_p30
**status**: task1_task5_done
**plan**:
- P30 Flutter Web Timing Artifact 계획을 추가하고, Flutter web release gate timing report를 GitHub Actions artifact로 생성/업로드하도록 CI와 정책 가드를 연결한다.
**changes**:
- docs/superpowers/plans/2026-05-14-flutter-web-timing-artifact-p30.md
  - P30 Task 1~5 구현계획과 RED-GREEN 검증 순서를 기록했다.
- scripts/ci_flutter_web_timing_artifact_workflow_test.sh
  - CI workflow가 `Flutter web timing report`, `Upload Flutter web timing report`, `actions/upload-artifact@v4`, artifact name/path를 포함하는지 검증한다.
- .github/workflows/ci.yml
  - `Flutter web release gate` 이후 timing report를 생성하고 `/tmp/manpasik_flutter_web_timing.env`를 `flutter-web-timing-report` artifact로 업로드한다.
- scripts/ci_web_gate_timing_collection_policy_test.sh
  - CI artifact upload 계약 marker와 workflow marker를 검증한다.
- docs/ci/flutter-web-gate-timing-collection.md
  - artifact name/path/action과 report 생성/업로드 step 이름을 정책 marker로 문서화했다.
- scripts/ci_gate_matrix_policy_test.sh
  - CI gate matrix가 timing artifact workflow guard를 참조하도록 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web timing artifact` row와 artifact marker를 추가했다.
- docs/audit/flutter-web-timing-artifact-p30.md
  - P30 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`: FAIL, workflow missing `Flutter web timing report`
- GREEN: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`: PASS
- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: FAIL, `ci_artifact_upload: true` marker missing
- GREEN: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, `Flutter web timing artifact` marker missing
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/flutter_web_timing_report_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n` targeted timing artifact scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Artifact upload is tied to successful `flutter-app` web release gate runs and keeps branch type plus runner context for measured cost review.
**next_steps**:
- 다음 단계는 P31 timing threshold promotion rule을 정의해 최소 5개 성공 sample 이후 gate 완화 조건을 명문화하는 것이다.
---

##_2026-05-13_Codex_flutter_web_timing_report_p29
**status**: task1_task3_done
**plan**:
- P29 Flutter Web Timing Report 계획을 추가하고, duration marker를 표준 key-value report artifact로 변환하는 스크립트와 정책 marker를 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-web-timing-report-p29.md
  - P29 Task 1~3 구현계획과 진행 상태를 기록했다.
- scripts/flutter_web_timing_report_test.sh
  - fixture log의 duration marker 2개를 report로 변환하고 marker 없는 log가 실패하는지 검증한다.
- scripts/flutter_web_timing_report.sh
  - duration marker를 추출해 `report_version`, `sample_count`, `duration_seconds_values`, latest/min/max duration, branch/runner context를 출력한다.
- scripts/ci_web_gate_timing_collection_policy_test.sh
  - report artifact format, report script/test, required fields marker를 검증한다.
- docs/ci/flutter-web-gate-timing-collection.md
  - `report_artifact_format: key_value_v1`와 report schema를 문서화했다.
- docs/audit/flutter-web-timing-report-p29.md
  - P29 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_timing_report_test.sh`: FAIL, report script missing
- GREEN: `bash scripts/flutter_web_timing_report_test.sh`: PASS
- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: FAIL, report artifact format marker missing
- GREEN: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash -n` targeted timing report scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Report format is `key_value_v1`; no JSON parser dependency is required.
**next_steps**:
- 다음 단계는 CI artifact upload 연결 또는 timing threshold 승격 조건을 더 엄격히 문서화하는 P30이다.
---

##_2026-05-13_Codex_flutter_web_timing_collection_policy_p28
**status**: task1_task3_done
**plan**:
- P28 Flutter Web Timing Collection Policy 계획을 추가하고, web release gate duration marker 수집 기준과 PR gate 완화 전 검토 조건을 문서/guard/CI matrix로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-web-timing-collection-policy-p28.md
  - P28 Task 1~3 구현계획과 진행 상태를 기록했다.
- scripts/ci_web_gate_timing_collection_policy_test.sh
  - timing collection policy marker, web release policy timing marker, workflow 연결을 검증한다.
- docs/ci/flutter-web-gate-timing-collection.md
  - duration marker 수집 query, 최소 5개 sample, advisory threshold, nightly split 검토 조건을 문서화했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter web timing collection policy` step을 추가했다.
- scripts/ci_gate_matrix_policy_test.sh
  - timing collection policy 문서/guard/matrix/workflow marker를 검증한다.
- docs/ci/ci-gate-matrix.md
  - `Flutter web timing collection policy` row와 sample/threshold marker를 추가했다.
- docs/audit/flutter-web-timing-collection-policy-p28.md
  - P28 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: FAIL, timing collection policy document missing
- GREEN: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, timing collection matrix marker missing
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n` targeted web timing policy scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Timing threshold remains advisory and `blocking_threshold_seconds` remains unset until measured samples exist.
**next_steps**:
- 다음 단계는 timing sample report/artifact 포맷 또는 threshold 승격 조건을 추가하는 P29다.
---

##_2026-05-13_Codex_flutter_web_release_gate_timing_p27
**status**: task1_task3_done
**plan**:
- P27 Flutter Web Release Gate Timing 계획을 추가하고, web release gate가 duration seconds marker를 출력하도록 구현한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-web-release-gate-timing-p27.md
  - P27 Task 1~3 구현계획과 진행 상태를 기록했다.
- scripts/flutter_web_release_gate_test.sh
  - policy-check 성공 fixture에서 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=42` marker 출력을 검증한다.
- scripts/flutter_web_release_gate.sh
  - 실제 build 경로에서 elapsed seconds를 측정하고 `FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS=<n>` marker를 출력한다.
  - duration env 값이 numeric이 아니면 실패한다.
- scripts/ci_web_gate_policy_test.sh
  - `timing_capture: true`와 `timing_marker: FLUTTER_WEB_RELEASE_GATE_DURATION_SECONDS` marker를 검증한다.
- scripts/ci_gate_matrix_policy_test.sh
  - CI gate matrix timing marker를 검증한다.
- docs/ci/flutter-web-release-gate-policy.md
  - web gate가 duration seconds marker를 출력해야 한다고 문서화했다.
- docs/ci/ci-gate-matrix.md
  - timing capture marker를 matrix marker block에 추가했다.
- docs/audit/flutter-web-release-gate-timing-p27.md
  - P27 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_release_gate_test.sh`: FAIL, duration marker missing
- GREEN: `bash scripts/flutter_web_release_gate_test.sh`: PASS
- RED: `bash scripts/ci_web_gate_policy_test.sh`: FAIL, timing capture marker missing
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, timing matrix marker missing
- GREEN: web policy and matrix guards PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n` targeted web timing scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- 실제 full web build는 이번 단계에서 재실행하지 않았고, parser/policy/common gate로 timing marker 구현을 검증했다.
**next_steps**:
- 다음 단계는 CI log에서 duration marker를 수집하는 절차 또는 timing threshold policy를 추가하는 P28이다.
---

##_2026-05-13_Codex_ci_gate_matrix_p26
**status**: task1_task2_done
**plan**:
- P26 CI Gate Matrix 계획을 추가하고, 릴리스/보안/evidence/Wasm 관련 gate의 blocking 여부, workflow 위치, policy 문서, guard script를 한 matrix로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-ci-gate-matrix-p26.md
  - P26 Task 1~2 구현계획과 진행 상태를 기록했다.
- scripts/ci_gate_matrix_policy_test.sh
  - matrix 문서 marker, 참조 policy/guard 파일 존재, workflow step 연결을 검증한다.
- docs/ci/ci-gate-matrix.md
  - SSOT, security, assay evidence, Flutter evidence, Flutter web release, Wasm readiness gate matrix를 추가했다.
- docs/audit/ci-gate-matrix-p26.md
  - P26 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`: FAIL, matrix document missing
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash -n` targeted CI policy scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Wasm readiness remains non-blocking while `wasm_release_target: false`.
- Flutter web release gate remains blocking while `release_target: js_web`.
**next_steps**:
- 다음 단계는 Flutter web release gate timing capture를 추가해 cost review 근거를 수치화하는 P27이다.
---

##_2026-05-13_Codex_flutter_web_gate_cost_policy_p25
**status**: task1_task2_done
**plan**:
- P25 Flutter Web Gate Cost Policy 계획을 추가하고, JS web release gate를 PR/release branch blocking으로 유지하는 비용 정책을 guard로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-web-gate-cost-policy-p25.md
  - P25 Task 1~2 구현계획과 진행 상태를 기록했다.
- scripts/ci_web_gate_policy_test.sh
  - `ci_execution_mode`, `release_branch_required`, `nightly_split`, `cost_review_required_before_relaxing`, `blocking_build_command` marker를 검증한다.
- docs/ci/flutter-web-release-gate-policy.md
  - JS web release artifact는 PR/release branch에서 계속 blocking gate로 유지하고, nightly 전환 전 비용 검토가 필요하다고 문서화했다.
- docs/audit/flutter-web-gate-cost-policy-p25.md
  - P25 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_gate_policy_test.sh`: FAIL, CI execution mode marker missing
- GREEN: `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash -n` targeted web policy scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- `nightly_split` remains false while JS web is the current release target.
**next_steps**:
- 다음 단계는 Flutter web release gate timing capture 또는 CI 전체 gate matrix 문서화다.
---

##_2026-05-13_Codex_flutter_evidence_ci_shard_policy_p24
**status**: task1_task2_done
**plan**:
- P24 Flutter Evidence CI Shard Policy 계획을 추가하고, P23에서 적용한 CI shard 운영 상태를 policy marker와 guard로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-evidence-ci-shard-policy-p24.md
  - P24 Task 1~2 구현계획과 진행 상태를 기록했다.
- scripts/flutter_evidence_ui_gate_policy_test.sh
  - `ci_execution_mode`, `ci_shard_steps`, `aggregate_ci_step`, `aggregate_local_gate_supported` marker를 검증한다.
  - policy의 CI shard category가 workflow `--category` command와 일치하는지 검증한다.
- docs/ci/flutter-evidence-ui-gate-policy.md
  - CI는 category shard step으로 실행하고 aggregate local gate는 개발자 검증용으로 유지한다는 정책을 문서화했다.
- docs/audit/flutter-evidence-ci-shard-policy-p24.md
  - P24 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: FAIL, CI shard policy marker missing
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `bash -n` targeted policy/shard scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- CI execution mode is now policy-guarded as `category_shards`.
**next_steps**:
- 다음 단계는 Flutter web release gate PR 필수 정책의 비용/분기 전략을 검토하는 P25다.
---

##_2026-05-13_Codex_ci_flutter_evidence_ui_shards_p23
**status**: task1_task2_done
**plan**:
- P23 CI Flutter Evidence UI Shards 계획을 추가하고, aggregate evidence gate CI step을 `contract`, `ui`, `data` shard step으로 분리한다.
**changes**:
- docs/superpowers/plans/2026-05-13-ci-flutter-evidence-ui-shards-p23.md
  - P23 Task 1~2 구현계획과 진행 상태를 기록했다.
- scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh
  - CI workflow가 세 category shard step과 `--category` 호출을 포함하고 aggregate-only command를 제거했는지 검증한다.
- .github/workflows/ci.yml
  - `Flutter evidence UI gate` aggregate step을 `contract`, `ui`, `data` 세 step으로 분리했다.
- docs/audit/ci-flutter-evidence-ui-shards-p23.md
  - P23 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: FAIL, shard workflow markers missing
- GREEN: `bash scripts/ci_flutter_evidence_ui_gate_shards_workflow_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash -n` targeted evidence CI scripts: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category contract`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category ui`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category data`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Aggregate-only evidence gate command was removed from CI to avoid duplicate execution cost.
**next_steps**:
- 다음 단계는 evidence policy 문서에 CI shard 운영 상태 marker를 추가하거나 Flutter web release gate 비용 정책을 PR/nightly로 조정할지 검토하는 P24다.
---

##_2026-05-13_Codex_flutter_evidence_ui_gate_shards_p22
**status**: task1_task3_done
**plan**:
- P22 Flutter Evidence UI Gate Shards 계획을 추가하고, P21에서 상한에 도달한 evidence gate를 `contract`, `ui`, `data` category shard로 분리한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-evidence-ui-gate-shards-p22.md
  - P22 Task 1~3 구현계획과 진행 상태를 기록했다.
- scripts/flutter_evidence_ui_gate_shard_test.sh
  - category 목록, shard별 count, 핵심 test 포함 여부, unknown category 실패를 검증한다.
- scripts/flutter_evidence_ui_gate.sh
  - `--list-categories`와 `--category contract|ui|data`를 추가했다.
  - `--list`/`--count`가 category option을 처리하도록 확장했다.
- scripts/flutter_evidence_ui_gate_policy_test.sh
  - `category_shards`와 `max_test_files_per_category` 정책 및 category별 count 상한을 검증한다.
- docs/ci/flutter-evidence-ui-gate-policy.md
  - `category_shards: contract,ui,data`와 `max_test_files_per_category: 5`를 추가했다.
- docs/audit/flutter-evidence-ui-gate-shards-p22.md
  - P22 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: FAIL, category support missing
- GREEN: `bash scripts/flutter_evidence_ui_gate_shard_test.sh`: PASS
- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: FAIL, category_shards marker missing
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash -n` targeted evidence gate scripts: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category contract`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category ui`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh --category data`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Category counts are contract=4, ui=4, data=5; all are within `max_test_files_per_category: 5`.
**next_steps**:
- 다음 단계는 CI에서 aggregate evidence gate를 category shard steps로 분리할지 결정하는 P23이다.
---

##_2026-05-13_Codex_flutter_evidence_ui_gate_policy_p21
**status**: task1_task2_done
**plan**:
- P21 Flutter Evidence UI Gate Policy 계획을 추가하고, evidence gate 범위/비용 상한/CI 정책 guard를 기계 검증으로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-evidence-ui-gate-policy-p21.md
  - P21 Task 1~2 구현계획과 진행 상태를 기록했다.
- scripts/flutter_evidence_ui_gate.sh
  - `--count` 모드를 추가해 gate test file 수를 출력한다.
- scripts/flutter_evidence_ui_gate_policy_test.sh
  - 정책 문서 marker, gate count 상한, CI policy step 연결을 검증한다.
- docs/ci/flutter-evidence-ui-gate-policy.md
  - `gate_scope: evidence_ui_contract`, `required_on_pull_request: true`, `duplicates_full_flutter_test: true`, `max_test_files: 13` 정책을 문서화했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Flutter evidence UI gate policy` step을 추가했다.
- docs/audit/flutter-evidence-ui-gate-policy-p21.md
  - P21 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: FAIL, policy document missing
- GREEN: `bash scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh scripts/flutter_evidence_ui_gate_policy_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Evidence UI gate는 현재 full Flutter test와 일부 중복되지만 named regulatory evidence regression gate로 유지한다.
**next_steps**:
- Evidence gate가 13개 test file을 넘으면 category shard 설계 또는 policy 상한 변경 근거를 먼저 작성한다.
---

##_2026-05-13_Codex_flutter_evidence_ui_gate_p20
**status**: task1_task3_done
**plan**:
- P20 Flutter Evidence UI Gate 계획을 완료하고, evidence contract/safe badge/REST mapper 회귀 테스트를 하나의 gate로 묶어 CI에서 실행한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-evidence-ui-gate-p20.md
  - P20 Task 1~3 체크리스트를 진행 상태에 맞춰 갱신했다.
- scripts/flutter_evidence_ui_gate_test.sh
  - `--list` 기반으로 핵심 evidence Flutter test 목록 포함 여부를 검증한다.
- scripts/flutter_evidence_ui_gate.sh
  - evidence 관련 Flutter test 묶음을 `flutter test --no-pub`으로 실행하고 `FLUTTER_EVIDENCE_UI_GATE_PASS` marker를 출력한다.
- scripts/ci_flutter_evidence_ui_gate_workflow_test.sh
  - CI workflow의 evidence gate step 이름과 script path를 검증한다.
- .github/workflows/ci.yml
  - Flutter job에 `Flutter evidence UI gate` step을 추가했다.
- docs/audit/flutter-evidence-ui-gate-p20.md
  - P20 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_evidence_ui_gate_test.sh`: FAIL, evidence gate script missing
- GREEN: `bash scripts/flutter_evidence_ui_gate_test.sh`: PASS
- RED: `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: FAIL, workflow step missing
- GREEN: `bash scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `bash scripts/flutter_evidence_ui_gate.sh`: PASS
- `bash -n scripts/flutter_evidence_ui_gate.sh scripts/flutter_evidence_ui_gate_test.sh scripts/ci_flutter_evidence_ui_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS
**notes**:
- Evidence UI gate는 Flutter 전체 테스트를 대체하지 않고, evidence 회귀 범위를 명시적으로 고정하는 추가 gate다.
**next_steps**:
- 다음 단계는 evidence gate 비용/분류 최적화 또는 DataHub 외 신규 화면의 evidence badge 회귀 범위 확대다.
---

##_2026-05-13_Codex_wasm_readiness_p19
**status**: task1_task2_done
**plan**:
- P19 Wasm Readiness 계획을 추가하고, Wasm을 아직 릴리스 타깃으로 보지 않는다는 상태와 dry-run blocker를 문서/guard로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-wasm-readiness-p19.md
  - P19 Task 1~2 구현계획과 검증 순서를 기록했다.
- scripts/ci_wasm_readiness_policy_test.sh
  - Wasm readiness 문서가 `wasm_release_target: false`와 현재 blocker package marker를 포함하는지 검증한다.
- docs/ci/flutter-wasm-readiness.md
  - 현재 release target은 JS web이고 Wasm은 non-blocking readiness 상태임을 기록했다.
  - dry-run blocker로 `flutter_secure_storage_web`, `share_plus`, `connectivity_plus`, `package:js`를 기록했다.
- docs/audit/wasm-readiness-p19.md
  - P19 TDD 기록, 결정 사항, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_wasm_readiness_policy_test.sh`: FAIL, readiness document missing
- GREEN: `bash scripts/ci_wasm_readiness_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash -n scripts/ci_wasm_readiness_policy_test.sh scripts/ci_web_gate_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
**notes**:
- Wasm은 아직 release target이 아니며 JS web release gate 정책을 유지한다.
**next_steps**:
- 다음 일반 강화 단계는 CI gate 비용/병렬화 최적화 또는 DataHub 외 화면의 evidence badge 회귀 확장이다.
---

##_2026-05-13_Codex_ci_web_gate_policy_p18
**status**: task1_task2_done
**plan**:
- P18 CI Web Gate Policy 계획을 추가하고, Flutter web release gate를 PR 필수로 유지한다는 정책을 문서와 guard test로 고정한다.
**changes**:
- docs/superpowers/plans/2026-05-13-ci-web-gate-policy-p18.md
  - P18 Task 1~2 구현계획과 검증 순서를 기록했다.
- scripts/ci_web_gate_policy_test.sh
  - `required_on_pull_request: true`, `wasm_dry_run_blocking: false`, `release_target: js_web` marker와 CI 연결을 검증한다.
- docs/ci/flutter-web-release-gate-policy.md
  - JS web artifact를 현재 release target으로 명시하고, Wasm dry-run warning은 non-blocking으로 추적한다고 기록했다.
- docs/audit/ci-web-gate-policy-p18.md
  - P18 TDD 기록, 결정 사항, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_web_gate_policy_test.sh`: FAIL, policy document missing
- GREEN: `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n` targeted CI/web gate scripts: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
**notes**:
- Flutter web gate는 PR 필수로 유지한다.
- Wasm dry-run warning은 현 단계에서 non-blocking이다.
**next_steps**:
- P19에서 Wasm strict compatibility readiness plan을 별도로 작성한다.
---

##_2026-05-13_Codex_ci_flutter_web_gate_p17
**status**: task1_task2_done
**plan**:
- P17 CI Flutter Web Gate 계획을 추가하고, P16 `flutter_web_release_gate.sh`를 GitHub Actions `flutter-app` job에 연결한다.
**changes**:
- docs/superpowers/plans/2026-05-13-ci-flutter-web-gate-p17.md
  - P17 Task 1~2 구현계획과 검증 순서를 기록했다.
- scripts/ci_flutter_web_gate_workflow_test.sh
  - CI workflow가 `Flutter web release gate` step과 `scripts/flutter_web_release_gate.sh` path를 포함하는지 검증한다.
- .github/workflows/ci.yml
  - Flutter job의 `Run tests` 이후 `Flutter web release gate` step을 추가했다.
- docs/audit/ci-flutter-web-gate-p17.md
  - P17 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/ci_flutter_web_gate_workflow_test.sh`: FAIL, workflow step missing
- GREEN: `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash scripts/flutter_web_release_gate_test.sh`: PASS
- `bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
**notes**:
- Flutter job working-directory가 `frontend/flutter-app`이므로 release gate는 `../../scripts/flutter_web_release_gate.sh`로 호출한다.
**next_steps**:
- P18에서 Flutter web gate를 PR 필수로 유지할지 nightly/release branch 전용으로 분리할지 CI 비용 정책을 결정한다.
---

##_2026-05-13_Codex_flutter_web_release_gate_p16
**status**: task1_task2_done
**plan**:
- P16 Flutter Web Release Gate 계획을 추가하고, JS web build 성공과 Wasm dry-run warning을 release gate에서 분리한다.
- 현재 릴리스 타깃은 `flutter build web --no-pub` JS web artifact로 두고, Wasm warning은 non-blocking marker로 추적한다.
**changes**:
- docs/superpowers/plans/2026-05-13-flutter-web-release-gate-p16.md
  - P16 Task 1~2 구현계획과 검증 순서를 기록했다.
- scripts/flutter_web_release_gate_test.sh
  - fixture log 기반 parser test를 추가했다.
  - build success marker와 Wasm warning marker를 검증하고, success marker 누락 로그는 실패하게 했다.
- scripts/flutter_web_release_gate.sh
  - 기본 모드에서 Flutter web build를 실행하고 build log policy를 검사한다.
  - `--policy-check <log>` 모드로 빠른 parser 검증을 지원한다.
  - `FLUTTER_WEB_RELEASE_GATE_PASS`와 `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN` marker를 출력한다.
- docs/audit/flutter-web-release-gate-p16.md
  - P16 TDD 기록, 실제 build gate 결과, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `bash scripts/flutter_web_release_gate_test.sh`: FAIL, release gate script missing
- GREEN: `bash scripts/flutter_web_release_gate_test.sh`: PASS
- `bash scripts/flutter_web_release_gate.sh`: PASS with `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`
- `flutter test --no-pub` targeted DataHub layout/evidence/domain/REST and shared badge tests: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh`: PASS
**notes**:
- Wasm은 아직 release blocking target이 아니며, 경고는 추적 marker로 남긴다.
**next_steps**:
- P17에서 `flutter_web_release_gate.sh`를 CI workflow에 연결할지 검토하고, Wasm strict target은 별도 platform compatibility 계획으로 분리한다.
---

##_2026-05-13_Codex_datahub_export_compat_p15
**status**: task1_task2_done
**plan**:
- P15 DataHub Export Compatibility 계획을 추가하고, export response의 `filePath`/`recordCount` camelCase fallback을 TDD로 보강한다.
**changes**:
- docs/superpowers/plans/2026-05-13-datahub-export-compat-p15.md
  - P15 Task 1~2 구현계획과 검증 순서를 기록했다.
- frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart
  - `/health-records/export/fhir` local server helper와 camelCase export response 보존 테스트를 추가했다.
- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart
  - `exportData`가 `file_path`/`filePath`/`fhir_json`/`fhirJson` 파일 경로와 `record_count`/`recordCount`를 모두 처리하게 했다.
- docs/audit/datahub-export-compat-p15.md
  - P15 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, camelCase `filePath` returned empty
- GREEN: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `flutter test --no-pub` targeted DataHub layout/evidence/domain/REST and shared badge tests: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P15 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only
- `git diff --check` targeted P15 tracked files: PASS
- P15 targeted trailing whitespace check: PASS
**notes**:
- DataHub REST compatibility는 history item, total count, export result까지 snake/camel fallback을 갖췄다.
**next_steps**:
- P16에서 Wasm dry-run compatibility warning의 지원 정책과 대응 범위를 별도 계획으로 분리한다.
---

##_2026-05-13_Codex_datahub_total_count_compat_p14
**status**: task1_task2_done
**plan**:
- P14 DataHub Total Count Compatibility 계획을 추가하고, `getTotalMeasurementCount`가 `total_count`와 legacy `totalCount`를 모두 읽는지 TDD로 보강한다.
**changes**:
- docs/superpowers/plans/2026-05-13-datahub-total-count-compat-p14.md
  - P14 Task 1~2 구현계획과 검증 순서를 기록했다.
- frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart
  - `totalCount` camelCase response에서 총 측정 수가 7로 보존되는지 검증했다.
- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart
  - `_intField` helper를 추가하고 `getTotalMeasurementCount`가 `total_count`/`totalCount`를 모두 읽게 했다.
- docs/audit/datahub-total-count-compat-p14.md
  - P14 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, `totalCount` returned 0
- GREEN: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `flutter test --no-pub` targeted DataHub layout/evidence/domain/REST and shared badge tests: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P14 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
**notes**:
- history item과 total count의 snake/camel compatibility 정책을 맞췄다.
**next_steps**:
- P15에서 DataHub export response fallback과 record count compatibility를 검토한다.
---

##_2026-05-13_Codex_datahub_layout_smoke_p13
**status**: validation_done
**plan**:
- P13 DataHub Layout Smoke 계획을 추가하고, 좁은 모바일 폭에서 evidence badge/trend header가 overflow 없이 표시되는지 widget smoke로 검증한다.
- 긴 metric label과 `research_only` evidence summary를 함께 주입해 레이아웃 압력을 높인다.
**changes**:
- docs/superpowers/plans/2026-05-13-datahub-layout-smoke-p13.md
  - P13 Task 1~2 구현계획과 검증 순서를 기록했다.
- frontend/flutter-app/test/features/data_hub/presentation/data_hub_layout_smoke_test.dart
  - 320px 모바일 폭에서 긴 metric label과 `연구용` 배지가 함께 표시되고 Flutter layout exception이 없는지 검증했다.
- docs/audit/datahub-layout-smoke-p13.md
  - P13 smoke 결과, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- Initial smoke: `flutter test --no-pub test/features/data_hub/presentation/data_hub_layout_smoke_test.dart`: PASS
- `flutter test --no-pub` targeted DataHub layout/evidence/domain/REST and shared badge tests: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P13 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only
**notes**:
- 예상했던 layout RED는 재현되지 않아 생산 코드 수정 없이 회귀 smoke만 추가했다.
**next_steps**:
- P14에서 DataHub `getTotalMeasurementCount`의 `totalCount` camelCase fallback을 TDD로 보강한다.
---

##_2026-05-13_Codex_datahub_timestamp_trend_p12
**status**: task1_task2_done
**plan**:
- P12 DataHub Timestamp Trend 계획을 추가하고, summary trend가 value-sort가 아니라 timestamp 순서로 계산되는지 RED 테스트로 고정한다.
- `getBiomarkerSummary`와 `getAllBiomarkerSummaries`가 같은 시간 순서 기준으로 falling/rising/stable을 계산하게 한다.
**changes**:
- docs/superpowers/plans/2026-05-13-datahub-timestamp-trend-p12.md
  - P12 Task 1~2 구현계획과 검증 순서를 기록했다.
- frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart
  - out-of-order history response에서 값이 시간 순서로 감소할 때 summary trend가 `falling`인지 검증했다.
- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart
  - trend 계산은 timestamp 정렬된 ordered values를 사용하고, min/max 통계는 별도 sorted values를 사용하게 분리했다.
  - `getAllBiomarkerSummaries`도 biomarker별 `TrendDataPoint` 목록을 timestamp 정렬해 latest/trend를 계산한다.
- docs/audit/datahub-timestamp-trend-p12.md
  - P12 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, decreasing timestamp sequence reported rising
- GREEN: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `flutter test --no-pub` targeted DataHub evidence/domain/REST and shared badge tests: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P10-P12 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter build web --no-pub`: PASS, WebAssembly dry-run compatibility warnings only
- `git diff --check` targeted P10-P12 tracked files: PASS
- P10-P12 targeted trailing whitespace check: PASS
**notes**:
- latest evidence/value와 trend가 모두 같은 timestamp 정렬 흐름을 공유한다.
- trend 임계값 5% 자체의 정책 타당성은 다음 별도 단계로 남겼다.
**next_steps**:
- P13에서 DataHub badge/trend layout overflow smoke 또는 trend threshold policy audit을 진행한다.
---

##_2026-05-13_Codex_datahub_evidence_ui_and_latest_p10_p11
**status**: task1_task2_task3_done
**plan**:
- P10/P11 DataHub Evidence UI And Latest Selection 계획을 추가하고, 화면 배지 표시와 timestamp 최신 evidence 선택을 순차 구현한다.
- DataHub 화면은 P7 `MeasurementEvidenceBadge`를 재사용해 새 의료 판정 문구 없이 evidence 상태를 표시한다.
- DataHub REST summary는 응답 순서가 아니라 `measured_at` 기준 최신 measurement의 evidence metadata를 선택한다.
**changes**:
- docs/superpowers/plans/2026-05-13-datahub-evidence-ui-and-latest-p10-p11.md
  - P10/P11 Task 1~3 구현계획과 검증 순서를 기록했다.
- frontend/flutter-app/test/features/data_hub/presentation/data_hub_evidence_badge_test.dart
  - `research_only` DataHub summary가 `연구용` 배지로 표시되고 금지 문구가 나오지 않는지 검증했다.
- frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart
  - `_HeroChartCard`와 `_DetailPanel`에 `MeasurementEvidenceBadge`를 연결했다.
- frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart
  - out-of-order history response에서 trend timestamp 정렬과 summary latest evidence 선택을 검증했다.
- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart
  - REST history item을 `TrendDataPoint`로 표준화하고 timestamp 기준 정렬/최신 evidence 선택을 적용했다.
- docs/audit/datahub-evidence-ui-and-latest-p10-p11.md
  - P10/P11 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`: FAIL, DataHub evidence badge not rendered
- GREEN: `flutter test --no-pub test/features/data_hub/presentation/data_hub_evidence_badge_test.dart`: PASS
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, timestamp latest evidence selection missing
- GREEN: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `flutter test --no-pub` targeted DataHub evidence/domain/REST and shared badge tests: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P10/P11 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P10/P11 tracked files: PASS
- P10/P11 targeted trailing whitespace check: PASS
**notes**:
- 배지 문구는 기존 `MeasurementEvidencePresentation` 경로만 사용한다.
- `getAllBiomarkerSummaries.latestValue`도 timestamp 최신 point 값으로 정렬했다.
**next_steps**:
- P12에서 DataHub trend 계산을 value-sort 기반이 아닌 timestamp 기반 산식으로 분리해 TDD로 강화한다.
---

##_2026-05-12_Codex_datahub_evidence_metadata_p9
**status**: task1_task2_task3_done
**plan**:
- P9 DataHub Evidence Metadata 계획을 추가하고, DataHub trend/summary 모델이 measurement evidence metadata를 보존하는지 RED tests로 고정한다.
- `TrendDataPoint`와 `BiomarkerSummary`에 보수적 evidence 기본값을 추가한다.
- DataHub REST repository가 history response의 snake_case/camelCase evidence fields를 trend/summary로 전달하게 한다.
**changes**:
- docs/superpowers/plans/2026-05-12-datahub-evidence-metadata-p9.md
  - P9 Task 1~3 구현계획과 완료 체크를 기록했다.
- frontend/flutter-app/lib/features/data_hub/domain/data_hub_repository.dart
  - `TrendDataPoint`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - `BiomarkerSummary`에 `latestEvidenceStatus`, `latestDiagnosticReady`, `latestEvidenceGaps`를 추가했다.
- frontend/flutter-app/test/features/data_hub/domain/data_hub_domain_test.dart
  - DataHub domain evidence 기본값과 explicit metadata 생성 계약을 검증했다.
- frontend/flutter-app/lib/features/data_hub/data/data_hub_repository_rest.dart
  - REST history snake_case/camelCase helper를 추가했다.
  - `getTrendData`, `getBiomarkerSummary`, `getAllBiomarkerSummaries`가 evidence metadata를 보존하게 했다.
- frontend/flutter-app/test/features/data_hub/data/data_hub_repository_rest_test.dart
  - snake_case/camelCase history evidence mapping과 summary latest evidence mapping을 검증했다.
- docs/audit/datahub-evidence-metadata-p9.md
  - P9 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart`: FAIL, DataHub evidence fields missing
- GREEN: `flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart`: PASS
- RED: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: FAIL, REST evidence mapping missing
- GREEN: `flutter test --no-pub test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter test --no-pub test/features/data_hub/domain/data_hub_domain_test.dart test/features/data_hub/data/data_hub_repository_rest_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P9 Flutter files: PASS
- `git diff --check` targeted P9 tracked files: PASS
- P9 targeted trailing whitespace check: PASS
**notes**:
- DataHub UI 표시는 아직 하지 않았고 모델/REST mapping 표면만 확장했다.
- summary latest evidence는 현재 history 응답 순서의 첫 항목을 기준으로 보존한다.
**next_steps**:
- P10에서 DataHub UI badge 표시 또는 measured_at 기반 latest evidence 정렬 계약을 별도 TDD로 진행한다.
---

##_2026-05-12_Codex_measurement_evidence_home_badge_p8
**status**: task1_task2_done
**plan**:
- P8 Measurement Evidence Home Badge 계획을 추가하고, HomeScreen이 latest measurement evidence를 표시하는지 RED screen widget test로 고정한다.
- Home hero card가 P7 `MeasurementEvidenceBadge`를 재사용해 `research_only`를 `연구용`으로 표시하게 한다.
- DataHub 확장은 별도 단계로 분리한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-home-badge-p8.md
  - P8 Task 1~2 구현계획과 완료 체크를 기록했다.
- frontend/flutter-app/test/features/home/presentation/home_measurement_evidence_badge_test.dart
  - `homeDashboardProvider` override로 `research_only` latest measurement를 주입하고 `연구용` 배지를 검증했다.
- frontend/flutter-app/lib/features/home/presentation/home_screen.dart
  - `_HeroBentoCard` 최근 측정 정보 아래에 `MeasurementEvidenceBadge`를 표시한다.
- docs/audit/measurement-evidence-home-badge-p8.md
  - P8 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart`: FAIL, HomeScreen badge not rendered
- GREEN: `flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter test --no-pub test/features/home/presentation/home_measurement_evidence_badge_test.dart test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P8 Flutter files: PASS
**notes**:
- Home 화면도 P4/P7과 같은 helper/widget 경로를 사용한다.
- DataHub trend/summary evidence 표시는 아직 하지 않았다.
**next_steps**:
- P9에서 DataHub trend/summary evidence metadata 확장 여부를 결정하고 별도 TDD로 진행한다.
---

##_2026-05-12_Codex_measurement_evidence_history_badge_p7
**status**: task1_task2_task3_done
**plan**:
- P7 Measurement Evidence History Badge 계획을 추가하고, P4 `MeasurementEvidencePresentation`을 재사용하는 compact badge부터 RED widget test로 고정한다.
- `MeasurementResultScreen` 최신 측정 카드가 P6 history item evidence fields를 안전한 배지로 표시하게 한다.
- 의료 판정 문구와 DataHub 확장은 별도 단계로 분리한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-history-badge-p7.md
  - P7 Task 1~3 구현계획과 완료 체크를 기록했다.
- frontend/flutter-app/lib/features/measurement/presentation/widgets/measurement_evidence_badge.dart
  - `MeasurementEvidencePresentation`을 사용하는 compact evidence badge 위젯을 추가했다.
- frontend/flutter-app/test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart
  - `research_only`가 `연구용`으로 표시되고 "정상", "위험", "진단", "확정" 표현이 나오지 않는지 검증했다.
- frontend/flutter-app/lib/features/measurement/presentation/measurement_result_screen.dart
  - 최신 측정 카드에 `MeasurementEvidenceBadge`를 연결했다.
- docs/audit/measurement-evidence-history-badge-p7.md
  - P7 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: FAIL, widget missing
- GREEN: `flutter test --no-pub test/features/measurement/presentation/widgets/measurement_evidence_badge_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P7 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P7 tracked files: PASS
- P7 targeted trailing whitespace check: PASS
**notes**:
- 배지 문구는 P4 helper에서만 생성되며 화면에서 별도 판정 문구를 만들지 않았다.
- DataHub/home dashboard 표시는 아직 하지 않았다.
**next_steps**:
- P8에서 DataHub/home dashboard 또는 full history screen smoke로 evidence badge 표시를 확장한다.
---

##_2026-05-12_Codex_measurement_evidence_history_app_surface_p6
**status**: task1_task2_task3_done
**plan**:
- P6 Measurement Evidence History App Surface 계획을 추가하고, gateway REST history와 Flutter history mappers가 P5 summary evidence fields를 보존하는지 RED 테스트로 먼저 고정한다.
- Gateway `/measurements/history` REST response가 evidence fields를 포함하게 mock contract를 보강한다.
- Flutter history domain item, native gRPC repository, REST repository가 `research_only`, `diagnostic_ready=false`, `clinical_lock_required` gap을 보존하게 한다.
- 감사 문서와 공유 컨텍스트를 갱신하고 SSOT/security/assay/Go/Flutter 게이트를 통과한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-history-app-surface-p6.md
  - P6 Task 1~3 구현계획과 완료 체크를 기록했다.
- backend/services/gateway/internal/handler/e2e_test.go
  - `TestE2E_GetMeasurementHistory`가 `evidence_status`, `diagnostic_ready`, `evidence_gaps`, `research_only`를 요구하게 했다.
  - mock `GetMeasurementHistory`가 evidence fields를 포함한 summary 1건을 반환하게 했다.
- frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart
  - `MeasurementHistoryItem`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - native gRPC `MeasurementSummary` evidence fields를 도메인 history item으로 매핑했다.
  - 테스트 주입용 optional `MeasurementHistoryCall` override를 추가했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart
  - REST history response의 snake_case/camelCase evidence fields를 모두 decode한다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart
  - native history evidence mapping contract를 추가했다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_rest_test.dart
  - REST snake_case/camelCase history evidence mapping contract를 추가했다.
- docs/audit/measurement-evidence-history-app-surface-p6.md
  - P6 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory`: FAIL, history evidence fields missing
- GREEN: `go test -count=1 ./services/gateway/internal/handler -run TestE2E_GetMeasurementHistory`: PASS
- RED: `flutter test --no-pub test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: FAIL, history domain fields and native override missing
- GREEN: `flutter test --no-pub test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P6 Flutter files: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P6 tracked files: PASS
- P6 targeted trailing whitespace check: PASS
**notes**:
- history item 기본값은 `unknown`, `false`, `[]`로 하위 호환 경로도 보수적으로 동작한다.
- UI 표시는 이번 범위에 넣지 않았고, 표시가 필요하면 P4 presentation helper를 재사용한다.
**next_steps**:
- P7에서 history 화면/data hub 표시 계층 또는 실제 gateway+measurement-service 통합 smoke로 확장한다.
---

##_2026-05-12_Codex_measurement_evidence_persistence_p5
**status**: task1_task2_task3_task4_done
**plan**:
- P5 Measurement Evidence Persistence 계획을 추가하고, history summary와 저장 스키마에서 evidence fields가 사라지는 문제를 RED 테스트로 먼저 고정한다.
- `MeasurementSummary` proto를 append-only로 확장하고 Go/Dart generated output을 재생성한다.
- measurement-service summary, gRPC history response, Postgres/Timescale init schema와 repository Store/GetHistory를 같은 evidence contract로 정렬한다.
- Flutter generated contract와 proto compile gate까지 확인한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-persistence-p5.md
  - P5 Task 1~4 구현계획과 완료 체크를 기록했다.
- backend/shared/proto/manpasik.proto
  - `MeasurementSummary`에 `evidence_status = 6`, `diagnostic_ready = 7`, `evidence_gaps = 8`을 추가했다.
- backend/shared/gen/go/v1/manpasik.pb.go
- backend/shared/gen/go/v1/manpasik_grpc.pb.go
- frontend/flutter-app/lib/generated/manpasik.pb.dart
- frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart
- frontend/flutter-app/lib/generated/manpasik.pbenum.dart
- frontend/flutter-app/lib/generated/manpasik.pbjson.dart
  - proto 변경에 맞춰 generated output을 재생성했다.
- backend/services/measurement-service/internal/service/measurement.go
  - `MeasurementSummary`에 evidence fields를 추가했다.
- backend/services/measurement-service/internal/handler/grpc.go
  - `GetMeasurementHistory` response summary에 evidence fields를 매핑했다.
- backend/services/measurement-service/internal/service/measurement_test.go
  - service history가 evidence fields를 보존하는지 검증했다.
- backend/services/measurement-service/internal/handler/grpc_stream_test.go
  - gRPC history response가 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 노출하는지 검증했다.
- backend/services/measurement-service/internal/repository/postgres/measurement.go
  - Store INSERT와 GetHistory SELECT/SCAN에 evidence fields를 추가했다.
- backend/services/measurement-service/internal/repository/postgres/measurement_schema_test.go
  - measurement init schema와 summary view가 evidence columns를 포함하는지 검증했다.
- infrastructure/database/init/04-measurement.sql
  - `measurement_data`에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`를 추가하고 summary view에 노출했다.
- frontend/flutter-app/test/generated/measurement_summary_evidence_contract_test.dart
  - Dart generated `MeasurementSummary` evidence 직렬화/역직렬화 계약을 추가했다.
- docs/audit/measurement-evidence-persistence-p5.md
  - P5 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `go test -count=1 ./services/measurement-service/internal/handler ./services/measurement-service/internal/service ./services/measurement-service/internal/repository/postgres`: FAIL, summary fields/schema missing
- RED: `flutter test --no-pub test/generated/measurement_summary_evidence_contract_test.dart`: FAIL, generated Dart summary evidence fields missing
- GREEN: `go test -count=1 ./services/measurement-service/internal/handler ./services/measurement-service/internal/service ./services/measurement-service/internal/repository/postgres`: PASS
- GREEN: `flutter test --no-pub test/generated/measurement_summary_evidence_contract_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted generated/test files: PASS
- `bash scripts/check_proto_generation_compile_gate.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted P5 tracked files: PASS
- P5 targeted trailing whitespace check: PASS
**notes**:
- proto field는 append-only라 기존 `MeasurementSummary` field 1~5를 변경하지 않았다.
- 기존 DB 데이터나 누락 status는 `unknown`, `DiagnosticReady=false`, empty gaps로 보수적으로 처리한다.
**next_steps**:
- P6에서 Flutter history domain/REST mapper가 summary evidence fields를 읽고, gateway REST history가 보존하는지 계약 테스트로 확장한다.
---

##_2026-05-12_Codex_measurement_evidence_safe_ui_p4
**status**: task1_task2_task3_task4_done
**plan**:
- P4 Measurement Evidence Safe UI 계획을 추가하고, evidence status를 UI copy로 변환하는 작은 domain helper부터 TDD로 잠근다.
- measurement golden path snapshot이 repository evidence fields를 전달하게 해 화면 표시 계층까지 데이터가 끊기지 않게 한다.
- MeasureScreen에는 상세 의료 문구가 아닌 compact badge label만 연결한다.
- 감사 문서와 공유 컨텍스트를 갱신하고 SSOT/security/assay/Flutter/Go 게이트를 통과한 뒤 종료한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-safe-ui-p4.md
  - P4 Task 1~4 구현계획과 stage gate 완료 체크를 기록했다.
- frontend/flutter-app/lib/features/measurement/domain/measurement_evidence_presentation.dart
  - `MeasurementEvidencePresentation` helper를 추가해 `research_only`, `clinical_locked`, `analytical_locked`, `unknown` 상태를 안전한 badge/detail copy로 변환한다.
- frontend/flutter-app/test/features/measurement/domain/measurement_evidence_presentation_test.dart
  - `research_only` copy가 "정상", "위험", "확정" 표현을 포함하지 않는지 검증했다.
- frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart
  - `MeasurementGoldenPathSnapshot`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - `serverProcessed`, `sessionEnded` snapshot이 서버 처리 결과의 evidence fields를 보존하게 했다.
- frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
  - 골든패스 snapshot이 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 전달하는지 검증했다.
- frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart
  - 서버 처리/세션 종료 상태에서 evidence status가 있으면 compact badge label을 표시하고, 없으면 기존 `SERVER`/`DONE` fallback을 유지한다.
- docs/audit/measurement-evidence-safe-ui-p4.md
  - P4 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart`: FAIL, helper missing
- GREEN: `flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart`: PASS
- RED: `flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: FAIL, snapshot evidence fields missing
- GREEN: `flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `flutter test --no-pub test/features/measurement/domain/measurement_evidence_presentation_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P4 Flutter files: PASS
- `go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: PASS
- `git diff --check` targeted tracked P4 files: PASS
- P4 targeted trailing whitespace check: PASS
**notes**:
- `research_only`는 `연구용` 배지로만 표시하며 진단/정상/위험/확정 표현으로 바꾸지 않았다.
- UI 상세 설명과 저장소 영속화는 아직 하지 않았다.
**next_steps**:
- P5에서는 evidence fields 저장 계층을 Timescale/Postgres migration과 history response contract부터 작성해 확장한다.
---

##_2026-05-12_Codex_measurement_evidence_app_surface_p3
**status**: task1_task2_task3_done
**plan**:
- P3 Measurement Evidence App Surface 계획을 추가한다.
- Gateway REST `/measurements/process` 응답이 gRPC evidence fields를 보존하는지 route contract로 잠근다.
- Flutter `ProcessMeasurementResult` 도메인 모델과 native/REST mapper가 evidence fields를 보존하게 한다.
- UI 표시와 저장소 영속화는 별도 단계로 분리한다.
**changes**:
- docs/superpowers/plans/2026-05-12-measurement-evidence-app-surface-p3.md
  - P3 Task 1~4 구현계획과 stage gate rules를 추가했다.
- backend/services/gateway/internal/handler/e2e_test.go
  - measurement process REST 응답에 `evidence_status`, `diagnostic_ready`, `evidence_gaps`가 포함되는지 검증했다.
  - mock stream response를 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap으로 갱신했다.
- frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart
  - `ProcessMeasurementResult`에 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 추가했다.
  - legacy 호환 기본값은 `unknown`, `false`, `[]`로 설정했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - native gRPC `MeasurementResult` evidence fields를 도메인 결과로 매핑했다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart
  - native repository가 evidence fields를 보존하는지 검증했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart
  - REST response의 snake_case/camelCase evidence keys를 모두 읽게 했다.
- frontend/flutter-app/test/features/measurement/data/measurement_process_gateway_mapper_test.dart
  - gateway mapper evidence decode 계약을 검증했다.
- docs/audit/measurement-evidence-app-surface-p3.md
  - P3 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
**quality_gate**:
- RED: `go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: FAIL, `research_only` missing
- GREEN: `go test -count=1 ./services/gateway/internal/handler -run TestE2E_ProcessMeasurementGoldenPath`: PASS
- RED: `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart`: FAIL, domain evidence fields missing
- GREEN: `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- RED: `flutter test --no-pub test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: FAIL, REST mapper evidence fields missing
- GREEN: `flutter test --no-pub test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_process_gateway_mapper_test.dart`: PASS
- `flutter analyze --no-pub --no-fatal-infos` targeted P3 Flutter files: PASS
- `git diff --check` targeted files: PASS
- P3 untracked file trailing whitespace check: PASS
**notes**:
- `research_only`는 진단 가능 표현으로 변환하지 않고 원문 상태로 전달한다.
- UI 표시와 DB 영속화는 아직 하지 않았다.
**next_steps**:
- P4에서 UI 표시 문구를 규정 친화적으로 추가하거나, 별도 persistence 단계에서 schema migration과 repository contract를 먼저 작성한다.
---

##_2026-05-11_Codex_measurement_evidence_contract_p2
**status**: task1_task2_done
**plan**:
- P2 Measurement Evidence Contract 계획을 추가하고, Task 1은 Go gRPC `MeasurementResult` 계약 확장으로 한정한다.
- RED-GREEN 순서로 handler response가 evidence status, diagnostic readiness, evidence gaps를 노출하는지 검증한다.
- Task 2에서는 Dart proto compile gate와 checked-in generated Dart output을 같은 evidence contract로 정렬한다.
**changes**:
- docs/superpowers/plans/2026-05-11-measurement-evidence-contract-p2.md
  - P2 Task 1~3 구현계획과 stage gate rules를 추가했다.
- backend/shared/proto/manpasik.proto
  - `MeasurementResult`에 `evidence_status = 7`, `diagnostic_ready = 8`, `evidence_gaps = 9`를 추가했다.
- backend/shared/gen/go/v1/manpasik.pb.go
  - Go generated `MeasurementResult` struct/getters/raw descriptor를 proto와 맞게 재생성했다.
- backend/shared/gen/go/v1/manpasik_grpc.pb.go
  - 동일 protoc 실행 결과를 반영했다.
- backend/services/measurement-service/internal/handler/grpc.go
  - `ProcessedResult` evidence fields를 gRPC `MeasurementResult` 응답으로 매핑했다.
- backend/services/measurement-service/internal/handler/grpc_stream_test.go
  - stream response가 `research_only`, `DiagnosticReady=false`, `clinical_lock_required` gap을 노출하는지 검증했다.
- docs/audit/measurement-evidence-contract-p2.md
  - Task 1/2 TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
- frontend/flutter-app/test/generated/measurement_result_evidence_contract_test.dart
  - checked-in Dart `MeasurementResult`가 evidence fields를 직렬화/역직렬화하는지 검증했다.
- frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh
  - 임시 generated Dart compile smoke가 `evidenceStatus`, `diagnosticReady`, `evidenceGaps`를 사용하게 했다.
- frontend/flutter-app/lib/generated/manpasik.pb.dart
- frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart
- frontend/flutter-app/lib/generated/manpasik.pbenum.dart
- frontend/flutter-app/lib/generated/manpasik.pbjson.dart
  - `bash scripts/generate_proto.sh`로 checked-in Dart generated output을 proto field 7~9와 맞게 재생성했다.
**quality_gate**:
- RED: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/handler -run TestStreamMeasurementStoresMeasurementAndFingerprint`: FAIL, generated `MeasurementResult` evidence fields missing
- GREEN: `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/handler -run TestStreamMeasurementStoresMeasurementAndFingerprint`: PASS
- RED: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/generated/measurement_result_evidence_contract_test.dart`: FAIL, checked-in Dart generated evidence fields missing
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
- `bash scripts/check_proto_generation_compile_gate.sh`: PASS, `PROTO_COMPILE_GATE_STATUS=passed`
- GREEN: `/mnt/d/우리집/flutter_cache/flutter/bin/flutter test --no-pub test/generated/measurement_result_evidence_contract_test.dart`: PASS
- `/mnt/d/우리집/flutter_cache/flutter/bin/flutter analyze --no-pub --no-fatal-infos test/generated/measurement_result_evidence_contract_test.dart lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/generated/manpasik.pbenum.dart lib/generated/manpasik.pbjson.dart`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `git diff --check` targeted files: PASS
- untracked P1/P2 file trailing whitespace check: PASS
**notes**:
- 새 field는 기존 `MeasurementResult` field 1~6을 건드리지 않는 backward-compatible 확장이다.
- 현재 실제 assay는 `research_only`로 응답되며 `DiagnosticReady=false`다.
**next_steps**:
- P3에서 REST/gateway 응답 또는 Flutter UI 표시 계층으로 evidence status를 확장한다.
- 저장 계층 영속화는 schema migration과 repository contract test를 먼저 작성한 뒤 진행한다.
---

##_2026-05-11_Codex_assay_evidence_factory_p1_task1
**status**: task1_task2_task3_done
**plan**:
- P1 Assay Evidence Factory를 단계별 구현계획으로 분리한다.
- Task 1에서는 실제 임상 성능을 주장하지 않고, assay evidence manifest와 diagnostic readiness gate만 구현한다.
- RED-GREEN-REVIEW 방식으로 evidence field 부재, clinical acceptance threshold 누락을 테스트로 먼저 잡는다.
- Task 3에서는 임상 잠금 assay의 evidence gaps를 CI에서 차단하는 좁은 gate를 추가한다.
**changes**:
- docs/superpowers/plans/2026-05-11-assay-evidence-factory-p1.md
  - Stage gate rules, Task 1~3 순차 구현계획, 다음 단계 지침을 추가했다.
- backend/shared/assay/registry.go
  - `EvidenceStatus`, `EvidenceManifest`, `AnalyticalPerformance`, `AcceptanceCriteria`를 추가했다.
  - 기본 건강 바이오마커 registry에 LOINC, UCUM, 필수 comparator method 힌트를 연결했다.
  - 실제 만파식 assay는 모두 `research_only`로 유지해 `IsDiagnosticReady()`가 false가 되게 했다.
  - clinical locked definition의 evidence gap 검증을 추가했다.
- backend/shared/assay/registry_test.go
  - glucose evidence manifest, research-only diagnostic gate, complete clinical lock, sensitivity/specificity threshold 누락 검증을 추가했다.
- docs/audit/assay-evidence-factory-p1.md
  - TDD 기록, 자체 코드리뷰, 품질 게이트, 다음 단계 지침을 문서화했다.
- backend/services/measurement-service/internal/service/measurement.go
  - 내부 `MeasurementData`와 `ProcessedResult`에 `EvidenceStatus`, `DiagnosticReady`, `EvidenceGaps`를 추가했다.
  - `applyAssaySemantics()`가 registry evidence status와 diagnostic readiness를 채우도록 연결했다.
- backend/services/measurement-service/internal/service/measurement_test.go
  - research-only assay 처리 결과가 `research_only`, `DiagnosticReady=false`, evidence gaps 존재를 노출하는지 검증했다.
- backend/shared/assay/registry.go
  - `ClinicalLockEvidenceGaps()`를 추가해 clinical-locked assay만 evidence gate 대상으로 수집한다.
- backend/shared/assay/registry_test.go
  - broken clinical lock fixture가 gate에 잡히는지, 현재 registry의 research-only assay들이 gate 실패 대상이 아닌지 검증했다.
- scripts/assay_evidence_gate.sh
  - CI에서 실행 가능한 assay evidence gate script를 추가했다.
- .github/workflows/ci.yml
  - `ssot-governance` job에 `Assay evidence gate` step을 추가했다.
**quality_gate**:
- RED 1: `go test -count=1 ./shared/assay`가 `Evidence`, `IsDiagnosticReady`, `EvidenceGaps` 부재로 FAIL
- RED 2: sensitivity/specificity threshold 누락 시 diagnostic-ready가 되면 안 된다는 테스트가 FAIL
- RED 3: measurement evidence fields 부재로 service target test가 build fail
- RED 4: `ClinicalLockEvidenceGaps` 부재로 `go test -count=1 ./shared/assay -run TestEvidenceGate`가 FAIL
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -run TestResolveIncludesEvidenceManifestWithoutDiagnosticClaim -count=1 ./shared/assay`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./services/measurement-service/internal/service -run TestProcessMeasurement_연구용_EvidenceStatus_노출`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay -run TestEvidenceGate`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- untracked P1 file trailing whitespace check: PASS
**notes**:
- LOINC/UCUM 연결은 식별자 정렬용이며 임상 성능 승인 주장이 아니다.
- `ReferenceMethod`는 clinical lock용 실제 증거 필드, `RequiredReference`는 research-only 상태에서 요구되는 comparator 힌트다.
**next_steps**:
- P2에서 measurement evidence surface를 저장 계층과 proto/gRPC/REST 외부 계약으로 확장한다.
- proto/gRPC 외부 계약 확장은 generated Go/Dart 재생성 게이트가 필요한 별도 단계로 진행한다.
---

##_2026-05-11_Codex_ssot_assay_semantics_p0
**status**: done
**plan**:
- 만파식 범용분석시스템 실현 가능성 강화를 위한 P0 실행계획을 수립한다.
- 최신 기술 스펙의 `alpha=0.98`, `MAX_CHANNELS=1792` 기준을 코드와 AGENTS 문서에 고정한다.
- CI에서 SSOT 상수 불일치를 먼저 차단하고 Rust clippy를 blocking gate로 전환한다.
- 측정 gRPC handler의 `primary_value/unit/confidence` 하드코딩을 제거하고 assay registry 기반 의미론으로 대체한다.
- Rust AI 추론은 production 모드에서 시뮬레이션 폴백을 금지한다.
- production profile 보안 설정과 runtime SSOT 상수를 release gate에서 검증한다.
**changes**:
- docs/superpowers/plans/2026-05-11-manpasik-feasibility-hardening-p0.md
  - P0 순차 실행계획을 추가하고 Task 1, Task 2 완료 상태를 반영했다.
- scripts/validate_ssot_constants.py
  - `ManPaSik_Tech_Spec_v2.4.3.md`를 기준으로 Rust core와 `AGENTS.md`의 핵심 상수 일치 여부를 검증한다.
- AGENTS.md
  - 차동측정 기본 alpha를 `0.98`로 정렬했다.
  - fingerprint 확장 경로를 `88 -> 448 -> 896 -> 1792`로 정렬했다.
  - `FingerprintClassifier` 입력/출력을 `1792 -> 30`으로 정렬했다.
- .github/workflows/ci.yml
  - `ssot-governance` job을 추가했다.
  - Rust clippy의 `continue-on-error`를 제거했다.
- backend/shared/assay/registry.go
- backend/shared/assay/registry_test.go
  - 카트리지별 assay 정의, alias 정규화, unknown cartridge 거부, channel completeness 기반 confidence 산출을 추가했다.
- backend/services/measurement-service/internal/handler/grpc.go
  - `SCorrected -> PrimaryValue`, `Unit=mg/dL`, `Confidence=0.95` 하드코딩을 제거했다.
- backend/services/measurement-service/internal/service/measurement.go
  - 세션 카트리지 타입 backfill 이후 assay registry로 측정 의미론을 적용한다.
- backend/services/measurement-service/internal/handler/grpc_stream_test.go
- backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go
- backend/services/measurement-service/internal/service/measurement_test.go
  - 새 assay-derived confidence와 CRP `mg/L` 단위 계약을 검증하도록 테스트를 보강했다.
- docs/audit/ssot-constants-gate.md
- docs/audit/measurement-assay-semantics.md
  - 변경 이유, 품질 게이트, 잔여 리스크를 문서화했다.
- rust-core/manpasik-engine/src/ai/mod.rs
  - `InferenceEngine::production()`과 `allow_simulation` 정책을 추가해 운영 모드에서 runtime-backed model 없이 시뮬레이션 추론을 반환하지 않게 했다.
- scripts/security_release_gate.sh
  - production 후보 설정에서 dev secret, 기본 Keycloak admin password, Elasticsearch security disabled, `DB_SSLMODE=disable` 패턴을 거부한다.
  - Kubernetes runtime config의 `MAX_FINGERPRINT_DIMENSION=1792`, `DEFAULT_ALPHA=0.98`도 검증한다.
- infrastructure/kubernetes/base/config/configmap.yaml
  - runtime config 상수를 `1792`, `0.98`로 최신 SSOT에 맞췄다.
- docs/audit/ai-prod-mock-free-gate.md
- docs/audit/security-release-gates.md
  - AI production mock-free gate와 security release gate의 잔여 리스크를 문서화했다.
**quality_gate**:
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `.github/workflows/ci.yml`에서 `continue-on-error` 제거 확인: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./shared/assay ./services/measurement-service/internal/handler ./services/measurement-service/internal/service`: PASS
- `cargo test -p manpasik-engine simulation -- --nocapture`: PASS
- `git diff --check` targeted files: PASS, 단 `AGENTS.md`는 기존 CRLF 파일이라 Git LF 치환 경고가 표시됨
**notes**:
- 이번 단계는 P0-1부터 P0-4까지 완료 범위다.
- assay registry는 하드코딩 제거용 최소 골격이며, 인허가 수준의 LoD/LoQ, linearity, interference, reference method, LOINC/UCUM 확장은 다음 단계에서 추가해야 한다.
**next_steps**:
- P1 assay evidence factory로 LoD/LoQ, linearity, interference, reference method, LOINC/UCUM, validation status를 registry에 확장한다.
- 실제 TFLite/ONNX runtime-backed model 연결 후 `has_runtime_backed_model()`을 interpreter readiness 검증으로 교체한다.
---

##_2026-05-02_Codex_flutter_auth_active_user_propagation
**status**: done
**plan**:
- Gateway/Auth에서 확보한 `user_id`가 Flutter `AuthState`에 머무르지 않고 `TenantInterceptor`의 active user 저장소까지 이어지게 한다.
- 로그인, 회원가입, 소셜 로그인, 게스트/데모 로그인 성공 시 `active_user_id`를 저장한다.
- 로그아웃과 인증 상태 실패 시 active user/tenant context를 제거한다.
- social login 직접 REST 경로도 snake_case와 legacy camelCase 응답을 모두 읽게 한다.
- AuthNotifier/TenantInterceptor 테스트로 후속 REST 요청의 `X-User-ID` 헤더 기반을 검증한다.
**changes**:
- frontend/flutter-app/lib/shared/providers/auth_provider.dart
  - `TenantInterceptor.setActiveUser()`를 로그인/회원가입/소셜/게스트/데모 로그인 성공 경로에 연결했다.
  - `TenantInterceptor.clear()`를 로그아웃과 `checkAuthStatus()` 미인증 경로에 연결했다.
  - `socialLogin()` 응답 매핑을 `access_token/accessToken`, `refresh_token/refreshToken`, `user_id/userId`, `display_name/displayName` 호환으로 보강했다.
  - 게스트/데모 로그인과 로그아웃은 active user 저장/해제를 기다릴 수 있도록 `Future<void>`로 정리했다.
- frontend/flutter-app/test/shared/providers/auth_notifier_test.dart
  - `SharedPreferences` mock 초기화를 추가했다.
  - 로그인/회원가입/게스트/데모 로그인 후 `TenantInterceptor.getActiveUser()` 값을 검증했다.
  - 로그아웃과 `checkAuthStatus()` 후 active user 및 active tenant가 제거되는지 검증했다.
- docs/audit/flutter-auth-active-user-propagation.md
  - 상세 구축계획, 변경 범위, 품질 게이트, 잔여 리스크를 문서화했다.
**quality_gate**:
- `flutter test --no-pub test/shared/providers/auth_notifier_test.dart test/core/network/tenant_interceptor_test.dart test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart test/features/auth/domain/auth_result_test.dart`: PASS
- `flutter analyze --no-pub lib/shared/providers/auth_provider.dart lib/core/network/tenant_interceptor.dart test/shared/providers/auth_notifier_test.dart test/core/network/tenant_interceptor_test.dart`: PASS
- `flutter analyze --no-pub lib/shared/providers/auth_provider.dart lib/features/auth/presentation/login_screen.dart lib/features/settings/presentation/settings_screen.dart test/shared/providers/auth_notifier_test.dart test/widget_test.dart test/core/providers/measurement_history_provider_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/gateway/cmd ./backend/services/gateway/internal/handler ./backend/services/auth-service/cmd ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler`: PASS
- `git diff --check` targeted files: PASS
**notes**:
- `TenantInterceptor`는 이미 저장된 active user를 모든 REST 요청의 `X-User-ID` 헤더로 주입한다.
- 이번 단계로 로그인 성공 후 후속 REST 호출의 사용자 컨텍스트 누락 위험을 줄였다.
**next_steps**:
- Docker Desktop WSL integration 활성화 후 `scripts/gateway_auth_compose_smoke.sh`와 `scripts/measure_service_compose_smoke.sh`를 실행한다.
- 또는 Gateway live smoke에 `X-User-ID`/tenant header propagation 검증을 추가한다.
---

##_2026-05-02_Codex_gateway_auth_compose_container_smoke_harness
**status**: harness_done_runtime_blocked_by_local_docker_compose
**plan**:
- 실제 binary smoke 다음 단계로 `auth-service + gateway` container image와 Compose runtime을 검증하는 좁은 smoke를 만든다.
- 전체 dev stack 대신 auth-service와 gateway 두 컨테이너만 띄워 REST auth user_id 계약을 검증한다.
- Gateway host HTTP 포트는 동적으로 예약해 충돌을 줄인다.
- Docker가 연결된 환경에서는 compose up, external endpoint Go smoke, compose down cleanup까지 한 줄로 수행한다.
- Docker가 현재 WSL에 연결되지 않은 환경에서는 blocked status를 명확히 출력한다.
**changes**:
- infrastructure/docker/docker-compose.gateway-auth-smoke.yml
  - `auth-service`와 `gateway` 두 컨테이너 smoke stack을 추가했다.
  - auth-service는 인메모리 repo 조건, `GRPC_PORT=:50051`, `METRICS_PORT=:9100`, smoke JWT 설정으로 기동한다.
  - gateway는 `AUTH_SERVICE_ADDR=auth-service:50051`, `HTTP_PORT=:8080`, smoke JWT 설정으로 기동한다.
  - gateway host HTTP port를 `MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_PORT`로 매핑한다.
  - auth-service `/health`, gateway `/health` healthcheck를 추가했다.
- scripts/gateway_auth_compose_smoke.sh
  - Docker/Compose/Go preflight를 추가했다.
  - 빈 host port를 동적으로 예약해 compose에 주입한다.
  - compose up/down cleanup과 `TestGatewayAuthExternalEndpointSmoke` 실행을 자동화했다.
  - 현재 WSL Docker Compose 연결 부재 시 `GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`를 출력한다.
- docs/audit/gateway-auth-compose-container-smoke.md
  - 상세 구축계획, 실행 방법, 품질 게이트, 현재 Docker blocked 상태를 문서화했다.
- docs/audit/gateway-auth-live-grpc-smoke.md
  - container compose smoke 하네스 구축 완료와 현재 실행 blocked 상태를 반영했다.
**quality_gate**:
- `bash -n scripts/gateway_auth_compose_smoke.sh`: PASS
- `scripts/gateway_auth_compose_smoke.sh`: BLOCKED, `GATEWAY_AUTH_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/gateway/cmd ./backend/services/gateway/internal/handler ./backend/services/auth-service/cmd ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler`: PASS
- `flutter test --no-pub test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart test/features/auth/domain/auth_result_test.dart`: PASS
- `flutter analyze --no-pub lib/features/auth/data/auth_repository_rest.dart lib/features/auth/data/auth_repository_impl.dart test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart`: PASS
- `git diff --check` targeted files: PASS
**notes**:
- 이번 단계는 container smoke 하네스 구축 완료이며, 현재 로컬 WSL에서는 Docker Compose runtime이 없어 실제 container PASS는 보류다.
- Docker Desktop WSL integration이 활성화되면 `scripts/gateway_auth_compose_smoke.sh`를 그대로 재실행하면 된다.
**next_steps**:
- Docker Desktop WSL integration 활성화 후 `scripts/gateway_auth_compose_smoke.sh`를 PASS로 확정한다.
- 이어서 `scripts/measure_service_compose_smoke.sh` 또는 full dev compose에서 auth + measurement smoke를 실행한다.
---

##_2026-05-02_Codex_gateway_auth_live_grpc_smoke
**status**: done
**plan**:
- 실제 auth-service binary를 인메모리 User/Token 저장소로 기동한다.
- 실제 gateway binary를 별도 프로세스로 기동하고 `AUTH_SERVICE_ADDR`로 auth-service gRPC 주소를 주입한다.
- Gateway REST `/auth/register -> /auth/login -> /auth/refresh`를 호출해 gRPC 경계를 지난 `user_id`가 계속 동일한지 검증한다.
- auth-service observability HTTP 포트를 smoke별 동적 포트로 격리할 수 있게 한다.
- 기존 Gateway/Auth/Flutter Auth 회귀를 함께 실행한다.
**changes**:
- backend/services/auth-service/cmd/main.go
  - metrics/health HTTP server 주소를 `METRICS_PORT` env로 override할 수 있게 했다. 기본값은 기존 `:9100` 유지.
- backend/services/gateway/cmd/main_smoke_test.go
  - `TestGatewayAuthLiveBinarySmoke`를 추가했다.
  - 테스트 내부에서 auth-service와 gateway production `main` binary를 `go build`로 생성한다.
  - auth-service는 동적 gRPC/metrics 포트, 인메모리 저장소, smoke JWT 설정으로 시작한다.
  - gateway는 동적 HTTP 포트와 실제 auth-service gRPC 주소로 시작한다.
  - REST register/login/refresh 응답의 `user_id` 일치와 access/refresh token 존재를 검증한다.
  - `MANPASIK_GATEWAY_AUTH_SMOKE_HTTP_ADDR`가 있으면 기존 외부 Gateway endpoint에도 동일 lifecycle을 실행하는 external smoke를 추가했다.
- docs/audit/gateway-auth-live-grpc-smoke.md
  - 상세 구축계획, 프로세스 경계, 품질 게이트, 잔여 리스크를 문서화했다.
- docs/audit/gateway-auth-user-id-rest-surface.md
  - live binary smoke 완료로 기존 프로세스 경계 리스크가 해소된 상태를 반영했다.
**quality_gate**:
- `MANPASIK_GO_BINARY=/home/kangjh3kang/sdk/go-go1.26.2/bin/go /home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v -count=1 ./backend/services/gateway/cmd -run TestGatewayAuthLiveBinarySmoke`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/gateway/cmd ./backend/services/gateway/internal/handler ./backend/services/auth-service/cmd ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler`: PASS
- `flutter test --no-pub test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart test/features/auth/domain/auth_result_test.dart`: PASS
- `flutter analyze --no-pub lib/features/auth/data/auth_repository_rest.dart lib/features/auth/data/auth_repository_impl.dart test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart`: PASS
- `git diff --check` targeted files: PASS
**notes**:
- Docker/Compose 없이도 실제 OS process 2개와 real TCP gRPC/HTTP stack을 통과한다.
- `TestGatewayAuthExternalEndpointSmoke`는 env 미설정 시 SKIP되며, compose나 외부 환경에서 같은 REST lifecycle을 재사용하기 위한 하네스다.
**next_steps**:
- Docker Desktop WSL integration 활성화 후 `scripts/measure_service_compose_smoke.sh`를 실행한다.
- 이후 full dev compose에서 gateway + auth-service container live smoke로 확장한다.
---

##_2026-05-01_Codex_gateway_auth_user_id_rest_surface
**status**: done
**plan**:
- Auth gRPC `LoginResponse.user_id` 보강을 Gateway REST 표면까지 확장한다.
- `/api/v1/auth/login`, `/api/v1/auth/refresh`, `/api/v1/auth/social-login` E2E mock 응답에 `user_id`를 채우고 테스트가 이를 필수 계약으로 검증하게 한다.
- 카카오 REST 로그인은 임시 `kakao_*` 식별자보다 auth-service가 확정한 `UserId`를 우선 사용하게 정리한다.
- Flutter REST AuthRepository가 Gateway snake_case와 legacy camelCase 응답을 모두 읽어 세션 ID 누락을 막게 한다.
- Go Gateway/Auth 회귀와 Flutter Auth 테스트/analyze로 검증한다.
**changes**:
- backend/services/gateway/internal/handler/auth_routes.go
  - 카카오 로그인 REST 응답의 `user_id`를 `resp.UserId` 우선, 기존 `kakaoID` fallback으로 변경했다.
  - 빈 auth-service user id에 대한 fallback helper `authUserID`를 추가했다.
- backend/services/gateway/internal/handler/e2e_test.go
  - mock auth client의 login/refresh/social-login 응답에 `UserId`를 추가했다.
  - `TestE2E_AuthLogin`, `TestE2E_AuthRefresh`, `TestE2E_AuthSocialLogin`이 `user_id` 값을 검증하도록 보강했다.
- frontend/flutter-app/lib/features/auth/data/auth_repository_rest.dart
  - `access_token/accessToken`, `refresh_token/refreshToken`, `user_id/userId`, `display_name/displayName`을 모두 읽는 매핑 helper를 추가했다.
  - 정상 Gateway 응답에서는 `unknown` 대신 실제 REST `user_id`를 AuthResult에 반영한다.
- frontend/flutter-app/test/features/auth/data/auth_repository_rest_test.dart
  - 로컬 HTTP 서버 기반 로그인 테스트를 추가해 Gateway snake_case `user_id`와 legacy camelCase `userId` 매핑을 검증했다.
- docs/audit/gateway-auth-user-id-rest-surface.md
  - 상세 구축계획, 변경 범위, 품질 게이트, 잔여 리스크를 문서화했다.
- docs/audit/auth-login-user-identity-surface.md
  - Gateway REST/E2E 사용자 식별자 리스크가 해소된 상태를 반영했다.
**quality_gate**:
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/gateway/internal/handler`: PASS
- `flutter test --no-pub test/features/auth/data/auth_repository_rest_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/gateway/cmd ./backend/services/gateway/internal/handler ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler`: PASS
- `flutter test --no-pub test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart test/features/auth/domain/auth_result_test.dart`: PASS
- `flutter analyze --no-pub lib/features/auth/data/auth_repository_rest.dart lib/features/auth/data/auth_repository_impl.dart test/features/auth/data/auth_repository_rest_test.dart test/features/auth/data/auth_proto_contract_test.dart`: PASS
- `git diff --check` targeted files: PASS
**notes**:
- Gateway `writeProtoJSON`은 snake_case proto name을 내보낼 수 있으므로 Flutter REST repository는 snake_case를 1순위로 읽는다.
- 기존 camelCase 응답을 내보내던 중간 Gateway/legacy 서버와의 호환도 유지했다.
- Docker Compose container smoke는 현재 WSL Docker integration 부재로 여전히 runtime blocked 상태다.
**next_steps**:
- Docker Desktop WSL integration 활성화 후 `scripts/measure_service_compose_smoke.sh`를 실행한다.
- 또는 Gateway auth live gRPC smoke를 추가해 실제 auth-service binary와 REST Gateway 사이의 로그인 user_id 전달을 별도 프로세스 경계에서 검증한다.
---

##_2026-05-01_Codex_auth_login_user_id_surface
**status**: done
**plan**:
- `LoginResponse`에 하위 호환 필드 `user_id = 5`를 추가한다.
- Go auth-service의 `TokenPair`와 gRPC 응답 매핑이 로그인/갱신/소셜 로그인 모두에서 같은 사용자 ID를 전달하게 한다.
- Flutter AuthRepository가 더 이상 정상 로그인에서 `unknown` 세션 ID를 사용하지 않게 한다.
- Go/Dart generated proto를 현재 `manpasik.proto` 기준으로 재생성하고 계약 테스트를 추가한다.
- Measure golden path 회귀를 함께 실행해 proto 확장이 기존 스트림 계약을 흔들지 않는지 확인한다.
**changes**:
- backend/shared/proto/manpasik.proto
  - `LoginResponse.user_id = 5`를 추가했다.
- backend/shared/gen/go/v1/manpasik.pb.go
- backend/shared/gen/go/v1/manpasik_grpc.pb.go
  - 현재 proto 기준 Go generated output을 재생성했다.
- frontend/flutter-app/lib/generated/manpasik.pb.dart
- frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart
- frontend/flutter-app/lib/generated/manpasik.pbenum.dart
- frontend/flutter-app/lib/generated/manpasik.pbjson.dart
  - official Dart generated output을 64,750라인 기준으로 재생성했다.
- backend/services/auth-service/internal/service/auth.go
  - `TokenPair.UserID`를 추가하고 `generateTokenPair()`가 실제 사용자 ID를 채우도록 했다.
- backend/services/auth-service/internal/handler/grpc.go
  - `Login`, `RefreshToken`, `SocialLogin` 응답에 `UserId`를 매핑했다.
- backend/services/auth-service/internal/service/auth_test.go
  - 로그인, 토큰 갱신, 소셜 로그인에서 `UserID`가 비지 않는지/기대 사용자와 일치하는지 검증했다.
- backend/services/auth-service/internal/handler/grpc_test.go
  - gRPC handler의 로그인/갱신 응답이 `UserId`를 전달하는지 검증했다.
- frontend/flutter-app/lib/features/auth/data/auth_repository_impl.dart
  - `res.userId`를 우선 사용하고, legacy empty response에서만 `unknown` fallback을 유지했다.
- frontend/flutter-app/test/features/auth/data/auth_proto_contract_test.dart
  - official generated `LoginResponse`가 `userId`를 바이너리 직렬화 왕복으로 보존하는지 검증했다.
- docs/audit/auth-login-user-identity-surface.md
  - 상세 구현계획, 변경 범위, 품질 게이트, 잔여 리스크를 문서화했다.
- docs/audit/dart-proto-official-generated-replacement.md
  - 이전 잔여 리스크였던 Auth `user_id` 공백이 해결된 상태를 반영했다.
- docs/audit/mock-retirement-register.md
  - official Dart generated compile gate 라인 수를 현재 64,750라인으로 갱신했다.
**quality_gate**:
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler`: PASS
- `flutter test --no-pub test/features/auth/data/auth_proto_contract_test.dart`: PASS
- `scripts/check_proto_generation_preflight.sh`: PASS, generated Dart 64,750 lines, `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
- `scripts/check_proto_generation_compile_gate.sh`: PASS, generated Dart 64,750 lines, `grpc 5.1.0`, `protobuf 6.0.0`
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/auth-service/internal/service ./backend/services/auth-service/internal/handler ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `flutter test --no-pub test/features/auth/data/auth_proto_contract_test.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/generated/manpasik.pbenum.dart lib/generated/manpasik.pbjson.dart lib/features/auth/data/auth_repository_impl.dart test/features/auth/data/auth_proto_contract_test.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `git diff --check` targeted files: PASS
**notes**:
- `user_id = 5`는 새 optional string 필드라 기존 클라이언트와 wire 호환된다.
- legacy 서버가 빈 `user_id`를 돌려주는 경우 Flutter는 여전히 `unknown` fallback으로 크래시 대신 낮은 품질 상태를 유지한다.
- Docker Compose container smoke는 이번 범위 밖이며, 현재 WSL Docker integration 부재로 여전히 blocked 상태다.
**next_steps**:
- Docker Desktop WSL integration 활성화 후 `scripts/measure_service_compose_smoke.sh`를 실행해 container smoke PASS를 확보한다.
- Gateway auth REST/E2E 응답에도 `user_id`를 노출할 필요가 있는지 점검한다.
---

##_2026-05-01_Codex_dart_proto_official_generated_replacement
**status**: done
**plan**:
- Flutter 앱 `grpc/protobuf` lock을 공식 `protoc_plugin 25.0.0` 산출물과 맞춘다.
- `grpc ^5.1.0`, `protobuf ^6.0.0`, `fixnum ^1.1.1`을 direct dependency로 고정한다.
- `protoc --dart_out=grpc` 공식 산출물을 checked-in `lib/generated`로 교체한다.
- 수동 stub에 기대던 소비자 차이는 official enum/client/response API에 맞춰 좁게 정렬한다.
- wire golden, repository, orchestrator, Go service 회귀로 교체 안전성을 확인한다.
**changes**:
- frontend/flutter-app/pubspec.yaml
  - `grpc`를 `^5.1.0`, `protobuf`를 `^6.0.0`으로 올리고 `fixnum ^1.1.1`을 direct dependency로 추가했다.
- frontend/flutter-app/pubspec.lock
  - `grpc 5.1.0`, `protobuf 6.0.0`, `fixnum 1.1.1` 정렬을 반영했다.
- frontend/flutter-app/lib/generated/manpasik.pb.dart
- frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart
- frontend/flutter-app/lib/generated/manpasik.pbenum.dart
- frontend/flutter-app/lib/generated/manpasik.pbjson.dart
  - 수동 스텁을 공식 `protoc_plugin 25.0.0` generated output 64,736라인으로 교체/확장했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - official pbgrpc export에 맞춰 중복 import를 제거했다.
- frontend/flutter-app/lib/features/auth/data/auth_repository_impl.dart
  - official `LoginResponse`에 `userId`가 없으므로 기존 fallback 성격의 `unknown`을 명시 사용했다.
- frontend/flutter-app/lib/features/devices/data/device_repository_impl.dart
  - `DeviceStatus` official enum mapping으로 정렬했다.
- frontend/flutter-app/lib/features/user/data/user_repository_impl.dart
  - `SubscriptionTier` enum을 도메인 `int`로 넘길 때 `.value`를 사용하도록 정렬했다.
- frontend/flutter-app/lib/shared/providers/chat_provider.dart
  - official client 이름 `AiInferenceServiceClient`로 정렬했다.
- frontend/flutter-app/lib/shared/providers/admin_settings_provider.dart
  - 중복 pb import를 제거했다.
- frontend/flutter-app/test/features/measurement/data/measurement_stream_wire_contract_test.dart
  - 테스트명을 official Dart proto wire contract로 갱신했다.
- docs/audit/dart-proto-official-generated-replacement.md
  - 상세 구현계획, 소비자 정렬, 품질 게이트, 잔여 리스크를 문서화했다.
- docs/audit/dart-proto-generation-preflight.md
- docs/audit/dart-proto-full-generated-compile-gate.md
- docs/audit/measure-stream-wire-smoke-gate.md
- docs/audit/mock-retirement-register.md
  - official checked-in generated 교체 완료 상태를 반영했다.
**quality_gate**:
- `flutter pub get`: PASS
- `scripts/check_proto_generation_preflight.sh`: PASS, `PROTO_PREFLIGHT_FULL_REPLACEMENT=ready_for_compile_gate`
- `scripts/check_proto_generation_compile_gate.sh`: PASS, generated Dart 64,736 lines, `grpc 5.1.0`, `protobuf 6.0.0`
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/generated/manpasik.pbenum.dart lib/generated/manpasik.pbjson.dart lib/features/measurement/data/measurement_repository_impl.dart lib/features/auth/data/auth_repository_impl.dart lib/features/user/data/user_repository_impl.dart lib/shared/providers/chat_provider.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
**notes**:
- official `LoginResponse` proto에는 `user_id`가 없어서 Auth 도메인 사용자 식별자는 임시 fallback `unknown`을 유지한다. 다음 단계에서 proto/API 차원 보강이 필요하다.
- Docker Compose container smoke는 여전히 현재 WSL Docker integration 부재로 runtime blocked 상태다.
**next_steps**:
- Auth login 사용자 식별자 surface를 proto/API/Flutter repository까지 보강한다.
- 또는 Docker Desktop WSL integration 활성화 후 `scripts/measure_service_compose_smoke.sh`를 실행해 container smoke PASS를 확보한다.
---

##_2026-05-01_Codex_dart_proto_full_generated_compile_gate
**status**: done
**plan**:
- checked-in 수동 `lib/generated` 스텁은 건드리지 않고, 임시 Dart package에서 공식 `protoc --dart_out=grpc` 전체 산출물을 생성한다.
- `protoc_plugin 25.0.0` 산출물과 호환되는 `grpc ^5.1.0`, `protobuf ^6.0.0`, `fixnum ^1.1.1` 조합으로 compile gate를 분리한다.
- 생성 산출물에 `MeasurementData`와 `streamMeasurement` surface가 있는지 확인한다.
- `dart pub get`, `dart analyze bin lib`, `dart run bin/compile_smoke.dart`를 통과시켜 전체 generated output의 컴파일 가능성을 증명한다.
**changes**:
- frontend/flutter-app/scripts/check_proto_generation_compile_gate.sh
  - 임시 package에 `manpasik.proto` 전체 Dart gRPC 산출물을 생성한다.
  - `MeasurementData`, `MeasurementResult`, `MeasurementServiceClient`를 실제로 참조하는 compile smoke를 추가했다.
  - resolved `grpc`/`protobuf` 버전과 generated Dart 라인 수를 출력한다.
- docs/audit/dart-proto-full-generated-compile-gate.md
  - 상세 구현계획, 격리 compile gate 범위, 품질 게이트, 잔여 앱 lock 정렬 단계를 문서화했다.
- docs/audit/dart-proto-generation-preflight.md
  - current app lock preflight는 계속 blocked지만, 격리 compile gate는 통과한 상태를 반영했다.
- docs/audit/measure-stream-wire-smoke-gate.md
  - 전체 generated compile gate 완료와 다음 앱 lock 정렬 단계를 반영했다.
- docs/audit/measure-grpc-compose-container-smoke-gate.md
  - proto compile gate 완료 후 남은 작업을 앱 lock 정렬/checked-in 교체로 좁혔다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 공식 Dart generated 64,736라인 compile gate 통과를 추가했다.
**quality_gate**:
- `bash -n scripts/check_proto_generation_compile_gate.sh`: PASS
- `scripts/check_proto_generation_preflight.sh`: PASS, current app lock replacement remains blocked by `protobuf 3.1.0` Timestamp import compatibility
- `scripts/check_proto_generation_compile_gate.sh`: PASS, `PROTO_COMPILE_GATE_STATUS=passed`, generated Dart 64,736 lines, `grpc 5.1.0`, `protobuf 6.0.0`
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
**notes**:
- 현재 앱 `pubspec.lock`는 여전히 `grpc 4.1.0`, `protobuf 3.1.0`이다.
- 이번 단계는 official generated output의 컴파일 가능성을 증명한 gate이며, checked-in generated file 교체는 앱 의존성 lock 정렬 후 별도 단계로 수행해야 한다.
**next_steps**:
- Flutter 앱 `pubspec.yaml`/`pubspec.lock`를 `grpc 5.1.0`, `protobuf 6.0.0` 계열로 실제 정렬한다.
- 그 다음 checked-in `lib/generated/manpasik.pb*.dart`를 공식 generated output으로 교체하고 기존 wire/orchestrator 회귀를 재실행한다.
---

##_2026-05-01_Codex_measure_grpc_compose_container_smoke_harness
**status**: harness_done_runtime_blocked_by_local_docker_compose
**plan**:
- 실제 service binary smoke 다음 단계로 container image와 compose runtime을 검증하는 narrow smoke를 만든다.
- 전체 dev dependency stack을 즉시 띄우기보다 measurement-service 단일 컨테이너의 build, env, port publish, HTTP health, gRPC health, measurement lifecycle을 먼저 고정한다.
- Docker가 연결된 환경에서는 스크립트 한 줄로 compose up, external endpoint Go smoke, compose down cleanup까지 수행한다.
- Docker가 현재 WSL에 연결되지 않은 환경에서는 blocked status를 명확히 출력한다.
**changes**:
- infrastructure/docker/docker-compose.measurement-smoke.yml
  - `backend/services/measurement-service/Dockerfile` 기반 `manpasik/measurement-service:smoke` build 정의를 추가했다.
  - `GRPC_PORT=:50054`, `HTTP_PORT=:8080`, `VERSION=compose-smoke`, `TENANCY_ENFORCED=false` runtime env를 고정했다.
  - host loopback dynamic port publish와 HTTP `/health` healthcheck를 추가했다.
- scripts/measure_service_compose_smoke.sh
  - Docker/Compose/Go preflight를 추가했다.
  - 빈 host port를 동적으로 예약해 compose에 주입한다.
  - compose up/down cleanup과 `TestMeasurementServiceExternalEndpointSmoke` 실행을 자동화했다.
  - 현재 WSL Docker Compose 연결 부재 시 `MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`를 출력한다.
- backend/services/measurement-service/cmd/main_smoke_test.go
  - `TestMeasurementServiceExternalEndpointSmoke`를 추가했다.
  - binary smoke와 external endpoint smoke가 공유하는 lifecycle 검증 helper를 만들었다.
  - 외부 HTTP health, gRPC health `SERVING`, `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` 검증을 추가했다.
- docs/audit/measure-grpc-compose-container-smoke-gate.md
  - 상세 구현계획, 실행 방법, 품질 게이트, 현재 Docker blocked 상태를 문서화했다.
- docs/audit/measure-grpc-service-binary-smoke-gate.md
  - container smoke 하네스 구축 완료와 현재 실행 blocked 상태를 반영했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 compose container smoke 하네스 구축과 Docker 환경 blocked 상태를 추가했다.
**quality_gate**:
- `bash -n scripts/measure_service_compose_smoke.sh`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/cmd -run 'TestMeasurementService(BinarySmoke|ExternalEndpointSmoke)'`: PASS, external endpoint smoke는 env 미지정 시 SKIP
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `scripts/measure_service_compose_smoke.sh`: BLOCKED, `MEASUREMENT_COMPOSE_SMOKE_STATUS=blocked_docker_compose_unavailable`
**next_steps**:
- Docker Desktop WSL integration을 활성화한 뒤 `scripts/measure_service_compose_smoke.sh`를 재실행해 실제 container smoke PASS를 확보한다.
- 이후 full dev compose의 `measurement-service + postgres + gateway` 경로 또는 Dart generated proto compile gate로 확장한다.
---

##_2026-05-01_Codex_measure_grpc_service_binary_smoke_gate
**status**: done
**plan**:
- `go build`로 `backend/services/measurement-service/cmd`의 실제 production `main` 바이너리를 임시 생성한다.
- 별도 OS process로 service binary를 기동하고 `GRPC_PORT`, `HTTP_PORT`, `VERSION`, shutdown timeout만 주입한다.
- 외부 DB/Milvus/Kafka/Elasticsearch env는 주입하지 않아 main wiring의 memory fallback 기동을 검증한다.
- HTTP `/health`, gRPC health `SERVING`, `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` lifecycle을 generated client로 검증한다.
**changes**:
- backend/services/measurement-service/cmd/main_smoke_test.go
  - `TestMeasurementServiceBinarySmoke`를 추가했다.
  - 테스트 내부에서 실제 `measurement-service` 바이너리를 build하고 별도 process로 실행한다.
  - HTTP health payload의 `healthy/measurement-service/smoke-test`를 검증한다.
  - gRPC health service의 `SERVING` 상태를 검증한다.
  - generated `MeasurementServiceClient`로 session 생성, stream frame 전송, history 저장 조회, session 종료를 검증한다.
  - `SIGTERM` 기반 graceful shutdown cleanup을 추가했다.
- docs/audit/measure-grpc-service-binary-smoke-gate.md
  - 상세 구현계획, main wiring 검증 범위, 품질 게이트, 다음 단계를 문서화했다.
- docs/audit/measure-grpc-helper-process-smoke-gate.md
  - 실제 service binary smoke 완료 상태를 반영했다.
- docs/audit/measure-grpc-tcp-loopback-smoke-gate.md
  - service binary smoke 완료와 dev compose container smoke 후속 단계를 반영했다.
- docs/audit/measure-grpc-transport-smoke-gate.md
  - transport 증거 사슬에 service binary smoke 완료를 추가했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 실제 service binary smoke를 추가하고 잔여 조치를 dev compose container smoke와 Dart proto generation 공식화로 좁혔다.
**quality_gate**:
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/cmd -run TestMeasurementServiceBinarySmoke`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/cmd ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
**next_steps**:
- dev compose에서 containerized measurement-service gRPC/HTTP health smoke를 추가한다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
---

##_2026-05-01_Codex_measure_grpc_helper_process_smoke_gate
**status**: done
**plan**:
- 현재 Go test binary를 별도 child process로 재실행해 parent process와 gRPC server process를 분리한다.
- child process는 memory repository 기반 `MeasurementService`와 `MeasurementHandler`를 구성하고 `127.0.0.1:0`에 `grpc.Server`를 기동한다.
- parent process는 generated Go `MeasurementServiceClient`로 접속해 `StartSession -> StreamMeasurement -> GetMeasurementHistory -> EndSession` lifecycle을 실행한다.
- 응답값, fingerprint vector, history 저장 조회, EOF, graceful shutdown을 검증한다.
**changes**:
- backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go
  - `TestStreamMeasurementAgainstHelperProcess`를 추가했다.
  - `TestMeasurementServiceProcessHelper`가 child test process에서 실제 TCP gRPC server를 기동한다.
  - parent test가 generated client로 `StartSession`, `StreamMeasurement`, `GetMeasurementHistory`, `EndSession`을 호출한다.
  - child process address file publish, interrupt signal, `GracefulStop()` cleanup을 검증 가능한 helper로 구성했다.
- docs/audit/measure-grpc-helper-process-smoke-gate.md
  - 상세 구현계획, helper-process cross-process 범위, 실제 production binary가 아닌 한계, 품질 게이트를 문서화했다.
- docs/audit/measure-grpc-tcp-loopback-smoke-gate.md
  - helper-process 후속 gate 완료 상태와 다음 service-process smoke 단계를 반영했다.
- docs/audit/measure-grpc-transport-smoke-gate.md
  - transport smoke 증거 사슬에 helper-process gate 완료를 추가했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 helper-process cross-process smoke를 추가하고 잔여 조치를 실제 service binary/dev compose smoke로 좁혔다.
**quality_gate**:
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/internal/handler -run TestStreamMeasurementAgainstHelperProcess`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
**next_steps**:
- 실제 `measurement-service` cmd binary 또는 dev compose service를 별도 process로 띄우는 service-process smoke를 추가한다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
---

##_2026-05-01_Codex_measure_grpc_tcp_loopback_smoke_gate
**status**: done
**plan**:
- `127.0.0.1:0` TCP listener를 열어 실제 OS network stack을 통과하는 gRPC server를 기동한다.
- generated Go `MeasurementServiceClient`로 TCP 주소에 접속해 `StreamMeasurement` bidi stream을 연다.
- 단일 frame이 아니라 같은 stream 안에서 2개 frame을 보내고 2개 응답과 EOF를 검증한다.
- measurement 저장소 2건 저장, session metadata backfill, corrected value, fingerprint vector 저장까지 검증한다.
**changes**:
- backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go
  - `TestStreamMeasurementOverTCPLoopbackHandlesMultipleFrames`를 추가했다.
  - TCP loopback gRPC stream에서 `Send` 2회, `CloseSend`, `Recv` 2회, EOF lifecycle을 검증한다.
  - 각 응답의 primary value, unit, confidence, fingerprint vector, `ProcessedAt`을 검증한다.
  - 저장된 measurement 2건과 마지막 fingerprint vector 저장 상태를 검증한다.
- docs/audit/measure-grpc-tcp-loopback-smoke-gate.md
  - 상세 구현계획, 실행 결과, 품질 게이트, 다음 단계를 문서화했다.
- docs/audit/measure-grpc-transport-smoke-gate.md
  - TCP loopback 후속 gate 완료 상태를 반영했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 TCP loopback multi-frame smoke를 추가했다.
**quality_gate**:
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -v ./backend/services/measurement-service/internal/handler -run TestStreamMeasurementOverTCP`: PASS
- `/home/kangjh3kang/sdk/go-go1.26.2/bin/go test -count=1 ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
**next_steps**:
- 실제 `measurement-service` binary를 별도 process로 띄우는 cross-process smoke를 추가한다.
- Dart proto toolchain 버전 정렬 후 전체 generated Dart compile gate를 실행한다.
---

##_2026-05-01_Codex_measure_grpc_transport_smoke_gate
**status**: done
**plan**:
- `bufconn`으로 실제 `grpc.Server`를 외부 포트 없이 기동한다.
- 기존 `MeasurementHandler`와 memory test repository를 서버에 등록한다.
- generated Go `MeasurementServiceClient`로 `StreamMeasurement`를 열어 client `Send`, `CloseSend`, `Recv`, EOF lifecycle을 검증한다.
- 응답값, measurement 저장, fingerprint vector 저장 side effect를 함께 확인한다.
**changes**:
- backend/services/measurement-service/internal/handler/grpc_stream_transport_smoke_test.go
  - `bufconn` 기반 in-process gRPC transport smoke test를 추가했다.
  - generated client/server 경로로 `StreamMeasurement` frame을 보내고 `MeasurementResult` 응답을 검증한다.
  - `ProcessedAt`, 저장된 measurement metadata, corrected value, raw channels, fingerprint vector 저장을 검증한다.
- docs/audit/measure-grpc-transport-smoke-gate.md
  - 상세 구현계획, transport smoke 범위, 품질 게이트, 다음 단계를 문서화했다.
- docs/audit/measure-stream-wire-smoke-gate.md
  - network-level 후속 gate 완료 상태를 반영했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 in-process gRPC transport smoke 완료를 추가했다.
**quality_gate**:
- `go test -count=1 ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
**next_steps**:
- 실제 dev compose 또는 로컬 service process를 대상으로 하는 cross-process gRPC smoke를 추가한다.
- Dart `protobuf`, `grpc`, `protoc_plugin` 버전 조합을 정렬해 전체 generated Dart compile gate를 실행한다.
---

##_2026-05-01_Codex_measure_stream_wire_smoke_gate
**status**: done
**plan**:
- Dart 수동 protobuf 스텁이 생성하는 `MeasurementData` binary payload를 golden hex로 고정한다.
- 같은 golden payload를 Go generated proto에서 unmarshal하고 실제 measurement-service `StreamMeasurement` handler에 넣어 저장/응답까지 확인한다.
- 전체 Dart generated proto 치환은 preflight script로 분리해 현재 toolchain 호환성을 덮어쓰기 없이 보고한다.
**changes**:
- frontend/flutter-app/test/features/measurement/data/measurement_stream_wire_contract_test.dart
  - Dart 수동 `MeasurementData` 스텁의 `writeToBuffer()` 결과를 Go proto golden hex와 비교한다.
  - round-trip decode로 raw channel, differential, env metadata 필드를 검증한다.
- backend/services/measurement-service/internal/handler/grpc_stream_wire_contract_test.go
  - 같은 golden hex를 Go `v1.MeasurementData`로 unmarshal한다.
  - unmarshal된 frame을 실제 `MeasurementHandler.StreamMeasurement()` fake stream에 넣어 measurement 저장, fingerprint 저장, stream response를 검증한다.
- frontend/flutter-app/scripts/check_proto_generation_preflight.sh
  - 임시 디렉터리에 Dart gRPC 산출물을 생성한다.
  - `MeasurementData`와 `streamMeasurement` 존재를 확인한다.
  - locked `protobuf 3.1.0` 기준 Timestamp import 호환성을 점검해 전체 치환 보류 사유를 출력한다.
- docs/audit/measure-stream-wire-smoke-gate.md
  - 상세 구현계획, 실행 결과, 품질 게이트, 다음 단계를 문서화했다.
- docs/audit/dart-proto-generation-preflight.md
  - 전체 Dart generated proto 전환 전 preflight 판정과 보류 사유를 문서화했다.
- docs/audit/measure-native-grpc-stream-binding.md
  - wire smoke/preflight 완료 상태를 후속 진행에 반영했다.
- docs/audit/mock-retirement-register.md
  - MR-009 증거에 Dart/Go golden wire smoke gate를 추가했다.
**quality_gate**:
- `go test ./backend/services/measurement-service/internal/handler ./backend/services/measurement-service/internal/service`: PASS
- `bash scripts/check_proto_generation_preflight.sh`: PASS, `StreamMeasurement` 생성 가능, 전체 치환은 locked `protobuf 3.1.0` Timestamp import 호환성으로 보류
- `flutter test --no-pub test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/core/providers/measurement_history_provider_test.dart test/features/measurement/domain/measurement_domain_test.dart test/features/domain_models_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart test/features/measurement/data/measurement_stream_wire_contract_test.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
**next_steps**:
- `protobuf`, `grpc`, `protoc_plugin` 버전 조합을 정렬해 전체 Dart generated proto compile gate를 실행한다.
- in-process 또는 dev compose 기반 network-level gRPC smoke test로 확장한다.
---

##_2026-05-01_Codex_measure_native_grpc_stream_binding_gate
**status**: done
**root_cause**:
- MR-009 local echo fallback은 제거됐지만 native Flutter `processMeasurement()`가 임시 Gateway REST bridge를 계속 사용하고 있었다.
- checked-in Dart gRPC 스텁은 수동 축약본이라 `MeasurementService.StreamMeasurement` client method와 관련 메시지가 빠져 있었다.
**changes**:
- frontend/flutter-app/lib/generated/manpasik.pb.dart
  - `MeasurementData`, `DifferentialCorrection`, `EnvironmentMeta`, `MeasurementResult` 수동 protobuf 메시지를 추가했다.
- frontend/flutter-app/lib/generated/manpasik.pbgrpc.dart
  - `MeasurementServiceClient.streamMeasurement()` bidi stream client method를 추가했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - native `processMeasurement()`를 Gateway REST process endpoint 대신 `MeasurementService.StreamMeasurement` gRPC stream 직결로 전환했다.
  - 테스트 주입용 `MeasurementStreamCall`을 추가해 stream frame 계약을 네트워크 없이 검증할 수 있게 했다.
- frontend/flutter-app/lib/core/providers/grpc_provider.dart
  - native measurement repository에서 REST process client 주입을 제거했다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart
  - native process path가 `MeasurementData` frame을 전송하고 `MeasurementResult`를 도메인 result로 매핑하는지 검증하도록 갱신했다.
- docs/audit/measure-native-grpc-stream-binding.md
  - direct gRPC stream 전환, 수동 스텁 보강 범위, Timestamp 제약, 품질 게이트를 문서화했다.
- docs/audit/mock-retirement-register.md
  - MR-009를 native gRPC stream 직결 완료 상태로 갱신했다.
- docs/audit/measure-mock-retirement-reconcile.md
  - MR-009 후속 진행 완료 상태와 다음 smoke test 단계를 반영했다.
- docs/audit/measure-history-stale-state.md
  - 다음 단계 문구를 direct stream 완료 이후 과제로 갱신했다.
**quality_gate**:
- `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart lib/core/providers/grpc_provider.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
- `flutter test --no-pub test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/features/measurement/domain/measurement_domain_test.dart test/features/domain_models_test.dart`: PASS
- `flutter analyze --no-pub lib/generated/manpasik.pb.dart lib/generated/manpasik.pbgrpc.dart lib/features/measurement/data/measurement_repository_impl.dart lib/core/providers/grpc_provider.dart lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_rest.dart lib/shared/providers/ecosystem_providers.dart lib/features/measurement/presentation/measurement_result_screen.dart lib/features/home/presentation/home_screen.dart test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart`: PASS
**notes**:
- 전체 `protoc --dart_out=grpc` 재생성은 `/tmp` 검증에서 약 6.4만 줄 산출물과 `pbenum/pbjson` 추가 파일을 만들었다.
- 현재 lock된 `protobuf 3.1.0` 환경은 최신 protoc plugin의 well-known Timestamp import 경로와 맞지 않아 이번 단계에서는 필요한 stream surface만 수동 스텁에 좁게 보강했다.
**next_steps**:
- Dart protobuf/protoc plugin 버전을 정렬해 전체 generated Dart 공식 산출물 전환 가능성을 검증한다.
- 실제 measurement-service integration smoke test를 추가한다.
---

##_2026-05-01_Codex_measure_history_stale_state_gate
**status**: done
**root_cause**:
- MR-001은 DemoMode mock history를 명시적으로 허용하되 RealMode 실패를 운영자가 인지할 수 있어야 한다.
- 기존 RealMode history 실패는 빈 목록처럼 보일 수 있어 "측정 기록 없음"과 "동기화 실패"가 구분되지 않는 위험이 있었다.
**changes**:
- frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart
  - `MeasurementHistoryResult`에 `isStale`, `errorMessage`, `hasError` 상태 계약을 추가했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart
  - REST history 조회 실패를 빈 정상 결과 대신 stale/error result로 반환하도록 변경했다.
- frontend/flutter-app/lib/core/providers/grpc_provider.dart
  - RealMode `measurementHistoryProvider` 예외를 stale/error result로 노출한다.
  - DemoMode mock history는 명시적 demo 상태에서만 fresh mock data로 유지한다.
- frontend/flutter-app/lib/shared/providers/ecosystem_providers.dart
  - home dashboard aggregate에 measurement history stale/error 상태를 연결했다.
- frontend/flutter-app/lib/features/measurement/presentation/measurement_result_screen.dart
  - history 동기화 실패 empty state, inline warning, retry action을 추가했다.
- frontend/flutter-app/lib/features/home/presentation/home_screen.dart
  - dashboard hero/stat 영역에서 측정 기록 동기화 확인 필요 상태를 표시한다.
- frontend/flutter-app/test/core/providers/measurement_history_provider_test.dart
  - RealMode history failure가 stale/error result로 노출되는지와 DemoMode mock history가 fresh 상태인지 검증했다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_rest_test.dart
  - REST history failure expectation을 stale/error contract로 갱신했다.
- docs/audit/mock-retirement-register.md
  - MR-001 판정을 RealMode 보강 완료로 갱신하고 품질 게이트를 추가했다.
- docs/audit/measure-history-stale-state.md
  - 이번 stale/error 상태 계약, 변경 파일, 품질 게이트, 다음 단계를 문서화했다.
- docs/audit/measure-mock-retirement-reconcile.md
  - MR-001 stale/error 보강 완료 후속 상태를 반영했다.
**quality_gate**:
- `flutter test --no-pub test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/domain/measurement_domain_test.dart test/features/domain_models_test.dart test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart`: PASS
- `flutter analyze --no-pub lib/features/measurement/domain/measurement_repository.dart lib/features/measurement/data/measurement_repository_rest.dart lib/core/providers/grpc_provider.dart lib/shared/providers/ecosystem_providers.dart lib/features/measurement/presentation/measurement_result_screen.dart lib/features/home/presentation/home_screen.dart test/core/providers/measurement_history_provider_test.dart test/features/measurement/data/measurement_repository_rest_test.dart`: PASS
**next_steps**:
- Dart gRPC generated binding에 `StreamMeasurement`가 포함되도록 proto generation 파이프라인을 정리한다.
- native process REST bridge를 gRPC stream 직결로 치환한다.
---

##_2026-05-01_Codex_measure_mock_retirement_reconcile
**status**: done
**root_cause**:
- Measure 경로의 MR-002 readiness gate와 audit persistence는 보강됐지만, 네이티브 Flutter `MeasurementRepositoryImpl.processMeasurement()`가 Dart gRPC stream binding 부재를 이유로 Rust 처리 결과를 로컬에서 그대로 반환하고 있었다.
- 이 상태에서는 네이티브 Measure가 measurement-service 저장소와 fingerprint vector 저장 경로를 타지 못해 MR-009급 P0 fallback으로 관리해야 했다.
**changes**:
- frontend/flutter-app/lib/features/measurement/data/measurement_process_gateway_mapper.dart
  - Gateway `/measurements/process` payload encode/decode 공통 mapper를 추가했다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - 네이티브 `processMeasurement()`의 local echo fallback을 제거했다.
  - `ManPaSikRestClient`를 통해 Gateway `POST /api/v1/measurements/process`로 측정 프레임을 전송하도록 변경했다.
  - `accessTokenProvider` 값을 REST Authorization header에 반영한다.
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart
  - Web/REST 구현도 동일 mapper를 사용하도록 정리했다.
- frontend/flutter-app/lib/core/providers/grpc_provider.dart
  - 네이티브 `MeasurementRepositoryImpl`에 shared `restClientProvider`를 주입했다.
- frontend/flutter-app/lib/core/network/tenant_interceptor.dart
  - SharedPreferences 접근 실패 시 REST 요청 자체가 막히지 않도록 fail-open 처리했다.
- frontend/flutter-app/test/features/measurement/data/measurement_repository_impl_test.dart
  - 네이티브 process path가 Gateway endpoint로 payload와 Authorization header를 전송하고 응답을 매핑하는지 검증했다.
- docs/audit/mock-retirement-register.md
  - MR-009를 등록하고 Measure 재대조 결과를 추가했다.
  - MR-002는 운영 차단 조건 충족, MR-009는 퇴역 완료, MR-001은 조건부 유지로 판정했다.
- docs/audit/measure-mock-retirement-reconcile.md
  - 이번 재대조 결과, 변경 파일, 품질 게이트, 다음 단계를 문서화했다.
**quality_gate**:
- `flutter test --no-pub test/features/measurement/data/measurement_repository_impl_test.dart test/features/measurement/data/measurement_repository_rest_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/core/network/tenant_interceptor_test.dart`: PASS
- `flutter analyze --no-pub lib/core/network/tenant_interceptor.dart lib/core/providers/grpc_provider.dart lib/features/measurement/data/measurement_process_gateway_mapper.dart lib/features/measurement/data/measurement_repository_impl.dart lib/features/measurement/data/measurement_repository_rest.dart test/features/measurement/data/measurement_repository_impl_test.dart`: PASS
**next_steps**:
- MR-001 RealMode measurement history 실패를 빈 목록으로 삼키지 않고 stale/error 상태로 노출한다.
- Dart gRPC generated binding에 `StreamMeasurement`가 포함되도록 proto generation 파이프라인을 정리한 뒤 native process REST bridge를 gRPC stream 직결로 치환한다.
---

##_2026-05-01_Codex_measure_trace_audit_persistence_gate
**status**: done
**root_cause**:
- Measure golden path trace remote sink는 Gateway까지 도달했지만, audit-service 쓰기 intake가 없어 Gateway가 deferred 감사 상태만 반환했다.
- PHI 안전성을 위해 `primary_value`를 저장/전달하지 않으면서 phase/context 메타데이터만 감사 로그로 영속화해야 했다.
**changes**:
- backend/services/audit-service/internal/service/audit.go
  - `RecordActionInput`과 `RecordActionWithMetadata()`를 추가해 metadata, user agent, timestamp를 저장할 수 있게 했다.
  - 기존 `RecordAction()`은 새 저장 경로를 호출하도록 유지해 기존 호출부와 호환되게 했다.
- backend/services/audit-service/internal/handler/http.go
  - audit-service metrics/health mux에 `POST /audit/events` HTTP intake를 추가했다.
  - `action`, `resource_type`, RFC3339 `occurred_at`을 검증한 뒤 기존 audit repository로 저장한다.
- backend/services/audit-service/internal/handler/http_test.go
  - HTTP audit intake의 저장 성공과 필수값 검증 테스트를 추가했다.
- backend/services/audit-service/cmd/main.go
  - `:9100` HTTP 서버에 `/metrics`, `/health`와 함께 감사 쓰기 intake를 등록했다.
- backend/services/gateway/internal/handler/audit_recorder.go
  - audit-service로 감사 이벤트를 POST하는 `HTTPAuditEventRecorder`를 추가했다.
- backend/services/gateway/internal/handler/audit_recorder_test.go
  - POST payload와 persisted 응답 매핑을 검증하는 recorder 계약 테스트를 추가했다.
- backend/services/gateway/internal/handler/rest_handler.go
  - `SetAuditEventRecorder()`로 audit recorder를 주입할 수 있게 했다.
- backend/services/gateway/cmd/main.go
  - `AUDIT_INTAKE_URL` 바인딩을 추가해 런타임에서 persisted trace audit write를 활성화할 수 있게 했다.
- backend/services/gateway/internal/handler/measurement_routes.go
  - Measure trace event를 `measure.trace.<phase>`, `measurement_trace`, session/event resource ID, request IP, user agent, PHI 최소화 metadata를 가진 audit record로 변환한다.
  - audit-service 저장 성공 시 응답에 `audit_status=persisted`와 `audit_entry_id`를 반환한다.
- backend/services/gateway/internal/handler/e2e_test.go
  - trace event E2E 테스트에서 persisted audit status, audit entry ID, audit metadata, `primary_value` 부재를 검증하도록 갱신했다.
- docker-compose.yml
  - Gateway용 `AUDIT_INTAKE_URL`과 audit-service 의존성을 추가했다.
- infrastructure/docker/docker-compose.dev.yml
  - dev audit-service 컨테이너, Gateway audit intake env, Gateway 의존성, `39-audit.sql` init 마운트를 추가했다.
- infrastructure/kubernetes/base/config/configmap.yaml
  - `AUDIT_SERVICE_ADDR`, `AUDIT_SERVICE_URL`, `AUDIT_INTAKE_URL`을 추가했다.
- docs/audit/measure-trace-audit-persistence.md
  - 영속화 계약, PHI 최소화 규칙, 인프라 연결, 품질 게이트 증적을 문서화했다.
**quality_gate**:
- `go test ./backend/services/audit-service/... ./backend/services/measurement-service/... ./backend/services/gateway/...`: PASS
- `docker compose -f docker-compose.yml config --quiet`: PASS
- `docker compose -f infrastructure/docker/docker-compose.dev.yml config --quiet`: PASS
- `git diff --check` for touched files: PASS
**notes**:
- `kubectl kustomize infrastructure/kubernetes/base`는 이번 변경과 무관한 기존 namespace transformation conflict로 BLOCKED 상태다.
- ConfigMap client dry-run은 `localhost:8080`의 로컬 Kubernetes API가 없어 완료하지 못했다.
**next_steps**:
- `docs/audit/mock-retirement-register.md`의 Measure 항목을 native Rust/DB 연결 상태와 대조하고 남은 mock/stub 항목을 퇴역 처리한다.
---

##_2026-05-01_Codex_measure_remote_trace_sink_gate
**status**: done
**root_cause**:
- Measure golden path trace events were only emitted to local `AppLogger`; operations needed a remote Gateway intake for phase-level failure triage.
- Because measurement values are PHI-sensitive, the remote observability payload needed to preserve execution context without shipping `primary_value`.
**changes**:
- frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart
  - Added `toRemoteObservabilityJson()` with `schema_version=measure_trace.v1`.
  - Remote payload redacts `primary_value` and only sends `has_primary_value`.
  - Added `MeasurementGoldenPathCompositeTraceSink`.
- frontend/flutter-app/lib/features/measurement/data/measurement_trace_sink_rest.dart
  - Added fire-and-forget REST trace sink that posts phase events to Gateway and logs send failures without blocking measurement execution.
- frontend/flutter-app/lib/core/services/rest_client.dart
  - Added `recordMeasurementTraceEvent()` for `POST /measurements/trace-events`.
- frontend/flutter-app/lib/core/providers/grpc_provider.dart
  - Added `measurementGoldenPathTraceSinkProvider` combining local logger and remote REST sink.
- frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart
  - Wired Measure CTA orchestrator to the composite trace sink provider.
- backend/services/gateway/internal/handler/measurement_routes.go
  - Added `POST /api/v1/measurements/trace-events`.
  - Validates `phase`, rejects negative `elapsed_ms`, and rejects any remote payload containing `primary_value`.
  - Returns `202 Accepted` with `audit_status` until audit-service write RPC exists.
- backend/services/gateway/internal/handler/e2e_test.go
  - Added trace event accept test and `primary_value` rejection test.
- frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
  - Verifies remote observability JSON redacts primary value.
- frontend/flutter-app/test/features/measurement/data/measurement_trace_sink_rest_test.dart
  - Verifies REST sink posts to `/api/v1/measurements/trace-events` with PHI-minimized payload.
- docs/audit/measure-golden-path-remote-trace-sink.md
  - Added remote trace sink contract, PHI minimization rule, quality gates, and environment note.
**quality_gate**:
- `gofmt -w backend/services/gateway/internal/handler/measurement_routes.go backend/services/gateway/internal/handler/e2e_test.go`: PASS
- `go test ./backend/services/measurement-service/... ./backend/services/gateway/...`: PASS
- `dart format` changed Dart files: PASS
- `flutter test --no-pub test/features/measurement/application/measurement_golden_path_orchestrator_test.dart test/features/measurement/data/measurement_trace_sink_rest_test.dart test/features/measurement/data/measurement_repository_rest_test.dart`: PASS
- `flutter analyze --no-pub` changed Measure/Gateway client files: PASS
**notes**:
- Plain `flutter pub get` currently fails because the WSL Flutter SDK cache is missing `bin/cache/flutter.version.json`; existing package config allowed `--no-pub` compile/test/analyze gates to pass.
**next_steps**:
- Add audit-service write RPC or a dedicated audit intake so Gateway `audit_status` can become persisted instead of deferred.
- Reconcile Measure rows in `docs/audit/mock-retirement-register.md` against native Rust and DB wiring.
---

##_2026-05-01_Codex_measurement_stream_contract_gate
**status**: done
**root_cause**:
- Gateway `/api/v1/measurements/process` E2E contract was covered, but the measurement-service `StreamMeasurement` handler still needed a direct unit gate proving the stream frame reaches measurement storage and fingerprint vector storage.
**changes**:
- backend/services/measurement-service/internal/handler/grpc_stream_test.go
  - Added a fake bidirectional gRPC stream for generated `MeasurementService_StreamMeasurementServer`.
  - Added `TestStreamMeasurementStoresMeasurementAndFingerprint` to verify active session metadata backfill, differential value mapping, environment metadata mapping, raw channel retention, fingerprint vector persistence, response fingerprint, and `processed_at`.
  - Added `TestStreamMeasurementRequiresSessionID` to verify empty `session_id` is rejected with `InvalidArgument` before store/vector side effects.
- docs/audit/measurement-service-stream-contract.md
  - Added StreamMeasurement persistence/vector contract, test evidence, and next-step handoff.
- CONTEXT.md
  - Added latest backend stream contract gate status and next steps.
**quality_gate**:
- `gofmt -w backend/services/measurement-service/internal/handler/grpc_stream_test.go`: PASS
- `go test ./backend/services/measurement-service/... ./backend/services/gateway/...`: PASS
**next_steps**:
- Forward Flutter Measure golden path trace events to remote observability or audit sink.
- Reconcile the Measure rows in `docs/audit/mock-retirement-register.md` against real native Rust and DB wiring.
---

##_2026-05-01_Codex_go_toolchain_and_measure_gateway_contract_gate
**status**: done
**root_cause**:
- Go quality gate was blocked because WSL had no Go binary and Windows Go could not reliably lock modules through the WSL UNC path.
- After Go was installed, the real gate exposed one measurement-service compile blocker and gateway E2E tests that still expected camelCase despite the REST proto JSON contract now using snake_case.
**changes**:
- Installed Go in WSL user space:
  - version: `go1.26.2 linux/amd64`
  - path: `/home/kangjh3kang/sdk/go-go1.26.2`
  - archive SHA256: `990e6b4bbba816dc3ee129eaeaf4b42f17c2800b88a2166c265ac1a200262282`
- backend/services/measurement-service/internal/service/measurement.go
  - Fixed malformed comment line that swallowed the `observations` declaration and caused `undefined: observations`.
- backend/services/gateway/internal/handler/e2e_test.go
  - Added mock `StreamMeasurement` client stream.
  - Added `/api/v1/measurements/process` E2E contract test proving REST payload is sent into gRPC `StreamMeasurement`.
  - Updated JSON key assertion helper to accept snake_case proto JSON while preserving existing camelCase test calls.
- docs/audit/go-toolchain-wsl-test-gate.md
  - Added install path, invocation, checksum, and quality gate evidence.
**quality_gate**:
- `go test ./backend/services/measurement-service/... ./backend/services/gateway/...`: PASS
**next_steps**:
- Add measurement-service `StreamMeasurement` unit test with store/vector assertions.
- Forward Flutter Measure trace events to remote observability or audit sink.
---

##_2026-05-01_Codex_measure_golden_path_observability_trace
**status**: done
**root_cause**:
- Measure golden path had execution and readiness gating, but each operational phase needed deterministic trace events so failures can be diagnosed by phase and proven in tests.
**changes**:
- frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart
  - Added `MeasurementGoldenPathTraceEvent` with phase, elapsed time, timestamp, session/cartridge IDs, engine mode, value/unit/confidence, failure reason, and diagnostic message.
  - Added injectable `MeasurementGoldenPathTraceSink`.
  - Added default `MeasurementGoldenPathLogger` bridge to `AppLogger` with warning level for failed events.
  - Orchestrator now emits trace events for readiness, cartridge, session, engine, server, end, and failed phases.
- frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
  - Happy path now verifies trace event phase sequence, engine mode, session ID, JSON payload, and log line.
  - Blocked readiness path now verifies failed trace event and failure reason.
- docs/audit/measure-golden-path-observability.md
  - Added phase trace contract and verification evidence.
**quality_gate**:
- dart format changed files: PASS
- flutter analyze changed observability files: PASS
- flutter test measurement golden path + REST repository + domain suite: PASS
- go test: NOT RERUN; previous environment blocker remains.
**next_steps**:
- Add gateway `/measurements/process` contract test when Go runtime is available in WSL or a non-UNC workspace.
- Forward Measure trace events to remote observability/audit sink after local logger contract stabilizes.
---

##_2026-05-01_Codex_measure_golden_path_readiness_gate
**status**: done
**root_cause**:
- Measure golden path had UI/Rust/repository/server wiring, but the runtime still needed an explicit harness gate so a release build cannot silently continue through Dart stub measurement when the native Rust engine is unavailable.
**changes**:
- frontend/flutter-app/lib/core/services/rust_ffi_stub.dart
  - Added `RustBridgeDiagnostics` with initialized/native/platform/release/version/release-stub-allowance state.
  - Added release stub policy via `MANPASIK_RELEASE_ALLOW_STUB_MEASUREMENT`.
- frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart
  - Added `readinessChecked` phase, `MeasurementEngineReadiness`, and `MeasurementGoldenPathReadinessException`.
  - `RustMeasurementEngine` now initializes/checks `RustBridge` before NFC/session/pipeline/server work.
  - Orchestrator blocks before `startSession` when readiness fails.
- frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart
  - Status text now surfaces the engine mode after readiness check.
- frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
  - Updated happy path phase trace.
  - Added blocked readiness test proving no session/process/end calls happen when the engine is not ready.
- docs/audit/measure-golden-path-readiness-gate.md
  - Added executable evidence note for MR-002 / Measure H1 readiness hardening.
**quality_gate**:
- dart format changed files: PASS
- flutter analyze changed Measure/Rust bridge files: PASS
- flutter test measurement golden path + REST repository + domain suite: PASS
- go test: NOT RERUN in this step; previous environment blocker remains Windows Go on WSL UNC lock and no WSL Go binary.
**next_steps**:
- Add gateway `/measurements/process` contract test when Go can run in native Linux or non-UNC workspace.
- Add lightweight observability events for readiness/cartridge/session/engine/server/end phases.
---

##_2026-05-01_Codex_harness_h0_and_measure_golden_path_foundation
**status**: done
**root_cause**:
- Sprint H0 산출물(harness-capillary-map, mock-retirement-register, menu-completion-scorecard)이 필요했고, Measure 골든 패스가 UI/Rust/Repository/Gateway/measurement-service 저장 경로로 분리되어 있지 않았음.
**changes**:
- docs/audit/harness-capillary-map.md
  - 핵심 메뉴별 Route -> Provider/Repository -> API/RPC -> Service -> Storage/Event 연결 맵 작성.
- docs/audit/mock-retirement-register.md
  - demo/mock/stub/fallback 지점을 P0/P1/P2로 분류하고 운영 은퇴 전략 정의.
- docs/audit/menu-completion-scorecard.md
  - 핵심 메뉴별 H0 완성도 점수와 H1/H2 보강 액션 정의.
- frontend/flutter-app/lib/features/measurement/application/measurement_golden_path_orchestrator.dart
  - NFC cartridge read, session start, Rust pipeline, backend process, session end를 묶는 MeasurementGoldenPathOrchestrator 추가.
- frontend/flutter-app/lib/features/measurement/domain/measurement_repository.dart
  - ProcessMeasurementRequest/Result와 processMeasurement 계약 추가.
- frontend/flutter-app/lib/core/services/rest_client.dart
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_rest.dart
- frontend/flutter-app/lib/features/measurement/data/measurement_repository_impl.dart
  - REST process endpoint 결선 및 native/gRPC binding 미지원 환경의 안전한 처리 결과 반환 추가.
- backend/services/gateway/internal/handler/measurement_routes.go
  - POST /api/v1/measurements/process 추가. REST 요청을 gRPC StreamMeasurement로 전달.
- backend/services/gateway/internal/handler/rest_handler.go
  - proto JSON 응답을 snake_case로 직렬화하도록 UseProtoNames=true 적용.
- backend/services/measurement-service/internal/handler/grpc.go
  - StreamMeasurement 구현. 수신 프레임을 service.ProcessMeasurement로 저장/처리.
- backend/services/measurement-service/internal/service/measurement.go
  - ProcessMeasurement에서 session active 검증, session 기반 device/user/cartridge 보강, 저장/벡터 저장 경로 정리.
- frontend/flutter-app/test/features/measurement/application/measurement_golden_path_orchestrator_test.dart
  - 골든 패스 단계 순서와 process payload 생성 검증 테스트 추가.
- frontend/flutter-app/test/helpers/fake_repositories.dart
  - MeasurementRepository 신규 계약 구현.
**quality_gate**:
- dart format changed measurement files: PASS
- gofmt changed Go files: PASS
- flutter test test/features/measurement/application/measurement_golden_path_orchestrator_test.dart: PASS
- flutter analyze changed measurement files: PASS
- go test: BLOCKED in current shell because Windows Go cannot lock WSL UNC go.mod/go.work files and WSL has no go binary installed.
**next_steps**:
- /measure 화면에 MeasurementGoldenPathOrchestrator를 실제 CTA로 연결하고, 운영 빌드에서 Rust native 미초기화 시 silent stub fallback을 차단하는 진단 UI를 추가.
---

##_2026-05-01_Codex_harness_engineering_platform_upgrade_plan
**status**: done
**root_cause**:
- 하네스 엔지니어링 관점의 1차 분석에서 UI/메뉴 범위 대비 실제 운용 연결성, 측정 골든 패스, readiness, mock/stub 은퇴, 규제 증적 자동화가 별도 실행 계획으로 관리되어야 함을 확인.
- 완성도 보강 작업이 화면 단위가 아니라 Route -> Provider -> Repository -> API -> Service -> Storage/Event -> Observability -> Test Evidence로 이어지는 혈관형 연결성 기준으로 재정의될 필요가 있음.
**changes**:
- docs/plan/harness-engineering-platform-upgrade-plan-2026-05-01.md
  - 만파식 플랫폼 완성/업그레이드를 위한 하네스 엔지니어링 상세 구현 계획 추가.
  - 8계층 하네스 모델, 메뉴별 보강 계획, Rust/Flutter/Go/Data/Ops/Compliance 구축 계획, Sprint H0~H5 실행안, Definition of Done 정의.
- CONTEXT.md
  - 최근 업데이트 상단에 하네스 엔지니어링 업그레이드 계획 수립 내역 반영.
**quality_gate**:
- 문서 변경만 수행. 빌드/테스트 실행 대상 없음.
**next_steps**:
- Sprint H0 산출물인 harness-capillary-map, mock-retirement-register, menu-completion-scorecard를 생성하고 Measure 골든 패스 구현으로 진입.
---

##_2026-02-23_Codex_navigation_storyline_medical_hub_rearchitecture_and_gate_revalidation
**status**: done
**root_cause**:
- 모바일에서 좌측 내비게이션 대체 수단이 부족해 주요 페이지 전환 접근성이 떨어졌고, 일부 경로에서 전역/로컬 상하단 크롬 표시 조건이 불안정해 디자인 일관성이 깨졌다.
- 의료 플로우가 단일 통화 진입 중심으로 단순화되면서 기관 선택/예약/약국 연계/데이터 공유의 실제 사용자 여정이 분절됐다.
- pop 불가 상황에서 화면별 뒤로가기 처리 방식이 제각각이라 홈 복귀 실패 케이스가 남아 있었다.
**changes**:
- frontend/flutter-app/lib/shared/utils/navigation_utils.dart
  - popOrHome() 공통 fallback 유틸 추가.
- frontend/flutter-app/lib/core/router/app_router.dart
  - 전역 크롬(상단/하단) 표시 로직을 경로 정책 중심으로 통합.
- frontend/flutter-app/lib/core/router/bottom_nav_visibility.dart
  - /medical/video-call 구간의 글로벌 크롬 숨김 정책 추가.
- frontend/flutter-app/lib/shared/widgets/mobile_scaffold.dart
  - 모바일 Drawer 기반 전체 메뉴 진입점 추가.
- frontend/flutter-app/lib/features/home/presentation/home_screen.dart
  - 전체 메뉴 바텀시트 및 이동 동선 정리(go 중심).
- frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart
  - 뒤로가기 fallback 및 메뉴 바텀시트 추가.
- frontend/flutter-app/lib/features/medical/presentation/medical_screen.dart
  - 의료 허브를 기관선택→예약→진료→약국전송 흐름으로 재구성.
- frontend/flutter-app/lib/features/medical/presentation/telemedicine_screen.dart
- frontend/flutter-app/lib/features/medical/presentation/facility_search_screen.dart
- frontend/flutter-app/lib/features/medical/presentation/prescription_detail_screen.dart
- frontend/flutter-app/lib/features/medical/presentation/consultation_result_screen.dart
  - 뒤로가기 동작을 popOrHome()로 통일.
- frontend/flutter-app/lib/shared/widgets/sanggam_orbit_frame.dart
- frontend/flutter-app/lib/shared/widgets/glass_dock_navigation.dart
- frontend/flutter-app/lib/features/home/presentation/home_screen.dart
- frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart
  - 상/하단 이미지 오버레이/fit 조정으로 흰 박스형 노출 리스크 완화.
- frontend/flutter-app/test/core/router/bottom_nav_visibility_test.dart
  - /medical/video-call/* 경로 케이스 포함 정책 테스트 보강.
- docs/reports/ui_storyline_audit_2026-02-23.md
  - 스토리라인/연결성/의료 여정 재점검 보고서 추가.
**quality_gate**:
- flutter analyze (변경 파일 대상): PASS
- flutter test test/core/router/bottom_nav_visibility_test.dart: PASS
- flutter build linux --debug: PASS
**notes**:
- webview_cef CMake 경고는 기존과 동일한 non-blocking 상태.
---
﻿##_2026-02-23_Codex_home_white_box_height_mismatch_hard_fix_and_gate_pass
**status**: done
**root_cause**:
- `/home` 및 `/measure`에서 라우트 문자열 판별이 어긋나면 전역/로컬 크롬이 중첩되어 상단/하단 높이 불일치 및 흰 박스 노출이 재발 가능.
- 홈 화면 상/하단 높이가 고정값(`36/100`) 기반이라 화면 크기 대비 비율이 깨지는 케이스 존재.
- 상/하단 레이어가 이미지 실패/밝은 소스에 취약해 배경 노출 시 시각적 흰 박스로 보일 수 있음.
**changes**:
- `frontend/flutter-app/lib/core/router/app_router.dart`
  - `navigationShell.currentIndex` 기반 강제 숨김 조건 추가(0: home, 2: measure).
  - 라우트 prefix 정책과 병행해 전역 상단/하단 크롬 중첩 방지.
- `frontend/flutter-app/lib/features/home/presentation/home_screen.dart`
  - 상단/하단 높이를 고정값에서 비율 기반(`header 9%`, `dock 14%`, min/max clamp)으로 전환.
  - 상/하단 레이어에 불투명 다크 베이스 추가로 흰 박스 노출 차단.
  - 하단 액션 버튼/네비게이션 배치를 컨텐츠 높이 기반으로 재계산.
- `frontend/flutter-app/lib/shared/widgets/mobile_scaffold.dart`
  - `hideGlobalTopFrame || hideGlobalDock` 시 다크 쉘 배경 강제.
- `frontend/flutter-app/lib/shared/widgets/desktop_scaffold.dart`
  - `hideGlobalTopFrame` 시 다크 쉘 배경 강제.
- `frontend/flutter-app/lib/features/data_hub/presentation/data_hub_screen.dart`
  - 빌드 차단 문법 오류(여분 bracket block) 제거.
  - `Colors.white05`를 유효 색상 상수(`Color(0x0DFFFFFF)`)로 교체.
**quality_gate**:
- `dart format` (changed files): PASS
- `flutter analyze` (changed files): PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- `webview_cef` CMake 경고는 비차단(non-blocking) 상태이며 기존과 동일.
---
##_2026-02-23_Codex_white_header_dock_box_root_cause_fix_route_chrome_and_asset_tint
**status**: done
**root_cause**:
- Top/Bottom white block issue was caused by two factors:
  1) Double chrome stacking: global shell frame/dock + local screen frame/dock rendered together on `/home` and `/measure`.
  2) Slim asset PNG files are RGB (no alpha), so direct full-bleed rendering exposed bright/solid backgrounds.
- Height mismatch was amplified by global top frame fixed height (`120`) over local header.
**changes**:
- `frontend/flutter-app/lib/core/router/bottom_nav_visibility.dart`
  - Added `/home` to global dock hide policy.
  - Added global top frame hide policy (`shouldHideGlobalTopFrame`) for `/home`, `/measure`.
  - Refactored prefix matching into shared helper.
- `frontend/flutter-app/lib/core/router/app_router.dart`
  - Added route-aware top frame gating (`hideGlobalTopFrameOnRoute`).
  - Applied route-aware shell composition to prevent duplicate top chrome.
  - Set scaffold fallback background color to deep sea tone.
- `frontend/flutter-app/lib/shared/widgets/sanggam_orbit_frame.dart`
  - Reduced top frame height `120 -> 84`.
  - Replaced direct bright overlay with dark gradient + multiply color filter + subtle gold separator line.
- `frontend/flutter-app/lib/shared/widgets/glass_dock_navigation.dart`
  - Replaced bright `screen + white` dock asset blend with `multiply` tint + controlled opacity.
  - Switched dock image fitting to width-based alignment.
- `frontend/flutter-app/lib/features/home/presentation/home_screen.dart`
  - Top header and bottom dock image layers converted to dark base + multiply tinted asset rendering.
  - Prevents white box appearance when source PNG lacks transparency.
- `frontend/flutter-app/test/core/router/bottom_nav_visibility_test.dart`
  - Rewritten tests for dock+top-frame policies, child-route coverage, and prefix-lookalike guards.
**quality_gate**:
- `dart format` (changed files): PASS
- `flutter analyze` (changed files): PASS
- `flutter test test/core/router/bottom_nav_visibility_test.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- Non-blocking `webview_cef` CMake warnings remain unchanged.
- Additional recommendation: regenerate/export true alpha PNG assets to simplify rendering stack.
---
##_2026-02-23_Codex_measure_design_mode_persistence_shared_preferences
**status**: done
**root_cause**:
- User requested immediate persistence of 3-mode design selection across app restart.
- Current implementation changed mode only in-memory.
**changes**:
- Updated `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`.
- Added `shared_preferences` integration:
  - restore on init via `_restoreSavedDesignMode()`
  - persist on mode switch via `_saveDesignMode()`
  - pref key: `measure_design_mode`
- Added mode parser `_modeFromStorage()` with safe default fallback (`royal_orbit`).
- Added `dart:async` and `unawaited()` for non-blocking preference save.
**quality_gate**:
- `dart format frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`: PASS
- `flutter analyze lib/features/measurement/presentation/measure_screen.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- Existing `webview_cef` CMake warnings are non-blocking and unchanged.
---
##_2026-02-23_Codex_measure_screen_three_design_modes_runtime_switch
**status**: done
**root_cause**:
- User requested applying all three design concepts as selectable modes.
- Existing measure screen used single fixed asset paths only.
**changes**:
- Updated `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`.
- Added `_SlimDesignMode` enum (`royalOrbit`, `jagaeWave`, `hanjiGold`) and mapping extension.
- Added runtime mode state in `MeasureScreen` and selection handler.
- Added top-bar mode selector popup (`_DesignModeMenu`) with Korean labels.
- Switched header/bottom/button assets to mode-based paths:
  - `assets/images/concepts/<mode>/header_3d_frame_slim.png`
  - `assets/images/concepts/<mode>/bottom_3d_dock_slim.png`
  - `assets/images/concepts/<mode>/btn_3d_action_slim.png`
- Added fallback to default slim assets in `assets/images/` on load failure.
**quality_gate**:
- `dart format measure_screen.dart`: PASS
- `flutter analyze lib/features/measurement/presentation/measure_screen.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- Existing non-blocking `webview_cef` CMake warnings remain unchanged.
---
##_2026-02-23_Codex_strict_prompt_v2_all_concepts_regenerated
**status**: done
**root_cause**:
- User selected option 3 to regenerate all concepts with stricter prompt constraints.
- Previous generations still mixed product-scene perspective artifacts.
**changes**:
- Updated `scripts/generate_design_concepts.py` with stricter composition controls.
- Added `STRICT_GLOBAL` and `STRICT_BY_KIND` rules for header/bottom/button-only generation.
- Added per-asset output sizes via `ASSET_TARGETS`:
  - `header_3d_frame_slim.png`: `1792x1024`
  - `bottom_3d_dock_slim.png`: `1792x1024`
  - `btn_3d_action_slim.png`: `1024x1024`
- Regenerated all concept sets:
  - `royal_orbit`
  - `jagae_wave`
  - `hanji_gold`
**quality_gate**:
- `python3 -m py_compile scripts/generate_design_concepts.py scripts/apply_design_concept.py`: PASS
- `python3 scripts/generate_design_concepts.py` (interactive shell): PASS
**notes**:
- Newly regenerated images are available under `frontend/flutter-app/assets/images/concepts/<concept>/`.
- Apply command remains `python3 scripts/apply_design_concept.py --concept <name>`.
---
##_2026-02-23_Codex_bottom_nav_visibility_policy_refactor_and_route_test_added
**status**: done
**root_cause**:
- `/measure`?먯꽌 ?꾩뿭 ?섎떒 ?ㅻ퉬? 痢≪젙 ?꾩슜 ?섎떒 ?덉씠?닿? ?숈떆???뚮뜑?섏뼱 以묐났/寃뱀묠 諛쒖깮.
- 寃쎈줈 議곌굔???쇱슦???뚯씪???섎뱶肄붾뵫?섏뼱 ?뺤옣/?뚭? ??묒씠 ?대젮?좎쓬.
**changes**:
- `frontend/flutter-app/lib/core/router/bottom_nav_visibility.dart` 異붽?
  - `kHideGlobalDockRoutePrefixes` ?곸닔 湲곕컲 ?뺤콉 ?꾩엯.
  - `shouldHideGlobalDock()` 寃쎈줈 ?먮퀎 ?⑥닔 ?꾩엯 (`exact or prefix/child`).
- `frontend/flutter-app/lib/core/router/app_router.dart` ?섏젙
  - ?꾩뿭 ?섎떒 ?ㅻ퉬 ?몄텧??`shouldHideGlobalDock(location)` ?뺤콉?쇰줈 遺꾨━.
- `frontend/flutter-app/test/core/router/bottom_nav_visibility_test.dart` 異붽?
  - `/measure`, `/measure/*`, 鍮꾩쑀??寃쎈줈(`/measured`, `/measurements`) 寃쎄퀎議곌굔 ?ы븿 ?⑥쐞 ?뚯뒪??
**quality_gate**:
- `dart format` (router/policy/test): PASS
- `flutter analyze` (router/policy/test): PASS
- `flutter test test/core/router/bottom_nav_visibility_test.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
---

##_2026-02-23_Codex_design_concept_assets_multi_set_generation_for_selection
**status**: done
**root_cause**:
- ?곷떒/?섎떒 ?붿옄??寃곌낵 ?몄감媛 而ㅼ꽌 ?⑥씪 ?먯뀑留뚯쑝濡?諛⑺뼢 寃곗젙???섍린 ?대젮?좎쓬.
**changes**:
- `scripts/generate_design_concepts.py` 異붽?
  - 而⑥뀎 3醫?`royal_orbit`, `jagae_wave`, `hanji_gold`)???곷떒/?섎떒/踰꾪듉 ?먯뀑 ?쇨큵 ?앹꽦 濡쒖쭅 援ы쁽.
  - 異쒕젰 寃쎈줈: `frontend/flutter-app/assets/images/concepts/<concept>/...`
- `scripts/apply_design_concept.py` 異붽?
  - ?좏깮??而⑥뀎 ?명듃瑜??쒖꽦 ?먯뀑(`assets/images/*.png`)?쇰줈 利됱떆 諛섏쁺?섎뒗 ?곸슜 ?ㅽ겕由쏀듃 援ы쁽.
- ?앹꽦 ?ㅽ뻾 ?꾨즺(Interactive shell):
  - `header_3d_frame_slim.png`, `bottom_3d_dock_slim.png`, `btn_3d_action_slim.png` 횞 3而⑥뀎
**quality_gate**:
- `python3 -m py_compile scripts/generate_design_concepts.py`: PASS
- `python3 -m py_compile scripts/apply_design_concept.py`: PASS
- `python3 scripts/generate_design_concepts.py` (interactive shell): PASS
**notes**:
- ?앹꽦 ?대?吏????붿뿉 泥⑤??섏뼱 而⑥뀎 ?좏깮 媛???곹깭.
---

##_2026-02-23_Codex_global_bottom_nav_hidden_on_measure_route_to_prevent_overlap
**status**: done
**root_cause**:
- `/measure`??Shell ?쇱슦???대????꾩뿭 `GlassDockNavigation`???쒖떆?섍퀬,
- `MeasureScreen`???먯껜 ?섎떒 ?꾪겕(`_SlimBottomLayer`)瑜??뚮뜑留곹빐 ?섎떒 UI媛 以묐났/寃뱀묠.
**changes**:
- `frontend/flutter-app/lib/core/router/app_router.dart`
  - `ScaffoldWithBottomNav`?먯꽌 ?꾩옱 ?쇱슦??`matchedLocation`)媛 `/measure`濡??쒖옉?섎㈃ ?꾩뿭 `bottomNavigationBar`瑜??④린?꾨줉 議곌굔 異붽?.
  - ?뺤쟻遺꾩꽍 寃쎄퀬 ?쒓굅瑜??꾪빐 誘몄궗??import/誘몄궗????誘몄궗??濡쒖쭅 ?뺣━.
**quality_gate**:
- `dart format app_router.dart`: PASS
- `flutter analyze app_router.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- 寃곌낵?곸쑝濡?痢≪젙 ?붾㈃?먯꽌??痢≪젙 ?꾩슜 ?섎떒 UI留??④퀬, 以묐났 ?ㅻ퉬寃뚯씠??寃뱀묠???쒓굅??
---

##_2026-02-23_Codex_slim_fit_assets_regeneration_rerun_success_interactive_shell
**status**: done
**root_cause**:
- 鍮꾨??뷀삎 ??`bash -lc`)?먯꽌??`OPENAI_API_KEY`媛 濡쒕뱶?섏? ?딆븘 ?먯뀑 ?ъ깮?깆뿉 ?ㅽ뙣?덉쓬.
**changes**:
- `bash -ic` ?ㅽ뻾 寃쎈줈濡?`scripts/generate_assets.py` ?ъ떎??
- ?꾨옒 3媛??ㅼ궗 ?먯뀑 ?ъ깮???꾨즺:
  - `frontend/flutter-app/assets/images/header_3d_frame_slim.png`
  - `frontend/flutter-app/assets/images/bottom_3d_dock_slim.png`
  - `frontend/flutter-app/assets/images/btn_3d_action_slim.png`
**quality_gate**:
- `python3 scripts/generate_assets.py` (interactive shell): PASS
- ?앹꽦 ?뚯씪 ??꾩뒪?ы봽/?⑸웾 ?뺤씤: PASS
**notes**:
- ?꾩옱 ?섍꼍?먯꽌???대?吏 ?앹꽦 ??`bash -ic` ?ㅽ뻾???덉젙??
---

##_2026-02-23_Codex_measure_screen_bottom_premium_gold_emboss_and_traditional_pattern
**status**: done_with_asset_regen_blocked_by_missing_api_key
**root_cause**:
- ?ъ슜???쇰뱶諛?湲곗??쇰줈 ?섎떒 ?덉씠?닿? ?⑥닚??怨좉툒媛먯씠 遺議깊뻽怨? ?꾩씠肄?諛곌꼍???뚯옱媛??묎컖 怨⑤뱶/?꾪넻臾몄뼇)??遺議깊뻽??
**changes**:
- `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
  - ?섎떒 諛곌꼍???꾪넻臾몄뼇 ?ㅻ쾭?덉씠(`_TraditionalPatternPainter`) 異붽?.
  - ?섎떒 ?곷떒 寃쎄퀎 怨⑤뱶 ?섏씠?쇱씠???쇱씤 異붽?.
  - ?섎떒/?곷떒 諛곗????묎컖 怨⑤뱶 ?꾩씠肄?`_EmbossedGoldIcon`) ?곸슜.
  - ?섎떒 ?≪뀡 踰꾪듉 以묒븰 ?꾩씠肄섏쓣 ?묎컖 怨⑤뱶 ?ㅽ??쇰줈 蹂寃?
  - 諛곗? ?쒓컖 ?ㅽ??쇱쓣 ?꾨━誘몄뾼 ??洹몃씪?곗씠???ｌ?/洹몃┝???쇰줈 蹂닿컯.
- `scripts/generate_assets.py`
  - ?섎떒 ?꾪겕 ?꾨＼?꾪듃???꾪넻臾몄뼇 誘몄꽭 媛곸씤 議곌굔 異붽?.
  - ?≪뀡 踰꾪듉 ?꾨＼?꾪듃???묎컖 怨⑤뱶 由대━??吏덇컧 議곌굔 異붽?.
**quality_gate**:
- `dart format measure_screen.dart`: PASS
- `python3 -m py_compile scripts/generate_assets.py`: PASS
- `flutter analyze measure_screen.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
- `python3 scripts/generate_assets.py`: FAIL (`OPENAI_API_KEY` 誘몄꽕??
**notes**:
- ?ㅼ궗 ?먯뀑 ?먮룞 ?ъ깮?깆? API ???먮뒗 fallback API ?ㅼ젙 ??利됱떆 ?ъ떎??媛??
---

##_2026-02-23_Codex_measure_screen_bottom_layer_70_percent_scale_for_readability
**status**: done
**root_cause**:
- ?ъ슜???쇰뱶諛?湲곗??쇰줈 ?섎떒 ?덉씠???믪씠媛 ?곷??곸쑝濡?而ㅼ꽌 以묒븰 ?곗씠??媛???곸뿭??以꾩뼱?쒕뒗 臾몄젣媛 ?덉뿀??
**changes**:
- `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
  - `_kDockVisualScale = 0.70` 異붽?(湲곗〈 ?섎떒 ?믪씠 ?鍮?70% ?ㅼ???.
  - `dockBodyHeight` 怨꾩궛?앹쓣 ?ㅼ???諛섏쁺 ?뺥깭濡?蹂寃?
  - ?섎떒 理쒖냼/理쒕? ?믪씠瑜?`64.0~96.0` 踰붿쐞濡?議곗젙.
**quality_gate**:
- `dart format measure_screen.dart`: PASS
- `flutter analyze measure_screen.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- `webview_cef` CMake 以묐났 ?쇱씤 寃쎄퀬??湲곗〈怨??숈씪??鍮꾩튂紐?寃쎄퀬?대ŉ 鍮뚮뱶 ?곗텧臾쇱뿉???곹뼢 ?놁쓬.
---

##_2026-02-23_Codex_slim_fit_assets_regeneration_rerun_success
**status**: done
**root_cause**:
- ?ъ슜???붿껌?쇰줈 ?щ┝??3醫??ㅼ궗 ?먯뀑 ?앹꽦 ?ㅽ겕由쏀듃 ?ъ떎?됱씠 ?꾩슂?덉쓬.
**changes**:
- `frontend/flutter-app/assets/images/header_3d_frame_slim.png` ?ъ깮??- `frontend/flutter-app/assets/images/bottom_3d_dock_slim.png` ?ъ깮??- `frontend/flutter-app/assets/images/btn_3d_action_slim.png` ?ъ깮??**quality_gate**:
- `python3 scripts/generate_assets.py`: PASS
- ?앹꽦 ?뚯씪 ??꾩뒪?ы봽 ?뺤씤: PASS
**notes**:
- ?대쾲 ?ъ떎?됱뿉??OpenAI ?대?吏 ?앹꽦???뺤긽 ?꾨즺??
---

##_2026-02-23_Codex_slim_fit_hybrid_architecture_measure_screen_rebuild
**status**: done_with_asset_generation_blocked_by_openai_billing_limit
**root_cause**:
- ?ъ슜???붿껌? ?щ┝??3醫??ㅼ궗 ?먯뀑 ?먮룞 ?앹꽦 + 痢≪젙 ?붾㈃??3???덉씠??Top 8~10 / Center Expanded / Bottom 12~15) ?ш뎄?깆씠?덉쓬.
- 湲곗〈 `scripts/generate_assets.py`??援ы삎 ?꾨＼?꾪듃(2醫?留?吏?먰뻽怨??щ┝???뚯씪紐낆쓣 ?앹꽦?섏? 紐삵뻽??
- OpenAI Images ?몄텧 ??`billing_hard_limit_reached`(HTTP 400)濡??ㅼ젣 ?좉퇋 ?대?吏 ?앹꽦??李⑤떒??
**changes**:
- `scripts/generate_assets.py`
  - ?щ┝??3媛??먯뀑(`header_3d_frame_slim.png`, `bottom_3d_dock_slim.png`, `btn_3d_action_slim.png`) ?꾨＼?꾪듃/異쒕젰 濡쒖쭅?쇰줈 援먯껜.
  - OpenAI API ?몄텧 ?ㅽ뙣 ???먮윭 蹂몃Ц??紐낇솗??異쒕젰?섎룄濡?蹂닿컯.
  - `FALLBACK_IMAGE_API_BASE` 湲곕컲 ?泥??대?吏 API 寃쎈줈 異붽?.
- `frontend/flutter-app/lib/features/measurement/presentation/measure_screen.dart`
  - ?붾㈃ 援ъ“瑜??щ┝??3???덉씠?대줈 ?꾨㈃ ?ъ옉??
  - Top Layer: `header_3d_frame_slim.png` 湲곕컲 ?뉗? ?ㅻ뜑 + 理쒖냼 ?곹깭 諭껋?.
  - Center Layer: 肄붾뱶 湲곕컲 DeepSea Blue 諛곌꼍 + `SizedBox.expand` ?꾩쁺??`CustomPaint(JagaeParticleSystem)` 1792 ?뚰떚???뚮뜑留?
  - Bottom Layer: `bottom_3d_dock_slim.png` 湲곕컲 ?щ┝ ??+ ?뚯텧??`btn_3d_action_slim.png` ?≪뀡 踰꾪듉.
  - ?ㅼ틪 ?ъ떆???명꽣?숈뀡 諛?吏꾪뻾瑜??곹깭 ?쒖떆 ?좎?.
- `frontend/flutter-app/pubspec.yaml`
  - `assets/images/header_3d_frame_slim.png`
  - `assets/images/bottom_3d_dock_slim.png`
  - `assets/images/btn_3d_action_slim.png`
  - ?깅줉 ?곹깭 ?뺤씤.
**quality_gate**:
- `python3 -m py_compile scripts/generate_assets.py`: PASS
- `python3 scripts/generate_assets.py`: FAIL (OpenAI billing hard limit reached)
- `dart format measure_screen.dart`: PASS
- `flutter analyze measure_screen.dart`: PASS
- `flutter test`: PASS
- `flutter build linux --debug`: PASS
**notes**:
- ?뚯씪? 湲곗〈 ?щ┝ ?먯뀑??`frontend/flutter-app/assets/images/`??議댁옱?섏뿬 ?고???李몄“???뺤긽.
- ?ㅼ궗 ?먯뀑 ?좉퇋 ?ъ깮?깆? OpenAI 寃곗젣 ?쒕룄 ?댁젣 ?먮뒗 ?泥?API ?ㅼ젙 ??利됱떆 ?ъ떎??媛??
---

##_2026-02-22_Codex_sanggam_orbit_image_pipeline_and_hybrid_hud_asset_integration
**status**: partial_done_blocked_by_env_key
**root_cause**:
- ?ъ슜???붿껌? DALL-E API 湲곕컲 ?ㅼ궗 ?먯뀑 ?앹꽦 ?먮룞?붿? Flutter ?섏씠釉뚮━???⑤꼸 ?곸슜?댁뿀??
- 湲곗〈 scripts/generate_assets.py???ㅼ젣 API ?몄텧 ?놁씠 ?뚮젅?댁뒪??붾쭔 湲곕줉?섎뒗 ?붾? ?뺥깭???
- ?ㅽ뻾 ?섍꼍??OPENAI_API_KEY媛 ?ㅼ젙?섏? ?딆븘 ???대?吏 ?앹꽦 ?④퀎媛 李⑤떒??
**changes**:
- scripts/generate_assets.py
  - OpenAI Images API(POST /v1/images/generations) ?ㅼ젣 ?몄텧 濡쒖쭅 援ы쁽.
  - sanggam_frame.png, jagae_button.png ?꾨＼?꾪듃? ?뚯씪紐? 1024x1024 ?앹꽦 ?ㅽ럺 異붽?.
  - b64_json ?곗꽑 泥섎━? URL fallback ?ㅼ슫濡쒕뱶 泥섎━ 異붽?.
- frontend/flutter-app/lib/shared/widgets/sanggam_hud_panel.dart
  - Stack 湲곕컲 ?섏씠釉뚮━???⑤꼸 援ы쁽.
  - Layer 1: Image.asset(assets/images/sanggam_frame.png)
  - Layer 2: BackdropFilter? ?ㅻ쾭?덉씠 洹몃씪?곗씠?? 肄섑뀗痢??덉씠??  - ?섎떒 ?≪뀡: Image.asset(assets/images/jagae_button.png) ?곗튂 ?≪뀡 ?곕룞
- frontend/flutter-app/pubspec.yaml
  - assets/images/ ?깅줉 ?곹깭 ?뺤씤 (?대? ?깅줉?섏뼱 異붽? ?섏젙 遺덊븘??
**quality_gate**:
- python3 -m py_compile scripts/generate_assets.py: PASS
- /home/kangjh3kang/flutter/bin/dart format lib/shared/widgets/sanggam_hud_panel.dart: PASS
- /home/kangjh3kang/flutter/bin/flutter analyze lib/shared/widgets/sanggam_hud_panel.dart: PASS
- python3 scripts/generate_assets.py: FAIL (OPENAI_API_KEY 誘몄꽕??
**notes**:
- ?앹꽦 ????뚯씪 assets/images/sanggam_frame.png, assets/images/jagae_button.png ??API ???ㅼ젙 ???숈씪 ?ㅽ겕由쏀듃 ?ъ떎?????앹꽦??
---
