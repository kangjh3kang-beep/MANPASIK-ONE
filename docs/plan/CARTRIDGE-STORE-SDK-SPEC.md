# 만파식 카트리지 스토어 & 개발자 SDK 상세 구현 기획명세

**문서번호**: MPK-CART-STORE-SDK-SPEC-v1.0  
**작성일**: 2026-02-12  
**목적**: Apple App Store 모델을 참고하여, 만파식 생태계 리더기와 호환되는 카트리지를 누구나 개발·등록·판매할 수 있는 카트리지 스토어와 오픈 SDK를 상세 기획한다.  
**상위 문서**: [COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0](COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md), [FINAL-MASTER-IMPLEMENTATION-PLAN](FINAL-MASTER-IMPLEMENTATION-PLAN.md), [MEASUREMENT-ANALYSIS-AI-SPEC](MEASUREMENT-ANALYSIS-AI-SPEC.md)

---

## 1. 관련 자료·사례 분석

### 1.1 벤치마크 플랫폼 비교

| 플랫폼 | 모델 | 핵심 특징 | 만파식 적용점 |
| --- | --- | --- | --- |
| **Apple App Store** | 앱 마켓플레이스 | 개발자 SDK(250K+ API) → 앱 심사 → 배포 → 수익 배분(70/30 또는 85/15). MFi 하드웨어 인증 프로그램. AccessorySetupKit으로 BLE 페어링 간소화 | 카트리지 SDK → 인증 심사 → 스토어 배포 → 수익 배분. NFC 인식으로 "플러그앤플레이" |
| **Biocartis (Idylla)** | 컨텐츠 파트너링 | 외부 개발자가 자사 어세이(시약)를 Idylla 하드웨어에 탑재. IP·브랜드·규제 주체 = 개발자. 75+개국 2,000+ 설치 기반 | 카트리지 개발자가 IP·브랜드 소유. 만파식 리더기 호환 인증만 통과하면 등록 가능 |
| **Fluxergy** | OEM 어세이 플랫폼 | "We bring the platform, you bring the assay." 공정한 수익 배분, IP 분리 명확. 개발자가 규제 인증 주체 | "We bring the reader, you bring the cartridge." 구조. 개발자 = 규제 주체 |
| **Fathym OpenBiotech** | 클라우드 바이오 플랫폼 | 오픈소스 하드웨어 호환, SDK·API·대시보드 제공. Azure 기반 데이터 파이프라인 | 클라우드 데이터 파이프라인 + SDK 제공 모델 참고. 개발자 콘솔·API·문서 |
| **OpenBCI** | 오픈소스 하드웨어 | Cyton(8ch) SDK 문서 공개. 서드파티 하드웨어(EmotiBit·MyoWare 등) 연동 | 오픈 하드웨어 인터페이스 명세 공개. 서드파티 센서 모듈 허용 가능성 |
| **DXRX (Diaceutics)** | 진단 마켓플레이스 | 진단 테스트 상업화·유통 마켓플레이스. 표준화된 구현 서비스 | 카트리지 상업화·유통 인프라 참고 |

### 1.2 설계 원칙 도출

1. **IP 분리**: 카트리지 설계·시약 IP는 개발자 소유. 만파식은 리더기·플랫폼·SDK IP 보유.
2. **인증 + 자유**: 안전·호환성 인증은 필수, 카트리지 아이디어·용도는 자유 (App Store 심사 모델).
3. **수익 배분**: Apple 모델(70/30 기본, 소규모 개발자 85/15 우대) 참고.
4. **규제 책임**: 의료용 카트리지는 개발자가 규제 인증(KFDA/CE/FDA) 주체. 만파식은 리더기 인증.
5. **플러그앤플레이**: NFC 태그로 카트리지 자동 인식 → 적합 측정 프로파일 로드 → 즉시 측정.

---

## 2. 카트리지 스토어 전체 구조

### 2.1 생태계 전체도

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                       만파식 카트리지 스토어 생태계                             │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌──────────────┐  │
│  │  카트리지     │    │  개발자       │    │  카트리지     │    │  사용자       │  │
│  │  개발자       │──►│  콘솔 &       │──►│  스토어       │──►│  (앱 내       │  │
│  │  (3rd Party) │    │  SDK          │    │  (마켓플레이스)│    │   구매·사용) │  │
│  └─────────────┘    └─────────────┘    └─────────────┘    └──────────────┘  │
│        │                   │                   │                   │          │
│        │                   │                   │                   │          │
│  ┌─────▼─────────────────▼───────────────────▼───────────────────▼────────┐ │
│  │                        만파식 플랫폼 백엔드                               │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐ │ │
│  │  │Cartridge │ │Developer │ │Store     │ │Review    │ │Payment &    │ │ │
│  │  │Registry  │ │Service   │ │Service   │ │Service   │ │Revenue Share│ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────────┘ │ │
│  └───────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        만파식 리더기 하드웨어                            │  │
│  │  NFC 인터페이스 │ BLE 통신 │ 전극 어레이 │ ADC │ MCU                    │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 핵심 서비스 구성

