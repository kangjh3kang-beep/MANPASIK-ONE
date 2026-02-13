# shop-service gRPC :50056
# Build from backend root: docker build -f services/shop-service/Dockerfile -t manpasik/shop-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/shop-service/ services/shop-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /shop-service ./services/shop-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /shop-service .
EXPOSE 50056
ENTRYPOINT ["./shop-service"]
