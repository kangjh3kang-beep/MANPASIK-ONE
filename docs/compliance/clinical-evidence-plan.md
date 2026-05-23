# ManPaSik 임상 근거 계획 (Clinical Evidence Plan)

> **문서번호** MPS-CEP-v1.0 · **표준** MEDDEV 2.7/1 Rev.4, FDA 510(k)
> **적용** Class II IVD, De Novo (해당 시)

---

## 1. 임상 근거 전략

### 1.1 Predicate 디바이스

| 항목 | Predicate | ManPaSik | 동등성 |
|------|-----------|----------|--------|
| 측정 원리 | 전기화학 (바이오센서) | 차동측정 전기화학 | 유사 |
| 시료 유형 | 전혈/혈장 | 전혈/혈장/호기/타액 | 확장 |
| 사용 환경 | POCT / Home | POCT / Home | 동등 |
| 연결성 | BLE | BLE + NFC | 향상 |
| AI 보조 | 없음 | 트렌드·코칭 | 신규 기능 |

### 1.2 근거 유형

| 근거 | 내용 | 시기 |
|------|------|------|
| **문헌 검토** | Predicate 디바이스 성능 데이터 | Phase 1 |
| **분석적 성능** | LoD, LoQ, 정밀도, 정확도, 간섭, 안정성 | Phase 2 |
| **방법 비교** | Reference method vs ManPaSik (n≥120) | Phase 3 |
| **임상 성능** | 목표 집단 (n≥200) | Phase 4 |

---

## 2. 분석적 성능 검증 프로토콜

### 2.1 혈당 (대표 항목, ISO 15197:2013 기준)

| 시험 | 기준 | 방법 | 표본 |
|------|------|------|------|
| **정확도** | ±15 mg/dL (<100) 또는 ±15% (≥100) | 방법 비교 (YSI 2300) | n≥200 |
| **정밀도** | CV < 5% (within-run), < 7.5% (between-day) | 반복 측정 (5×5×5) | 3 농도 |
| **검출한계 (LoD)** | 신호 > 3σ blank | 공시료 반복 | n≥20 |
| **정량한계 (LoQ)** | CV < 20% at LoQ | 저농도 반복 | n≥20 |
| **간섭** | 편향 < 10% | CLSI EP7 | 간섭물 15종+ |
| **안정성** | 유효기간 내 사양 유지 | 가속 (40°C/75%RH) | 3/6/12/24개월 |

### 2.2 Parkes Consensus Error Grid

- Zone A: 임상적으로 정확한 결과 ≥95%
- Zone A+B: 임상적으로 수용 가능 ≥99%
- Zone C-E: 0% 목표

---

## 3. IEC 62366 사용성 시험 실행 계획

### 3.1 형성적 평가 (Formative, n=8~12)

**일정**: Phase 2 완료 후 4주
**환경**: 시뮬레이터 + 프로토타입 앱
**과업**:
1. 카트리지 삽입 및 NFC 인식 (≤30초)
2. 시료 적용 가이드 따라 측정 시작 (≤60초)
3. 결과 카드 해석 (3중 인코딩 이해도)
4. 긴급 알림 수신 및 대응

**측정 항목**:
- 과업 완료율: 목표 ≥80%
- 치명적 use error: 0건
- SUS 점수: ≥60
- 과업 소요 시간: 중앙값 기록

### 3.2 총괄적 평가 (Summative, n=15~20)

**일정**: Phase 3 완료 후 8주
**환경**: 실제 디바이스 + 정식 앱
**참가자**: 고령자 5명+ 포함, P1~P3 분포
**추가 과업**:
5. elder 모드에서 전체 측정 완주
6. 주치의 데이터 공유 동의/철회
7. 정기 배송 구독 설정

**수용 기준**:
- 과업 완료율: ≥90%
- 치명적 use error: 0건
- SUS 점수: ≥68
- 결과 오해율: <5% (3중 인코딩 효과 검증)

---

## 4. 510(k) 제출 체크리스트

| # | 항목 | 상태 | 문서 위치 |
|---|------|------|----------|
| 1 | Device Description | 🟡 초안 | docs/compliance/technical-file-structure.md |
| 2 | Predicate Comparison | 🟡 초안 | docs/compliance/predicate-device-research.md |
| 3 | Software Documentation (IEC 62304) | ✅ 완성 | docs/compliance/iec62304-*.md |
| 4 | Risk Analysis (ISO 14971) | ✅ 완성 | docs/compliance/iso14971-*.md |
| 5 | Software Safety Classification | ✅ 완성 | docs/compliance/software-safety-classification.md |
| 6 | Cybersecurity (FDA 524B) | 🟡 초안 | docs/security/stride-threat-model.md |
| 7 | Analytical Performance | ❌ 미완 | 실측 데이터 필요 |
| 8 | Clinical Performance | ❌ 미완 | 임상 시험 필요 |
| 9 | Biocompatibility | ❌ 미완 | 카트리지 소재 시험 필요 |
| 10 | Labeling/IFU | ❌ 미완 | eIFU + UDI-DI |
| 11 | HFE/Usability (IEC 62366) | 🟡 계획 | docs/compliance/iec62366-hfe-plan.md |
| 12 | AI/ML Model Cards (PCCP) | ✅ 완성 | docs/compliance/ai-model-cards/ |
| 13 | DHF Summary | 🟡 부분 | docs/compliance/technical-file-structure.md |
| 14 | SOUP List | ✅ 완성 | docs/compliance/compliance-gap-resolution.md |
| 15 | RTM | 🟡 부분 | docs/compliance/rtm-risk-linking.md |
| 16 | V&V Plan | ✅ 완성 | docs/compliance/vnv-master-plan.md |
| 17 | DPIA | ✅ 완성 | docs/compliance/dpia-template.md |

**현재 완성도**: 8/17 완성 (47%), 5/17 초안 (29%), 4/17 미완 (24%)
**미완 항목은 실측 데이터 + HW 프로토타입 의존**
