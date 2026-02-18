# 만파식(MANPASIK) 전체 시스템 구축 현황 Gap 분석

> Agent 분석 산출물 | 2026-02-14 | v1.0

## 1. 전체 구현 현황 요약

| 계층 | 완성도 | 상세 |
|------|--------|------|
| 기획/설계 문서 | 95% | 94개 문서 (11 Layer), 기능정의/NFR/규정준수/UX 완비 |
| 백엔드 (Go gRPC) | 90% | 22/27 서비스 구현, 168 RPC, 5개 서비스 디렉토리만 존재 |
| Proto 정의 | 95% | 3,228줄, 21 서비스, 41 enum, 331 message |
| 데이터베이스 | 90% | 25개 SQL 초기화, 80+ 테이블, TimescaleDB/Milvus/ES 연동 |
| 인프라 | 85% | Docker Compose 25+ 컨테이너, K8s base, Prometheus+Grafana |
| Rust 코어 엔진 | 70% | 8 모듈 5,637줄, FFI 브릿지 비활성 상태 |
| 프론트엔드 (Flutter) | 35% | 6/16 기능 부분구현, 나머지 스텁/미구현 |
| 테스트 | 60% | E2E 9파일, Security 5파일, Load test 준비됨 |

## 2. 프론트엔드 Gap 매트릭스

| 기능 | Phase | 라우트 | UI | 데이터바인딩 | 백엔드연동 |
|------|-------|--------|-----|-------------|----------|
| 온보딩 | 1 | ❌ | ❌ | ❌ | ❌ |
| 인증 | 1 | ✅ | ✅ | 부분 | Email만 |
| 홈 대시보드 | 1 | ✅ | ✅ | 부분 | 부분 |
| 측정 | 1 | ✅ | ✅ | Mock | Stub |
| 디바이스 | 1 | ✅ | ✅ | 부분 | Stub |
| 설정 | 1 | ✅ | ✅ | 테마/로케일 | 부분 |
| 데이터허브 | 2 | ✅ | 기간선택만 | ❌ | ❌ |
| 마켓 | 2 | ✅ | 스텁 | ❌ | ❌ |
| 백과사전 | 2 | ❌ | ❌ | ❌ | ❌ |
| AI코치 | 2 | ✅ | 카테고리만 | ❌ | ❌ |
| 커뮤니티 | 3 | ✅ | 탭구조 | ❌ | ❌ |
| 원격진료 | 3 | ✅ | 그리드 | ❌ | ❌ |
| 가족관리 | 3 | ✅ | 카드 | ❌ | ❌ |
| 고객지원 | 1 | ❌ | ❌ | ❌ | ❌ |

## 3. 미구현 백엔드 서비스

| 서비스 | Phase | 용도 |
|--------|-------|------|
| analytics-service | 2 | BI 대시보드 |
| emergency-service | 3 | 119 연동 |
| iot-gateway-service | 4 | IoT 허브 |
| marketplace-service | 4 | 3rd-party 마켓 |
| nlp-service | 3 | NLP/AI 고도화 |
| gateway (REST) | 1 | gRPC→REST 브릿지 |

## 4. Rust FFI Gap

- flutter_rust_bridge 코드 주석처리
- BLE/NFC Mock 데이터만 반환
- CRDT 오프라인 동기화 미테스트

## 5. 인프라 Gap

- Terraform IaC 빈 파일
- CI/CD 미구성
- TURN/STUN 서버 없음
- FCM/APNS 미연동
- Toss PG 미연동

## 6. 추가 아이디어 (20개)

### UX 개선
1. 측정 위젯 (홈화면)
2. Apple Watch/WearOS 컴패니언
3. 음성 안내 모드
4. 측정 타임라인
5. 가족 비상 SOS 버튼

### 비즈니스
6. 카트리지 정기배송
7. B2B 대시보드
8. 약국 파트너 포탈
9. 보험사 API
10. 리퍼럴 프로그램

### 기술
11. GraphQL Gateway
12. Feature Flag
13. Edge Computing
14. ETL 파이프라인
15. Chaos Engineering

### 규정/보안
16. SOC 2 Type II
17. SBOM
18. 침투 테스트
19. 키 로테이션
20. 감사 로그 불변 저장소

## 7. 우선순위 로드맵

- 즉시(1-2주): Rust FFI, 온보딩, 측정결과 AI, REST Gateway
- 단기(3-4주): 데이터허브 차트, 마켓 결제, 가족관리, 커뮤니티
- 중기(5-8주): WebRTC 원격진료, 119 연동, 음식분석, CI/CD
- 장기(9-16주): B2B, HealthKit, 카트리지 구독, Federated Learning
