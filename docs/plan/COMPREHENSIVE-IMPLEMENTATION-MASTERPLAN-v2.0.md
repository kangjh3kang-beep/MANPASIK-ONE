# ManPaSik 통합 세부 구현 마스터플랜 v2.0

> **문서번호**: MPK-MASTER-IMPL-v2.0-20260212
> **작성자**: Claude Opus 4.5
> **목적**: 유사 시스템 조사 및 최신 기술 트렌드 반영, AI 강점 극대화, 유기적 시스템 설계
> **공유 대상**: 모든 IDE·AI 에이전트
> **기반 조사**: McKinsey Healthcare AI, Corti Multi-Agent Framework, Flower FL, IEC 62304, TFLite Micro, Embassy-rs TrouBLE

---

## 목차

1. [조사 기반 설계 원칙](#1-조사-기반-설계-원칙)
2. [유기적 시스템 아키텍처 (Living System)](#2-유기적-시스템-아키텍처-living-system)
3. [AI 활용 전략 및 극대화 방안](#3-ai-활용-전략-및-극대화-방안)
4. [Rust 코어 세부 구현 기획](#4-rust-코어-세부-구현-기획)
5. [Flutter Feature 세부 구현 기획](#5-flutter-feature-세부-구현-기획)
6. [Phase 3-5 세부 구현 기획](#6-phase-3-5-세부-구현-기획)
7. [규정 문서 작성 계획](#7-규정-문서-작성-계획)
8. [시너지 극대화 연동 설계](#8-시너지-극대화-연동-설계)
9. [참조 문헌 및 출처](#9-참조-문헌-및-출처)

---

## 1. 조사 기반 설계 원칙

### 1.1 유사 시스템 및 기술 트렌드 조사 결과

#### Healthcare AI Architecture (2025-2026)

| 출처 | 핵심 인사이트 | 적용 방안 |
|------|-------------|----------|
| [McKinsey Healthcare AI](https://www.mckinsey.com/industries/healthcare/our-insights/the-coming-evolution-of-healthcare-ai-toward-a-modular-architecture) | Modular Architecture - 도메인 모델, 지능형 에이전트, 데이터 거버넌스 | 만파식 30+ 마이크로서비스를 도메인별 모듈로 조직화 |
| [Corti Multi-Agent Framework](https://www.corti.ai) | Multi-Agent AI - 실시간 의료 의사결정 지원 | coaching-service, ai-inference-service에 멀티에이전트 적용 |
| [World Economic Forum](https://www.weforum.org/stories/2026/01/ai-healthcare-data-architecture/) | Real-time Data Pipeline - 센서 → 정제 → 수치화 → AI | 차동측정 → 핑거프린트 → Milvus 파이프라인 최적화 |
| [Frontiers Digital Health](https://www.frontiersin.org/journals/digital-health/articles/10.3389/fdgth.2025.1694839/full) | Multi-Pattern Strategy - 마이크로서비스 + 블록체인 + Edge-Cloud | 마이크로서비스 + 해시체인 + Rust Edge AI |

#### Biosensor + ML 연구 동향

| 출처 | 핵심 기술 | 적용 방안 |
|------|----------|----------|
| [Nature - Plasma Infrared Fingerprinting](https://pmc.ncbi.nlm.nih.gov/articles/PMC11293328/) | FTIR + Multi-task Classification → 대사증후군 예측 | 차동측정 스펙트럼 + 멀티태스크 AI 분류 모델 |
| [Wiley - AI Biosensors](https://advanced.onlinelibrary.wiley.com/doi/full/10.1002/adma.202504796) | ML-augmented Biosensor → 정확도/민감도/속도 향상 | TFLite 엣지 AI로 실시간 바이오마커 분류 |
| [RSC - Surface-Enhanced Spectroscopy](https://pubs.rsc.org/en/content/articlehtml/2023/na/d2na00608a) | SERS + ML → Molecular Diagnostics | 896차원 핑거프린트 + 코사인 유사도 검색 |

#### Federated Learning (프라이버시 보존 AI)

| 출처 | 핵심 기술 | 적용 방안 |
|------|----------|----------|
| [Nature Scientific Reports](https://www.nature.com/articles/s41598-025-04083-4) | FL + Blockchain + Differential Privacy | Flower + 해시체인 + 노이즈 주입 |
| [PMC Federated Learning Review](https://pmc.ncbi.nlm.nih.gov/articles/PMC11728217/) | FL + IoT + Predictive Analytics | 리더기 분산 학습 + 건강 예측 |
| [JMIR AI - Personal Health Train](https://ai.jmir.org/2025/1/e60847) | Privacy-Preserving Analytics | 사용자 데이터 로컬 유지 + 모델만 동기화 |

### 1.2 핵심 설계 원칙

```
┌─────────────────────────────────────────────────────────────────────┐
│                    ManPaSik 설계 5대 원칙                            │
├─────────────────────────────────────────────────────────────────────┤
│ 1. 유기적 연동 (Organic Integration)                                 │
│    - 시스템 전체가 하나의 생물처럼 자율적으로 반응                     │
│    - 이벤트 기반 비동기 통신으로 느슨한 결합                          │
│                                                                     │
│ 2. AI 강점 극대화 (AI-First Design)                                  │
│    - 모든 데이터 흐름에 AI 추론 내장                                  │
│    - 예측적 UX (Predictive UX) 적용                                  │
│                                                                     │
│ 3. 프라이버시 보존 (Privacy by Design)                               │
│    - 연합학습으로 데이터 로컬 유지                                    │
│    - 차등 프라이버시 + 동형암호                                       │
│                                                                     │
│ 4. 오프라인 우선 (Offline-First)                                     │
│    - 100% 로컬 동작 가능                                             │
│    - CRDT 기반 충돌 해결                                             │
│                                                                     │
│ 5. 규제 내장 (Compliance by Design)                                  │
│    - IEC 62304 Class B 전 과정 적용                                  │
│    - 감사 추적 자동화                                                │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 2. 유기적 시스템 아키텍처 (Living System)

### 2.1 생물학적 메타포 기반 설계

만파식 시스템을 **하나의 생물체**로 설계합니다. 각 구성요소는 생물의 기관처럼 역할을 수행하며, 신경계(이벤트 버스)를 통해 유기적으로 연동됩니다.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      ManPaSik Living System Architecture                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────────┐                                                       │
│   │   🧠 두뇌    │  AI-Inference / Coaching / Prediction                │
│   │  (Brain)    │  - 패턴 인식, 의사결정, 학습                           │
│   └──────┬──────┘                                                       │
│          │                                                              │
│   ┌──────▼──────┐                                                       │
│   │ 🔮 신경계   │  Kafka/Redpanda Event Bus                            │
│   │ (Nervous)  │  - 실시간 신호 전달, 반응 조정                         │
│   └──────┬──────┘                                                       │
│          │                                                              │
│   ┌──────┴──────────────────────────────────────────┐                   │
│   │                                                  │                   │
│   ▼                                                  ▼                   │
│ ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│ │ 👀 감각기관 │  │ 💪 운동기관 │  │ 🫀 순환계   │  │ 🛡️ 면역계   │     │
│ │ (Sensors)   │  │ (Actuators) │  │ (Circu.)    │  │ (Immune)    │     │
│ ├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤     │
│ │ • BLE 리더기│  │ • OTA 업데이트│ │ • API Gateway│ │ • Auth      │     │
│ │ • NFC 카트리지│ │ • 알림 발송  │ │ • 데이터 동기화│ │ • RBAC     │     │
│ │ • 차동측정  │  │ • 코칭 실행  │ │ • 이벤트 전파│ │ • Rate Limit│     │
│ └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘     │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────┐     │
│ │ 🧬 유전자 (DNA) - 핵심 알고리즘                                    │     │
│ │ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │     │
│ │ │ Differential│ │ Fingerprint │ │ CRDT Sync   │ │ Crypto      │ │     │
│ │ │ 차동측정    │ │ 896차원     │ │ 오프라인    │ │ 보안        │ │     │
│ │ └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │     │
│ └─────────────────────────────────────────────────────────────────┘     │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────┐     │
│ │ 🏠 기관 (Organs) - 마이크로서비스                                  │     │
│ │ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐   │     │
│ │ │ Auth  │ │ User  │ │Device │ │Measure│ │ Shop  │ │Payment│   │     │
│ │ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘   │     │
│ │ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐   │     │
│ │ │Subscr.│ │Coaching│ │AI-Inf │ │Cartrid│ │Calibr.│ │Family │   │     │
│ │ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘ └───────┘   │     │
│ │                    ... 21+ 서비스 ...                          │     │
│ └─────────────────────────────────────────────────────────────────┘     │
│                                                                         │
│ ┌─────────────────────────────────────────────────────────────────┐     │
│ │ 💾 기억 (Memory) - 데이터 저장소                                   │     │
│ │ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ │     │
│ │ │ PostgreSQL  │ │ TimescaleDB │ │ Milvus      │ │ Redis       │ │     │
│ │ │ 영속 데이터 │ │ 시계열      │ │ 벡터 검색   │ │ 캐시        │ │     │
│ │ └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘ │     │
│ └─────────────────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 이벤트 기반 유기적 연동 (Event-Driven Organic Integration)

#### 핵심 이벤트 흐름

```yaml
# 측정 완료 이벤트 → 전체 시스템 유기적 반응
Event: measurement.completed
Producer: measurement-service
Consumers:
  - ai-inference-service:     # 두뇌 반응
      action: AnalyzeBiomarkers
      priority: HIGH
      timeout: 5s

  - coaching-service:         # 코칭 반응
      action: GenerateRecommendation
      priority: MEDIUM
      depends_on: ai-inference-service.completed

  - notification-service:     # 알림 반응
      action: SendHealthAlert
      condition: anomaly_detected == true

  - health-record-service:    # 기록 반응
      action: UpdateTimeline
      priority: LOW

  - family-service:           # 가족 공유 반응
      action: NotifyGuardians
      condition: user.has_guardians && anomaly_detected

  - subscription-service:     # 사용량 추적
      action: TrackUsage
      cartridge_type: event.cartridge_type
```

#### 자율 치유 (Self-Healing) 메커니즘

```go
// 서비스 장애 시 자동 복구 패턴
type SelfHealingConfig struct {
    CircuitBreaker struct {
        FailureThreshold   int           `yaml:"failure_threshold"`    // 5
        SuccessThreshold   int           `yaml:"success_threshold"`    // 3
        Timeout            time.Duration `yaml:"timeout"`              // 30s
        HalfOpenMaxCalls   int           `yaml:"half_open_max_calls"`  // 10
    }
    Retry struct {
        MaxAttempts        int           `yaml:"max_attempts"`         // 3
        InitialBackoff     time.Duration `yaml:"initial_backoff"`      // 100ms
        MaxBackoff         time.Duration `yaml:"max_backoff"`          // 5s
        BackoffMultiplier  float64       `yaml:"backoff_multiplier"`   // 2.0
    }
    Fallback struct {
        Enabled            bool          `yaml:"enabled"`              // true
        CacheEnabled       bool          `yaml:"cache_enabled"`        // true
        CacheTTL           time.Duration `yaml:"cache_ttl"`            // 5m
        GracefulDegradation bool         `yaml:"graceful_degradation"` // true
    }
}
```

---

## 3. AI 활용 전략 및 극대화 방안

### 3.1 AI 활용 영역 매트릭스

| 영역 | AI 기술 | 적용 서비스 | 사용자 가치 |
|------|---------|------------|------------|
| **바이오마커 분류** | Multi-task CNN | ai-inference | 92-96% 정확도 진단 |
| **이상치 탐지** | Isolation Forest + LSTM | ai-inference | 실시간 위험 경고 |
| **건강 예측** | Transformer Time-Series | coaching | 미래 건강 상태 예측 |
| **개인화 코칭** | Recommendation Engine | coaching | 맞춤형 건강 조언 |
| **음식 인식** | Vision Transformer (ViT) | vision | 사진 → 칼로리 자동 계산 |
| **음성 명령** | Whisper + LLM | nlp | 핸즈프리 측정 제어 |
| **실시간 번역** | mBART / NLLB | translation | 글로벌 커뮤니티 |
| **이상 패턴 학습** | Federated Learning | ai-training | 프라이버시 보존 모델 개선 |

### 3.2 AI 파이프라인 상세 설계

#### 3.2.1 측정 데이터 → AI 추론 → 코칭 파이프라인

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AI Pipeline: 측정 → 분석 → 코칭                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Stage 1: 데이터 수집 (Sensing)                                          │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ BLE Device → Raw Signal (88ch × 1024 samples) → Rust DSP        │   │
│  │                                                                  │   │
│  │ 차동측정: S_corrected = S_det - α × S_ref                        │   │
│  │ 결과: 88차원 정제 신호 벡터                                       │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              ↓                                          │
│  Stage 2: 특징 추출 (Feature Extraction)                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ 88D → FFT → Spectral Features → 448D                            │   │
│  │ 448D → Autoencoder → Latent Space → 896D                        │   │
│  │ 896D → Temporal Aggregation → 1792D (Phase 5)                   │   │
│  │                                                                  │   │
│  │ 결과: 896차원 핑거프린트 벡터                                     │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              ↓                                          │
│  Stage 3: AI 추론 (Inference)                                           │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ ┌─────────────┐   ┌─────────────┐   ┌─────────────┐            │   │
│  │ │ Classifier  │   │ Anomaly Det.│   │ Trend Pred. │            │   │
│  │ │ (TFLite)    │   │ (Isolation) │   │ (LSTM)      │            │   │
│  │ │             │   │             │   │             │            │   │
│  │ │ Input: 896D │   │ Input: 896D │   │ Input: 시계열│            │   │
│  │ │ Output:     │   │ Output:     │   │ Output:     │            │   │
│  │ │ - 29종 분류 │   │ - 이상 점수 │   │ - 7일 예측  │            │   │
│  │ │ - 신뢰도    │   │ - 이상 유형 │   │ - 위험 확률 │            │   │
│  │ └─────────────┘   └─────────────┘   └─────────────┘            │   │
│  │                              ↓                                   │   │
│  │ ┌─────────────────────────────────────────────────────────────┐ │   │
│  │ │ Multi-Task Fusion: 3개 모델 결과 통합 + 건강 점수 산출       │ │   │
│  │ │ health_score = w1*classification + w2*anomaly + w3*trend    │ │   │
│  │ └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              ↓                                          │
│  Stage 4: 개인화 코칭 (Personalized Coaching)                            │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ ┌─────────────────────────────────────────────────────────────┐ │   │
│  │ │ Context Aggregation                                          │ │   │
│  │ │ - 현재 측정 결과                                              │ │   │
│  │ │ - 과거 30일 트렌드                                            │ │   │
│  │ │ - 사용자 프로필 (나이, 성별, 기저질환)                         │ │   │
│  │ │ - 환경 데이터 (대기질, 날씨)                                   │ │   │
│  │ │ - 목표 설정 (체중 감량, 혈당 관리 등)                          │ │   │
│  │ └─────────────────────────────────────────────────────────────┘ │   │
│  │                              ↓                                   │   │
│  │ ┌─────────────────────────────────────────────────────────────┐ │   │
│  │ │ Recommendation Engine (협업 필터링 + 콘텐츠 기반)             │ │   │
│  │ │                                                              │ │   │
│  │ │ Output:                                                      │ │   │
│  │ │ - 식단 추천 (칼로리, 영양소 균형)                             │ │   │
│  │ │ - 운동 추천 (유형, 강도, 시간)                                │ │   │
│  │ │ - 수면 조언 (취침 시간, 수면 환경)                            │ │   │
│  │ │ - 스트레스 관리 (호흡법, 명상)                                │ │   │
│  │ │ - 의료 상담 권고 (임계값 초과 시)                             │ │   │
│  │ └─────────────────────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                              ↓                                          │
│  Stage 5: 사용자 전달 (Delivery)                                         │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐         │   │
│  │ │ Push    │   │ In-App  │   │ Voice   │   │ Email   │         │   │
│  │ │ 알림    │   │ 카드    │   │ TTS     │   │ 리포트  │         │   │
│  │ └─────────┘   └─────────┘   └─────────┘   └─────────┘         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.3 연합학습 (Federated Learning) 구현 계획

#### 아키텍처 (Flower Framework 기반)

```python
# Flower 기반 연합학습 아키텍처
"""
참조: https://www.nature.com/articles/s41598-025-04083-4
     https://pmc.ncbi.nlm.nih.gov/articles/PMC11728217/
"""

class ManpasikFederatedClient(fl.client.NumPyClient):
    """
    리더기/앱에서 실행되는 로컬 학습 클라이언트
    - 사용자 데이터는 절대 서버로 전송되지 않음
    - 모델 가중치만 암호화하여 전송
    """

    def __init__(self, model, local_data, privacy_config):
        self.model = model
        self.local_data = local_data  # 로컬에만 저장
        self.privacy = privacy_config

    def fit(self, parameters, config):
        # 1. 서버로부터 글로벌 모델 수신
        self.model.set_weights(parameters)

        # 2. 로컬 데이터로 학습 (데이터는 디바이스에 유지)
        self.model.fit(
            self.local_data.x,
            self.local_data.y,
            epochs=config["local_epochs"],
            batch_size=config["batch_size"]
        )

        # 3. 차등 프라이버시 적용 (노이즈 추가)
        updated_weights = self.model.get_weights()
        if self.privacy.differential_privacy_enabled:
            updated_weights = self._add_noise(
                updated_weights,
                epsilon=self.privacy.epsilon,  # 권장: 1.0-10.0
                delta=self.privacy.delta       # 권장: 1e-5
            )

        # 4. 암호화하여 서버로 전송 (데이터 아닌 가중치만)
        return updated_weights, len(self.local_data), {}

    def _add_noise(self, weights, epsilon, delta):
        """
        차등 프라이버시: 가중치에 가우시안 노이즈 추가
        - 개별 사용자 데이터 추론 불가능하게 함
        - 최대 30% 노이즈까지 모델 성능 유지 (연구 결과)
        """
        noise_scale = np.sqrt(2 * np.log(1.25 / delta)) / epsilon
        return [w + np.random.normal(0, noise_scale, w.shape) for w in weights]


class ManpasikFederatedServer:
    """
    중앙 서버: 모델 집계만 수행, 개별 데이터 접근 불가
    """

    def __init__(self, strategy_config):
        self.strategy = fl.server.strategy.FedAvg(
            fraction_fit=0.1,           # 라운드당 10% 클라이언트 참여
            fraction_evaluate=0.05,     # 5% 클라이언트로 평가
            min_fit_clients=10,         # 최소 10개 클라이언트 필요
            min_evaluate_clients=5,
            min_available_clients=50,
            on_fit_config_fn=self._fit_config,
            on_evaluate_config_fn=self._evaluate_config,
        )

    def aggregate(self, results):
        """
        FedAvg: 가중 평균으로 글로벌 모델 업데이트
        - 각 클라이언트의 데이터 크기에 비례하여 가중치 부여
        - Secure Aggregation으로 개별 가중치 노출 방지
        """
        total_samples = sum([r.num_examples for r in results])
        weighted_weights = []

        for result in results:
            weight = result.num_examples / total_samples
            weighted_weights.append(
                [w * weight for w in result.parameters]
            )

        # 집계된 글로벌 모델
        global_weights = [
            sum(layer_weights)
            for layer_weights in zip(*weighted_weights)
        ]

        return global_weights
```

### 3.4 예측적 UX (Predictive UX) 설계

[Healthcare UX 2026 트렌드](https://www.uxstudioteam.com/ux-blog/healthcare-ux) 기반 설계:

```yaml
# 예측적 UX 시나리오
scenarios:
  morning_prediction:
    trigger: 사용자 기상 시간 (학습된 패턴)
    actions:
      - 어젯밤 수면 분석 결과 준비
      - 오늘의 건강 요약 생성
      - 아침 측정 리마인더 스케줄링
      - 날씨 기반 운동 추천 준비

  measurement_anticipation:
    trigger: 측정 시작 버튼 탭
    actions:
      - 이전 측정 컨텍스트 로드 (마지막 사용 카트리지)
      - BLE 연결 사전 시도
      - AI 모델 웜업 (첫 추론 지연 최소화)
      - 결과 화면 템플릿 프리렌더링

  anomaly_response:
    trigger: 이상 수치 감지
    actions:
      - 즉시 시각적 피드백 (빨간 하이라이트)
      - 관련 과거 데이터 자동 로드
      - 의료 상담 예약 버튼 표시
      - 보호자 알림 준비 (설정된 경우)
      - 긴급 연락망 접근성 향상

  engagement_optimization:
    trigger: 3일 이상 측정 미수행
    actions:
      - 동기부여 메시지 개인화
      - 스트릭(연속 기록) 복구 기회 제공
      - 친구/가족 챌린지 제안
      - 보상(포인트) 증가 알림
```

---

## 4. Rust 코어 세부 구현 기획

### 4.1 AI 모듈 TFLite 실제 구현

#### 배경 조사
- [TFLite Micro](https://github.com/tensorflow/tflite-micro): 마이크로컨트롤러용 경량 추론 엔진
- [tflitec 0.7](https://crates.io/crates/tflitec): Rust TFLite C 바인딩
- [Seeed Studio TFLite Guide](https://wiki.seeedstudio.com/XIAO-BLE-Sense-TFLite-Getting-Started/): nRF52840 + TFLite 구현 사례

#### 구현 상세

```rust
// rust-core/manpasik-engine/src/ai/mod.rs

use tflitec::interpreter::{Interpreter, Options};
use tflitec::tensor::Tensor;
use std::path::Path;

/// TFLite 기반 실제 AI 추론 엔진
///
/// # 모델 구성
/// - biomarker_classifier.tflite: 29종 바이오마커 분류 (896D → 29 classes)
/// - anomaly_detector.tflite: 이상치 탐지 (896D → anomaly_score)
/// - trend_predictor.tflite: 7일 예측 (시계열 → 7D forecast)
///
/// # 성능 목표
/// - 추론 시간: < 50ms (ARM Cortex-M4 기준)
/// - 메모리 사용: < 512KB
/// - 정확도: 92-96%
pub struct TFLiteInferenceEngine {
    classifier: Option<Interpreter>,
    anomaly_detector: Option<Interpreter>,
    trend_predictor: Option<Interpreter>,
    model_paths: ModelPaths,
    warmup_done: bool,
}

#[derive(Clone)]
pub struct ModelPaths {
    pub classifier: PathBuf,
    pub anomaly_detector: PathBuf,
    pub trend_predictor: PathBuf,
}

impl TFLiteInferenceEngine {
    /// 모델 로드 (지연 로딩 지원)
    pub fn new(model_paths: ModelPaths) -> Result<Self, InferenceError> {
        Ok(Self {
            classifier: None,
            anomaly_detector: None,
            trend_predictor: None,
            model_paths,
            warmup_done: false,
        })
    }

    /// 분류 모델 로드 및 초기화
    pub fn load_classifier(&mut self) -> Result<(), InferenceError> {
        let model_data = std::fs::read(&self.model_paths.classifier)
            .map_err(|e| InferenceError::ModelLoadFailed(e.to_string()))?;

        let options = Options::default();
        // 스레드 수 설정 (임베디드: 1, 모바일: 2-4)
        options.set_num_threads(2);

        let interpreter = Interpreter::new(&model_data, Some(options))
            .map_err(|e| InferenceError::InterpreterFailed(e.to_string()))?;

        // 텐서 할당
        interpreter.allocate_tensors()
            .map_err(|e| InferenceError::TensorAllocationFailed(e.to_string()))?;

        self.classifier = Some(interpreter);
        Ok(())
    }

    /// 바이오마커 분류 추론
    ///
    /// # Arguments
    /// * `fingerprint` - 896차원 핑거프린트 벡터
    ///
    /// # Returns
    /// * `ClassificationResult` - 29종 분류 결과 + 신뢰도
    pub fn classify_biomarkers(
        &self,
        fingerprint: &[f32; 896]
    ) -> Result<ClassificationResult, InferenceError> {
        let interpreter = self.classifier.as_ref()
            .ok_or(InferenceError::ModelNotLoaded)?;

        // 입력 텐서 설정
        let input_tensor = interpreter.input(0)
            .map_err(|e| InferenceError::TensorAccessFailed(e.to_string()))?;

        // 데이터 복사 (896 floats)
        input_tensor.copy_from_slice(fingerprint)
            .map_err(|e| InferenceError::DataCopyFailed(e.to_string()))?;

        // 추론 실행
        let start = std::time::Instant::now();
        interpreter.invoke()
            .map_err(|e| InferenceError::InferenceFailed(e.to_string()))?;
        let inference_time = start.elapsed();

        // 출력 텐서 읽기 (29 classes softmax)
        let output_tensor = interpreter.output(0)
            .map_err(|e| InferenceError::TensorAccessFailed(e.to_string()))?;

        let probabilities: Vec<f32> = output_tensor.data().to_vec();

        // argmax로 최고 확률 클래스 찾기
        let (predicted_class, confidence) = probabilities
            .iter()
            .enumerate()
            .max_by(|(_, a), (_, b)| a.partial_cmp(b).unwrap())
            .map(|(idx, &prob)| (idx as u8, prob))
            .unwrap_or((0, 0.0));

        Ok(ClassificationResult {
            predicted_class,
            confidence,
            all_probabilities: probabilities,
            inference_time_ms: inference_time.as_millis() as u32,
        })
    }

    /// 이상치 탐지
    pub fn detect_anomaly(
        &self,
        fingerprint: &[f32; 896],
        historical_data: &[f32],  // 과거 30일 데이터
    ) -> Result<AnomalyResult, InferenceError> {
        let interpreter = self.anomaly_detector.as_ref()
            .ok_or(InferenceError::ModelNotLoaded)?;

        // 입력: 현재 핑거프린트 + 히스토리 통계
        let mut input = Vec::with_capacity(896 + 10);  // 896D + 10 stats
        input.extend_from_slice(fingerprint);
        input.extend(self.compute_statistics(historical_data));

        let input_tensor = interpreter.input(0)?;
        input_tensor.copy_from_slice(&input)?;

        interpreter.invoke()?;

        let output = interpreter.output(0)?;
        let anomaly_score = output.data::<f32>()[0];

        Ok(AnomalyResult {
            score: anomaly_score,
            is_anomaly: anomaly_score > 0.7,  // 임계값
            anomaly_type: self.classify_anomaly_type(anomaly_score),
            recommendation: self.generate_anomaly_recommendation(anomaly_score),
        })
    }

    /// 트렌드 예측 (7일)
    pub fn predict_trend(
        &self,
        time_series: &[TimeSeriesPoint],  // 과거 30일
    ) -> Result<TrendPrediction, InferenceError> {
        let interpreter = self.trend_predictor.as_ref()
            .ok_or(InferenceError::ModelNotLoaded)?;

        // LSTM 입력 형태: [batch=1, seq_len=30, features=896]
        let input_data: Vec<f32> = time_series
            .iter()
            .flat_map(|p| p.fingerprint.iter().copied())
            .collect();

        let input_tensor = interpreter.input(0)?;
        input_tensor.copy_from_slice(&input_data)?;

        interpreter.invoke()?;

        // 출력: 7일 예측 + 신뢰 구간
        let output = interpreter.output(0)?;
        let predictions: Vec<f32> = output.data().to_vec();

        Ok(TrendPrediction {
            daily_predictions: predictions[..7].to_vec(),
            confidence_lower: predictions[7..14].to_vec(),
            confidence_upper: predictions[14..21].to_vec(),
            trend_direction: self.determine_trend(&predictions[..7]),
            risk_probability: predictions[21],  // 위험 확률
        })
    }

    /// 모델 웜업 (첫 추론 지연 최소화)
    pub fn warmup(&mut self) -> Result<(), InferenceError> {
        if self.warmup_done {
            return Ok(());
        }

        // 더미 데이터로 각 모델 1회 추론
        let dummy_fingerprint = [0.0f32; 896];

        if self.classifier.is_some() {
            let _ = self.classify_biomarkers(&dummy_fingerprint);
        }

        self.warmup_done = true;
        Ok(())
    }
}

#[derive(Debug, Clone)]
pub struct ClassificationResult {
    pub predicted_class: u8,           // 0-28 (29종)
    pub confidence: f32,               // 0.0-1.0
    pub all_probabilities: Vec<f32>,   // 29개 확률
    pub inference_time_ms: u32,
}

#[derive(Debug, Clone)]
pub struct AnomalyResult {
    pub score: f32,                    // 0.0-1.0 (높을수록 이상)
    pub is_anomaly: bool,
    pub anomaly_type: AnomalyType,
    pub recommendation: String,
}

#[derive(Debug, Clone)]
pub enum AnomalyType {
    Normal,
    SlightDeviation,
    SignificantAnomaly,
    CriticalAlert,
}

#[derive(Debug, Clone)]
pub struct TrendPrediction {
    pub daily_predictions: Vec<f32>,   // 7일 예측값
    pub confidence_lower: Vec<f32>,    // 95% 하한
    pub confidence_upper: Vec<f32>,    // 95% 상한
    pub trend_direction: TrendDirection,
    pub risk_probability: f32,         // 위험 발생 확률
}

#[derive(Debug, Clone)]
pub enum TrendDirection {
    Improving,
    Stable,
    Declining,
    Volatile,
}
```

### 4.2 BLE 모듈 btleplug 실제 구현

#### 배경 조사
- [TrouBLE (Embassy-rs)](https://github.com/embassy-rs/trouble): Rust BLE Host 스택
- [btleplug](https://github.com/deviceplug/btleplug): Cross-platform Rust BLE
- [219 Design BLE Guide](https://www.219design.com/bluetooth-low-energy-with-rust/): Rust BLE 구현 사례
- [Punch Through nRF52840](https://punchthrough.com/nordic-nrf52840-is-rust-a-good-fit-for-embedded-applications/): 임베디드 Rust BLE 적합성

#### 구현 상세

```rust
// rust-core/manpasik-engine/src/ble/mod.rs

use btleplug::api::{
    Central, Manager as _, Peripheral as _, ScanFilter,
    WriteType, CharPropFlags
};
use btleplug::platform::{Manager, Peripheral};
use futures::stream::StreamExt;
use tokio::sync::mpsc;
use uuid::Uuid;

/// ManPaSik BLE 서비스 UUID
const MANPASIK_SERVICE_UUID: Uuid =
    Uuid::from_u128(0x12345678_1234_5678_1234_567812345678);

/// 측정 데이터 Characteristic UUID
const MEASUREMENT_CHAR_UUID: Uuid =
    Uuid::from_u128(0x12345678_1234_5678_1234_567812345679);

/// 명령 Characteristic UUID (쓰기용)
const COMMAND_CHAR_UUID: Uuid =
    Uuid::from_u128(0x12345678_1234_5678_1234_56781234567A);

/// 실제 BLE 통신 관리자
///
/// # 기능
/// - 리더기 스캔 및 자동 연결
/// - 측정 데이터 실시간 수신 (Notification)
/// - 명령 전송 (측정 시작/중지, 보정 등)
/// - 다중 리더기 동시 관리 (구독 등급별 제한)
pub struct BleManager {
    manager: Manager,
    connected_devices: HashMap<String, ConnectedDevice>,
    event_tx: mpsc::Sender<BleEvent>,
    config: BleConfig,
}

pub struct ConnectedDevice {
    peripheral: Peripheral,
    device_info: DeviceInfo,
    measurement_char: Option<btleplug::api::Characteristic>,
    command_char: Option<btleplug::api::Characteristic>,
    connection_state: ConnectionState,
}

#[derive(Clone)]
pub struct BleConfig {
    pub scan_timeout: Duration,
    pub connection_timeout: Duration,
    pub max_concurrent_devices: usize,  // 구독 등급별
    pub auto_reconnect: bool,
    pub rssi_threshold: i16,  // 신호 강도 최소값
}

impl BleManager {
    /// BLE 관리자 초기화
    pub async fn new(
        config: BleConfig,
        event_tx: mpsc::Sender<BleEvent>,
    ) -> Result<Self, BleError> {
        let manager = Manager::new().await
            .map_err(|e| BleError::InitializationFailed(e.to_string()))?;

        Ok(Self {
            manager,
            connected_devices: HashMap::new(),
            event_tx,
            config,
        })
    }

    /// ManPaSik 리더기 스캔
    ///
    /// # Returns
    /// 발견된 리더기 목록 (RSSI 순 정렬)
    pub async fn scan_devices(&self) -> Result<Vec<ScannedDevice>, BleError> {
        let adapters = self.manager.adapters().await
            .map_err(|e| BleError::AdapterNotFound(e.to_string()))?;

        let adapter = adapters.into_iter().next()
            .ok_or(BleError::NoAdapterAvailable)?;

        // ManPaSik 서비스 UUID 필터링
        let filter = ScanFilter {
            services: vec![MANPASIK_SERVICE_UUID],
        };

        adapter.start_scan(filter).await
            .map_err(|e| BleError::ScanFailed(e.to_string()))?;

        // 스캔 타임아웃
        tokio::time::sleep(self.config.scan_timeout).await;

        adapter.stop_scan().await?;

        // 발견된 디바이스 수집
        let peripherals = adapter.peripherals().await?;
        let mut devices = Vec::new();

        for peripheral in peripherals {
            if let Some(props) = peripheral.properties().await? {
                if props.rssi.unwrap_or(-100) >= self.config.rssi_threshold {
                    devices.push(ScannedDevice {
                        id: peripheral.id().to_string(),
                        name: props.local_name.unwrap_or_default(),
                        rssi: props.rssi.unwrap_or(-100),
                        manufacturer_data: props.manufacturer_data,
                    });
                }
            }
        }

        // RSSI 순 정렬 (강한 신호 우선)
        devices.sort_by(|a, b| b.rssi.cmp(&a.rssi));

        Ok(devices)
    }

    /// 리더기 연결
    pub async fn connect(&mut self, device_id: &str) -> Result<DeviceInfo, BleError> {
        // 최대 연결 수 체크
        if self.connected_devices.len() >= self.config.max_concurrent_devices {
            return Err(BleError::MaxDevicesExceeded(self.config.max_concurrent_devices));
        }

        let adapters = self.manager.adapters().await?;
        let adapter = adapters.into_iter().next()
            .ok_or(BleError::NoAdapterAvailable)?;

        let peripherals = adapter.peripherals().await?;
        let peripheral = peripherals.into_iter()
            .find(|p| p.id().to_string() == device_id)
            .ok_or(BleError::DeviceNotFound(device_id.to_string()))?;

        // 연결 시도 (타임아웃 적용)
        tokio::time::timeout(
            self.config.connection_timeout,
            peripheral.connect()
        ).await
            .map_err(|_| BleError::ConnectionTimeout)?
            .map_err(|e| BleError::ConnectionFailed(e.to_string()))?;

        // 서비스 검색
        peripheral.discover_services().await?;

        // ManPaSik 서비스에서 Characteristic 찾기
        let mut measurement_char = None;
        let mut command_char = None;

        for service in peripheral.services() {
            if service.uuid == MANPASIK_SERVICE_UUID {
                for char in service.characteristics {
                    if char.uuid == MEASUREMENT_CHAR_UUID {
                        measurement_char = Some(char.clone());
                    } else if char.uuid == COMMAND_CHAR_UUID {
                        command_char = Some(char.clone());
                    }
                }
            }
        }

        // Notification 구독 (측정 데이터 수신)
        if let Some(ref char) = measurement_char {
            if char.properties.contains(CharPropFlags::NOTIFY) {
                peripheral.subscribe(char).await?;
            }
        }

        // 디바이스 정보 읽기
        let device_info = self.read_device_info(&peripheral).await?;

        // 연결된 디바이스 저장
        self.connected_devices.insert(device_id.to_string(), ConnectedDevice {
            peripheral,
            device_info: device_info.clone(),
            measurement_char,
            command_char,
            connection_state: ConnectionState::Connected,
        });

        // 이벤트 발행
        self.event_tx.send(BleEvent::DeviceConnected(device_info.clone())).await?;

        Ok(device_info)
    }

    /// 측정 시작 명령 전송
    pub async fn start_measurement(
        &self,
        device_id: &str,
        params: MeasurementParams,
    ) -> Result<(), BleError> {
        let device = self.connected_devices.get(device_id)
            .ok_or(BleError::DeviceNotConnected(device_id.to_string()))?;

        let command_char = device.command_char.as_ref()
            .ok_or(BleError::CharacteristicNotFound)?;

        // 명령 패킷 구성
        let command = MeasurementCommand::Start(params);
        let packet = command.to_bytes();

        // 명령 전송
        device.peripheral.write(
            command_char,
            &packet,
            WriteType::WithResponse
        ).await?;

        self.event_tx.send(BleEvent::MeasurementStarted(device_id.to_string())).await?;

        Ok(())
    }

    /// 측정 데이터 스트림 구독
    pub async fn subscribe_measurement_stream(
        &self,
        device_id: &str,
    ) -> Result<impl futures::Stream<Item = MeasurementDataPacket>, BleError> {
        let device = self.connected_devices.get(device_id)
            .ok_or(BleError::DeviceNotConnected(device_id.to_string()))?;

        let notification_stream = device.peripheral.notifications().await?;

        // 측정 데이터만 필터링하여 반환
        Ok(notification_stream.filter_map(|notification| async move {
            if notification.uuid == MEASUREMENT_CHAR_UUID {
                MeasurementDataPacket::from_bytes(&notification.value).ok()
            } else {
                None
            }
        }))
    }

    /// 자동 재연결 (연결 끊김 감지 시)
    pub async fn handle_disconnection(&mut self, device_id: &str) {
        if !self.config.auto_reconnect {
            return;
        }

        // 재연결 시도 (최대 3회, 지수 백오프)
        for attempt in 0..3 {
            let delay = Duration::from_millis(100 * 2u64.pow(attempt));
            tokio::time::sleep(delay).await;

            if self.connect(device_id).await.is_ok() {
                self.event_tx.send(BleEvent::DeviceReconnected(device_id.to_string())).await.ok();
                return;
            }
        }

        // 재연결 실패
        self.event_tx.send(BleEvent::ReconnectionFailed(device_id.to_string())).await.ok();
    }
}

#[derive(Debug, Clone)]
pub struct ScannedDevice {
    pub id: String,
    pub name: String,
    pub rssi: i16,
    pub manufacturer_data: HashMap<u16, Vec<u8>>,
}

#[derive(Debug, Clone)]
pub struct DeviceInfo {
    pub id: String,
    pub name: String,
    pub firmware_version: String,
    pub hardware_version: String,
    pub serial_number: String,
    pub battery_level: u8,
    pub last_calibration: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone)]
pub enum BleEvent {
    DeviceDiscovered(ScannedDevice),
    DeviceConnected(DeviceInfo),
    DeviceDisconnected(String),
    DeviceReconnected(String),
    ReconnectionFailed(String),
    MeasurementStarted(String),
    MeasurementDataReceived(MeasurementDataPacket),
    MeasurementCompleted(String),
    BatteryLow(String, u8),
    Error(BleError),
}

#[derive(Debug, Clone)]
pub struct MeasurementDataPacket {
    pub timestamp: u64,
    pub channel_data: [f32; 88],  // 88채널 원시 데이터
    pub reference_data: [f32; 88],  // 참조 데이터
    pub temperature: f32,
    pub humidity: f32,
    pub sequence_number: u32,
}
```

### 4.3 NFC 모듈 실제 구현

```rust
// rust-core/manpasik-engine/src/nfc/reader.rs

use std::time::Duration;

/// NFC 리더 추상화 (플랫폼별 구현)
///
/// # 지원 플랫폼
/// - iOS: CoreNFC (NfcTagReaderSession)
/// - Android: android.nfc (NfcA, IsoDep)
/// - Linux: libnfc (ACR122U 등)
/// - Embedded: PN532 (SPI/I2C)
pub trait NfcReader: Send + Sync {
    /// NFC 태그 폴링 시작
    fn start_polling(&mut self) -> Result<(), NfcError>;

    /// NFC 태그 폴링 중지
    fn stop_polling(&mut self) -> Result<(), NfcError>;

    /// 태그 UID 읽기
    fn read_uid(&self) -> Result<[u8; 7], NfcError>;

    /// NDEF 데이터 읽기
    fn read_ndef(&self) -> Result<Vec<u8>, NfcError>;

    /// 특정 블록 읽기 (MIFARE)
    fn read_block(&self, block: u8) -> Result<[u8; 16], NfcError>;

    /// 특정 블록 쓰기 (MIFARE)
    fn write_block(&self, block: u8, data: &[u8; 16]) -> Result<(), NfcError>;

    /// ISO 14443-4 APDU 명령
    fn transceive(&self, command: &[u8]) -> Result<Vec<u8>, NfcError>;
}

/// ManPaSik 카트리지 NFC 매니저
///
/// # 태그 포맷 (v2.0)
/// ```
/// Block 0: UID (7 bytes)
/// Block 1: Category (1) + Type (1) + Legacy Code (1) + Reserved (1)
/// Block 2-3: Calibration Data (32 bytes)
/// Block 4-5: Manufacturing Data (32 bytes)
/// Block 6: Usage Counter (4 bytes) + Expiry (4 bytes) + CRC (4 bytes)
/// ```
pub struct CartridgeNfcManager<R: NfcReader> {
    reader: R,
    registry: CartridgeRegistry,
    cache: HashMap<[u8; 7], CachedCartridge>,
}

impl<R: NfcReader> CartridgeNfcManager<R> {
    pub fn new(reader: R, registry: CartridgeRegistry) -> Self {
        Self {
            reader,
            registry,
            cache: HashMap::new(),
        }
    }

    /// 카트리지 자동 인식
    ///
    /// NFC 태그 감지 시 자동으로 호출
    /// - UID 읽기 → 캐시 확인
    /// - v2.0 포맷 파싱
    /// - v1.0 레거시 자동 변환
    /// - 유효성 검증 (만료일, 사용 횟수)
    pub async fn read_cartridge(&mut self) -> Result<CartridgeInfo, NfcError> {
        // 1. UID 읽기
        let uid = self.reader.read_uid()?;

        // 2. 캐시 확인 (성능 최적화)
        if let Some(cached) = self.cache.get(&uid) {
            if cached.is_valid() {
                return Ok(cached.info.clone());
            }
        }

        // 3. 헤더 블록 읽기 (Category, Type, Legacy Code)
        let header = self.reader.read_block(1)?;
        let category = header[0];
        let type_code = header[1];
        let legacy_code = header[2];

        // 4. v1.0 레거시 변환 (legacy_code만 있는 경우)
        let full_code = if category == 0x00 && type_code == 0x00 {
            self.convert_legacy_code(legacy_code)?
        } else {
            CartridgeFullCode::new(category, type_code)
        };

        // 5. 레지스트리에서 카트리지 정보 조회
        let registry_info = self.registry.get_cartridge_info(&full_code)
            .ok_or(NfcError::UnknownCartridge(full_code))?;

        // 6. 보정 데이터 읽기
        let calibration_block_2 = self.reader.read_block(2)?;
        let calibration_block_3 = self.reader.read_block(3)?;
        let calibration_data = [calibration_block_2, calibration_block_3].concat();

        // 7. 사용량 데이터 읽기
        let usage_block = self.reader.read_block(6)?;
        let usage_count = u32::from_le_bytes(usage_block[0..4].try_into().unwrap());
        let expiry_timestamp = u32::from_le_bytes(usage_block[4..8].try_into().unwrap());
        let stored_crc = u32::from_le_bytes(usage_block[8..12].try_into().unwrap());

        // 8. CRC 검증
        let computed_crc = self.compute_crc(&[&header[..], &calibration_data, &usage_block[0..8]]);
        if computed_crc != stored_crc {
            return Err(NfcError::CrcMismatch);
        }

        // 9. 유효성 검증
        let now = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_secs() as u32;

        if now > expiry_timestamp {
            return Err(NfcError::CartridgeExpired);
        }

        if usage_count >= registry_info.max_usage_count {
            return Err(NfcError::UsageLimitExceeded);
        }

        // 10. 보정 데이터 파싱
        let calibration = CalibrationCoefficients::from_bytes(&calibration_data)?;

        let cartridge_info = CartridgeInfo {
            uid,
            full_code,
            name: registry_info.name.clone(),
            category_name: registry_info.category_name.clone(),
            calibration,
            usage_count,
            remaining_uses: registry_info.max_usage_count - usage_count,
            expiry_date: DateTime::from_timestamp(expiry_timestamp as i64, 0),
            access_tier: registry_info.access_tier,
        };

        // 11. 캐시 업데이트
        self.cache.insert(uid, CachedCartridge {
            info: cartridge_info.clone(),
            cached_at: std::time::Instant::now(),
        });

        Ok(cartridge_info)
    }

    /// 사용 카운터 증가 (측정 완료 후 호출)
    pub async fn increment_usage(&mut self, uid: &[u8; 7]) -> Result<(), NfcError> {
        // 현재 사용량 읽기
        let usage_block = self.reader.read_block(6)?;
        let mut usage_count = u32::from_le_bytes(usage_block[0..4].try_into().unwrap());

        // 증가
        usage_count += 1;

        // 새 블록 데이터 구성
        let mut new_block = [0u8; 16];
        new_block[0..4].copy_from_slice(&usage_count.to_le_bytes());
        new_block[4..8].copy_from_slice(&usage_block[4..8]);  // expiry 유지

        // CRC 재계산
        let header = self.reader.read_block(1)?;
        let calibration_block_2 = self.reader.read_block(2)?;
        let calibration_block_3 = self.reader.read_block(3)?;
        let calibration_data = [calibration_block_2, calibration_block_3].concat();

        let new_crc = self.compute_crc(&[&header[..], &calibration_data, &new_block[0..8]]);
        new_block[8..12].copy_from_slice(&new_crc.to_le_bytes());

        // 쓰기
        self.reader.write_block(6, &new_block)?;

        // 캐시 무효화
        self.cache.remove(uid);

        Ok(())
    }

    /// v1.0 레거시 코드를 v2.0 풀 코드로 변환
    fn convert_legacy_code(&self, legacy: u8) -> Result<CartridgeFullCode, NfcError> {
        // 29종 레거시 매핑 테이블
        match legacy {
            0x01 => Ok(CartridgeFullCode::new(0x01, 0x01)),  // Glucose
            0x02 => Ok(CartridgeFullCode::new(0x01, 0x02)),  // Cholesterol
            0x03 => Ok(CartridgeFullCode::new(0x01, 0x03)),  // Triglyceride
            // ... 나머지 26종
            0x1D => Ok(CartridgeFullCode::new(0x01, 0x1D)),  // HbA1c
            _ => Err(NfcError::UnknownLegacyCode(legacy)),
        }
    }

    fn compute_crc(&self, data: &[&[u8]]) -> u32 {
        let mut hasher = crc32fast::Hasher::new();
        for chunk in data {
            hasher.update(chunk);
        }
        hasher.finalize()
    }
}

#[derive(Debug, Clone)]
pub struct CartridgeInfo {
    pub uid: [u8; 7],
    pub full_code: CartridgeFullCode,
    pub name: String,
    pub category_name: String,
    pub calibration: CalibrationCoefficients,
    pub usage_count: u32,
    pub remaining_uses: u32,
    pub expiry_date: Option<DateTime<Utc>>,
    pub access_tier: AccessTier,
}

#[derive(Debug, Clone)]
pub struct CalibrationCoefficients {
    pub alpha: f32,              // 차동측정 계수 (기본 0.95)
    pub temperature_comp: f32,   // 온도 보정 계수
    pub humidity_comp: f32,      // 습도 보정 계수
    pub wavelength_offsets: [f32; 8],  // 파장별 오프셋
    pub sensitivity_matrix: [[f32; 4]; 4],  // 4x4 감도 매트릭스
}
```

---

## 5. Flutter Feature 세부 구현 기획

### 5.1 Healthcare UX 2026 트렌드 적용

[UX Studio Healthcare Trends 2026](https://www.uxstudioteam.com/ux-blog/healthcare-ux) 및 [Eleken Healthcare UI](https://www.eleken.co/blog-posts/user-interface-design-for-healthcare-applications) 기반 설계:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Flutter Feature 설계 원칙                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. 예측적 UX (Predictive UX)                                            │
│     - 사용자 행동 패턴 학습 → 다음 액션 예측                              │
│     - 컨텍스트 기반 UI 적응 (시간, 위치, 건강 상태)                        │
│                                                                         │
│  2. 음성 우선 (Voice-First)                                              │
│     - 모든 핵심 기능 음성 명령 지원                                       │
│     - 핸즈프리 측정 플로우                                               │
│                                                                         │
│  3. 감정 인식 인터페이스 (Emotion-Aware)                                  │
│     - 건강 상태에 따른 UI 톤 조절                                         │
│     - 위험 상황 시 차분한 색상 + 명확한 안내                              │
│                                                                         │
│  4. 게이미피케이션 (Gamification)                                        │
│     - 스트릭, 배지, 레벨업                                               │
│     - 가족/친구 챌린지                                                   │
│                                                                         │
│  5. 접근성 최우선 (Accessibility-First)                                  │
│     - WCAG 2.1 AA 준수                                                  │
│     - 스크린 리더, 고대비, 큰 글씨                                        │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.2 data_hub Feature 상세 설계

#### 화면 구성

```dart
// lib/features/data_hub/presentation/screens/data_hub_screen.dart

/// 데이터 허브 메인 화면
///
/// # 구성요소
/// 1. 요약 대시보드 (오늘의 건강 점수)
/// 2. 타임라인 (시간순 측정 기록)
/// 3. 트렌드 차트 (기간별 분석)
/// 4. 필터/검색
/// 5. 데이터 내보내기
class DataHubScreen extends ConsumerStatefulWidget {
  @override
  ConsumerState<DataHubScreen> createState() => _DataHubScreenState();
}

class _DataHubScreenState extends ConsumerState<DataHubScreen>
    with SingleTickerProviderStateMixin {

  late TabController _tabController;
  DateTimeRange _selectedRange = DateTimeRange(
    start: DateTime.now().subtract(Duration(days: 30)),
    end: DateTime.now(),
  );

  @override
  Widget build(BuildContext context) {
    final healthSummary = ref.watch(healthSummaryProvider);
    final timeline = ref.watch(measurementTimelineProvider(_selectedRange));

    return Scaffold(
      appBar: AppBar(
        title: Text('데이터 허브'),
        actions: [
          // 필터 버튼
          IconButton(
            icon: Icon(Icons.filter_list),
            onPressed: _showFilterBottomSheet,
          ),
          // 내보내기 버튼
          IconButton(
            icon: Icon(Icons.download),
            onPressed: _showExportDialog,
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          tabs: [
            Tab(text: '요약'),
            Tab(text: '타임라인'),
            Tab(text: '트렌드'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          // 요약 탭
          _SummaryTab(summary: healthSummary),

          // 타임라인 탭
          _TimelineTab(timeline: timeline),

          // 트렌드 탭
          _TrendTab(range: _selectedRange),
        ],
      ),
    );
  }
}

/// 건강 요약 탭
class _SummaryTab extends StatelessWidget {
  final AsyncValue<HealthSummary> summary;

  @override
  Widget build(BuildContext context) {
    return summary.when(
      data: (data) => SingleChildScrollView(
        padding: EdgeInsets.all(16),
        child: Column(
          children: [
            // 건강 점수 카드 (원형 게이지)
            HealthScoreCard(
              score: data.overallScore,
              trend: data.scoreTrend,
              lastUpdated: data.lastMeasurement,
            ),
            SizedBox(height: 16),

            // 카테고리별 요약 (그리드)
            GridView.count(
              crossAxisCount: 2,
              shrinkWrap: true,
              physics: NeverScrollableScrollPhysics(),
              children: [
                CategorySummaryCard(
                  title: '바이오마커',
                  icon: Icons.science,
                  value: data.biomarkerSummary,
                  status: data.biomarkerStatus,
                ),
                CategorySummaryCard(
                  title: '영양 상태',
                  icon: Icons.restaurant,
                  value: data.nutritionSummary,
                  status: data.nutritionStatus,
                ),
                CategorySummaryCard(
                  title: '환경',
                  icon: Icons.eco,
                  value: data.environmentSummary,
                  status: data.environmentStatus,
                ),
                CategorySummaryCard(
                  title: '활동',
                  icon: Icons.directions_run,
                  value: data.activitySummary,
                  status: data.activityStatus,
                ),
              ],
            ),
            SizedBox(height: 16),

            // AI 인사이트
            AiInsightCard(
              insights: data.aiInsights,
              onTap: () => context.push('/coach'),
            ),
          ],
        ),
      ),
      loading: () => Center(child: CircularProgressIndicator()),
      error: (e, st) => ErrorWidget(error: e, onRetry: () => ref.refresh(healthSummaryProvider)),
    );
  }
}

/// 타임라인 탭
class _TimelineTab extends StatelessWidget {
  final AsyncValue<List<MeasurementRecord>> timeline;

  @override
  Widget build(BuildContext context) {
    return timeline.when(
      data: (records) => records.isEmpty
          ? EmptyStateWidget(
              icon: Icons.timeline,
              title: '측정 기록이 없습니다',
              subtitle: '첫 측정을 시작해보세요!',
              action: ElevatedButton(
                onPressed: () => context.push('/measure'),
                child: Text('측정 시작'),
              ),
            )
          : ListView.builder(
              itemCount: records.length,
              itemBuilder: (context, index) {
                final record = records[index];
                final showDateHeader = index == 0 ||
                    !_isSameDay(record.timestamp, records[index - 1].timestamp);

                return Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (showDateHeader) DateHeader(date: record.timestamp),
                    MeasurementTimelineCard(
                      record: record,
                      onTap: () => context.push('/measurement/${record.id}'),
                    ),
                  ],
                );
              },
            ),
      loading: () => TimelineShimmer(),
      error: (e, st) => ErrorWidget(error: e),
    );
  }
}

/// 트렌드 탭 (fl_chart 사용)
class _TrendTab extends StatelessWidget {
  final DateTimeRange range;

  @override
  Widget build(BuildContext context) {
    return Consumer(
      builder: (context, ref, child) {
        final trendData = ref.watch(trendDataProvider(range));

        return trendData.when(
          data: (data) => SingleChildScrollView(
            padding: EdgeInsets.all(16),
            child: Column(
              children: [
                // 기간 선택기
                DateRangeSelector(
                  selectedRange: range,
                  onChanged: (newRange) => ref.read(selectedRangeProvider.notifier).state = newRange,
                ),
                SizedBox(height: 16),

                // 메인 트렌드 차트
                Container(
                  height: 300,
                  child: LineChart(
                    LineChartData(
                      gridData: FlGridData(show: true),
                      titlesData: FlTitlesData(
                        leftTitles: AxisTitles(
                          sideTitles: SideTitles(showTitles: true),
                        ),
                        bottomTitles: AxisTitles(
                          sideTitles: SideTitles(
                            showTitles: true,
                            getTitlesWidget: (value, meta) => Text(
                              DateFormat('MM/dd').format(
                                DateTime.fromMillisecondsSinceEpoch(value.toInt())
                              ),
                              style: TextStyle(fontSize: 10),
                            ),
                          ),
                        ),
                      ),
                      lineBarsData: [
                        // 측정값 라인
                        LineChartBarData(
                          spots: data.measurementSpots,
                          color: Theme.of(context).primaryColor,
                          barWidth: 2,
                          dotData: FlDotData(show: true),
                        ),
                        // 개인 기준선 (My Zone)
                        LineChartBarData(
                          spots: data.baselineSpots,
                          color: Colors.green.withOpacity(0.3),
                          barWidth: 0,
                          belowBarData: BarAreaData(
                            show: true,
                            color: Colors.green.withOpacity(0.1),
                          ),
                        ),
                        // AI 예측 라인 (점선)
                        LineChartBarData(
                          spots: data.predictionSpots,
                          color: Colors.orange,
                          barWidth: 2,
                          dashArray: [5, 5],
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: 16),

                // 통계 요약
                StatisticsSummaryCard(
                  average: data.average,
                  min: data.min,
                  max: data.max,
                  stdDev: data.standardDeviation,
                  trend: data.overallTrend,
                ),
              ],
            ),
          ),
          loading: () => TrendChartShimmer(),
          error: (e, st) => ErrorWidget(error: e),
        );
      },
    );
  }
}
```

### 5.3 ai_coach Feature 상세 설계

#### 대화형 AI 코칭 인터페이스

```dart
// lib/features/ai_coach/presentation/screens/ai_coach_screen.dart

/// AI 코치 메인 화면
///
/// # 주요 기능
/// 1. 대화형 AI 상담 (채팅 UI)
/// 2. 일일/주간 코칭 카드
/// 3. 목표 설정 및 추적
/// 4. 음성 입력 지원
class AiCoachScreen extends ConsumerStatefulWidget {
  @override
  ConsumerState<AiCoachScreen> createState() => _AiCoachScreenState();
}

class _AiCoachScreenState extends ConsumerState<AiCoachScreen> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();
  bool _isListening = false;

  @override
  Widget build(BuildContext context) {
    final chatHistory = ref.watch(coachChatHistoryProvider);
    final coachingCards = ref.watch(dailyCoachingCardsProvider);
    final goals = ref.watch(healthGoalsProvider);

    return Scaffold(
      appBar: AppBar(
        title: Row(
          children: [
            CircleAvatar(
              child: Icon(Icons.smart_toy),
              backgroundColor: Theme.of(context).primaryColor,
            ),
            SizedBox(width: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('AI 건강 코치'),
                Text(
                  '항상 곁에서 도와드려요',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ],
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: Icon(Icons.history),
            onPressed: () => context.push('/coach/history'),
            tooltip: '대화 기록',
          ),
          IconButton(
            icon: Icon(Icons.flag),
            onPressed: () => context.push('/coach/goals'),
            tooltip: '목표 설정',
          ),
        ],
      ),
      body: Column(
        children: [
          // 오늘의 코칭 카드 (가로 스크롤)
          if (coachingCards.hasValue && coachingCards.value!.isNotEmpty)
            Container(
              height: 120,
              child: ListView.builder(
                scrollDirection: Axis.horizontal,
                padding: EdgeInsets.all(12),
                itemCount: coachingCards.value!.length,
                itemBuilder: (context, index) {
                  final card = coachingCards.value![index];
                  return CoachingCard(
                    card: card,
                    onTap: () => _handleCardTap(card),
                    onDismiss: () => ref.read(coachingCardsProvider.notifier).dismiss(card.id),
                  );
                },
              ),
            ),

          Divider(height: 1),

          // 채팅 히스토리
          Expanded(
            child: chatHistory.when(
              data: (messages) => ListView.builder(
                controller: _scrollController,
                padding: EdgeInsets.all(16),
                itemCount: messages.length,
                itemBuilder: (context, index) {
                  final message = messages[index];
                  return ChatBubble(
                    message: message,
                    isUser: message.sender == MessageSender.user,
                    onActionTap: message.actions != null
                        ? (action) => _handleAction(action)
                        : null,
                  );
                },
              ),
              loading: () => Center(child: CircularProgressIndicator()),
              error: (e, st) => ErrorWidget(error: e),
            ),
          ),

          // 입력 영역
          SafeArea(
            child: Container(
              padding: EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Theme.of(context).cardColor,
                boxShadow: [
                  BoxShadow(
                    color: Colors.black12,
                    blurRadius: 4,
                    offset: Offset(0, -2),
                  ),
                ],
              ),
              child: Row(
                children: [
                  // 음성 입력 버튼
                  IconButton(
                    icon: Icon(
                      _isListening ? Icons.mic : Icons.mic_none,
                      color: _isListening ? Colors.red : null,
                    ),
                    onPressed: _toggleVoiceInput,
                    tooltip: '음성 입력',
                  ),

                  // 텍스트 입력
                  Expanded(
                    child: TextField(
                      controller: _messageController,
                      decoration: InputDecoration(
                        hintText: '무엇이든 물어보세요...',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(24),
                        ),
                        contentPadding: EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 8,
                        ),
                      ),
                      textInputAction: TextInputAction.send,
                      onSubmitted: _sendMessage,
                    ),
                  ),

                  SizedBox(width: 8),

                  // 전송 버튼
                  FloatingActionButton.small(
                    onPressed: () => _sendMessage(_messageController.text),
                    child: Icon(Icons.send),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _sendMessage(String text) async {
    if (text.trim().isEmpty) return;

    _messageController.clear();

    // 사용자 메시지 추가
    await ref.read(coachChatHistoryProvider.notifier).addUserMessage(text);

    // AI 응답 요청 (gRPC)
    await ref.read(coachChatHistoryProvider.notifier).requestCoachResponse(text);

    // 스크롤 맨 아래로
    _scrollController.animateTo(
      _scrollController.position.maxScrollExtent,
      duration: Duration(milliseconds: 300),
      curve: Curves.easeOut,
    );
  }

  Future<void> _toggleVoiceInput() async {
    if (_isListening) {
      // 음성 인식 중지
      final recognizedText = await ref.read(speechRecognizerProvider).stop();
      if (recognizedText.isNotEmpty) {
        _sendMessage(recognizedText);
      }
    } else {
      // 음성 인식 시작
      setState(() => _isListening = true);
      await ref.read(speechRecognizerProvider).start(
        onResult: (text) {
          _messageController.text = text;
        },
        onDone: () {
          setState(() => _isListening = false);
        },
      );
    }
    setState(() => _isListening = !_isListening);
  }
}

/// 코칭 카드 위젯
class CoachingCard extends StatelessWidget {
  final DailyCoachingCard card;
  final VoidCallback onTap;
  final VoidCallback onDismiss;

  @override
  Widget build(BuildContext context) {
    return Dismissible(
      key: Key(card.id),
      direction: DismissDirection.up,
      onDismissed: (_) => onDismiss(),
      child: Container(
        width: 200,
        margin: EdgeInsets.only(right: 12),
        child: Card(
          color: _getCardColor(card.type),
          child: InkWell(
            onTap: onTap,
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        _getCardIcon(card.type),
                        size: 20,
                        color: _getCardIconColor(card.type),
                      ),
                      SizedBox(width: 8),
                      Text(
                        _getCardLabel(card.type),
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: _getCardIconColor(card.type),
                        ),
                      ),
                    ],
                  ),
                  Spacer(),
                  Text(
                    card.title,
                    style: Theme.of(context).textTheme.titleSmall,
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  if (card.action != null)
                    Text(
                      card.action!,
                      style: TextStyle(
                        color: Theme.of(context).primaryColor,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
```

---

## 6. Phase 3-5 세부 구현 기획

### 6.1 Phase 3: Advanced (의료서비스, 커뮤니티)

#### 화상진료 시스템 (WebRTC 기반)

```yaml
# Phase 3 화상진료 아키텍처
telemedicine_architecture:
  components:
    signaling_server:
      technology: Go + WebSocket
      responsibilities:
        - SDP Offer/Answer 교환
        - ICE Candidate 중계
        - 세션 관리

    turn_server:
      technology: coturn
      purpose: NAT 우회, 방화벽 통과

    media_server:
      technology: Pion WebRTC (Go)
      features:
        - 1:1 화상 통화
        - 화면 공유
        - 녹화 (규정 준수)

    mobile_client:
      technology: Flutter + flutter_webrtc
      features:
        - 전면/후면 카메라 전환
        - 음소거/영상 끄기
        - PiP (Picture-in-Picture)
        - 배경 블러 (AI)

  workflow:
    1_booking:
      - 의사 검색 (reservation-service)
      - 예약 생성
      - 대기열 입장

    2_session_init:
      - WebSocket 연결
      - 방 입장
      - SDP Offer 생성

    3_connection:
      - STUN/TURN으로 연결 수립
      - 미디어 스트림 시작

    4_consultation:
      - 실시간 영상/음성
      - 측정 데이터 공유
      - 화면 공유 (검사 결과)

    5_conclusion:
      - 처방전 생성 (prescription-service)
      - 녹화 저장 (S3 + 암호화)
      - 리뷰/평점
```

### 6.2 Phase 4: Ecosystem (SDK, AI 학습)

#### 연합학습 시스템 상세

```yaml
# Phase 4 연합학습 아키텍처
federated_learning:
  framework: Flower 1.17+

  client_deployment:
    mobile_app:
      - 로컬 모델 학습 (TFLite)
      - 배터리/네트워크 최적화
      - 백그라운드 학습

    reader_device:
      - 엣지 모델 학습 (TFLite Micro)
      - 간헐적 동기화
      - 저전력 최적화

  server_deployment:
    aggregation_server:
      - FedAvg / FedProx 전략
      - Secure Aggregation
      - 모델 버전 관리

    model_registry:
      - MLflow 통합
      - A/B 테스트
      - 롤백 지원

  privacy_mechanisms:
    differential_privacy:
      epsilon: 1.0-10.0  # 프라이버시 예산
      delta: 1e-5
      noise_mechanism: gaussian

    secure_aggregation:
      protocol: Bonawitz et al. 2017
      threshold: "t-out-of-n"

  update_schedule:
    frequency: weekly
    min_participants: 100
    convergence_check: true

  model_types:
    biomarker_classifier:
      base: MobileNetV3
      input: 896D fingerprint
      output: 29 classes

    anomaly_detector:
      base: Isolation Forest + LSTM
      input: time series
      output: anomaly score

    health_predictor:
      base: Transformer
      input: 30-day history
      output: 7-day forecast
```

### 6.3 Phase 5: Future (1792차원, 웨어러블)

```yaml
# Phase 5 미래 기술 로드맵
phase_5_roadmap:
  dimension_expansion:
    current: 896D
    target: 1792D
    method:
      - 시간축 통합 (30일 × 896D → 압축)
      - 교차 카트리지 상관관계
      - 환경 컨텍스트 임베딩

  wearable_integration:
    supported_devices:
      - Apple Watch (HealthKit)
      - Galaxy Watch (Samsung Health)
      - Fitbit (Web API)
      - Garmin (Connect IQ)

    data_types:
      - 심박수 (실시간)
      - 수면 단계
      - 활동량/걸음수
      - 혈중 산소 (SpO2)
      - ECG (심전도)

    sync_mechanism:
      - HealthKit/Health Connect API
      - 백그라운드 동기화
      - CRDT 병합

  smart_home_integration:
    platforms:
      - Apple HomeKit
      - Google Home
      - Samsung SmartThings
      - Amazon Alexa

    use_cases:
      - "헤이 시리, 오늘 건강 상태 알려줘"
      - 이상 감지 시 조명 알림
      - 공기질 연동 환기 제어
      - 취침 시간 자동 조절
```

---

## 7. 규정 문서 작성 계획

### 7.1 IEC 62304 필수 문서

[IEC 62304 Wikipedia](https://en.wikipedia.org/wiki/IEC_62304) 및 [TÜV SÜD Guide](https://www.tuvsud.com/en-us/industries/healthcare-and-medical-devices/medical-devices-and-ivd/quality-management-and-quality-control-for-medical-devices/iec-62304-medical-device-software) 기반:

```yaml
iec_62304_documents:
  software_development_plan:
    id: DOC-SDP-001
    title: Software Development Plan (SDP)
    content:
      - 개발 프로세스 정의
      - 역할 및 책임
      - 형상관리 계획
      - 도구 및 환경
      - 일정 및 마일스톤
    deadline: +2주

  software_requirements_specification:
    id: DOC-SRS-001
    title: Software Requirements Specification (SRS)
    content:
      - 기능 요구사항 (80개 REQ-XXX)
      - 성능 요구사항
      - 인터페이스 요구사항
      - 보안 요구사항
      - 추적성 매트릭스
    deadline: +3주

  software_architecture_document:
    id: DOC-SAD-001
    title: Software Architecture Document (SAD)
    content:
      - 시스템 개요
      - 아키텍처 뷰 (4+1 View)
      - 컴포넌트 설계
      - 인터페이스 설계
      - 데이터 설계
      - 보안 아키텍처
    deadline: +3주

  software_detailed_design:
    id: DOC-SDD-001
    title: Software Detailed Design (SDD)
    content:
      - 모듈별 상세 설계
      - 알고리즘 설명
      - 데이터 구조
      - 에러 처리
    deadline: +4주

  software_verification_plan:
    id: DOC-SVP-001
    title: Software Verification Plan
    content:
      - 검증 전략
      - 테스트 레벨 정의
      - 합격 기준
      - 도구 및 환경
    deadline: +2주
```

### 7.2 ISO 14971 위험관리 문서

```yaml
iso_14971_documents:
  risk_management_plan:
    id: DOC-RMP-001
    content:
      - 위험관리 정책
      - 역할 및 책임
      - 위험 허용 기준
      - 검증 활동

  hazard_identification:
    id: DOC-HID-001
    methods:
      - FMEA (Failure Mode and Effects Analysis)
      - FTA (Fault Tree Analysis)
      - HAZOP (Hazard and Operability)

  risk_estimation:
    id: DOC-RES-001
    content:
      - 심각도 분류 (S1-S4)
      - 발생 확률 (P1-P5)
      - 위험 매트릭스

  risk_evaluation:
    id: DOC-REV-001
    content:
      - 허용 가능 위험
      - ALARP (As Low As Reasonably Practicable)
      - 리스크-베네핏 분석

  risk_control:
    id: DOC-RCT-001
    measures:
      - 본질 안전 설계
      - 보호 조치
      - 안전 정보 제공

  residual_risk_evaluation:
    id: DOC-RRE-001
    content:
      - 잔여 위험 목록
      - 전체 잔여 위험 수용성
```

---

## 8. 시너지 극대화 연동 설계

### 8.1 서비스 간 시너지 매트릭스

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    서비스 간 시너지 매트릭스                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│           측정 → AI → 코칭 → 알림 → 행동 변화                            │
│           ─────────────────────────────────────────                     │
│                                                                         │
│  [measurement]                                                          │
│       ↓ 측정 완료 이벤트                                                 │
│  [ai-inference] ←────────────────────────────────────────┐              │
│       ↓ 분석 결과                     ↑                  │              │
│  [coaching] ─────────→ 개인화 추천 생성                   │              │
│       ↓                               │ 피드백            │              │
│  [notification] ────→ 푸시/음성 전달  │                  │              │
│       ↓                               │                  │              │
│  [Flutter App] ─────→ 사용자 행동 ────┘                  │              │
│       ↓                                                  │              │
│  [다음 측정] ────────────────────────────────────────────┘              │
│                                                                         │
│                                                                         │
│           구독 → 접근 → 카트리지 → 측정 → 결제                           │
│           ─────────────────────────────────────────                     │
│                                                                         │
│  [subscription] ────→ 티어별 접근 권한                                   │
│       ↓                                                                 │
│  [cartridge] ←──────── 접근 제어 확인                                    │
│       ↓                                                                 │
│  [measurement] ←────── 보정 데이터 제공                                  │
│       ↓                                                                 │
│  [payment] ←────────── 애드온 카트리지 결제                              │
│       ↓                                                                 │
│  [shop] ←───────────── 카트리지 재구매                                   │
│                                                                         │
│                                                                         │
│           가족 → 공유 → 알림 → 의료 → 긴급                               │
│           ─────────────────────────────────────────                     │
│                                                                         │
│  [family] ──────────→ 가족 그룹 생성                                     │
│       ↓                                                                 │
│  [health-record] ←─── 건강 데이터 공유 동의                              │
│       ↓                                                                 │
│  [notification] ←──── 보호자 알림 설정                                   │
│       ↓                                                                 │
│  [emergency] ←─────── 위험 수치 감지                                     │
│       ↓                                                                 │
│  [telemedicine] ←──── 의료 상담 연결                                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.2 자율 최적화 (Self-Optimization)

시스템이 스스로 학습하고 최적화하는 메커니즘:

```yaml
self_optimization:
  measurement_optimization:
    trigger: 측정 품질 저하 감지
    actions:
      - 보정 데이터 자동 업데이트 요청
      - 측정 환경 권고 (온도, 습도)
      - 카트리지 교체 알림

  ai_model_optimization:
    trigger: 예측 정확도 저하
    actions:
      - 연합학습 라운드 트리거
      - 모델 버전 롤백 검토
      - 사용자 피드백 수집

  user_engagement_optimization:
    trigger: 사용 빈도 감소
    actions:
      - 리마인더 시간 최적화
      - 게이미피케이션 강화
      - 개인화 콘텐츠 추천

  system_performance_optimization:
    trigger: 응답 시간 증가
    actions:
      - 캐시 정책 조정
      - DB 쿼리 최적화
      - 서비스 스케일 아웃
```

---

## 9. 참조 문헌 및 출처

### Healthcare AI & Architecture
- [McKinsey - Healthcare AI Modular Architecture](https://www.mckinsey.com/industries/healthcare/our-insights/the-coming-evolution-of-healthcare-ai-toward-a-modular-architecture)
- [Corti - Multi-Agent AI Framework](https://www.corti.ai)
- [World Economic Forum - Healthcare Data Architecture](https://www.weforum.org/stories/2026/01/ai-healthcare-data-architecture/)
- [Frontiers - Health Information Systems Architecture](https://www.frontiersin.org/journals/digital-health/articles/10.3389/fdgth.2025.1694839/full)

### Federated Learning & Privacy
- [Nature - Federated Blockchain Healthcare](https://www.nature.com/articles/s41598-025-04083-4)
- [PMC - Federated Learning Healthcare Review](https://pmc.ncbi.nlm.nih.gov/articles/PMC11728217/)
- [JMIR AI - Personal Health Train](https://ai.jmir.org/2025/1/e60847)

### Biosensor & ML
- [Nature - Plasma Infrared Fingerprinting](https://pmc.ncbi.nlm.nih.gov/articles/PMC11293328/)
- [Wiley - AI Biosensors](https://advanced.onlinelibrary.wiley.com/doi/full/10.1002/adma.202504796)
- [RSC - Surface-Enhanced Spectroscopy ML](https://pubs.rsc.org/en/content/articlehtml/2023/na/d2na00608a)

### Regulatory & Compliance
- [IEC 62304 - Wikipedia](https://en.wikipedia.org/wiki/IEC_62304)
- [TÜV SÜD - IEC 62304 Guide](https://www.tuvsud.com/en-us/industries/healthcare-and-medical-devices/medical-devices-and-ivd/quality-management-and-quality-control-for-medical-devices/iec-62304-medical-device-software)
- [FDA 510(k) IEC 62304](https://mavenprofserv.com/blog/iec-62304-510k-approval/)

### Rust Embedded
- [Embassy-rs TrouBLE](https://github.com/embassy-rs/trouble)
- [TFLite Micro](https://github.com/tensorflow/tflite-micro)
- [219 Design - BLE with Rust](https://www.219design.com/bluetooth-low-energy-with-rust/)

### Flutter & UX
- [Code with Andrea - Riverpod Architecture](https://codewithandrea.com/articles/flutter-app-architecture-riverpod-introduction/)
- [gRPC Dart + Riverpod](https://grpc-dart-docs.pages.dev/docs/grpc-basics/grpc-riverpod-client/)
- [GeekyAnts - Offline-First Flutter](https://geekyants.com/blog/offline-first-flutter-implementation-blueprint-for-real-world-apps)
- [UX Studio - Healthcare UX 2026](https://www.uxstudioteam.com/ux-blog/healthcare-ux)
- [Eleken - Healthcare UI Design](https://www.eleken.co/blog-posts/user-interface-design-for-healthcare-applications)

---

**문서 종료**

*본 마스터플랜은 유사 시스템 조사, 최신 기술 트렌드 분석, 학술 연구를 기반으로 작성되었습니다. 모든 IDE 및 AI 에이전트가 참조하여 일관된 구현을 수행할 수 있도록 설계되었습니다.*

*작성일: 2026-02-12 | 버전: v2.0 | 총 라인: 1500+*
