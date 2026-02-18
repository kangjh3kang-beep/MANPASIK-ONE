# ManPaSik REST API v1.0 Reference

> Base URL: `http://gateway:8080/api/v1`

## Authentication

모든 보호된 엔드포인트는 `Authorization: Bearer <token>` 헤더가 필요합니다.

### POST /auth/register
신규 사용자 등록

```json
// Request
{ "email": "user@example.com", "password": "...", "display_name": "홍길동" }
// Response 201
{ "user_id": "uuid", "email": "...", "access_token": "...", "refresh_token": "..." }
```

### POST /auth/login
로그인 (JWT 토큰 발급)

```json
// Request
{ "email": "...", "password": "..." }
// Response 200
{ "access_token": "...", "refresh_token": "...", "expires_in": 3600 }
```

### POST /auth/refresh
토큰 갱신

### POST /auth/logout
로그아웃

### POST /auth/reset-password
비밀번호 재설정 요청

### POST /auth/social-login
소셜 로그인 (Google/Apple)

---

## Users

### GET /users/:id/profile
사용자 프로필 조회

### PUT /users/:id/profile
프로필 수정

### POST /users/:id/avatar
아바타 업로드 (multipart/form-data)

### POST /users/onboarding/profile
온보딩 프로필 설정

### POST /users/onboarding/consent
온보딩 약관 동의

### POST /users/onboarding/complete
온보딩 완료

### GET /users/onboarding/status
온보딩 상태 조회

### PUT /users/:id/emergency-settings
긴급 설정 저장

---

## Measurements

### POST /measurements/sessions
측정 세션 생성

```json
// Request
{ "device_id": "...", "user_id": "...", "cartridge_id": "...", "channels": [0.5, ...] }
// Response 201
{ "session_id": "...", "result_value": 95.5, "result_unit": "mg/dL", "ai_analysis": {...} }
```

### POST /measurements/sessions/:id/end
측정 세션 종료

### GET /measurements/history
측정 이력 조회 (`?user_id=&limit=20&offset=0`)

---

## Devices

### POST /devices
디바이스 등록

### GET /devices
디바이스 목록 (`?user_id=`)

### GET /devices/:id
디바이스 상세

### DELETE /devices/:id
디바이스 삭제

### GET /devices/:id/firmware/check
펌웨어 업데이트 확인

---

## Cartridges

### POST /cartridges/read
NFC 카트리지 읽기

### POST /cartridges/validate
카트리지 검증

### POST /cartridges/usage
카트리지 사용 기록

### GET /cartridges/types
카트리지 타입 목록

### GET /cartridges/:id/remaining
잔여 사용횟수

---

## Health Records

### POST /health-records
건강 기록 생성

### GET /health-records
건강 기록 목록 (`?user_id=&type_filter=&limit=20`)

### GET /health-records/:id
건강 기록 상세

### POST /health-records/reports
건강 리포트 생성

### POST /health-records/export/fhir
FHIR 형식 내보내기

### POST /health-records/import
외부 건강 데이터 가져오기

---

## Reservations

### POST /reservations
예약 생성

### GET /reservations
예약 목록 (`?user_id=`)

### GET /reservations/:id
예약 상세

---

## Telemedicine

### GET /telemedicine/doctors
의료진 검색 (`?specialty=`)

### POST /telemedicine/consultations
진료 상담 생성

### GET /telemedicine/consultations/:id/result
진료 결과 조회

---

## Prescriptions

### GET /prescriptions
처방전 목록

### POST /prescriptions/:id/pharmacy
약국 선택

### POST /prescriptions/:id/send
약국 전송

### GET /prescriptions/token/:token
처방전 조회 (약국용)

---

## Family

### POST /family/groups
가족 그룹 생성

### GET /family/groups
가족 그룹 목록 (`?user_id=`)

### POST /family/groups/:id/invite
가족 초대

### POST /family/invites/accept
초대 수락

### GET /family/groups/:id/members
그룹 멤버 목록

### PUT /family/groups/:id/members/:memberId
멤버 역할/모드 수정

### GET /family/dashboard
가족 대시보드

### GET /family/guardian/dashboard
보호자 대시보드 (`?period=7d`)

### GET /family/members/:id/health-data
멤버 건강 데이터 (`?period=7d`)

### POST /family/alerts
긴급 알림 발송

### GET /family/alerts
알림 목록 (`?type=emergency`)

### GET /family/groups/:id/report
가족 건강 리포트

---

## Community

### GET /posts
게시글 목록 (`?category=&query=&limit=20`)

### POST /posts
게시글 작성

### GET /posts/:id
게시글 상세

### POST /posts/:id/like
좋아요

### POST /posts/:id/bookmark
북마크

### POST /posts/:id/comments
댓글 작성

### GET /community/challenges
챌린지 목록

