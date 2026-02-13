# AI/ML 모델 설계 전략 (만파식)

> **문서 ID**: MPS-AI-MDS-001
> **작성일**: 2026-02-09
> **작성자**: Claude (AI/ML Architecture)
> **참조**: FDA AI/ML SaMD Action Plan, IEC 62304 Class B, IMDRF SaMD N41
> **상태**: 전략 개요 (상세 설계는 Phase 2에서 진행)

---

## 1. AI/ML 시스템 개요

### 1.1 만파식 AI 아키텍처

```
┌─────────────────────────────────────────────────────────────┐
│                    AI/ML 3계층 아키텍처                        │
│                                                               │
│  [엣지 AI] ─────── [클라우드 AI] ─────── [연합학습]            │
│  (TFLite)           (서버 추론)           (프라이버시 보존)     │
│  리더기+앱에서       고급 분석,             모델 업데이트,       │
│  100% 오프라인       대규모 모델           데이터 미이동        │
│                                                               │
│  Phase 1 (MVP)      Phase 2              Phase 3              │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 모델 유형 (5종)

| # | 모델 | 입력 | 출력 | 용도 | Phase |
|---|------|------|------|------|-------|
| M1 | **Calibration** | 88차원 원시 데이터 | 88차원 보정 데이터 | 채널별 보정값 산출 | 1 (MVP) |
| M2 | **FingerprintClassifier** | 896차원 핑거프린트 | 29 클래스 확률 | 카트리지/시료 타입 분류 | 1 (MVP) |
| M3 | **AnomalyDetection** | 88차원 원시 데이터 | 1 스코어 (0~1) | 이상 측정 탐지 | 1 (MVP) |
| M4 | **ValuePredictor** | 88차원 보정 데이터 | 1 예측값 | 바이오마커 수치 예측 | 1 (MVP) |
| M5 | **QualityAssessment** | 88차원 원시 데이터 | 3 클래스 (좋음/보통/나쁨) | 측정 품질 판정 | 1 (MVP) |

---

## 2. 엣지 AI 전략 (Phase 1)

### 2.1 모델 아키텍처 설계 원칙

```
제약 조건:
  - 모바일 디바이스에서 실행 (ARM Cortex-A, Apple A-series)
  - TFLite 런타임 사용 (Rust tflitec 크레이트)
  - 추론 시간: < 100ms (사용자 체감)
  - 모델 크기: < 5MB (카트리지별 모델 다운로드 고려)
  - 100% 오프라인 동작 필수

설계 원칙:
  1. 경량 아키텍처 (MobileNet v3, EfficientNet-Lite 계열)
  2. 양자화 (INT8 Post-Training Quantization)
  3. 프루닝 (Structured Pruning, 50%+ 파라미터 감소)
  4. Knowledge Distillation (서버 모델 → 엣지 모델)
```

### 2.2 모델별 상세 설계

#### M1: Calibration Model

```
아키텍처: Dense Neural Network (DNN)
  Input: [88] float32 (원시 채널 데이터)
  Hidden: [64] → ReLU → [64] → ReLU
  Output: [88] float32 (보정된 채널 데이터)
  
학습 데이터:
  - 표준 용액(Standard Solution) 측정 데이터
  - (원시, 기대값) 쌍
  - 최소 10,000 쌍 / 카트리지 타입당

손실 함수: MSE (Mean Squared Error)
평가 지표: MAE, R², 채널별 오차 분포
모델 크기 목표: < 500KB (INT8)
```

#### M2: Fingerprint Classifier

```
아키텍처: 1D-CNN + Global Average Pooling
  Input: [896] float32 (896차원 핑거프린트)
  Conv1D: 128 filters, kernel=5 → BatchNorm → ReLU
  Conv1D: 64 filters, kernel=3 → BatchNorm → ReLU
  GlobalAvgPool
  Dense: [29] → Softmax
  
학습 데이터:
  - 각 카트리지 타입별 최소 1,000 샘플
  - 총 29,000+ 핑거프린트
  - 다양한 환경 조건 (온도, 습도, 배터리 수준)

손실 함수: Categorical Cross-Entropy + Label Smoothing
평가 지표: Top-1 Accuracy, Confusion Matrix, F1-Score (macro)
정확도 목표: ≥ 99% (Top-1), ≥ 99.9% (Top-2)
모델 크기 목표: < 2MB (INT8)
```

#### M3: Anomaly Detection

```
아키텍처: Autoencoder (Reconstruction-based)
  Encoder: [88] → [64] → [32] → [16] (latent)
  Decoder: [16] → [32] → [64] → [88]
  
