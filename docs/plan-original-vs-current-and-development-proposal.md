# 만파식 기획안 원본 대비 검증 및 기획안 발전 수립안

**문서번호**: MPK-PLAN-VERIFY-v1.0-20260210  
**작성일**: 2026년 2월 10일  
**대상 원본**: ManPaSik AI Ecosystem Construction Plan v1.0 FINAL (MPK-ECO-PLAN-v1.0-20260208-FINAL)  
**목적**: 원본 기획안에 대한 현재 시스템의 부합도 검증, 보완사항 정리, 기획안 발전을 위한 수립안 제안

---

## 1. 총괄 요약

| 구분 | 요약 |
|------|------|
| **원본 부합도** | 핵심 기술·아키텍처·규정 방향은 원본에 충실. MSA·사이트맵·스토리보드·AI 에이전트·데이터 패킷 상세는 부분 반영 또는 미반영. |
| **보완 필요** | MSA 도메인 축소 현황 명시, 데이터 패킷(패밀리C) 문서화, 사이트맵/스토리보드 공식 문서화, SDK·AI 에이전트 설계 문서 추가. |
| **발전 제안** | 기획안 v1.1에서 “현행 시스템 매핑표”, “단계별 MSA 확장 로드맵”, “기획–구현 추적성 매트릭스”를 추가하여 완성도와 실행 연계를 강화. |

---

## 2. 원본 섹션별 검증 (부합 여부)

### 2.1 I. 총괄 개요 (Executive Summary)

| 항목 | 원본 내용 | 현재 반영 | 부합도 | 비고 |
|------|----------|----------|--------|------|
| 프로젝트 정의 | 차동측정 기반 범용 분석, 88→896차원, 전자코/전자혀 융합 | CONTEXT·AGENTS에 동일 정의 | ✅ 부합 | |
| 측정 차원 | 88→448→896 | Rust FingerprintBuilder, MAX_CHANNELS=896 | ✅ 부합 | |
| 리더기/카트리지 | 29종 카트리지, 동시 리더기 최대 10대 | 29종 CartridgeType, 무제한 확장으로 규칙 명시 | ✅ 부합 | 원본 "10대"는 구독별 제한으로 해석 가능 |
| 오프라인 100% | 완전 구동 | CRDT 동기화, 오프라인 모드 설계 | ✅ 부합 | |
| SaaS 구독 티어 | Basic Safety / Bio-Optimization / Clinical Guard (3단계) | Free + Basic + Pro + Clinical (4단계) | 🟡 부분 | 무료 티어 추가는 사업 판단 사항, 원본과 충돌 없음 |
| 데이터 보안 | TPM + 해시체인 + AES-256 | AES-256-GCM, SHA-256 해시체인 구현 | ✅ 부합 | TPM은 하드웨어 스펙, SW는 해시체인·암호화 반영 |
| 기술 통합 기반 | 제안서_펄스팩, 896Dim 제조설계서, cursor 어플 기술문서 등 | AGENTS·규정 문서에서 프로토콜·카트리지 반영 | 🟡 부분 | 원본 문서 자체는 repo 외부 참조 |

**결론**: 총괄 개요의 핵심 수치·개념은 현재 시스템에 충실히 반영됨.

---

### 2.2 II. 생태계 철학 및 비전

| 항목 | 원본 내용 | 현재 반영 | 부합도 |
|------|----------|----------|--------|
| 홍익인간 이념 | "세상의 모든 파동을 분석하여 인간을 이롭게" | AGENTS.md 프로젝트 아이덴티티에 명시 | ✅ 부합 |
| 비전 계층도 | 건강/환경/식품/산업/안전 5대 영역 | 카트리지 29종이 해당 영역 매핑 | ✅ 부합 |
| 유기적 성장 모델 | 88→448→896→1792, 생물체 비유 | 88→448→896→1792 4단계 성장 경로 구현 (1792=Phase 5 궁극 확장) | ✅ 부합 |

**결론**: 철학·비전은 문서와 카트리지 체계에 부합.

---

### 2.3 III. 기술 스택 선정

