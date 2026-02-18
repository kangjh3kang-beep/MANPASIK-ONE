# Agent D — 서비스 보완 항목표 (Service Gap Matrix)

> Sprint 1 Phase 1 기준 | 21개 서비스 전체 분석 | 2026-02-14

## 1. 개요

이 문서는 현재 21개 마이크로서비스의 Phase 1 구현 상태를 분석하고, Proto 확장 및 Phase 3 핸들러 구현 시 보완이 필요한 항목을 정리합니다.

**우선순위 기준:**
- **P1 (Critical)**: Sprint 1 내 반드시 완료
- **P2 (Important)**: Sprint 2에서 완료
- **P3 (Nice-to-have)**: 추후 백로그

## 2. Agent A 서비스 (의료/예약)

### reservation-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 1 | Region 계층 구조 | 서비스 로직 구현됨 | Proto 반영 + DB region 테이블 | P1 |
| 2 | Haversine 거리 계산 | 함수 구현됨 | SearchFacilities RPC에 lat/lon 파라미터 추가 | P1 |
| 3 | Doctor 도메인 | 구조체 + Repository 정의됨 | ListDoctors, SelectDoctor RPC 추가 | P1 |
| 4 | TimeSlot 조회 | GetAvailableSlots 구현됨 | GetDoctorAvailability RPC 추가 | P1 |
| 5 | 예약 이벤트 발행 | EventPublisher 연동됨 | Kafka topic 스키마 문서화 | P2 |
| 6 | 예약 알림 연동 | 미구현 | notification-service 연동 (appointment_reminder) | P2 |
| 7 | 다국어 시설명 | Facility에 region 필드만 | 시설명 다국어 지원 (translation-service 연동) | P3 |

### telemedicine-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 8 | ConsultationStatus | Proto stub 삭제됨 | manpasik.proto에 정식 반영 | P1 |
| 9 | VideoSession 관리 | 기본 구현 | WebRTC signaling 상세 구현 | P2 |
| 10 | 진료→처방 연동 | 미구현 | consultation 완료 시 prescription-service 호출 | P1 |
| 11 | 원격진료 기록 | 미구현 | health-record-service에 진료 기록 저장 | P2 |

## 3. Agent B 서비스 (처방/약국)

### prescription-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 12 | FulfillmentType | 서비스 로직 구현됨 | Proto Enum 반영 | P1 |
| 13 | DispensaryStatus 상태 머신 | 전이 규칙 구현됨 | gRPC 핸들러 매핑 | P1 |
| 14 | FulfillmentToken | 6자 토큰 생성 구현됨 | Proto 메시지 + RPC 반영 | P1 |
| 15 | 약물 상호작용 검사 | CheckDrugInteraction 구현됨 | 실제 약물 DB 연동 (외부 API) | P2 |
| 16 | 복약 알림 | GetMedicationReminders 구현됨 | notification-service 연동 (시간별 푸시) | P1 |
| 17 | 처방전 만료 관리 | ExpiresAt 필드 있음 | 크론잡/스케줄러로 만료 처리 | P2 |
| 18 | 약국 검색 | 미구현 | reservation-service의 Facility(Pharmacy) 활용 | P2 |
| 19 | 배송 추적 | FulfillmentType만 정의 | 배송 상태 추적 (외부 배송 API 연동) | P3 |

## 4. Agent C 서비스 (데이터 공유/FHIR)

### health-record-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 20 | FHIR R4 Export | ExportToFHIR 구현됨 | LOINC 코드 매핑 정확도 개선 | P1 |
| 21 | FHIR R4 Import | ImportFromFHIR 구현됨 | 중복 체크 + FHIRResourceID 기반 upsert | P2 |
| 22 | DataSharingConsent | CRUD 구현됨 | Proto 반영 + 핸들러 | P1 |
| 23 | DataAccessLog | 접근 로그 기록됨 | 5년 보관 정책 + 아카이빙 전략 | P2 |
| 24 | HealthSummary | 기간별 요약 구현됨 | AI 코칭 서비스 연동 (insight 생성) | P3 |
| 25 | 동의 만료 관리 | ExpiresAt 필드 있음 | 만료 시 자동 상태 전환 스케줄러 | P2 |
| 26 | scope 세분화 | string 배열로 구현 | scope enum 정의 + 표준화 | P2 |

## 5. Agent D 서비스 (기반 서비스)

### admin-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 27 | Region 기반 관리자 | ListByRegion 구현됨 | Proto 반영 + 핸들러 | P1 |
| 28 | AuditLogDetail (확장) | OldValue/NewValue 추적 구현됨 | Proto 반영 + 조회 RPC | P1 |
| 29 | SystemStats | 인메모리 더미 데이터 | 실제 DB 집계 쿼리 연동 | P2 |
| 30 | 역할 기반 접근 제어 | AdminRole 정의됨 | gRPC 인터셉터에서 역할 검증 | P1 |
| 31 | 대시보드 실시간 통계 | 미구현 | Redis pub/sub 또는 SSE 스트림 | P3 |
| 32 | 다국가 관리자 | CountryCode/RegionCode 있음 | 지역별 권한 범위 제한 로직 | P2 |

