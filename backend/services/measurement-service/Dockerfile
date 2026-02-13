# measurement-service gRPC :50054
# Build from backend root: docker build -f services/measurement-service/Dockerfile -t manpasik/measurement-service:dev .
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ shared/
COPY services/measurement-service/ services/measurement-service/
RUN go mod tidy && CGO_ENABLED=0 go build -o /measurement-service ./services/measurement-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /measurement-service .
EXPOSE 50054
ENTRYPOINT ["./measurement-service"]
