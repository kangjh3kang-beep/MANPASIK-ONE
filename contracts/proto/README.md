# ManPaSik Protocol Contracts (SSOT)

Proto 파일의 단일 진실원천(Single Source of Truth)입니다.

## 원칙
- 모든 서비스는 이 디렉토리의 계약에서 코드를 생성합니다
- 비호환 변경 시 어댑터 1세대 의무 제공 (H2 후방호환)
- 변경은 `contracts/` 에서만 수행하고 코드 생성으로 전파
- Proto 메시지의 `reserved` 필드로 전방 호환성 보장

## 현재 Proto 위치 (통합 대상)
- `backend/shared/proto/manpasik.proto` — 주 서비스 정의 (33개 서비스)
- `backend/shared/proto/health.proto` — 헬스체크
- `backend/protos/measurement.proto` — 측정 패킷 상세
- `backend/protos/ai_model.proto` — AI 모델 정의

## 코드 생성 대상
- Go: `backend/shared/gen/go/v1/`
- Dart: `frontend/flutter-app/lib/generated/`
- Rust: `rust-core/manpasik-engine/src/generated/` (예정)

## 버전 정책
- 필드 번호는 절대 재사용 금지
- 삭제된 필드는 `reserved` 처리
- E12-IF 12핀 관련 필드명 사용 금지 (CSI v1.0 16핀만 유효)
