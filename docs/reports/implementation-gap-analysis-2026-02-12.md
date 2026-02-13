# ManPaSik 구현 현황 vs 마스터플랜 Gap 분석 보고서

> **작성일**: 2026-02-12
> **작성자**: Claude Opus 4.5
> **목적**: 현재 구현 현황과 마스터플랜(v2.0) 비교, 미구현 사항 식별 및 구현 계획 수립
> **공유 대상**: 모든 IDE·AI 에이전트

---

## 1. 종합 현황 대시보드

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ManPaSik 구현 현황 대시보드 (2026-02-12)                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│   ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐        │
│   │  백엔드 서비스    │  │   Rust 코어      │  │   Flutter 앱     │        │
│   │  ████████░░ 79%  │  │  ███████░░░ 78%  │  │  ████░░░░░░ 42%  │        │
│   │  22/28 서비스    │  │  8/8 모듈        │  │  6/12 Feature    │        │
│   │                  │  │  (하드웨어 I/O   │  │  (P0: 2개 미구현)│        │
│   │  ✅ Phase 1-3    │  │   3개 미완성)    │  │                  │        │
│   │  ❌ 추가 6개     │  │                  │  │                  │        │
│   └──────────────────┘  └──────────────────┘  └──────────────────┘        │
│                                                                             │
│   전체 시스템 완성도: ████████░░░░░░░░░░░░ 66%                              │
│                                                                             │
│   P0 긴급 미구현: 11개 항목 (Rust 3 + Flutter 2 + Backend 6)                │
│   예상 소요 기간: 6-8주 (풀타임 개발팀 기준)                                   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. 영역별 상세 비교

### 2.1 백엔드 서비스 (Go Microservices)

| Phase | 서비스 | 마스터플랜 | 현재 상태 | 코드량 | Gap |
|-------|--------|-----------|----------|--------|-----|
| **Phase 1** | auth-service | ✅ 필수 | ✅ 완료 | 1,135줄 | - |
| | user-service | ✅ 필수 | ✅ 완료 | 1,175줄 | - |
| | device-service | ✅ 필수 | ✅ 완료 | 1,060줄 | - |
| | measurement-service | ✅ 필수 | ✅ 완료 | 1,696줄 | - |
| **Phase 2** | subscription-service | ✅ 필수 | ✅ 완료 | 1,478줄 | - |
| | shop-service | ✅ 필수 | ✅ 완료 | 1,495줄 | - |
| | payment-service | ✅ 필수 | ✅ 완료 | 1,075줄 | - |
| | ai-inference-service | ✅ 필수 | ✅ 완료 | 1,430줄 | - |
| | cartridge-service | ✅ 필수 | ✅ 완료 | 2,410줄 | - |
| | calibration-service | ✅ 필수 | ✅ 완료 | 1,995줄 | - |
| | coaching-service | ✅ 필수 | ✅ 완료 | 3,024줄 | - |
| **Phase 3+** | admin-service | ✅ 필수 | ✅ 완료 | 2,448줄 | - |
| | community-service | ✅ 필수 | ✅ 완료 | 1,834줄 | - |
| | family-service | ✅ 필수 | ✅ 완료 | 1,907줄 | - |
| | health-record-service | ✅ 필수 | ✅ 완료 | 2,883줄 | - |
| | notification-service | ✅ 필수 | ✅ 완료 | 1,836줄 | - |
| | prescription-service | ✅ 필수 | ✅ 완료 | 2,370줄 | - |
| | reservation-service | ✅ 필수 | ✅ 완료 | 2,688줄 | - |
| | telemedicine-service | ✅ 필수 | ✅ 완료 | 1,353줄 | - |
| | translation-service | ✅ 필수 | ✅ 완료 | 923줄 | - |
| | video-service | ✅ 필수 | ✅ 완료 | 1,144줄 | - |
| **추가** | analytics-service | 🔶 권장 | ❌ 미구현 | 0줄 | **P1** |
| | emergency-service | 🔴 필수 | ❌ 미구현 | 0줄 | **P0** |
| | iot-gateway-service | 🔶 권장 | ❌ 미구현 | 0줄 | **P2** |
| | marketplace-service | 🔶 권장 | ❌ 미구현 | 0줄 | **P2** |
| | nlp-service | 🔶 권장 | ❌ 미구현 | 0줄 | **P2** |
| | vision-service | 🔶 권장 | ❌ 미구현 | 0줄 | **P2** |

