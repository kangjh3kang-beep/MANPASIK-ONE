# device-service gRPC :50053
# Build from backend root: docker build -f services/device-service/Dockerfile -t manpasik/device-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/device-service/ services/device-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /device-service ./services/device-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /device-service .
EXPOSE 50053
ENTRYPOINT ["./device-service"]
