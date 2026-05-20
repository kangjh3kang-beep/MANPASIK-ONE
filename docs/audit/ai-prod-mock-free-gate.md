# AI Production Mock-Free Gate

## 목적

Rust AI 추론 엔진은 모델 파일 로드 여부와 무관하게 시뮬레이션 추론으로 폴백할 수 있었다. 이는 데모와 단위 테스트에는 유용하지만, 운영/인허가 경로에서는 AI 성능을 실제 모델 성능처럼 보이게 만드는 P0 리스크다.

이번 변경은 production 모드 엔진을 별도로 만들고, runtime-backed model이 없으면 시뮬레이션 폴백 대신 명시적 오류를 반환하게 한다.

## 변경 범위

- `rust-core/manpasik-engine/src/ai/mod.rs`
  - `InferenceEngine.allow_simulation` 필드를 추가했다.
  - `InferenceEngine::production(model_type)` 생성자를 추가해 운영 모드에서는 simulation fallback을 비활성화한다.
  - `with_simulation_allowed()`와 `simulation_allowed()`를 추가해 demo/test 정책을 명시적으로 다룰 수 있게 했다.
  - 현재 실제 TFLite/ONNX runtime-backed model 연결은 없으므로 production predict는 `simulation fallback disabled` 오류를 반환한다.
  - 기존 기본 생성자와 `calibration()`, `classifier()`, `anomaly_detector()`는 테스트/데모 호환을 위해 simulation을 계속 허용한다.

## 품질 게이트

- `cargo test -p manpasik-engine simulation -- --nocapture`: PASS

## 잔여 리스크

- `has_runtime_backed_model()`은 현재 `false`를 반환한다. 실제 TFLite/ONNX interpreter가 연결되면 이 함수가 interpreter readiness를 확인하도록 바뀌어야 한다.
- production entrypoint가 반드시 `InferenceEngine::production()` 또는 `with_simulation_allowed(false)`를 사용하도록 Flutter bridge/Go AI service 경계에 추가 게이트가 필요하다.
- 모델 레지스트리, 모델 서명 검증, PCCP 변경 관리, drift monitoring은 다음 AI governance 단계에서 추가해야 한다.
