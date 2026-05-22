# gateway Dockerfile
# 포트: HTTP 8080

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo     -o /gateway ./services/gateway/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /gateway /usr/local/bin/gateway

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/api/v1/health || exit 1

CMD ["gateway"]
