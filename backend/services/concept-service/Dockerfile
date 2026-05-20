FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /concept-service ./services/concept-service/cmd

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /concept-service /concept-service
EXPOSE 50078
ENTRYPOINT ["/concept-service"]
