# 다음 단계 세부 계획 (B-3 Toss PG 연동)

> **기준일**: 2026-02-12  
> **목표**: B-3 Toss 실제 결제 승인·취소 연동 완료 및 검증

---

## 1. 범위

| 항목 | 내용 |
|------|------|
| **B-3 본 구현** | Toss Payments HTTP 클라이언트 구현, ConfirmPayment/RefundPayment 연동, main 조건부 주입 |
| **Proto 확장** | ConfirmPaymentRequest에 `payment_key`(optional) 추가 — Toss 리다이렉트 콜백에서 받은 키 전달용 |
| **비목표** | Flutter/웹 결제 UI, Toss 결제창 연동(클라이언트) — 별도 태스크 |

---

## 2. 세부 작업 순서

### Step 1: Toss HTTP 클라이언트 구현
- **파일**: `backend/services/payment-service/internal/pg/toss.go`
- **내용**:
  - 구조체: `TossClient{ baseURL, secretKey, httpClient }`, 기본 baseURL `https://api.tosspayments.com`
  - `Confirm(ctx, paymentKey, orderId string, amountKRW int32) (pgTransactionID string, err error)`: `POST /v1/payments/confirm` Body JSON `{"paymentKey": paymentKey, "orderId": orderId, "amount": amountKRW}`. 인증 `Authorization: Basic base64(secretKey + ":")`. 응답에서 트랜잭션 식별자 추출(예: `paymentKey` 또는 응답 필드)하여 반환.
  - `Cancel(ctx, paymentKey, reason string) error`: `POST /v1/payments/{paymentKey}/cancel` Body JSON `{"cancelReason": reason}`. 2xx면 nil.
  - 환경변수: 시크릿은 호출측에서 주입, 코드에 미기재.
- **의존**: 없음.

### Step 2: Proto에 payment_key 추가
- **파일**: `backend/shared/proto/manpasik.proto`
- **내용**: `ConfirmPaymentRequest`에 `string payment_key = 4;` 추가. (Toss 콜백에서 받은 결제 키, PG 승인 호출 시 사용)
- **재생성**: `make proto` 또는 프로젝트 규칙에 따른 Go/Flutter 코드 생성. 기존 필드 호환 유지.

### Step 3: ConfirmPayment 서비스 로직 변경
- **파일**: `backend/services/payment-service/internal/service/payment.go`
- **내용**:
  - 시그니처 확장: `ConfirmPayment(ctx, paymentID, pgTransactionID, pgProvider, paymentKey string)`. (paymentKey는 optional, 빈 문자열 가능)
  - 로직:  
    1) 기존처럼 payment 조회, Pending 검증.  
    2) **s.pgGateway != nil && paymentKey != ""** 이면: `txID, err := s.pgGateway.Confirm(ctx, paymentKey, payment.OrderID, payment.AmountKRW)`. err 시 결제 실패 이벤트 발행 후 에러 반환. 성공 시 pgTransactionID = txID, pgProvider = "toss"(고정 또는 파라미터 유지).  
    3) 그 외: 기존처럼 인자 pgTransactionID, pgProvider 사용.  
    4) 이후 DB 갱신·이벤트 발행은 기존과 동일.
- **호환**: paymentKey == "" 이면 기존 동작과 동일.

### Step 4: RefundPayment 서비스 로직 변경
- **파일**: `backend/services/payment-service/internal/service/payment.go`
- **내용**:
  - **s.pgGateway != nil && payment.PgTransactionID != ""** 이면: `err := s.pgGateway.Cancel(ctx, payment.PgTransactionID, reason)`. (Toss는 취소 시 paymentKey 사용. 우리가 저장한 PgTransactionID가 Toss paymentKey와 동일하다고 가정.) 실패 시 DB 변경 없이 에러 반환. 성공 시 기존 환불 레코드·상태 갱신·이벤트.
  - 그 외: 기존처럼 DB만 갱신.
- **참고**: Toss 응답에서 받은 `paymentKey`를 PgTransactionID로 저장해 두면 취소 시 동일 값으로 호출 가능.

### Step 5: Handler 및 Proto 생성 코드 반영
- **파일**: `backend/services/payment-service/internal/handler/grpc.go`
- **내용**: `ConfirmPayment` 호출 시 `req.PaymentKey` 전달. `ConfirmPayment(ctx, req.PaymentId, req.PgTransactionId, req.PgProvider, req.PaymentKey)`.
- **Proto 생성**: 수정 후 `make proto` 실행해 `ConfirmPaymentRequest`에 `payment_key` 반영. Go 패키지 재생성.

### Step 6: main 조건부 주입
- **파일**: `backend/services/payment-service/cmd/main.go`
- **내용**: `TOSS_SECRET_KEY` 환경변수 존재 시 `pg.NewTossClient(cfg.TossSecretKey, cfg.TossAPIURL)`(또는 유사 생성자), 없으면 `pg.NewNoopGateway()`. `paySvc.SetPaymentGateway(...)`.
- **설정**: `shared/config`에 `TossSecretKey`, `TossAPIURL`(optional) 필드 추가 및 `LoadFromEnv`에서 읽기.

### Step 7: 테스트·문서·빌드
- **단위 테스트**: `payment_test.go` — pgGateway가 Noop일 때 기존 ConfirmPayment/RefundPayment 동작 유지. 필요 시 Toss 클라이언트 모크으로 Toss 호출 경로 테스트.
- **문서**: `docs/plan/B3-toss-pg-integration.md` 완료 기준 체크, `CHANGELOG.md`에 B-3 본 구현 항목 추가.
- **빌드**: `go build ./...`, `go test ./...` (payment-service 포함).

---

## 3. 완료 기준 체크리스트

- [x] TossClient 구현 (Confirm/Cancel), noop 대비 인터페이스 충족
- [x] Proto ConfirmPaymentRequest.payment_key 추가, 생성 코드 반영
- [x] ConfirmPayment: paymentKey 있을 때 pgGateway.Confirm 호출 후 DB/이벤트
- [x] RefundPayment: pgGateway 있을 때 pgGateway.Cancel 호출 후 DB/이벤트
- [x] main: TOSS_SECRET_KEY 유무에 따라 Toss vs Noop
- [x] config: Toss 시크릿/URL 환경변수 로드
- [x] 기존 단위 테스트 통과, 빌드·vet 통과

---

## 4. 위험·가정

- **Toss 응답 형식**: 승인 응답에서 트랜잭션 ID로 사용할 필드(paymentKey 또는 고유 ID)는 Toss 문서 기준으로 선택. 필요 시 응답 구조 확인 후 매핑.
- **취소 API**: Toss 취소 시 사용하는 키가 승인 시 받은 paymentKey와 동일하다고 가정. 우리는 PgTransactionID에 paymentKey를 저장해 두면 됨.
- **Proto 재생성**: `make proto` 실행 시 기존 생성 코드 덮어쓰기. Flutter 등 다른 클라이언트도 재생성 필요할 수 있음.

---

**마지막 업데이트**: 2026-02-12
