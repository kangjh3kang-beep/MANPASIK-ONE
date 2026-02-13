# ai-inference-service gRPC :50058
# Build from backend root: docker build -f services/ai-inference-service/Dockerfile -t manpasik/ai-inference-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/ai-inference-service/ services/ai-inference-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /ai-inference-service ./services/ai-inference-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /ai-inference-service .
EXPOSE 50058
ENTRYPOINT ["./ai-inference-service"]