| 계층 | 원본 | 현재 | 부합도 | 비고 |
|------|------|------|--------|------|
| 모바일 | Flutter 3.x | Flutter 3.x, Riverpod, go_router | ✅ 부합 | |
| 웹 | Next.js 14+ | Next.js 14+ | ✅ 부합 | |
| 로컬 DB | SQLite + Sled | 스펙에는 SQLite + Hive, Rust는 Sled 언급 | 🟡 부분 | Flutter 스펙은 Hive, 원본 Sled는 Rust 측 |
| 코어 엔진 | Rust no_std, TFLite, btleplug, nfc_manager, rustfft, ring, tokio, flutter_rust_bridge, sled | differential·fingerprint·ai·ble·nfc·dsp·crypto·sync, flutter_rust_bridge | ✅ 부합 | ring→aes-gcm 등 구현체만 대응 |
| 백엔드 | Kong, Go+gRPC, Keycloak, PostgreSQL, TimescaleDB, Redis, Elasticsearch, MinIO, Kafka, K8s, Prometheus+Grafana, ELK | 동일 스택 (Docker Compose 15+ 서비스) | ✅ 부합 | |
| AI/ML | PyTorch, XGBoost, Transformer, YOLOv8, Whisper, VITS, NLLB-200, Claude API, TFLite, Flower, MLflow, Milvus | TFLite(엣지), Milvus, ml-model-design-spec에 5종 모델·연합학습 | 🟡 부분 | 학습·번역·음성 등은 Phase 2+ |
| 실시간 통신 | WebRTC(Janus), Whisper 자막, NLLB 번역, Matrix, FCM, HLS/DASH | 미구현 | ❌ 미반영 | Phase 3 이후 |

**결론**: Phase 1 MVP에 해당하는 기술 스택은 원본에 부합. AI/ML·실시간 통신은 단계별 반영 예정으로 보완 문서화 권장.

---

### 2.4 IV. 전체 시스템 아키텍처

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| MSA 도메인 | 8개 도메인, 30+ 서비스 (auth, user, family, device, measurement, cartridge, calibration, ai-inference, ai-training, coaching, vision, nlp, health-record, telemedicine, reservation, prescription, shop, payment, subscription, marketplace, admin, inventory, logistics, analytics, community, translation, video, notification, iot-gateway, location, emergency) | Proto·구현: auth, user, device, measurement 4개 서비스. family는 user/DB 스키마에 포함 | 🟡 부분 | 원본의 “완전한 MSA”는 장기 목표, 현재는 MVP 4서비스 |
| 데이터 흐름 | 카트리지→RAFE→STM32→BLE→Rust(차동측정·DSP·AI·핑거프린트)→로컬 저장→API Gateway→measurement/ai/coaching/health-record→TimescaleDB+Milvus→UI | AGENTS·CONTEXT 아키텍처도 동일 흐름, Rust 8모듈·Go 4서비스·TimescaleDB·Milvus 반영 | ✅ 부합 | |
| 오프라인-온라인 | CRDT 동기화, Delta Sync, 충돌 해소 LWW + Vector Clock | sync 모듈 CRDT(GCounter, LWWRegister, ORSet), SyncManager | ✅ 부합 | |

**결론**: 아키텍처 방향과 데이터 흐름·오프라인 동기화는 원본 부합. 30+ MSA는 “현재 4개 + 단계별 확장”으로 기획안에 명시하는 것이 좋음.

---

### 2.5 V. 핵심 기능 상세 설계 (5.1~5.14)

