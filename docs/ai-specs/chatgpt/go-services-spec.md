# Go 마이크로서비스 구현 스펙 (ChatGPT용)

## 프로젝트 정보

- **프로젝트명**: ManPaSik Backend
- **언어**: Go 1.22+
- **프레임워크**: gRPC + Gin (HTTP Gateway)
- **DB**: PostgreSQL 16
- **ORM**: sqlc 또는 GORM
- **인증**: Keycloak (OIDC)

---

## 요청 작업

### 1. 공통 서비스 구조

각 서비스는 다음 구조를 따릅니다:

```
backend/services/{service-name}/
├── cmd/
│   └── main.go              # 진입점
├── internal/
│   ├── handler/             # gRPC/HTTP 핸들러
│   ├── service/             # 비즈니스 로직
│   ├── repository/          # DB 접근
│   └── model/               # 도메인 모델
├── pkg/                     # 외부 공개 패키지
├── Dockerfile
└── Makefile
```

---

### 2. auth-service 구현

**기능:**
- Keycloak 연동 (OIDC)
- JWT 토큰 검증
- 세션 관리 (Redis)
- MFA 지원

**gRPC API (이미 정의됨: `shared/proto/manpasik.proto`):**
- 참조: UserService.GetProfile, GetSubscription

**main.go 예시:**
```go
package main

import (
    "context"
    "log"
    "net"
    "os"
    "os/signal"
    "syscall"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
    
    pb "github.com/manpasik/backend/shared/proto/v1"
)

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    
    // 헬스체크 등록
    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(s, healthServer)
    healthServer.SetServingStatus("auth-service", grpc_health_v1.HealthCheckResponse_SERVING)
    
    // 서비스 등록
    // pb.RegisterUserServiceServer(s, &userService{})
    
    // Graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        s.GracefulStop()
    }()
    
    log.Printf("auth-service listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
```

---

### 3. user-service 구현

**기능:**
- 사용자 프로필 CRUD
- 구독 관리
- 가족 그룹 관리

**DB 스키마:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    avatar_url TEXT,
    language VARCHAR(10) DEFAULT 'ko',
    timezone VARCHAR(50) DEFAULT 'Asia/Seoul',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    tier VARCHAR(20) NOT NULL, -- 'free', 'basic', 'pro', 'clinical'
    started_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ,
    max_devices INT DEFAULT 3,
    max_family_members INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE family_groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE family_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    family_id UUID REFERENCES family_groups(id),
    user_id UUID REFERENCES users(id),
    role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member'
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(family_id, user_id)
);
```

---

### 4. device-service 구현

**기능:**
- 리더기 등록/해제
- 펌웨어 OTA 관리
- 디바이스 상태 추적

**DB 스키마:**
```sql
CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id VARCHAR(100) UNIQUE NOT NULL, -- BLE MAC 또는 시리얼
    user_id UUID REFERENCES users(id),
    name VARCHAR(100),
    firmware_version VARCHAR(20),
    last_seen TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'offline',
    battery_percent INT,
    registered_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE device_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID REFERENCES devices(id),
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

### 5. measurement-service 구현

**기능:**
- 측정 세션 관리
- 측정 데이터 저장 (TimescaleDB)
- 핑거프린트 벡터 저장 (Milvus)

**TimescaleDB 스키마:**
```sql
CREATE TABLE measurements (
    time TIMESTAMPTZ NOT NULL,
    session_id UUID NOT NULL,
    device_id UUID NOT NULL,
    user_id UUID,
    cartridge_type VARCHAR(50),
    raw_channels DOUBLE PRECISION[],
    s_det DOUBLE PRECISION,
    s_ref DOUBLE PRECISION,
    alpha DOUBLE PRECISION,
    s_corrected DOUBLE PRECISION,
    primary_value DOUBLE PRECISION,
    unit VARCHAR(20),
    confidence FLOAT,
    fingerprint_vector REAL[], -- 88/448/896 차원
    temp_c FLOAT,
    humidity_pct FLOAT,
    battery_pct INT
);

SELECT create_hypertable('measurements', 'time');

-- 인덱스
CREATE INDEX idx_measurements_session ON measurements (session_id);
CREATE INDEX idx_measurements_user ON measurements (user_id, time DESC);
```

---

## 산출물

1. `backend/services/auth-service/cmd/main.go`
2. `backend/services/user-service/cmd/main.go`
3. `backend/services/device-service/cmd/main.go`
4. `backend/services/measurement-service/cmd/main.go`
5. SQL 마이그레이션 파일들

---

**작성자**: Antigravity (Gemini)
**작성일**: 2026-02-09
