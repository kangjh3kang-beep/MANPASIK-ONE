FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /data-platform-service ./services/data-platform-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /data-platform-service /data-platform-service
EXPOSE 50080
ENTRYPOINT ["/data-platform-service"]
