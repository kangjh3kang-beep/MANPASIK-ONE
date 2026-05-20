# Flutter Web Timing Artifact P30 Audit

## Scope

P30은 P29의 `key_value_v1` timing report를 GitHub Actions 아티팩트로 남기도록 CI 연결을 추가한 단계다.
대상 범위는 `flutter-app` job의 성공한 `Flutter web release gate` 이후 report 생성과 업로드, 그리고 이를 검증하는 정책/매트릭스 가드다.

## Changes

- `.github/workflows/ci.yml`
  - `Flutter web timing report` step을 추가해 `/tmp/manpasik_flutter_web_build.log`에서 report를 생성한다.
  - `Upload Flutter web timing report` step을 추가해 `/tmp/manpasik_flutter_web_timing.env`를 `flutter-web-timing-report` 아티팩트로 업로드한다.
- `scripts/ci_flutter_web_timing_artifact_workflow_test.sh`
  - workflow에 report 생성, branch type, runner context, artifact upload marker가 있는지 검증한다.
- `scripts/ci_web_gate_timing_collection_policy_test.sh`
  - timing collection 정책이 CI artifact upload 계약을 명시하도록 marker를 확장했다.
- `docs/ci/flutter-web-gate-timing-collection.md`
  - artifact name/path/action과 생성/업로드 step 이름을 정책 marker로 기록했다.
- `scripts/ci_gate_matrix_policy_test.sh`
  - CI gate matrix가 timing artifact workflow guard를 참조하는지 검증한다.
- `docs/ci/ci-gate-matrix.md`
  - `Flutter web timing artifact` row와 artifact marker를 추가했다.
- `docs/superpowers/plans/2026-05-14-flutter-web-timing-artifact-p30.md`
  - RED-GREEN 구현 순서와 검증 명령을 기록했다.

## TDD Record

- RED: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`
  - 실패: `workflow missing marker: Flutter web timing report`
- GREEN: `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`
  - 통과: `CI_FLUTTER_WEB_TIMING_ARTIFACT_WORKFLOW_PASS`
- RED: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`
  - 실패: `missing timing collection marker: ci_artifact_upload: true`
- GREEN: `bash scripts/ci_web_gate_timing_collection_policy_test.sh`
  - 통과: `CI_WEB_GATE_TIMING_COLLECTION_POLICY_PASS`
- RED: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 실패: `missing matrix marker: Flutter web timing artifact`
- GREEN: `bash scripts/ci_gate_matrix_policy_test.sh`
  - 통과: `CI_GATE_MATRIX_POLICY_PASS`

## Verification

- `bash scripts/ci_flutter_web_timing_artifact_workflow_test.sh`: PASS
- `bash scripts/ci_web_gate_timing_collection_policy_test.sh`: PASS
- `bash scripts/ci_gate_matrix_policy_test.sh`: PASS
- `bash scripts/flutter_web_timing_report_test.sh`: PASS
- `bash scripts/ci_web_gate_policy_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n scripts/ci_flutter_web_timing_artifact_workflow_test.sh scripts/ci_web_gate_timing_collection_policy_test.sh scripts/ci_gate_matrix_policy_test.sh scripts/flutter_web_timing_report.sh scripts/flutter_web_timing_report_test.sh scripts/ci_web_gate_policy_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `git diff --check` targeted files: PASS
- trailing whitespace check on targeted files: PASS

## Self Review

- The artifact is generated only after the release gate step succeeds, so failed web builds do not create misleading timing samples.
- The report keeps branch type and runner context to prevent mixing pull request samples with release branch samples during cost review.
- The workflow guard checks exact script, log path, output path, upload action, artifact name, and artifact path.

## Residual Notes

- The full `flutter build web --no-pub` was not re-run in this step. This step changed CI wiring and policy guards, so verification focused on workflow/policy scripts and common release gates.
- P31 should define timing threshold promotion rules after at least 5 successful artifact samples are available.
