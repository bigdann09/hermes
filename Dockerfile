FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go && \
    go build -o notification-worker cmd/workers/notifications/main.go && \
    go build -o cli cmd/cli/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/notification-worker .
COPY --from=builder /app/cli .

EXPOSE 3002

CMD ["./main"]