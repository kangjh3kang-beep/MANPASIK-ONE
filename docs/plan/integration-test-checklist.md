# Integration Test Checklist — Sprint 1

> Agent E 산출물 | E2E 시나리오 + 통합 테스트 체크리스트 | 2026-02-14

## 1. 서비스 호출 매트릭스

| Caller → Callee | Auth | Meas | Resv | Prsc | HlthRec | Fam | Noti | Admin | Tele |
|-----------------|------|------|------|------|---------|-----|------|-------|------|
| **Auth** | — | | | | | | | L | |
| **Measurement** | V | — | | | W | | E | | |
| **Reservation** | V | | — | | | | E | | G |
| **Prescription** | V | | | — | | | E | | |
| **HealthRecord** | V | | | | — | R | | L | |
| **Family** | V | R | | | R | — | E | | |
| **Notification** | V | | | | | | — | | |
| **Admin** | V | R | R | R | R | | E | — | |
| **Telemedicine** | V | | G | W | | | E | | — |

**범례:** V=ValidateToken, W=Write, R=Read, E=Event(Kafka), G=gRPC call, L=Log

## 2. E2E 시나리오

### Scenario 1: 완전한 진료 여정 (Happy Path)

```
1. 사용자 등록 → 로그인
2. 시설 검색 (지역+전문과)
3. 의사 목록 조회
4. 의사 선택 + 시간대 조회
5. 예약 생성
6. → [알림] 예약 확인 푸시
7. 원격 진료 시작
8. 진료 완료 → 처방전 생성
9. → [알림] 처방전 발행 푸시
10. 약국 선택 + 수령 방식 설정
11. 처방전 약국 전송 + 조제 토큰 발급
12. → [알림] 처방전 전송 완료 푸시
13. 약국에서 토큰으로 조회
14. 조제 상태 전이: pending → preparing → ready
15. → [알림] 조제 완료 푸시
16. 조제 상태: ready → dispensed
17. → [알림] 수령 완료 인앱
```

**검증 포인트:**
- [ ] 각 단계의 gRPC 응답 코드가 OK
- [ ] 예약 상태 전이가 올바름 (PENDING → CONFIRMED → COMPLETED)
- [ ] 조제 상태 전이가 올바름 (pending → preparing → ready → dispensed)
- [ ] 모든 알림이 정확한 템플릿으로 발송됨
- [ ] 처방전 상태가 ACTIVE → DISPENSED로 변경됨

### Scenario 2: 데이터 공유 + FHIR Export

```
1. 사용자 A: 건강 기록 생성 (검사 결과, 활력징후)
2. 사용자 A: 데이터 공유 동의 생성 (provider_id, scope: [blood_glucose, blood_pressure])
3. 제공자: ShareWithProvider(consent_id)
4. → FHIR R4 Bundle JSON 생성 + 접근 로그 기록
5. 사용자 A: 동의 목록 조회
6. 사용자 A: 접근 로그 조회 → 제공자 접근 기록 확인
7. 사용자 A: 동의 철회
8. 제공자: ShareWithProvider 재시도 → 거부됨
```

**검증 포인트:**
- [ ] FHIR Bundle JSON이 유효한 R4 형식
- [ ] scope 내 데이터만 공유됨
- [ ] DataAccessLog에 접근 기록 존재
- [ ] 철회 후 공유 시도 시 PERMISSION_DENIED
- [ ] 만료된 동의로 공유 시도 시 PERMISSION_DENIED

### Scenario 3: 가족 데이터 공유 + 접근 제어

```
1. 사용자 A: 가족 그룹 생성 (Owner)
2. 사용자 A: 사용자 B 초대 (Guardian 역할)
3. 사용자 B: 초대 수락
4. 사용자 A: 공유 설정 (blood_glucose만, 최근 7일, 승인 불필요)
5. 사용자 B: ValidateSharingAccess(blood_glucose) → 허용
6. 사용자 B: ValidateSharingAccess(heart_rate) → 거부 (AllowedBiomarkers)
7. 사용자 B: GetSharedHealthData(group, days=7) → 사용자 A 요약
8. 사용자 A: RequireApproval = true로 변경
9. 사용자 B: ValidateSharingAccess(blood_glucose) → 거부 (승인 필요)
```

