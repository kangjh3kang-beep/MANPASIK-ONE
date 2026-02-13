# 카트리지 무한확장 체계 및 등급별 접근 제어 명세서

**문서번호**: MPK-SPEC-CART-v1.0-20260211  
**기반**: MPK-ECO-PLAN-v1.1-COMPLETE §V.5.9, §V.5.14  
**목적**: 카트리지를 무한 확장 가능한 레지스트리 구조로 재설계하고, 구독 등급별 분석 카트리지 접근 범위를 정의한다.

---

## 1. 설계 원칙

1. **무한확장(Open-Ended)**: 카트리지 종류에 상한이 없다. 신규 카테고리·타입은 코드 배포 없이 서버 레지스트리(DB) 등록만으로 추가한다.
2. **하위호환**: 기존 29종(v1.0 레거시) 카트리지는 코드 변경 없이 동작한다.
3. **등급 연동**: 사용자의 구독 티어(Free/Basic/Pro/Clinical)에 따라 사용 가능한 카트리지 범위가 결정된다.
4. **동적 정책**: 카트리지-등급 접근 정책은 DB 기반으로 운영되어, 관리자가 코드 배포 없이 정책을 변경할 수 있다.
5. **SDK/서드파티**: 외부 개발자가 카트리지를 설계·제조하여 마켓플레이스에 등록할 수 있는 확장 경로를 제공한다.

---

## 2. 카트리지 코드 체계 (Cartridge Code Architecture)

### 2.1 2-Byte 계층형 코드

```
┌────────────────────────────────────────────────┐
│   Cartridge Full Code (2 bytes = 16 bits)      │
│   ┌──────────────┬───────────────────────────┐ │
│   │ Category (8b)│ TypeIndex (8b)            │ │
│   │ 0x00 ~ 0xFF  │ 0x00 ~ 0xFF              │ │
│   └──────────────┴───────────────────────────┘ │
│   → 256 카테고리 × 256 타입/카테고리            │
│   → 총 65,536 종 수용                          │
└────────────────────────────────────────────────┘
```

### 2.2 확장 코드 (4-Byte, Phase 4+)

```
┌────────────────────────────────────────────────────────┐
│   Extended Code (4 bytes = 32 bits) — 미래 확장용       │
│   ┌──────────────┬──────────────┬────────────────────┐ │
│   │ Category(16b)│ TypeIndex(16b)│                    │ │
│   │ 0x0000~0xFFFF│ 0x0000~0xFFFF│                    │ │
│   └──────────────┴──────────────┴────────────────────┘ │
│   → 65,536 카테고리 × 65,536 타입/카테고리              │
│   → 총 4,294,967,296 (약 43억) 종 수용                  │
└────────────────────────────────────────────────────────┘
```

### 2.3 카테고리 코드 할당표

| 범위 | 카테고리 | 설명 | Phase | 현재 타입 수 |
|------|---------|------|-------|-------------|
| `0x01` | **HealthBiomarker** | 건강 바이오마커 (혈액/타액/체액) | 1 | 14종 (확장 가능) |
| `0x02` | **Environmental** | 환경 모니터링 (수질/공기/방사능) | 1 | 4종 (확장 가능) |
| `0x03` | **FoodSafety** | 식품 안전 (농약/신선도/알레르겐) | 1 | 4종 (확장 가능) |
| `0x04` | **ElectronicSensor** | 전자코/전자혀/EHD | 1 | 3종 (확장 가능) |
| `0x05` | **AdvancedAnalysis** | 고급 분석 (비표적/다중패널) | 2 | 3종 (확장 가능) |
| `0x06` | **Industrial** | 산업용 분석 (화학물질/중금속/유해가스) | 3 | 예비 |
| `0x07` | **Veterinary** | 수의학 (동물 혈액/바이오마커) | 3 | 예비 |
| `0x08` | **Pharmaceutical** | 제약 (약물 성분/농도 분석) | 3 | 예비 |
| `0x09` | **Agricultural** | 농업 (토양/비료/작물 분석) | 4 | 예비 |
| `0x0A` | **Cosmetic** | 화장품 (성분/피부 타입) | 4 | 예비 |
| `0x0B` | **Forensic** | 법의학 (체액/약물/독물) | 4 | 예비 |
| `0x0C` | **Marine** | 해양 (해수/양식장/선박 연료) | 4 | 예비 |
| `0x0D`~`0xEF` | **Reserved** | 미래 확장 예비 | — | — |
| `0xF0`~`0xFD` | **ThirdParty** | SDK/서드파티 마켓플레이스 | 4 | 동적 할당 |
| `0xFE` | **Beta** | 베타/실험용 (Clinical 전용) | 2 | 동적 |
| `0xFF` | **CustomResearch** | 맞춤형 연구용 | 1 | 1종 |

