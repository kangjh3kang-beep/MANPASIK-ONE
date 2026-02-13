# health-record-service Dockerfile
# 포트: 50064

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /health-record-service ./services/health-record-service/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /health-record-service /usr/local/bin/health-record-service

EXPOSE 50064

CMD ["health-record-service"]
