# 만파식 AI 비서·주치의 통합 세부 기획명세 (AI Assistant Master Specification)

**문서번호**: MPK-AI-ASSISTANT-SPEC-v1.0  
**작성일**: 2026-02-12  
**목적**: 사용자가 **텍스트 및 음성 명령**을 통해 모든 기능 수행·설정 변경을 AI 비서(주치의)에게 위임하고, 만파식 생태계 전체를 AI 비서를 통해 손쉽게 운영·관리할 수 있는 시스템의 상세 기획. 관련 자료·논문 조사·분석을 반영하고, 전체 시스템과 유기적으로 통합·학습·성장하도록 설계.  
**상위 문서**: [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md), [FINAL-MASTER-IMPLEMENTATION-PLAN](FINAL-MASTER-IMPLEMENTATION-PLAN.md)

---

## 1. 관련 자료·논문 조사 및 분석

### 1.1 헬스케어 음성·대화형 AI 비서

- **Voice Assistants for Health Self-Management (2024, arXiv)**: 고령자를 위한 건강 자가관리 음성 비서. **개인화**, **맥락 적응**, **사용자 자율성 존중**을 핵심 설계 원칙으로 제시. 의사 방문 노트 디브리핑·맞춤 복약 리마인더.
- **Talk2Care (IMWUT 2024)**: 의료 제공자와 고령자 간 소통을 LLM 기반 음성 비서로 연결. 환자 측 대화형 음성 입력·제공자 측 요약 대시보드. 시간·노력 절감과 함께 **신뢰·안전** 이슈 강조.
- **Personal Health Data Queries (PHD_Query_Response)**: 음성 비서가 **개인 건강 데이터**를 답할 때 사용자는 **완전한 문장 응답**을 선호하며, 기술적·효율적 상호작용을 사회적 친밀도보다 중시.
- **AI-Powered Conversational Agents as Virtual Health Carers (JMIR 2024)**: 비전염성 질환 원격 관리에서의 잠재력과 **신뢰·정확성·안전** 한계. 기능적 요인(유용성·내용 신뢰도·전문가 대비 품질), 보안·프라이버시, 개인 기술 수용도가 신뢰에 영향.

**반영 사항**: AI 비서 응답은 **완전한 문장·기술적 정확성** 우선. 개인화·맥락 적응·사용자 자율성(확인 후 실행) 필수. 의료·설정 변경 시 **확인 플로우**·감사 로그.

### 1.2 작업 지향 대화·함수 호출·헬스케어 에이전트

- **Conversational Health Agents (openCHA, arXiv 2310.02374)**: LLM 기반 **대화형 건강 에이전트(CHA)**. 외부 데이터·지식베이스·분석 모델을 **오케스트레이터**가 계획·실행. 다단계 문제 해결·개인화 대화·멀티모달 분석.
- **Task-Oriented Dialogue through Function Calling (ACL 2025)**: LLM이 **함수 호출**로 외부 도구를 선택적 호출·실행. 정적 지식만 의존하지 않아 확장 가능.
- **CoALM (arXiv 2502.08820)**: 다중 턴 대화와 **도구 사용**을 하나의 모델로 통합. 대화 추론과 API 사용을 인터리브한 학습으로 도메인 특화 모델 대비 성능.
- **Multi-turn Function Calling (BUTTON, arXiv 2410.12952)**: 복잡한 질의를 **다단계·다중 턴** 함수 호출로 처리. compositional instruction tuning.

**반영 사항**: AI 비서는 **의도 → 도구/API 선택 → 실행 → 결과 요약** 오케스트레이션. 단일 턴뿐 아니라 “측정 시작해줘 → (카트리지 확인) → 어떤 카트리지로 할까요?” 등 **다중 턴·정보 수집** 지원. Function Calling(OpenAI/Anthropic 등) 또는 자체 intent–action 매핑.

### 1.3 음성 제어·NLU·온디바이스

- **Speech-Controlled Health Record Management (IWSDS 2025)**: 건강·케어 기록 관리용 음성 인식. 의료 맥락 **Whisper 파인튜닝**으로 WER 16.8→1.0. SNR·오디오 스트림 검사로 환각 감소.
- **On-Device LLMs for Home Assistant (WNUT 2025)**: **의도 검출**과 **응답 생성**을 온디바이스 LLM으로 동시 수행. 노이즈·도메인 외 의도에서 80–86% 정확도. 5–6초/쿼리 수준으로 원샷 명령에 적합.
- **Voice Assistants for Health (Older Adults)**: 접근성·건강 인식·복약 순응 지원. 사용자 자율성·맞춤형 안내.

