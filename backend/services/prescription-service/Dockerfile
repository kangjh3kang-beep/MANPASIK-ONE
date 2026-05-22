FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /prescription-service ./services/prescription-service/cmd/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /prescription-service /prescription-service
EXPOSE 50069
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO /dev/null http://localhost:9100/health 2>/dev/null || exit 1

CMD ["/prescription-service"]
