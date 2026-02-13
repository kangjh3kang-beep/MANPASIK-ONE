# telemedicine-service Dockerfile
# 포트: 50065

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /telemedicine-service ./services/telemedicine-service/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /telemedicine-service /usr/local/bin/telemedicine-service

EXPOSE 50065

CMD ["telemedicine-service"]
