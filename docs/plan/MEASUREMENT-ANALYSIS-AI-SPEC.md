# 만파식 측정·분석·AI 확장 상세 구현 기획명세

**문서번호**: MPK-MEASURE-AI-SPEC-v1.0  
**작성일**: 2026-02-12  
**목적**: 88차원~1792차원 원시 데이터 추출·분석·가공·진단 전체 파이프라인과 AI 분석 확장 기능을 관련 논문·연구 자료를 바탕으로 상세하고 전문적·체계적으로 기획하여 구현 계획을 수립한다.  
**상위 문서**: [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md), [FINAL-MASTER-IMPLEMENTATION-PLAN](FINAL-MASTER-IMPLEMENTATION-PLAN.md)

---

## 1. 관련 자료·논문·연구 분석

### 1.1 다차원 바이오임피던스 분광 (EIS) + 머신러닝

| 출처 | 핵심 내용 | 만파식 설계 반영 |
| --- | --- | --- |
| ACS Omega 2025 — "ML-Enabled Multidimensional Data via Multi-Resonance Architecture" | 다중 공진 구조에서 생성되는 고차원 임피던스 데이터를 ML로 처리하면 단일 주파수 대비 검출 정확도·감도 대폭 향상 | 88→448→896→1792 다단계 차원 확장 경로의 학술적 근거. 다중 주파수 스윕 + 다중 채널(전자코/전자혀) 융합이 곧 "다중 공진" 전략 |
| Science Advances 2025 — "ML-Enhanced EIS for Real-Time Cellular Monitoring" | EIS 원시 스펙트럼에 ML을 적용해 실시간 세포 시공간 동태 비침습 모니터링 성공 | 실시간 스트리밍 파이프라인(StreamMeasurement)에서 패킷 단위 ML 추론 → 실시간 이상 알림의 이론 근거 |
| MDPI JLPEA 2025 — "EIS Data Analysis with ANNs for Resource-Constrained Devices" | 소형 ANN이 MCU에서도 1~3% 정확도로 EIS 데이터 피팅 가능. 메모리 KB 수준 | 온디바이스(Rust TFLite) 엣지 추론 전략 확인. manpasik-engine AI 모듈에서 경량 ANN 활용 |
| IEEE TBME 2023 — "Adipose Tissue Characterization with EIS + ML" | 다주파수 임피던스 데이터를 SVM·RF·XGBoost로 분류, 지방 조직 특성화 정확도 95%+ | 88차원(기본) 데이터에서도 전통 ML 앙상블로 빠른 초기 진단 가능 확인 |
| Sensors 2022 — "Dimensionality Reduction and Prediction of Impedance Data" | PCA + LSTM으로 고차원 임피던스 시계열 예측. PCA 대비 낮은 계산 복잡도 | 1792차원 궁극 단계에서 PCA/오토인코더로 잠재 표현 추출 후 LSTM 시계열 예측 적용 |

### 1.2 EIS 기반 질병 진단·POCT

| 출처 | 핵심 내용 | 만파식 설계 반영 |
| --- | --- | --- |
| Biosensors & Bioelectronics 2024 — "Applications of EIS in Disease Diagnosis: A Review" | EIS가 전극 표면 전기화학 프로세스 탐지에 유리. 혈당·콜레스테롤·DNA·단백질 바이오마커 검출에 POCT 적합 | 카트리지 타입별 바이오마커 대상과 EIS 전극 설계 근거. 88차원(기본 건강) 카트리지 검증 |
| Nature SR 2025 — "Classifying Metal Passivity from EIS with Interpretable ML" | 입력 정규화 + PCA → k-NN/얕은 신경망으로 EIS 분류. 정규화 없으면 클러스터링·신뢰도 저조 | DSP 전처리 단계에서 **정규화 필수** 확인. `normalize_signal()` 단계를 파이프라인에 강제 포함 |
| Frontiers 2025 — "ML-Driven Calibration in Impedimetric Biosensors" | ML 기반 보정 곡선이 ECM 기반 피팅보다 예측 오차 낮음. 다중 ML 모델 비교 | 차동 보정(α) 이후 ML 보정 레이어 추가. `ModelType::Calibration` (88→88) 활용 근거 |

### 1.3 고차원 신호 처리·딥러닝

