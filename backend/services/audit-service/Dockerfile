# audit-service gRPC :50072
# Build from backend root: docker build -f services/audit-service/Dockerfile -t manpasik/audit-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/audit-service/ services/audit-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /audit-service ./services/audit-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /audit-service .
EXPOSE 50072
ENTRYPOINT ["./audit-service"]