### 2.4 레거시 호환 매핑 (v1.0 → v2.0)

| 레거시 코드 (1byte) | 카테고리 | 신규 코드 (2byte) | 타입명 |
|---------------------|---------|------------------|--------|
| `0x01` | 0x01 HealthBiomarker | `0x01:0x01` | Glucose |
| `0x02` | 0x01 HealthBiomarker | `0x01:0x02` | LipidPanel |
| `0x03` | 0x01 HealthBiomarker | `0x01:0x03` | HbA1c |
| `0x04` | 0x01 HealthBiomarker | `0x01:0x04` | UricAcid |
| `0x05` | 0x01 HealthBiomarker | `0x01:0x05` | Creatinine |
| `0x06` | 0x01 HealthBiomarker | `0x01:0x06` | VitaminD |
| `0x07` | 0x01 HealthBiomarker | `0x01:0x07` | VitaminB12 |
| `0x08` | 0x01 HealthBiomarker | `0x01:0x08` | Ferritin |
| `0x09` | 0x01 HealthBiomarker | `0x01:0x09` | Tsh |
| `0x0A` | 0x01 HealthBiomarker | `0x01:0x0A` | Cortisol |
| `0x0B` | 0x01 HealthBiomarker | `0x01:0x0B` | Testosterone |
| `0x0C` | 0x01 HealthBiomarker | `0x01:0x0C` | Estrogen |
| `0x0D` | 0x01 HealthBiomarker | `0x01:0x0D` | Crp |
| `0x0E` | 0x01 HealthBiomarker | `0x01:0x0E` | Insulin |
| `0x20` | 0x02 Environmental | `0x02:0x01` | WaterQuality |
| `0x21` | 0x02 Environmental | `0x02:0x02` | IndoorAirQuality |
| `0x22` | 0x02 Environmental | `0x02:0x03` | Radon |
| `0x23` | 0x02 Environmental | `0x02:0x04` | Radiation |
| `0x30` | 0x03 FoodSafety | `0x03:0x01` | PesticideResidue |
| `0x31` | 0x03 FoodSafety | `0x03:0x02` | FoodFreshness |
| `0x32` | 0x03 FoodSafety | `0x03:0x03` | Allergen |
| `0x33` | 0x03 FoodSafety | `0x03:0x04` | DateDrug |
| `0x40` | 0x04 ElectronicSensor | `0x04:0x01` | ENose |
| `0x41` | 0x04 ElectronicSensor | `0x04:0x02` | ETongue |
| `0x42` | 0x04 ElectronicSensor | `0x04:0x03` | EhdGas |
| `0x50` | 0x05 AdvancedAnalysis | `0x05:0x01` | NonTarget448 |
| `0x51` | 0x05 AdvancedAnalysis | `0x05:0x02` | NonTarget896 |
| `0x52` | 0x05 AdvancedAnalysis | `0x05:0x03` | NonTarget1792 (1792차원 궁극, Phase 5) |
| `0x53` | 0x05 AdvancedAnalysis | `0x05:0x04` | MultiBiomarker |
| `0xFF` | 0xFF CustomResearch | `0xFF:0x01` | CustomResearch |

---

## 3. 구독 등급별 카트리지 접근 제어 (Tier-Based Cartridge Access Control)