이상 스코어: Reconstruction Error (MSE)
  - 정상: 낮은 재구성 오차
  - 이상: 높은 재구성 오차

학습: 정상 데이터만 사용 (비지도)
임계값: 99.5 percentile of normal reconstruction error
평가 지표: AUROC, Precision@Recall=0.95
모델 크기 목표: < 500KB (INT8)
```

#### M4: Value Predictor

```
아키텍처: DNN + Dropout (바이오마커별 개별 모델)
  Input: [88] float32 (보정된 채널 데이터)
  Hidden: [128] → ReLU → Dropout(0.2) → [64] → ReLU
  Output: [1] float32 (예측값, e.g., mg/dL)
  
카트리지 타입별 14개 모델 (바이오마커 카트리지)
  - Glucose: mg/dL (정상범위 70-100)
  - LipidPanel: Total/HDL/LDL/TG mg/dL
  - HbA1c: % (정상범위 4.0-5.6)
  - etc.

학습 데이터:
  - (보정 데이터, 기준법 측정값) 쌍
  - 기준법: 병원급 IVD 기기 (대조군)
  - 최소 500 쌍 / 바이오마커

평가 지표: MAE, MAPE, Bland-Altman 분석, 회귀 상관 (R² ≥ 0.95)
모델 크기 목표: < 300KB (INT8) / 모델
```

#### M5: Quality Assessment

```
아키텍처: DNN + Softmax
  Input: [88] float32 (원시 채널 데이터)
  Features: 통계적 특성 추출 (mean, std, min, max, skewness, kurtosis)
  Hidden: [32] → ReLU → [16] → ReLU
  Output: [3] Softmax (Good/Fair/Poor)
  
판정 기준:
  - Good: 모든 채널 SNR > 10dB, 드리프트 < 1%
  - Fair: SNR > 5dB, 드리프트 < 3%
  - Poor: 그 외 (재측정 권고)

모델 크기 목표: < 200KB (INT8)
```

---

## 3. 모델 검증 전략 (규정 준수)

### 3.1 FDA AI/ML SaMD 프레임워크 준수

| 요구사항 | 대응 방안 | 상태 |
|---------|---------|------|
| SaMD Pre-Specifications (SPS) | 모델 변경 범위 사전 정의 | ❌ 작성 필요 |
| Algorithm Change Protocol (ACP) | 모델 업데이트 시 검증 절차 | ❌ 작성 필요 |
| Good Machine Learning Practice (GMLP) | 학습/검증 분리, 편향 평가 | 🟡 설계 중 |
| Real-World Performance | 시판 후 성능 모니터링 | ❌ 계획 필요 |

### 3.2 모델 검증 프로세스

```
1. 데이터 수집 및 전처리
   ├── 데이터 품질 검사 (결측, 이상치, 분포)
   ├── 학습/검증/테스트 분할 (60/20/20 또는 K-Fold)
   └── 데이터 편향 평가 (인구통계학적 대표성)

2. 모델 학습
   ├── 하이퍼파라미터 튜닝 (Grid/Random/Bayesian)
   ├── 과적합 방지 (Early Stopping, Dropout, Regularization)
   └── 앙상블/교차검증

3. 내부 검증 (Analytical Validation)
   ├── 정확도/정밀도/재현율/F1
   ├── ROC/AUC (분류 모델)
   ├── Bland-Altman 분석 (수치 예측 모델)
   └── 보정 곡선 (Calibration)

4. 외부 검증 (Clinical Validation)
   ├── 독립적 테스트 데이터셋 (다른 기관/인구)
   ├── 기준법 비교 (Hospital-grade IVD)
   └── 하위 집단 분석 (연령, 성별, 인종)

5. 양자화 후 검증 (Deployment Validation)
   ├── FP32 vs INT8 정확도 비교 (< 1% 차이)
   ├── 추론 시간 벤치마크 (< 100ms)
   └── 엣지 디바이스 호환성
