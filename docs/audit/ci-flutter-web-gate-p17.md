# CI Flutter Web Gate Audit (P17)

## 목적

P16에서 Flutter web release gate를 추가했지만, CI workflow에는 아직 연결되지 않았다. 이번 단계는 GitHub Actions `flutter-app` job에서 analyze/test 이후 release gate를 실행하도록 연결하고, workflow guard script로 이 연결을 고정한다.

## 변경 사항

- `.github/workflows/ci.yml`
  - `flutter-app` job의 `Run tests` 이후 `Flutter web release gate` step을 추가했다.
  - step은 `bash ../../scripts/flutter_web_release_gate.sh`를 실행한다.
- `scripts/ci_flutter_web_gate_workflow_test.sh`
  - CI workflow에 step name과 release gate script path가 포함되어 있는지 검증한다.

## TDD 기록

- RED: `bash scripts/ci_flutter_web_gate_workflow_test.sh`
  - 실패 이유: CI workflow에 `Flutter web release gate` step이 없음.
- GREEN: 같은 테스트 재실행 PASS.

## 자체 코드리뷰

- Flutter job의 working directory가 `frontend/flutter-app`이므로 `../../scripts/flutter_web_release_gate.sh` 경로를 사용했다.
- release gate script 자체가 repository root를 계산하므로 CI와 로컬 경로 양쪽에서 같은 정책을 적용한다.
- 기존 analyze/test 단계는 유지하고 build gate만 후속 단계로 추가했다.

## 품질 게이트

- `bash scripts/flutter_web_release_gate_test.sh`: PASS
- `bash scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh scripts/ci_flutter_web_gate_workflow_test.sh`: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS

## 다음 단계 지침

- P18에서는 CI 전체 비용을 고려해 Flutter web release gate를 PR 필수 게이트로 유지할지, nightly/release branch 전용으로 분리할지 결정한다.
- Wasm strict compatibility는 여전히 별도 제품 타깃 결정이 필요하다.
