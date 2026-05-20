# Flutter Web Release Gate Audit (P16)

## 목적

P13-P15에서 `flutter build web --no-pub`는 반복적으로 성공했지만, Flutter의 Wasm dry-run compatibility warning이 함께 출력됐다. 이번 단계는 현재 릴리스 타깃을 JS web build artifact로 명확히 두고, Wasm compatibility warning은 non-blocking warning으로 추적하는 release gate를 추가한다.

## 변경 사항

- `scripts/flutter_web_release_gate.sh`
  - 기본 모드에서 `frontend/flutter-app`의 `flutter build web --no-pub`를 실행한다.
  - build log에 `Built build/web` marker가 없으면 실패한다.
  - `Wasm dry run findings:`가 있으면 `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`을 출력하되 실패하지 않는다.
  - `--policy-check <log>` 모드로 parser만 빠르게 검증할 수 있다.
- `scripts/flutter_web_release_gate_test.sh`
  - fixture log로 pass-with-warning과 missing-build-marker failure를 검증한다.

## TDD 기록

- RED: `bash scripts/flutter_web_release_gate_test.sh`
  - 실패 이유: `scripts/flutter_web_release_gate.sh`가 없음.
- GREEN: 같은 테스트 재실행 PASS.

## 실제 build gate 결과

- `bash scripts/flutter_web_release_gate.sh`: PASS
- 출력 marker:
  - `FLUTTER_WEB_RELEASE_GATE_PASS`
  - `FLUTTER_WEB_RELEASE_GATE_WARN_WASM_DRY_RUN`

## 자체 코드리뷰

- JS web build artifact 성공 여부는 release blocking으로 유지했다.
- Wasm dry-run warning은 현재 명시적으로 추적하되 non-blocking으로 분리했다.
- 향후 Wasm을 공식 릴리스 타깃으로 채택하면 strict Wasm gate를 별도 추가해야 한다.

## 품질 게이트

- `bash scripts/flutter_web_release_gate_test.sh`: PASS
- `bash scripts/flutter_web_release_gate.sh`: PASS with Wasm warning marker
- `flutter test --no-pub` targeted DataHub layout/evidence/domain/REST and shared badge tests: PASS
- `python3 scripts/validate_ssot_constants.py`: PASS
- `bash scripts/security_release_gate.sh`: PASS
- `bash scripts/assay_evidence_gate.sh`: PASS
- `bash -n scripts/flutter_web_release_gate.sh scripts/flutter_web_release_gate_test.sh`: PASS

## 다음 단계 지침

- P17에서는 새 Flutter web release gate를 CI workflow에 연결할지 검토한다.
- Wasm 지원을 제품 요구사항으로 승격하려면 `flutter_secure_storage_web`, `share_plus`, `connectivity_plus`, `js` 사용 경로를 대체/분기하는 별도 호환성 계획이 필요하다.
