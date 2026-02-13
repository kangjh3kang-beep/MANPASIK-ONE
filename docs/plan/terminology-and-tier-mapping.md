# 용어·티어 통일 및 카트리지 접근 제어 (Terminology, Tier & Cartridge Access Mapping)

**문서번호**: MPK-PLAN-TERM-v2.0  
**목적**: 구독 티어 명칭 통일, 카트리지 무한확장 체계의 등급별 접근 정책 공식 기준. 모든 신규 문서는 이 표를 기준으로 사용.

---

## 1. 구독 티어 대응표

| 원본 (기획안 v1.0) | 현재 (시스템·CONTEXT) | 월 요금 | Enum/코드 |
|-------------------|------------------------|---------|-----------|
| — | **Free** | 무료 | SUBSCRIPTION_TIER_FREE (0) |
| Basic Safety | **Basic** / Basic Safety | ₩9,900/월 | SUBSCRIPTION_TIER_BASIC (1) |
| Bio-Optimization | **Pro** / Bio-Optimization | ₩29,900/월 | SUBSCRIPTION_TIER_PRO (2) |
| Clinical Guard | **Clinical** / Clinical Guard | ₩59,900/월 | SUBSCRIPTION_TIER_CLINICAL (3) |

- **문서·UI**: "Basic Safety", "Bio-Optimization", "Clinical Guard"는 원본·마케팅용 명칭으로 유지 가능.
- **코드·API·DB**: `free` / `basic` / `pro` / `clinical` (소문자) 또는 위 Enum 사용.

---

## 2. 핵심 용어 정리

| 용어 | 정의 | 사용처 |
|------|------|--------|
| MANPASIK | 프로젝트·제품명. "Calming All Waves" | 전체 |
| MANPASIK World | 통합 생태계 브랜드명 | UX, 마케팅 |
| 차동측정 | S_corrected = S_det - α × S_ref | 기술 문서, 코드 |
| 패밀리C | 데이터 패킷·보안 모듈 관련 제조/인증 패밀리 | 데이터 패킷 표준 문서 |
| My Zone | 개인 기준선 (개인 정상 범위) | 앱, 코칭 |
| 리더기 | 만파식 측정 디바이스 (하드웨어) | 전체 |
| 카트리지 | 교체형 측정 모듈 (**무한확장**, 기본 29종 + 레지스트리 무제한 추가) | 전체 |
| 카트리지 카테고리 | 분석 영역별 카트리지 분류 (HealthBiomarker, Environmental 등) | 기술, API, DB |
| 접근 레벨 | 등급별 카트리지 사용 권한 (INCLUDED/LIMITED/ADD_ON/RESTRICTED/BETA) | 구독, 정책 |

---

## 3. 등급별 카트리지 접근 정책 (Tier-Cartridge Access Matrix)

### 3.1 접근 레벨 범례

| 기호 | 접근 레벨 | 설명 |
|------|----------|------|
| ✅ | INCLUDED | 구독에 포함, 무제한 사용 |
| ⏱ | LIMITED | 구독에 포함, 일/월 사용 횟수 제한 |
| 💰 | ADD_ON | 별도 구매 시 사용 가능 (건당/팩 과금) |
| 🔒 | RESTRICTED | 사용 불가 (상위 등급 필요) |
| 🧪 | BETA | 베타 테스트용 (신청 기반) |

### 3.2 기본 정책 매트릭스

| 카테고리 코드 | 카테고리명 | Free | Basic | Pro | Clinical |
|-------------|-----------|------|-------|-----|----------|
| 0x01 | **HealthBiomarker (기본 3종)** ¹ | ⏱ 일3회 | ✅ | ✅ | ✅ |
| 0x01 | **HealthBiomarker (나머지 11종)** | 🔒 | ✅ | ✅ | ✅ |
| 0x02 | **Environmental** | 🔒 | 💰 | ✅ | ✅ |
| 0x03 | **FoodSafety** | 🔒 | 💰 | ✅ | ✅ |
| 0x04 | **ElectronicSensor** | 🔒 | 🔒 | ✅ | ✅ |
| 0x05 | **AdvancedAnalysis** | 🔒 | 🔒 | 💰 | ✅ |
| 0x06 | **Industrial** | 🔒 | 🔒 | 🔒 | ✅ |
| 0x07 | **Veterinary** | 🔒 | 🔒 | 💰 | ✅ |
| 0x08 | **Pharmaceutical** | 🔒 | 🔒 | 🔒 | ✅ |
| 0x09 | **Agricultural** | 🔒 | 🔒 | 💰 | ✅ |
| 0x0A | **Cosmetic** | 🔒 | 💰 | ✅ | ✅ |
| 0x0B | **Forensic** | 🔒 | 🔒 | 🔒 | ✅ |
| 0x0C | **Marine** | 🔒 | 🔒 | 💰 | ✅ |
| 0xF0~0xFD | **ThirdParty (SDK)** | 🔒 | 🔒 | 💰 | 💰 |
| 0xFE | **Beta** | 🔒 | 🔒 | 🔒 | 🧪 |
| 0xFF | **CustomResearch** | 🔒 | 🔒 | 🔒 | ✅ |

> ¹ 기본 3종: Glucose(혈당), LipidPanel(지질패널), HbA1c(당화혈색소). Free 등급에서도 핵심 건강 모니터링 가능.

### 3.3 정책 적용 규칙

1. **우선순위**: 타입별 오버라이드 > 카테고리별 정책 > 글로벌 기본값(RESTRICTED)
2. **동적 관리**: 정책은 DB(`cartridge_tier_access`) 기반으로 관리자가 코드 배포 없이 변경 가능
3. **유효 기간**: 정책에 `effective_from`/`effective_until`을 지정하여 프로모션·시즌 한정 정책 운영 가능
4. **애드온 구매**: ADD_ON 카트리지는 쇼핑몰(shop-service)에서 사용권 팩 구매 후 사용
5. **오프라인**: 마지막 동기화 시점의 정책 캐시를 사용, 다음 온라인 시 갱신

### 3.4 등급별 핵심 차별화 요약

| 등급 | 카트리지 접근 | 리더기 | 기타 핵심 기능 |
|------|-------------|--------|-------------|
| **Free** | 기본 건강 3종(제한) | 1대 | 기본 측정, 이력 조회 |
| **Basic** | 건강 14종 + 환경·식품(ADD_ON) | 3대 | 데이터 내보내기, 가족 2명 |
| **Pro** | 표준 전체 + 고급(ADD_ON) | 5대 | AI 코칭, 트렌드 분석, 건강 점수 |
| **Clinical** | 전체 무제한 + 베타 | 10대 | 화상진료, 의료진 매칭, FHIR |

---

## 4. 문서 작성 시 준수

- 신규 기획·설계·API 문서에서 구독 티어를 쓸 때는 위 **1. 구독 티어 대응표**를 따름.
- 카트리지 관련 문서에서 "29종"이라고 쓰지 않고 "기본 29종 + 무한 확장"으로 표기.
- 카트리지 접근 정책은 위 **3. 등급별 카트리지 접근 정책**을 기준으로 함.
- 원본 v1.0에서 "3단계"라고 한 것은 유료 3단계(Basic Safety, Bio-Optimization, Clinical Guard)를 의미하며, "4단계"는 Free를 포함한 현재 시스템 기준.

---

**참조**: 원본 기획안 5.1/5.9, `docs/specs/cartridge-system-spec.md`, CONTEXT.md, AGENTS.md §4.5, `docs/plan/MPK-ECO-PLAN-v1.1-COMPLETE.md`