| 출처 | 핵심 내용 | 만파식 설계 반영 |
| --- | --- | --- |
| Springer Energy 2025 — "Innovative DL for EIS: Attention + GRU" | Gramian Angular Field(GAF)로 EIS→이미지 변환 후 CNN+GRU+어텐션 → 정확도 극대화 | 896/1792차원 고차원 핑거프린트를 GAF 이미지로 변환 → CNN 분류 모델 추가 옵션 |
| MDPI Energies 2025 — "CGMA-Net: Conv-Gated Multi-Attention" | Conv + Multi-Head Attention + GRU 융합. RMSE/MAE 1mAh 이내 | 시계열 핑거프린트(1792차원 = 896×2 시간 윈도우)에 CGMA-Net 아키텍처 참고 |
| GitHub 2024 — "Feature Extraction Paper Code" | 임피던스 특징 추출 자동화 코드 | 특징 엔지니어링 자동화 파이프라인에 참고. `AutoFeatureExtractor` 모듈 설계 |

### 1.4 설계 반영 원칙 요약

1. **정규화 필수**: 원시 EIS 데이터는 반드시 정규화 후 분석(Nature SR 2025 확인).
2. **경량 온디바이스 추론**: 소형 ANN/TFLite로 MCU/모바일에서 1~3% 정확도 달성 가능(MDPI JLPEA 2025).
3. **다차원 = 다중공진**: 채널 수 증가가 정확도 향상과 직결(ACS Omega 2025).
4. **시계열 확장**: 시간축(이전 시점) 결합이 예측력 강화(LSTM/GRU 논문 다수).
5. **ML 보정 > ECM 보정**: 차동 보정 후 ML 보정 레이어 추가가 최적(Frontiers 2025).
6. **GAF 이미지 변환**: 고차원 스펙트럼을 이미지화하면 CNN 활용 가능(Springer 2025).

---

## 2. 차원별 데이터 구조 및 원시 데이터 스펙

### 2.1 4단계 차원 성장 경로

```text
Phase 1 (단세포) ── 88차원 ── 기본 건강 바이오마커
     ↓
Phase 2 (다세포) ── 448차원 ── 전자코(8ch) + 전자혀(8ch) + 기본 88ch 다중 주파수 확장
     ↓
Phase 3 (유기체) ── 896차원 ── 완전 센서 융합 (전체 채널 풀 스윕)
     ↓
Phase 5 (생태계) ── 1792차원 ── 896ch × 2 시간 윈도우 (현재+이전)
```

### 2.2 원시 데이터 패킷 상세

```text
┌─────────────────────────────────────────────────────────────────┐
│ PacketPayload                                                    │
├─────────────────────────────────────────────────────────────────┤
│ raw_channels: Vec<f64>    ← 채널 수: 88 / 448 / 896 / 1792     │
│ ├─ [0..87]   기본 전극 배열 (8전극 × 11주파수)                    │
│ ├─ [88..95]  전자코 채널 (8 MOX 센서)         ← Phase 2+        │
│ ├─ [96..103] 전자혀 채널 (8 ISFET 센서)       ← Phase 2+        │
│ ├─ [104..447] 확장 주파수 스윕 (344채널)       ← Phase 2+        │
│ ├─ [448..895] 고밀도 융합 채널                  ← Phase 3+        │
│ └─ [896..1791] 이전 시점 896차원 스냅샷         ← Phase 5         │
├─────────────────────────────────────────────────────────────────┤
│ env_meta: EnvironmentMeta                                        │
│ ├─ temperature_c: f32     (온도 °C)                              │
│ ├─ humidity_pct: f32      (상대습도 %)                           │
│ └─ pressure_hpa: f32      (기압 hPa)                             │
├─────────────────────────────────────────────────────────────────┤
│ state_meta: StateMeta                                            │
│ ├─ battery_pct: u8        (배터리 %)                             │
│ ├─ signal_quality: f32    (SNR dB)                               │
│ └─ self_test_ok: bool     (자가 진단)                            │
├─────────────────────────────────────────────────────────────────┤
│ timestamp_us: u64          (마이크로초 Unix 타임스탬프)            │
│ sequence_no: u32           (패킷 시퀀스 번호)                     │
│ cartridge_uid: String      (NFC 카트리지 고유 ID)                │
│ session_id: String         (측정 세션 ID)                        │
└─────────────────────────────────────────────────────────────────┘
```

