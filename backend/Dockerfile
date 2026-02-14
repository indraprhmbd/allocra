FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env.example .env

EXPOSE 8080

CMD ["./main"]