```

### 3.3 모델 편향(Bias) 평가

| 평가 항목 | 방법 | 기준 |
|----------|------|------|
| 인구통계학적 편향 | 연령/성별/인종별 성능 비교 | 하위 집단 간 성능 차이 < 5% |
| 데이터 편향 | 학습 데이터 분포 분석 | 대표 인구 비율 반영 |
| 측정 조건 편향 | 온도/습도/배터리별 성능 | 정상 운용 범위 내 성능 유지 |

---

## 4. 연합학습 전략 (Phase 3)

### 4.1 아키텍처

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Device A │     │ Device B │     │ Device C │
│ 로컬 학습 │     │ 로컬 학습 │     │ 로컬 학습 │
└─────┬────┘     └─────┬────┘     └─────┬────┘
      │ gradient        │ gradient        │ gradient
      └────────────┬────┘────────────────┘
                   ▼
            ┌─────────────┐
            │ Aggregation │  Secure Aggregation
            │ Server      │  (Paillier HE / SecAgg)
            └──────┬──────┘
                   │ global model update
                   ▼
            ┌─────────────┐
            │ Updated     │  차분 프라이버시
            │ Global Model│  (ε = 8.0, δ = 10⁻⁵)
            └─────────────┘
```

### 4.2 프라이버시 보존 기법

| 기법 | 적용 대상 | 파라미터 | 효과 |
|------|---------|---------|------|
| Secure Aggregation | gradient 전송 | Paillier 준동형암호 | 서버도 개별 gradient 볼 수 없음 |
| 차분 프라이버시 | 모델 업데이트 | ε = 8.0, δ = 10⁻⁵ | 개인 데이터 추론 불가 |
| Gradient Clipping | 학습 과정 | max_norm = 1.0 | 특이값 정보 유출 방지 |
| Client Selection | 참여 디바이스 | 무작위 30% 선택/라운드 | 특정 디바이스 추적 방지 |

---

## 5. 모델 수명주기 관리

### 5.1 Predetermined Change Control Plan (PCCP)

> FDA PCCP 가이던스에 따라 사전 승인된 변경 범위를 정의

| 변경 유형 | SPS (사전 사양) | ACP (변경 프로토콜) | 규제 영향 |
|----------|----------------|-------------------|---------|
| 학습 데이터 추가 (동일 분포) | ≤ 20% 데이터 증가 | 재학습 → 검증 → 성능 비교 | 사전 승인 범위 |
| 하이퍼파라미터 미세 조정 | 사전 정의된 범위 내 | 검증 → 성능 비교 | 사전 승인 범위 |
| 모델 아키텍처 변경 | ❌ 범위 초과 | 전체 재검증 | 신규 제출 필요 |
| 새 카트리지 타입 추가 | 분류기 클래스 추가 | 신규 클래스 검증 + 기존 성능 유지 | 변경 신고/보충 |
| 연합학습 글로벌 모델 업데이트 | 사전 승인된 집계 프로토콜 | 글로벌 모델 검증 | 사전 승인 범위 |

### 5.2 모델 버전 관리

```
모델 ID 형식: MPS-{ModelType}-v{Major}.{Minor}.{Patch}-{Quantization}

예시:
  MPS-CALIBRATION-v1.0.0-INT8
  MPS-CLASSIFIER-v1.2.1-FP16
  MPS-ANOMALY-v2.0.0-INT8

버전 규칙:
  Major: 아키텍처 변경 (규제 재제출 트리거)
  Minor: 학습 데이터 업데이트 (PCCP 범위 내)
  Patch: 양자화/최적화 (성능 동등 검증)
```

### 5.3 모델 배포 파이프라인

```
모델 학습 (서버)
  → 내부 검증 (자동)
  → 외부 검증 (수동 승인)
  → 양자화 (INT8/FP16)
  → 엣지 검증 (자동)
  → 서명 (코드 사이닝)
  → 앱 번들 또는 OTA (보안 채널)
  → 디바이스 배포
  → 배포 후 모니터링 (30일)
```

---

## 6. 다음 단계

### Phase 1 (MVP)
1. M1~M5 모델 아키텍처 확정 및 학습 파이프라인 구축
2. 시뮬레이션 데이터 기반 초기 모델 학습
3. TFLite 변환 + INT8 양자화 파이프라인
4. Rust `tflitec` 연동 (ai 모듈 시뮬레이션 → 실제 모드 전환)
5. 모델 검증 프로토콜 수립

### Phase 2 (Core)
6. 실제 측정 데이터 수집 시작 (임상 전)
7. 모델 재학습 + 성능 최적화
8. 서버 AI 모델 개발 (고급 분석)
9. FDA SPS/ACP 문서 작성

### Phase 3 (Advanced)
10. 연합학습 프레임워크 구축
11. Secure Aggregation + 차분 프라이버시 구현
12. 다기관 연합학습 파일럿

---

**Document Version**: 1.0.0
**Next Review**: 2026-02-23
**Approval Required**: AI/ML Lead, Regulatory Affairs