### 2.3 88차원 기본 전극 배열 상세

```text
8 전극 (Working × 4, Counter × 2, Reference × 2)
× 11 주파수 포인트 (10Hz, 100Hz, 1kHz, 5kHz, 10kHz, 50kHz, 100kHz, 500kHz, 1MHz, 5MHz, 10MHz)
= 88 임피던스 값 (Z_real + Z_imag 쌍으로 44 복소수 → 88 실수)

각 값: 복소 임피던스 Z = Z_real + j·Z_imag
단위: Ω (옴)
해상도: 24-bit ADC, 유효 분해능 ~20-bit
샘플링 레이트: 100 SPS (패킷당)
```

### 2.4 448차원 확장 (Phase 2)

```text
88 (기본 전극)
+ 8  (전자코 MOX 센서: 저항 변화율 ΔR/R₀)
+ 8  (전자혀 ISFET 센서: 전위 변화 mV)
+ 344 (확장 주파수 스윕: 8전극 × 43추가주파수)
────
= 448차원

추가 주파수: 0.1Hz~100MHz 범위 로그 스케일 43포인트
전자코: 에탄올, 암모니아, 아세톤, H₂S, NO₂, CO, CH₄, VOC 감응
전자혀: pH, Na⁺, K⁺, Ca²⁺, Cl⁻, 글루코스, 유산, 요산 감응
```

### 2.5 896차원 (Phase 3) 및 1792차원 (Phase 5)

```text
896차원 = 448차원 + 448차원(고밀도 교차 전극·듀얼 주파수 인터리빙)
- 교차 전극: 8×8 매트릭스에서 비대각 쌍 28개 × 16주파수 = 448

1792차원 = 896차원(현재 시점 t₀) + 896차원(이전 시점 t₋₁)
- 시간 윈도우: Δt = 카트리지별 측정 간격 (기본 15초~180초)
- 시간축 확장으로 변화율(Δ), 기울기(slope), 안정성(stability) 특징 내재
```

---

## 3. 전체 측정·분석 파이프라인 상세

### 3.1 엔드투엔드 파이프라인 흐름