| 서비스 | 역할 | 주요 API |
| --- | --- | --- |
| **CartridgeRegistryService** | 카트리지 타입 등록·메타데이터·버전 관리 | RegisterCartridgeType, UpdateCartridgeType, GetCartridgeSpec, ListCartridgeTypes |
| **DeveloperService** | 개발자 계정·인증·팀·키 관리 | RegisterDeveloper, CreateApiKey, GetDeveloperProfile, ListTeamMembers |
| **StoreService** | 스토어 리스팅·검색·구매·다운로드 | ListStoreItems, SearchCartridges, GetStoreItem, PurchaseCartridge, GetPurchaseHistory |
| **ReviewService** | 카트리지 심사·인증·승인 워크플로 | SubmitForReview, GetReviewStatus, ApproveCartridge, RejectCartridge, RequestRevision |
| **RevenueService** | 수익 배분·정산·대시보드 | GetSalesReport, GetPayoutHistory, ConfigureRevenueSplit, RequestPayout |
| **CartridgeAnalyticsService** | 사용 통계·평점·피드백 | GetUsageStats, GetRatings, SubmitReview, GetDeveloperAnalytics |

---

## 3. 개발자 SDK 상세

### 3.1 SDK 구성 (ManPaSik Cartridge Development Kit — CDK)

```text
manpasik-cdk/
├── docs/                        # 문서
│   ├── getting-started.md       # 빠른 시작 가이드
│   ├── hardware-spec.md         # 리더기 하드웨어 인터페이스 명세
│   ├── electrical-spec.md       # 전극 배열·ADC·주파수 스펙
│   ├── nfc-protocol.md          # NFC 태그 데이터 포맷
│   ├── calibration-guide.md     # 보정 프로세스 가이드
│   ├── ai-model-guide.md        # AI 모델 패키징 가이드
│   ├── review-guidelines.md     # 심사 가이드라인
│   ├── revenue-model.md         # 수익 배분 정책
│   └── regulatory-guide.md      # 규제 인증 안내
│
├── tools/                       # 개발 도구
│   ├── cdk-cli/                 # CLI 도구 (Rust)
│   │   ├── init                 # 새 카트리지 프로젝트 생성
│   │   ├── validate             # 스펙 검증
│   │   ├── simulate             # 시뮬레이터로 테스트
│   │   ├── calibrate            # 보정 계수 산출
│   │   ├── package              # 패키지 빌드
│   │   ├── submit               # 심사 제출
│   │   └── publish              # 스토어 게시
│   │
│   ├── simulator/               # 소프트웨어 시뮬레이터
│   │   ├── virtual-reader/      # 가상 리더기 (실제 하드웨어 없이 개발)
│   │   ├── signal-generator/    # 합성 EIS 신호 생성
│   │   └── test-harness/        # 자동화 테스트 프레임워크
│   │
│   └── hardware-dev-kit/        # 하드웨어 개발 키트 (물리)
│       ├── reference-board/     # 레퍼런스 전극 보드 설계
│       └── nfc-programmer/      # NFC 태그 프로그래머
│
├── sdk/                         # 소프트웨어 라이브러리
│   ├── rust/                    # Rust SDK (핵심)
│   │   ├── cartridge-spec/      # 카트리지 스펙 정의 라이브러리
│   │   ├── calibration/         # 보정 알고리즘 라이브러리
│   │   ├── signal-processing/   # DSP 라이브러리 (manpasik-engine 서브셋)
│   │   └── ai-model/            # AI 모델 패키징 라이브러리
│   │
│   ├── python/                  # Python SDK (연구·프로토타이핑)
│   │   ├── manpasik_cdk/        # 메인 패키지
│   │   ├── notebooks/           # Jupyter 예제 노트북
│   │   └── ml-tools/            # ML 모델 훈련·변환 도구
│   │
│   └── web-api/                 # REST API 클라이언트
│       ├── openapi.yaml         # OpenAPI 명세
│       └── client-libs/         # TypeScript/Python/Go 클라이언트
│
├── templates/                   # 카트리지 프로젝트 템플릿
│   ├── basic-biomarker/         # 기본 바이오마커 (88차원)
│   ├── e-nose/                  # 전자코 (448차원)
│   ├── e-tongue/                # 전자혀 (448차원)
│   ├── multi-panel/             # 다중 바이오마커 패널 (896차원)
│   └── custom/                  # 빈 템플릿
│
└── examples/                    # 예제 카트리지
    ├── glucose-basic/           # 혈당 기본 측정
    ├── water-quality/           # 수질 검사
    ├── food-freshness/          # 식품 신선도
    └── air-voc/                 # 공기 VOC 측정
```