**반영 사항**: 음성 입력은 **STT(Whisper 또는 파인튜닝)** → 텍스트 → 동일 NLU/오케스트레이터. 응답은 TTS 또는 10.8 같은 음성. 온디바이스 옵션은 엣지 배포 시 고려(지연·정확도 트레이드오프).

### 1.4 에이전트 오케스트레이션·피드백 학습

- **IBM watsonx Orchestrate**: 사용자 요청 이해 → 대화로 정보 수집 → **이산 액션** 실행(조건·분기). AI 어시스턴트에 “액션” 연결.
- **Azure OpenAI Assistants API**: **다중 도구** 병렬 접근, 코드 해석·함수 호출 등. 플랜·실행·검증 분리 가능.
- **TypeAgent (Microsoft)**: 자연어 → **논리 구조**로 변환, LLM과 기존 소프트웨어 안전 결합. **대화 메모리·RAG**로 과거 상호작용 학습·맥락 적용.
- **Intent (Augment)**: 코디네이터 에이전트가 계획 제안 → 구현 에이전트 병렬 실행 → 검증 에이전트가 결과 검증 → **사람 검토** 전 제출. 피드백 루프로 품질 개선.

**반영 사항**: **오케스트레이터**가 “의도+엔티티” → 실행 계획(단일/다단계) → 기존 gRPC·REST 호출. **확인 필요 액션**(결제·설정 변경·삭제)은 사용자 확인 후 실행. **대화 메모리**로 문맥·선호도 축적 → 10.9 익명 학습과 분리된 “개인 세션 맥락”으로 활용. 학습·성장은 2.5절.

---

## 2. AI 비서·주치의 범위 및 목표

### 2.1 목표

- 사용자가 **직접 수행해야 하는 모든 이벤트**(기능 실행·설정 변경·조회·관리)를 **텍스트 및 음성 명령**으로 AI 비서(주치의)에게 요청하면, AI가 해당 기능을 **대신 수행**하거나 **설정을 변경**하여, 만파식 생태계를 **AI 비서를 통해 손쉽게 운영·관리**할 수 있게 한다.
- AI 비서는 기존 **CoachingService·AiInferenceService·10.2 AI 주치의**와 동일 페르소나로, “건강 코칭”뿐 아니라 **시스템 전역 액션 실행**으로 확장된다.
- 전체 시스템에 **유기적으로 통합**되고, 사용 패턴·피드백·익명 집계(10.9)를 통해 **학습하고 성장**한다.

### 2.2 사용자 직접 수행 이벤트 망라 (AI 비서 대행 대상)

| 도메인 | 사용자 이벤트(예) | 대행 시 AI 비서 동작 |
| --- | --- | --- |
| **측정** | 측정 시작·종료·이력 조회 | StartSession(파라미터 수집)·EndSession·GetMeasurementHistory 호출 후 결과 요약 |
| **건강·코칭** | 목표 설정·코칭 요청·리포트 조회 | SetHealthGoal·GenerateCoaching·GetWeeklyReport 등 호출·요약 |
| **식단·칼로리** | 식사 기록·칼로리 분석 요청 | LogMeal·AnalyzeFoodImage(이미지 업로드 시) 호출 |
| **기기** | 기기 등록·목록·OTA·컨셉 할당 | RegisterDevice·ListDevices·RequestOtaUpdate·AssignDeviceToConcept |
| **마켓·결제** | 상품 검색·장바구니·주문·결제 | ListProducts·AddToCart·CreateOrder·CreatePayment(확인 후)·ConfirmPayment |
| **예약·의료** | 시설 검색·예약·화상 입장·처방 조회 | SearchFacilities·GetAvailableSlots·CreateReservation·JoinRoom·ListPrescriptions |
| **커뮤니티** | 글·댓글·챌린지 | CreatePost·CreateComment·JoinChallenge |
| **가족** | 초대·공유 설정·가족 대시 조회 | InviteMember·SetSharingPreferences·GetSharedHealthData |
| **설정** | 프로필·구독·알림·테마·언어·긴급 연락망 | UpdateProfile·UpdateNotificationPreferences·시스템 설정 변경(테마/언어 등) |
| **알림** | 알림 목록·읽음 처리 | ListNotifications·MarkAsRead |
| **목적별·컨셉** | 컨셉 전환·대시 조회·리더기/카트리지 할당 | ListConcepts·GetStatsForConcept·AssignDeviceToConcept |
| **관리자**(역할 시) | 설정 조회·변경·감사 로그 | ListSystemConfigs·BulkSetConfigs·GetAuditLog(역할 검증) |

