# ManPaSik 이벤트 스키마 명세서 (Kafka Event Schema Specification)

**문서번호**: MPK-EVENT-SCHEMA-v1.0  
**갱신일**: 2026-02-12  
**목적**: Kafka(Redpanda) 토픽별 이벤트 메시지의 JSON 스키마를 정의하여 서비스 간 계약(Contract)을 확립  
**적용**: Go 백엔드 전 서비스 (`backend/shared/events/`)

---

## 1. 이벤트 공통 엔벨로프 (Common Envelope)

모든 이벤트는 아래 공통 구조를 따른다.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["event_id", "event_type", "version", "timestamp", "source", "payload"],
  "properties": {
    "event_id":    { "type": "string", "format": "uuid", "description": "이벤트 고유 ID (UUID v4)" },
    "event_type":  { "type": "string", "pattern": "^manpasik\\.[a-z_]+\\.[a-z_]+$", "description": "이벤트 타입 (dot notation)" },
    "version":     { "type": "string", "pattern": "^\\d+\\.\\d+$", "description": "스키마 버전 (major.minor)" },
    "timestamp":   { "type": "string", "format": "date-time", "description": "이벤트 발생 시각 (ISO 8601, UTC)" },
    "source":      { "type": "string", "description": "발행 서비스명 (예: measurement-service)" },
    "correlation_id": { "type": "string", "format": "uuid", "description": "요청 추적 ID (선택, X-Request-ID)" },
    "user_id":     { "type": "string", "format": "uuid", "description": "관련 사용자 ID (선택)" },
    "payload":     { "type": "object", "description": "이벤트별 페이로드" }
  }
}
```

**Go 구조체**:
```go
type EventEnvelope struct {
    EventID       string          `json:"event_id"`
    EventType     string          `json:"event_type"`
    Version       string          `json:"version"`
    Timestamp     time.Time       `json:"timestamp"`
    Source        string          `json:"source"`
    CorrelationID string          `json:"correlation_id,omitempty"`
    UserID        string          `json:"user_id,omitempty"`
    Payload       json.RawMessage `json:"payload"`
}
```

---

## 2. 토픽 설계 총괄

| 토픽명 | 파티션 | 보존 기간 | Producer | Consumer(s) | Phase |
|--------|--------|---------|----------|-------------|-------|
| `manpasik.measurement.completed` | 6 | 30일 | measurement-service | ai-inference, coaching, notification, analytics | 1 |
| `manpasik.measurement.session.started` | 3 | 7일 | measurement-service | device-service (상태 갱신) | 1 |
| `manpasik.measurement.session.ended` | 3 | 7일 | measurement-service | coaching, health-record | 1 |
| `manpasik.payment.completed` | 3 | 90일 | payment-service | subscription, shop, notification | 2 |
| `manpasik.payment.failed` | 3 | 90일 | payment-service | notification, admin | 2 |
| `manpasik.subscription.changed` | 3 | 30일 | subscription-service | cartridge, notification, user | 2 |
| `manpasik.cartridge.verified` | 3 | 30일 | cartridge-service | measurement, notification | 2 |
| `manpasik.cartridge.depleted` | 3 | 30일 | cartridge-service | notification, shop (추천) | 2 |
| `manpasik.notification.send` | 6 | 3일 | 여러 서비스 | notification-service | 2 |
| `manpasik.user.registered` | 3 | 30일 | auth-service | user, notification, coaching | 1 |
| `manpasik.user.profile.updated` | 3 | 7일 | user-service | coaching, health-record | 1 |
| `manpasik.device.registered` | 3 | 30일 | device-service | notification, admin | 1 |
| `manpasik.device.status.changed` | 6 | 7일 | device-service | notification (오프라인 경고) | 1 |
| `manpasik.ai.risk.detected` | 3 | 30일 | ai-inference-service | notification, coaching, emergency | 2 |
| `manpasik.reservation.created` | 3 | 7일 | reservation-service | notification, telemedicine | 3 |
| `manpasik.prescription.created` | 3 | 30일 | prescription-service | notification, health-record | 3 |
| `manpasik.community.post.created` | 3 | 7일 | community-service | notification, translation | 3 |
| `manpasik.dlq` | 1 | 30일 | 전체 (실패 시) | admin (수동 재처리) | 1 |

---

## 3. 이벤트별 페이로드 스키마

### 3.1 측정 도메인

#### `manpasik.measurement.completed` (v1.0)
```json
{
  "session_id": "uuid",
  "device_id": "uuid",
  "cartridge_type": "0x01",
  "cartridge_category": "health_biomarker",
  "primary_value": 95.5,
  "unit": "mg/dL",
  "confidence": 0.98,
  "fingerprint_dimension": 88,
  "fingerprint_vector_hash": "sha256:abc123...",
  "measurement_count": 10,
  "duration_ms": 15000,
  "environment": {
    "temp_c": 25.3,
    "humidity_pct": 45.2,
    "pressure_kpa": 101.3
  }
}
```

#### `manpasik.measurement.session.started` (v1.0)
```json
{
  "session_id": "uuid",
  "device_id": "uuid",
  "cartridge_id": "uuid",
  "cartridge_type": "0x01"
}
```

#### `manpasik.measurement.session.ended` (v1.0)
```json
{
  "session_id": "uuid",
  "total_measurements": 10,
  "duration_ms": 15000,
  "status": "completed"
}
```

### 3.2 결제/구독 도메인

#### `manpasik.payment.completed` (v1.0)
```json
{
  "payment_id": "uuid",
  "order_id": "uuid",
  "amount": 9900,
  "currency": "KRW",
  "payment_method": "card",
  "pg_provider": "toss",
  "pg_transaction_id": "toss_txn_abc123",
  "payment_type": "subscription",
  "subscription_id": "uuid"
}
```

#### `manpasik.payment.failed` (v1.0)
```json
{
  "payment_id": "uuid",
  "order_id": "uuid",
  "amount": 9900,
  "currency": "KRW",
  "error_code": "INSUFFICIENT_FUNDS",
  "error_message": "잔액 부족",
  "retry_count": 1,
  "max_retries": 3
}
```

#### `manpasik.subscription.changed` (v1.0)
```json
{
  "subscription_id": "uuid",
  "previous_tier": "free",
  "new_tier": "basic",
  "change_type": "upgrade",
  "effective_at": "2026-03-01T00:00:00Z",
  "max_devices": 3,
  "max_family_members": 5,
  "ai_coaching_enabled": true,
  "telemedicine_enabled": false
}
```

### 3.3 카트리지 도메인

#### `manpasik.cartridge.verified` (v1.0)
```json
{
  "cartridge_id": "uuid",
  "cartridge_uid": "NFC_UID_8bytes_hex",
  "cartridge_type": "0x01",
  "category": "health_biomarker",
  "lot_id": "LOT20260",
  "expiry_date": "20270301",
  "remaining_uses": 48,
  "max_uses": 50,
  "alpha_coefficient": 0.95,
  "verification_status": "valid"
}
```

#### `manpasik.cartridge.depleted` (v1.0)
```json
{
  "cartridge_id": "uuid",
  "cartridge_type": "0x01",
  "category": "health_biomarker",
  "total_uses": 50,
  "suggested_replacement": "0x01"
}
```

### 3.4 사용자/디바이스 도메인

#### `manpasik.user.registered` (v1.0)
```json
{
  "email_hash": "sha256:...",
  "display_name": "홍길동",
  "language": "ko",
  "timezone": "Asia/Seoul",
  "subscription_tier": "free",
  "registration_method": "email"
}
```

#### `manpasik.device.status.changed` (v1.0)
```json
{
  "device_id": "uuid",
  "serial_number": "MPK-2026-001",
  "previous_status": "online",
  "new_status": "offline",
  "battery_percent": 15,
  "signal_strength": -80,
  "firmware_version": "1.2.0",
  "last_seen": "2026-02-12T10:30:00Z"
}
```

### 3.5 AI/위험 도메인

#### `manpasik.ai.risk.detected` (v1.0)
```json
{
  "risk_id": "uuid",
  "measurement_session_id": "uuid",
  "risk_level": "high",
  "risk_score": 0.87,
  "risk_type": "hyperglycemia",
  "biomarker": "glucose",
  "measured_value": 280.5,
  "unit": "mg/dL",
  "reference_range": { "min": 70.0, "max": 140.0 },
  "recommended_action": "immediate_medical_attention",
  "alert_targets": ["user", "guardian", "doctor"]
}
```

### 3.6 알림 도메인

#### `manpasik.notification.send` (v1.0)
```json
{
  "notification_id": "uuid",
  "target_user_id": "uuid",
  "channels": ["push", "email"],
  "priority": "high",
  "category": "measurement_alert",
  "title": "혈당 수치 이상 감지",
  "body": "최근 측정 결과 혈당 수치가 280.5 mg/dL로 정상 범위(70-140)를 초과했습니다.",
  "data": {
    "action": "open_measurement_result",
    "session_id": "uuid"
  },
  "locale": "ko"
}
```

### 3.7 의료 도메인 (Phase 3)

#### `manpasik.reservation.created` (v1.0)
```json
{
  "reservation_id": "uuid",
  "facility_id": "uuid",
  "facility_type": "hospital",
  "doctor_id": "uuid",
  "appointment_time": "2026-03-15T14:00:00+09:00",
  "service_type": "telemedicine",
  "reason": "혈당 상담"
}
```

#### `manpasik.prescription.created` (v1.0)
```json
{
  "prescription_id": "uuid",
  "doctor_id": "uuid",
  "medications": [
    {
      "name": "메트포르민",
      "dosage": "500mg",
      "frequency": "2회/일",
      "duration_days": 30
    }
  ],
  "diagnosis_code": "E11.9"
}
```

---

## 4. 스키마 버전 관리 규칙

1. **하위 호환 변경** (minor 버전 증가): 선택 필드 추가, 설명 변경
   - 예: v1.0 → v1.1 (필드 추가, 기존 Consumer 영향 없음)
2. **비호환 변경** (major 버전 증가): 필수 필드 추가/제거/변경
   - 예: v1.x → v2.0 (새 토픽 생성, 구 토픽 병행 운영 후 폐기)
3. **Consumer는 알 수 없는 필드를 무시**해야 함 (Forward Compatibility)
4. **Producer는 필수 필드를 반드시 포함**해야 함 (Strict Validation)

---

## 5. Dead Letter Queue (DLQ) 처리

### DLQ 메시지 구조
```json
{
  "original_topic": "manpasik.measurement.completed",
  "original_event": { "...원본 이벤트 전체..." },
  "error": {
    "code": "DESERIALIZATION_FAILED",
    "message": "필수 필드 session_id 누락",
    "consumer": "coaching-service",
    "retry_count": 3,
    "first_failed_at": "2026-02-12T10:00:00Z",
    "last_failed_at": "2026-02-12T10:15:00Z"
  }
}
```

### DLQ 처리 정책
- **자동 재시도**: 3회, 지수 백오프 (1초 → 2초 → 4초)
- **3회 실패 시**: DLQ 토픽으로 이동
- **수동 재처리**: admin-service 대시보드에서 DLQ 메시지 확인 및 재투입
- **알림**: DLQ 메시지 발생 시 admin에게 알림 발송

---

## 6. Go 구현 가이드

### Producer 예시
```go
// measurement-service에서 측정 완료 이벤트 발행
event := events.NewEvent("manpasik.measurement.completed", "measurement-service", userID, payload)
if err := eventBus.Publish(ctx, "manpasik.measurement.completed", event); err != nil {
    log.Error("이벤트 발행 실패", zap.Error(err))
    // DLQ 또는 로컬 재시도 큐에 저장
}
```

### Consumer 예시
```go
// coaching-service에서 측정 완료 이벤트 소비
eventBus.Subscribe("manpasik.measurement.completed", func(ctx context.Context, event events.EventEnvelope) error {
    var payload MeasurementCompletedPayload
    if err := json.Unmarshal(event.Payload, &payload); err != nil {
        return fmt.Errorf("페이로드 파싱 실패: %w", err) // DLQ로 이동
    }
    return coachingService.ProcessMeasurement(ctx, event.UserID, payload)
})
```

---

**참조**: backend/shared/events/, docs/plan/PHASE11-17_MASTER_IMPLEMENTATION_PLAN.md §11-B