### 3.1 접근 레벨 정의

| 접근 레벨 | 코드 | 설명 |
|-----------|------|------|
| **INCLUDED** | `included` | 구독에 포함, 무제한 사용 |
| **LIMITED** | `limited` | 구독에 포함, 일/월 사용 횟수 제한 |
| **ADD_ON** | `add_on` | 별도 구매 시 사용 가능 (건당 과금 또는 팩 구매) |
| **RESTRICTED** | `restricted` | 해당 등급에서 사용 불가 (상위 등급 필요) |
| **BETA** | `beta` | 베타 테스트용 (Clinical 등급만 신청 가능) |

### 3.2 기본 정책 매트릭스 (Default Tier-Cartridge Access)

| 카테고리 | Free | Basic Safety | Bio-Optimization (Pro) | Clinical Guard |
|---------|------|-------------|----------------------|----------------|
| **HealthBiomarker 기본 3종** (Glucose, LipidPanel, HbA1c) | ✅ INCLUDED (일 3회) | ✅ INCLUDED | ✅ INCLUDED | ✅ INCLUDED |
| **HealthBiomarker 나머지** (UricAcid~Insulin, 11종) | 🔒 RESTRICTED | ✅ INCLUDED | ✅ INCLUDED | ✅ INCLUDED |
| **Environmental** (수질/공기/라돈/방사능) | 🔒 RESTRICTED | 💰 ADD_ON | ✅ INCLUDED | ✅ INCLUDED |
| **FoodSafety** (농약/신선도/알레르겐/데이트약물) | 🔒 RESTRICTED | 💰 ADD_ON | ✅ INCLUDED | ✅ INCLUDED |
| **ElectronicSensor** (전자코/전자혀/EHD) | 🔒 RESTRICTED | 🔒 RESTRICTED | ✅ INCLUDED | ✅ INCLUDED |
| **AdvancedAnalysis** (비표적448/896/다중패널) | 🔒 RESTRICTED | 🔒 RESTRICTED | 💰 ADD_ON | ✅ INCLUDED |
| **Industrial** (산업용) | 🔒 RESTRICTED | 🔒 RESTRICTED | 🔒 RESTRICTED | ✅ INCLUDED |
| **Veterinary** (수의학) | 🔒 RESTRICTED | 🔒 RESTRICTED | 💰 ADD_ON | ✅ INCLUDED |
| **Pharmaceutical** (제약) | 🔒 RESTRICTED | 🔒 RESTRICTED | 🔒 RESTRICTED | ✅ INCLUDED |
| **ThirdParty** (SDK/서드파티) | 🔒 RESTRICTED | 🔒 RESTRICTED | 💰 ADD_ON | 💰 ADD_ON |
| **Beta** (베타/실험용) | 🔒 RESTRICTED | 🔒 RESTRICTED | 🔒 RESTRICTED | 🧪 BETA |
| **CustomResearch** (맞춤 연구) | 🔒 RESTRICTED | 🔒 RESTRICTED | 🔒 RESTRICTED | ✅ INCLUDED |

### 3.3 세분화 정책 (타입 레벨 오버라이드)

카테고리 단위 정책 외에, 개별 타입에 대한 오버라이드가 가능합니다.

```
정책 적용 우선순위:
  1. 타입별 오버라이드 (type-level override)
  2. 카테고리별 정책 (category-level policy)
  3. 글로벌 기본값 (global default = RESTRICTED)
```

**예시**: HealthBiomarker 카테고리의 Free 등급 기본 정책은 RESTRICTED이지만, Glucose·LipidPanel·HbA1c 3종은 타입별 오버라이드로 INCLUDED(LIMITED, 일 3회)로 설정.

### 3.4 접근 제어 데이터 모델