```text
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          만파식 측정·분석 파이프라인                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ① 카트리지 인식 (NFC)                                                          │
│  ├─ ReadCartridge → cartridge_uid, cartridge_type, category                     │
│  ├─ ValidateCartridge → calibration_coefficients, required_channels, meas_secs  │
│  └─ CheckCartridgeAccess → 구독 등급 확인                                       │
│                                                                                 │
│  ② 세션 시작                                                                    │
│  ├─ StartSession(session_type, cartridge_uid, device_id, concept_id)           │
│  ├─ 채널 수 설정 (88/448/896/1792)                                             │
│  └─ BLE 연결 확인 + 자가 진단 (state_meta.self_test_ok)                         │
│                                                                                 │
│  ③ 원시 데이터 수집 (BLE 스트리밍)                                              │
│  ├─ StreamMeasurement 양방향 gRPC 스트림                                        │
│  ├─ 패킷 수신 주기: 100ms (10 SPS)                                             │
│  ├─ 패킷 무결성 검사 (sequence_no 연속성, CRC)                                  │
│  └─ 환경 메타(온도·습도·기압) 동시 수집                                          │
│                                                                                 │
│  ④ 단계 A — 전처리 (Preprocessing)                                              │
│  │  4-A1. 이상치 제거 (Z-score > 3σ → 보간)                                    │
│  │  4-A2. 노이즈 필터링 (Band-pass 1Hz~15MHz, Notch 50/60Hz)                   │
│  │  4-A3. 기저선 보정 (Baseline drift removal)                                  │
│  │  4-A4. 정규화 (Min-Max → [0,1] 또는 Z-score)  ← 논문 확인: 필수            │
│  │  4-A5. 결측값 처리 (BLE 패킷 유실 시 스플라인 보간)                           │
│  │  4-A6. 윈도우 함수 (Hamming/Hann) 적용 후 FFT                               │
│  │                                                                              │
│  ⑤ 단계 B — 차동 측정 보정 (Differential Correction)                            │
│  │  5-B1. S_corrected = S_det − α × S_ref                                     │
│  │  5-B2. α = 카트리지별 calibration_coefficients.alpha (기본 0.95)             │
│  │  5-B3. 채널별 오프셋 보정: S − offset_i                                     │
│  │  5-B4. 채널별 게인 보정: S × gain_i                                         │
│  │  5-B5. 온도 보정: S × (1 + temp_coeff × (T − T_ref))                       │
│  │  5-B6. ML 보정 레이어: Calibration 모델(88→88) 적용  ← Frontiers 2025      │
│  │                                                                              │
│  ⑥ 단계 C — 특징 추출 (Feature Extraction)                                      │
│  │  6-C1. 주파수 도메인: FFT 크기/위상, 피크 주파수, 대역 에너지                  │
│  │  6-C2. 시간 도메인: RMS, 분산, 왜도, 첨도, 교차율                             │
│  │  6-C3. 통계: 채널 간 상관계수 매트릭스, PCA 상위 k 성분                       │
│  │  6-C4. 임피던스 특화: Nyquist 곡선 특징점, Cole-Cole 파라미터                 │
│  │  6-C5. GAF 이미지 변환 (896/1792차원 시)  ← Springer 2025                   │
│  │  6-C6. 자동 특징 추출: AutoFeatureExtractor (PCA + LSTM 잠재 벡터)           │
│  │                                                                              │
│  ⑦ 단계 D — 핑거프린트 생성 (Fingerprint Build)                                 │
│  │  7-D1. FingerprintBuilder.build(measurement_type)                           │
│  │  7-D2. 차원 매핑:                                                            │
│  │       Basic(88) → 88-dim 벡터                                               │
│  │       Enhanced(448) → 88 + eNose(8) + eTongue(8) + 확장주파수(344)          │
│  │       Full(896) → 448 + 교차전극(448)                                       │
│  │       Ultimate(1792) → 896(t₀) ⊕ 896(t₋₁)                                 │
│  │  7-D3. L2 정규화 → 단위 벡터화                                              │
│  │  7-D4. Milvus 벡터DB 저장 (유사 패턴 검색용)                                 │
│  │                                                                              │
│  ⑧ 단계 E — AI 추론·진단 (Inference & Diagnosis)                                │
│  │  8-E1. 모델 선택 (카트리지 타입 → ModelType 매핑)                             │
│  │  8-E2. 온디바이스 추론 (TFLite/ONNX, manpasik-engine)                       │
│  │  8-E3. 서버 추론 (AiInferenceService, GPU 가속) — 복잡 모델                  │
│  │  8-E4. 앙상블: 온디바이스 결과 + 서버 결과 가중 평균                           │
│  │  8-E5. 신뢰도 평가 (confidence 0.0~1.0)                                     │
│  │  8-E6. 결과 해석 (완전 문장 + 등급 + 한 줄 요약)                              │
│  │                                                                              │
│  ⑨ 단계 F — 후처리 및 전달                                                      │
│  │  9-F1. 결과 구조화: MeasurementResult                                       │
│  │  9-F2. 타임라인 저장: TimescaleDB (시계열)                                    │
│  │  9-F3. 핑거프린트 저장: Milvus (벡터 유사도 검색)                              │
│  │  9-F4. 이벤트 발행: measurement.completed → Kafka                            │
│  │  9-F5. 코칭 연동: CoachingService 트리거                                     │
│  │  9-F6. 알림: 이상치 감지 시 NotificationService                              │
│  │  9-F7. AI 주치의 연동: AssistantService 세션 맥락 갱신                        │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 파이프라인 단계별 Rust 모듈 매핑

| 파이프라인 단계 | Rust 모듈 | 주요 함수/구조체 |
| --- | --- | --- |
| ① 카트리지 인식 | `nfc/mod.rs` | `read_tag()`, `CartridgeInfo`, `CartridgeFullCode` |
| ② 세션 시작 | (gRPC) | `MeasurementService.StartSession` |
| ③ 원시 수집 | `ble/mod.rs` | `scan()`, `connect()`, `PacketPayload` |
| ④ 전처리 | `dsp/mod.rs` | `filter()`, `fft()`, `normalize_signal()`, `moving_average()`, `peak_detection()` |
| ⑤ 차동 보정 | `differential/mod.rs` | `DifferentialEngine.correct()`, `CorrectionParams` |
| ⑥ 특징 추출 | `dsp/mod.rs` + 신규 `feature/mod.rs` | `signal_quality()`, `AutoFeatureExtractor`, `gaf_transform()` |
| ⑦ 핑거프린트 | `fingerprint/mod.rs` | `FingerprintBuilder.build()`, `similarity()`, `euclidean_distance()` |
| ⑧ AI 추론 | `ai/mod.rs` | `AiEngine.infer()`, `ModelType`, `InferenceResult` |
| ⑨ 후처리 | (gRPC + Kafka) | `MeasurementService.EndSession`, 이벤트 발행 |

---

## 4. AI 분석 확장 기능 상세

### 4.1 모델 아키텍처 로드맵

| 단계 | 모델 | 입력 | 출력 | 아키텍처 | 배포 위치 |
| --- | --- | --- | --- | --- | --- |
| Phase 1 | Calibration | 88-dim | 88-dim (보정값) | 경량 MLP (3층, 128-64-88) | 온디바이스 (TFLite) |
| Phase 1 | BasicClassifier | 88-dim | 5클래스 (정상/경계/주의/경고/위험) | SVM/XGBoost + 소프트 투표 | 온디바이스 |
| Phase 1 | AnomalyDetector | 88-dim | 이상 스코어 (0~1) | Isolation Forest / AutoEncoder | 온디바이스 |
| Phase 2 | EnhancedClassifier | 448-dim | 15클래스 (바이오마커 조합) | CNN-1D (Conv→Pool→FC) | 서버 + 온디바이스(경량) |
| Phase 2 | FoodAnalyzer | 이미지+텍스트 | 영양소 벡터 | Vision Transformer + MLP | 서버 (GPU) |
| Phase 3 | FullFusionClassifier | 896-dim | 30클래스 | CGMA-Net (Conv+Attention+GRU) | 서버 (GPU) |
| Phase 3 | PatternMatcher | 896-dim → Milvus | Top-K 유사 패턴 | 벡터 유사도 + 메타 필터 | 서버 (Milvus) |
| Phase 5 | UltimateDiagnostic | 1792-dim | 30+α 클래스 + 연속값 | Dual-Encoder (공간+시간) + Transformer | 서버 (GPU 클러스터) |
| Phase 5 | TemporalPredictor | 1792-dim 시계열 | 미래 바이오마커 예측 | Bi-GRU + Attention | 서버 (GPU) |
| 전 Phase | QualityAssessment | 88-dim | 3클래스 (좋음/보통/나쁨) | 경량 RF | 온디바이스 |

### 4.2 온디바이스 vs 서버 추론 전략

```text
┌──────────────────────────────────────────────────────────┐
│                      추론 전략 의사결정 트리                │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  입력 차원 ≤ 88?                                         │
│  ├─ YES → 온디바이스 추론 (TFLite, <50ms)                │
│  │        ├─ 오프라인 가능                                │
│  │        ├─ Calibration + BasicClassifier + Anomaly     │
│  │        └─ 결과 즉시 표시 + 서버 비동기 검증            │
│  │                                                       │
│  └─ NO → 네트워크 확인                                   │
│          ├─ 온라인 → 서버 추론 (GPU, <500ms)              │
│          │          ├─ 448/896/1792 모델                  │
│          │          └─ 앙상블 결과 반환                    │
│          │                                                │
│          └─ 오프라인 → 경량 온디바이스 폴백                │
│                      ├─ 448→PCA→88 축소 후 Basic 모델     │
│                      └─ "서버 연결 시 정밀 분석" 안내      │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