| 절 | 원본 요약 | 현재 반영 | 부합도 |
|----|----------|----------|--------|
| 5.1 SaaS 구독·개인 주치의 | Basic/Bio-Optim/Clinical 3티어, 개인 주치의 AI 아키텍처 | 구독 4티어(Free 포함), UserService.GetSubscription, TierConfig | ✅ 부합 | |
| 5.2 다수 리더기 | 최대 10대 BLE, Wi-Fi Direct, 위치 대시보드 | DeviceService, ListDevices, 구독별 max_devices | ✅ 부합 | 위치/대시보드는 Flutter·Phase 2 |
| 5.3 계층형 관리자 | 총괄→국가→지역→지점→판매점, 기능 매트릭스 | 미구현 (admin-service 없음) | ❌ 미반영 | Phase 3, 기획안에는 유지 |
| 5.4 화상진료/병원·약국 예약 | WebRTC, 건강 데이터 공유, 예약 흐름 | 미구현 | ❌ 미반영 | Phase 3 |
| 5.5 통합 쇼핑몰 | 카트리지·리더기·구독·전문가 서비스 | 미구현 | ❌ 미반영 | Phase 2+ |
| 5.6 SDK·카트리지 마켓플레이스 | manpasik-sdk 구조, 검증·수익분배 워크플로우 | sdk/ 디렉토리 예정, 문서 없음 | ❌ 미반영 | Phase 4, 기획안 유지 |
| 5.7 건강관리 코칭 | AI 코칭, 음식/운동 분석, 게이미피케이션 | ml-model-design-spec, coaching 설계 | 🟡 부분 | 구현은 Phase 2 |
| 5.8 오프라인 완전 구동 | 기능 매트릭스(온라인/오프라인/구현방식) | CRDT, 로컬 저장, TFLite 로컬 | ✅ 부합 | 매트릭스 문서화 권장 |
| 5.9 카트리지 자동인식 | NFC 자동읽기, 타입별 프로토콜, 다중 리더기 관리 | NFC 태그 구조, CartridgeType 29종, BLE 패킷 | ✅ 부합 | |
| 5.10 글로벌 규제 | 5개국 규제 매트릭스, TPM·해시체인·AES·감사추적 | regulatory-compliance-checklist, STRIDE, AES-256-GCM, 해시체인 | ✅ 부합 | |
| 5.11 글로벌 커뮤니티·실시간 번역 | NLLB-200, Whisper, 실시간 자막/번역 | 미구현 | ❌ 미반영 | Phase 3+ |
| 5.12 자기학습 성장 AI | 연합학습, 지속학습, 성장 지표 | ml-model-design-spec (연합학습, PCCP) | 🟡 부분 | 설계만 반영 |
| 5.13 음성 명령/접근성 | Whisper STT, NLU, TTS, 접근성 기능 | 미구현 | ❌ 미반영 | Phase 4·5 |
| 5.14 유기적 확장 | Phase 4+ 확장 영역, SDK·OTA·E12-IF | 29종 카트리지, 확장 경로 | ✅ 부합 | |

**결론**: Phase 1에 해당하는 측정·오프라인·카트리지·규제는 원본에 부합. 관리자·쇼핑·SDK·번역·음성 등은 단계별 반영이면 되며, 기획안에 “기능–Phase 매핑”을 명시하면 좋음.

---

### 2.6 VI. 사이트맵 & VII. 스토리보드

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| 사이트맵 | 홈, 측정, 데이터허브, AI코치, 마켓, 커뮤니티, 의료서비스, 기기관리, 가족, 설정, 관리자 포탈 | frontend.mdc·flutter-ui-spec에 라우트 12개 등 | 🟡 부분 | 원본과 동일 트리 구조의 단일 “사이트맵” 문서 없음 |
| 스토리보드 | "첫 측정" 6장면, "음식 촬영→칼로리" 3장면 | 없음 | ❌ 미반영 | UX·QA·인수테스트에 유리하므로 문서화 권장 |

**결론**: 라우트·화면은 스펙에 분산 반영. 원본 수준의 **사이트맵 전용 문서**와 **스토리보드**를 추가하면 원본 충실도와 실행 연계가 좋아짐.

---

### 2.7 VIII. AI 에이전트 자동화 설계

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| UI 에이전트 | 음성 명령, 자연어 대화, 음식/운동 영상, 다국어 번역 | 미구현 | ❌ 미반영 |
| 측정 자동화 에이전트 | 카트리지 인식·프로토콜 선택, RAFE 구성, 차동측정, 재측정 트리거 | BLE/NFC·차동측정·프로토콜은 Rust에 구현 | 🟡 부분 | “에이전트”로서의 자동화 시나리오 문서 없음 |
| 건강 관리 에이전트 | 기준선, 위험 예측, 코칭 생성, 리마인더, 의료기관 추천 | coaching·위험 예측은 설계 단계 | 🟡 부분 |
| 시스템 관리 에이전트 | OTA, AI 모델 업데이트, 보정 동기화, 재고 발주 | DeviceService.RequestOtaUpdate 등 | 🟡 부분 |
| 긴급 대응 에이전트 | 위험 감지 알림, 연락망 연락, AI 음성통화, 119 연동 | 미구현 | ❌ 미반영 |