```
cartridge_tier_access:
  - tier:           SubscriptionTier (0~3)
  - category_code:  u8 (카테고리, 0x00 = 전체 카테고리)
  - type_index:     u8 (타입 인덱스, 0x00 = 카테고리 내 전체)
  - access_level:   enum (included, limited, add_on, restricted, beta)
  - daily_limit:    int (일일 사용 제한, 0 = 무제한, limited일 때만 적용)
  - monthly_limit:  int (월간 사용 제한, 0 = 무제한)
  - addon_price_krw: int (add_on일 때 건당/팩당 가격)
  - priority:       int (오버라이드 우선순위, 높을수록 우선)
  - is_active:      bool
  - effective_from: timestamp
  - effective_until: timestamp (null = 무기한)
```

### 3.5 접근 검증 흐름

```
측정 시작 요청 (StartSession)
  │
  ├─→ 카트리지 NFC 읽기 → CartridgeInfo (category_code, type_index)
  │
  ├─→ 사용자 구독 조회 → SubscriptionTier
  │
  ├─→ 접근 정책 조회 (우선순위: 타입별 → 카테고리별 → 기본값)
  │     │
  │     ├─ INCLUDED / LIMITED → ✅ 허용 (LIMITED는 잔여 횟수 차감)
  │     ├─ ADD_ON → 사용자 애드온 구매 여부 확인
  │     │     ├─ 구매함 → ✅ 허용 (잔여 횟수 차감)
  │     │     └─ 미구매 → ❌ 차단 + 구매 유도 UI
  │     ├─ BETA → Clinical 등급 + 베타 옵트인 확인
  │     └─ RESTRICTED → ❌ 차단 + 상위 등급 안내
  │
  └─→ 측정 진행 or 차단 응답
```

---

## 4. NFC 태그 데이터 구조 (v2.0)

### 4.1 확장 태그 레이아웃 (80+ 바이트)

```
Offset  Length  Field                Description
──────  ──────  ──────────────────   ─────────────────────────
[0-7]   8       cartridge_uid        카트리지 UID (고유 식별자)
[8]     1       category_code        카테고리 코드 (0x01~0xFF)
[9]     1       type_index           타입 인덱스 (카테고리 내 순번)
[10]    1       legacy_code          v1.0 호환 코드 (0x00이면 v2.0 전용)
[11]    1       version              태그 포맷 버전 (0x01=v1.0, 0x02=v2.0)
[12-19] 8       lot_id               제조 로트 ID (ASCII)
[20-27] 8       expiry_date          유효 기간 (YYYYMMDD)
[28-29] 2       remaining_uses       잔여 사용 횟수 (u16 LE)
[30-31] 2       max_uses             최대 사용 횟수 (u16 LE)
[32]    1       required_channels_hi 필요 채널 수 상위 바이트
[33]    1       required_channels_lo 필요 채널 수 하위 바이트
[34]    1       measurement_secs     측정 시간 (초)
[35]    1       flags                플래그 (비트: 0=인증필요, 1=보정필수, ...)
[36-43] 8       alpha_coefficient    α 계수 (f64 LE)
[44-51] 8       temp_coefficient     온도 보정 계수 (f64 LE)
[52-59] 8       humidity_coefficient 습도 보정 계수 (f64 LE)
[60-63] 4       checksum             CRC-32 체크섬
[64+]   var     extended_calibration 확장 보정 데이터 (가변)
```

### 4.2 v1.0 → v2.0 자동 변환 규칙

```
if tag.version == 0x01 (v1.0):
    category_code = legacy_category_map[tag[8]]
    type_index    = legacy_type_map[tag[8]]
    legacy_code   = tag[8]
elif tag.version == 0x02 (v2.0):
    category_code = tag[8]
    type_index    = tag[9]
    legacy_code   = tag[10]  // 0x00이면 레거시 매핑 없음
```

---

## 5. 카트리지 레지스트리 (Cartridge Registry)

### 5.1 서버 레지스트리 DB 스키마

