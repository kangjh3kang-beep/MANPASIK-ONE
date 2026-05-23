# ManPaSik AI 모델카드 (FDA AI/ML SaMD 가이던스 준수)

> **문서번호** MPS-AI-MC-v1.0 · **표준** FDA AI/ML 2023 · IEC 62304 Class B
> **기존 구현 참조**: `rust-core/manpasik-engine/src/ai/mod.rs` (BiasDetector, ExplainabilityEngine, SafetyValidator, TrendAnalyzer)

---

## 1. 캘리브레이션 모델 (Calibration)

| 항목 | 내용 |
|------|------|
| **목적** | 88채널 원시 센서 데이터를 보정된 88채널 출력으로 변환 |
| **입력/출력** | 88-dim float → 88-dim float |
| **모델 유형** | Linear regression + temperature compensation |
| **훈련 데이터** | 합성 데이터 (896-dim fingerprint 기반, 5-fold CV) [현재 합성, 실측 데이터 수집 필요] |
| **성능 지표** | RMSE per fold, R² > 0.95 목표 |
| **편향 평가** | BiasDetector (rust-core ai/mod.rs) — 인구통계 그룹별 정확도/민감도 격차 < 5% |
| **드리프트 모니터링** | SignalQualityAssessor (manpasik-core diagnostics) — EWMA/CUSUM 기반 |
| **PCCP 변경 범위** | 보정 계수 재학습: SPS(동일 성능 기준), 재훈련 프로토콜 필요 |
| **롤백** | 이전 모델 버전으로 자동 롤백 (A/B 파티션) |

## 2. 핑거프린트 분류기 (Fingerprint Classifier)

| 항목 | 내용 |
|------|------|
| **목적** | 1792-dim 핑거프린트를 30개 바이오마커 클래스로 분류 |
| **입력/출력** | 1792-dim float → 30-class softmax |
| **모델 유형** | RandomForest (ONNX export) |
| **훈련 데이터** | 합성 896-dim → 1792 temporal expansion [실측 데이터 수집 필요] |
| **성능 지표** | 5-fold stratified CV, accuracy per fold, confusion matrix |
| **편향 평가** | BiasDetector — 연령/성별/인종별 per-class accuracy 격차 분석 |
| **드리프트 모니터링** | AUC 5% 하락 시 알림, 월간 성능 리포트 |
| **PCCP 변경 범위** | 새 바이오마커 클래스 추가: ACP(신규 510(k) 필요 여부 평가) |

## 3. 이상탐지기 (Anomaly Detector)

| 항목 | 내용 |
|------|------|
| **목적** | Out-of-distribution 측정값 탐지, 비정상 카트리지/센서 감지 |
| **입력/출력** | 88-dim float → anomaly score (0~1) |
| **모델 유형** | Isolation Forest + 통계적 Mahalanobis 폴백 |
| **성능 지표** | Precision/Recall at threshold, F1 > 0.90 목표 |
| **편향 평가** | 드문 바이오마커 유형에 대한 false positive 비율 모니터링 |
| **드리프트 모니터링** | 이상 탐지율 변화 추적 (월간) |
| **PCCP** | 임계값 변경: SPS, 모델 교체: ACP |

## 4. 건강점수 예측기 (Health Score Predictor)

| 항목 | 내용 |
|------|------|
| **목적** | 다항목 측정 데이터로부터 종합 건강 점수(0~100) 산출 |
| **입력/출력** | 측정값 배열 + 메타데이터 → score (0~100) |
| **모델 유형** | XGBoost regressor |
| **안전 검증** | SafetyValidator (ai/mod.rs) — confidence ≥ 0.85, 기준범위 검사, 위험값 탐지 |
| **설명 가능성** | ExplainabilityEngine — 기여 채널 Top-K, 방향(긍정/부정) |
| **편향 평가** | 연령/성별별 점수 분포 편향 분석 |
| **PCCP** | 가중치 재학습: SPS, 새 입력 추가: ACP |

## 5. 트렌드 분석기 (Risk Trend Analyzer)

| 항목 | 내용 |
|------|------|
| **목적** | 시계열 측정 데이터의 추세 분석 및 미래 예측 |
| **입력/출력** | 시계열 배열 → trend direction + prediction ± CI |
| **모델 유형** | Linear regression + EWMA (TrendAnalyzer, ai/mod.rs) |
| **성능 지표** | R², RMSE, 예측 구간 커버리지 (95% CI) |
| **안전 경고** | 악화 추세 시 anomaly.detected 이벤트 발행 |
| **PCCP** | 예측 윈도우 변경: SPS, 비선형 모델 전환: ACP |

---

## 공통 사항

- **모든 모델 출력에 confidence + uncertainty 필수** (과대표현 금지, CLAUDE.md §1)
- **의료 진단 단정 금지** → 인간 주치의 핸드오프 (Human-in-the-loop)
- **편향 임계값**: 인구통계 그룹 간 accuracy gap < 5% (BiasDetector 기본값)
- **모델 버전 형식**: MPS-{ModelType}-v{Major}.{Minor}.{Patch}
- **기존 구현 위치**:
  - BiasDetector: `rust-core/manpasik-engine/src/ai/mod.rs` (lines 911-996)
  - ExplainabilityEngine: `rust-core/manpasik-engine/src/ai/mod.rs` (lines 1030+)
  - SafetyValidator: `rust-core/manpasik-engine/src/ai/mod.rs`
  - TrendAnalyzer: `rust-core/manpasik-engine/src/ai/mod.rs`
  - SignalQualityAssessor: `manpasik-core/src/diagnostics/system_health.rs`
