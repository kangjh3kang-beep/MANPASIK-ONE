# user-service gRPC :50052
# Build from backend root: docker build -f services/user-service/Dockerfile -t manpasik/user-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/user-service/ services/user-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /user-service ./services/user-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /user-service .
EXPOSE 50052
ENTRYPOINT ["./user-service"]