위 표는 FINAL-MASTER I.2 매트릭스의 “사용자 행위”를 AI가 **대행 가능한 이벤트**로 매핑한 것이다. 신규 기능 추가 시 동일하게 “사용자 이벤트 → AI 비서 액션”으로 확장.

---

## 3. 시스템 아키텍처 — 유기적 통합

### 3.1 전역 흐름

```text
[사용자] ──텍스트 입력 또는 음성(마이크)──► [STT] ──► [NLU/Intent·Entity]
                                                          │
                                                          ▼
[응답] ◄── TTS(선택) 또는 텍스트 ◄── [응답 생성] ◄── [오케스트레이터] ──► [기존 gRPC/REST]
                                                                              MeasurementService,
                                                                              CoachingService,
                                                                              ShopService, ...
```

- **진입점**: 앱 전역 “AI 비서” 버튼/플로팅·홈 상단 “말하기/쓰기”·AiCoachScreen 내 대화. 음성은 마이크 → STT → 동일 파이프라인.
- **NLU/Intent·Entity**: 사용자 발화에서 “의도”(측정_시작, 설정_알림_변경, 주문_하기 등)와 “엔티티”(카트리지 타입, 날짜, 상품 ID 등) 추출. LLM 기반 또는 규칙+LLM 하이브리드.
- **오케스트레이터**: 의도에 맞는 **도구(함수)** 선택. 도구 = 기존 서비스의 RPC/API 래퍼(StartSession, SetHealthGoal, AddToCart 등). 파라미터 부족 시 **다중 턴**으로 질문(“어떤 카트리지로 측정할까요?”). 확인 필요 액션은 “○○ 할까요?” 후 사용자 승인 시 실행.
- **실행**: Gateway 경유 또는 백엔드 내부 호출로 해당 서비스 호출. 결과(성공/실패/데이터)를 오케스트레이터에 반환.
- **응답 생성**: 결과를 **완전한 문장**으로 요약(1.1 참고). “측정을 시작했어요. 약 90초 후 결과를 알려드릴게요.” TTS 사용 시 10.8 같은 음성 옵션.

### 3.2 기존 시스템과의 통합

- **CoachingService·AiInferenceService**: AI 비서가 “코칭 요청”을 받으면 기존 GenerateCoaching·GetRecommendations 호출. “건강 요약 알려줘” → GetHealthSummary·GenerateDailyReport 조합.
- **10.2 AI 주치의**: 동일 “주치의” 페르소나. AI 비서 = **주치의의 행동 인터페이스**(말/글로 명령하면 주치의가 시스템을 조작).
- **인증·RBAC**: 모든 도구 호출은 **현재 사용자 JWT·역할**로 실행. 관리자 전용 액션(ListSystemConfigs 등)은 역할 검증 후만 노출.
- **이벤트**: 측정 시작·결제·예약 등 실행 시 기존 이벤트(measurement.completed, payment.completed 등) 그대로 발행. AI 비서는 “트리거”만 수행.

### 3.3 학습·성장 (유기적 통합)

- **개인 세션 맥락**: 대화 이력(최근 N턴)·선호(“항상 혈당 카트리지로 측정”)를 **세션/사용자별 메모리**에 저장. TypeAgent 스타일 RAG로 “이전에 사용자가 ~라고 했음” 반영. 개인 데이터이므로 암호화·접근 제어.
- **익명 학습(10.9)**: “의도–액션–결과” 집계를 **익명**으로만 수집(어떤 의도가 많았는지, 실패율, 다중 턴 비율). 모델·정책 개선에 활용. 개인 식별 없음.
- **피드백**: “잘 했어/잘못됐어” 또는 암묵적(사용자가 즉시 수동으로 같은 작업 재시도). 피드백 시 해당 턴·의도·결과를 (익명) 학습 파이프라인에 반영해 **응답 품질·의도 해석** 개선.
- **성장**: 주기적 재학습·A/B 테스트로 “새 의도 추가”“새 도구 연결”“확인 플로우 개선” 반영. 10.9 생물형 AI 생태계와 동일 원리로 “공통 AI”가 성장하고, 개인 맥락은 세션 메모리로만 유지.

---

## 4. 의도·도구·확인 정책 상세

### 4.1 의도(Intent) 분류 체계