```sql
-- 카트리지 카테고리 (무한 확장)
CREATE TABLE cartridge_categories (
    code          SMALLINT PRIMARY KEY,        -- 0x01~0xFF (카테고리 코드)
    name_en       VARCHAR(100) NOT NULL,
    name_ko       VARCHAR(100) NOT NULL,
    description   TEXT DEFAULT '',
    icon_url      VARCHAR(500) DEFAULT '',
    sort_order    INTEGER DEFAULT 0,
    is_active     BOOLEAN DEFAULT TRUE,
    phase         INTEGER DEFAULT 1,           -- 도입 Phase
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

-- 카트리지 타입 (무한 확장, 카테고리당 최대 256종)
CREATE TABLE cartridge_types (
    category_code SMALLINT NOT NULL REFERENCES cartridge_categories(code),
    type_index    SMALLINT NOT NULL,           -- 0x01~0xFF (카테고리 내 타입 순번)
    legacy_code   SMALLINT DEFAULT 0,          -- v1.0 호환 코드 (0이면 없음)
    name_en       VARCHAR(100) NOT NULL,
    name_ko       VARCHAR(100) NOT NULL,
    description   TEXT DEFAULT '',
    required_channels  INTEGER NOT NULL DEFAULT 88,
    measurement_secs   INTEGER NOT NULL DEFAULT 15,
    unit               VARCHAR(30) DEFAULT '',   -- 측정 단위 (mg/dL, ppm 등)
    reference_range    VARCHAR(100) DEFAULT '',  -- 정상 범위
    is_active     BOOLEAN DEFAULT TRUE,
    is_beta       BOOLEAN DEFAULT FALSE,
    phase         INTEGER DEFAULT 1,
    manufacturer  VARCHAR(200) DEFAULT 'ManPaSik',  -- 제조사 (서드파티 확장)
    sdk_vendor_id VARCHAR(100) DEFAULT '',          -- SDK 벤더 ID (서드파티)
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (category_code, type_index)
);

-- 구독 등급별 카트리지 접근 정책 (동적 관리)
CREATE TABLE cartridge_tier_access (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier            INTEGER NOT NULL,              -- 0: Free, 1: Basic, 2: Pro, 3: Clinical
    category_code   SMALLINT NOT NULL DEFAULT 0,   -- 0이면 전체 카테고리
    type_index      SMALLINT NOT NULL DEFAULT 0,   -- 0이면 카테고리 내 전체 타입
    access_level    VARCHAR(20) NOT NULL DEFAULT 'restricted',
                    -- included, limited, add_on, restricted, beta
    daily_limit     INTEGER DEFAULT 0,             -- 0 = 무제한
    monthly_limit   INTEGER DEFAULT 0,             -- 0 = 무제한
    addon_price_krw INTEGER DEFAULT 0,             -- add_on 시 건당 가격
    priority        INTEGER DEFAULT 0,             -- 높을수록 우선 (타입별 오버라이드 > 카테고리별)
    is_active       BOOLEAN DEFAULT TRUE,
    effective_from  TIMESTAMPTZ DEFAULT NOW(),
    effective_until TIMESTAMPTZ,                   -- NULL = 무기한
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (tier, category_code, type_index)
);

-- 사용자별 애드온 구매 내역
CREATE TABLE cartridge_addon_purchases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    category_code   SMALLINT NOT NULL,
    type_index      SMALLINT NOT NULL,
    remaining_uses  INTEGER NOT NULL DEFAULT 0,    -- 잔여 사용 횟수
    total_purchased INTEGER NOT NULL DEFAULT 0,    -- 총 구매 횟수
    price_krw       INTEGER NOT NULL DEFAULT 0,
    purchased_at    TIMESTAMPTZ DEFAULT NOW(),
    expires_at      TIMESTAMPTZ                    -- 유효 기간 (NULL = 무기한)
);

-- 카트리지 사용 로그 (감사 추적 + 사용량 추적)
CREATE TABLE cartridge_usage_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    session_id      UUID NOT NULL,
    cartridge_uid   VARCHAR(100) NOT NULL,
    category_code   SMALLINT NOT NULL,
    type_index      SMALLINT NOT NULL,
    tier_at_usage   INTEGER NOT NULL,              -- 사용 시점의 구독 등급
    access_level    VARCHAR(20) NOT NULL,          -- 적용된 접근 레벨
    used_at         TIMESTAMPTZ DEFAULT NOW()
);
```