**백엔드 요약**:
- ✅ 완료: 22개 서비스 (약 45,000줄)
- ❌ 미구현: 6개 서비스
- **완성도: 79%**

---

### 2.2 Rust 코어 엔진

| 모듈 | 마스터플랜 기능 | 현재 상태 | 코드량 | 구현율 | Gap |
|------|---------------|----------|--------|--------|-----|
| **differential** | 차동측정 (α=0.95) | ✅ 완료 | 208줄 | 100% | - |
| **fingerprint** | 88→448→896→1792차원 | ✅ 완료 | 359줄 | 100% | - |
| **dsp** | FFT, 필터링, 피크 검출 | ✅ 완료 | 696줄 | 100% | - |
| **crypto** | AES-256-GCM, 해시체인 | ✅ 완료 | 574줄 | 100% | - |
| **sync** | CRDT 오프라인 동기화 | ✅ 완료 | 892줄 | 100% | - |
| **ai** | TFLite 엣지 추론 | ⚠️ 부분 | 376줄 | 50% | **P0** |
| | - 구조/시뮬레이션 | ✅ 완료 | - | - | - |
| | - 실제 TFLite 추론 | ❌ 미구현 | - | - | **50줄 추가** |
| **ble** | btleplug BLE 5.0 | ⚠️ 부분 | 395줄 | 40% | **P0** |
| | - 스캔 기능 | ✅ 완료 | - | - | - |
| | - 연결/데이터 전송 | ❌ 미구현 | - | - | **100줄 추가** |
| **nfc** | ISO 14443A 카트리지 | ⚠️ 부분 | 1,163줄 | 90% | **P0** |
| | - 데이터 파싱 (v1+v2) | ✅ 완료 | - | - | - |
| | - 하드웨어 I/O | ❌ 미구현 | - | - | **50줄 추가** |

**Rust 요약**:
- ✅ 완전 구현: 5개 모듈 (differential, fingerprint, dsp, crypto, sync)
- ⚠️ 부분 구현: 3개 모듈 (ai, ble, nfc) - **하드웨어 I/O 계층만 미완성**
- 총 코드: 5,165줄 (69개 테스트)
- **완성도: 78%**

---

### 2.3 Flutter 앱 Features

| Feature | 마스터플랜 | 현재 상태 | 코드량 | 아키텍처 | Gap |
|---------|-----------|----------|--------|----------|-----|
| **auth** | 인증/로그인 | ✅ 완료 | 639줄 | Clean Arch | - |
| **measurement** | 측정 화면 | ✅ 완료 | 931줄 | Clean Arch | - |
| **devices** | 디바이스 관리 | ✅ 완료 | 346줄 | Clean Arch | - |
| **home** | 홈 대시보드 | ⚠️ 부분 | 152줄 | Presentation만 | **P1** |
| **settings** | 설정 | ⚠️ 부분 | 255줄 | Presentation만 | **P1** |
| **user** | 프로필 | ⚠️ 부분 | 109줄 | Domain/Data만 | **P1** |
| **data_hub** | 데이터 분석/시각화 | ❌ 미구현 | 0줄 | - | **P0** |
| **ai_coach** | AI 건강 코칭 | ❌ 미구현 | 0줄 | - | **P0** |
| **market** | 쇼핑몰/마켓 | ❌ 미구현 | 0줄 | - | **P2** |
| **community** | 글로벌 커뮤니티 | ❌ 미구현 | 0줄 | - | **P2** |
| **family** | 가족 그룹 | ❌ 미구현 | 0줄 | - | **P2** |
| **medical** | 의료 서비스 | ❌ 미구현 | 0줄 | - | **P2** |

**Flutter 요약**:
- ✅ 완전 구현: 3개 Feature (auth, measurement, devices)
- ⚠️ 부분 구현: 3개 Feature (home, settings, user)
- ❌ 미구현: 6개 Feature (data_hub, ai_coach, market, community, family, medical)
- 총 코드: 2,432줄
- **완성도: 42%**

---

## 3. P0 긴급 미구현 항목 (11개)

### 3.1 Rust 코어 (3개)