- **측정_***: 측정_시작, 측정_이력_조회, 측정_결과_요약.
- **건강_***: 건강_요약, 코칭_요청, 목표_설정, 리포트_조회, 추천_요청.
- **식단_***: 식사_기록, 칼로리_분석_요청.
- **기기_***: 기기_등록, 기기_목록, 기기_OTA, 기기_컨셉_할당.
- **마켓_***: 상품_검색, 장바구니_추가, 주문_생성, 결제_진행.
- **예약_***: 시설_검색, 예약_생성, 예약_목록, 화상_입장.
- **처방_***: 처방_목록, 복약_리마인더_조회.
- **커뮤니티_***: 글_작성, 댓글_작성, 챌린지_참가.
- **가족_***: 가족_초대, 공유_설정, 가족_건강_조회.
- **설정_***: 프로필_수정, 알림_설정, 테마_변경, 언어_변경, 긴급_연락망_설정.
- **알림_***: 알림_목록, 알림_읽음처리.
- **컨셉_***: 컨셉_전환, 대시_통계_조회, 리더기_할당.
- **관리자_***: (역할=admin) 설정_조회, 설정_변경, 감사_로그_조회.
- **일반_***: 인사, 도움말, 취소, 확인_승인.

엔티티(slot): device_id, cartridge_type, concept_id, date_range, product_id, facility_id, reservation_id, theme, language 등. 의도별 필수/선택 슬롯 정의.

### 4.2 도구(Tool) 매핑

- 각 의도에 대해 **호출할 RPC/API** 1:1 또는 1:N(다단계). 예: 측정_시작 → CartridgeService.ReadCartridge(선택) + MeasurementService.StartSession(device_id, cartridge_uid, concept_id).
- **파라미터 바인딩**: 엔티티에서 추출한 값 + 기본값(현재 사용자, 오늘 날짜 등). 부족 시 다중 턴 질문.
- **도구 결과** → 오케스트레이터가 성공/실패/데이터 해석 후 응답 문장 생성.

### 4.3 확인 정책

- **확인 필수**: 결제(ConfirmPayment), 구독 해지, 설정 대량 변경(BulkSetConfigs), 가족 멤버 제거, 예약 취소 등. “○○할까요?” → 사용자 “예/응/해줘” 등으로 승인 시에만 실행.
- **확인 생략 가능**: 조회(이력·요약·목록), 측정 시작(이미 “측정 시작해줘”로 명시적 요청), 알림 읽음 처리 등. 정책으로 “위험도 낮은 액션”만 생략.
- **취소**: 사용자 “아니야/취소” 시 실행하지 않고 “취소했어요” 응답.

---

## 5. 음성·텍스트 통일 및 접근성

- **텍스트**: 채팅 UI에서 입력 → 동일 NLU·오케스트레이터. 응답은 텍스트.
- **음성**: 마이크 → STT(Whisper 또는 파인튜닝, 1.3 참고) → 텍스트 → 동일 파이프라인. 응답은 TTS(기본) 또는 10.8 “같은 음성” 옵션. “음성만 쓰기” 사용자도 모든 기능 수행 가능.
- **접근성**: 스크린 리더·키보드만으로도 “AI 비서 열기 → 텍스트 명령” 가능. 음성은 핸즈프리·시각 장애 사용자 지원.

---

## 6. 구현·구축 제안 (요약)

- **서비스**: AssistantService(또는 AiInferenceService 확장). 진입: ProcessUserInput(text or audio_url) → intent·entities → orchestrate → tool calls → response_text(, audio_url). 대화 세션·메모리 저장(assistant_sessions, assistant_turns).
- **DB**: assistant_sessions(session_id, user_id, started_at), assistant_turns(turn_id, session_id, role: user/assistant, content, intent, tool_calls, tool_results, created_at). 개인 데이터 암호화·보존 기간 정책.
- **클라이언트**: Flutter “AI 비서” 플로팅 또는 전용 화면. 텍스트 필드 + 마이크 버튼. STT/TTS는 플러그인 또는 백엔드 연동.
- **학습**: 턴 로그를 익명화 파이프라인(10.9)에 “의도·성공여부·다중턴여부”만 집계 입력. 주기적 의도 분류기·응답 생성기 재학습.

---

## 7. 참조

- **Blueprint**: Part 10.2 AI 주치의, 10.4·10.8 음성·번역, Part 5 인터페이스, Part 6 페이지.
- **FINAL-MASTER**: I.2 기능↔API 매트릭스, II.1~II.10 사용자 이벤트별 세부.
- **규제**: 의료·건강 데이터 음성 처리 시 개인정보보호법·동의. AI 비서 행위에 대한 “보조” 용도 명시(진단·처방 판단은 의료진 책임).

---

**문서 끝.**
