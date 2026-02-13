# coaching-service gRPC :50061
# Build from backend root: docker build -f services/coaching-service/Dockerfile -t manpasik/coaching-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/coaching-service/ services/coaching-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /coaching-service ./services/coaching-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /coaching-service .
EXPOSE 50061
ENTRYPOINT ["./coaching-service"]
