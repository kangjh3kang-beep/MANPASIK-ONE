# marketplace-service Dockerfile
# 포트: 8080

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /marketplace-service ./services/marketplace-service/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /marketplace-service /usr/local/bin/marketplace-service

EXPOSE 8080

CMD ["marketplace-service"]
