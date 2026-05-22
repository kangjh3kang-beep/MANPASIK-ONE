FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /cartridge-store-service ./services/cartridge-store-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /cartridge-store-service /cartridge-store-service
EXPOSE 50079
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO /dev/null http://localhost:9100/health 2>/dev/null || exit 1

ENTRYPOINT ["/cartridge-store-service"]