| # | 항목 | 상세 | 예상 작업량 | 우선순위 |
|---|------|------|-----------|----------|
| R1 | AI TFLite 실제 추론 | `InferenceEngine::load_model()` + tflitec 통합 | 50-100줄 | **P0** |
| R2 | BLE 연결/데이터 전송 | `BleManager::connect()`, `start_measurement()` | 100-150줄 | **P0** |
| R3 | NFC 하드웨어 I/O | `NfcReader::read_cartridge()` 플랫폼별 구현 | 50-100줄 | **P0** |

### 3.2 Flutter 앱 (2개)

| # | 항목 | 상세 | 예상 화면 | 우선순위 |
|---|------|------|----------|----------|
| F1 | data_hub Feature | 측정 이력 분석, 그래프/차트, 건강 통계 | 4-6개 화면 | **P0** |
| F2 | ai_coach Feature | AI 코칭, 개인화 권장사항, 건강 목표 | 3-5개 화면 | **P0** |

### 3.3 백엔드 서비스 (6개)

| # | 항목 | 상세 | 예상 작업량 | 우선순위 |
|---|------|------|-----------|----------|
| B1 | emergency-service | 긴급상황 감지, 119 연동, 보호자 알림 | 2,000줄 | **P0** |
| B2 | analytics-service | 측정 통계, 리포팅, 대시보드 | 2,000줄 | **P1** |
| B3 | iot-gateway-service | IoT 기기 관리, 프로토콜 변환 | 1,500줄 | **P2** |
| B4 | marketplace-service | 서드파티 카트리지 마켓 | 2,500줄 | **P2** |
| B5 | nlp-service | 자연어 처리, 음성 명령 | 1,500줄 | **P2** |
| B6 | vision-service | 이미지 분석 (카트리지 검사 등) | 1,500줄 | **P2** |

---

## 4. 구현 로드맵

### 4.1 Phase 구분

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        구현 로드맵 (2026-02-12 기준)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  Week 1-2 (P0 긴급)                                                         │
│  ├── R1: Rust AI TFLite 추론 구현                                           │
│  ├── R2: Rust BLE 연결/전송 구현                                             │
│  ├── R3: Rust NFC 하드웨어 I/O 구현                                          │
│  └── B1: emergency-service 백엔드                                           │
│                                                                             │
│  Week 3-4 (P0 Flutter)                                                      │
│  ├── F1: Flutter data_hub Feature                                           │
│  │   ├── domain/data 레이어                                                 │
│  │   ├── 측정 이력 분석 화면                                                 │
│  │   ├── 그래프/차트 (fl_chart)                                             │
│  │   └── 건강 통계 대시보드                                                  │
│  └── F2: Flutter ai_coach Feature                                           │
│      ├── coaching-service gRPC 연동                                         │
│      ├── AI 권장사항 화면                                                    │
│      └── 건강 목표 관리                                                      │
│                                                                             │
│  Week 5-6 (P1 개선)                                                         │
│  ├── B2: analytics-service 백엔드                                           │
│  ├── Flutter home/settings/user 아키텍처 개선                                │
│  └── E2E 통합 테스트 강화                                                    │
│                                                                             │
│  Week 7-8 (P2 확장)                                                         │
│  ├── B3-B6: 추가 백엔드 서비스                                               │
│  ├── Flutter market/community/family/medical                                │
│  └── 성능 최적화 및 프로덕션 준비                                             │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 상세 구현 계획

#### Week 1-2: Rust 하드웨어 I/O 완성

**R1: AI TFLite 실제 추론 (3일)**

```rust
// 목표: manpasik-engine/src/ai/mod.rs 수정
impl InferenceEngine {
    pub fn load_model(&mut self, model_path: &Path) -> Result<(), InferenceError> {
        #[cfg(feature = "ai")]
        {
            use tflitec::model::Model;
            use tflitec::interpreter::Interpreter;

            let model = Model::from_file(model_path)?;
            let interpreter = Interpreter::new(&model, None)?;
            interpreter.allocate_tensors()?;

            self.interpreter = Some(interpreter);
            self.model_loaded = true;
        }
        Ok(())
    }

    pub fn predict(&self, input: &[f32]) -> Result<InferenceResult, InferenceError> {
        #[cfg(feature = "ai")]
        if let Some(ref interpreter) = self.interpreter {
            let input_tensor = interpreter.input(0)?;
            input_tensor.copy_from_buffer(input)?;

            interpreter.invoke()?;

            let output_tensor = interpreter.output(0)?;
            let output: Vec<f32> = output_tensor.data().to_vec();
            // ... 결과 처리
        }
    }
}
```

