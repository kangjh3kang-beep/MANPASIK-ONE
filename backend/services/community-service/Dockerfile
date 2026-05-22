# community-service gRPC :50067
# Build from backend root: docker build -f services/community-service/Dockerfile -t manpasik/community-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/community-service/ services/community-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /community-service ./services/community-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /community-service .
EXPOSE 50067
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO /dev/null http://localhost:9100/health 2>/dev/null || exit 1

ENTRYPOINT ["./community-service"]
