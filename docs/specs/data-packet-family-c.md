# 만파식 표준 데이터 패킷 (패밀리C 준수)

**문서번호**: MPK-SPEC-DATA-v1.0  
**기준**: 원본 기획안 IX. 데이터 아키텍처, 패밀리C 공통증  
**목적**: 측정 데이터 패킷의 단일 표준 정의 및 gRPC Proto 매핑

---

## 1. 패킷 구조 개요

모든 측정 데이터는 아래 3부로 구성된 표준 패킷을 따른다.

| 구분 | 내용 | 무결성 |
|------|------|--------|
| **header** | 디바이스·세션·카트리지·펌웨어·타임스탬프 | — |
| **payload** | 원시 채널, 차동 보정 결과, 환경 메타, 상태 메타 | 해시체인 단위 |
| **footer** | checksum, schema_ver, transform_log | 서명·검증 |

---

## 2. Header

| 필드 | 타입 | 필수 | 설명 |
|------|------|------|------|
| device_id | string | O | 리더기 고유 식별자 (예: MPK-2026-001) |
| lot_id | string | O | 제조 로트 ID (예: LOT-20260208) |
| fw_ver | string (SemVer) | O | 펌웨어 버전 (예: 3.2.0) |
| cartridge_id | string | O | 카트리지 고유 ID (예: CART-GLU-001) |
| session_id | string (UUID v4) | O | 측정 세션 식별자 |
| timestamp | ISO 8601 | O | 패킷 생성 시각 (UTC) |

**Proto 매핑**: `StartSessionRequest`(device_id, cartridge_id, user_id) 및 세션·디바이스 메타데이터. 개별 패킷 header는 `StreamMeasurement` 호출 시 `session_id` + 메시지별 `timestamp`로 전달.

---

## 3. Payload

### 3.1 raw_channel

- **타입**: `double[]` (또는 채널 수에 따른 고정 길이)
- **의미**: 검출 전극 원시 신호 (채널별)
- **Proto**: `MeasurementData.raw_channels` (repeated double)

### 3.2 result (차동 보정 결과)

| 필드 | 타입 | 설명 |
|------|------|------|
| primary_value | double | 주 결과값 (예: 혈당 mg/dL) |
| unit | string | 단위 (mg/dL, ppm 등) |
| confidence | double | 신뢰도 0.0~1.0 |
| differential_correction | object | 차동 보정 상세 |

**differential_correction**:

| 필드 | 타입 | 설명 |
|------|------|------|
| s_det | double | 검출 전극 신호 |
| s_ref | double | 기준 전극 신호 |
| alpha | double | 보정 계수 (기본 0.95) |
| s_corrected | double | S_det - α × S_ref |

**Proto**: `DifferentialCorrection`, `MeasurementResult`(primary_value, unit, confidence), `MeasurementData.differential`

### 3.3 env_meta (환경 메타)

| 필드 | 타입 | 단위 | 설명 |
|------|------|------|------|
| temp_C | float | °C | 주변 온도 |
| humidity_pct | float | % | 습도 |
| pressure_kPa | float | kPa | 기압 |

**Proto**: `EnvironmentMeta` (temp_c, humidity_pct, pressure_kpa)

### 3.4 state_meta (상태 메타)

| 필드 | 타입 | 설명 |
|------|------|------|
| battery_pct | u8 | 배터리 잔량 (%) |
| signal_quality | enum | high / medium / low |
| self_diagnostic | string | pass / fail / warning |
| uncertainty | float | 측정 불확도 (예: 0.03) |

**Proto**: 현재 `MeasurementData`에 직접 필드 없음. 확장 시 `StateMeta` 메시지 추가 권장. (추적성: 원본 payload.state_meta)

---

## 4. Footer

| 필드 | 타입 | 설명 |
|------|------|------|
| checksum | string | SHA-256 해시 (hex). payload 또는 header+payload 대상. |
| schema_ver | string | 패킷 스키마 버전 (예: 2.0) |
| transform_log | array | 변환 단계별 해시 (무결성 체인) |

**transform_log** 항목 예:

| step | 설명 |
|------|------|
| raw_to_filtered | 원시 → 필터링 |
| filtered_to_corrected | 필터링 → 차동 보정 |
| corrected_to_result | 보정 → 결과값·핑거프린트 |

**Proto**: gRPC 스트림에서는 전송 구간에서 TLS로 무결성 보장. 저장 시(DB·파일) footer는 별도 컬럼 또는 해시체인 테이블로 유지 권장. (IEC 62304, 21 CFR Part 11 감사 추적)

---

## 5. JSON 예시 (로컬·파일 저장용)

```json
{
  "header": {
    "device_id": "MPK-2026-001",
    "lot_id": "LOT-20260208",
    "fw_ver": "3.2.0",
    "cartridge_id": "CART-GLU-001",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2026-02-08T09:00:00Z"
  },
  "payload": {
    "raw_channel": [1.234, 0.012, null],
    "result": {
      "primary_value": 105.0,
      "unit": "mg/dL",
      "confidence": 0.96,
      "differential_correction": {
        "s_det": 1.234,
        "s_ref": 0.012,
        "alpha": 0.95,
        "s_corrected": 1.2226
      }
    },
    "env_meta": {
      "temp_C": 23.5,
      "humidity_pct": 45,
      "pressure_kPa": 101.3
    },
    "state_meta": {
      "battery_pct": 85,
      "signal_quality": "high",
      "self_diagnostic": "pass",
      "uncertainty": 0.03
    }
  },
  "footer": {
    "checksum": "SHA256:a1b2c3...",
    "schema_ver": "2.0",
    "transform_log": [
      {"step": "raw_to_filtered", "hash": "..."},
      {"step": "filtered_to_corrected", "hash": "..."},
      {"step": "corrected_to_result", "hash": "..."}
    ]
  }
}
```

---

## 6. Proto와의 정합성

| 패밀리C 필드 | Proto 메시지·필드 | 비고 |
|-------------|-------------------|------|
| header.* | StartSessionRequest, MeasurementData.session_id, timestamp | device_id/cartridge_id는 세션 시작 시 전달 |
| payload.raw_channel | MeasurementData.raw_channels | |
| payload.result.differential_correction | DifferentialCorrection | |
| payload.result.primary_value, unit, confidence | MeasurementResult | |
| payload.env_meta | EnvironmentMeta | |
| payload.state_meta | (미정의) | StateMeta 메시지 추가 시 매핑 |
| footer | (전송 레이어) | 저장·감사 시 별도 스키마 권장 |

---

**참조**: 원본 기획안 IX절, `backend/shared/proto/manpasik.proto`, `docs/compliance/iso14971-risk-management-plan.md`
