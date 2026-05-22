# cartridge-service gRPC :50059
# Build from backend root: docker build -f services/cartridge-service/Dockerfile -t manpasik/cartridge-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/cartridge-service/ services/cartridge-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /cartridge-service ./services/cartridge-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /cartridge-service .
EXPOSE 50059
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO /dev/null http://localhost:9100/health 2>/dev/null || exit 1

ENTRYPOINT ["./cartridge-service"]