### 3.2 카트리지 스펙 파일 (Cartridge Specification)

모든 카트리지는 `cartridge.toml` 스펙 파일로 정의:

```toml
[cartridge]
name = "glucose-basic"
display_name = "기본 혈당 측정 카트리지"
version = "1.0.0"
developer_id = "dev_abc123"
category = "health.biomarker.glucose"
description = "모세혈관 혈액에서 혈당 수치를 빠르고 정확하게 측정합니다."
license = "proprietary"

[hardware]
required_channels = 88
electrode_config = "standard-8"        # 표준 8전극 배열
frequency_range = { min_hz = 10, max_hz = 10_000_000 }
frequency_points = 11                  # 11개 주파수 포인트
measurement_duration_sec = 15
sample_type = "capillary_blood"        # 검체 타입
storage_temp_c = { min = 2, max = 30 }
shelf_life_months = 24

[nfc]
tag_format = "manpasik-v2"
code_category = 0x01                   # 건강·바이오마커
code_type_index = 0x01                 # 혈당
extended_code = false                  # 2-byte 기본 코드

[calibration]
alpha = 0.95                           # 차동 보정 계수
temp_coefficient = 0.002               # 온도 보정
channel_offsets = [0.0, 0.0, ...]      # 88개 오프셋
channel_gains = [1.0, 1.0, ...]        # 88개 게인
calibration_model = "models/calibration_v1.tflite"

[ai]
primary_model = "models/glucose_predictor_v1.tflite"
model_type = "value_predictor"         # 값 예측 (연속값)
input_dim = 88
output_dim = 1
output_unit = "mg/dL"
output_range = { min = 20, max = 600 }
confidence_threshold = 0.8
quality_model = "models/quality_gate_v1.tflite"

[ai.classification]
enabled = true
model = "models/glucose_classifier_v1.tflite"
classes = ["정상", "경계", "주의", "위험"]
thresholds = [100, 126, 200]           # mg/dL 기준

[display]
result_template = "glucose_result"     # 결과 화면 템플릿
primary_value_label = "혈당"
primary_value_unit = "mg/dL"
icon = "assets/icon_glucose.svg"
color_scheme = "health_blue"

[pricing]
model = "per_unit"                     # per_unit | subscription_addon | free
unit_price_krw = 3000                  # 카트리지 단가
subscription_included = ["premium", "clinical"]  # 이 구독에 포함

[regulatory]
classification = "class_2_ivd"         # 체외진단 의료기기 등급
certifications = ["KFDA"]             # 보유 인증
clinical_validation = "studies/glucose_clinical_v1.pdf"
intended_use = "자가 혈당 모니터링 보조 기기용 카트리지"
```

### 3.3 개발 워크플로

