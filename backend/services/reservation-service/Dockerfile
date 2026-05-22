# reservation-service Dockerfile
# 포트: 50066

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /reservation-service ./services/reservation-service/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /reservation-service /usr/local/bin/reservation-service

EXPOSE 50066

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO /dev/null http://localhost:9100/health 2>/dev/null || exit 1

CMD ["reservation-service"]