### POST /community/challenges/:id/join
챌린지 참여

### GET /community/qna
Q&A 목록

---

## Market

### GET /products
상품 목록 (`?category=&limit=20`)

### GET /products/:id
상품 상세

### GET /products/:id/reviews
상품 리뷰

### POST /products/:id/reviews
리뷰 작성

### POST /cart
장바구니 추가

### GET /cart/:userId
장바구니 조회

### PATCH /cart/:productId
장바구니 수량 변경

### DELETE /cart/:productId
장바구니 상품 삭제

### POST /orders
주문 생성

### GET /orders
주문 목록 (`?user_id=`)

### GET /orders/:id
주문 상세

### GET /orders/:id/tracking
배송 추적

---

## Subscriptions

### GET /subscriptions/plans
구독 플랜 목록

### GET /subscriptions/plans/compare
플랜 비교표

### GET /subscriptions/:userId
현재 구독 조회

### POST /subscriptions
구독 생성

### POST /subscriptions/upgrade
구독 업그레이드

### DELETE /subscriptions/:id
구독 취소

---

## Payments

### POST /payments
결제 생성

### POST /payments/:id/confirm
결제 확인 (PG 콜백)

### GET /payments/:id
결제 상세

---

## Notifications

### GET /notifications
알림 목록 (`?user_id=&unread_only=true`)

### GET /notifications/unread-count
읽지 않은 알림 수

### POST /notifications/:id/read
알림 읽음 처리

### GET /notifications/alerts/:id
알림 상세

---

## AI / Coaching

### POST /ai/analyze
측정 AI 분석

### GET /ai/health-score/:userId
건강 점수 (`?days=30`)

### POST /ai/predict-trend
트렌드 예측

### GET /ai/models
AI 모델 목록

### POST /ai/food-analyze
음식 이미지 분석 (multipart)

### POST /ai/exercise-analyze
운동 영상 분석 (multipart)

### POST /coaching/goals
건강 목표 설정

### GET /coaching/goals/:userId
건강 목표 목록

### POST /coaching/generate
AI 코칭 생성

### GET /coaching/daily-report/:userId
일일 건강 리포트

### GET /coaching/recommendations/:userId
맞춤 추천

---

## Sync (Offline)

### POST /sync/measurements
오프라인 측정 일괄 동기화

### GET /sync/status
동기화 상태

### GET /sync/conflicts
충돌 목록

### POST /sync/conflicts/resolve
충돌 해결

### POST /sync/settings
오프라인 설정 동기화

---

## Admin (관리자 전용)

### GET /admin/users
사용자 관리 목록

### PUT /admin/users/:id/role
사용자 역할 변경

### POST /admin/users/bulk
일괄 사용자 작업

### GET /admin/stats
시스템 통계

### GET /admin/audit-log
감사 로그

### GET /admin/metrics
시스템 메트릭

### GET /admin/hierarchy
조직 계층 구조

### GET /admin/compliance
규제 체크리스트

### GET /admin/compliance/gdpr
GDPR 체크리스트

### GET /admin/compliance/pipa
PIPA 체크리스트

### GET /admin/compliance/hipaa
HIPAA 체크리스트

### POST /admin/compliance/reports
규제 보고서 생성

### GET /admin/system/config
시스템 설정

### PUT /admin/system/config
시스템 설정 변경

### POST /admin/cartridges
카트리지 타입 등록

---

## Compliance (사용자 자기 데이터 관리)

### GET /users/me/audit-logs
내 감사 로그

### POST /users/me/data-export
데이터 내보내기 (GDPR Right to Access)

### GET /users/me/data-export/status
내보내기 상태

### POST /users/me/data-deletion
데이터 삭제 (GDPR Right to Erasure)

### GET /users/me/consents
동의 상태

### PATCH /users/me/consents
동의 변경

### GET /users/me/consents/history
동의 이력

---

## Translations

### GET /translations
번역 조회 (`?locale=ko`)

### POST /translations/translate
텍스트 번역

---

## Video (WebRTC)

### POST /video/rooms/:roomId/join
영상 통화 입장

### POST /video/rooms/:roomId/leave
영상 통화 퇴장

---

## Calibration

### POST /calibration/factory
공장 캘리브레이션 등록

### POST /calibration/field
현장 캘리브레이션

### GET /calibration/:deviceId/status
캘리브레이션 상태

### GET /calibration/models
캘리브레이션 모델 목록

---

## Support

### POST /support/inquiries
1:1 문의 등록

---

## Error Response Format

```json
{
  "error": {
    "code": "INVALID_TOKEN",
    "message": "유효하지 않은 인증 토큰",
    "details": {}
  }
}
```

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 409 | Conflict |
| 429 | Too Many Requests |
| 500 | Internal Server Error |