### 5.2 Rust 엣지 레지스트리 (로컬 캐시)

오프라인 환경에서도 카트리지 검증이 가능하도록 서버 레지스트리를 로컬에 캐시합니다.

```
로컬 레지스트리 (Sled/SQLite):
  cartridge_registry:
    - 서버에서 주기적 동기화 (온라인 시)
    - 최소 기본 29종은 펌웨어에 내장 (fallback)
    - 신규 타입은 OTA 또는 동기화로 추가
  tier_access_cache:
    - 사용자 구독 등급 + 카트리지 접근 정책 캐시
    - 오프라인 시 마지막 동기화 데이터 사용
    - 정책 변경은 다음 온라인 동기화에 반영
```

---

## 6. 서드파티/SDK 카트리지 확장 경로

### 6.1 카트리지 마켓플레이스 워크플로우 (Phase 4)

```
서드파티 개발자
  │
  ├─ 1. manpasik-sdk로 카트리지 프로토콜 개발
  ├─ 2. 카트리지 스펙 제출 (카테고리, 채널, 보정식, NFC 레이아웃)
  ├─ 3. 만파식 검증팀 리뷰 (정확도, 안전성, 규제)
  ├─ 4. 승인 → cartridge_types 레지스트리 등록
  │     - category_code = 0xF0~0xFD (ThirdParty 범위)
  │     - manufacturer = 서드파티 벤더명
  │     - sdk_vendor_id = 벤더 ID
  ├─ 5. 마켓플레이스 게시
  └─ 6. 수익 분배 (만파식:벤더 = 30:70 기본)
```

### 6.2 동적 카테고리 할당

서드파티가 기존 카테고리에 맞지 않는 완전히 새로운 영역을 제안하면:

1. `0xF0~0xFD` 범위에서 새 카테고리 코드 할당
2. 카테고리가 충분히 성숙하면 `0x0D~0xEF` Reserved 범위로 공식 승격
3. 승격 시 레거시 매핑 유지 (이전 코드도 계속 동작)

---

## 7. 확장 예시 (로드맵)

### Phase 1 (현재): 29종 기본 카트리지

| 카테고리 | 타입 수 | 등록 상태 |
|---------|---------|----------|
| HealthBiomarker (0x01) | 14 | ✅ 등록 |
| Environmental (0x02) | 4 | ✅ 등록 |
| FoodSafety (0x03) | 4 | ✅ 등록 |
| ElectronicSensor (0x04) | 3 | ✅ 등록 |
| AdvancedAnalysis (0x05) | 3 | ✅ 등록 |
| CustomResearch (0xFF) | 1 | ✅ 등록 |
| **합계** | **29** | |

### Phase 2 확장 예시 (+15종)

| 카테고리 | 신규 타입 | 예시 |
|---------|----------|------|
| HealthBiomarker | +5 | ProBNP(심부전), CEA(종양), PSA(전립선), Procalcitonin(패혈증), Troponin(심근경색) |
| Environmental | +3 | SoilHeavyMetal(토양중금속), MicroPlastic(미세플라스틱), Asbestos(석면) |
| FoodSafety | +3 | Mycotoxin(곰팡이독소), Antibiotic(항생제잔류), HeavyMetal(중금속) |
| AdvancedAnalysis | +1 | MultiEnvironment(복합환경) — **NonTarget1792(1792차원)는 v1.0 기본 등록 완료** |
| ElectronicSensor | +2 | ENoseAdvanced(16채널), ETongueAdvanced(16채널) |

### Phase 3~4 확장 예시 (+30종 이상)