**결론**: “AI 에이전트”는 원본의 자동화 시나리오·역할 구분이 현재 문서 체계에 없음. 기획안 VIII을 유지하고, **에이전트–Phase–서비스 매핑** 문서를 추가하는 것을 권장.

---

### 2.8 IX. 데이터 아키텍처

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| 데이터 패킷 (패밀리C) | header(device_id, lot_id, fw_ver, cartridge_id, session_id, timestamp), payload(raw_channel, result{differential_correction}, env_meta, state_meta), footer(checksum, schema_ver, transform_log) | Proto에 DifferentialCorrection, EnvironmentMeta, MeasurementData 등 유사 필드 | 🟡 부분 | **footer(checksum, schema_ver, transform_log)** 및 state_meta는 Proto/문서에 미상세 |
| DB 스키마 | users, families, devices, cartridges, measurements, measurement_results, fingerprint_vectors, health_records, ai_coaching_logs, food_analysis_logs, exercise_logs, subscriptions, orders, admin_hierarchy, inventory, community_posts, translations, audit_logs | AGENTS에 users, subscriptions, family_groups, devices, device_events, measurements(TimescaleDB), Milvus, Redis | 🟡 부분 | cartridges, health_records, audit_logs 등은 미정의 또는 별도 설계 필요 |

**결론**: 패밀리C 표준 패킷의 **전체 구조(header/payload/footer)와 transform_log·state_meta**를 한 문서로 고정하고, DB 테이블 목록을 원본과 매핑·갱신하는 것이 좋음.

---

### 2.9 X. 보안 및 규제 준수

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| 보안 계층 | TPM, 펌웨어 보안 부팅, BLE AES-CCM, HTTPS TLS 1.3, AES-256-GCM, 해시체인, 감사추적, OIDC, MFA, RBAC, 프라이버시 | STRIDE, 데이터 보호 정책, AES-256-GCM·해시체인 구현, Keycloak OIDC, RBAC 매트릭스 | ✅ 부합 | TPM·보안 부팅은 하드웨어 영역 |
| 규제 | 5개국 인증·데이터보호·품질·SW수명주기·사이버보안·임상 | regulatory-compliance-checklist, IEC 62304 Class B, ISO 14971, V&V 마스터 플랜 | ✅ 부합 | |

**결론**: 보안·규제 설계는 원본에 부합. 구현 단계에서 감사추적·키관리 상세화만 보완하면 됨.

---

### 2.10 XI. 개발 로드맵 & XII. 비용 분석

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| Phase 1 (1–4개월) | 기본 UI, BLE/NFC 연동, 차동측정, 88차원, 시각화, 클라우드 연동, 오프라인 기본 | Rust 코어·FFI 완료, Go 4서비스 비즈니스 로직 35%, Flutter 대기, Docker·규정 문서 | ✅ 부합 | 주차별 진행은 CONTEXT 체크리스트와 대체로 일치 |
| Phase 2~5 | Core(5–8), Advanced(9–12), Ecosystem(13–18), Evolution(19–24) | AI_COLLABORATION·CONTEXT에 Phase 구분 | ✅ 부합 | |
| 비용·인력 | 약 67억원, 피크 32명 | AGENTS에 동일 수치 | ✅ 부합 | |

**결론**: 로드맵·비용은 원본과 일치. “현재 Phase 1 내에서의 주차–산출물” 매핑만 명시하면 추적이 쉬워짐.

---

### 2.11 XIII. 가상 시뮬레이션 검증 & XIV. 전문가 패널 & XV. 결론

| 항목 | 원본 | 현재 | 부합도 |
|------|------|------|--------|
| 부하 시뮬레이션 | 100,000 동시 사용자, API/측정/AI/TimescaleDB/WebSocket/Milvus | 미구현 | ❌ 미반영 | 인프라·성능 검증 시 재사용 권장 |
| 오프라인 72시간 시뮬레이션 | 1,000회 측정, 동기화 큐 1,000건 | 미구현 | ❌ 미반영 | |
| 896차원 검색 시뮬레이션 | 100만 벡터, Milvus IVF_SQ8 | 미구현 | ❌ 미반영 | |
| 전문가 패널 검증 | 25인 만장일치, 12개 항목 점수 | 원본 문서 보존 시 그대로 유효 | ✅ 부합 | |
| 결론·차별화 포인트 | 5대 가치, 7대 차별화 | CONTEXT·AGENTS에 반영 | ✅ 부합 | |

