# admin-service gRPC :50068
# Build from backend root: docker build -f services/admin-service/Dockerfile -t manpasik/admin-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/admin-service/ services/admin-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /admin-service ./services/admin-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /admin-service .
EXPOSE 50068
ENTRYPOINT ["./admin-service"]