**검증 포인트:**
- [ ] Owner가 자동으로 멤버에 추가됨
- [ ] Guardian만 초대 가능
- [ ] 최대 멤버 수 초과 시 에러
- [ ] AllowedBiomarkers 필터가 정확히 동작
- [ ] MeasurementDaysLimit가 적용됨
- [ ] RequireApproval=true일 때 접근 거부
- [ ] Owner/Guardian은 공유 설정 무관하게 조회 가능

### Scenario 4: 관리자 감사 + 시스템 설정

```
1. SuperAdmin: 관리자 생성 (Moderator 역할)
2. → 감사 로그: admin_create
3. SuperAdmin: 역할 변경 (Moderator → Admin)
4. → 감사 로그: role_change (OldValue: moderator, NewValue: admin)
5. Admin: 시스템 설정 변경 (security.jwt_ttl_hours: 24 → 12)
6. → 감사 로그: config_update (OldValue: 24, NewValue: 12)
7. SuperAdmin: 감사 로그 조회 (전체)
8. SuperAdmin: 감사 로그 조회 (admin_id 필터)
9. SuperAdmin: 감사 로그 조회 (action 필터)
10. SuperAdmin: 관리자 비활성화
11. → 감사 로그: admin_deactivate
```

**검증 포인트:**
- [ ] 모든 관리자 액션에 감사 로그 생성됨
- [ ] OldValue/NewValue가 정확히 기록됨
- [ ] 감사 로그 필터링이 올바름 (admin_id, action)
- [ ] 비활성화된 관리자는 IsActive=false

### Scenario 5: 약물 상호작용 + 복약 알림

```
1. 처방전 생성 (medications: [aspirin, warfarin])
2. CheckDrugInteraction([aspirin, warfarin]) → MAJOR 상호작용 경고
3. 의사가 위험 인지 후 처방 유지
4. GetMedicationReminders(user) → 시간별 복약 목록
5. 08:00 알림: aspirin (1일 2회 BID → 08:00, 20:00)
6. 08:00 알림: warfarin (1일 1회 QD → 08:00)
```

**검증 포인트:**
- [ ] 상호작용 검사 결과가 정확한 severity
- [ ] 복약 알림 시간이 frequency에 따라 정확
- [ ] BID → 2회, TID → 3회, QD → 1회

## 3. 실패 시나리오 (Negative Tests)

### 3.1 인증/권한 실패

| # | 시나리오 | 예상 결과 |
|---|---------|----------|
| N1 | 인증 없이 API 호출 | UNAUTHENTICATED |
| N2 | 만료된 JWT로 API 호출 | UNAUTHENTICATED |
| N3 | 일반 사용자가 Admin API 호출 | PERMISSION_DENIED |
| N4 | 다른 사용자의 처방전 조회 | PERMISSION_DENIED |
| N5 | 비멤버가 가족 데이터 조회 | PERMISSION_DENIED |

### 3.2 입력 검증 실패

| # | 시나리오 | 예상 결과 |
|---|---------|----------|
| N6 | 빈 user_id로 예약 생성 | INVALID_ARGUMENT |
| N7 | 잘못된 facility_id로 검색 | NOT_FOUND |
| N8 | 만료된 조제 토큰으로 조회 | FAILED_PRECONDITION (ErrTokenExpired) |
| N9 | 잘못된 상태 전이 (pending → dispensed) | INVALID_ARGUMENT |
| N10 | 이미 취소된 예약 재취소 | ALREADY_EXISTS (ErrConflict) |

### 3.3 비즈니스 규칙 실패

| # | 시나리오 | 예상 결과 |
|---|---------|----------|
| N11 | 비활성 의사에게 예약 | FAILED_PRECONDITION |
| N12 | 다른 시설 의사 선택 | INVALID_ARGUMENT |
| N13 | 약국 미설정 처방전 전송 | INVALID_ARGUMENT |
| N14 | 배송 수령인데 주소 없음 | INVALID_ARGUMENT |
| N15 | 이미 철회된 동의 재철회 | ALREADY_EXISTS |
| N16 | Owner가 그룹 탈퇴 시도 | FAILED_PRECONDITION |
| N17 | 최대 멤버 초과 초대 | RESOURCE_EXHAUSTED |

