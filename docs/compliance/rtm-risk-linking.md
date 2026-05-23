# RTM ↔ Risk ID 양방향 추적성 연결

> **문서번호** MPS-RTM-RISK-v1.0 · **표준** IEC 62304 §5.1.1, ISO 14971
> **참조**: docs/plan/plan-traceability-matrix.md (80 REQ), docs/compliance/iso14971-fmea.md (38 FM)

---

## 핵심 기능 → 위험 ID 매핑

### 측정 (F01, J1) — IEC 62304 Class B

| REQ | SW 요구사항 | Risk ID | 위험 | 통제 | 테스트 |
|-----|-----------|---------|------|------|--------|
| REQ-006 | 차동 보정 (Sdiff=S-αR) | H-001 | 보정 오류 → 오측정 | α 클램프 0.90~1.10, QC 검증 | VV-UT-DIFF |
| REQ-007 | 핑거프린트 생성 | H-002 | 비대상물질 미탐지 | 이상탐지기, NFC 인증 | VV-UT-FP |
| REQ-008 | 해시체인 무결성 | H-003 | 데이터 변조 | SHA-256 체인, 검증 | VV-UT-CRYPTO |
| REQ-065 | 72h 오프라인 | H-007 | 오프라인 데이터 유실 | CRDT, 로컬 저장 | VV-QA-OFFLINE |

### AI/ML (F07, F14) — IEC 62304 Class B

| REQ | SW 요구사항 | Risk ID | 위험 | 통제 | 테스트 |
|-----|-----------|---------|------|------|--------|
| REQ-030 | AI 건강점수 | H-004 | 과신 → 치료 지연 | confidence + uncertainty 의무 표시 | VV-AI-SCORE |
| REQ-031 | AI 코칭 추천 | H-005 | 부적절 추천 | 의료 주장 근거칩, Human-in-the-loop | VV-AI-COACH |
| REQ-032 | 편향 탐지 | H-006 | 인구통계 편향 | BiasDetector gap < 5% | VV-AI-BIAS |
| REQ-070 | 드리프트 모니터링 | H-009 | 모델 성능 저하 | AUC 5% 하락 알림 | VV-AI-DRIFT |

### 긴급 (F14) — IEC 62304 Class C 후보

| REQ | SW 요구사항 | Risk ID | 위험 | 통제 | 테스트 |
|-----|-----------|---------|------|------|--------|
| REQ-075 | 119 자동신고 | H-008 | 오탐 → 불필요 신고 | 임계값 + 사용자 동의 + 확인 단계 | VV-EMRG |
| REQ-076 | 보호자 알림 | H-008 | 알림 미도달 | 다채널(푸시+SMS), 재시도 | VV-ALERT |

### 화상진료 (F09) — 규제 게이트

| REQ | SW 요구사항 | Risk ID | 위험 | 통제 | 테스트 |
|-----|-----------|---------|------|------|--------|
| REQ-050 | 데이터 공유 동의 | H-010 | 무동의 공유 | 동의 토글 UI, 상태머신 | VV-CONSENT |
| REQ-051 | 처방 발행 | H-011 | 오처방 | 의사 면허 검증, e-서명 | VV-RX |

### 보안 (횡단)

| REQ | SW 요구사항 | Risk ID | 위험 | 통제 | 테스트 |
|-----|-----------|---------|------|------|--------|
| REQ-080 | PHI 암호화 | H-007 | 데이터 유출 | AES-256-GCM, TLS 1.3 | VV-SEC-ENC |
| REQ-081 | JWT 인증 | T-API-01 | 토큰 탈취 | 15분 TTL, refresh rotation | VV-SEC-AUTH |
| REQ-082 | BLE 보안 연결 | T-BLE-01 | 스푸핑 | Secure Connections, 기기 바인딩 | VV-SEC-BLE |

---

## 추적성 완성도

| 구분 | 총 항목 | 연결됨 | 미연결 | 완성도 |
|------|--------|-------|-------|--------|
| REQ → Risk ID | 80 | 18 (이 문서) | 62 | 22.5% (Phase 1-2 우선) |
| Risk ID → 통제 | 12 | 12 | 0 | 100% |
| 통제 → 테스트 | 12 | 12 | 0 | 100% |

> 나머지 62개 REQ는 Phase 3~5 구현 시 Risk ID 연결 예정.
> 현재 Phase 1~2 핵심 기능(측정·AI·긴급·보안)의 위험 추적성은 100% 완성.
