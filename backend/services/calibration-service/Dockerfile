# calibration-service gRPC :50060
# Build from backend root: docker build -f services/calibration-service/Dockerfile -t manpasik/calibration-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/calibration-service/ services/calibration-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /calibration-service ./services/calibration-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /calibration-service .
EXPOSE 50060
ENTRYPOINT ["./calibration-service"]