**결론**: 시뮬레이션은 원본 기획안의 “검증 방법”으로 유지하고, Phase 2 이후 실제 부하·오프라인·벡터 검색 테스트 시 재사용하는 수준이면 충분.

---

## 3. 보완사항 정리 (원본에 충실하기 위해 필요한 작업)

### 3.1 문서 보완

| 우선순위 | 보완 항목 | 설명 |
|----------|----------|------|
| P0 | **데이터 패킷 표준 문서** | 원본 IX 절 패밀리C 패킷(header/payload/footer, transform_log, state_meta)를 `docs/data-architecture/` 또는 `docs/specs/`에 단일 문서로 정의. Proto와의 필드 매핑 명시. |
| P0 | **사이트맵 공식 문서** | 원본 VI 절 트리 구조를 그대로 `docs/sitemap.md`(또는 `docs/ux/sitemap.md`)로 두고, 라우트 ID·화면명·Phase 매핑 추가. |
| P1 | **스토리보드 문서** | 원본 VII 절 "첫 측정"·"음식 촬영→칼로리" 스토리보드를 `docs/ux/storyboard-*.md`로 보존. 추가 시나리오(리더기 추가, 구독 변경 등)는 Phase별로 확장. |
| P1 | **오프라인 기능 매트릭스** | 원본 5.8 표(기능×온라인/오프라인/구현방식)를 `docs/specs/offline-capability-matrix.md`로 정리. |
| P1 | **AI 에이전트–Phase 매핑** | 원본 VIII 절 에이전트 목록을 유지하고, 각 에이전트를 담당 Phase·서비스(또는 모듈)와 매핑한 표 추가. |
| P2 | **MSA 확장 로드맵** | 원본 4.1의 30+ 서비스를 Phase 2/3/4별로 “도입 시점·담당 서비스명”으로 정리한 표. |

### 3.2 구현·스키마 보완

| 우선순위 | 보완 항목 | 설명 |
|----------|----------|------|
| P0 | **감사 추적(audit_logs)** | IEC 62304·원본 요구사항에 맞춰 PHI 접근·주요 이벤트용 audit_logs 스키마 및 저장 정책(10년 보존) 정의 및 반영. |
| P0 | **tenant_id (멀티테넌시)** | 프로젝트 규칙상 tenant_id 필수이므로, users, devices, measurements 등 핵심 테이블에 tenant_id 추가 및 문서화. |
| P1 | **Proto와 패밀리C 정합성** | MeasurementData 등에 state_meta, footer 용 checksum/schema_ver/transform_log 대응 필드 또는 메시지 추가 검토. |
| P1 | **cartridges·health_records 테이블** | 원본 IX DB 목록에 맞춰 cartridges, health_records 테이블 설계(또는 기존 스키마와의 매핑) 문서화. |

### 3.3 기획안 자체 보완 (원본 v1.0 유지 + 부록)

- **현행 시스템 매핑표**: “원본 절 번호 ↔ 현재 문서/코드 위치” 한 페이지 표로 정리해 기획안 부록 또는 별도 `docs/plan-traceability.md`로 유지.
- **구독 티어 공식 명칭**: 원본 3단계(Basic Safety, Bio-Optimization, Clinical Guard)와 현재 4단계(Free, Basic, Pro, Clinical)의 대응 관계를 기획안 또는 CONTEXT에 한 줄로 명시.

---

## 4. 기획안 발전 수립안 (더 완벽한 시스템 구축을 위한 제안)

다음 내용을 반영하면 기획안이 “실행 가능한 마스터 문서”로 한 단계 발전합니다.

### 4.1 기획안 v1.1 권장 구조 (추가·갱신 섹션)

1. **현행 시스템 매핑표 (신설)**  
   - 원본 I~XV 각 절(또는 주요 표/목록)에 대해 “현재 반영 위치”(문서 경로, repo 경로, Phase)를 한 표로 정리.  
   - 예: `IV. 아키텍처 4.1 MSA` → `backend/services/`, `AGENTS.md` §3, Phase 1: 4개 서비스, Phase 2–4: 확장 로드맵 참조.

