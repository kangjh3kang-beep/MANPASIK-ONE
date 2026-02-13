# ManPaSik Go 서비스 멀티스테이지 Dockerfile
# 최적화된 빌드: 작은 이미지, 보안 강화

# ============================================================================
# Stage 1: 빌드
# ============================================================================
FROM golang:1.24-alpine AS builder

# 빌드 인자
ARG SERVICE_NAME=gateway
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown

# 빌드 종속성
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# 모듈 캐시 (레이어 최적화)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 소스 복사
COPY backend/ ./

# 빌드
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/${SERVICE_NAME} \
    ./${SERVICE_NAME}/cmd || \
    go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/${SERVICE_NAME} \
    ./services/${SERVICE_NAME}/cmd

# ============================================================================
# Stage 2: 런타임
# ============================================================================
FROM scratch

# 메타데이터
LABEL org.opencontainers.image.title="ManPaSik Service"
LABEL org.opencontainers.image.vendor="ManPaSik"
LABEL org.opencontainers.image.source="https://github.com/manpasik/manpasik"

# 인증서 및 타임존
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 바이너리
COPY --from=builder /app/* /app/

# 비루트 사용자 (UID 65534 = nobody)
USER 65534:65534

# 환경변수
ENV TZ=Asia/Seoul

# 헬스체크
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/service", "health"]

# 진입점
ENTRYPOINT ["/app/service"]
