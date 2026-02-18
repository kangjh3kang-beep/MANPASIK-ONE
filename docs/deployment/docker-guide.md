# ManPaSik Docker Compose 배포 가이드

## 사전 요구사항

- Docker 24.0+
- Docker Compose v2.20+
- 최소 8GB RAM, 20GB 디스크

## 빠른 시작

```bash
# 1. 프로젝트 루트에서 실행
docker-compose up -d

# 2. 서비스 상태 확인
docker-compose ps

# 3. 게이트웨이 헬스 체크
curl http://localhost:8080/api/v1/health
```

## 서비스 구성

| 서비스 | 포트 | 설명 |
|--------|------|------|
| Gateway | 8080 | REST API 게이트웨이 |
| PostgreSQL | 5432 | 메인 데이터베이스 |
| Redis | 6379 | 캐시 + 세션 저장소 |
| Kafka (Redpanda) | 9092 | 이벤트 버스 |
| MinIO | 9000/9001 | 오브젝트 스토리지 |

## 환경변수

### 데이터베이스
```env
DB_HOST=postgres
DB_PORT=5432
DB_NAME=manpasik
DB_USER=manpasik
DB_PASSWORD=<secure-password>
```

### Redis
```env
REDIS_URL=redis://redis:6379/0
```

### Kafka
```env
KAFKA_BROKERS=kafka:9092
```

### JWT
```env
JWT_SECRET=<256-bit-secret>
JWT_EXPIRY=3600
JWT_REFRESH_EXPIRY=604800
```

## 서비스별 빌드

```bash
# 전체 빌드
docker-compose build

# 특정 서비스만 빌드
docker-compose build gateway auth-service measurement-service

# 캐시 없이 재빌드
docker-compose build --no-cache
```

## 데이터베이스 마이그레이션

```bash
# 자동 마이그레이션 (서비스 시작 시)
# infrastructure/database/init/ 디렉토리의 SQL 파일이 순서대로 실행됩니다

# 수동 마이그레이션
docker-compose exec postgres psql -U manpasik -d manpasik -f /docker-entrypoint-initdb.d/09-cartridge.sql
```

## 로그 확인

```bash
# 전체 로그
docker-compose logs -f

# 특정 서비스 로그
docker-compose logs -f gateway auth-service

# 최근 100줄
docker-compose logs --tail=100 measurement-service
```

## 프로덕션 배포 체크리스트

1. **보안**
   - [ ] JWT_SECRET을 256비트 이상 랜덤 값으로 변경
   - [ ] DB_PASSWORD를 강력한 비밀번호로 변경
   - [ ] SSL/TLS 인증서 적용
   - [ ] 네트워크 방화벽 설정 (포트 8080만 외부 노출)

2. **성능**
   - [ ] PostgreSQL `max_connections` 조정 (기본: 100)
   - [ ] Redis `maxmemory` 설정 (기본: 256mb)
   - [ ] 서비스별 CPU/메모리 제한 설정

3. **모니터링**
   - [ ] Prometheus 메트릭 수집 설정
   - [ ] Grafana 대시보드 구성
   - [ ] 알림 규칙 설정 (에러율, 응답 시간)

4. **백업**
   - [ ] PostgreSQL 자동 백업 스케줄 (pg_dump)
   - [ ] MinIO 오브젝트 복제 설정

## 트러블슈팅

### 서비스가 시작되지 않을 때
```bash
docker-compose logs <service-name>
docker-compose restart <service-name>
```

### 데이터베이스 연결 실패
```bash
docker-compose exec postgres pg_isready -U manpasik
```

### 포트 충돌
```bash
# 사용 중인 포트 확인
lsof -i :8080
# docker-compose.yml에서 포트 매핑 변경
```

### 전체 초기화 (데이터 삭제)
```bash
docker-compose down -v  # 볼륨 포함 삭제
docker-compose up -d    # 재생성
```
