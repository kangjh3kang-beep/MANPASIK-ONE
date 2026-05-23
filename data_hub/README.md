# ManPaSik Data Hub — 개방 데이터 허브

> **에픽** E10 · **표준** Beacon v2, DUO, GA4GH Passport/AAI, Five Safes

## 구조

```
data_hub/
├── tre/       — Trusted Research Environment (compute-to-data)
├── beacon/    — Beacon v2 API (데이터 발견)
├── duo/       — Data Use Ontology (접근 등급)
└── passport/  — GA4GH Passport/AAI (연합 인증)
```

## 데이터 접근 3-tier

| 등급 | 데이터 | 접근 방법 | 요건 |
|------|--------|----------|------|
| Tier 1 | 집계 통계 | Beacon API | 공개 |
| Tier 2 | 비식별 데이터셋 | DUO 동의 기반 | 연구 윤리 승인 |
| Tier 3 | 원시 데이터 | TRE 워크벤치 | IRB + DAC 승인 |

## Five Safes 원칙

1. Safe People — 연구자 인증 (GA4GH Passport)
2. Safe Projects — 연구 윤리 승인 (IRB)
3. Safe Data — 비식별화 (k-anonymity, differential privacy)
4. Safe Settings — 보안 환경 (TRE, compute-to-data)
5. Safe Outputs — 결과 심사 (집계만 반출)

## 연합학습

- Flower/NVIDIA FLARE 기반 연합 노드
- 원시 데이터 이동 없이 가중치만 교환
- 차등 프라이버시 적용
