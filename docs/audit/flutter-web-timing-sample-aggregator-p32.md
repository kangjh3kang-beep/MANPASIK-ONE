# Flutter Web Timing Sample Aggregator P32 Audit

## Scope

P32는 exported `flutter-web-timing-report` artifact 5개 이상을 offline으로 집계해 measured cost review에 필요한 median, p95, worst-case 값을 산출하는 단계다.
CI 안에서 artifact를 직접 내려받지는 않고, export된 `.env` 파일 묶음을 입력으로 받는 repo-local script와 테스트를 제공한다.

## Changes

- `scripts/flutter_web_timing_sample_aggregate_test.sh`
  - 5개 fixture artifact를 생성해 sample count, branch types, median, p95, worst-case를 검증한다.
  - `--min-samples 6` 조건에서 실패하는지 검증한다.
- `scripts/flutter_web_timing_sample_aggregate.sh`
  - `.env` artifact directory를 읽어 `latest_duration_seconds`, `branch_type`, `runner_context`를 수집한다.
  - 최소 sample 수, `pull_request` 및 `release_branch` sample 존재, runner context 존재를 검증한다.
  - sorted duration values, median, nearest-rank p95, worst-case를 key-value report로 출력한다.
- `scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - promotion policy가 aggregator script/test와 aggregate output schema를 문서화했는지 검증한다.
- `docs/ci/flutter-web-timing-threshold-promotion.md`
  - offline aggregator 사용법과 aggregate output fields를 문서화했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-sample-aggregator-p32.md`
  - P32 RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`
  - 실패: `scripts/flutter_web_timing_sample_aggregate.sh` missing
- GREEN: `bash scripts/flutter_web_timing_sample_aggregate_test.sh`
  - 통과: `FLUTTER_WEB_TIMING_SAMPLE_AGGREGATE_TEST_PASS`
- RED: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - 실패: `missing promotion marker: sample_aggregator_script: scripts/flutter_web_timing_sample_aggregate.sh`
- GREEN: `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`
  - 통과: `CI_WEB_TIMING_THRESHOLD_PROMOTION_POLICY_PASS`

## Verification

- `bash scripts/flutter_web_timing_sample_aggregate_test.sh`: PASS
- `bash scripts/ci_web_timing_threshold_promotion_policy_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash -n scripts/flutter_web_timing_sample_aggregate.sh scripts/flutter_web_timing_sample_aggregate_test.sh scripts/ci_web_timing_threshold_promotion_policy_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- Aggregator intentionally reads exported artifact files, not GitHub APIs, so it stays deterministic and testable without network access.
- Nearest-rank p95 is used for a small sample set because it avoids fractional interpolation ambiguity in shell.
- The script fails when either `pull_request` or `release_branch` samples are missing, preventing one-sided cost review.

## Residual Notes

- P33 can add a review fixture directory convention and an audit template for recording real 5-sample promotion reviews.
- Full Flutter web build was not re-run in this step because this change is an offline artifact aggregation utility and policy integration.
