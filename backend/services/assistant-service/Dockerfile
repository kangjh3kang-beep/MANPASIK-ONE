# assistant-service gRPC :50077
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/assistant-service/ services/assistant-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /assistant-service ./services/assistant-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /assistant-service .
EXPOSE 50077
ENTRYPOINT ["./assistant-service"]