**R2: BLE 연결/데이터 전송 (4일)**

```rust
// 목표: manpasik-engine/src/ble/mod.rs 수정
impl BleManager {
    pub async fn connect(&mut self, device_id: &str) -> Result<(), BleError> {
        #[cfg(feature = "ble")]
        {
            use btleplug::api::{Central, Peripheral};

            let manager = Manager::new().await?;
            let adapter = manager.adapters().await?.into_iter().next()?;

            let peripheral = adapter.peripherals().await?
                .into_iter()
                .find(|p| p.id().to_string() == device_id)?;

            peripheral.connect().await?;
            peripheral.discover_services().await?;

            self.connected_devices.insert(device_id.to_string(), ConnectedDevice {
                peripheral,
                state: ConnectionState::Connected,
            });
        }
        Ok(())
    }

    pub async fn start_measurement(&self, device_id: &str) -> Result<(), BleError> {
        if let Some(device) = self.connected_devices.get(device_id) {
            let char = device.peripheral
                .characteristics()
                .iter()
                .find(|c| c.uuid == MEASUREMENT_CONTROL_CHAR)?;

            device.peripheral.write(char, &[CMD_START_MEASUREMENT], WriteType::WithResponse).await?;
        }
        Ok(())
    }
}
```

**R3: NFC 하드웨어 I/O (3일)**

```rust
// 목표: manpasik-engine/src/nfc/mod.rs 수정
impl NfcReader {
    pub async fn read_cartridge(&mut self) -> Result<CartridgeInfo, NfcError> {
        #[cfg(target_os = "android")]
        {
            // Android NFC API 호출 (Flutter FFI를 통해)
            let tag_data = self.platform_read_nfc_tag().await?;
            return self.parse_tag_data(&tag_data);
        }

        #[cfg(target_os = "ios")]
        {
            // iOS Core NFC 호출 (Flutter FFI를 통해)
            let tag_data = self.platform_read_nfc_tag().await?;
            return self.parse_tag_data(&tag_data);
        }

        Err(NfcError::PlatformNotSupported)
    }
}
```

**B1: emergency-service (3일)**

```go
// backend/services/emergency-service/
// 주요 기능:
// 1. 긴급상황 감지 (이상 바이오마커 패턴)
// 2. 119 연동 API
// 3. 보호자 긴급 알림 (family-service 연계)
// 4. 응급 의료 정보 전송

type EmergencyService struct {
    proto.UnimplementedEmergencyServiceServer
    repo       EmergencyRepository
    familySvc  proto.FamilyServiceClient
    notifySvc  proto.NotificationServiceClient
    kafka      *events.KafkaClient
}

func (s *EmergencyService) TriggerEmergency(ctx context.Context, req *proto.TriggerEmergencyRequest) (*proto.TriggerEmergencyResponse, error) {
    // 1. 응급 상황 기록
    emergency := s.repo.CreateEmergency(ctx, req.UserId, req.EmergencyType, req.Location)

    // 2. 보호자 알림
    guardians, _ := s.familySvc.GetGuardians(ctx, &proto.GetGuardiansRequest{UserId: req.UserId})
    for _, g := range guardians.Guardians {
        s.notifySvc.SendUrgentNotification(ctx, &proto.SendUrgentNotificationRequest{
            UserId:  g.Id,
            Type:    "EMERGENCY",
            Title:   "긴급 상황 발생",
            Message: fmt.Sprintf("%s님에게 긴급 상황이 발생했습니다", req.UserName),
        })
    }

    // 3. 119 연동 (한국 기준)
    if req.Call119 {
        s.call119API(emergency)
    }

    return &proto.TriggerEmergencyResponse{EmergencyId: emergency.Id}, nil
}
```

#### Week 3-4: Flutter P0 Features

**F1: data_hub Feature 구조**

```
lib/features/data_hub/
├── domain/
│   ├── entities/
│   │   ├── measurement_stats.dart
│   │   └── health_trend.dart
│   ├── repositories/
│   │   └── analytics_repository.dart
│   └── usecases/
│       ├── get_measurement_stats.dart
│       └── get_health_trends.dart
├── data/
│   ├── datasources/
│   │   └── analytics_remote_datasource.dart
│   ├── models/
│   │   ├── measurement_stats_model.dart
│   │   └── health_trend_model.dart
│   └── repositories/
│       └── analytics_repository_impl.dart
└── presentation/
    ├── providers/
    │   └── analytics_provider.dart
    ├── screens/
    │   ├── data_hub_screen.dart         # 메인 허브
    │   ├── measurement_history_screen.dart
    │   ├── health_trends_screen.dart
    │   └── detailed_stats_screen.dart
    └── widgets/
        ├── trend_chart.dart             # fl_chart
        ├── stats_card.dart
        └── period_selector.dart
```