```text
┌──────────────────────────────────────────────────────────────────┐
│                    카트리지 개발 워크플로                          │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1️⃣ 개발자 등록                                                  │
│  └─ developer.manpasik.com 가입 → 개인/기업/연구기관 유형 선택     │
│     → 이메일 인증 → API 키 발급 → NDA 동의 (하드웨어 스펙 접근)    │
│                                                                  │
│  2️⃣ 프로젝트 생성                                                │
│  └─ cdk init --template basic-biomarker --name my-cartridge      │
│     → cartridge.toml 스켈레톤 + 디렉토리 구조 자동 생성            │
│                                                                  │
│  3️⃣ 카트리지 설계                                                │
│  ├─ 하드웨어: 전극 패턴·시약 선택·NFC 태그 설계                    │
│  ├─ cartridge.toml: 채널 수·주파수·측정 시간·검체 타입 정의         │
│  └─ (선택) 물리 프로토타입: HDK(Hardware Dev Kit) 활용              │
│                                                                  │
│  4️⃣ 보정 (Calibration)                                           │
│  ├─ 표준 용액 / 기준 물질로 N회 반복 측정                          │
│  ├─ cdk calibrate --data calibration_data.csv                    │
│  │   → α, offset, gain, temp_coeff 자동 산출                     │
│  └─ ML 보정 모델 훈련 (Python SDK notebooks)                      │
│                                                                  │
│  5️⃣ AI 모델 개발 (선택)                                          │
│  ├─ Python SDK + Jupyter 노트북으로 모델 개발                     │
│  ├─ TFLite/ONNX로 변환                                           │
│  ├─ cdk validate --model models/my_model.tflite                  │
│  │   → 입력/출력 차원 검증, 추론 시간 벤치마크                      │
│  └─ 결과 템플릿 정의 (display 섹션)                                │
│                                                                  │
│  6️⃣ 시뮬레이션 테스트                                             │
│  ├─ cdk simulate --signals synthetic_glucose.csv                 │
│  │   → 가상 리더기에서 전체 파이프라인 실행                         │
│  ├─ cdk test --suite regression                                  │
│  │   → 자동화 회귀 테스트                                         │
│  └─ 시뮬레이터 UI에서 결과 화면 미리보기                            │
│                                                                  │
│  7️⃣ 패키지 & 제출                                                │
│  ├─ cdk package → .mpk 파일 생성 (스펙+모델+에셋 번들)            │
│  ├─ cdk submit --package my-cartridge-1.0.0.mpk                  │
│  │   → 자동 검증 (스펙 포맷, 모델 호환성, 보안 스캔)               │
│  └─ ReviewService 심사 큐에 등록                                  │
│                                                                  │
│  8️⃣ 심사 (Review)                                                │
│  ├─ 자동 심사: 스펙 유효성, 모델 안전성, NFC 코드 충돌 검사         │
│  ├─ 기술 심사: 전극 호환성, DSP 파이프라인 정합성, AI 정확도 벤치   │
│  ├─ 안전 심사: 검체 타입 위험도, 사용자 안내 적절성                  │
│  └─ (의료용) 규제 서류 확인: 인증 증서, 임상 데이터                  │
│                                                                  │
│  9️⃣ 게시 & 배포                                                  │
│  ├─ 심사 승인 → 스토어 리스팅 활성화                               │
│  ├─ 사용자 앱에 카트리지 스토어 노출                                │
│  └─ NFC 태그 프로비저닝 키 발급 (양산 시)                           │
│                                                                  │
│  🔟 판매 & 정산                                                   │
│  ├─ 사용자 구매 → 물리 카트리지 배송 또는 구독 포함                 │
│  ├─ 수익 배분: 개발자 70% / 만파식 30% (기본)                      │
│  │   소규모 개발자: 개발자 85% / 만파식 15% (연 매출 1억 미만)       │
│  └─ 월 정산 → 개발자 계좌 입금                                     │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

---

## 4. 카트리지 스토어 (사용자 측) 상세

### 4.1 스토어 UI 구성

```text
카트리지 스토어 (앱 내 탭)
├── 🏠 홈
│   ├── 배너: 추천·신규·인기 카트리지
│   ├── "내 구독에 포함된 카트리지" 섹션
│   ├── 카테고리 가로 스크롤 (건강·환경·식품·산업·연구·기타)
│   └── 개인 추천 ("최근 측정 기반 추천")
│
├── 🔍 검색 & 탐색
│   ├── 검색바 (이름·키워드·바이오마커)
│   ├── 필터: 카테고리·가격·평점·차원·검체타입·인증·구독포함
│   ├── 정렬: 인기순·최신순·평점순·가격순
│   └── 태그 클라우드 (혈당, 수질, VOC, 콜레스테롤, ...)
│
├── 📋 카트리지 상세 페이지
│   ├── 상단: 아이콘·이름·개발자·평점·다운로드 수
│   ├── 설명: 용도·측정 대상·검체 타입·측정 시간
│   ├── 스펙: 차원·주파수·정확도·인증
│   ├── 스크린샷: 결과 화면 미리보기
│   ├── 리뷰: 사용자 후기·평점
│   ├── 가격: 단품 구매 / 구독 포함 여부
│   ├── 호환성: 지원 리더기 모델
│   └── 구매/장바구니 버튼
│
├── 📦 내 카트리지
│   ├── 등록된 카트리지 목록 (NFC 스캔 이력)
│   ├── 잔여 사용 횟수
│   ├── 구매 이력
│   └── 구독 포함 카트리지
│
└── 🏷️ 카테고리
    ├── 건강·바이오마커 (혈당, 콜레스테롤, 요산, CRP, ...)
    ├── 호르몬·내분비 (코르티솔, 갑상선, ...)
    ├── 감염·면역 (CRP, PCT, IgG/IgM, ...)
    ├── 전자코·가스 (VOC, 구취, 환경가스, ...)
    ├── 전자혀·수질 (pH, 이온, 중금속, ...)
    ├── 식품·음료 (신선도, 첨가물, 알레르겐, ...)
    ├── 환경 (대기질, 토양, ...)
    ├── 산업·품질관리 (원료 순도, 제품 검사, ...)
    ├── 연구·교육 (범용 전극, 실험용, ...)
    └── 기타·특수 (커스텀, 비표적, ...)
