# ManPaSik Menu Completion Scorecard

**작성일**: 2026-05-01  
**작성자**: Codex  
**점수 기준**: route 20, real data 20, persistence 15, failure 15, offline 10, observability 10, test/trace 10

| Menu | Score | 주요 결손 | H1/H2 조치 |
|---|---:|---|---|
| Home | 64 | demo provider, 실패 상태 은닉 | 통합 dashboard state, stale 표시 |
| Measure | 58 | 오케스트레이션 분리 부족, process 결선 부족 | H1 골든 패스 구현 |
| Data Hub | 66 | chart와 history 계약 보강 필요 | history/summary/export 결선 |
| Devices | 62 | BLE 결과와 서버 등록 상태 분리 부족 | device state model 표준화 |
| AI Coach | 55 | 추천 근거 trace 부족 | reason code/source measurement 연결 |
| Market | 61 | entitlement 연결 보강 필요 | order-payment-entitlement E2E |
| Medical | 57 | case lifecycle 연결 부족 | MedicalCase 모델 도입 |
| Family | 54 | alert/event escalation 보강 필요 | family emergency event flow |
| Settings/Admin | 63 | dependency health/audit trace 부족 | readiness/audit/config approval |

## 릴리스 목표

핵심 메뉴는 85점 이상, Measure/Auth/Payment/Medical/Emergency는 92점 이상을 목표로 한다.

