# payment-service gRPC :50057
# Build from backend root: docker build -f services/payment-service/Dockerfile -t manpasik/payment-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/payment-service/ services/payment-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /payment-service ./services/payment-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /payment-service .
EXPOSE 50057
ENTRYPOINT ["./payment-service"]
