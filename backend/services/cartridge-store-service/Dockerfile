FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /cartridge-store-service ./services/cartridge-store-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /cartridge-store-service /cartridge-store-service
EXPOSE 50079
ENTRYPOINT ["/cartridge-store-service"]
