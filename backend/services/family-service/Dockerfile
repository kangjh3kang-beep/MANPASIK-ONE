# family-service Dockerfile
# 포트: 50063

FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -o /family-service ./services/family-service/cmd/main.go

# --- Runtime ---
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /family-service /usr/local/bin/family-service

EXPOSE 50063

CMD ["family-service"]