**F2: ai_coach Feature 구조**

```
lib/features/ai_coach/
├── domain/
│   ├── entities/
│   │   ├── coaching_recommendation.dart
│   │   └── health_goal.dart
│   ├── repositories/
│   │   └── coaching_repository.dart
│   └── usecases/
│       ├── get_daily_coaching.dart
│       └── manage_health_goals.dart
├── data/
│   ├── datasources/
│   │   └── coaching_remote_datasource.dart  # gRPC coaching-service
│   └── repositories/
│       └── coaching_repository_impl.dart
└── presentation/
    ├── providers/
    │   └── coaching_provider.dart
    ├── screens/
    │   ├── ai_coach_screen.dart         # 메인 코칭 화면
    │   ├── recommendations_screen.dart
    │   └── goals_screen.dart
    └── widgets/
        ├── recommendation_card.dart
        ├── goal_progress.dart
        └── ai_avatar.dart
```

---

## 5. 리소스 요구사항

### 5.1 개발 리소스

| 단계 | 기간 | 인력 | 주요 작업 |
|------|------|------|----------|
| Week 1-2 | 2주 | 백엔드 1명 + Rust 1명 | 하드웨어 I/O, emergency-service |
| Week 3-4 | 2주 | Flutter 2명 | data_hub, ai_coach Feature |
| Week 5-6 | 2주 | 풀스택 2명 | 아키텍처 개선, 통합 테스트 |
| Week 7-8 | 2주 | 풀스택 2명 | 추가 서비스, 성능 최적화 |

### 5.2 기술 의존성

| 기술 | 버전 | 용도 |
|------|------|------|
| tflitec | 2.14+ | Rust TFLite 바인딩 |
| btleplug | 0.11+ | Rust BLE 라이브러리 |
| fl_chart | 0.69+ | Flutter 차트 라이브러리 |
| flutter_riverpod | 2.4+ | 상태 관리 |

---

## 6. 위험 요소 및 대응 방안

| 위험 | 영향 | 확률 | 대응 방안 |
|------|------|------|----------|
| TFLite Rust 바인딩 호환성 | 높음 | 중간 | tract-onnx 대안 검토 |
| BLE 플랫폼별 차이 | 중간 | 높음 | Flutter 네이티브 플러그인 사용 |
| NFC iOS 제한 | 높음 | 높음 | Core NFC + Background Tag Reading |
| 테스트 커버리지 | 중간 | 중간 | CI/CD에 최소 80% 커버리지 게이트 |

---

## 7. 성공 기준 (Exit Criteria)

### 7.1 P0 완료 기준

- [ ] Rust AI: TFLite 모델 로드 및 추론 성공 (5개 모델 타입)
- [ ] Rust BLE: 리더기 연결 및 측정 데이터 수신 성공
- [ ] Rust NFC: 카트리지 태그 읽기 성공 (v1.0 + v2.0)
- [ ] Flutter data_hub: 측정 이력 그래프 표시 성공
- [ ] Flutter ai_coach: AI 권장사항 표시 성공
- [ ] Backend emergency: 긴급 알림 전송 성공

### 7.2 전체 완료 기준

- [ ] 시스템 완성도 95% 이상
- [ ] E2E 테스트 전체 통과
- [ ] 성능: 측정→분석→코칭 3초 이내
- [ ] 오프라인 모드: 100% 로컬 동작

---

## 8. 다음 단계

1. **즉시**: Rust AI TFLite 구현 시작 (R1)
2. **이번 주**: BLE/NFC 하드웨어 I/O 병렬 진행 (R2, R3)
3. **다음 주**: Flutter data_hub Feature 착수 (F1)
4. **2주 후**: ai_coach Feature + emergency-service 완료

---

**문서 종료**

*이 보고서는 모든 IDE 및 AI 에이전트가 참조하여 개발 진행 상황을 파악하고 협업할 수 있도록 작성되었습니다.*