## 4. 데이터 일관성 검증

### 4.1 이벤트 정합성

| 이벤트 | 검증 항목 |
|--------|----------|
| reservation.created | reservation DB에 레코드 존재 + 알림 발송됨 |
| prescription.sent_to_pharmacy | fulfillment_token 생성됨 + DispensaryStatus=pending |
| prescription.dispensed | PrescriptionStatus=DISPENSED + DispensedAt 설정됨 |
| consent.revoked | ConsentStatus=REVOKED + RevokedAt 설정됨 |

### 4.2 상태 일관성

- 처방전 상태와 조제 상태 동기화: `DispensaryDispensed → PrescriptionStatus.Dispensed`
- 동의 만료와 상태 동기화: `ExpiresAt < now → ConsentStatus.Expired`
- 초대 만료와 상태 동기화: `ExpiresAt < now → InvitationStatus.Expired`

## 5. 성능 기준

| 항목 | 기준 | 측정 방법 |
|------|------|----------|
| 단일 gRPC 호출 응답 | < 100ms (p95) | Load test |
| FHIR Export (100 records) | < 500ms | Benchmark |
| FHIR Import (50 entries) | < 1s | Benchmark |
| 시설 검색 (+ Haversine) | < 200ms | Load test |
| Kafka 이벤트 전달 | < 1s (end-to-end) | Trace |
| 동시 사용자 처리 | 100 concurrent | Load test |

## 6. 테스트 인프라 요구사항

### 6.1 필수 서비스

- PostgreSQL (각 서비스별 DB)
- Redis (캐시 + 세션)
- Kafka (이벤트 버스)
- 21개 gRPC 서비스

### 6.2 테스트 데이터

- 시드 사용자: 5명 (환자 3, 의사 1, 관리자 1)
- 시드 시설: 3개 (병원, 의원, 약국)
- 시드 지역: KR/seoul/gangnam, US/new-york/manhattan
- 시드 약물: 5개 (aspirin, warfarin, metformin, lisinopril, omeprazole)

### 6.3 기존 E2E 테스트 파일 현황

| 파일 | 내용 |
|------|------|
| `tests/e2e/auth_flow_test.go` | 인증 플로우 |
| `tests/e2e/device_management_test.go` | 디바이스 관리 |
| `tests/e2e/measurement_flow_test.go` | 측정 플로우 |
| `tests/e2e/medical_service_test.go` | 의료 서비스 |
| `tests/e2e/community_family_test.go` | 커뮤니티/가족 |
| `tests/e2e/admin_test.go` | 관리자 |
| `tests/e2e/helpers_test.go` | 테스트 헬퍼 |

Sprint 1 Phase 3 완료 후 추가 필요:
- `tests/e2e/prescription_pharmacy_test.go` — 처방/약국 전체 플로우
- `tests/e2e/fhir_data_sharing_test.go` — FHIR Export/Import + 동의 플로우
- `tests/e2e/family_sharing_test.go` — 가족 데이터 공유 세부 시나리오

## 7. 체크리스트 요약

### Sprint 1 완료 기준

- [ ] 21/21 서비스 빌드 PASS
- [ ] 21/21 서비스 유닛 테스트 PASS
- [ ] Proto 확장 병합 완료 (protoc 에러 없음)
- [ ] Phase 3 핸들러 구현 완료
- [ ] E2E Scenario 1 (진료 여정) PASS
- [ ] E2E Scenario 2 (데이터 공유) PASS
- [ ] E2E Scenario 3 (가족 공유) PASS
- [ ] E2E Scenario 4 (관리자 감사) PASS
- [ ] E2E Scenario 5 (약물 상호작용) PASS
- [ ] 모든 Negative Test PASS
- [ ] 성능 기준 충족