2. **단계별 MSA 확장 로드맵 (IV절 보강)**  
   - Phase 1: auth, user, device, measurement (현행).  
   - Phase 2: subscription, shop, payment, coaching, ai-inference 등.  
   - Phase 3: telemedicine, reservation, community, translation, video, notification, admin 등.  
   - Phase 4: marketplace, inventory, logistics, analytics, iot-gateway, location, emergency 등.  
   - 위를 표로 정리해 “원본 4.1 도메인 맵”과 일대일로 연결.

3. **기획–구현 추적성 매트릭스 (신설 또는 부록)**  
   - 요구사항 ID(원본 절·표·항목 번호) ↔ 설계 문서 ↔ 코드/Proto/DB ↔ 테스트·V&V 문서.  
   - 의료기기 인허가·감사 대비용으로 유지.

4. **사이트맵·스토리보드 공식 위치**  
   - 기획안 VI·VII을 “원본 명세”로 두고, “실제 반영 문서”를 `docs/ux/sitemap.md`, `docs/ux/storyboard-*.md`로 고정. 기획안에 해당 경로만 명시해도 됨.

5. **데이터 패킷 표준 단일 참조**  
   - IX절 패밀리C 패킷을 `docs/specs/data-packet-family-c.md`(또는 동일 성격 문서)로 고정하고, 기획안 IX에서는 “상세는 해당 문서 참조”로 통일.

6. **시뮬레이션 검증 재사용 계획 (XIII절 보강)**  
   - 부하·오프라인·896차원 시뮬레이션을 “Phase 2/3 인프라·성능 검증 시 수행”으로 명시하고, 수용 기준(예: 지연, 처리량)만 기획안에 유지.

7. **에이전트 자동화 (VIII) – Phase 매핑**  
   - 각 에이전트 블록(UI, 측정, 건강관리, 시스템관리, 긴급대응)에 “담당 Phase·주요 서비스/모듈” 컬럼을 추가한 표.  
   - 구현 담당(Flutter/Rust/Go/외부) 표기 시 추적 용이.

8. **용어·티어 통일**  
   - Basic Safety / Bio-Optimization / Clinical Guard ↔ Free·Basic·Pro·Clinical 대응표를 기획안 5.1 또는 총괄 개요에 추가.  
   - 이후 신규 문서는 이 표를 기준으로 사용.

### 4.2 제안 요약

- **원본 v1.0**: 그대로 “승인된 기획 명세”로 보존.  
- **v1.1**: 위 4.1 항목을 “추가·보강”만 반영한 개정안으로 두고,  
  - 현행 시스템 매핑,  
  - 단계별 MSA 확장 로드맵,  
  - 기획–구현 추적성 매트릭스,  
  - 사이트맵/스토리보드/데이터패킷/오프라인 매트릭스 등 **실행 연계 문서 참조**  
  를 명시하면, “원본에 충실하면서도 구현과 단계별 확장이 연결된” 완성도 높은 기획안이 됩니다.

---

## 5. 결론

- **원본 부합도**: 핵심 기술(차동측정, 88→896차원, 오프라인, CRDT, 보안·규제 방향), 아키텍처(데이터 흐름, 클라우드 스택), 로드맵·비용·철학은 현재 시스템에 잘 반영되어 있음.  
- **보완**: 데이터 패킷(패밀리C) 상세, 사이트맵·스토리보드 공식 문서화, MSA 확장 현황·로드맵, AI 에이전트–Phase 매핑, DB(tenant_id, audit_logs, cartridges, health_records) 정리.  
- **발전**: 기획안 v1.1에서 “현행 매핑표”, “단계별 MSA 확장”, “추적성 매트릭스”, “실행 연계 문서 경로”를 추가하여 원본에 충실한 동시에 실행·인허가 대비를 강화하는 것을 권장합니다.

---

**문서 종료**

*본 검증은 MPK-ECO-PLAN-v1.0-20260208-FINAL 원본과 현재 코드베이스(CONTEXT.md, AGENTS.md, backend, rust-core, docs) 및 관련 규정·보안 문서를 대조하여 작성되었습니다.*
