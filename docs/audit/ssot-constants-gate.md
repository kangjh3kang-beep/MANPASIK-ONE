# SSOT 상수 검증 게이트

## 목적

만파식의 분석 파이프라인은 `alpha`, fingerprint dimension, classifier input/output 같은 상수가 문서와 코드 사이에서 어긋나면 측정 보정, AI 추론, 규제 문서가 동시에 흔들린다. 이번 게이트는 `ManPaSik_Tech_Spec_v2.4.3.md`를 SSOT로 삼아 Rust core와 `AGENTS.md`가 같은 값을 문서화하는지 CI에서 검증한다.

## 변경 범위

- `scripts/validate_ssot_constants.py`
  - `ALPHA_DEFAULT`와 `FINGERPRINT_DIM_MAX`를 최신 기술 스펙에서 읽는다.
  - `rust-core/manpasik-engine/src/lib.rs`의 `DEFAULT_ALPHA`, `MAX_CHANNELS`가 스펙과 일치하는지 검증한다.
  - `AGENTS.md`에 `alpha=0.98`, `MAX_CHANNELS=1792`가 반영됐는지 검증한다.
  - `alpha=0.95`, `MAX_CHANNELS=896` 같은 구형 기준이 다시 들어오면 실패한다.
- `.github/workflows/ci.yml`
  - `ssot-governance` job을 추가해 언어별 빌드 전에 상수 불일치를 잡는다.
  - Rust `cargo clippy -- -D warnings`를 blocking gate로 전환했다.
- `AGENTS.md`
  - 차동측정 기본 `alpha`를 `0.98`로 정렬했다.
  - fingerprint 확장 경로를 `88 -> 448 -> 896 -> 1792`로 정렬했다.
  - `FingerprintClassifier` 입력/출력을 `1792 -> 30`으로 정렬했다.

## 품질 게이트

- `python3 scripts/validate_ssot_constants.py`: PASS
- `.github/workflows/ci.yml`에서 `continue-on-error` 제거 확인: PASS

## 잔여 리스크

- 현재 게이트는 가장 위험한 P0 상수만 검증한다. CSI connector, manifest size, cartridge count 등은 다음 SSOT 확장 단계에서 같은 스크립트에 추가해야 한다.
- GitHub Actions YAML 파싱은 로컬 WSL 기본 환경에 PyYAML이 없어 별도 파서 실행 대신 구조적 patch와 grep 확인으로 제한했다.
