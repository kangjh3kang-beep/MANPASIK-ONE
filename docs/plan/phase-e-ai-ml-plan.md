# Phase E: AI/ML — 엣지추론 · 예측 · 896차원 · 연합학습

## 개요
Phase E는 AI/ML 파이프라인의 핵심 GAP을 보강하는 단계입니다.
정밀 실사 결과, 기존 구현이 예상보다 높은 수준(~55%)이었으며,
실제 GAP은 3개 영역에 집중되었습니다.

## 사전 실사 결과

### 이미 구현되어 있던 항목 (GAP 아님)
| 항목 | 위치 | 상태 |
|------|------|------|
| FedAvg 연합학습 | `data-platform-service/platform.go` L476-625 | ✅ 완전 구현 + 6 테스트 |
| K-Anonymity | `data-platform-service/platform.go` | ✅ 완전 구현 + 4 테스트 |
| 1792차원 핑거프린트 | `manpasik-engine/fingerprint/` | ✅ 88→448→896→1792 빌더 |
| InferenceEngine (TFLite 시뮬레이션) | `manpasik-engine/ai/` | ✅ 12 테스트 |
| SafetyValidator / BiasDetector | `manpasik-engine/ai/` | ✅ FM-AI-002 / FDA SaMD |
| ExplainabilityEngine | `manpasik-engine/ai/` | ✅ Permutation Importance |
| LLM 연동 (graceful degradation) | `ai-inference-service/` | ✅ 8 테스트 |
| 규칙 기반 추천 시스템 | `ai-inference-service/` | ✅ 카테고리별 규칙 |
| 위험 에스컬레이션 | `ai-inference-service/` | ✅ Critical→알림 |

### 실제 GAP 3건
1. **E-1**: Rust TrendAnalyzer — ValuePredictor가 시뮬레이션 전용, 통계적 추세 분석 부재
2. **E-2**: Go DP-FedAvg — FedAvg 존재하나 차등 개인정보보호(DP) 미적용
3. **E-3**: Go AI 모델 레지스트리 — 모델 버전 관리가 하드코딩

## 구현 내역

### E-1: Rust AI TrendAnalyzer + FFI (완료)

**파일**: `rust-core/manpasik-engine/src/ai/mod.rs`

| 구성요소 | 설명 |
|----------|------|
| `TrendDirection` enum | Improving / Worsening / Stable / Insufficient |
| `LinearRegressionResult` | slope, intercept, R², slope_std_error, 95% CI |
| `TrendAnalysis` | 방향, 회귀, 예측, 예측구간, MA, CV, 요약 |
| `TrendAnalyzer` | OLS 선형회귀 + 바이오마커별 방향 판정 |

- 바이오마커 인식: glucose, hba1c, cholesterol 등 → "lower is better"
- 예측 구간: t분포 기반 95% PI
- 한국어 요약 자동 생성

**FFI**: `rust-core/flutter-bridge/src/lib.rs`
- `TrendAnalysisDto` + `analyze_trend()` FFI 함수

**테스트**: 8건 추가
- manpasik-engine: improving, stable, insufficient, worsening, perfect R², CI (6건)
- flutter-bridge: analyze_trend_improving, analyze_trend_insufficient (2건)

### E-2: Go DP-FedAvg 차등 개인정보보호 (완료)

**파일**: `backend/services/data-platform-service/internal/service/platform.go`

| 구성요소 | 설명 |
|----------|------|
| `DPConfig` | Epsilon, Sensitivity, ClipNorm |
| `DefaultDPConfig()` | ε=1.0, Δf=1.0, C=1.0 |
| `AggregateModelWithDP()` | L2 클리핑 → FedAvg → Laplace 노이즈 |

- L2 노름 클리핑: `‖w‖₂ > C → w * C/‖w‖₂`
- Laplace 노이즈: `Lap(Δf/ε)` (결정적 시뮬레이션, 배포 시 crypto/rand)
- 메트릭 기록: dp_epsilon, dp_clip_norm, noise_scale

**테스트**: 3건 추가
- TestFederatedAggregator_WithDP
- TestFederatedAggregator_DP_InvalidEpsilon
- TestFederatedAggregator_DP_StrongPrivacy

### E-3: Go AI 모델 레지스트리 관리 (완료)

**파일**: `backend/services/ai-inference-service/internal/service/inference.go`

| 메서드 | 설명 |
|--------|------|
| `RegisterModel()` | 새 모델 등록, 기존 동일 타입 → Deprecated |
| `UpdateModelStatus()` | 모델 상태 변경 (Active/Training/Deprecated) |
| `GetActiveModels()` | Active 상태 모델만 필터링 반환 |

**테스트**: 10건 추가
- RegisterModel: Success, DeprecatesExisting, NilModel, EmptyName, EmptyVersion
- UpdateModelStatus: Success, NotFound
- GetActiveModels: 기본, deprecated 후, 등록 후

### 버그 수정
- **Go L2 클리핑**: `ClipNorm / L2²` → `ClipNorm / √(L2²)` 수정
- **Rust 안정 테스트**: 비대칭 테스트 데이터 → 대칭 배치로 수정

## 검증 결과

| 구분 | 결과 |
|------|------|
| Rust manpasik-engine | 91 PASS |
| Rust flutter-bridge | 15 PASS |
| Go ai-inference-service | ALL PASS |
| Go data-platform-service | ALL PASS |

## 변경 파일 목록
1. `rust-core/manpasik-engine/src/ai/mod.rs` — TrendAnalyzer + 6 테스트
2. `rust-core/flutter-bridge/src/lib.rs` — TrendAnalysisDto + analyze_trend FFI + 2 테스트
3. `backend/services/data-platform-service/internal/service/platform.go` — DPConfig + AggregateModelWithDP + L2 버그 수정
4. `backend/services/data-platform-service/internal/service/platform_test.go` — 3 DP 테스트
5. `backend/services/ai-inference-service/internal/service/inference.go` — RegisterModel + UpdateModelStatus + GetActiveModels
6. `backend/services/ai-inference-service/internal/service/inference_test.go` — 10 모델 레지스트리 테스트
