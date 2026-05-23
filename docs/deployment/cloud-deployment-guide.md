# ManPaSik 클라우드 배포 가이드

> **에픽** E1 · **인프라** Kubernetes (EKS) / Railway / Fly.io

---

## 1. 배포 전략

### 1.1 단계적 전환

| 단계 | 환경 | 데이터 | 상태 |
|------|------|--------|------|
| **현재** | Cloudflare Workers (Web) | Next.js API Route (데모) | ✅ 운영 |
| **Phase A** | + Supabase (PostgreSQL) | 실제 DB 연동 | 대기 |
| **Phase B** | + Railway/Fly.io (Go Gateway) | 37 gRPC 서비스 | 대기 |
| **Phase C** | + EKS (Kubernetes) | 프로덕션 전체 | 대기 |

### 1.2 Supabase 연동 (Phase A)

현재 환경변수:
```
VITE_SUPABASE_URL=https://nwjozczygkw...
VITE_SUPABASE_ANON_KEY=eyJhbGciOiJ...
```

Next.js API Route에서 Supabase 직접 연동:
```typescript
// apps/web/app/api/v1/[...path]/route.ts 수정
import { createClient } from '@supabase/supabase-js'

const supabase = createClient(
  process.env.VITE_SUPABASE_URL!,
  process.env.VITE_SUPABASE_ANON_KEY!
)

// 데모 데이터 대신 Supabase 조회
const data = await supabase.from('measurements').select('*')
```

### 1.3 Railway 배포 (Phase B)

```yaml
# railway.toml
[build]
  builder = "dockerfile"
  dockerfilePath = "backend/services/gateway/Dockerfile"

[deploy]
  healthcheckPath = "/api/v1/health"
  healthcheckTimeout = 30
```

### 1.4 Kubernetes 배포 (Phase C)

이미 구성됨: `infrastructure/kubernetes/`
- base/: 서비스 41개 매니페스트
- overlays/: dev/staging/production

---

## 2. 환경변수 매트릭스

| 변수 | 개발 | 스테이징 | 프로덕션 |
|------|------|---------|---------|
| DB_HOST | localhost | supabase | RDS/Supabase |
| REDIS_HOST | localhost | Upstash | ElastiCache |
| JWT_SECRET | dev-secret | vault | Vault/AWS SM |
| S3_ENDPOINT | localhost:9000 | Cloudflare R2 | S3/R2 |

---

## 3. next.config.mjs 프록시 설정 (Go 백엔드 연결 시)

```javascript
// Go Gateway 연결 후 데모 API Route 대체
module.exports = {
  async rewrites() {
    return [
      {
        source: '/api/v1/:path*',
        destination: `${process.env.GATEWAY_URL}/api/v1/:path*`,
      },
    ]
  },
}
```
