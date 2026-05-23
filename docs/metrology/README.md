# ManPaSik 계량학 사양 (Metrology Specification)

> **문서번호** MPS-MET-v1.0 · **표준** ISO 17511, ISO 15197, GUM (JCGM 100), UCUM
> **기존 구현 참조**: `backend/services/calibration-service/`, `manpasik-core/src/diagnostics/system_health.rs`

---

## 1. GUM 불확실성 버짓 (Measurement Uncertainty Budget)

차동 측정 공식: **Sdiff = S_det - α × S_ref**

### 1.1 Type A 불확실성 (통계적)

| 성분 | 기호 | 산출 방법 | 예상 범위 |
|------|------|----------|----------|
| 반복성 | u_rep | n회 반복 측정의 표준편차 / √n | ±0.5~2% |
| 재현성 | u_repro | 서로 다른 리더기/카트리지 조합 | ±1~3% |
| 시료 변동 | u_sample | 시료 내 불균일성 | ±0.3~1% |

### 1.2 Type B 불확실성 (체계적)

| 성분 | 기호 | 원인 | 산출 근거 |
|------|------|------|----------|
| α 보정 계수 | u_alpha | α = 0.98 ± δ (범위 0.90~1.10) | 제조사 사양 / 보정 데이터 |
| ADC 양자화 | u_adc | ADS1256 24-bit 분해능 | 1 LSB / √12 |
| 온도 보상 | u_temp | 온도 변동에 의한 센서 감도 변화 | 온도 계수 × ΔT |
| 기준물질 | u_ref | CRM 인증 불확실성 | 인증서 (ISO 17511) |
| 카트리지 배치 | u_lot | 배치 간 제조 편차 | lot별 QC 데이터 |

### 1.3 합성 표준 불확실성

```
u_c(Sdiff) = √(u_rep² + u_repro² + u_sample² + u_alpha² + u_adc² + u_temp² + u_ref² + u_lot²)
```

확장 불확실성 (95% CI): **U = k × u_c** (k = 2, 정규분포 가정)

### 1.4 구현 위치
- `contracts/packet_schema/standard_packet.json` → payload.uncertainty 필드
- `rust-core/manpasik-engine/src/differential/mod.rs` → SNR 계산
- `manpasik-core/src/diagnostics/system_health.rs` → SignalQualityAssessor

---

## 2. UCUM 단위 바인딩 테이블

> 참조: `contracts/mapping_registry/loinc_mapping.json`

| 바이오마커 | UCUM 코드 | 관례 단위 | 변환 계수 |
|-----------|----------|----------|----------|
| blood_glucose | mg/dL | mg/dL ↔ mmol/L | ÷18.0182 |
| hemoglobin_a1c | % | % (DCCT) ↔ mmol/mol (IFCC) | (% - 2.15) × 10.929 |
| cholesterol_total | mg/dL | mg/dL ↔ mmol/L | ÷38.67 |
| triglycerides | mg/dL | mg/dL ↔ mmol/L | ÷88.57 |
| blood_pressure | mm[Hg] | mmHg | — |
| heart_rate | /min | bpm | — |
| body_temperature | Cel | °C ↔ °F | °F = °C × 9/5 + 32 |
| oxygen_saturation | % | % | — |
| creatinine | mg/dL | mg/dL ↔ μmol/L | ×88.4 |
| uric_acid | mg/dL | mg/dL ↔ μmol/L | ×59.48 |
| cortisol | ug/dL | μg/dL ↔ nmol/L | ×27.59 |

**주의**: 단위 변환 오류는 임상 사고를 유발합니다. 모든 변환은 이 테이블 기준으로 수행하고, 검증 테스트를 포함해야 합니다.

---

## 3. Westgard QC 규칙

### 3.1 적용 규칙

| 규칙 | 설명 | 조치 |
|------|------|------|
| **1-2s** | 1개 QC 결과가 ±2s 초과 | 경고 (warning) |
| **1-3s** | 1개 QC 결과가 ±3s 초과 | 거부 (reject run) |
| **2-2s** | 연속 2개 QC가 같은 방향 ±2s 초과 | 거부 (systematic error) |
| **R-4s** | 연속 2개 QC 범위가 4s 초과 | 거부 (random error) |
| **4-1s** | 연속 4개 QC가 같은 방향 ±1s 초과 | 경고 (trend) |
| **10x** | 연속 10개 QC가 같은 쪽 | 거부 (bias) |

### 3.2 QC 물질 추적

| 항목 | 사양 |
|------|------|
| QC 물질 수준 | Level 1 (정상), Level 2 (비정상) |
| QC 실행 빈도 | 매 배치 시작, 8시간마다, 보정 후 |
| QC 목표 CV% | < 5% (전기화학), < 10% (면역) |
| QC 데이터 보관 | 2년 (Levey-Jennings 차트 자동 생성) |

### 3.3 구현 위치
- `backend/services/calibration-service/` → QC 검증 로직
- `manpasik-core/src/diagnostics/` → 드리프트 탐지 (EWMA/CUSUM)

---

## 4. 보정 추적성 체인 (ISO 17511)

```
SI 단위 (NIST/BIPM)
    ↓
1차 기준물질 (CRM, JCTLM 등재)
    ↓
2차 기준물질 (제조사 보정용)
    ↓
작업 기준물질 (카트리지 lot별)
    ↓
현장 보정 (Field Calibration)
    ↓
측정 결과 (Sdiff + uncertainty)
```

### 4.1 보정 수명 주기

| 단계 | 주체 | 빈도 | 기록 |
|------|------|------|------|
| 공장 보정 | 제조 | 출하 시 | POST /calibration/factory |
| 현장 보정 | 사용자/기사 | QC 실패 시 | POST /calibration/field |
| 자동 보정 확인 | 시스템 | 매 측정 전 | GET /calibration/{deviceId}/status |
| 드리프트 감지 | AI | 실시간 | SignalQualityAssessor → anomaly.detected |

---

## 5. 데이터 품질 게이트 (Kahn DQ Framework)

| 차원 | 검증 항목 | 임계값 |
|------|----------|--------|
| **Conformance** | 필수 필드 존재, UCUM 단위 유효 | 100% |
| **Completeness** | confidence + uncertainty 필드 비어있지 않음 | 100% |
| **Plausibility** | 값이 생리적 범위 내 (혈당 20~600 mg/dL) | 99.5% |
| **Timeliness** | 측정~저장 지연 < 5초 | P95 |
