# AI 에이전트 자동화 – Phase 매핑

**문서번호**: MPK-PLAN-AGENT-v1.0  
**기준**: 원본 기획안 VIII. AI 에이전트 자동화 설계  
**목적**: 에이전트 블록별 담당 Phase·주요 서비스/모듈·구현 담당의 공식 기준

---

## 1. 에이전트 블록 ↔ Phase · 서비스/모듈

| 에이전트 블록 | 원본 기능 요약 | 담당 Phase | 주요 서비스/모듈 | 구현 담당 |
|---------------|----------------|------------|------------------|-----------|
| **사용자 인터페이스** | 음성 명령, 자연어 대화, 음식/운동 영상 분석, 다국어 번역 | 2~4 | Flutter UI, nlp-service, translation-service, vision-service | Flutter, Go, 외부 API |
| **측정 자동화** | 카트리지 자동 인식, 프로토콜 선택, RAFE 구성, 차동측정, 재측정 트리거, 정숙구간 동기화 | 1 | Rust: ble, nfc, differential, ai; device-service | Rust, Go |
| **건강 관리** | 개인 기준선, 위험 예측, 맞춤 코칭, 측정 리마인더, 의료기관 추천·예약 | 2~3 | coaching-service, ai-inference-service, reservation-service | Go, AI/ML |
| **시스템 관리** | OTA, AI 모델 업데이트, 보정 동기화, 재고 발주, 자가진단·복구 | 1~2 | device-service (OTA), calibration-service, inventory (Phase 4) | Go, Rust |
| **긴급 대응** | 위험 수치 알림, 긴급 연락망 순차 연락, AI 음성통화, 119 연동, 위치 공유 | 3~4 | emergency-service, notification-service, location-service | Go, Flutter, 외부 |

---

## 2. Phase별 에이전트 구현 우선순위

| Phase | 에이전트 | 구현 초점 |
|-------|----------|-----------|
| 1 | 측정 자동화, 시스템 관리(OTA) | BLE/NFC/차동측정/CRDT, 디바이스 등록·OTA |
| 2 | 건강 관리(기본), 시스템 관리(보정) | 코칭·ai-inference, 보정·구독·쇼핑 |
| 3 | 사용자 인터페이스(일부), 긴급 대응 | 화상진료, 예약, 커뮤니티, 알림, 관리자 |
| 4 | 사용자 인터페이스(음성·번역), 긴급 대응(119) | SDK·마켓, 음성·번역, 119·위치·긴급 |

---

## 3. 원본 VIII 절 구조와의 대응

- **사용자 인터페이스 에이전트**: 원본 8.1 첫 블록 → Phase 2~4, Flutter + nlp/translation/vision 서비스.
- **측정 자동화 에이전트**: 원본 8.1 두 번째 블록 → Phase 1, Rust 엔진 + device-service.
- **건강 관리 에이전트**: 원본 8.1 세 번째 블록 → Phase 2~3, coaching, ai-inference, reservation.
- **시스템 관리 에이전트**: 원본 8.1 네 번째 블록 → Phase 1~2, device (OTA), calibration.
- **긴급 대응 에이전트**: 원본 8.1 다섯 번째 블록 → Phase 3~4, emergency, notification, location.

---

**참조**: 원본 기획안 VIII, `docs/plan/msa-expansion-roadmap.md`, `docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md`
