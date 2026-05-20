# ManPaSik Harness Capillary Map

**작성일**: 2026-05-01  
**작성자**: Codex  
**목적**: 메뉴, 상태, API, 서비스, 저장소, 이벤트, 테스트를 하나의 실행 맵으로 연결한다.

## 핵심 메뉴 연결 맵

| Menu | Route | Flutter State/Repository | API/RPC | Service | Storage/Event | H0 판정 |
|---|---|---|---|---|---|---|
| Home | `/home` | `measurementHistoryProvider`, `deviceListProvider` | history, devices, coaching | measurement, device, coaching | PostgreSQL/Timescale, Redis | 부분 결선 |
| Measure | `/measure`, `/measure/result` | `MeasurementGoldenPathOrchestrator`, `MeasurementRepository` | StartSession, StreamMeasurement/process, EndSession | measurement-service, RustBridge | Timescale, Milvus, Kafka | P0 골든 패스 대상 |
| Data Hub | `/data`, `/data/monitoring` | DataHub providers/repository | history, summaries, devices | measurement, data-platform, device | Timescale, ES, Redis | 부분 결선 |
| Devices | `/devices`, `/devices/:id` | `DeviceRepository` | register/list/status/ota | device-service, iot-gateway | PostgreSQL, Redis, Kafka | 부분 결선 |
| AI Coach | `/coach`, `/chat`, `/coach/food` | AI/coach repositories | analyze, coaching, vision | ai-inference, coaching, vision | PostgreSQL, LLM provider | 부분 결선 |
| Market | `/market/*` | Market repository | products, cart, orders, payments | shop, payment, subscription | PostgreSQL, Kafka | 부분 결선 |
| Medical | `/medical/*` | Medical repository | reservations, prescriptions, video | reservation, telemedicine, prescription | PostgreSQL, WebRTC provider | 부분 결선 |
| Family | `/family/*` | Family repository | family groups, alerts, reports | family, emergency, notification | PostgreSQL, Kafka, FCM | 부분 결선 |
| Settings/Admin | `/settings/*`, `/admin/*` | settings/admin providers | config, audit, export, compliance | admin, audit, health-record | PostgreSQL, MinIO, audit log | 부분 결선 |

## Measure 골든 패스 맵

```text
/measure
  -> MeasurementGoldenPathOrchestrator
  -> RustMeasurementEngine
      -> RustBridge.nfcReadCartridge
      -> RustBridge.runMeasurementPipeline
  -> MeasurementRepository.startSession
  -> MeasurementRepository.processMeasurement
      -> REST /api/v1/measurements/process
      -> Gateway handleProcessMeasurement
      -> gRPC StreamMeasurement
      -> measurement-service ProcessMeasurement
      -> MeasurementRepository.Store
      -> VectorRepository.StoreFingerprint
  -> MeasurementRepository.endSession
  -> Data Hub history
```

## 완료 기준

1. Route가 실제 Provider 또는 Orchestrator를 호출한다.
2. Provider는 demo/mock 여부를 명시적으로 드러낸다.
3. Repository는 실패를 빈 값으로 숨기지 않고 실패/오프라인/스테일 상태로 전달한다.
4. API/RPC는 서비스 저장소까지 도달한다.
5. 저장, 이벤트, 관측성, 테스트 증거가 연결된다.

