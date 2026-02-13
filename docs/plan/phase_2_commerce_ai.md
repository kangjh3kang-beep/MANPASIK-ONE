# Phase 2: Core (커머스 + AI, Month 5–8)

> **전제**: Phase 1D 완료 (E2E 플로우 검증, Docker 빌드·기동 정상, gRPC 4서비스 연동)
> **목표**: SaaS 구독·결제 기반 구축, AI 추론 파이프라인 구현, 카트리지 관리 체계 수립

---

## 1. 범위 및 우선순위

| 순서 | 서비스 | 도메인 | 포트 | 의존성 | 목표 시점 |
|------|--------|--------|------|--------|-----------|
| 1 | subscription-service | 커머스 | 50055 | auth, user | Month 5 Week 1-2 |
| 2 | shop-service | 커머스 | 50056 | auth, subscription | Month 5 Week 3-4 |
| 3 | payment-service | 커머스 | 50057 | auth, subscription, shop | Month 6 Week 1-2 |
| 4 | ai-inference-service | AI | 50058 | auth, measurement | Month 6 Week 3-4 |
| 5 | cartridge-service | 측정 | 50059 | auth, device | Month 6 Week 3-4 |
| 6 | coaching-service | AI | 50060 | auth, ai-inference, measurement | Month 7 Week 1-2 |
| 7 | calibration-service | 측정 | 50061 | auth, cartridge, measurement | Month 7 Week 3-4 |

---

## 2. 서비스별 상세

### 2.1 subscription-service (포트 50055)

**역할**: SaaS 구독 관리, 티어별 기능 제어

**gRPC API:**
| RPC | 설명 |
|-----|------|
| CreateSubscription | 구독 생성 (회원가입 시 Free 자동) |
| GetSubscription | 구독 정보 조회 |
| UpdateSubscription | 티어 변경 (업/다운그레이드) |
| CancelSubscription | 구독 해지 |
| CheckFeatureAccess | 기능 접근 권한 확인 |
| ListSubscriptionPlans | 구독 플랜 목록 조회 |

**구독 티어:**
| 티어 | 월 가격 | 리더기 | 가족 | AI 코칭 | 화상진료 |
|------|---------|--------|------|---------|----------|
| Free | ₩0 | 1대 | 0명 | ❌ | ❌ |
| Basic | ₩9,900 | 3대 | 2명 | ❌ | ❌ |
| Pro | ₩29,900 | 5대 | 5명 | ✅ | ❌ |
| Clinical | ₩59,900 | 10대 | 10명 | ✅ | ✅ |

**DB 테이블:**
- `subscriptions`: 사용자별 구독 상태
- `subscription_plans`: 플랜 정의 (불변 참조)
- `subscription_history`: 변경 이력 (감사 추적)

---

### 2.2 shop-service (포트 50056)

**역할**: 상품 관리, 장바구니, 주문 처리

**gRPC API:**
| RPC | 설명 |
|-----|------|
| ListProducts | 상품 목록 (카트리지, 리더기, 액세서리) |
| GetProduct | 상품 상세 |
| AddToCart | 장바구니 추가 |
| GetCart | 장바구니 조회 |
| RemoveFromCart | 장바구니 항목 제거 |
| CreateOrder | 주문 생성 |
| GetOrder | 주문 상세 |
| ListOrders | 주문 이력 |

**DB 테이블:**
- `products`: 상품 정보 (카트리지, 리더기, 액세서리)
- `cart_items`: 장바구니
- `orders`: 주문
- `order_items`: 주문 항목

---

### 2.3 payment-service (포트 50057)

**역할**: PG 연동, 결제 처리, 구독 자동 결제

**gRPC API:**
| RPC | 설명 |
|-----|------|
| CreatePayment | 결제 요청 (일회성·구독) |
| ConfirmPayment | 결제 확인 (PG 콜백) |
| GetPayment | 결제 상세 |
| ListPayments | 결제 이력 |
| RefundPayment | 환불 처리 |
| RegisterPaymentMethod | 결제 수단 등록 |

**PG 연동**: Toss Payments / NHN KCP (추상화 인터페이스)

**DB 테이블:**
- `payments`: 결제 기록
- `payment_methods`: 결제 수단
- `refunds`: 환불 기록

---

### 2.4 ai-inference-service (포트 50058)

