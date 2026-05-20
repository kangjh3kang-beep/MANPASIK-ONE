# digital-twin-service gRPC :50073
# Build from backend root: docker build -f services/digital-twin-service/Dockerfile -t manpasik/digital-twin-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/digital-twin-service/ services/digital-twin-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /digital-twin-service ./services/digital-twin-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /digital-twin-service .
EXPOSE 50073
ENTRYPOINT ["./digital-twin-service"]
