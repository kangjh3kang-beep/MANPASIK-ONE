# ManPaSik SDK — Core Module

만파식 플랫폼 SDK의 핵심 모듈입니다.

## 구성

```
sdk/manpasik-sdk/
├── core/          ← 핵심 측정·계약 라이브러리
├── cartridge/     — 카트리지 개발 도구 (CSI v1.0 14신호핀)
├── ai/            — AI 모델 통합 (TFLite/ONNX)
├── docs/          — SDK 문서
├── examples/      — 예제 코드
└── tools/
    ├── packager/  — 카트리지 패키징
    ├── simulator/ — 측정 시뮬레이터 (44 SKU)
    └── validator/ — 검증 도구
```

## 핵심 계약

- 표준 데이터 패킷: `contracts/packet_schema/standard_packet.json`
- LOINC 매핑: `contracts/mapping_registry/loinc_mapping.json`
- Kafka 이벤트: `contracts/events/kafka_topics.json`
- OpenAPI: `contracts/openapi/gateway-v1.yaml`

## 시뮬레이터 사용법

```bash
cd tools/simulator
python cartridge_simulator.py --sku glucose --alpha 0.98 --dim 88 --count 5
```

## 수익 배분

카트리지 스토어 수익 배분: 개발자 70% / 플랫폼 30% [추정, 확정 전 사람 승인 필요]
(config/baseline_params/pricing.json 참조)