### 4.3 AI 확장 기능 목록

| 기능 | 설명 | 입력 | 모델 | 트리거 |
| --- | --- | --- | --- | --- |
| 실시간 이상 탐지 | 스트리밍 중 패킷 단위 이상 스코어 계산 | 88-dim 패킷 | AnomalyDetector | StreamMeasurement 패킷마다 |
| 바이오마커 정량 | 혈당·콜레스테롤·요산 등 연속값 예측 | 88/448-dim | ValuePredictor (회귀) | EndSession 후 |
| 다중 바이오마커 패널 | 여러 바이오마커 동시 분류·정량 | 448/896-dim | EnhancedClassifier | EndSession 후 |
| 핑거프린트 유사도 검색 | 과거 패턴과 비교, 변화 추세 | 896/1792-dim | Milvus ANN 검색 | EndSession 후 |
| 시계열 예측 | 다음 측정 시 바이오마커 예상값 | 1792-dim 시계열 | TemporalPredictor | 3회 이상 측정 후 |
| 식품 이미지 분석 | 사진에서 영양소 추정 | 이미지+텍스트 | FoodAnalyzer (ViT) | 사용자 촬영 시 |
| 품질 게이트 | 측정 품질 자동 평가·재측정 권고 | 88-dim | QualityAssessment | 전처리 후 즉시 |
| GAF 시각화 | 고차원 데이터를 이미지로 변환·시각 표시 | 896/1792-dim | GAF Transform (비ML) | 결과 화면에서 |
| 개인화 보정 | 사용자 이력 기반 보정 모델 미세조정 | 과거 N회 데이터 | FineTuned Calibration | 10회 측정 후 자동 |
| 교차 검증 | 온디바이스 결과와 서버 결과 비교·신뢰도 조정 | 양측 결과 | 앙상블 가중 평균 | 서버 결과 도착 시 |

