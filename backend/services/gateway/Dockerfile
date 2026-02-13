# REST API Gateway HTTP :8080
# Build from backend root: docker build -f services/gateway/Dockerfile -t manpasik/gateway:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY gateway/ gateway/
RUN go mod tidy && CGO_ENABLED=0 go build -o /gateway ./gateway/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /gateway .
EXPOSE 8080
ENTRYPOINT ["./gateway"]