```

### 4.2 카트리지 구매·사용 흐름

```text
구매 흐름:
1. 스토어에서 카트리지 선택
2. "구매" 또는 "장바구니 추가"
3. 결제 (PaymentService → Toss PG)
4. 물리 카트리지 배송 (ShopService 주문·배송 연동)
   OR 구독 포함 카트리지 → 자동 활성화

사용 흐름:
1. 물리 카트리지를 리더기에 삽입
2. NFC 태그 자동 인식 → CartridgeRegistryService.GetCartridgeSpec
3. 카트리지 스펙(보정계수·AI모델·차원·측정시간) 자동 다운로드
4. 측정 파이프라인 자동 설정 (MEASUREMENT-ANALYSIS-AI-SPEC 참조)
5. 측정 시작 → 결과 표시 (카트리지 정의 결과 템플릿)
6. 잔여 사용 횟수 차감 → 소진 시 재구매 안내
```

### 4.3 카트리지 자동 인식 상세

```text
NFC 태그 데이터 포맷 (ManPaSik Cartridge Tag v2):
┌──────────────────────────────────────────────────┐
│ Byte 0-1: Magic (0xMP)                           │
│ Byte 2:   Version (0x02)                          │
│ Byte 3:   Category Code (8-bit)                   │
│ Byte 4:   Type Index (8-bit)                      │
│ Byte 5-8: Extended Code (4-byte, Phase 4+)        │
│ Byte 9-24: Cartridge UID (128-bit UUID)           │
│ Byte 25-26: Manufacturing Date (16-bit)           │
│ Byte 27-28: Expiry Date (16-bit)                  │
│ Byte 29: Max Uses (uint8, 0=unlimited)            │
│ Byte 30: Current Uses (uint8)                     │
│ Byte 31: Lot Number (uint8)                       │
│ Byte 32-47: Calibration Hash (128-bit)            │
│ Byte 48-63: Developer Signature (128-bit HMAC)    │
└──────────────────────────────────────────────────┘

인식 절차:
1. NFC Read → Magic 검증 (0xMP)
2. Version 확인 → 호환 프로토콜 선택
3. Category + TypeIndex → 로컬 캐시에서 스펙 조회
   ├─ 캐시 히트 → 즉시 로드
   └─ 캐시 미스 → CartridgeRegistryService.GetCartridgeSpec(code) 다운로드
4. Developer Signature 검증 (HMAC, 위변조 방지)
5. Expiry Date 확인 → 유효기간 초과 시 경고
6. Current Uses < Max Uses 확인 → 소진 시 "새 카트리지 필요"
7. Calibration Hash로 보정 데이터 무결성 검증
8. 스펙 로드 완료 → 측정 준비 ("이 카트리지로 측정하시겠습니까?")
```

---

## 5. 개발자 콘솔 (Developer Console)

### 5.1 기능 구성

```text
developer.manpasik.com
├── 대시보드
│   ├── 판매 요약 (일/주/월 매출, 다운로드 수)
│   ├── 활성 카트리지 수
│   ├── 평균 평점
│   └── 심사 상태 알림
│
├── 카트리지 관리
│   ├── 등록된 카트리지 목록
│   ├── 새 카트리지 등록 (웹 UI 또는 cdk-cli)
│   ├── 버전 관리 (업데이트 제출)
│   ├── 심사 상태 추적
│   └── 게시/비게시 토글
│
├── 분석 (Analytics)
│   ├── 사용 통계 (측정 횟수, 성공률, 평균 측정 시간)
│   ├── 사용자 분포 (지역, 구독 등급)
│   ├── AI 모델 성능 (신뢰도 분포, 오분류율)
│   ├── 품질 게이트 통과율
│   └── 사용자 피드백·리뷰
│
├── 수익 (Revenue)
│   ├── 매출 상세 (카트리지별, 기간별)
│   ├── 수익 배분 내역
│   ├── 정산 이력
│   ├── 세금 서류
│   └── 출금 요청
│
├── SDK & 도구
│   ├── API 키 관리
│   ├── SDK 다운로드
│   ├── 문서 (하드웨어 스펙, 가이드)
│   ├── 시뮬레이터 웹 버전
│   └── NFC 태그 프로비저닝 키 발급
│
├── 팀 관리
│   ├── 팀원 초대·역할 설정 (관리자/개발자/분석가)
│   └── 권한 관리
│
└── 설정
    ├── 개발자 프로필
    ├── 알림 설정
    ├── 결제 정보
    └── NDA·계약 관리