**역할**: 실시간 AI 추론, 모델 서빙, 건강 분석

**gRPC API:**
| RPC | 설명 |
|-----|------|
| AnalyzeMeasurement | 측정 데이터 AI 분석 |
| GetHealthScore | 건강 점수 산출 |
| PredictTrend | 트렌드 예측 |
| GetModelInfo | 모델 정보 조회 |
| StreamAnalysis | 실시간 분석 스트림 |

**AI 모델 (5종):**
1. BiomarkerClassifier: 바이오마커 분류
2. AnomalyDetector: 이상치 탐지
3. TrendPredictor: 트렌드 예측
4. HealthScorer: 건강 점수 산출
5. FoodCalorieEstimator: 음식 칼로리 추정 (Phase 2 후반)

---

### 2.5 cartridge-service (포트 50059)

**역할**: NFC 카트리지 인증, 보정 데이터, 사용 추적

**gRPC API:**
| RPC | 설명 |
|-----|------|
| AuthenticateCartridge | 카트리지 NFC 인증 |
| GetCartridgeInfo | 카트리지 정보 조회 |
| RecordUsage | 사용 횟수 기록 |
| GetCalibrationData | 보정 데이터 조회 |
| ListCartridgeTypes | 카트리지 종류 목록 (29종) |

---

### 2.6 coaching-service (포트 50060)

**역할**: AI 기반 건강 코칭, 개인화 추천

**gRPC API:**
| RPC | 설명 |
|-----|------|
| GetCoachingAdvice | 건강 조언 생성 |
| GetDailyPlan | 일일 건강 계획 |
| GetWeeklySummary | 주간 요약 |
| SetGoals | 건강 목표 설정 |
| GetGoalProgress | 목표 진행 상황 |

---

### 2.7 calibration-service (포트 50061)

**역할**: 보정 모델 관리, 팩토리 보정

**gRPC API:**
| RPC | 설명 |
|-----|------|
| RunCalibration | 보정 실행 |
| GetCalibrationStatus | 보정 상태 |
| ApplyCalibrationModel | 보정 모델 적용 |
| GetCalibrationHistory | 보정 이력 |

---

## 3. 공통 패턴

### 3.1 서비스 구조 (기존 패턴 유지)
```
backend/services/{service-name}/
├── cmd/main.go                    # gRPC 서버 진입점
├── internal/
│   ├── handler/grpc.go            # gRPC 핸들러
│   ├── service/{name}.go          # 비즈니스 로직
│   ├── service/{name}_test.go     # 유닛 테스트
│   └── repository/
│       ├── memory/{name}.go       # 인메모리 (개발/테스트)
│       └── postgres/{name}.go     # PostgreSQL (프로덕션)
└── Dockerfile                     # golang:1.24-alpine
```

### 3.2 포트 할당 확장
| 서비스 | 포트 | Phase |
|--------|------|-------|
| auth-service | 50051 | 1 |
| user-service | 50052 | 1 |
| device-service | 50053 | 1 |
| measurement-service | 50054 | 1 |
| subscription-service | 50055 | 2 |
| shop-service | 50056 | 2 |
| payment-service | 50057 | 2 |
| ai-inference-service | 50058 | 2 |
| cartridge-service | 50059 | 2 |
| coaching-service | 50060 | 2 |
| calibration-service | 50061 | 2 |

### 3.3 Proto 확장
- `manpasik.proto`에 Phase 2 서비스 gRPC 정의 추가
- `make proto` 후 생성 코드 반영

### 3.4 Docker Compose
- `docker-compose.dev.yml`에 Phase 2 서비스 추가
- 공통 환경변수 (DB, JWT 등) 유지

---

## 4. Phase 2 Gate 통과 기준

- [ ] 7개 서비스 빌드·테스트 통과
- [ ] 서비스별 유닛 테스트 80%+ 커버리지
- [ ] E2E 플로우: 구독 생성 → 상품 주문 → 결제 → AI 분석
- [ ] Docker Compose 전체 기동 (11서비스)
- [ ] 보안: JWT 인증·RBAC 전 서비스 적용
- [ ] 문서 갱신: CHANGELOG, CONTEXT, QUALITY_GATES

---

**문서 버전**: 1.0  
**최종 업데이트**: 2026-02-10
