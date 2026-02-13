# B-3: Toss PG 결제 연동 계획

> **목적**: payment-service에서 실제 Toss Payments API를 호출하여 결제 승인·취소·환불을 처리한다.  
> **상태**: 착수 (2026-02-12). 검증 통과 후 본 문서 기준으로 구현 진행.

## 1. 현재 상태

- **CreatePayment**: 결제 레코드 생성(Pending). PG 호출 없음.
- **ConfirmPayment**: `payment_id`, `pg_transaction_id`, `pg_provider`를 받아 DB만 갱신(Completed). PG 사전 승인/검증 없음.
- **RefundPayment**: 환불 레코드 생성 및 결제 상태 갱신. PG 환불 API 호출 없음.

## 2. 목표

- **결제 요청**: 클라이언트가 Toss로 결제 요청 후 콜백으로 오는 `paymentKey`, `orderId`, `amount`를 우리 `payment_id`(또는 order_id)와 매핑하여 **Toss API로 결제 승인** 호출 후 성공 시 ConfirmPayment 로직 수행.
- **결제 취소/환불**: Toss Payments **취소 API** 호출 후 성공 시 기존 RefundPayment 로직 수행.
- **보안**: 시크릿 키는 환경변수(`TOSS_SECRET_KEY` 등)로만 주입. 코드/저장소에 키 미기재.

## 3. Toss Payments API 참고

- **승인**: `POST https://api.tosspayments.com/v1/payments/confirm` — body: `paymentKey`, `orderId`, `amount`.
- **취소**: `POST https://api.tosspayments.com/v1/payments/{paymentKey}/cancel` — body: `cancelReason`, `canceledAt`(optional).
- **인증**: `Authorization: Basic base64(secretKey + ":")`.
- **문서**: [Toss Payments 개발자센터](https://docs.tosspayments.com/).

## 4. 구현 방향

| 단계 | 내용 |
|------|------|
| 1 | **PG 클라이언트 인터페이스** — `Confirm(ctx, paymentKey, orderId, amountKRW)` (승인), `Cancel(ctx, paymentKey, reason)` (취소). payment-service의 service 레이어에서 주입. |
| 2 | **Toss 구현체** — HTTP 클라이언트, `TOSS_SECRET_KEY`, `TOSS_API_URL`(기본값 production). 실패 시 에러 반환, 서비스에서 DB 롤백/상태 복구. |
| 3 | **ConfirmPayment 변경** — PG 클라이언트가 설정된 경우: Toss 승인 호출 후 성공 시에만 DB Completed + 이벤트 발행. 미설정 시: 기존처럼 pg_transaction_id만 받아 완료 처리(개발/테스트). |
| 4 | **RefundPayment 변경** — PG 클라이언트가 설정된 경우: Toss 취소 API 호출 후 성공 시 환불 레코드·상태 갱신. |
| 5 | **환경변수** — `TOSS_SECRET_KEY`(필수 for 실연동), `TOSS_API_URL`(선택). 없으면 No-op PG 클라이언트 사용. |
| 6 | **테스트** — Toss 샌드박스 키로 단위 테스트 또는 통합 테스트. 금액/orderId 일치 검증. |

## 5. 파일 위치 제안

- **인터페이스**: `backend/services/payment-service/internal/service/payment.go` — `PaymentGateway` interface.
- **Toss 구현**: `backend/services/payment-service/internal/pg/toss.go` — Toss HTTP 클라이언트.
- **No-op**: `backend/services/payment-service/internal/pg/noop.go` — PG 미설정 시 사용.
- **main**: `cmd/main.go` — `TOSS_SECRET_KEY` 유무에 따라 Toss 또는 No-op 주입.

## 6. 완료 기준

- [x] PaymentGateway 인터페이스 정의 및 서비스 주입
- [x] Toss Confirm/Cancel HTTP 호출 구현
- [x] ConfirmPayment에서 Toss 승인 연동 (설정 시)
- [x] RefundPayment에서 Toss 취소 연동 (설정 시)
- [x] 시크릿 키 환경변수만 사용, 테스트/문서 반영

---

**마지막 업데이트**: 2026-02-12 (B-3 본 구현 완료)
