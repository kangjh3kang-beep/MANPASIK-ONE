# D-5: Predicate Device 조사 및 FDA 510(k) 제출 자료

> **문서 ID**: MPK-REG-D5  
> **버전**: 1.0  
> **작성일**: 2026-02-13  
> **작성자**: ManPaSik Regulatory Affairs Team  
> **상태**: 초안 (Draft)

---

## 섹션 목록

1. [ManPaSik 제품 분류](#1-manpasik-제품-분류)
2. [유사 기기 후보 5개](#2-유사-기기-후보-5개)
3. [Substantial Equivalence 분석 매트릭스](#3-substantial-equivalence-분석-매트릭스)
4. [Substantial Equivalence (SE) 평가](#4-substantial-equivalence-se-평가)
5. [선정 Predicate 결정 및 근거](#5-선정-predicate-결정-및-근거)
6. [참조 문헌](#6-참조-문헌)

---

## 1. ManPaSik 제품 분류

### 1.1 제품 개요

| 항목 | 내용 |
|------|------|
| **제품명** | ManPaSik (만파식) Point-of-Care Clinical Chemistry Analyzer |
| **제품 코드** | QKQ (Clinical Chemistry Analyzer, General Purpose) |
| **규제 등급** | Class II (Special Controls) |
| **적용 규정** | 21 CFR 862.2 (Clinical Chemistry and Clinical Toxicology Devices) |
| **적용 상세 규정** | 21 CFR 862.2160 — Discrete Photometric Chemistry Analyzer for Clinical Use |
| **제출 경로** | 510(k) Premarket Notification |
| **CLIA 분류 목표** | CLIA-waived (이후 waiver 검토 가능) |

### 1.2 제품 개요

ManPaSik은 소형 휴대형 Point-of-Care (POC) 임상화학 분석기로, 초소량 혈액 샘플(약 50–70 µL)을 사용하여 광학 및 전기화학 바이오마커를 정량 측정하는 체외진단(IVD) 의료기기입니다.

**핵심 특징:**
- NFC/BLE 기반 스마트 카트리지 인증 시스템
- AI 기반 측정 보정 및 코칭 서비스
- 스마트폰/태블릿 연동 (Flutter 앱)
- 클라우드 기반 데이터 동기화 및 원격진료 연동
- 896차원 측정 핑거프린트 기반 오염 탐지 및 품질 보증

**측정 항목:**
- 혈당 (Glucose)
- 지질패널 (Total Cholesterol, HDL, LDL, Triglycerides)
- 당화혈색소 (HbA1c)
- 간기능 (ALT, AST, GGT)
- 신기능 (Creatinine, BUN)
- 전해질 (Na+, K+, Cl-)
- CRP (C-Reactive Protein)

### 1.3 사용 목적 (Intended Use)

ManPaSik 임상화학 분석 시스템은 의료 전문가와 환자자가관리에 사용될 의료 환경(병원, 클리닉, 약국, 가정 등에서 혈액 채취 가능한 검체)에서 혈액의 정량적 임상화학 바이오마커를 측정함으로써 질병 진단·예방·치료 지원에 활용됩니다.

---

## 2. 유사 기기 후보 5개

### 2.1 Abbott i-STAT System

| 항목 | 내용 |
|------|------|
| **제조사** | Abbott Point of Care Inc. (Abbott Laboratories) |
| **510(k) 번호** | K153357 (i-STAT Alinity), K970720 (i-STAT 1 원형) |
| **규제 등급** | Class II |
| **제품 코드** | QKQ (Clinical Chemistry Analyzer) |
| **적용 규정** | 21 CFR 862.2160 |
| **CLIA 분류** | CLIA-waived (일부 테스트유닛) |

**제품 설명:**  
i-STAT는 휴대형 혈액분석 플랫폼으로, 스마트 카트리지 인증 및 일회용 전기화학 센서를 사용하여 2–3방울 혈액으로 약 2분 내에 결과를 제공합니다. 카트리지별로 전해질, 혈가스, 응고, 심장마커 등 광범위한 검사 메뉴를 제공합니다.

**측정 항목:**  
Na+, K+, Cl-, TCO2, BUN, Creatinine, Glucose, Lactate, Hematocrit, pH, pCO2, pO2, PT/INR, ACT, cTnI, BNP, CK-MB 등

| ManPaSik과의 유사점 | ManPaSik과의 차이점 |
|-------------------|---------------------|
| 동일 POC 플랫폼 + 스마트 카트리지 인증 | i-STAT는 전기화학 센서 중심, ManPaSik은 광학+NFC 이중 인증 |
| 초소량 혈액 샘플 사용 | ManPaSik은 AI 기반 측정 보정 알고리즘 탑재 |
| 동일 분류 임상화학 바이오마커 측정 | ManPaSik은 스마트폰 연동 휴대용 플랫폼 (i-STAT는 의료용 핸드헬드) |
| 클라우드/EMR 연동 가능 | ManPaSik은 896차원 핑거프린트 오염탐지 체계 적용 |
| Class II, 동일 제품코드(QKQ) | — |

---

### 2.2 Roche cobas b 101 System

| 항목 | 내용 |
|------|------|
| **제조사** | Roche Diagnostics GmbH |
| **510(k) 번호** | K131544 (일부 cobas 시리즈), CE 인증 주임 |
| **규제 등급** | Class II (EU: Class IIa) |
| **제품 코드** | QKQ / JJE |
| **적용 규정** | 21 CFR 862.2 / EU IVDR |

**제품 설명:**  
cobas b 101는 POC 검사 플랫폼으로 CRP, HbA1c, 지질패널 등을 4–6분 내에 측정합니다. 터치스크린 디스플레이, 5,000건 결과 저장, USB/RS422 연동을 갖춘 데스크톱형(234×135×184mm, 2kg) 기기입니다.

**측정 항목:**  
HbA1c, CRP, Total Cholesterol, HDL, LDL, Triglycerides

| ManPaSik과의 유사점 | ManPaSik과의 차이점 |
|-------------------|---------------------|
| 소형 POC 플랫폼, 유사한 크기/중량 | cobas b 101는 데스크톱형, ManPaSik은 스마트폰 연동 휴대용 |
| HbA1c + 지질패널 + CRP 측정 | 측정 항목 범위가 ManPaSik보다 협소 (간기능, 신기능, 전해질 미포함) |
| 터치스크린 UI | ManPaSik은 AI 보정 및 클라우드 동기화 기능 보유 |
| 스마트 카트리지/바코드 인증 사용 | cobas b 101는 CE 인증 주임, FDA 510(k)는 cobas 시리즈 일괄 확인 필요 |
| 병원 관리 시스템 및 IT 연동 | — |

---

### 2.3 Samsung LABGEO PT10

| 항목 | 내용 |
|------|------|
| **제조사** | Samsung Healthcare (Samsung Electronics) |
| **510(k) 번호** | K142498 (LABGEO IB10, 미국형 변형) |
| **규제 등급** | Class II |
| **제품 코드** | QKQ |
| **적용 규정** | 21 CFR 862.2160 |
| **글로벌 인증** | KFDA, CE 인증 |

**제품 설명:**  
LABGEO PT10은 휴대형 POC 임상화학 분석기로, 70 µL 혈액으로 자동 혈장분리를 수행해 7–10분 내에 최대 22종의 임상화학 바이오마커를 측정합니다. 4.3인치 터치스크린, 5,000건 결과 저장, 내장 열전사 프린터 등이 특징입니다.

**측정 항목:**  
ALT, AST, GGT, ALP, Total Bilirubin, Direct Bilirubin, Total Protein, Albumin, BUN, Creatinine, Uric Acid, Total Cholesterol, HDL, LDL, Triglycerides, Glucose, HbA1c, CRP, Amylase, Lipase 등 22종

| ManPaSik과의 유사점 | ManPaSik과의 차이점 |
|-------------------|---------------------|
| 동일 POC 임상화학 분석 플랫폼 | PT10은 데스크톱형 분석기, ManPaSik은 스마트폰 연동 플랫폼 |
| 유사한 샘플 용량 (70 µL vs ManPaSik 50–70 µL) | ManPaSik은 NFC/BLE 기반 무선 인증 카트리지 |
| 광범위한 임상화학 바이오마커 (22종) | ManPaSik은 AI 기반 코칭 서비스 및 원격 보정 |
| 자동 혈장분리 기능 | ManPaSik은 클라우드 동기화 및 원격진료 연동 플랫폼 |
| 소형/휴대형(2kg) | — |
| 터치스크린 디스플레이 | — |

---

### 2.4 PTS Diagnostics CardioChek Plus

| 항목 | 내용 |
|------|------|
| **제조사** | PTS Diagnostics (Polymer Technology Systems, Inc.) |
| **510(k) 번호** | K193406 |
| **규제 등급** | Class II |
| **제품 코드** | NBW, CGA, CHH, LBR, JGY |
| **적용 규정** | 21 CFR 862.1345 (Glucose), 21 CFR 862.1450 (Lipid) |
| **CLIA 분류** | CLIA-waived |

**제품 설명:**  
CardioChek Plus는 CLIA 면제 POC 분석기로, 40 µL의 손가락 채혈 혈액으로 지질패널과 혈당을 약 90초 내에 측정합니다. 스마트 카트리지/테스트스트립을 사용하며, 실온 보관이 가능한 소형 핸드헬드형입니다.

**측정 항목:**  
Total Cholesterol, HDL, Triglycerides, Glucose, Ketone (계산: LDL, TC/HDL Ratio, Non-HDL, VLDL)

| ManPaSik과의 유사점 | ManPaSik과의 차이점 |
|-------------------|---------------------|
| CLIA-waived POC 분석기 | CardioChek는 지질+혈당 전문, ManPaSik은 광범위한 바이오마커 |
| 초소량 혈액(손가락 채혈) 사용 | ManPaSik은 플랫폼 기능 확장 (클라우드, 원격진료) |
| 스마트 카트리지/테스트스트립 인증 | ManPaSik은 AI 보정 플랫폼 탑재 |
| 지질 + 혈당 측정 | 측정 폭 차이: CardioChek는 심혈관 중심, ManPaSik은 광범위 화학+전기화학 |

---

### 2.5 Siemens DCA Vantage Analyzer

| 항목 | 내용 |
|------|------|
| **제조사** | Siemens Healthcare Diagnostics Inc. |
| **510(k) 번호** | K071466 |
| **규제 등급** | Class II |
| **제품 코드** | LCP (HbA1c), CGX (Creatinine), JIR (Albumin) |
| **적용 규정** | 21 CFR 862.1165 (HbA1c), 21 CFR 862.1245 (Creatinine) |
| **CLIA 분류** | CLIA-waived (HbA1c) |

**제품 설명:**  
DCA Vantage는 반자동 POC 분석기로, 1 µL 혈액으로 HbA1c를 6–7분 내에 정량 측정합니다. 동일 플랫폼에서 소변 알부민·크레아티닌 및 Albumin/Creatinine Ratio(ACR)를 측정합니다. 4인치 터치스크린, 4,000건 결과 저장, 프린터 내장 등이 특징입니다. NGSP·IFCC 인증 당뇨 및 신장 검사를 사용합니다.

**측정 항목:**  
HbA1c, Urine Albumin, Urine Creatinine, Albumin/Creatinine Ratio (ACR)

| ManPaSik과의 유사점 | ManPaSik과의 차이점 |
|-------------------|---------------------|
| POC 환경 사용 (병원, 클리닉) | DCA Vantage는 HbA1c/신기능 전문, ManPaSik은 광범위 바이오마커 |
| HbA1c 정량 측정 | ManPaSik은 혈액 단일 채집, DCA는 혈액+소변 측정 |
| CLIA-waived 적용 | ManPaSik은 스마트폰 연동 플랫폼, DCA는 의료용 데스크톱 |
| 소형 데스크톱형 분석기 | ManPaSik은 AI/클라우드 기능 보유 |
| 병원 관리 시스템 및 연동 | — |

---

## 3. Substantial Equivalence 분석 매트릭스

### 3.1 사용 목적 비교

| 비교 항목 | ManPaSik | i-STAT (K153357) | cobas b 101 | LABGEO PT10 | CardioChek Plus (K193406) | DCA Vantage (K071466) |
|----------|-----------|-------------------|-------------|-------------|---------------------------|------------------------|
| **사용 환경** | POC (병원, 클리닉, 약국, 가정) | POC (병원, 클리닉) | POC (병원, 클리닉) | POC (병원, 클리닉) | POC (병원, 클리닉, 약국) | POC (병원, 클리닉) |
| **검체 종류** | 혈액 (정맥채취/캐피럴) | 혈액 | 혈액/캐피럴 | 혈액 | 혈액 (손가락 채혈) | 혈액/소변 |
| **사용자** | 의료 전문가/자가관리 | 의료 전문가 | 의료 전문가 | 의료 전문가 | 의료 전문가/자가관리 | 의료 전문가 |
| **측정 목적** | 진단, 예방, 치료 지원 | 진단, 예방, 치료 지원 | 진단, 예방, 치료 지원 | 진단, 예방, 치료 지원 | 스크리닝, 예방 | 진단, 예방, 치료 지원 |
| **처방 필요** | Rx Only | Rx Only | Rx Only | Rx Only | OTC (일부) | Rx Only |

### 3.2 기술적 특성 비교

| 비교 항목 | ManPaSik | i-STAT | cobas b 101 | LABGEO PT10 | CardioChek Plus | DCA Vantage |
|----------|-----------|--------|-------------|-------------|-----------------|-------------|
| **측정 방식** | 광학+전기화학 이중 인증 | 전기화학 (센서팁형) | 광학 (반사광 측정) | 광학 (포토메트릭) | 광학 (반사광 측정) | 면역화학+화학 |
| **샘플 용량** | 50–70 µL | 17–95 µL | ~2 µL | 70 µL | 40 µL | 1 µL |
| **측정 시간** | 5–10분 | 2분 | 4–6분 | 7–10분 | 90초 | 6–7분 |
| **분석 바이오마커 수** | 15+ | 30+ (테스트유닛별) | 6 | 22 | 5+계산값 | 4 |
| **인증 방식** | NFC 인증 내장 | 리더기 인증 | 카트리지/바코드 | 이더넷/바코드 | 카트리지/테스트스트립 | 의료용 인증 |
| **크기(mm)** | ~150×100×50 | 213×63×43 (핸드헬드) | 234×135×184 | 140×206×205 | 134×80×30 (핸드헬드) | 325×235×200 |
| **중량** | <0.5 kg (휴대용) | 0.52 kg | 2 kg | 2 kg | 0.17 kg | 4.5 kg |
| **전원** | 배터리(휴대용) | 배터리(휴대용) | AC 전원 | AC 전원/배터리 | 배터리(AAA) | AC 전원 |
| **연동** | NFC, BLE, Wi-Fi, 클라우드 | Wi-Fi, USB | USB, RS422 | USB | USB | USB, 프린터 |
| **결과 저장** | 클라우드 (규제 준수) | 5,000건 | 5,000건 | 5,000건 | 200건 | 4,000건 |
| **AI 기능** | 있음 (원격 보정, QC) | 없음 | 없음 | 없음 | 없음 | 없음 |

### 3.3 성능 비교 (주요 분석 바이오마커 정밀도 목표)

| 분석 바이오마커 | ManPaSik 목표 CV% | i-STAT CV% | cobas b 101 CV% | LABGEO PT10 CV% | CardioChek CV% |
|----------------|-------------------|------------|-----------------|-----------------|----------------|
| Glucose | TBD | 2–4% | N/A | 3–5% | 3–5% |
| Total Cholesterol | TBD | N/A | 3–5% | 3–5% | 3–5% |
| HbA1c | TBD | N/A | 2–3% | 3–5% | N/A |
| Creatinine | TBD | 3–5% | N/A | 4–6% | N/A |
| ALT | TBD | N/A | N/A | 5–8% | N/A |
| Triglycerides | TBD | N/A | 3–5% | 4–6% | 3–5% |
| CRP | TBD | N/A | 3–5% | 5–8% | N/A |

> **참고**: CV% = 변동계수 (Coefficient of Variation). ManPaSik 목표값은 추후 확정하며, 실제 성능은 분석적 검증 (Analytical Verification) 시정에 따라 확정 예정.

---

## 4. Substantial Equivalence (SE) 평가

### 4.1 SE 판단 기준

FDA 510(k) 경로에서 Substantial Equivalence를 입증하려면 다음을 충족해야 합니다:

1. **동일한 사용 목적 (Intended Use)**: 제출 기기가 predicate와 동일한 사용 목적을 가짐  
2. **동일한 기술적 특성 (Technological Characteristics)**: 동일하거나  
3. **상이한 기술적 특성임에도 상당·등가·안전성**: 상이한 기술이 사용목적 관점에서 상당·등가적 안전성 문제를 야기하지 않음

### 4.2 후보별 SE 적합성 요약

| Predicate 후보 | 사용 목적 적합 | 기술적 유사성 | SE 적합성 | 비고 |
|----------------|----------------|---------------|----------|------|
| **Abbott i-STAT** (K153357) | ◎ | ◎ | **최우선** | 동일 제품코드(QKQ), 최우수 predicate |
| **Roche cobas b 101** | ◎ | ◎ | **중간-최우선** | 유사 POC 용도, FDA 공개 K번호 확인 필요 |
| **Samsung LABGEO PT10** | ◎ | ◎ | **높음** | 기술적 유사성 우수, 추가 KFDA 인증 |
| **PTS CardioChek Plus** (K193406) | ◎ | △ | **중간** | CLIA-waived 경로 참고, 측정 범위 협소 |
| **Siemens DCA Vantage** (K071466) | ◎ | △ | **중간** | HbA1c 전문, 광범위 바이오마커 대비 제한적 |

### 4.3 신규 기술적 특성에 대한 비교 검증

ManPaSik은 기존 predicate 대비 다음과 같은 신규 기술적 특성을 포함합니다:

| 신규 기술 | 안전성 위험 | 이득 | 검증 방식 |
|----------|------------|------|----------|
| **NFC 기반 카트리지 인증 방식** | 낮음 (기존 RFID/바코드와 유사) | 없음 (동등 기능) | 전기적 안전성 검증 (IEC 60601-1) |
| **AI 측정 보정 알고리즘** | 중간 (소프트웨어 변경가능성) | 높음 (보정 정확도) | 분석적 검증 및 소프트웨어 V&V (IEC 62304) |
| **클라우드 기반 데이터 동기화** | 중간 (사이버보안) | 있음 | 사이버보안 대응 가이드 FDA Pre-Cert 참고 |
| **스마트폰 연결 플랫폼** | 있음 | 있음 | 사용성 공학 (IEC 62366-1), 소프트웨어 검증 |
| **896차원 핑거프린트 오염탐지** | 있음 | 높음 (품질 보증) | 임상 성능 검증, 임상 프로토콜 설계 등 |

> 신규 기술적 특성이 존재하나, 대부분 predicate의 안전성·유효성 범위를 확장하는 방향이며, 추가적인 위험평가를 통해 설명 가능합니다. 특히 AI 보정 및 오염탐지 기반 검증은 **성능 증폭** 요건으로, 추가적인 임상·분석 검증이 필요합니다.

---

## 5. 선정 Predicate 결정 및 근거

### 5.1 Primary Predicate: Abbott i-STAT Alinity System (K153357)

**선정 근거:**

1. **동일 규제 분류**: 제품 코드 QKQ, 21 CFR 862.2160 등 동일한 규제 경로  
2. **동일 사용 목적**: POC 환경에서 혈액 기반 다양한 임상화학 바이오마커 정량 측정  
3. **유사 기술 아키텍처**: 스마트 카트리지 인증 + 휴대형 설계, 센서팁형 바이오마커 측정  
4. **FDA 검증 이력**: 다수의 510(k) clearance 이력을 갖춘 검증된 predicate  
5. **광범위한 검사 메뉴**: 혈액, 전해질, 혈가스 등 다양한 바이오마커 범위

### 5.2 Secondary Predicate (참조): Samsung LABGEO PT10 (K142498)

**선정 근거:**

1. **기술적 유사성 우수**: 소형 POC, 자동 혈장분리, 70 µL 샘플, 22종 바이오마커  
2. **추가 규제 인증**: KFDA(식약처) 승인 등을 통한 추가 시장 진입 참고 가능  
3. **성능 비교 참고**: 다수의 peer-reviewed 성능 비교 논문 존재

### 5.3 보조 참조: PTS CardioChek Plus (K193406), Siemens DCA Vantage (K071466)

**참조 근거:**
- **CardioChek Plus**: CLIA-waived 경로, 초소량 채혈, 소형 핸드헬드 설계 참고  
- **DCA Vantage**: HbA1c·신기능 검사 POC 설계, CLIA-waived 검사 구조 참고  

---

## 6. 참조 문헌

1. FDA 510(k) Database: https://www.accessdata.fda.gov/scripts/cdrh/cfdocs/cfpmn/pmn.cfm  
2. 21 CFR Part 862 — Clinical Chemistry and Clinical Toxicology Devices  
3. FDA Guidance: "The 510(k) Program: Evaluating Substantial Equivalence" (2014)  
4. FDA Guidance: "Content of Premarket Submissions for Management of Cybersecurity in Medical Devices" (2023)  
5. CLSI EP05-A3: Evaluation of Precision of Quantitative Measurement Procedures  
6. CLSI EP09-A3: Measurement Procedure Comparison and Bias Estimation Using Patient Samples  
7. IEC 62304:2006+A1:2015 — Medical device software — Software life cycle processes  
8. ISO 14971:2019 — Medical devices — Application of risk management to medical devices  

---

*본 문서는 FDA 510(k) 제출을 위한 사전 조사 자료이며, 최종 제출 시 FDA Pre-Submission (Q-Sub) 회의를 통해 업데이트가 확정될 수 있습니다.*