```

### 5.2 개발자 등급

| 등급 | 요건 | 혜택 |
| --- | --- | --- |
| **Explorer** (무료) | 이메일 인증 | SDK 다운로드, 시뮬레이터, 문서 접근. 스토어 등록 불가 |
| **Developer** (연 $99) | 개인 신원 확인 | 카트리지 3종 등록, 스토어 게시, 수익 배분 70/30 |
| **Professional** (연 $299) | 사업자 확인 | 카트리지 무제한 등록, 우선 심사, 수익 배분 75/25 |
| **Enterprise** (맞춤) | 기업 계약 | 맞춤 수익 배분, 전담 지원, 대량 NFC 키, API 우선 할당량 |
| **Research** (무료) | 교육·연구기관 | Developer 동등 + 논문용 데이터 접근, 비상업 조건 |

---

## 6. 심사 프로세스 (Review Process)

### 6.1 심사 단계

```text
제출 (.mpk 패키지)
    ↓
[1단계] 자동 검증 (Automated)          — 즉시
├─ cartridge.toml 스키마 유효성
├─ NFC 코드 고유성 (기존 코드 충돌 검사)
├─ AI 모델 포맷 검증 (TFLite/ONNX 로드 가능)
├─ 모델 입력/출력 차원 일관성
├─ 보안 스캔 (악성 코드, 데이터 유출 시도)
├─ 에셋 크기 제한 (모델 < 50MB, 총 < 100MB)
└─ 결과: Pass → 2단계 / Fail → 즉시 반려 + 사유

    ↓
[2단계] 기술 심사 (Technical)          — 1~3 영업일
├─ 전극 배열 호환성 (리더기 하드웨어 매칭)
├─ DSP 파이프라인 정합성 (채널 수, 주파수 범위)
├─ 보정 데이터 품질 (R² ≥ 0.95, CV < 10%)
├─ AI 모델 성능 벤치마크 (표준 데이터셋 대비)
├─ 시뮬레이터 통과 테스트
└─ 결과: Pass → 3단계 / Revision → 수정 요청

    ↓
[3단계] 안전·콘텐츠 심사 (Safety)      — 1~2 영업일
├─ 검체 타입 안전성 (혈액·타액·소변·환경 등)
├─ 사용자 안내 문구 적절성 ("진단 보조" 명시)
├─ 결과 해석 가이드 정확성
├─ 연령 제한·경고 문구
└─ 결과: Pass → 게시 / Revision → 수정 요청

    ↓
[의료용 추가] 규제 심사 (Regulatory)    — 별도
├─ KFDA/CE/FDA 인증서 확인
├─ 임상 검증 데이터 검토
├─ GMP 제조 시설 확인
└─ 결과: Pass → 의료용 배지 부여 / Fail → 일반용만 허용

    ↓
승인 → 스토어 게시 🎉
```

### 6.2 심사 기준 상세

| 영역 | 기준 | 합격 조건 |
| --- | --- | --- |
| 보정 품질 | 결정계수 (R²) | ≥ 0.95 |
| 보정 품질 | 변동계수 (CV) | < 10% |
| AI 정확도 | 분류: F1-Score | ≥ 0.85 (표준 데이터셋) |
| AI 정확도 | 회귀: MAPE | < 15% |
| 추론 성능 | 온디바이스 지연 | < 100ms (88-dim), < 500ms (448-dim) |
| 추론 성능 | 서버 지연 | < 1000ms (896/1792-dim) |
| 모델 크기 | TFLite/ONNX 파일 | < 50MB |
| 보안 | 코드 스캔 | 악성 코드 0건 |
| NFC | 코드 고유성 | 기존 등록 코드 비충돌 |

---

## 7. 수익 모델 및 정산

### 7.1 수익 배분 구조

```text
카트리지 판매 수익 흐름:
─────────────────────────────────────────
사용자 결제 (100%)
    │
    ├─ PG 수수료 (약 3%)
    │
    └─ 순 수익 (97%)
        ├─ 개발자 몫 (70% 기본, 소규모 85%)
        └─ 만파식 몫 (30% 기본, 소규모 15%)
            ├─ 플랫폼 운영
            ├─ 리더기 하드웨어 유지
            └─ 인프라·심사·지원

구독 포함 카트리지:
─────────────────────────────────────────
구독료 중 카트리지 비중 배분
├─ 구독 등급별 "카트리지 풀" 예산 책정
├─ 사용 횟수 비례 배분 (측정 1회당 개발자에게 X원)
└─ 월별 정산