### 4.4 모델 생명주기 관리

```text
1. 훈련 (Training)
   ├─ 초기: 공개 EIS 데이터셋 + 시뮬레이션 데이터
   ├─ 점진: 익명 집계 데이터 (10.9 생물형 AI 연동)
   └─ 개인화: Federated Learning (로컬 미세조정, 가중치만 업로드)

2. 검증 (Validation)
   ├─ Hold-out 테스트셋 + K-fold 교차검증
   ├─ A/B 테스트 (기존 모델 vs 신규 모델)
   └─ 임상 검증 (규제 요구 시)

3. 배포 (Deployment)
   ├─ 서버: AiInferenceService 모델 레지스트리 + 카나리 배포
   ├─ 온디바이스: OTA 모델 업데이트 (DeviceService.RequestOtaUpdate)
   └─ 버전 관리: model_id, version, cartridge_type_compatibility

4. 모니터링 (Monitoring)
   ├─ 추론 지연 시간 (p50, p95, p99)
   ├─ 신뢰도 분포 히스토그램
   ├─ 데이터 드리프트 감지 (입력 분포 변화)
   └─ 피드백 루프: 사용자 "결과 정확했나요?" → 라벨 보강

5. 폐기/교체 (Deprecation)
   ├─ 신규 모델 A/B 승인 후 점진 교체
   └─ 구 모델 호환 기간 (최소 6개월)
```

---

## 5. 카트리지-측정-AI 연계 상세

### 5.1 카트리지 타입별 측정·AI 설정 매트릭스

| 카트리지 카테고리 | 대표 타입 | 차원 | 측정 시간 | 주파수 범위 | AI 모델 | 출력 |
| --- | --- | --- | --- | --- | --- | --- |
| 기본 건강 | 혈당, 콜레스테롤, 요산 | 88 | 15초 | 10Hz~10MHz | Calibration + BasicClassifier | 수치 + 5등급 |
| 감염·면역 | CRP, PCT, IgG/IgM | 88 | 20초 | 10Hz~10MHz | Calibration + BasicClassifier | 수치 + 양/음성 |
| 전자코 | 호흡 VOC, 구취, 환경 가스 | 448 | 30초 | 0.1Hz~100MHz + MOX | EnhancedClassifier | 가스 프로필 + 이상 |
| 전자혀 | 수질, 식품 신선도 | 448 | 30초 | 0.1Hz~100MHz + ISFET | EnhancedClassifier | 이온 프로필 + 등급 |
| 호르몬 패널 | 코르티솔, 갑상선 | 448 | 60초 | 0.1Hz~100MHz | EnhancedClassifier | 다중 수치 + 패널 등급 |
| 종합 건강 | 다중 바이오마커 동시 | 896 | 90초 | 전 범위 + 교차전극 | FullFusionClassifier | 30종 프로필 |
| 궁극 비표적 | 미지 물질 탐지·신약 | 1792 | 180초 | 전 범위 + 시간축 | UltimateDiagnostic | 패턴 ID + 유사도 |

### 5.2 측정 품질 게이트

