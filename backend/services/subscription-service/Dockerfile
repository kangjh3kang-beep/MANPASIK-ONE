# subscription-service gRPC :50055
# Build from backend root: docker build -f services/subscription-service/Dockerfile -t manpasik/subscription-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/subscription-service/ services/subscription-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /subscription-service ./services/subscription-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /subscription-service .
EXPOSE 50055
ENTRYPOINT ["./subscription-service"]