소규모 개발자 우대:
─────────────────────────────────────────
연 매출 1억원 미만 → 개발자 85% / 만파식 15%
연 매출 1억원 이상 → 개발자 70% / 만파식 30%
(Apple Small Business Program 모델)
```

### 7.2 가격 정책 옵션

| 모델 | 설명 | 예시 |
| --- | --- | --- |
| **Per-Unit** | 물리 카트리지 개당 판매 | 혈당 카트리지 3,000원/개 |
| **Subscription-Included** | 특정 구독에 포함 (무제한 또는 월 N회) | Premium 구독에 혈당·콜레스테롤 포함 |
| **Subscription-AddOn** | 구독자 추가 결제 옵션 | 호르몬 패널 +5,000원/월 |
| **Free** | 무료 배포 (연구·교육·프로모션) | 대학 교육용 범용 전극 |
| **B2B** | 기업 대량 구매 | 공장 품질관리용 카트리지 |

---

## 8. 기술 인프라

### 8.1 신규 백엔드 서비스

| 서비스 | 역할 | 기술 |
| --- | --- | --- |
| `CartridgeRegistryService` | 카트리지 타입 등록·스펙 관리·버전 | gRPC, PostgreSQL `cartridge_types`, `cartridge_versions` |
| `DeveloperService` | 개발자 계정·인증·API키·팀 | gRPC, PostgreSQL `developers`, `api_keys`, `teams` |
| `StoreService` | 스토어 리스팅·검색·구매 | gRPC, Elasticsearch(검색), PostgreSQL `store_items`, `purchases` |
| `ReviewService` | 심사 워크플로·자동 검증 | gRPC, PostgreSQL `review_submissions`, `review_results` |
| `RevenueService` | 수익 계산·정산·대시보드 | gRPC, PostgreSQL `revenue_transactions`, `payouts` |
| `CartridgeAnalyticsService` | 사용 통계·평점·피드백 | gRPC, TimescaleDB `cartridge_usage_stats`, PostgreSQL `cartridge_reviews` |

### 8.2 데이터베이스 스키마 (주요 테이블)

```sql
-- 카트리지 타입 레지스트리
CREATE TABLE cartridge_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID REFERENCES developers(id),
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200) NOT NULL,
    category_code SMALLINT NOT NULL,      -- 8-bit
    type_index SMALLINT NOT NULL,         -- 8-bit
    extended_code INTEGER,                -- 4-byte (Phase 4+)
    version VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',   -- draft/in_review/approved/rejected/deprecated
    spec_json JSONB NOT NULL,             -- cartridge.toml → JSON
    model_bundle_url TEXT,                -- MinIO URL (.mpk)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(category_code, type_index, version)
);

-- 개발자
CREATE TABLE developers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    developer_type VARCHAR(20) NOT NULL,  -- individual/company/research
    company_name VARCHAR(200),
    tier VARCHAR(20) DEFAULT 'explorer',  -- explorer/developer/professional/enterprise/research
    api_key_hash VARCHAR(64),
    revenue_share_pct NUMERIC(5,2) DEFAULT 70.00,
    payout_account JSONB,                 -- 정산 계좌 정보 (암호화)
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 스토어 아이템
CREATE TABLE store_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartridge_type_id UUID REFERENCES cartridge_types(id),
    price_krw INTEGER,
    price_model VARCHAR(20) NOT NULL,     -- per_unit/subscription_included/addon/free/b2b
    subscription_tiers TEXT[],            -- 포함 구독 등급
    featured BOOLEAN DEFAULT FALSE,
    download_count INTEGER DEFAULT 0,
    avg_rating NUMERIC(3,2) DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 심사
CREATE TABLE review_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartridge_type_id UUID REFERENCES cartridge_types(id),
    developer_id UUID REFERENCES developers(id),
    package_url TEXT NOT NULL,            -- .mpk MinIO URL
    status VARCHAR(20) DEFAULT 'pending', -- pending/auto_check/tech_review/safety_review/regulatory/approved/rejected
    auto_check_result JSONB,
    tech_review_result JSONB,
    safety_review_result JSONB,
    reviewer_notes TEXT,
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

-- 수익 거래
CREATE TABLE revenue_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    developer_id UUID REFERENCES developers(id),
    cartridge_type_id UUID REFERENCES cartridge_types(id),
    transaction_type VARCHAR(20) NOT NULL, -- sale/subscription_usage/refund
    gross_amount INTEGER NOT NULL,         -- 총액 (원)
    platform_fee INTEGER NOT NULL,         -- 플랫폼 수수료
    developer_amount INTEGER NOT NULL,     -- 개발자 정산액
    period_month VARCHAR(7),               -- YYYY-MM
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 카트리지 사용 통계
CREATE TABLE cartridge_usage_stats (
    time TIMESTAMPTZ NOT NULL,
    cartridge_type_id UUID NOT NULL,
    measurement_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    avg_confidence NUMERIC(5,4),
    avg_measurement_sec NUMERIC(6,2),
    quality_pass_rate NUMERIC(5,4)
);
SELECT create_hypertable('cartridge_usage_stats', 'time');