```text
패킷 수신 후 즉시 (단계 ④ 전처리 직후):
─────────────────────────────────────────
1. SNR ≥ 20dB?          → NO: "환경 노이즈 높음. 위치 변경 권장"
2. 패킷 유실률 < 5%?     → NO: "BLE 연결 불안정. 기기 가까이"
3. 온도 15~40°C?         → NO: "온도 범위 초과. 실내 이동 권장"
4. 자가진단 OK?           → NO: "기기 이상. 재시작 또는 교체"
5. 배터리 > 10%?          → NO: "배터리 부족. 충전 후 측정"
─────────────────────────────────────────
모두 통과 → "측정 품질 양호" → 다음 단계 진행
하나라도 실패 → 경고 표시 + "그래도 진행?" 옵션 (결과에 품질 주의 마크)
```

---

## 6. 데이터 저장 및 흐름

### 6.1 저장소별 역할

| 저장소 | 데이터 | 보존 기간 | 접근 패턴 |
| --- | --- | --- | --- |
| TimescaleDB | 측정 시계열 (패킷별 원시값, 보정값, 결과) | 사용자 설정 (기본 2년) | 시간 범위 쿼리, 트렌드 |
| Milvus | 핑거프린트 벡터 (88~1792-dim) | 무제한 | ANN 유사도 검색 |
| PostgreSQL | 세션 메타, 카트리지 사용 로그, 결과 요약 | 무제한 | 관계형 조인 |
| Redis | 진행 중 세션 상태, 실시간 패킷 버퍼 | TTL 24시간 | 키-값 고속 |
| MinIO | GAF 이미지, 보고서 PDF | 사용자 설정 | 오브젝트 GET |
| Kafka | measurement.completed 등 이벤트 | 7일 (리텐션) | 토픽 구독 |

### 6.2 데이터 흐름 다이어그램

```text
BLE 리더기 ──패킷──► Rust Engine (온디바이스)
                     ├─ 전처리 + 차동보정 + 특징추출
                     ├─ 핑거프린트 생성
                     ├─ 온디바이스 AI 추론 (88-dim)
                     └─ 결과 ──► Flutter UI (즉시 표시)
                               │
                               ▼ (비동기)
                     Gateway ──► MeasurementService
                                ├─ TimescaleDB 저장
                                ├─ Milvus 핑거프린트 인덱싱
                                ├─ AiInferenceService (서버 모델)
                                │  └─ 앙상블 결과 → 클라이언트 갱신
                                ├─ Kafka: measurement.completed
                                │  ├─ CoachingService → 코칭 메시지
                                │  ├─ NotificationService → 알림
                                │  └─ AssistantService → 세션 맥락
                                └─ PostgreSQL (세션 메타)
```

---

## 7. 구현 Phase별 로드맵

| Phase | 차원 | Rust 모듈 | AI 모델 | 백엔드 | 프론트엔드 |
| --- | --- | --- | --- | --- | --- |
| 1 (MVP) | 88 | dsp, differential, fingerprint(Basic), ai(Calibration/Basic/Anomaly/Quality) | TFLite 경량 4종 | MeasurementService, TimescaleDB | MeasurementScreen, 결과 카드 |
| 2 | 448 | + feature(AutoExtractor), fingerprint(Enhanced) | + EnhancedClassifier, FoodAnalyzer | + AiInferenceService GPU | + 전자코/전자혀 UI, 고급 차트 |
| 3 | 896 | + fingerprint(Full), gaf_transform | + FullFusionClassifier, PatternMatcher | + Milvus ANN 검색 | + GAF 시각화, 유사 패턴 |
| 5 | 1792 | + fingerprint(Ultimate), temporal 확장 | + UltimateDiagnostic, TemporalPredictor | + GPU 클러스터, 시계열 예측 | + 예측 대시보드, 트렌드 |

---

## 8. 참조

- COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md: Part 3 (측정 서비스), Part 6.5 (데이터허브), Part 10.2 (AI 주치의), Part 10.12 (카트리지)
- FINAL-MASTER-IMPLEMENTATION-PLAN.md: I.2 추적성, II.2 측정 플로우
- AI-ASSISTANT-MASTER-SPEC.md: AI 비서의 측정 관련 의도·도구 매핑
- Rust 소스: `rust-core/manpasik-engine/src/` (dsp, differential, fingerprint, ai, ble, nfc)
- **논문/연구**: §1 표 참조