### notification-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 33 | PredefinedTemplates | 12개 템플릿 정의됨 | Proto 반영 + SendFromTemplate RPC | P1 |
| 34 | 채널 자동 선택 | selectBestChannel 구현됨 | Proto에 채널 enum 반영 | P1 |
| 35 | FCM 푸시 전송 | PushSender 인터페이스 | 실제 FCM 연동 구현 | P2 |
| 36 | SMTP 이메일 전송 | EmailSender 인터페이스 | 실제 SMTP/SES 연동 구현 | P2 |
| 37 | Quiet Hours | 필드 정의됨 | 시간대별 알림 지연/차단 로직 | P2 |
| 38 | 다국어 알림 | Language 필드 있음 | translation-service 연동 | P3 |
| 39 | 배치 알림 | 미구현 | 대량 사용자 일괄 발송 (프로모션 등) | P3 |

### family-service

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 40 | SharingPreferences | 바이오마커별 세분화 구현됨 | Proto 반영 + 핸들러 | P1 |
| 41 | ValidateSharingAccess | 다단계 접근 검증 구현됨 | Proto 반영 + 핸들러 | P1 |
| 42 | SharedHealthData | Guardian/Owner 우선 접근 구현됨 | 실제 measurement-service 연동 | P2 |
| 43 | 가족 초대 알림 | 미구현 | notification-service 연동 | P2 |
| 44 | 그룹 삭제 | 미구현 | Owner 전용 그룹 해산 기능 | P2 |
| 45 | 멤버 수 제한 | MaxMembers 필드 있음 | 구독 티어별 MaxMembers 차등 | P3 |

## 6. 기타 서비스 (Phase 1 미변경)

### measurement-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 46 | FHIR Observation 변환 | health-record-service의 MeasurementToFHIRObservation 연동 | P2 |
| 47 | 측정 완료 이벤트 | notification-service measurement_complete 템플릿 연동 | P2 |

### coaching-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 48 | 건강 요약 기반 코칭 | health-record-service HealthSummary 활용 | P2 |
| 49 | 가족 공유 코칭 | family-service 공유 설정 연동 | P3 |

### cartridge-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 50 | 카트리지 교체 알림 | notification-service 연동 | P2 |

### user-profile-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 51 | 건강 프로필 FHIR Patient 매핑 | FHIR Patient 리소스 생성 | P2 |

### community-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 52 | 커뮤니티 알림 | notification-service community 템플릿 연동 | P3 |

### translation-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 53 | 알림 번역 | notification-service Language 기반 번역 | P3 |

### video-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 54 | 원격진료 비디오 | telemedicine-service 비디오 세션 연동 | P2 |

### vision-service

| # | 항목 | 보완 내용 | 우선순위 |
|---|------|----------|----------|
| 55 | 테스트 파일 없음 | 기본 테스트 작성 필요 | P2 |

## 7. 횡단 관심사 (Cross-cutting Concerns)

| # | 항목 | 현재 상태 | 보완 내용 | 우선순위 |
|---|------|----------|----------|----------|
| 56 | gRPC 인터셉터 인증 | 미구현 | JWT 토큰 검증 인터셉터 | P1 |
| 57 | Rate Limiting | 미구현 | API별 요청 제한 | P2 |
| 58 | 분산 트레이싱 | 미구현 | OpenTelemetry + Jaeger | P2 |
| 59 | 헬스체크 표준화 | 서비스별 상이 | gRPC Health Checking Protocol | P1 |
| 60 | 에러 코드 표준화 | apperrors 패키지 사용 | gRPC Status Code 매핑 일관성 | P1 |
| 61 | HIPAA/PIPA 규정 준수 | DataSharingConsent 기본 구현 | 감사 로그 전 서비스 통합 | P1 |
| 62 | Proto 버전 관리 | 단일 proto 파일 | 서비스별 proto 분리 또는 버전 태깅 | P3 |

## 8. 우선순위 요약

| 우선순위 | 항목 수 | 주요 내용 |
|----------|---------|----------|
| **P1 (Critical)** | 20개 | Proto 반영, 핸들러 구현, 인증 인터셉터, 에러 표준화 |
| **P2 (Important)** | 27개 | 실제 외부 서비스 연동, 스케줄러, 서비스간 호출 |
| **P3 (Nice-to-have)** | 15개 | 다국어, 배치 처리, 실시간 대시보드 |
| **합계** | **62개** | |

## 9. Sprint 1 잔여 작업 (P1 항목)

1. **Proto 병합** — Agent A~D 제안서 기반 manpasik.proto 업데이트
2. **protoc 재생성** — `*.pb.go` 파일 재생성
3. **Phase 3 핸들러** — 신규 RPC를 handler/grpc.go에 구현
4. **gRPC 인터셉터** — JWT 인증 미들웨어
5. **에러 코드 매핑** — apperrors → gRPC Status Code
6. **헬스체크** — 모든 서비스 gRPC Health Checking Protocol 적용
7. **HIPAA/PIPA 감사** — 전 서비스 감사 로그 통합