-- 카트리지 리뷰
CREATE TABLE cartridge_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cartridge_type_id UUID REFERENCES cartridge_types(id),
    user_id UUID REFERENCES users(id),
    rating SMALLINT CHECK (rating BETWEEN 1 AND 5),
    title VARCHAR(200),
    body TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 8.3 프론트엔드 라우트

| 경로 | 화면 | 설명 |
| --- | --- | --- |
| `/store` | CartridgeStoreScreen | 스토어 홈 (배너·카테고리·추천) |
| `/store/search` | StoreSearchScreen | 검색·필터·정렬 |
| `/store/category/:id` | StoreCategoryScreen | 카테고리별 리스트 |
| `/store/item/:id` | StoreItemDetailScreen | 카트리지 상세 |
| `/store/my-cartridges` | MyCartridgesScreen | 내 카트리지 관리 |
| `/store/purchase-history` | PurchaseHistoryScreen | 구매 이력 |

---

## 9. 보안 및 위변조 방지

### 9.1 카트리지 인증 체인

```text
만파식 루트 CA (HSM 보관)
    │
    ├─ 개발자 서명 키 (개발자 등록 시 발급)
    │   └─ 카트리지 HMAC 서명 (NFC 태그 Byte 48-63)
    │
    └─ NFC 프로비저닝 키 (양산 시 발급)
        └─ 태그 쓰기 권한 제어

검증 흐름:
1. NFC 태그 읽기
2. Developer Signature (HMAC-SHA256) 검증
   → 키: 개발자 공개키 (CartridgeRegistryService에서 조회)
   → 데이터: Byte 0-47
3. Calibration Hash 검증
   → 서버 스펙의 calibration_coefficients SHA-256과 비교
4. 모두 통과 → 정품 인증 완료
   하나라도 실패 → "인증되지 않은 카트리지" 경고
```

### 9.2 보안 정책

| 위협 | 대응 |
| --- | --- |
| 위조 카트리지 | NFC HMAC 서명 검증, 서버 크로스체크 |
| NFC 태그 복제 | 카트리지 UID 고유성 + 사용 횟수 서버 동기화 |
| 악성 AI 모델 | 심사 시 샌드박스 실행, 시스템 콜 제한, 모델 서명 |
| 데이터 유출 | AI 모델 입출력 감사, 네트워크 접근 차단 (온디바이스) |
| 가격 조작 | 서버 사이드 가격 검증, PG 결제 금액 교차 확인 |

---

## 10. 구현 Phase별 로드맵

| Phase | 범위 | 핵심 산출물 |
| --- | --- | --- |
| Phase 1 (MVP) | 1st-party 카트리지만. 내부 등록·스토어 기본 UI | CartridgeRegistryService, StoreService(기본), 스토어 홈·상세 |
| Phase 2 | 개발자 등록·SDK 공개 (Explorer/Developer). 시뮬레이터. 기본 심사 | DeveloperService, ReviewService(자동+기술), cdk-cli, Python SDK |
| Phase 3 | 스토어 전면 오픈. 수익 배분·정산. 고급 분석 | RevenueService, CartridgeAnalyticsService, 개발자 콘솔 전체 |
| Phase 4 | Enterprise/Research 등급. B2B. 4-byte 확장 코드. 글로벌 확장 | 대량 NFC 키, 국제 결제, 다국어 스토어 |
| Phase 5 | AI 마켓플레이스 (카트리지 + AI 모델 별도 거래). 커뮤니티 기여 | AI 모델 스토어, 오픈소스 기여 프로그램 |

---

## 11. 참조

- COMPREHENSIVE-SYSTEM-BLUEPRINT-v3.0.md: Part 10.12 (카트리지 시스템), Part 3.5 (마켓·결제)
- MEASUREMENT-ANALYSIS-AI-SPEC.md: 측정·분석 파이프라인, 카트리지-AI 연계
- FINAL-MASTER-IMPLEMENTATION-PLAN.md: I.2 추적성, II.4 마켓 세부
- Apple Developer — App Store Review Guidelines, Accessory Design Guidelines
- Biocartis — Content Partnering Model
- Fluxergy — OEM Assay Platform
- Fathym OpenBiotech — Cloud-Native Biosensing Platform
- OpenBCI — Open-Source Hardware SDK
