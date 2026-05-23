# 만파식(MMUP) 빌드 글로벌 규칙

## 규칙 계층 원칙 (확장성 보장)
- **불변 규칙(헌법)**: 정체성·계약(인터페이스·패킷)·안전/규제 불변식·표기 규칙만 글로벌 규칙에 고정한다.
- **베이스라인 파라미터(외부화·교체가능)**: 특정 금액·수수료율·환율 수치·부품 IC·RTOS·대수·차원/정확도 수치는 **규칙에 박지 않고** `config/baseline_params`(버전화·[추정] 플래그·12-factor/시크릿 외부화, EX)로 둔다. 변경은 §4 진화 원칙·사람 승인이며 **규칙 개정 불요**. 코드·규칙은 이를 값이 아닌 **역량·키로 참조**한다.

## 정체성·언어 (불변)
- 플랫폼 공식명: MPS 만파식 다중측정원리 유니버설 POCT 플랫폼(MMUP). 금지표현: "유니버설 POCT 분석장치"·"소형 정량 면역분석기"·"혈액 전용 분석기".
- 산출 문서·UI 카피·주석은 한글, 코드·식별자·경로는 영어. **금액은 ₩ 표기, 외화 괄호 병기**(환율 수치는 파라미터 `fx.krw_per_usd`, 현행 [추정] 1,480).
- 불확실성 표기: [검증됨]/[추정]/[미검증]/[추측] + 출처.

## 불변 계약·불변식 (변경 시 어댑터+버전업·사람 승인)
- 카트리지 인터페이스 계약: CSI v1.0 — 14 신호핀(7×2) + 2 비도금 마운트홀, 1.27mm 피치. ("12핀/E12-IF" 표기 금지). 부품(MECF-08 등)은 파라미터·footprint 검증(§아키텍처 불변식 커넥터 규칙).
- 차동측정 식: Sdiff_n = S_n − α_n × R_n. (계수 α는 파라미터)
- 측정 차원 곡선: 88 → 448 → 896 → 1792차원+ (가변 길이 스키마·N차원 개방; 마일스톤 수치는 파라미터).
- 표준 데이터 패킷·전역 ID·이벤트·계약 우선(아래 아키텍처 불변식).

## 베이스라인 파라미터 (현행값·교체가능·[추정], `config/baseline_params`)
> 아래는 **현재 작업기준점**일 뿐 불변 규칙이 아니다. 변경은 규칙 개정 없이 파라미터 갱신 + 사람 승인.
- 카트리지 SKU: 현행 44종(Layer 3) — 커머스·재고·구독은 **확장 수용 구조**로 설계(수치 하드코딩 금지).
- 측정 정확도(결합): 현행 92~98% [검증됨/SSOT, 임상 미검증]. (학습 향상 92→95→98%는 별개 축)
- 리더기 등록 대수: 하드 캡 없음(구독 티어). 동시 BLE 연결은 SoC 한계(Apple 7~10 / Nordic ~20)를 허브·스케줄링·PAwR로 우회.
- 가격·수수료: 구독가·판매수수료·수익배분(예 70:30)·예산은 모두 파라미터·[추정], UI/로직에 하드코딩 금지.
- 부품·RTOS 후보: MCU(STM32F405급)·BLE(nRF52832/40급)·NFC(PN7150/PN7642급)·RTOS(FreeRTOS→SafeRTOS/ThreadX·Zephyr) — 전부 **HAL/역량 명세 뒤 현재 후보**, 데이터시트·HW SSOT로 확정.

## 기술 스택 (현행 채택·버전화, 부품/IC는 HAL 뒤·교체가능)
- 앱: Flutter 3.x + Riverpod + Material 3 + flutter_localizations + WCAG 2.2 AA.
- 코어: Rust(no_std 차동측정, btleplug, nfc, rustfft, ring, tokio, sled, flutter_rust_bridge).
- 펌웨어: 리더기 임베디드 OS(MMUP-OS) — RTOS 후보(파라미터)·embedded-hal·no_std Rust(§2.1·E-FW).
- 백엔드: Go + gRPC, Kong, Keycloak(OIDC/MFA/RBAC), PostgreSQL 16, TimescaleDB, Redis, Milvus, Kafka, MinIO, Kubernetes(EKS).
- AI: PyTorch, XGBoost, YOLOv8, Whisper, VITS, NLLB-200, SeamlessM4T, Flower/NVIDIA FLARE, MLflow.
- 통신: WebRTC(Janus), MQTT, FCM/APNs.
> 스택·라이브러리는 입증된 우위 시 §4 진화 원칙으로 교체 가능(고정 아님). 외부 호출은 어댑터/HAL 뒤로.

## 아키텍처 불변식 (모세혈관 정합성)
- 계약 우선: `contracts/`(proto·packet schema·OMOP/FHIR 매핑·Kafka 이벤트·OpenAPI)가 단일 진실원천. 코드는 계약에서 생성.
- 표준 데이터 패킷: header/payload/footer + transform_log(해시체인). 모든 측정 결과에 confidence·uncertainty 필수.
- 전역 ID: 단말·카트리지·user·org·dataset 전역 고유 ID + 토큰화 매핑.
- 후방호환: 핀맵·커넥터·프로토콜·API 변경 시 어댑터 1세대 의무 제공.
- 커넥터 핀 계약: CSI 핀맵은 **신호핀(기능)으로 정의**하고 정렬핀·키잉·웰드탭 등 기구물은 별도 기계 사양으로 분리. 신호핀 수는 SSOT v1.2 기준 **14(7×2)**이며, 제조사 nominal 16(8×2)과의 정합은 **footprint·데이터시트로 상시 확인**한다. 핀 확장(예: 40핀 MECF-20)은 계약층 변경 → SSOT 버전업·사람 승인.
- 오프라인 우선: 측정·기본 AI·로컬 저장은 100% 오프라인. 동기화는 CRDT+Vector Clock.

## 안전·사실성 게이트 (절대 준수)
- 모든 측정·코칭·번역·예측 산출에 신뢰도·불확실성·출처 동반. 단정 금지.
- 의료: AI 진단·처방 단정 금지 → 인간 주치의/의료진(Human-in-the-loop)으로 핸드오프.
- 비의료 웰니스 코칭과 의료(화상진료)는 분리. 영양제·키트는 질병 효능 단정 금지(건기식법/FDA·FTC).
- 긴급(119)·자동신고·결제·데이터 공유·자기학습은 사용자 동의 + 휴먼인더루프 + 임계값.
- AI 비서: 어조(따뜻함)와 사실성·신뢰성 분리, 반(反)아첨, 과의존·과신 방지, AI임을 고지.
- 외부 인용 ID(PMID·K-number·특허)는 1차 출처 교차검증. 창작 인용은 Critical 위반.
- 규제: QMSR·IVDR·MFDS·HIPAA/GDPR/PIPA·21 CFR Part 11·ISO 13485/14971·IEC 62304 설계 초기 내장.

## 작업 방식
- 작은 단위로 변경, 각 작업에 테스트·수용 기준. 게이트 통과 전 다음 단계 금지.
- 사람 최종 승인: SSOT·BOM·Gate 통과를 에이전트가 단독 확정하지 않는다.