| 신규 카테고리 | 예시 타입 |
|-------------|----------|
| Industrial (0x06) | ChemicalAgent, GasLeak, LubricantQuality, WeldingFume |
| Veterinary (0x07) | CanineBlood, FelineBlood, EquineBlood, LivestockPathogen |
| Pharmaceutical (0x08) | DrugPurity, DrugConcentration, DrugStability, Counterfeit |
| Agricultural (0x09) | SoilNutrient, FertilizerQuality, PlantDisease, PestPresence |
| Cosmetic (0x0A) | SkinType, IngredientPurity, PreservativeLevel, Allergenicity |

---

## 8. API 설계

### 8.1 Proto 메시지 (gRPC)

```protobuf
// 카트리지 카테고리
message CartridgeCategory {
  int32 code = 1;         // 카테고리 코드 (0x01~0xFF)
  string name_en = 2;
  string name_ko = 3;
  string description = 4;
  int32 type_count = 5;   // 등록된 타입 수
  bool is_active = 6;
}

// 카트리지 타입 정보
message CartridgeTypeInfo {
  int32 category_code = 1;
  int32 type_index = 2;
  int32 legacy_code = 3;
  string name_en = 4;
  string name_ko = 5;
  string description = 6;
  int32 required_channels = 7;
  int32 measurement_secs = 8;
  string unit = 9;
  string reference_range = 10;
  bool is_active = 11;
  bool is_beta = 12;
  string manufacturer = 13;
}

// 카트리지 접근 검증 요청
message CheckCartridgeAccessRequest {
  string user_id = 1;
  int32 category_code = 2;
  int32 type_index = 3;
}

// 카트리지 접근 검증 응답
message CheckCartridgeAccessResponse {
  bool allowed = 1;
  string access_level = 2;       // included, limited, add_on, restricted, beta
  int32 remaining_daily = 3;     // 일일 잔여 횟수 (-1 = 무제한)
  int32 remaining_monthly = 4;   // 월간 잔여 횟수 (-1 = 무제한)
  SubscriptionTier required_tier = 5;
  SubscriptionTier current_tier = 6;
  string message = 7;
  int32 addon_price_krw = 8;     // add_on인 경우 가격
}

// 사용자별 접근 가능 카트리지 목록 조회
message ListAccessibleCartridgesRequest {
  string user_id = 1;
}

message ListAccessibleCartridgesResponse {
  repeated CartridgeAccessEntry entries = 1;
}

message CartridgeAccessEntry {
  CartridgeTypeInfo type_info = 1;
  string access_level = 2;
  int32 remaining_daily = 3;
  int32 remaining_monthly = 4;
}
```

### 8.2 RPC 추가 (SubscriptionService 확장)

```protobuf
service SubscriptionService {
  // ... 기존 RPC ...
  
  // 카트리지 접근 권한 확인
  rpc CheckCartridgeAccess(CheckCartridgeAccessRequest) returns (CheckCartridgeAccessResponse);
  
  // 사용자별 접근 가능 카트리지 목록
  rpc ListAccessibleCartridges(ListAccessibleCartridgesRequest) returns (ListAccessibleCartridgesResponse);
}
```

### 8.3 MeasurementService StartSession 확장

```protobuf
message StartSessionRequest {
  string device_id = 1;
  string cartridge_id = 2;       // NFC UID
  string user_id = 3;
  int32 cartridge_category = 4;  // [신규] 카테고리 코드
  int32 cartridge_type_index = 5; // [신규] 타입 인덱스
}
```

---

## 9. 참조

| 문서 | 경로 |
|------|------|
| 기획안 v1.1 | docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md |
| 구독 티어 매핑 | docs/plan/terminology-and-tier-mapping.md |
| MSA 확장 로드맵 | docs/plan/msa-expansion-roadmap.md |
| 데이터 패킷 표준 | docs/specs/data-packet-family-c.md |
| NFC 모듈 구현 | rust-core/manpasik-engine/src/nfc/mod.rs |
| gRPC Proto | backend/shared/proto/manpasik.proto |
| 구독 서비스 | backend/services/subscription-service/ |

---

**문서 종료**

*본 명세서는 카트리지 체계의 무한 확장과 등급별 접근 제어를 위한 기준 문서이며, 개발 진행에 따라 갱신한다.*
