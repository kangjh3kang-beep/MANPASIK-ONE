# auth-service gRPC :50051
# Build from backend root: docker build -f services/auth-service/Dockerfile -t manpasik/auth-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/auth-service/ services/auth-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /auth-service ./services/auth-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /auth-service .
EXPOSE 50051
ENTRYPOINT ["./auth-service"]
