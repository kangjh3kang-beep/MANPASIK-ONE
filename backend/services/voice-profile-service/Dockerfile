FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /voice-profile-service ./services/voice-profile-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /voice-profile-service /voice-profile-service
EXPOSE 50081
ENTRYPOINT ["/voice-profile-service"]
